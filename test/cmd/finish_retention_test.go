package cmd_test

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gittower/git-flow-next/test/testutil"
)

// TestFinishFeatureBranchDefaultLocalDeletion tests that feature branches are deleted locally by default.
// Steps:
// 1. Sets up a test repository and initializes git-flow
// 2. Creates a feature branch
// 3. Adds changes to the feature branch
// 4. Finishes the feature branch
// 5. Verifies the local branch is deleted
func TestFinishFeatureBranchDefaultLocalDeletion(t *testing.T) {
	// Setup
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Initialize git-flow with defaults and create branches
	output, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
	if err != nil {
		t.Fatalf("Failed to initialize git-flow: %v\nOutput: %s", err, output)
	}

	// Create a feature branch
	output, err = testutil.RunGitFlow(t, dir, "feature", "start", "my-feature")
	if err != nil {
		t.Fatalf("Failed to create feature branch: %v\nOutput: %s", err, output)
	}

	// Create a test file
	testutil.WriteFile(t, dir, "test.txt", "test content")

	// Commit the changes
	_, err = testutil.RunGit(t, dir, "add", "test.txt")
	if err != nil {
		t.Fatalf("Failed to add file: %v", err)
	}
	_, err = testutil.RunGit(t, dir, "commit", "-m", "Add test file")
	if err != nil {
		t.Fatalf("Failed to commit file: %v", err)
	}

	// Finish the feature branch
	output, err = testutil.RunGitFlow(t, dir, "feature", "finish", "my-feature")
	if err != nil {
		t.Fatalf("Failed to finish feature branch: %v\nOutput: %s", err, output)
	}

	// Verify that local feature branch is deleted
	if testutil.BranchExists(t, dir, "feature/my-feature") {
		t.Error("Expected local feature branch to be deleted by default")
	}
}

// TestFinishFeatureBranchDefaultRemoteDeletion tests that feature branches are deleted both locally and remotely by default.
// Steps:
// 1. Sets up a test repository and initializes git-flow
// 2. Creates a feature branch
// 3. Adds changes to the feature branch
// 4. Adds a remote repository
// 5. Finishes the feature branch
// 6. Verifies both local and remote branches are deleted
func TestFinishFeatureBranchDefaultRemoteDeletion(t *testing.T) {
	// Setup
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Initialize git-flow with defaults and create branches
	output, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
	if err != nil {
		t.Fatalf("Failed to initialize git-flow: %v\nOutput: %s", err, output)
	}

	// Create a feature branch
	output, err = testutil.RunGitFlow(t, dir, "feature", "start", "my-feature")
	if err != nil {
		t.Fatalf("Failed to create feature branch: %v\nOutput: %s", err, output)
	}

	// Create a test file
	testutil.WriteFile(t, dir, "test.txt", "test content")

	// Commit the changes
	_, err = testutil.RunGit(t, dir, "add", "test.txt")
	if err != nil {
		t.Fatalf("Failed to add file: %v", err)
	}
	_, err = testutil.RunGit(t, dir, "commit", "-m", "Add test file")
	if err != nil {
		t.Fatalf("Failed to commit file: %v", err)
	}

	// Add a remote repository
	remoteDir, err := testutil.AddRemote(t, dir, "origin", true)
	if err != nil {
		t.Fatalf("Failed to add remote: %v", err)
	}
	defer testutil.CleanupTestRepo(t, remoteDir)

	// Push the feature branch to remote
	_, err = testutil.RunGit(t, dir, "push", "origin", "feature/my-feature")
	if err != nil {
		t.Fatalf("Failed to push feature branch: %v", err)
	}

	// Finish the feature branch
	output, err = testutil.RunGitFlow(t, dir, "feature", "finish", "my-feature")
	if err != nil {
		t.Fatalf("Failed to finish feature branch: %v\nOutput: %s", err, output)
	}

	// Verify that local feature branch is deleted
	if testutil.BranchExists(t, dir, "feature/my-feature") {
		t.Error("Expected local feature branch to be deleted by default")
	}

	// Verify that remote feature branch is deleted
	_, err = testutil.RunGit(t, dir, "fetch", "origin")
	if err != nil {
		t.Fatalf("Failed to fetch from remote: %v", err)
	}
	if testutil.BranchExists(t, dir, "origin/feature/my-feature") {
		t.Error("Expected remote feature branch to be deleted by default")
	}
}

// TestFinishFeatureBranchKeepLocal tests that the keep-local option preserves the local branch when finishing.
// Steps:
// 1. Sets up a test repository and initializes git-flow
// 2. Creates a feature branch
// 3. Adds changes to the feature branch
// 4. Finishes the feature branch with the keeplocal option
// 5. Verifies the branch is merged into develop
// 6. Verifies the local feature branch is preserved
func TestFinishFeatureBranchKeepLocal(t *testing.T) {
	// Setup
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Initialize git-flow with defaults
	output, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
	if err != nil {
		t.Fatalf("Failed to initialize git-flow: %v\nOutput: %s", err, output)
	}

	// Create a feature branch
	output, err = testutil.RunGitFlow(t, dir, "feature", "start", "keep-local-test")
	if err != nil {
		t.Fatalf("Failed to create feature branch: %v\nOutput: %s", err, output)
	}

	// Create a test file
	testutil.WriteFile(t, dir, "test.txt", "feature content")

	// Commit the changes
	_, err = testutil.RunGit(t, dir, "add", "test.txt")
	if err != nil {
		t.Fatalf("Failed to add file: %v", err)
	}
	_, err = testutil.RunGit(t, dir, "commit", "-m", "Add test file")
	if err != nil {
		t.Fatalf("Failed to commit file: %v", err)
	}

	// Get path to the git-flow binary
	gitFlowPath, err := filepath.Abs(filepath.Join("..", "..", "git-flow"))
	if err != nil {
		t.Fatalf("Failed to get absolute path to git-flow: %v", err)
	}

	// Finish the feature branch with keeplocal option
	cmd := exec.Command(gitFlowPath, "feature", "finish", "keep-local-test", "--keeplocal")
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to finish feature branch: %v\nOutput: %s", err, stdout.String()+stderr.String())
	}
	output = stdout.String() + stderr.String()

	// Verify that we're now on develop branch
	currentBranch := testutil.GetCurrentBranch(t, dir)
	if currentBranch != "develop" {
		t.Errorf("Expected to be on develop branch after finish, got %s", currentBranch)
	}

	// Verify that the changes are in develop
	if !testutil.FileExists(t, dir, "test.txt") {
		t.Error("Expected test.txt to exist in develop branch")
	}

	// Verify that the local branch still exists
	if !testutil.BranchExists(t, dir, "feature/keep-local-test") {
		t.Error("Expected feature branch to still exist with --keeplocal option")
	}

	// Checkout the feature branch and verify it still has the content
	_, err = testutil.RunGit(t, dir, "checkout", "feature/keep-local-test")
	if err != nil {
		t.Fatalf("Failed to checkout feature branch: %v", err)
	}

	// Verify the file content in the feature branch
	content := testutil.ReadFile(t, dir, "test.txt")
	if content != "feature content" {
		t.Errorf("Expected file content to be 'feature content', got '%s'", content)
	}
}

// TestFinishFeatureBranchKeepRemote tests that the keep-remote option preserves the remote branch when finishing.
// Steps:
// 1. Sets up a test repository and initializes git-flow
// 2. Creates a feature branch
// 3. Adds changes to the feature branch
// 4. Adds a remote and pushes the branch
// 5. Finishes the feature branch with the keepremote option
// 6. Verifies the branch is merged into develop
// 7. Verifies the local feature branch is deleted
// 8. Verifies the remote feature branch is preserved
func TestFinishFeatureBranchKeepRemote(t *testing.T) {
	// Setup
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Initialize git-flow with defaults
	output, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
	if err != nil {
		t.Fatalf("Failed to initialize git-flow: %v\nOutput: %s", err, output)
	}

	// Create a feature branch
	output, err = testutil.RunGitFlow(t, dir, "feature", "start", "keep-remote-test")
	if err != nil {
		t.Fatalf("Failed to create feature branch: %v\nOutput: %s", err, output)
	}

	// Create a test file
	testutil.WriteFile(t, dir, "test.txt", "feature content")

	// Commit the changes
	_, err = testutil.RunGit(t, dir, "add", "test.txt")
	if err != nil {
		t.Fatalf("Failed to add file: %v", err)
	}
	_, err = testutil.RunGit(t, dir, "commit", "-m", "Add test file")
	if err != nil {
		t.Fatalf("Failed to commit file: %v", err)
	}

	// Add a remote repository
	remoteDir, err := testutil.AddRemote(t, dir, "origin", true)
	if err != nil {
		t.Fatalf("Failed to add remote: %v", err)
	}
	defer testutil.CleanupTestRepo(t, remoteDir)

	// Push the feature branch to remote
	_, err = testutil.RunGit(t, dir, "push", "origin", "feature/keep-remote-test")
	if err != nil {
		t.Fatalf("Failed to push feature branch: %v", err)
	}

	// Get path to the git-flow binary
	gitFlowPath, err := filepath.Abs(filepath.Join("..", "..", "git-flow"))
	if err != nil {
		t.Fatalf("Failed to get absolute path to git-flow: %v", err)
	}

	// Finish the feature branch with keepremote option
	cmd := exec.Command(gitFlowPath, "feature", "finish", "keep-remote-test", "--keepremote")
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to finish feature branch: %v\nOutput: %s", err, stdout.String()+stderr.String())
	}
	output = stdout.String() + stderr.String()

	// Verify that we're now on develop branch
	currentBranch := testutil.GetCurrentBranch(t, dir)
	if currentBranch != "develop" {
		t.Errorf("Expected to be on develop branch after finish, got %s", currentBranch)
	}

	// Verify that the changes are in develop
	if !testutil.FileExists(t, dir, "test.txt") {
		t.Error("Expected test.txt to exist in develop branch")
	}

	// Verify that the local branch is deleted
	if testutil.BranchExists(t, dir, "feature/keep-remote-test") {
		t.Error("Expected local feature branch to be deleted")
	}

	// Fetch from remote to update references
	_, err = testutil.RunGit(t, dir, "fetch", "origin")
	if err != nil {
		t.Fatalf("Failed to fetch from remote: %v", err)
	}

	// Verify that the remote branch still exists
	remoteExists := testutil.RemoteBranchExists(t, dir, "origin", "feature/keep-remote-test")
	if !remoteExists {
		t.Error("Expected remote feature branch to still exist with --keepremote option")
	}

	// Try to checkout the remote branch and verify it has the content
	_, err = testutil.RunGit(t, dir, "checkout", "-b", "verify-remote", "origin/feature/keep-remote-test")
	if err != nil {
		t.Fatalf("Failed to checkout remote branch: %v", err)
	}

	// Verify the file content from the remote branch
	content := testutil.ReadFile(t, dir, "test.txt")
	if content != "feature content" {
		t.Errorf("Expected file content to be 'feature content', got '%s'", content)
	}
}
