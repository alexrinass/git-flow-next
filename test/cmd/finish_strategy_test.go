package cmd_test

import (
	"strings"
	"testing"

	"github.com/gittower/git-flow-next/test/testutil"
)

// TestFinishWithRebaseFlag tests that the --rebase flag works correctly.
// Steps:
// 1. Sets up a test repository and initializes git-flow with defaults
// 2. Creates a feature branch with commits
// 3. Makes changes in develop branch
// 4. Finishes the feature branch with --rebase flag
// 5. Verifies the branch was rebased before merging
func TestFinishWithRebaseFlag(t *testing.T) {
	// Setup
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Initialize git-flow with defaults
	output, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
	if err != nil {
		t.Fatalf("Failed to initialize git-flow: %v\nOutput: %s", err, output)
	}

	// Create and switch to feature branch
	output, err = testutil.RunGitFlow(t, dir, "feature", "start", "rebase-test")
	if err != nil {
		t.Fatalf("Failed to start feature branch: %v\nOutput: %s", err, output)
	}

	// Add a commit to feature branch
	testutil.WriteFile(t, dir, "feature.txt", "feature content")
	_, err = testutil.RunGit(t, dir, "add", "feature.txt")
	if err != nil {
		t.Fatalf("Failed to add feature file: %v", err)
	}
	_, err = testutil.RunGit(t, dir, "commit", "-m", "Add feature file")
	if err != nil {
		t.Fatalf("Failed to commit feature file: %v", err)
	}

	// Switch to develop and add a commit
	_, err = testutil.RunGit(t, dir, "checkout", "develop")
	if err != nil {
		t.Fatalf("Failed to switch to develop: %v", err)
	}
	testutil.WriteFile(t, dir, "develop.txt", "develop content")
	_, err = testutil.RunGit(t, dir, "add", "develop.txt")
	if err != nil {
		t.Fatalf("Failed to add develop file: %v", err)
	}
	_, err = testutil.RunGit(t, dir, "commit", "-m", "Add develop file")
	if err != nil {
		t.Fatalf("Failed to commit develop file: %v", err)
	}

	// Finish feature branch with --rebase flag
	output, err = testutil.RunGitFlow(t, dir, "feature", "finish", "rebase-test", "--rebase")
	if err != nil {
		t.Fatalf("Failed to finish feature branch: %v\nOutput: %s", err, output)
	}

	// Verify that the output indicates rebase strategy was used
	if !strings.Contains(output, "Merging using strategy: rebase") {
		t.Errorf("Expected output to indicate rebase strategy, got: %s", output)
	}

	// Verify that both files exist in develop
	_, err = testutil.RunGit(t, dir, "checkout", "develop")
	if err != nil {
		t.Fatalf("Failed to checkout develop: %v", err)
	}

	if !testutil.FileExists(t, dir, "feature.txt") {
		t.Error("Expected feature.txt to exist in develop branch")
	}
	if !testutil.FileExists(t, dir, "develop.txt") {
		t.Error("Expected develop.txt to exist in develop branch")
	}
}

// TestFinishWithSquashFlag tests that the --squash flag works correctly.
// Steps:
// 1. Sets up a test repository and initializes git-flow with defaults
// 2. Creates a feature branch with multiple commits
// 3. Finishes the feature branch with --squash flag
// 4. Verifies the commits were squashed into a single commit
func TestFinishWithSquashFlag(t *testing.T) {
	// Setup
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Initialize git-flow with defaults
	output, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
	if err != nil {
		t.Fatalf("Failed to initialize git-flow: %v\nOutput: %s", err, output)
	}

	// Create and switch to feature branch
	output, err = testutil.RunGitFlow(t, dir, "feature", "start", "squash-test")
	if err != nil {
		t.Fatalf("Failed to start feature branch: %v\nOutput: %s", err, output)
	}

	// Add multiple commits to feature branch
	testutil.WriteFile(t, dir, "file1.txt", "first file")
	_, err = testutil.RunGit(t, dir, "add", "file1.txt")
	if err != nil {
		t.Fatalf("Failed to add first file: %v", err)
	}
	_, err = testutil.RunGit(t, dir, "commit", "-m", "Add first file")
	if err != nil {
		t.Fatalf("Failed to commit first file: %v", err)
	}

	testutil.WriteFile(t, dir, "file2.txt", "second file")
	_, err = testutil.RunGit(t, dir, "add", "file2.txt")
	if err != nil {
		t.Fatalf("Failed to add second file: %v", err)
	}
	_, err = testutil.RunGit(t, dir, "commit", "-m", "Add second file")
	if err != nil {
		t.Fatalf("Failed to commit second file: %v", err)
	}

	// Count commits in feature branch before finish
	beforeLog, err := testutil.RunGit(t, dir, "log", "--oneline", "develop..HEAD")
	if err != nil {
		t.Fatalf("Failed to get commit log: %v", err)
	}
	beforeCommitCount := len(strings.Split(strings.TrimSpace(beforeLog), "\n"))
	if beforeCommitCount != 2 {
		t.Fatalf("Expected 2 commits in feature branch, got %d", beforeCommitCount)
	}

	// Finish feature branch with --squash flag (keep branch to avoid deletion issues)
	output, err = testutil.RunGitFlow(t, dir, "feature", "finish", "squash-test", "--squash", "--keep")
	if err != nil {
		t.Fatalf("Failed to finish feature branch: %v\nOutput: %s", err, output)
	}

	// Verify that the output indicates squash strategy was used
	if !strings.Contains(output, "Merging using strategy: squash") {
		t.Errorf("Expected output to indicate squash strategy, got: %s", output)
	}

	// Verify that both files exist in develop
	_, err = testutil.RunGit(t, dir, "checkout", "develop")
	if err != nil {
		t.Fatalf("Failed to checkout develop: %v", err)
	}

	if !testutil.FileExists(t, dir, "file1.txt") {
		t.Error("Expected file1.txt to exist in develop branch")
	}
	if !testutil.FileExists(t, dir, "file2.txt") {
		t.Error("Expected file2.txt to exist in develop branch")
	}
}

// TestFinishWithNoFFFlag tests that the --no-ff flag works correctly.
// Steps:
// 1. Sets up a test repository and initializes git-flow with defaults
// 2. Creates a feature branch with commits
// 3. Finishes the feature branch with --no-ff flag
// 4. Verifies a merge commit was created even for fast-forward case
func TestFinishWithNoFFFlag(t *testing.T) {
	// Setup
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Initialize git-flow with defaults
	output, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
	if err != nil {
		t.Fatalf("Failed to initialize git-flow: %v\nOutput: %s", err, output)
	}

	// Create and switch to feature branch
	output, err = testutil.RunGitFlow(t, dir, "feature", "start", "no-ff-test")
	if err != nil {
		t.Fatalf("Failed to start feature branch: %v\nOutput: %s", err, output)
	}

	// Add a commit to feature branch
	testutil.WriteFile(t, dir, "feature.txt", "feature content")
	_, err = testutil.RunGit(t, dir, "add", "feature.txt")
	if err != nil {
		t.Fatalf("Failed to add feature file: %v", err)
	}
	_, err = testutil.RunGit(t, dir, "commit", "-m", "Add feature file")
	if err != nil {
		t.Fatalf("Failed to commit feature file: %v", err)
	}

	// Finish feature branch with --no-ff flag
	output, err = testutil.RunGitFlow(t, dir, "feature", "finish", "no-ff-test", "--no-ff")
	if err != nil {
		t.Fatalf("Failed to finish feature branch: %v\nOutput: %s", err, output)
	}

	// Switch to develop branch to check merge commit
	_, err = testutil.RunGit(t, dir, "checkout", "develop")
	if err != nil {
		t.Fatalf("Failed to checkout develop: %v", err)
	}

	// Check that the latest commit has multiple parents (indicating a merge commit)
	parents, err := testutil.RunGit(t, dir, "rev-list", "--parents", "-n", "1", "HEAD")
	if err != nil {
		t.Fatalf("Failed to get parent commits: %v", err)
	}

	// A merge commit should have at least 2 parents (so at least 3 hashes in the output)
	parentCount := len(strings.Fields(strings.TrimSpace(parents)))
	if parentCount < 3 {
		t.Errorf("Expected merge commit with multiple parents, got %d parents", parentCount-1)
	}

	// Verify that feature file exists in develop
	if !testutil.FileExists(t, dir, "feature.txt") {
		t.Error("Expected feature.txt to exist in develop branch")
	}
}

// TestMergeStrategyFlagOverridesConfig tests that command-line flags override configuration.
// Steps:
// 1. Sets up a test repository and initializes git-flow with defaults
// 2. Sets configuration to use merge strategy
// 3. Creates a feature branch
// 4. Finishes with --rebase flag to override config
// 5. Verifies rebase strategy was used instead of configured merge
func TestMergeStrategyFlagOverridesConfig(t *testing.T) {
	// Setup
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Initialize git-flow with defaults
	output, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
	if err != nil {
		t.Fatalf("Failed to initialize git-flow: %v\nOutput: %s", err, output)
	}

	// Set configuration to use merge strategy for features
	_, err = testutil.RunGit(t, dir, "config", "gitflow.branch.feature.upstreamstrategy", "merge")
	if err != nil {
		t.Fatalf("Failed to set upstream strategy config: %v", err)
	}

	// Verify the config was set
	strategy, err := testutil.RunGit(t, dir, "config", "--get", "gitflow.branch.feature.upstreamstrategy")
	if err != nil {
		t.Fatalf("Failed to get upstream strategy config: %v", err)
	}
	if strings.TrimSpace(strategy) != "merge" {
		t.Fatalf("Expected merge strategy in config, got: %s", strings.TrimSpace(strategy))
	}

	// Create and switch to feature branch
	output, err = testutil.RunGitFlow(t, dir, "feature", "start", "config-override-test")
	if err != nil {
		t.Fatalf("Failed to start feature branch: %v\nOutput: %s", err, output)
	}

	// Add a commit to feature branch
	testutil.WriteFile(t, dir, "feature.txt", "feature content")
	_, err = testutil.RunGit(t, dir, "add", "feature.txt")
	if err != nil {
		t.Fatalf("Failed to add feature file: %v", err)
	}
	_, err = testutil.RunGit(t, dir, "commit", "-m", "Add feature file")
	if err != nil {
		t.Fatalf("Failed to commit feature file: %v", err)
	}

	// Finish feature branch with --rebase flag (should override config)
	output, err = testutil.RunGitFlow(t, dir, "feature", "finish", "config-override-test", "--rebase")
	if err != nil {
		t.Fatalf("Failed to finish feature branch: %v\nOutput: %s", err, output)
	}

	// Verify that rebase strategy was used (override config)
	if !strings.Contains(output, "Merging using strategy: rebase") {
		t.Errorf("Expected rebase strategy to override config, got: %s", output)
	}

	// Verify that feature file exists in develop
	_, err = testutil.RunGit(t, dir, "checkout", "develop")
	if err != nil {
		t.Fatalf("Failed to checkout develop: %v", err)
	}

	if !testutil.FileExists(t, dir, "feature.txt") {
		t.Error("Expected feature.txt to exist in develop branch")
	}
}

// TestMergeStrategyConfigUsedByDefault tests that branch configuration is used when no flags are provided.
// Steps:
// 1. Sets up a test repository and initializes git-flow with defaults
// 2. Sets configuration to use rebase strategy
// 3. Creates a feature branch
// 4. Finishes without any flags
// 5. Verifies rebase strategy from config was used
func TestMergeStrategyConfigUsedByDefault(t *testing.T) {
	// Setup
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Initialize git-flow with defaults
	output, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
	if err != nil {
		t.Fatalf("Failed to initialize git-flow: %v\nOutput: %s", err, output)
	}

	// Set configuration to use rebase strategy for features
	_, err = testutil.RunGit(t, dir, "config", "gitflow.branch.feature.upstreamstrategy", "rebase")
	if err != nil {
		t.Fatalf("Failed to set upstream strategy config: %v", err)
	}

	// Create and switch to feature branch
	output, err = testutil.RunGitFlow(t, dir, "feature", "start", "config-default-test")
	if err != nil {
		t.Fatalf("Failed to start feature branch: %v\nOutput: %s", err, output)
	}

	// Add a commit to feature branch
	testutil.WriteFile(t, dir, "feature.txt", "feature content")
	_, err = testutil.RunGit(t, dir, "add", "feature.txt")
	if err != nil {
		t.Fatalf("Failed to add feature file: %v", err)
	}
	_, err = testutil.RunGit(t, dir, "commit", "-m", "Add feature file")
	if err != nil {
		t.Fatalf("Failed to commit feature file: %v", err)
	}

	// Switch to develop and add a commit
	_, err = testutil.RunGit(t, dir, "checkout", "develop")
	if err != nil {
		t.Fatalf("Failed to switch to develop: %v", err)
	}
	testutil.WriteFile(t, dir, "develop.txt", "develop content")
	_, err = testutil.RunGit(t, dir, "add", "develop.txt")
	if err != nil {
		t.Fatalf("Failed to add develop file: %v", err)
	}
	_, err = testutil.RunGit(t, dir, "commit", "-m", "Add develop file")
	if err != nil {
		t.Fatalf("Failed to commit develop file: %v", err)
	}

	// Finish feature branch without any flags (should use config)
	output, err = testutil.RunGitFlow(t, dir, "feature", "finish", "config-default-test")
	if err != nil {
		t.Fatalf("Failed to finish feature branch: %v\nOutput: %s", err, output)
	}

	// Verify that rebase strategy from config was used
	if !strings.Contains(output, "Merging using strategy: rebase") {
		t.Errorf("Expected rebase strategy from config to be used, got: %s", output)
	}

	// Verify that both files exist in develop
	_, err = testutil.RunGit(t, dir, "checkout", "develop")
	if err != nil {
		t.Fatalf("Failed to checkout develop: %v", err)
	}

	if !testutil.FileExists(t, dir, "feature.txt") {
		t.Error("Expected feature.txt to exist in develop branch")
	}
	if !testutil.FileExists(t, dir, "develop.txt") {
		t.Error("Expected develop.txt to exist in develop branch")
	}
}
