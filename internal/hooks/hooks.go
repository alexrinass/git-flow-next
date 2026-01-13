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

// runHook executes a hook script and returns the result.
func runHook(gitDir string, phase HookPhase, branchType string, action HookAction, ctx HookContext) HookResult {
	hookName := fmt.Sprintf("%s-flow-%s-%s", phase, branchType, action)
	hookPath := filepath.Join(gitDir, "hooks", hookName)

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
