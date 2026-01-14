package cmd_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gittower/git-flow-next/test/testutil"
)

// createHookScript creates an executable hook script in the repository's hooks directory.
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

// =============================================================================
// Pre-Hook Tests - Verify hooks can block operations
// =============================================================================

// TestStartPreHookBlocks tests that a failing pre-hook prevents branch creation.
func TestStartPreHookBlocks(t *testing.T) {
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Initialize git-flow
	_, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
	if err != nil {
		t.Fatalf("Failed to initialize git-flow: %v", err)
	}

	// Create a pre-hook that fails
	script := `#!/bin/sh
echo "Pre-hook blocking operation" >&2
exit 1
`
	createHookScript(t, dir, "pre-flow-feature-start", script)

	// Try to start a feature - should fail
	output, err := testutil.RunGitFlow(t, dir, "feature", "start", "blocked-feature")
	if err == nil {
		t.Fatal("Expected feature start to fail due to pre-hook, but it succeeded")
	}

	// Verify the error mentions the hook
	if !strings.Contains(output, "pre-hook") {
		t.Errorf("Expected error to mention pre-hook, got: %s", output)
	}

	// Verify branch was NOT created
	if testutil.BranchExists(t, dir, "feature/blocked-feature") {
		t.Error("Branch should not have been created when pre-hook failed")
	}
}

// TestFinishPreHookBlocks tests that a failing pre-hook prevents finish operation.
func TestFinishPreHookBlocks(t *testing.T) {
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Initialize git-flow and create a feature branch
	_, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
	if err != nil {
		t.Fatalf("Failed to initialize git-flow: %v", err)
	}

	_, err = testutil.RunGitFlow(t, dir, "feature", "start", "my-feature")
	if err != nil {
		t.Fatalf("Failed to start feature: %v", err)
	}

	// Add a commit to the feature branch
	testutil.WriteFile(t, dir, "feature.txt", "feature content")
	_, _ = testutil.RunGit(t, dir, "add", "feature.txt")
	_, _ = testutil.RunGit(t, dir, "commit", "-m", "Add feature file")

	// Create a pre-hook that fails
	script := `#!/bin/sh
echo "Pre-hook blocking finish" >&2
exit 1
`
	createHookScript(t, dir, "pre-flow-feature-finish", script)

	// Try to finish - should fail
	output, err := testutil.RunGitFlow(t, dir, "feature", "finish", "my-feature")
	if err == nil {
		t.Fatal("Expected feature finish to fail due to pre-hook, but it succeeded")
	}

	// Verify the error mentions the hook
	if !strings.Contains(output, "pre-hook") {
		t.Errorf("Expected error to mention pre-hook, got: %s", output)
	}

	// Verify branch still exists (finish was blocked)
	if !testutil.BranchExists(t, dir, "feature/my-feature") {
		t.Error("Feature branch should still exist when pre-hook blocked finish")
	}
}

// TestDeletePreHookBlocks tests that a failing pre-hook prevents delete operation.
func TestDeletePreHookBlocks(t *testing.T) {
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Initialize git-flow and create a feature branch
	_, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
	if err != nil {
		t.Fatalf("Failed to initialize git-flow: %v", err)
	}

	_, err = testutil.RunGitFlow(t, dir, "feature", "start", "to-delete")
	if err != nil {
		t.Fatalf("Failed to start feature: %v", err)
	}

	// Switch back to develop so we can delete the feature branch
	_, _ = testutil.RunGit(t, dir, "checkout", "develop")

	// Create a pre-hook that fails
	script := `#!/bin/sh
echo "Pre-hook blocking delete" >&2
exit 1
`
	createHookScript(t, dir, "pre-flow-feature-delete", script)

	// Try to delete - should fail
	output, err := testutil.RunGitFlow(t, dir, "feature", "delete", "to-delete")
	if err == nil {
		t.Fatal("Expected feature delete to fail due to pre-hook, but it succeeded")
	}

	// Verify the error mentions the hook
	if !strings.Contains(output, "pre-hook") {
		t.Errorf("Expected error to mention pre-hook, got: %s", output)
	}

	// Verify branch still exists
	if !testutil.BranchExists(t, dir, "feature/to-delete") {
		t.Error("Feature branch should still exist when pre-hook blocked delete")
	}
}

// =============================================================================
// Post-Hook Tests - Verify hooks run after operations
// =============================================================================

// TestStartPostHookRuns tests that post-hook executes after successful branch creation.
func TestStartPostHookRuns(t *testing.T) {
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Initialize git-flow
	_, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
	if err != nil {
		t.Fatalf("Failed to initialize git-flow: %v", err)
	}

	// Create a marker file path
	markerFile := filepath.Join(dir, "post-hook-executed.txt")

	// Create a post-hook that creates a marker file
	script := `#!/bin/sh
echo "BRANCH=$BRANCH" > "` + markerFile + `"
echo "BRANCH_NAME=$BRANCH_NAME" >> "` + markerFile + `"
echo "BRANCH_TYPE=$BRANCH_TYPE" >> "` + markerFile + `"
`
	createHookScript(t, dir, "post-flow-feature-start", script)

	// Start a feature
	_, err = testutil.RunGitFlow(t, dir, "feature", "start", "post-hook-test")
	if err != nil {
		t.Fatalf("Failed to start feature: %v", err)
	}

	// Verify post-hook ran by checking marker file
	content, err := os.ReadFile(markerFile)
	if err != nil {
		t.Fatalf("Post-hook did not run - marker file not found: %v", err)
	}

	// Verify environment variables were passed correctly
	contentStr := string(content)
	if !strings.Contains(contentStr, "BRANCH=feature/post-hook-test") {
		t.Errorf("Expected BRANCH=feature/post-hook-test in hook output, got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "BRANCH_NAME=post-hook-test") {
		t.Errorf("Expected BRANCH_NAME=post-hook-test in hook output, got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "BRANCH_TYPE=feature") {
		t.Errorf("Expected BRANCH_TYPE=feature in hook output, got: %s", contentStr)
	}
}

// TestFinishPostHookReceivesExitCode tests that post-hook receives correct exit code.
func TestFinishPostHookReceivesExitCode(t *testing.T) {
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Initialize git-flow and create a feature branch
	_, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
	if err != nil {
		t.Fatalf("Failed to initialize git-flow: %v", err)
	}

	_, err = testutil.RunGitFlow(t, dir, "feature", "start", "exit-code-test")
	if err != nil {
		t.Fatalf("Failed to start feature: %v", err)
	}

	// Add a commit
	testutil.WriteFile(t, dir, "test.txt", "test content")
	_, _ = testutil.RunGit(t, dir, "add", "test.txt")
	_, _ = testutil.RunGit(t, dir, "commit", "-m", "Add test file")

	// Create a marker file path
	markerFile := filepath.Join(dir, "post-finish-executed.txt")

	// Create a post-hook that records the exit code
	script := `#!/bin/sh
echo "EXIT_CODE=$EXIT_CODE" > "` + markerFile + `"
`
	createHookScript(t, dir, "post-flow-feature-finish", script)

	// Finish the feature
	_, err = testutil.RunGitFlow(t, dir, "feature", "finish", "exit-code-test")
	if err != nil {
		t.Fatalf("Failed to finish feature: %v", err)
	}

	// Verify post-hook ran and received exit code 0
	content, err := os.ReadFile(markerFile)
	if err != nil {
		t.Fatalf("Post-hook did not run - marker file not found: %v", err)
	}

	if !strings.Contains(string(content), "EXIT_CODE=0") {
		t.Errorf("Expected EXIT_CODE=0 in hook output, got: %s", string(content))
	}
}

// =============================================================================
// Version Filter Tests - Verify filters modify branch names
// =============================================================================

// TestStartVersionFilterModifiesBranchName tests that version filter modifies the branch name.
func TestStartVersionFilterModifiesBranchName(t *testing.T) {
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Initialize git-flow
	_, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
	if err != nil {
		t.Fatalf("Failed to initialize git-flow: %v", err)
	}

	// Create a version filter that adds 'v' prefix
	script := `#!/bin/sh
VERSION="$1"
if [ "${VERSION#v}" = "$VERSION" ]; then
    echo "v$VERSION"
else
    echo "$VERSION"
fi
`
	createHookScript(t, dir, "filter-flow-release-start-version", script)

	// Start a release with version "1.0.0" - filter should change to "v1.0.0"
	output, err := testutil.RunGitFlow(t, dir, "release", "start", "1.0.0")
	if err != nil {
		t.Fatalf("Failed to start release: %v\nOutput: %s", err, output)
	}

	// Verify the filter message is shown
	if !strings.Contains(output, "Version filter changed") {
		t.Errorf("Expected output to mention version filter, got: %s", output)
	}

	// Verify the branch was created with the filtered name
	if !testutil.BranchExists(t, dir, "release/v1.0.0") {
		t.Error("Expected release/v1.0.0 branch to exist (filtered from 1.0.0)")
	}

	// Verify original name branch was NOT created
	if testutil.BranchExists(t, dir, "release/1.0.0") {
		t.Error("release/1.0.0 should not exist - filter should have changed it to v1.0.0")
	}
}

// TestVersionFilterPassedToHooks tests that filtered version is passed to hooks.
func TestVersionFilterPassedToHooks(t *testing.T) {
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Initialize git-flow
	_, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
	if err != nil {
		t.Fatalf("Failed to initialize git-flow: %v", err)
	}

	markerFile := filepath.Join(dir, "hook-received-version.txt")

	// Create a version filter that adds 'v' prefix
	filterScript := `#!/bin/sh
VERSION="$1"
echo "v$VERSION"
`
	createHookScript(t, dir, "filter-flow-release-start-version", filterScript)

	// Create a post-hook that records the version it received
	hookScript := `#!/bin/sh
echo "VERSION=$VERSION" > "` + markerFile + `"
echo "BRANCH_NAME=$BRANCH_NAME" >> "` + markerFile + `"
`
	createHookScript(t, dir, "post-flow-release-start", hookScript)

	// Start a release
	_, err = testutil.RunGitFlow(t, dir, "release", "start", "2.0.0")
	if err != nil {
		t.Fatalf("Failed to start release: %v", err)
	}

	// Verify hook received filtered version
	content, err := os.ReadFile(markerFile)
	if err != nil {
		t.Fatalf("Post-hook did not run: %v", err)
	}

	contentStr := string(content)
	// The hook should receive the filtered version
	if !strings.Contains(contentStr, "VERSION=v2.0.0") {
		t.Errorf("Expected VERSION=v2.0.0, got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "BRANCH_NAME=v2.0.0") {
		t.Errorf("Expected BRANCH_NAME=v2.0.0, got: %s", contentStr)
	}
}

// =============================================================================
// Tag Message Filter Tests - Verify filters modify tag messages
// =============================================================================

// TestFinishTagMessageFilterModifiesTag tests that tag message filter modifies the tag message.
func TestFinishTagMessageFilterModifiesTag(t *testing.T) {
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Initialize git-flow
	_, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
	if err != nil {
		t.Fatalf("Failed to initialize git-flow: %v", err)
	}

	// Start a release
	_, err = testutil.RunGitFlow(t, dir, "release", "start", "3.0.0")
	if err != nil {
		t.Fatalf("Failed to start release: %v", err)
	}

	// Add a commit to the release branch
	testutil.WriteFile(t, dir, "release.txt", "release content")
	_, _ = testutil.RunGit(t, dir, "add", "release.txt")
	_, _ = testutil.RunGit(t, dir, "commit", "-m", "Add release file")

	// Create a tag message filter that adds custom prefix
	script := `#!/bin/sh
VERSION="$1"
echo "Custom Release: $VERSION - Modified by filter"
`
	createHookScript(t, dir, "filter-flow-release-finish-tag-message", script)

	// Finish the release
	_, err = testutil.RunGitFlow(t, dir, "release", "finish", "3.0.0")
	if err != nil {
		t.Fatalf("Failed to finish release: %v", err)
	}

	// Verify tag was created
	output, err := testutil.RunGit(t, dir, "tag", "-l")
	if err != nil {
		t.Fatalf("Failed to list tags: %v", err)
	}
	if !strings.Contains(output, "3.0.0") {
		t.Errorf("Expected tag 3.0.0 to exist, got: %s", output)
	}

	// Verify tag message was modified by filter
	output, err = testutil.RunGit(t, dir, "tag", "-l", "-n1", "3.0.0")
	if err != nil {
		t.Fatalf("Failed to get tag message: %v", err)
	}
	if !strings.Contains(output, "Custom Release") || !strings.Contains(output, "Modified by filter") {
		t.Errorf("Expected tag message to be modified by filter, got: %s", output)
	}
}

// =============================================================================
// Worktree Integration Test - Verify hooks work from worktree
// =============================================================================

// TestStartWithHooksInWorktree tests that hooks work when running commands from a worktree.
func TestStartWithHooksInWorktree(t *testing.T) {
	// Setup main repository
	mainRepo := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, mainRepo)

	// Initialize git-flow in main repo
	_, err := testutil.RunGitFlow(t, mainRepo, "init", "--defaults")
	if err != nil {
		t.Fatalf("Failed to initialize git-flow: %v", err)
	}

	// Create a worktree
	worktreePath, err := os.MkdirTemp("", "git-flow-cmd-worktree-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory for worktree: %v", err)
	}
	defer os.RemoveAll(worktreePath)
	os.RemoveAll(worktreePath)

	_, err = testutil.RunGit(t, mainRepo, "worktree", "add", worktreePath, "-b", "worktree-branch")
	if err != nil {
		t.Fatalf("Failed to create worktree: %v", err)
	}

	// Create marker file in worktree directory
	markerFile := filepath.Join(worktreePath, "worktree-hook-executed.txt")

	// Create hooks in main repo's .git/hooks (shared location)
	preScript := `#!/bin/sh
echo "pre-hook-ran" > "` + markerFile + `"
`
	createHookScript(t, mainRepo, "pre-flow-feature-start", preScript)

	postScript := `#!/bin/sh
echo "post-hook-ran" >> "` + markerFile + `"
echo "BRANCH=$BRANCH" >> "` + markerFile + `"
`
	createHookScript(t, mainRepo, "post-flow-feature-start", postScript)

	// Run git-flow command from worktree
	output, err := testutil.RunGitFlow(t, worktreePath, "feature", "start", "worktree-feature")
	if err != nil {
		t.Fatalf("Failed to start feature from worktree: %v\nOutput: %s", err, output)
	}

	// Verify hooks ran (marker file was created)
	content, err := os.ReadFile(markerFile)
	if err != nil {
		t.Fatalf("Hooks did not run from worktree - marker file not found: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "pre-hook-ran") {
		t.Error("Pre-hook did not run from worktree")
	}
	if !strings.Contains(contentStr, "post-hook-ran") {
		t.Error("Post-hook did not run from worktree")
	}
	if !strings.Contains(contentStr, "BRANCH=feature/worktree-feature") {
		t.Errorf("Hook did not receive correct branch, got: %s", contentStr)
	}

	// Verify branch was created
	if !testutil.BranchExists(t, worktreePath, "feature/worktree-feature") {
		t.Error("Feature branch should have been created from worktree")
	}
}

// =============================================================================
// Publish Hook Tests - Verify hooks run for publish operations
// =============================================================================

// TestPublishPreHookBlocks tests that a failing pre-hook prevents publish operation.
func TestPublishPreHookBlocks(t *testing.T) {
	dir, remoteDir := testutil.SetupTestRepoWithRemote(t)
	defer testutil.CleanupTestRepo(t, dir)
	defer testutil.CleanupTestRepo(t, remoteDir)

	// Start a feature branch
	_, err := testutil.RunGitFlow(t, dir, "feature", "start", "publish-test")
	if err != nil {
		t.Fatalf("Failed to start feature: %v", err)
	}

	// Create a pre-hook that fails
	script := `#!/bin/sh
echo "Pre-hook blocking publish" >&2
exit 1
`
	createHookScript(t, dir, "pre-flow-feature-publish", script)

	// Try to publish - should fail
	output, err := testutil.RunGitFlow(t, dir, "feature", "publish", "publish-test")
	if err == nil {
		t.Fatal("Expected feature publish to fail due to pre-hook, but it succeeded")
	}

	// Verify the error mentions the hook
	if !strings.Contains(output, "pre-hook") {
		t.Errorf("Expected error to mention pre-hook, got: %s", output)
	}

	// Verify branch was NOT pushed to remote
	if testutil.RemoteBranchExists(t, dir, "origin", "feature/publish-test") {
		t.Error("Branch should not have been pushed when pre-hook failed")
	}
}

// TestPublishPostHookRuns tests that post-hook executes after successful publish.
func TestPublishPostHookRuns(t *testing.T) {
	dir, remoteDir := testutil.SetupTestRepoWithRemote(t)
	defer testutil.CleanupTestRepo(t, dir)
	defer testutil.CleanupTestRepo(t, remoteDir)

	// Start a feature branch
	_, err := testutil.RunGitFlow(t, dir, "feature", "start", "publish-hook-test")
	if err != nil {
		t.Fatalf("Failed to start feature: %v", err)
	}

	// Create a marker file path
	markerFile := filepath.Join(dir, "publish-hook-executed.txt")

	// Create a post-hook that creates a marker file
	script := `#!/bin/sh
echo "BRANCH=$BRANCH" > "` + markerFile + `"
echo "BRANCH_NAME=$BRANCH_NAME" >> "` + markerFile + `"
echo "ORIGIN=$ORIGIN" >> "` + markerFile + `"
`
	createHookScript(t, dir, "post-flow-feature-publish", script)

	// Publish the feature
	_, err = testutil.RunGitFlow(t, dir, "feature", "publish", "publish-hook-test")
	if err != nil {
		t.Fatalf("Failed to publish feature: %v", err)
	}

	// Verify post-hook ran by checking marker file
	content, err := os.ReadFile(markerFile)
	if err != nil {
		t.Fatalf("Post-hook did not run - marker file not found: %v", err)
	}

	// Verify environment variables were passed correctly
	contentStr := string(content)
	if !strings.Contains(contentStr, "BRANCH=feature/publish-hook-test") {
		t.Errorf("Expected BRANCH=feature/publish-hook-test in hook output, got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "BRANCH_NAME=publish-hook-test") {
		t.Errorf("Expected BRANCH_NAME=publish-hook-test in hook output, got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "ORIGIN=origin") {
		t.Errorf("Expected ORIGIN=origin in hook output, got: %s", contentStr)
	}
}

// =============================================================================
// Track Hook Tests - Verify hooks run for track operations
// =============================================================================

// TestTrackPreHookBlocks tests that a failing pre-hook prevents track operation.
func TestTrackPreHookBlocks(t *testing.T) {
	dir, remoteDir := testutil.SetupTestRepoWithRemote(t)
	defer testutil.CleanupTestRepo(t, dir)
	defer testutil.CleanupTestRepo(t, remoteDir)

	// Start and publish a feature branch to create it on remote
	_, err := testutil.RunGitFlow(t, dir, "feature", "start", "track-test")
	if err != nil {
		t.Fatalf("Failed to start feature: %v", err)
	}
	_, err = testutil.RunGitFlow(t, dir, "feature", "publish", "track-test")
	if err != nil {
		t.Fatalf("Failed to publish feature: %v", err)
	}

	// Delete local branch and switch to develop
	_, _ = testutil.RunGit(t, dir, "checkout", "develop")
	_, _ = testutil.RunGit(t, dir, "branch", "-D", "feature/track-test")

	// Create a pre-hook that fails
	script := `#!/bin/sh
echo "Pre-hook blocking track" >&2
exit 1
`
	createHookScript(t, dir, "pre-flow-feature-track", script)

	// Try to track - should fail
	output, err := testutil.RunGitFlow(t, dir, "feature", "track", "track-test")
	if err == nil {
		t.Fatal("Expected feature track to fail due to pre-hook, but it succeeded")
	}

	// Verify the error mentions the hook
	if !strings.Contains(output, "pre-hook") {
		t.Errorf("Expected error to mention pre-hook, got: %s", output)
	}

	// Verify local branch was NOT created
	if testutil.BranchExists(t, dir, "feature/track-test") {
		t.Error("Local branch should not have been created when pre-hook failed")
	}
}

// TestTrackPostHookRuns tests that post-hook executes after successful track.
func TestTrackPostHookRuns(t *testing.T) {
	dir, remoteDir := testutil.SetupTestRepoWithRemote(t)
	defer testutil.CleanupTestRepo(t, dir)
	defer testutil.CleanupTestRepo(t, remoteDir)

	// Start and publish a feature branch to create it on remote
	_, err := testutil.RunGitFlow(t, dir, "feature", "start", "track-hook-test")
	if err != nil {
		t.Fatalf("Failed to start feature: %v", err)
	}
	_, err = testutil.RunGitFlow(t, dir, "feature", "publish", "track-hook-test")
	if err != nil {
		t.Fatalf("Failed to publish feature: %v", err)
	}

	// Delete local branch and switch to develop
	_, _ = testutil.RunGit(t, dir, "checkout", "develop")
	_, _ = testutil.RunGit(t, dir, "branch", "-D", "feature/track-hook-test")

	// Create a marker file path
	markerFile := filepath.Join(dir, "track-hook-executed.txt")

	// Create a post-hook that creates a marker file
	script := `#!/bin/sh
echo "BRANCH=$BRANCH" > "` + markerFile + `"
echo "BRANCH_NAME=$BRANCH_NAME" >> "` + markerFile + `"
echo "ORIGIN=$ORIGIN" >> "` + markerFile + `"
`
	createHookScript(t, dir, "post-flow-feature-track", script)

	// Track the feature
	_, err = testutil.RunGitFlow(t, dir, "feature", "track", "track-hook-test")
	if err != nil {
		t.Fatalf("Failed to track feature: %v", err)
	}

	// Verify post-hook ran by checking marker file
	content, err := os.ReadFile(markerFile)
	if err != nil {
		t.Fatalf("Post-hook did not run - marker file not found: %v", err)
	}

	// Verify environment variables were passed correctly
	contentStr := string(content)
	if !strings.Contains(contentStr, "BRANCH=feature/track-hook-test") {
		t.Errorf("Expected BRANCH=feature/track-hook-test in hook output, got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "BRANCH_NAME=track-hook-test") {
		t.Errorf("Expected BRANCH_NAME=track-hook-test in hook output, got: %s", contentStr)
	}
}

// =============================================================================
// Delete Post-Hook Test - Verify post-hook runs after delete
// =============================================================================

// TestDeletePostHookRuns tests that post-hook executes after successful delete.
func TestDeletePostHookRuns(t *testing.T) {
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Initialize git-flow and create a feature branch
	_, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
	if err != nil {
		t.Fatalf("Failed to initialize git-flow: %v", err)
	}

	_, err = testutil.RunGitFlow(t, dir, "feature", "start", "to-delete-hook")
	if err != nil {
		t.Fatalf("Failed to start feature: %v", err)
	}

	// Switch back to develop so we can delete the feature branch
	_, _ = testutil.RunGit(t, dir, "checkout", "develop")

	// Create a marker file path
	markerFile := filepath.Join(dir, "delete-hook-executed.txt")

	// Create a post-hook that creates a marker file
	script := `#!/bin/sh
echo "BRANCH=$BRANCH" > "` + markerFile + `"
echo "BRANCH_NAME=$BRANCH_NAME" >> "` + markerFile + `"
echo "EXIT_CODE=$EXIT_CODE" >> "` + markerFile + `"
`
	createHookScript(t, dir, "post-flow-feature-delete", script)

	// Delete the feature
	_, err = testutil.RunGitFlow(t, dir, "feature", "delete", "to-delete-hook")
	if err != nil {
		t.Fatalf("Failed to delete feature: %v", err)
	}

	// Verify post-hook ran by checking marker file
	content, err := os.ReadFile(markerFile)
	if err != nil {
		t.Fatalf("Post-hook did not run - marker file not found: %v", err)
	}

	// Verify environment variables were passed correctly
	contentStr := string(content)
	if !strings.Contains(contentStr, "BRANCH=feature/to-delete-hook") {
		t.Errorf("Expected BRANCH=feature/to-delete-hook in hook output, got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "EXIT_CODE=0") {
		t.Errorf("Expected EXIT_CODE=0 in hook output, got: %s", contentStr)
	}
}

// =============================================================================
// Custom Branch Type Hook Tests - Verify hooks work with custom branch configs
// =============================================================================

// TestCustomBranchTypePreHookBlocks tests that hooks work with custom branch types.
func TestCustomBranchTypePreHookBlocks(t *testing.T) {
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Initialize git-flow
	_, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
	if err != nil {
		t.Fatalf("Failed to initialize git-flow: %v", err)
	}

	// Configure a custom branch type "bugfix"
	_, _ = testutil.RunGit(t, dir, "config", "gitflow.branch.bugfix.prefix", "bugfix/")
	_, _ = testutil.RunGit(t, dir, "config", "gitflow.branch.bugfix.parent", "develop")

	// Create a pre-hook that fails for the custom branch type
	script := `#!/bin/sh
echo "Pre-hook blocking bugfix start" >&2
exit 1
`
	createHookScript(t, dir, "pre-flow-bugfix-start", script)

	// Try to start a bugfix - should fail
	output, err := testutil.RunGitFlow(t, dir, "bugfix", "start", "custom-blocked")
	if err == nil {
		t.Fatal("Expected bugfix start to fail due to pre-hook, but it succeeded")
	}

	// Verify the error mentions the hook
	if !strings.Contains(output, "pre-hook") {
		t.Errorf("Expected error to mention pre-hook, got: %s", output)
	}

	// Verify branch was NOT created
	if testutil.BranchExists(t, dir, "bugfix/custom-blocked") {
		t.Error("Branch should not have been created when pre-hook failed")
	}
}

// TestCustomBranchTypePostHookReceivesCorrectType tests that post-hooks receive the correct branch type.
// Uses "support" branch type which is a standard git-flow type but configured with custom prefix.
func TestCustomBranchTypePostHookReceivesCorrectType(t *testing.T) {
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Initialize git-flow
	_, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
	if err != nil {
		t.Fatalf("Failed to initialize git-flow: %v", err)
	}

	// Configure support branch with a custom prefix
	_, _ = testutil.RunGit(t, dir, "config", "gitflow.branch.support.prefix", "sup/")
	_, _ = testutil.RunGit(t, dir, "config", "gitflow.branch.support.parent", "main")

	// Create a marker file path
	markerFile := filepath.Join(dir, "custom-hook-executed.txt")

	// Create a post-hook that records the branch type
	script := `#!/bin/sh
echo "BRANCH=$BRANCH" > "` + markerFile + `"
echo "BRANCH_NAME=$BRANCH_NAME" >> "` + markerFile + `"
echo "BRANCH_TYPE=$BRANCH_TYPE" >> "` + markerFile + `"
echo "BASE_BRANCH=$BASE_BRANCH" >> "` + markerFile + `"
`
	createHookScript(t, dir, "post-flow-support-start", script)

	// Start a support branch
	_, err = testutil.RunGitFlow(t, dir, "support", "start", "lts-1.0")
	if err != nil {
		t.Fatalf("Failed to start support: %v", err)
	}

	// Verify post-hook ran by checking marker file
	content, err := os.ReadFile(markerFile)
	if err != nil {
		t.Fatalf("Post-hook did not run - marker file not found: %v", err)
	}

	// Verify environment variables were passed correctly for support type
	contentStr := string(content)
	if !strings.Contains(contentStr, "BRANCH=sup/lts-1.0") {
		t.Errorf("Expected BRANCH=sup/lts-1.0 in hook output, got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "BRANCH_NAME=lts-1.0") {
		t.Errorf("Expected BRANCH_NAME=lts-1.0 in hook output, got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "BRANCH_TYPE=support") {
		t.Errorf("Expected BRANCH_TYPE=support in hook output, got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "BASE_BRANCH=main") {
		t.Errorf("Expected BASE_BRANCH=main in hook output, got: %s", contentStr)
	}
}

// TestHotfixVersionFilter tests that version filters work with hotfix branch type.
// This verifies filters work with branch types beyond just release.
func TestHotfixVersionFilter(t *testing.T) {
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Initialize git-flow
	_, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
	if err != nil {
		t.Fatalf("Failed to initialize git-flow: %v", err)
	}

	// Create a version filter that adds 'hotfix-' prefix for hotfixes
	script := `#!/bin/sh
VERSION="$1"
if [ "${VERSION#hotfix-}" = "$VERSION" ]; then
    echo "hotfix-$VERSION"
else
    echo "$VERSION"
fi
`
	createHookScript(t, dir, "filter-flow-hotfix-start-version", script)

	// Start a hotfix - filter should change "1.0.1" to "hotfix-1.0.1"
	output, err := testutil.RunGitFlow(t, dir, "hotfix", "start", "1.0.1")
	if err != nil {
		t.Fatalf("Failed to start hotfix: %v\nOutput: %s", err, output)
	}

	// Verify the filter message is shown
	if !strings.Contains(output, "Version filter changed") {
		t.Errorf("Expected output to mention version filter, got: %s", output)
	}

	// Verify the branch was created with the filtered name
	if !testutil.BranchExists(t, dir, "hotfix/hotfix-1.0.1") {
		t.Error("Expected hotfix/hotfix-1.0.1 branch to exist (filtered from 1.0.1)")
	}

	// Verify original name branch was NOT created
	if testutil.BranchExists(t, dir, "hotfix/1.0.1") {
		t.Error("hotfix/1.0.1 should not exist - filter should have changed it to hotfix-1.0.1")
	}
}
