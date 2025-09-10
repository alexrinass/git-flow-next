# GIT-FLOW-LIST(1)

## NAME

git-flow-list - List topic branches

## SYNOPSIS

**git-flow** *topic* **list** [*pattern*]

## DESCRIPTION

List all topic branches of the specified type. This command works with any topic branch type (feature, release, hotfix, support, or custom types defined in your configuration).

The list command displays all local branches that match the topic branch type's prefix pattern, with optional filtering by name pattern.

## ARGUMENTS

*topic*
: The topic branch type (feature, release, hotfix, support, or any configured custom type)

*pattern*
: Optional pattern to filter branch names. Supports shell-style wildcards.

## OUTPUT FORMAT

The list command displays branches in this format:
```
  branch-name-1
* branch-name-2    (current branch, marked with asterisk)
  branch-name-3
```

Branch names are shown without the prefix for readability.

## EXAMPLES

### Basic Usage

List all feature branches:
```bash
git flow feature list
```

List all release branches:
```bash
git flow release list
```

List all hotfix branches:
```bash
git flow hotfix list
```

### Pattern Filtering

List features matching pattern:
```bash
git flow feature list "user-*"
```

List releases for version 1.x:
```bash
git flow release list "1.*"
```

List hotfixes containing "security":
```bash
git flow hotfix list "*security*"
```

### Custom Branch Types

List custom branch type:
```bash
git flow bugfix list
git flow epic list
```

## PATTERN MATCHING

Patterns support shell-style globbing:

- **`*`** - Matches any characters
- **`?`** - Matches single character  
- **`[abc]`** - Matches any character in brackets
- **`[a-z]`** - Matches any character in range

Examples:
```bash
git flow feature list "api-*"        # Starts with "api-"
git flow feature list "*-test"       # Ends with "-test"
git flow feature list "user-?"       # "user-" plus one character
git flow feature list "[a-m]*"       # Starts with letters a through m
```

## BRANCH STATUS INDICATORS

**Current Branch**
: Marked with asterisk (`*`) when you're currently on that branch

**Tracking Information**
: Shows if branch has remote tracking (implementation dependent)

## WORKFLOW INTEGRATION

### Before Starting Work
```bash
# Check existing features before starting new one
git flow feature list
git flow feature start new-feature
```

### During Development
```bash
# See all active releases
git flow release list

# Check specific pattern
git flow feature list "*api*"
```

### Cleanup Planning
```bash
# List all features to identify candidates for deletion
git flow feature list
git flow feature delete old-feature
```

## EMPTY RESULTS

When no branches match the criteria:
```bash
$ git flow feature list
No feature branches found.
```

With pattern filtering:
```bash
$ git flow feature list "nonexistent-*"
No feature branches matching 'nonexistent-*' found.
```

## CONFIGURATION

List behavior is controlled by branch type configuration:

### Prefix Configuration
```bash
git config gitflow.branch.feature.prefix feature/
git config gitflow.branch.release.prefix release/
git config gitflow.branch.hotfix.prefix hotfix/
```

The list command automatically uses the configured prefix to identify branches of each type.

## REMOTE BRANCHES

The list command shows only local branches by default. To see remote branches:

```bash
# Show remote feature branches
git branch -r | grep feature/

# Or use Git directly
git ls-remote --heads origin | grep feature/
```

## EXIT STATUS

**0**
: Successful listing (even if no branches found)

**1**
: Invalid topic branch type

**2**
: Git operation failed

**3**
: Invalid pattern syntax

## SEE ALSO

**git-flow**(1), **git-flow-start**(1), **git-flow-delete**(1), **git-branch**(1), **git-ls-remote**(1)

## NOTES

- Only shows local branches - use **git branch -r** for remote branches
- Branch names are displayed without the prefix for readability
- Pattern matching uses shell-style globbing, not regex
- Empty results are not considered an error condition
- Current branch is highlighted with an asterisk
- Custom topic branch types work exactly like built-in types