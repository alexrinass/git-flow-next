package hooks

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// RunPreHook executes a pre-hook script. Returns an error if the hook fails (non-zero exit).
// If the hook does not exist or is not executable, it returns nil (no error).
func RunPreHook(gitDir string, branchType string, action HookAction, ctx HookContext) error {
	result := runHook(gitDir, HookPre, branchType, action, ctx)
	if result.Error != nil {
		return result.Error
	}
	if result.Executed && result.ExitCode != 0 {
		if result.Output != "" {
			return fmt.Errorf("pre-hook '%s-flow-%s-%s' failed with exit code %d:\n%s",
				HookPre, branchType, action, result.ExitCode, result.Output)
		}
		return fmt.Errorf("pre-hook '%s-flow-%s-%s' failed with exit code %d",
			HookPre, branchType, action, result.ExitCode)
	}
	return nil
}

// RunPostHook executes a post-hook script. The result is returned but errors do not
// cause the operation to fail. If the hook does not exist or is not executable,
// it returns a result with Executed=false.
func RunPostHook(gitDir string, branchType string, action HookAction, ctx HookContext) HookResult {
	return runHook(gitDir, HookPost, branchType, action, ctx)
}

// getHooksDir returns the directory where hooks are stored.
// For regular repositories, this is gitDir/hooks.
// For worktrees, hooks are shared in the main repository's git directory.
// Worktree git directories have a path like: /repo/.git/worktrees/<name>
// In this case, hooks should be found at /repo/.git/hooks (the common git dir).
func getHooksDir(gitDir string) string {
	// Check if this is a worktree git directory
	// Worktree git dirs have a path pattern: .../worktrees/<worktree-name>
	// In this case, we need to go up two levels to get to the common git dir
	commonDir := getCommonGitDir(gitDir)
	return filepath.Join(commonDir, "hooks")
}

// getCommonGitDir returns the common git directory from a git directory path.
// For regular repositories, this returns the gitDir unchanged.
// For worktrees (where gitDir is like /repo/.git/worktrees/<name>), this returns /repo/.git
func getCommonGitDir(gitDir string) string {
	// Normalize the path
	gitDir = filepath.Clean(gitDir)

	// Check if this looks like a worktree git dir
	// Pattern: .../worktrees/<worktree-name>
	parent := filepath.Dir(gitDir)
	parentBase := filepath.Base(parent)

	if parentBase == "worktrees" {
		// This is a worktree git directory
		// Go up two levels: worktrees/<name> -> .git
		return filepath.Dir(parent)
	}

	// Not a worktree, return as-is
	return gitDir
}

// runHook executes a hook script and returns the result.
func runHook(gitDir string, phase HookPhase, branchType string, action HookAction, ctx HookContext) HookResult {
	hookName := fmt.Sprintf("%s-flow-%s-%s", phase, branchType, action)
	hooksDir := getHooksDir(gitDir)
	hookPath := filepath.Join(hooksDir, hookName)

	// Check if hook exists
	info, err := os.Stat(hookPath)
	if os.IsNotExist(err) {
		return HookResult{Executed: false}
	}
	if err != nil {
		return HookResult{Executed: false, Error: err}
	}

	// Check if executable
	if info.Mode()&0111 == 0 {
		// Not executable, skip silently
		return HookResult{Executed: false}
	}

	// Build environment variables
	env := buildHookEnv(ctx, phase)

	// Execute hook
	cmd := exec.Command(hookPath)
	cmd.Env = env
	cmd.Dir = filepath.Dir(gitDir) // Repository root

	output, err := cmd.CombinedOutput()

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return HookResult{
				Executed: true,
				ExitCode: 1,
				Output:   string(output),
				Error:    err,
			}
		}
	}

	return HookResult{
		Executed: true,
		ExitCode: exitCode,
		Output:   string(output),
		Error:    nil,
	}
}

// buildHookEnv builds environment variables for hook execution.
func buildHookEnv(ctx HookContext, phase HookPhase) []string {
	env := os.Environ()
	env = append(env,
		fmt.Sprintf("BRANCH=%s", ctx.FullBranch),
		fmt.Sprintf("BRANCH_NAME=%s", ctx.BranchName),
		fmt.Sprintf("BRANCH_TYPE=%s", ctx.BranchType),
		fmt.Sprintf("BASE_BRANCH=%s", ctx.BaseBranch),
		fmt.Sprintf("ORIGIN=%s", ctx.Origin),
	)

	if ctx.Version != "" {
		env = append(env, fmt.Sprintf("VERSION=%s", ctx.Version))
	}

	// For post-hooks, include the exit code of the operation
	if phase == HookPost {
		env = append(env, fmt.Sprintf("EXIT_CODE=%d", ctx.ExitCode))
	}

	return env
}

// WithHooks wraps an operation with pre and post hooks.
// The pre-hook is run before the operation. If it fails, the operation is not executed.
// The post-hook is run after the operation, regardless of success or failure.
// The context's ExitCode is set based on the operation result before running the post-hook.
func WithHooks(gitDir string, branchType string, action HookAction, ctx HookContext, operation func() error) error {
	// Run pre-hook
	if err := RunPreHook(gitDir, branchType, action, ctx); err != nil {
		return err
	}

	// Run operation
	opErr := operation()

	// Set exit code for post-hook
	if opErr != nil {
		ctx.ExitCode = 1
	} else {
		ctx.ExitCode = 0
	}

	// Run post-hook (ignore errors from post-hook)
	result := RunPostHook(gitDir, branchType, action, ctx)
	if result.Executed && result.Output != "" {
		// Print post-hook output for visibility
		fmt.Print(result.Output)
	}

	return opErr
}
