# Coding Guidelines

This document outlines the coding standards and conventions used in the git-flow-next project. Following these guidelines ensures consistency and maintainability across the codebase.

## Go Code Style

### Package Organization

**Package Naming:**
- Use single word, lowercase names: `config`, `git`, `errors`, `util`
- Avoid underscores or mixed case
- Keep package names short and descriptive

**File Organization:**
- One file per major command in `cmd/`
- Group related functionality within packages
- Use descriptive filenames: `config.go`, `repo.go`, `mergestate.go`

### Import Organization

Organize imports in three distinct groups with blank lines between:

```go
import (
    // 1. Standard library packages (alphabetical)
    "fmt"
    "os"
    "strings"

    // 2. Third-party packages (alphabetical)
    "github.com/spf13/cobra"

    // 3. Local packages (alphabetical)
    "github.com/gittower/git-flow-next/internal/config"
    "github.com/gittower/git-flow-next/internal/errors"
    "github.com/gittower/git-flow-next/internal/git"
)
```

### Naming Conventions

**Functions:**
- PascalCase for exported functions: `LoadConfig()`, `FinishCommand()`
- camelCase for private functions: `executeFinish()`, `handleContinue()`
- Use descriptive verb names: `ValidateBranchName()`, `CreateTag()`

**Variables:**
- camelCase consistently: `branchName`, `mergeStrategy`, `configValue`
- Avoid abbreviations unless they're well-known: `cfg` for config, `err` for error
- Use descriptive names: `parentBranch` not `parent`

**Constants:**
- PascalCase for exported: `ExitCodeSuccess`, `DefaultTimeout`
- Use groups with descriptive prefixes:

```go
// Step constants for state machine
const (
    stepMerge          = "merge"
    stepCreateTag      = "create_tag"
    stepUpdateChildren = "update_children"
    stepDeleteBranch   = "delete_branch"
)
```

**Types:**
- PascalCase for exported structs: `BranchConfig`, `TagOptions`
- Use descriptive, unambiguous names
- Suffix with purpose when helpful: `BranchRetentionOptions`, `MergeState`

### Struct Definitions

Structure fields logically and provide clear documentation:

```go
// TagOptions contains options for tag creation when finishing a branch
type TagOptions struct {
    ShouldTag   *bool  // Whether to create a tag (nil means use config default)
    ShouldSign  *bool  // Whether to sign the tag (nil means use config default)
    SigningKey  string // GPG signing key to use
    Message     string // Custom tag message
    MessageFile string // Path to file containing tag message
    TagName     string // Custom tag name
}
```

**Guidelines:**
- Group related fields together
- Use pointer types for optional boolean values (`*bool`)
- Document when `nil` has special meaning
- Align field comments for readability

## Error Handling

### Custom Error Types

Define structured errors with specific types and exit codes:

```go
type BranchNotFoundError struct {
    BranchName string
}

func (e *BranchNotFoundError) Error() string {
    return fmt.Sprintf("branch '%s' not found", e.BranchName)
}

func (e *BranchNotFoundError) ExitCode() ExitCode {
    return ExitCodeBranchNotFound
}
```

### Error Handling Pattern

Always handle errors explicitly with appropriate context:

```go
output, err := git.GetConfig(configKey)
if err != nil {
    return &errors.GitError{
        Operation: fmt.Sprintf("get config '%s'", configKey),
        Err:       err,
    }
}
```

**Guidelines:**
- Never ignore errors (`_ = someFunction()` is prohibited)
- Wrap errors with context using structured error types
- Return specific error types for different failure conditions
- Use `fmt.Errorf()` with `%w` verb for error wrapping when appropriate

### Exit Codes

Define meaningful exit codes for different error conditions:

```go
const (
    ExitCodeSuccess               ExitCode = 0
    ExitCodeGeneral              ExitCode = 1
    ExitCodeNotInitialized       ExitCode = 2
    ExitCodeBranchNotFound       ExitCode = 3
    ExitCodeMergeConflict        ExitCode = 4
    ExitCodeInvalidBranchType    ExitCode = 5
    ExitCodeUncommittedChanges   ExitCode = 6
)
```

## Command Structure

### Universal Command Architecture

All commands must follow this exact three-layer pattern:

```go
// Layer 1: Cobra Command Handler (in init() or command creation)
RunE: func(cmd *cobra.Command, args []string) error {
    // Parse flags and arguments
    param1, _ := cmd.Flags().GetString("flag1")
    param2, _ := cmd.Flags().GetBool("flag2")
    
    // Call command wrapper
    CommandName(branchType, name, param1, param2)
    return nil
}

// Layer 2: Command Wrapper - Handles error conversion and exit codes
func CommandName(params...) {
    if err := executeCommand(params...); err != nil {
        var exitCode errors.ExitCode
        if flowErr, ok := err.(errors.Error); ok {
            exitCode = flowErr.ExitCode()
        } else {
            exitCode = errors.ExitCodeGitError // Default fallback
        }
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(int(exitCode))
    }
}

// Layer 3: Execute Function - Contains actual business logic
func executeCommand(params...) error {
    // All business logic here
    // Return structured errors only
    return nil
}
```

### Configuration Loading Pattern

**CRITICAL**: Load configuration once at the beginning and pass through all calls:

```go
func executeCommand(params...) error {
    // Load config ONCE at the start
    cfg, err := config.LoadConfig()
    if err != nil {
        return &errors.GitError{Operation: "load configuration", Err: err}
    }
    
    // Pass cfg to all subsequent function calls
    return doWork(cfg, params...)
}
```

**Guidelines:**
- Never call `config.LoadConfig()` multiple times within a command
- Pass `cfg` parameter to all helper functions that need configuration
- Cache configuration lookups within the command execution

### Standard Command Initialization

Every command (except `init`) must start with this exact pattern:

```go
// Validate that git-flow is initialized
initialized, err := config.IsInitialized()
if err != nil {
    return &errors.GitError{Operation: "check if git-flow is initialized", Err: err}
}
if !initialized {
    return &errors.NotInitializedError{}
}
```

### Branch Type Validation

For topic branch commands:

```go
// Get branch configuration
branchConfig, ok := cfg.Branches[branchType]
if !ok {
    return &errors.InvalidBranchTypeError{BranchType: branchType}
}
```

### Branch Name Resolution

For commands accepting branch names:

```go
// Construct full branch name if needed
fullBranchName := name
if branchConfig.Prefix != "" && !strings.HasPrefix(name, branchConfig.Prefix) {
    fullBranchName = branchConfig.Prefix + name
}

// Check if branch exists
if err := git.BranchExists(fullBranchName); err != nil {
    return &errors.BranchNotFoundError{BranchName: fullBranchName}
}
```

### Current Branch Detection

When commands can work on current branch:

```go
var branchName string
if name == "" {
    currentBranch, err := git.GetCurrentBranch()
    if err != nil {
        return &errors.GitError{Operation: "get current branch", Err: err}
    }
    branchName = currentBranch
} else {
    branchName = name
}
```

### Configuration Override Pattern

For configurable options with multiple precedence levels:

```go
// 1. Start with branch configuration default
shouldTag := branchConfig.Tag

// 2. Check for branch-specific config override
branchSpecificConfig, err := git.GetConfig(fmt.Sprintf("gitflow.%s.finish.notag", branchType))
if err == nil && branchSpecificConfig == "true" {
    shouldTag = false
}

// 3. Command-line flags override everything
if tagOptions != nil && tagOptions.ShouldTag != nil {
    shouldTag = *tagOptions.ShouldTag
}
```

### Input Validation

Always validate inputs early in the execute function:

```go
// Validate required inputs
if name == "" {
    return &errors.EmptyBranchNameError{}
}

// Validate branch exists before operations
if err := git.BranchExists(branchName); err != nil {
    return &errors.BranchNotFoundError{BranchName: branchName}
}
```

### Options Struct Pattern

For commands with multiple flags, use structured options:

```go
type TagOptions struct {
    ShouldTag   *bool  // nil means use config default
    ShouldSign  *bool  // nil means use config default
    SigningKey  string
    Message     string
    TagName     string
}

type BranchRetentionOptions struct {
    Keep        *bool // Whether to keep the branch
    KeepRemote  *bool // Whether to keep remote branch
    KeepLocal   *bool // Whether to keep local branch
    ForceDelete *bool // Whether to force delete
}
```

### Command Function Signatures

Use consistent signatures:
- Command wrapper: `func CommandName(branchType, name string, options...)`
- Execute function: `func executeCommand(branchType, name string, options...) error`

## Git Operations

### Wrapper Functions

All Git operations must go through wrapper functions in `internal/git/`:

```go
// internal/git/repo.go
func Checkout(branch string) error {
    cmd := exec.Command("git", "checkout", branch)
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("failed to checkout branch '%s': %w", branch, err)
    }
    return nil
}
```

**Guidelines:**
- Never call `git` commands directly in command functions
- Provide clear error messages with context
- Handle common Git errors appropriately
- Use consistent parameter validation

## State Management

### Persistent State

For complex multi-step operations, use persistent state:

```go
type MergeState struct {
    Action          string   `json:"action"`
    BranchType      string   `json:"branch_type"`
    BranchName      string   `json:"branch_name"`
    CurrentStep     string   `json:"current_step"`
    ParentBranch    string   `json:"parent_branch"`
    MergeStrategy   string   `json:"merge_strategy"`
    FullBranchName  string   `json:"full_branch_name"`
    ChildBranches   []string `json:"child_branches"`
    UpdatedBranches []string `json:"updated_branches"`
}
```

**Guidelines:**
- Use JSON tags for serialization
- Include all information needed to resume operations
- Provide clear field names and types
- Document state transitions

## Output and User Communication

### Output Patterns

Use consistent patterns for user communication:

```go
// Regular progress output
fmt.Printf("Merging using strategy: %s\n", strategy)

// Error output to stderr
fmt.Fprintf(os.Stderr, "Error: %s\n", err)

// Success messages with context
fmt.Printf("Successfully finished branch '%s' and updated %d child branches\n", 
    branchName, len(updatedBranches))
```

### User-Friendly Messages

Provide clear, actionable error messages:

```go
func (e *MergeConflictError) Error() string {
    return fmt.Sprintf("Merge conflicts detected. Resolve conflicts and run 'git flow %s finish --continue %s'", 
        e.BranchType, e.BranchName)
}
```

## Testing

### Test Organization

- Tests mirror source structure in `test/` directory
- Use descriptive test names: `TestFinishFeatureBranchWithMergeConflict`
- Group related tests in the same file

### Test Structure

Follow consistent test structure:

```go
func TestExample(t *testing.T) {
    // Setup
    dir := testutil.SetupTestRepo(t)
    defer testutil.CleanupTestRepo(t, dir)
    
    // Execute
    output, err := testutil.RunGitFlow(t, dir, "feature", "start", "test")
    
    // Assert
    if err != nil {
        t.Fatalf("Expected success, got error: %v\nOutput: %s", err, output)
    }
    
    // Verify state
    if !testutil.BranchExists(t, dir, "feature/test") {
        t.Error("Expected feature branch to be created")
    }
}
```

### Test Utilities

Use shared test utilities for common operations:

```go
// Setup and cleanup
dir := testutil.SetupTestRepo(t)
defer testutil.CleanupTestRepo(t, dir)

// Git operations
_, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
testutil.WriteFile(t, dir, "file.txt", "content")

// Assertions
if !testutil.BranchExists(t, dir, "feature/test") {
    t.Error("Expected branch to exist")
}
```

## Documentation

### Function Documentation

Document all exported functions with clear descriptions:

```go
// LoadConfig loads the git-flow configuration from git config.
// It reads all gitflow.branch.* configuration keys and constructs
// a Config struct with branch definitions and settings.
func LoadConfig() (*Config, error) {
    // Implementation
}
```

### Package Documentation

Provide package-level documentation at the top of main package files:

```go
// Package config provides git-flow configuration management.
// It handles loading branch type definitions and workflow settings
// from Git configuration files.
package config
```

### Complex Logic Documentation

Document complex algorithms and state machines:

```go
// The finish operation progresses through sequential steps:
// 1. merge: Merge topic branch into parent
// 2. create_tag: Create tag if configured  
// 3. update_children: Update dependent child branches
// 4. delete_branch: Clean up topic branch
//
// Each step can be interrupted by conflicts and resumed with --continue.
func finish(state *mergestate.MergeState) error {
    // Implementation
}
```

## Quality Standards

### Mandatory Change Requirements

**CRITICAL**: When changing code, you MUST always complete all three steps:

1. **Create/Adjust Tests**
   - Add new tests for new functionality
   - Update existing tests when behavior changes
   - Follow test naming conventions: `TestFunctionNameWithScenario`
   - Use test utilities from `testutil/` package
   - Ensure tests cover both success and error cases

2. **Run All Tests**
   - Execute `go test ./...` to run complete test suite
   - Fix any test failures before committing
   - Ensure no regressions in existing functionality
   - Verify new tests pass consistently

3. **Update Documentation**
   - Update relevant `.md` files for user-facing changes
   - Update function/package documentation for API changes
   - Update CLAUDE.md for development workflow changes
   - Update configuration examples if config changes

**Test Command Reference:**
```bash
# Run all tests
go test ./...

# Run tests for specific package
go test ./test/cmd/
go test ./test/internal/

# Run specific test file with verbose output
go test -v ./test/cmd/init_test.go
```

### Code Reviews

All code changes must:
- Follow these coding guidelines
- **Complete all three mandatory requirements above**
- Have clear commit messages following COMMIT_GUIDELINES.md
- Pass all existing tests
- Include comprehensive test coverage

### Static Analysis

The project uses:
- `go fmt` for consistent formatting
- `go vet` for basic static analysis
- `golint` for style checking (when available)

### Performance Considerations

- Minimize Git operations in large repositories
- Cache configuration lookups where appropriate
- Use efficient data structures for branch operations
- Avoid unnecessary string allocations in hot paths

These guidelines help maintain code quality and ensure consistency across the git-flow-next project. When in doubt, follow the patterns established in existing code and prioritize clarity and maintainability.