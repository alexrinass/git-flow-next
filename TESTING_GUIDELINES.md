# Testing Guidelines

This document outlines the testing conventions and patterns used in the git-flow-next project.

## Test Function Naming

Follow descriptive naming patterns that clearly indicate what is being tested:

```go
// Basic functionality
func TestStartFeatureBranch(t *testing.T)
func TestFinishFeatureBranch(t *testing.T)
func TestInitWithDefaults(t *testing.T)

// Configuration-specific tests
func TestStartWithCustomConfig(t *testing.T)
func TestInitWithAVHConfig(t *testing.T)

// Error conditions
func TestUpdateWithMergeConflict(t *testing.T)
func TestFinishWithMergeConflict(t *testing.T)
func TestDeleteNonExistentRemoteBranch(t *testing.T)

// Feature-specific tests
func TestFinishFeatureBranchWithFetchFlag(t *testing.T)
func TestStartWithFetchFlag(t *testing.T)
```

## Test Comments and State Documentation

**REQUIRED**: Every test function must have a structured comment that follows this exact pattern:

1. **First line**: Brief description of what the test validates
2. **Steps:** section with numbered list of test actions
3. **Expected outcomes** embedded in the steps

### Comment Pattern (MANDATORY)

```go
// TestFinishWithMergeConflict tests the behavior when finishing a branch with merge conflicts.
// Steps:
// 1. Sets up a test repository and initializes git-flow
// 2. Creates a feature branch
// 3. Adds conflicting changes to both feature and develop branches
// 4. Attempts to finish the feature branch
// 5. Verifies the operation fails with merge conflict
func TestFinishWithMergeConflict(t *testing.T) {
    // Test implementation...
}
```

### Additional Comment Guidelines

- **Be specific**: Include exact branch names, commands, and expected results
- **Number all steps**: Use sequential numbering (1, 2, 3...)  
- **Include verification**: Always specify what is being verified
- **Use active voice**: "Creates a branch" not "A branch is created"
- **Match test structure**: Comments should reflect actual test implementation

### Examples of Good Test Comments

```go
// TestConfigAddBase tests adding base branch configurations.
// Steps:
// 1. Sets up a test repository and initializes git-flow
// 2. Adds various base branches with different configurations  
// 3. Verifies branches are created and configuration is saved correctly
// 4. Tests error conditions like duplicate branches and invalid parents
func TestConfigAddBase(t *testing.T) { ... }

// TestStartFeatureBranch tests the start command for feature branches.
// Steps:
// 1. Sets up a test repository and initializes git-flow with defaults
// 2. Runs 'git flow feature start my-feature'
// 3. Verifies feature/my-feature branch is created
// 4. Verifies branch is based on develop branch
func TestStartFeatureBranch(t *testing.T) { ... }
```

## One Test Case Per Function Rule

**CRITICAL**: Each test function must test exactly one scenario or behavior. Never use table-driven tests with multiple test cases in a single function.

### ❌ BAD: Multiple Test Cases in One Function

```go
func TestMergeStrategyConfigHierarchy(t *testing.T) {
    testCases := []struct {
        name           string
        branchConfig   string
        commandConfig  string
        flag           string
        expectedResult string
    }{
        {"branch_default_merge", "merge", "", "", "merge"},
        {"command_overrides_branch", "merge", "rebase=true", "", "rebase"},
        {"flag_overrides_config", "merge", "rebase=true", "--no-rebase", "merge"},
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Test implementation...
        })
    }
}
```

### ✅ GOOD: One Test Case Per Function

```go
// TestMergeStrategyBranchDefault tests that branch default strategy is used.
// Steps:
// 1. Sets up repository with branch configured for merge strategy
// 2. Runs finish command without overrides
// 3. Verifies merge strategy was used
func TestMergeStrategyBranchDefault(t *testing.T) {
    // Single test scenario implementation...
}

// TestMergeStrategyCommandConfigOverride tests command config overriding branch default.
// Steps:
// 1. Sets up repository with branch configured for merge strategy
// 2. Sets gitflow.feature.finish.rebase=true in config
// 3. Runs finish command without flags
// 4. Verifies rebase strategy was used instead of merge
func TestMergeStrategyCommandConfigOverride(t *testing.T) {
    // Single test scenario implementation...
}

// TestMergeStrategyFlagOverridesConfig tests command flag overriding config.
// Steps:
// 1. Sets up repository with rebase configured via command config
// 2. Runs finish command with --no-rebase flag
// 3. Verifies merge strategy was used instead of rebase
func TestMergeStrategyFlagOverridesConfig(t *testing.T) {
    // Single test scenario implementation...
}
```

### Why One Test Case Per Function?

1. **Clear failure identification** - When a test fails, you immediately know which specific scenario failed
2. **Focused debugging** - Each test tests one thing, making debugging straightforward
3. **Maintainable test suite** - Easy to modify, disable, or extend individual test scenarios
4. **Better test names** - Function names clearly describe what is being tested
5. **Proper test isolation** - Each test sets up exactly what it needs, nothing more
6. **Follows git-flow-next philosophy** - Explicit and readable over clever and compact

### Exception: Helper Functions

Table-driven tests are acceptable in helper functions that test pure functions with multiple input/output combinations:

```go
func TestValidateBranchName(t *testing.T) {
    // This is acceptable for pure validation functions
    testCases := []struct {
        name    string
        input   string
        isValid bool
    }{
        {"valid_name", "feature-branch", true},
        {"invalid_spaces", "feature branch", false},
        {"invalid_colon", "feature:branch", false},
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := util.ValidateBranchName(tc.input)
            assert.Equal(t, tc.isValid, result == nil)
        })
    }
}
```

**Rule**: Only use table-driven tests for testing pure functions with simple input/output validation, never for complex integration scenarios.

## Temporary Git Repository Testing

All tests use temporary Git repositories created through test utilities.

### Default Branch Configuration

When using `git flow init --defaults`, the following branches and settings are configured:

#### Base Branches
- **main**: Production branch (root branch)
- **develop**: Integration branch (auto-updates from main)

#### Topic Branch Types  
- **feature**: Prefix `feature/`, parent `develop`, starts from `develop`
- **release**: Prefix `release/`, parent `main`, starts from `develop` 
- **hotfix**: Prefix `hotfix/`, parent `main`, starts from `main`

#### Default Merge Strategies
- **Feature finish**: `merge` into `develop`
- **Release finish**: `merge` into `main` (then auto-update `develop`)
- **Hotfix finish**: `merge` into `main` (then auto-update `develop`)
- **Feature update**: `rebase` from `develop`
- **Release update**: `merge` from `main`
- **Hotfix update**: `rebase` from `main`

#### Default Tag Settings
- **Feature**: No tags created
- **Release**: Tags created on finish
- **Hotfix**: Tags created on finish

### Testing with Default Configuration

- **Use git-flow defaults by default**: Initialize repositories with `git flow init --defaults`
- **Use feature branches for topic branch testing**: Feature branches are the most common topic branch type and should be used for general topic branch functionality tests
- **Only create standalone branches when required**: If your test specifically needs a non-git-flow branch or custom configuration, create it explicitly
- **Adjust configuration when needed**: Use `git config` commands to modify behavior for specific test scenarios

```go
// Standard test setup - use this for most tests
func TestExample(t *testing.T) {
    dir := testutil.SetupTestRepo(t)
    defer testutil.CleanupTestRepo(t, dir)
    
    // Initialize with defaults - creates main, develop, and configures feature/release/hotfix
    output, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
    if err != nil {
        t.Fatalf("Failed to initialize git-flow: %v\nOutput: %s", err, output)
    }
    
    // Use feature branches for topic branch testing
    output, err = testutil.RunGitFlow(t, dir, "feature", "start", "test-branch")
    // ... test implementation
}

// Example: Testing with modified merge strategy
func TestFeatureFinishWithRebase(t *testing.T) {
    dir := testutil.SetupTestRepo(t)
    defer testutil.CleanupTestRepo(t, dir)
    
    // Initialize with defaults
    output, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
    if err != nil {
        t.Fatalf("Failed to initialize git-flow: %v", err)
    }
    
    // Modify merge strategy for this test
    _, err = testutil.RunGit(t, dir, "config", "gitflow.feature.finish.merge", "rebase")
    if err != nil {
        t.Fatalf("Failed to configure rebase strategy: %v", err)
    }
    
    // Now test with the modified configuration
    output, err = testutil.RunGitFlow(t, dir, "feature", "start", "rebase-test")
    // ... test implementation
}
```

### Basic Test Setup

```go
func TestExample(t *testing.T) {
    // Setup temporary repository
    dir := testutil.SetupTestRepo(t)
    defer testutil.CleanupTestRepo(t, dir)
    
    // Initialize git-flow with defaults
    output, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
    if err != nil {
        t.Fatalf("Failed to initialize git-flow: %v\nOutput: %s", err, output)
    }
    
    // Test implementation...
}
```

### Git Test Helper Functions

Available in `test/testutil/git.go`:

- `SetupTestRepo(t *testing.T) string` - Creates temporary Git repository
- `CleanupTestRepo(t *testing.T, dir string)` - Removes temporary repository
- `RunGit(t *testing.T, dir string, args ...string)` - Executes Git commands
- `RunGitFlow(t *testing.T, dir string, args ...string)` - Executes git-flow commands
- `WriteFile(t *testing.T, dir string, name string, content string)` - Creates files
- `BranchExists(t *testing.T, dir string, branch string) bool` - Checks branch existence
- `GetCurrentBranch(t *testing.T, dir string) string` - Gets current branch name

## Creating Git Scenarios for Tests

### Merge Conflicts

To create merge conflicts for testing:

1. **Both Modified Conflict (Preferred)**:
   - Create a file with initial content on the base branch
   - Create a new branch and modify the file
   - Switch back to base branch and modify the same file differently
   - Attempting to merge will create a conflict

```go
// Create feature branch
output, err := testutil.RunGitFlow(t, dir, "feature", "start", "conflict-test")

// Create file in feature branch
testutil.WriteFile(t, dir, "test.txt", "feature content")
_, err = testutil.RunGit(t, dir, "add", "test.txt")
_, err = testutil.RunGit(t, dir, "commit", "-m", "Add test.txt in feature")

// Switch to develop and create conflicting content
_, err = testutil.RunGit(t, dir, "checkout", "develop")
testutil.WriteFile(t, dir, "test.txt", "develop content")
_, err = testutil.RunGit(t, dir, "add", "test.txt")
_, err = testutil.RunGit(t, dir, "commit", "-m", "Add test.txt in develop")

// Attempt to finish - should fail with conflict
output, err = testutil.RunGitFlow(t, dir, "feature", "finish", "conflict-test")
```

2. **Multiple Rebase Conflicts**:
   - Repeat the above scenario with additional commits and files
   - Each conflicting commit will create a separate rebase step

### Using Remotes

To test remote functionality:

1. **Create Bare Remote Repository**:
   - Use `testutil.AddRemote()` to create a local bare repository
   - Add it as a remote using the file path
   - Push branches to establish tracking

```go
// Add a remote repository
remoteDir, err := testutil.AddRemote(t, dir, "origin", true)
if err != nil {
    t.Fatalf("Failed to add remote: %v", err)
}
defer testutil.CleanupTestRepo(t, remoteDir)

// Create feature branch and make changes
output, err = testutil.RunGitFlow(t, dir, "feature", "start", "fetch-test")
testutil.WriteFile(t, dir, "feature.txt", "feature content")
_, err = testutil.RunGit(t, dir, "add", "feature.txt")
_, err = testutil.RunGit(t, dir, "commit", "-m", "Add feature file")

// Test operations with remote (fetch, push, etc.)
output, err = testutil.RunGitFlow(t, dir, "feature", "finish", "fetch-test", "--fetch")
```

## Test Organization

### Directory Structure

```
test/
├── cmd/              # Command-level integration tests
├── internal/         # Internal package unit tests
└── testutil/         # Test utilities and helpers
```

### Error Handling in Tests

Always include comprehensive error checking with detailed failure messages:

```go
output, err := testutil.RunGitFlow(t, dir, "feature", "start", "test-branch")
if err != nil {
    t.Fatalf("Failed to start feature branch: %v\nOutput: %s", err, output)
}
```

### State Verification

Verify both Git state and application-specific state:

```go
// Check Git state
if !testutil.BranchExists(t, dir, "feature/test-branch") {
    t.Error("Expected feature branch to exist")
}

// Check application state
state, err := testutil.LoadMergeState(t, dir)
if state.Action != "finish" {
    t.Errorf("Expected action to be 'finish', got '%s'", state.Action)
}
```

## Best Practices

1. **Always use testutil helpers** - Never execute Git commands directly
2. **Include setup/cleanup** - Use defer to ensure cleanup happens
3. **Test error conditions** - Verify failures behave correctly
4. **Check intermediate state** - Don't just test final outcomes
5. **Use descriptive assertions** - Include context in error messages
6. **Test with remotes when relevant** - Many Git operations behave differently with remotes