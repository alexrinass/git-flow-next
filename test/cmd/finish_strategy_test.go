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

// TestFinishWithSquashMessageFlag tests that the --squash-message flag customizes the commit message.
// Steps:
// 1. Sets up a test repository and initializes git-flow with defaults
// 2. Creates a feature branch with multiple commits
// 3. Finishes the feature branch with --squash and --squash-message flags
// 4. Verifies the custom commit message was used
func TestFinishWithSquashMessageFlag(t *testing.T) {
	// Setup
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Initialize git-flow with defaults
	output, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
	if err != nil {
		t.Fatalf("Failed to initialize git-flow: %v\nOutput: %s", err, output)
	}

	// Create and switch to feature branch
	output, err = testutil.RunGitFlow(t, dir, "feature", "start", "squash-msg-test")
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

	// Finish feature branch with custom squash message
	customMessage := "feat: add login functionality"
	output, err = testutil.RunGitFlow(t, dir, "feature", "finish", "squash-msg-test", "--squash", "--squash-message", customMessage, "--keep")
	if err != nil {
		t.Fatalf("Failed to finish feature branch: %v\nOutput: %s", err, output)
	}

	// Verify that squash strategy was used
	if !strings.Contains(output, "Merging using strategy: squash") {
		t.Errorf("Expected output to indicate squash strategy, got: %s", output)
	}

	// Verify the commit message in develop
	_, err = testutil.RunGit(t, dir, "checkout", "develop")
	if err != nil {
		t.Fatalf("Failed to checkout develop: %v", err)
	}

	// Get the last commit message
	commitMsg, err := testutil.RunGit(t, dir, "log", "-1", "--format=%s")
	if err != nil {
		t.Fatalf("Failed to get commit message: %v", err)
	}

	if strings.TrimSpace(commitMsg) != customMessage {
		t.Errorf("Expected commit message '%s', got '%s'", customMessage, strings.TrimSpace(commitMsg))
	}
}

// TestSquashMessageIgnoredWithoutSquashStrategy tests that --squash-message is ignored when not using squash strategy.
// Steps:
// 1. Sets up a test repository and initializes git-flow with defaults
// 2. Creates a feature branch with commits
// 3. Finishes with --squash-message but WITHOUT --squash flag
// 4. Verifies the commit message is NOT the custom squash message (uses merge commit message instead)
func TestSquashMessageIgnoredWithoutSquashStrategy(t *testing.T) {
	// Setup
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Initialize git-flow with defaults
	output, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
	if err != nil {
		t.Fatalf("Failed to initialize git-flow: %v\nOutput: %s", err, output)
	}

	// Create and switch to feature branch
	output, err = testutil.RunGitFlow(t, dir, "feature", "start", "no-squash-test")
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

	// Finish feature branch with squash-message but WITHOUT --squash flag (using default merge strategy)
	ignoredMessage := "this message should be ignored"
	output, err = testutil.RunGitFlow(t, dir, "feature", "finish", "no-squash-test", "--squash-message", ignoredMessage, "--no-ff", "--keep")
	if err != nil {
		t.Fatalf("Failed to finish feature branch: %v\nOutput: %s", err, output)
	}

	// Verify that merge strategy was used (not squash)
	if strings.Contains(output, "Merging using strategy: squash") {
		t.Errorf("Expected merge strategy, but squash was used: %s", output)
	}

	// Verify the commit message in develop is NOT the squash message
	_, err = testutil.RunGit(t, dir, "checkout", "develop")
	if err != nil {
		t.Fatalf("Failed to checkout develop: %v", err)
	}

	// Get the last commit message
	commitMsg, err := testutil.RunGit(t, dir, "log", "-1", "--format=%s")
	if err != nil {
		t.Fatalf("Failed to get commit message: %v", err)
	}

	// The commit message should be a merge commit message, NOT the squash message
	if strings.TrimSpace(commitMsg) == ignoredMessage {
		t.Errorf("Squash message should be ignored when not using squash strategy, but got: '%s'", commitMsg)
	}

	// It should be a merge commit message
	if !strings.Contains(commitMsg, "Merge branch") {
		t.Errorf("Expected merge commit message, got: '%s'", commitMsg)
	}
}

// TestSquashMessagePreservedAfterConflict tests that --squash-message is preserved in merge state after conflict.
// Steps:
// 1. Sets up a test repository and initializes git-flow with defaults
// 2. Creates a feature branch with conflicting changes
// 3. Attempts to finish with --squash --squash-message (causes conflict)
// 4. Resolves conflict and continues WITHOUT --squash-message
// 5. Verifies the original custom squash message was used (preserved from state)
func TestSquashMessagePreservedAfterConflict(t *testing.T) {
	// Setup
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Initialize git-flow with defaults
	output, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
	if err != nil {
		t.Fatalf("Failed to initialize git-flow: %v\nOutput: %s", err, output)
	}

	// Create conflicting content in develop
	testutil.RunGit(t, dir, "checkout", "develop")
	testutil.WriteFile(t, dir, "conflict.txt", "Develop version\nLine 2\nLine 3")
	testutil.RunGit(t, dir, "add", "conflict.txt")
	testutil.RunGit(t, dir, "commit", "-m", "Develop changes")

	// Create feature with conflicting changes
	output, err = testutil.RunGitFlow(t, dir, "feature", "start", "squash-conflict")
	if err != nil {
		t.Fatalf("Failed to create feature: %v\nOutput: %s", err, output)
	}

	testutil.WriteFile(t, dir, "conflict.txt", "Feature version\nLine 2 modified\nLine 3")
	testutil.RunGit(t, dir, "add", "conflict.txt")
	testutil.RunGit(t, dir, "commit", "-m", "Feature changes")

	// Meanwhile, add more changes to develop to create conflict
	testutil.RunGit(t, dir, "checkout", "develop")
	testutil.WriteFile(t, dir, "conflict.txt", "Develop version updated\nLine 2\nLine 3 modified")
	testutil.RunGit(t, dir, "add", "conflict.txt")
	testutil.RunGit(t, dir, "commit", "-m", "More develop changes")

	// Try to finish feature with squash (will conflict)
	customMessage := "feat: squash merged after conflict resolution"
	testutil.RunGit(t, dir, "checkout", "feature/squash-conflict")
	output, _ = testutil.RunGitFlow(t, dir, "feature", "finish", "squash-conflict", "--squash", "--squash-message", customMessage)

	// Verify conflict detected
	if !strings.Contains(output, "conflict") {
		t.Fatal("Expected merge conflict to be detected")
	}

	// Verify merge state exists with squash strategy and message
	state, err := testutil.LoadMergeState(t, dir)
	if err != nil || state == nil {
		t.Fatal("Expected merge state to exist after conflict")
	}

	if state.MergeStrategy != "squash" {
		t.Errorf("Expected MergeStrategy to be 'squash', got: %s", state.MergeStrategy)
	}

	if state.SquashMessage != customMessage {
		t.Errorf("Expected SquashMessage in state to be '%s', got: '%s'", customMessage, state.SquashMessage)
	}

	// Resolve conflict
	testutil.WriteFile(t, dir, "conflict.txt", "Resolved version\nLine 2 resolved\nLine 3 resolved")
	testutil.RunGit(t, dir, "add", "conflict.txt")

	// Continue finish operation WITHOUT --squash-message (should use preserved message from state)
	output, err = testutil.RunGitFlow(t, dir, "feature", "finish", "--continue", "squash-conflict")
	if err != nil {
		t.Fatalf("Failed to continue finish: %v\nOutput: %s", err, output)
	}

	// Verify success
	if !strings.Contains(output, "Successfully finished") {
		t.Error("Expected successful finish message")
	}

	// Verify the commit message in develop uses the preserved message
	testutil.RunGit(t, dir, "checkout", "develop")
	commitMsg, err := testutil.RunGit(t, dir, "log", "-1", "--format=%s")
	if err != nil {
		t.Fatalf("Failed to get commit message: %v", err)
	}

	if strings.TrimSpace(commitMsg) != customMessage {
		t.Errorf("Expected commit message '%s', got '%s'", customMessage, strings.TrimSpace(commitMsg))
	}
}

// TestSquashMessageOverrideOnContinue tests that --squash-message can be overridden on continue.
// Steps:
// 1. Sets up a test repository and initializes git-flow with defaults
// 2. Creates a feature branch with conflicting changes
// 3. Attempts to finish with --squash --squash-message (causes conflict)
// 4. Resolves conflict and continues WITH a different --squash-message
// 5. Verifies the new message was used (overriding the preserved one)
func TestSquashMessageOverrideOnContinue(t *testing.T) {
	// Setup
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)

	// Initialize git-flow with defaults
	output, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
	if err != nil {
		t.Fatalf("Failed to initialize git-flow: %v\nOutput: %s", err, output)
	}

	// Create conflicting content in develop
	testutil.RunGit(t, dir, "checkout", "develop")
	testutil.WriteFile(t, dir, "conflict.txt", "Develop version\nLine 2\nLine 3")
	testutil.RunGit(t, dir, "add", "conflict.txt")
	testutil.RunGit(t, dir, "commit", "-m", "Develop changes")

	// Create feature with conflicting changes
	output, err = testutil.RunGitFlow(t, dir, "feature", "start", "squash-override")
	if err != nil {
		t.Fatalf("Failed to create feature: %v\nOutput: %s", err, output)
	}

	testutil.WriteFile(t, dir, "conflict.txt", "Feature version\nLine 2 modified\nLine 3")
	testutil.RunGit(t, dir, "add", "conflict.txt")
	testutil.RunGit(t, dir, "commit", "-m", "Feature changes")

	// Meanwhile, add more changes to develop to create conflict
	testutil.RunGit(t, dir, "checkout", "develop")
	testutil.WriteFile(t, dir, "conflict.txt", "Develop version updated\nLine 2\nLine 3 modified")
	testutil.RunGit(t, dir, "add", "conflict.txt")
	testutil.RunGit(t, dir, "commit", "-m", "More develop changes")

	// Try to finish feature with squash (will conflict)
	originalMessage := "feat: original message"
	testutil.RunGit(t, dir, "checkout", "feature/squash-override")
	output, _ = testutil.RunGitFlow(t, dir, "feature", "finish", "squash-override", "--squash", "--squash-message", originalMessage)

	// Verify conflict detected
	if !strings.Contains(output, "conflict") {
		t.Fatal("Expected merge conflict to be detected")
	}

	// Resolve conflict
	testutil.WriteFile(t, dir, "conflict.txt", "Resolved version\nLine 2 resolved\nLine 3 resolved")
	testutil.RunGit(t, dir, "add", "conflict.txt")

	// Continue finish operation WITH a different --squash-message (should override)
	overrideMessage := "feat: overridden message after conflict"
	output, err = testutil.RunGitFlow(t, dir, "feature", "finish", "--continue", "--squash-message", overrideMessage, "squash-override")
	if err != nil {
		t.Fatalf("Failed to continue finish: %v\nOutput: %s", err, output)
	}

	// Verify success
	if !strings.Contains(output, "Successfully finished") {
		t.Error("Expected successful finish message")
	}

	// Verify the commit message in develop uses the override message
	testutil.RunGit(t, dir, "checkout", "develop")
	commitMsg, err := testutil.RunGit(t, dir, "log", "-1", "--format=%s")
	if err != nil {
		t.Fatalf("Failed to get commit message: %v", err)
	}

	if strings.TrimSpace(commitMsg) != overrideMessage {
		t.Errorf("Expected commit message '%s', got '%s'", overrideMessage, strings.TrimSpace(commitMsg))
	}

	// Make sure it's NOT the original message
	if strings.TrimSpace(commitMsg) == originalMessage {
		t.Errorf("Expected override message to be used, but got original message: '%s'", commitMsg)
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
