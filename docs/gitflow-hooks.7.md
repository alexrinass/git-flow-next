# GITFLOW-HOOKS(7)

## NAME

gitflow-hooks - Git-flow hooks and filters for customizing workflow behavior

## SYNOPSIS

Custom scripts in `.git/hooks/` to extend and customize git-flow operations.

## DESCRIPTION

Git-flow supports two types of extension points:

**Filters** transform input/output values during operations (e.g., modifying version numbers or tag messages).

**Hooks** execute scripts before and after operations (e.g., running CI checks before starting a release).

All scripts are located in the `.git/hooks/` directory following specific naming conventions.

## FILTERS

Filters transform values passed to git-flow commands. They receive input as arguments or via environment variables and output the transformed value to stdout.

### Filter Naming Convention

```
filter-flow-{type}-{action}-{target}
```

Where:
- `{type}` is the branch type (release, hotfix)
- `{action}` is the git-flow action (start, finish)
- `{target}` is what is being filtered (version, tag-message)

### Available Filters

| Filter Name | Command | Purpose |
|-------------|---------|---------|
| `filter-flow-release-start-version` | `git flow release start` | Modify version number |
| `filter-flow-hotfix-start-version` | `git flow hotfix start` | Modify version number |
| `filter-flow-release-finish-tag-message` | `git flow release finish` | Customize tag message |
| `filter-flow-hotfix-finish-tag-message` | `git flow hotfix finish` | Customize tag message |

### Version Filters

Version filters receive the version as the first argument (`$1`) and should output the modified version to stdout.

**Example: Add 'v' prefix to versions**

```bash
#!/bin/sh
# .git/hooks/filter-flow-release-start-version
VERSION="$1"
if [ "${VERSION#v}" = "$VERSION" ]; then
    echo "v$VERSION"
else
    echo "$VERSION"
fi
```

**Example: Enforce semantic versioning**

```bash
#!/bin/sh
# .git/hooks/filter-flow-release-start-version
VERSION="$1"

# Validate semver format
if ! echo "$VERSION" | grep -qE '^v?[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?$'; then
    echo "Error: Version must follow semantic versioning (e.g., 1.0.0)" >&2
    exit 1
fi

echo "$VERSION"
```

### Tag Message Filters

Tag message filters receive:
- `$1` - Version number
- `$2` - Original tag message

Environment variables:
- `BRANCH_TYPE` - The branch type (release, hotfix)
- `BRANCH_NAME` - The short branch name
- `BASE_BRANCH` - The parent/base branch
- `VERSION` - The version number

**Example: Append changelog to tag message**

```bash
#!/bin/sh
# .git/hooks/filter-flow-release-finish-tag-message
VERSION="$1"
MESSAGE="$2"

# Get recent commits since last tag
LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
if [ -n "$LAST_TAG" ]; then
    CHANGELOG=$(git log --oneline "${LAST_TAG}..HEAD" 2>/dev/null)
else
    CHANGELOG=$(git log --oneline -10 2>/dev/null)
fi

echo "${MESSAGE}

Changes in this release:
${CHANGELOG}"
```

### Filter Behavior

- If a filter does not exist, the original value is used
- If a filter is not executable, it is skipped silently
- If a filter exits with non-zero status, the operation fails with an error
- If a filter outputs nothing, the original value is used
- Filter output is trimmed of leading/trailing whitespace

## HOOKS

Hooks are scripts that execute before (pre) or after (post) git-flow operations. Pre-hooks can prevent operations from proceeding.

### Hook Naming Convention

```
{pre,post}-flow-{type}-{action}
```

Where:
- `{pre,post}` indicates when the hook runs
- `{type}` is the branch type (feature, release, hotfix, support)
- `{action}` is the git-flow action (start, finish, publish, track, delete)

### Available Hooks

| Hook Pattern | Operations |
|--------------|------------|
| `{pre,post}-flow-feature-{action}` | start, finish, publish, track, delete |
| `{pre,post}-flow-release-{action}` | start, finish, publish, track, delete |
| `{pre,post}-flow-hotfix-{action}` | start, finish, publish, delete |
| `{pre,post}-flow-support-{action}` | start, finish, publish, delete |

### Environment Variables

Hooks receive context via environment variables:

| Variable | Description |
|----------|-------------|
| `BRANCH` | Full branch name (e.g., `feature/my-feature`) |
| `BRANCH_NAME` | Short name (e.g., `my-feature`) |
| `BRANCH_TYPE` | Type (e.g., `feature`) |
| `BASE_BRANCH` | Parent branch (e.g., `develop`) |
| `ORIGIN` | Remote name |
| `VERSION` | Version (for release/hotfix) |
| `EXIT_CODE` | Post-hooks only: exit code of the operation |

### Pre-hooks

Pre-hooks run before the operation. If a pre-hook exits with a non-zero status, the operation is aborted.

**Example: Verify CI is passing before release**

```bash
#!/bin/sh
# .git/hooks/pre-flow-release-start

if command -v gh &> /dev/null; then
    STATUS=$(gh run list --branch "${BASE_BRANCH:-develop}" --limit 1 --json conclusion -q '.[0].conclusion')
    if [ "$STATUS" != "success" ]; then
        echo "Error: CI is not passing on ${BASE_BRANCH:-develop}"
        exit 1
    fi
fi

exit 0
```

**Example: Check for uncommitted changes**

```bash
#!/bin/sh
# .git/hooks/pre-flow-feature-start

if ! git diff-index --quiet HEAD --; then
    echo "Error: You have uncommitted changes. Please commit or stash them first."
    exit 1
fi

exit 0
```

### Post-hooks

Post-hooks run after the operation completes, regardless of success or failure. The `EXIT_CODE` environment variable indicates the operation result.

**Example: Send notification on release completion**

```bash
#!/bin/sh
# .git/hooks/post-flow-release-finish

if [ "$EXIT_CODE" -eq 0 ]; then
    echo "Release $VERSION completed successfully!"
    # Send Slack notification, update ticket system, etc.
fi
```

**Example: Update documentation after feature finish**

```bash
#!/bin/sh
# .git/hooks/post-flow-feature-finish

if [ "$EXIT_CODE" -eq 0 ]; then
    echo "Feature $BRANCH_NAME merged. Consider updating documentation."
fi
```

### Hook Behavior

- If a hook does not exist, the operation proceeds normally
- If a hook is not executable, it is skipped silently
- Pre-hooks that exit non-zero abort the operation
- Post-hooks always run (success or failure), their exit codes are ignored
- Hook output is displayed to the user

## CREATING HOOK SCRIPTS

1. Create the script in `.git/hooks/` with the appropriate name
2. Make the script executable: `chmod +x .git/hooks/<script-name>`
3. Test the script manually before relying on it

### Tips

- Always start scripts with a shebang (`#!/bin/sh` or `#!/bin/bash`)
- Use `exit 0` for success, `exit 1` (or any non-zero) for failure
- Write error messages to stderr: `echo "Error: ..." >&2`
- Test scripts in isolation before enabling them
- Consider using version control for your hook scripts (see below)

## SHARING HOOKS

Since `.git/hooks/` is not tracked by Git, consider these approaches for sharing hooks:

### Using a hooks directory in the repository

```bash
# Store hooks in a tracked directory
mkdir .githooks

# Configure Git to use this directory
git config core.hooksPath .githooks
```

### Using symbolic links

```bash
# Create symbolic links in .git/hooks/
ln -s ../../.githooks/pre-flow-release-start .git/hooks/
```

## EXAMPLES

### Complete Version Filter

```bash
#!/bin/sh
# .git/hooks/filter-flow-release-start-version
# Validates and normalizes version numbers

VERSION="$1"

# Remove 'v' prefix for validation
CLEAN_VERSION="${VERSION#v}"

# Validate semver format (major.minor.patch)
if ! echo "$CLEAN_VERSION" | grep -qE '^[0-9]+\.[0-9]+\.[0-9]+$'; then
    echo "Error: Invalid version format. Use MAJOR.MINOR.PATCH" >&2
    exit 1
fi

# Always output with 'v' prefix
echo "v$CLEAN_VERSION"
```

### Complete Pre-hook

```bash
#!/bin/sh
# .git/hooks/pre-flow-feature-start
# Ensures feature names follow conventions

NAME="$BRANCH_NAME"

# Check for valid characters
if ! echo "$NAME" | grep -qE '^[a-z0-9-]+$'; then
    echo "Error: Feature name must contain only lowercase letters, numbers, and hyphens" >&2
    exit 1
fi

# Check minimum length
if [ ${#NAME} -lt 3 ]; then
    echo "Error: Feature name must be at least 3 characters long" >&2
    exit 1
fi

exit 0
```

## SEE ALSO

**git-flow**(1), **git-flow-start**(1), **git-flow-finish**(1), **githooks**(5)

## REFERENCES

- [git-flow-avh Hooks and Filters](https://github.com/petervanderdoes/gitflow-avh/wiki/Reference:-Hooks-and-Filters)
- [Git Hooks Documentation](https://git-scm.com/docs/githooks)
