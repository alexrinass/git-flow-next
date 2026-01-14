package hooks_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gittower/git-flow-next/internal/hooks"
	"github.com/gittower/git-flow-next/test/testutil"
)

// createHookScript creates an executable hook script in the hooks directory.
func createHookScript(t *testing.T, dir, name, content string) {
	t.Helper()
	hooksDir := filepath.Join(dir, ".git", "hooks")
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		t.Fatalf("Failed to create hooks directory: %v", err)
	}

	scriptPath := filepath.Join(hooksDir, name)
	if err := os.WriteFile(scriptPath, []byte(content), 0755); err != nil {
		t.Fatalf("Failed to create script %s: %v", name, err)
	}
}

// createNonExecutableHookScript creates a non-executable hook script.
func createNonExecutableHookScript(t *testing.T, dir, name, content string) {
	t.Helper()
	hooksDir := filepath.Join(dir, ".git", "hooks")
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		t.Fatalf("Failed to create hooks directory: %v", err)
	}

	scriptPath := filepath.Join(hooksDir, name)
	if err := os.WriteFile(scriptPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create script %s: %v", name, err)
	}
}

// TestPreHookSuccess tests that pre-hook runs and allows operation to proceed.
func TestPreHookSuccess(t *testing.T) {
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Create a pre-hook script that succeeds
	script := `#!/bin/sh
echo "Pre-hook executed for $BRANCH_TYPE/$BRANCH_NAME"
exit 0
`
	createHookScript(t, dir, "pre-flow-feature-start", script)

	gitDir := filepath.Join(dir, ".git")
	ctx := hooks.HookContext{
		BranchType: "feature",
		BranchName: "my-feature",
		FullBranch: "feature/my-feature",
		BaseBranch: "develop",
		Origin:     "origin",
	}

	err := hooks.RunPreHook(gitDir, "feature", hooks.HookActionStart, ctx)
	if err != nil {
		t.Fatalf("RunPreHook failed unexpectedly: %v", err)
	}
}

// TestPreHookFails tests that failing pre-hook returns error.
func TestPreHookFails(t *testing.T) {
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Create a pre-hook script that fails
	script := `#!/bin/sh
echo "Error: CI is not passing" >&2
exit 1
`
	createHookScript(t, dir, "pre-flow-release-start", script)

	gitDir := filepath.Join(dir, ".git")
	ctx := hooks.HookContext{
		BranchType: "release",
		BranchName: "1.0.0",
		FullBranch: "release/1.0.0",
		BaseBranch: "main",
		Origin:     "origin",
		Version:    "1.0.0",
	}

	err := hooks.RunPreHook(gitDir, "release", hooks.HookActionStart, ctx)
	if err == nil {
		t.Fatal("Expected error for failing pre-hook, got nil")
	}
	if !strings.Contains(err.Error(), "exit code 1") {
		t.Errorf("Expected error to mention exit code, got: %v", err)
	}
}

// TestPreHookNonExistent tests that non-existent pre-hook does not cause error.
func TestPreHookNonExistent(t *testing.T) {
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	gitDir := filepath.Join(dir, ".git")
	ctx := hooks.HookContext{
		BranchType: "feature",
		BranchName: "my-feature",
		FullBranch: "feature/my-feature",
		BaseBranch: "develop",
		Origin:     "origin",
	}

	err := hooks.RunPreHook(gitDir, "feature", hooks.HookActionStart, ctx)
	if err != nil {
		t.Fatalf("Expected no error for non-existent hook, got: %v", err)
	}
}

// TestPreHookNonExecutable tests that non-executable pre-hook is skipped.
func TestPreHookNonExecutable(t *testing.T) {
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Create a non-executable pre-hook script
	script := `#!/bin/sh
exit 1
`
	createNonExecutableHookScript(t, dir, "pre-flow-feature-start", script)

	gitDir := filepath.Join(dir, ".git")
	ctx := hooks.HookContext{
		BranchType: "feature",
		BranchName: "my-feature",
		FullBranch: "feature/my-feature",
		BaseBranch: "develop",
		Origin:     "origin",
	}

	// Should not error because script is not executable
	err := hooks.RunPreHook(gitDir, "feature", hooks.HookActionStart, ctx)
	if err != nil {
		t.Fatalf("Expected no error for non-executable hook, got: %v", err)
	}
}

// TestPostHookExecutes tests that post-hook runs after operation.
func TestPostHookExecutes(t *testing.T) {
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Create a post-hook script that outputs info
	script := `#!/bin/sh
echo "Post-hook: $BRANCH_TYPE/$BRANCH_NAME completed with exit code $EXIT_CODE"
`
	createHookScript(t, dir, "post-flow-feature-finish", script)

	gitDir := filepath.Join(dir, ".git")
	ctx := hooks.HookContext{
		BranchType: "feature",
		BranchName: "my-feature",
		FullBranch: "feature/my-feature",
		BaseBranch: "develop",
		Origin:     "origin",
		ExitCode:   0,
	}

	result := hooks.RunPostHook(gitDir, "feature", hooks.HookActionFinish, ctx)
	if !result.Executed {
		t.Fatal("Expected post-hook to execute")
	}
	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", result.ExitCode)
	}
	if !strings.Contains(result.Output, "exit code 0") {
		t.Errorf("Expected output to contain exit code, got: %s", result.Output)
	}
}

// TestPostHookNonExistent tests that non-existent post-hook returns Executed=false.
func TestPostHookNonExistent(t *testing.T) {
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	gitDir := filepath.Join(dir, ".git")
	ctx := hooks.HookContext{
		BranchType: "feature",
		BranchName: "my-feature",
		FullBranch: "feature/my-feature",
		BaseBranch: "develop",
		Origin:     "origin",
		ExitCode:   0,
	}

	result := hooks.RunPostHook(gitDir, "feature", hooks.HookActionFinish, ctx)
	if result.Executed {
		t.Fatal("Expected Executed=false for non-existent hook")
	}
}

// TestPostHookFailureIgnored tests that post-hook failure is ignored.
func TestPostHookFailureIgnored(t *testing.T) {
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Create a post-hook script that fails
	script := `#!/bin/sh
exit 1
`
	createHookScript(t, dir, "post-flow-feature-finish", script)

	gitDir := filepath.Join(dir, ".git")
	ctx := hooks.HookContext{
		BranchType: "feature",
		BranchName: "my-feature",
		FullBranch: "feature/my-feature",
		BaseBranch: "develop",
		Origin:     "origin",
		ExitCode:   0,
	}

	result := hooks.RunPostHook(gitDir, "feature", hooks.HookActionFinish, ctx)
	if !result.Executed {
		t.Fatal("Expected post-hook to execute")
	}
	if result.ExitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", result.ExitCode)
	}
	// Post-hook failure should not cause an error in result.Error
	if result.Error != nil {
		t.Errorf("Expected no error for post-hook failure, got: %v", result.Error)
	}
}

// TestWithHooksRunsPreAndPost tests the WithHooks wrapper.
func TestWithHooksRunsPreAndPost(t *testing.T) {
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Create marker file to track execution
	markerFile := filepath.Join(dir, "hook-markers.txt")

	// Create pre-hook
	preScript := `#!/bin/sh
echo "pre" >> "` + markerFile + `"
`
	createHookScript(t, dir, "pre-flow-feature-start", preScript)

	// Create post-hook
	postScript := `#!/bin/sh
echo "post-$EXIT_CODE" >> "` + markerFile + `"
`
	createHookScript(t, dir, "post-flow-feature-start", postScript)

	gitDir := filepath.Join(dir, ".git")
	ctx := hooks.HookContext{
		BranchType: "feature",
		BranchName: "my-feature",
		FullBranch: "feature/my-feature",
		BaseBranch: "develop",
		Origin:     "origin",
	}

	operationRan := false
	err := hooks.WithHooks(gitDir, "feature", hooks.HookActionStart, ctx, func() error {
		operationRan = true
		return nil
	})

	if err != nil {
		t.Fatalf("WithHooks failed: %v", err)
	}
	if !operationRan {
		t.Fatal("Operation did not run")
	}

	// Check marker file
	content, err := os.ReadFile(markerFile)
	if err != nil {
		t.Fatalf("Failed to read marker file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) != 2 {
		t.Fatalf("Expected 2 lines in marker file, got %d: %v", len(lines), lines)
	}
	if lines[0] != "pre" {
		t.Errorf("Expected first line to be 'pre', got '%s'", lines[0])
	}
	if lines[1] != "post-0" {
		t.Errorf("Expected second line to be 'post-0', got '%s'", lines[1])
	}
}

// TestWithHooksPreHookBlocks tests that failing pre-hook blocks operation.
func TestWithHooksPreHookBlocks(t *testing.T) {
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Create pre-hook that fails
	preScript := `#!/bin/sh
echo "blocking operation" >&2
exit 1
`
	createHookScript(t, dir, "pre-flow-release-start", preScript)

	gitDir := filepath.Join(dir, ".git")
	ctx := hooks.HookContext{
		BranchType: "release",
		BranchName: "1.0.0",
		FullBranch: "release/1.0.0",
		BaseBranch: "main",
		Origin:     "origin",
		Version:    "1.0.0",
	}

	operationRan := false
	err := hooks.WithHooks(gitDir, "release", hooks.HookActionStart, ctx, func() error {
		operationRan = true
		return nil
	})

	if err == nil {
		t.Fatal("Expected error from failing pre-hook")
	}
	if operationRan {
		t.Fatal("Operation should not have run when pre-hook failed")
	}
}

// TestHookReceivesEnvironmentVariables tests that hooks receive expected env vars.
func TestHookReceivesEnvironmentVariables(t *testing.T) {
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Create a hook that outputs environment variables
	script := `#!/bin/sh
echo "BRANCH=$BRANCH"
echo "BRANCH_NAME=$BRANCH_NAME"
echo "BRANCH_TYPE=$BRANCH_TYPE"
echo "BASE_BRANCH=$BASE_BRANCH"
echo "ORIGIN=$ORIGIN"
echo "VERSION=$VERSION"
`
	createHookScript(t, dir, "pre-flow-release-start", script)

	gitDir := filepath.Join(dir, ".git")
	ctx := hooks.HookContext{
		BranchType: "release",
		BranchName: "2.0.0",
		FullBranch: "release/2.0.0",
		BaseBranch: "main",
		Origin:     "origin",
		Version:    "2.0.0",
	}

	err := hooks.RunPreHook(gitDir, "release", hooks.HookActionStart, ctx)
	if err != nil {
		t.Fatalf("RunPreHook failed: %v", err)
	}
}

// TestHookDifferentActions tests hooks for different actions.
func TestHookDifferentActions(t *testing.T) {
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	actions := []hooks.HookAction{
		hooks.HookActionStart,
		hooks.HookActionFinish,
		hooks.HookActionPublish,
		hooks.HookActionTrack,
		hooks.HookActionDelete,
	}

	for _, action := range actions {
		t.Run(string(action), func(t *testing.T) {
			hookName := "pre-flow-feature-" + string(action)
			script := `#!/bin/sh
echo "action: ` + string(action) + `"
exit 0
`
			createHookScript(t, dir, hookName, script)

			gitDir := filepath.Join(dir, ".git")
			ctx := hooks.HookContext{
				BranchType: "feature",
				BranchName: "test",
				FullBranch: "feature/test",
				BaseBranch: "develop",
				Origin:     "origin",
			}

			err := hooks.RunPreHook(gitDir, "feature", action, ctx)
			if err != nil {
				t.Errorf("RunPreHook failed for action %s: %v", action, err)
			}
		})
	}
}

// TestHooksWorkInGitWorktree tests that hooks execute correctly within a git worktree.
// Git worktrees have a separate git directory structure where the worktree-specific
// git dir is at /main-repo/.git/worktrees/<worktree-name>/.
// Hooks should be found in the shared hooks directory of the main repository.
func TestHooksWorkInGitWorktree(t *testing.T) {
	// Setup main repository
	mainRepo := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, mainRepo)

	// Create a worktree
	worktreePath, err := os.MkdirTemp("", "git-flow-worktree-hooks-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory for worktree: %v", err)
	}
	defer os.RemoveAll(worktreePath)

	// Remove the directory so git worktree add can create it
	os.RemoveAll(worktreePath)

	// Create worktree on a new branch
	_, err = testutil.RunGit(t, mainRepo, "worktree", "add", worktreePath, "-b", "worktree-branch")
	if err != nil {
		t.Fatalf("Failed to create worktree: %v", err)
	}

	// Create a marker file to verify hook execution
	markerFile := filepath.Join(worktreePath, "hook-executed.txt")

	// Create a pre-hook in the main repository's hooks directory (shared location)
	// This is where Git looks for hooks, even when running from a worktree
	script := `#!/bin/sh
echo "Hook executed in worktree" > "` + markerFile + `"
echo "BRANCH_TYPE=$BRANCH_TYPE"
echo "BRANCH_NAME=$BRANCH_NAME"
exit 0
`
	createHookScript(t, mainRepo, "pre-flow-feature-start", script)

	// Get the worktree's git directory by running git rev-parse from the worktree
	// Save current directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Change to worktree directory to get its git dir
	if err := os.Chdir(worktreePath); err != nil {
		t.Fatalf("Failed to change to worktree directory: %v", err)
	}

	// Get git directory from worktree's perspective
	worktreeGitDirOutput, err := testutil.RunGit(t, worktreePath, "rev-parse", "--git-dir")
	if err != nil {
		os.Chdir(oldDir)
		t.Fatalf("Failed to get worktree git directory: %v", err)
	}
	worktreeGitDir := strings.TrimSpace(worktreeGitDirOutput)

	// Restore original directory
	if err := os.Chdir(oldDir); err != nil {
		t.Fatalf("Failed to restore directory: %v", err)
	}

	// The worktree git dir should contain "worktrees" in the path
	if !strings.Contains(worktreeGitDir, "worktrees") {
		t.Errorf("Expected worktree git dir to contain 'worktrees', got: %s", worktreeGitDir)
	}

	// Run the pre-hook using the worktree's git directory
	ctx := hooks.HookContext{
		BranchType: "feature",
		BranchName: "test-worktree-feature",
		FullBranch: "feature/test-worktree-feature",
		BaseBranch: "develop",
		Origin:     "origin",
	}

	err = hooks.RunPreHook(worktreeGitDir, "feature", hooks.HookActionStart, ctx)
	if err != nil {
		t.Fatalf("RunPreHook failed in worktree: %v", err)
	}

	// Verify the hook was executed by checking the marker file
	if _, err := os.Stat(markerFile); os.IsNotExist(err) {
		t.Error("Hook was not executed - marker file not found")
	}
}

// TestPostHookWorksInGitWorktree tests that post-hooks execute correctly in a worktree.
func TestPostHookWorksInGitWorktree(t *testing.T) {
	// Setup main repository
	mainRepo := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, mainRepo)

	// Create a worktree
	worktreePath, err := os.MkdirTemp("", "git-flow-worktree-posthook-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory for worktree: %v", err)
	}
	defer os.RemoveAll(worktreePath)
	os.RemoveAll(worktreePath)

	_, err = testutil.RunGit(t, mainRepo, "worktree", "add", worktreePath, "-b", "posthook-branch")
	if err != nil {
		t.Fatalf("Failed to create worktree: %v", err)
	}

	// Create a post-hook in the main repository's hooks directory
	script := `#!/bin/sh
echo "Post-hook: exit_code=$EXIT_CODE branch=$BRANCH_TYPE/$BRANCH_NAME"
exit 0
`
	createHookScript(t, mainRepo, "post-flow-feature-finish", script)

	// Get the worktree's git directory
	oldDir, _ := os.Getwd()
	os.Chdir(worktreePath)
	worktreeGitDirOutput, err := testutil.RunGit(t, worktreePath, "rev-parse", "--git-dir")
	os.Chdir(oldDir)
	if err != nil {
		t.Fatalf("Failed to get worktree git directory: %v", err)
	}
	worktreeGitDir := strings.TrimSpace(worktreeGitDirOutput)

	ctx := hooks.HookContext{
		BranchType: "feature",
		BranchName: "completed-feature",
		FullBranch: "feature/completed-feature",
		BaseBranch: "develop",
		Origin:     "origin",
		ExitCode:   0,
	}

	result := hooks.RunPostHook(worktreeGitDir, "feature", hooks.HookActionFinish, ctx)
	if !result.Executed {
		t.Error("Expected post-hook to execute in worktree")
	}
	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", result.ExitCode)
	}
	if !strings.Contains(result.Output, "exit_code=0") {
		t.Errorf("Expected output to contain exit_code=0, got: %s", result.Output)
	}
}

// TestWithHooksInWorktree tests the WithHooks wrapper in a worktree context.
func TestWithHooksInWorktree(t *testing.T) {
	// Setup main repository
	mainRepo := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, mainRepo)

	// Create a worktree
	worktreePath, err := os.MkdirTemp("", "git-flow-worktree-withhooks-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory for worktree: %v", err)
	}
	defer os.RemoveAll(worktreePath)
	os.RemoveAll(worktreePath)

	_, err = testutil.RunGit(t, mainRepo, "worktree", "add", worktreePath, "-b", "withhooks-branch")
	if err != nil {
		t.Fatalf("Failed to create worktree: %v", err)
	}

	// Create marker file to track execution order
	markerFile := filepath.Join(worktreePath, "withhooks-markers.txt")

	// Create pre-hook
	preScript := `#!/bin/sh
echo "pre" >> "` + markerFile + `"
`
	createHookScript(t, mainRepo, "pre-flow-release-start", preScript)

	// Create post-hook
	postScript := `#!/bin/sh
echo "post-$EXIT_CODE" >> "` + markerFile + `"
`
	createHookScript(t, mainRepo, "post-flow-release-start", postScript)

	// Get worktree git directory
	oldDir, _ := os.Getwd()
	os.Chdir(worktreePath)
	worktreeGitDirOutput, err := testutil.RunGit(t, worktreePath, "rev-parse", "--git-dir")
	os.Chdir(oldDir)
	if err != nil {
		t.Fatalf("Failed to get worktree git directory: %v", err)
	}
	worktreeGitDir := strings.TrimSpace(worktreeGitDirOutput)

	ctx := hooks.HookContext{
		BranchType: "release",
		BranchName: "1.0.0",
		FullBranch: "release/1.0.0",
		BaseBranch: "main",
		Origin:     "origin",
		Version:    "1.0.0",
	}

	operationRan := false
	err = hooks.WithHooks(worktreeGitDir, "release", hooks.HookActionStart, ctx, func() error {
		operationRan = true
		return nil
	})

	if err != nil {
		t.Fatalf("WithHooks failed in worktree: %v", err)
	}
	if !operationRan {
		t.Fatal("Operation did not run in worktree")
	}

	// Verify execution order from marker file
	content, err := os.ReadFile(markerFile)
	if err != nil {
		t.Fatalf("Failed to read marker file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) != 2 {
		t.Fatalf("Expected 2 lines in marker file, got %d: %v", len(lines), lines)
	}
	if lines[0] != "pre" {
		t.Errorf("Expected first line to be 'pre', got '%s'", lines[0])
	}
	if lines[1] != "post-0" {
		t.Errorf("Expected second line to be 'post-0', got '%s'", lines[1])
	}
}
