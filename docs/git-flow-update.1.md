# GIT-FLOW-UPDATE(1)

## NAME

git-flow-update - Update topic branches with parent changes

## SYNOPSIS

**git-flow** *topic* **update** [*name*] [*options*]

**git-flow update** [*name*] [*options*]

## DESCRIPTION

Update a topic branch with the latest changes from its parent branch using the configured downstream merge strategy. This command works with any topic branch type (feature, release, hotfix, support, or custom types).

The update operation merges or rebases changes from the parent branch into the topic branch, keeping the topic branch current with the latest development.

## ARGUMENTS

*topic*
: The topic branch type (feature, release, hotfix, support, or any configured custom type)

*name*
: Name of the topic branch to update. If omitted, the current branch is used (when using shorthand **git-flow update**)

## OPTIONS

**--rebase**
: Force rebase strategy instead of the configured downstream strategy

## MERGE STRATEGIES

The merge strategy used when updating is determined by configuration:

**merge** (default for most branches)
: Standard Git merge preserving both branch histories

**rebase**
: Rebase topic branch commits on top of parent branch for linear history

The strategy is configured via:
1. **Branch defaults**: `gitflow.branch.<type>.downstreamStrategy`  
2. **Command overrides**: `gitflow.<type>.downstreamStrategy`
3. **Flag override**: `--rebase` always forces rebase strategy

## EXAMPLES

### Basic Usage

Update current topic branch (shorthand):
```bash
git flow update
```

Update specific feature branch:
```bash
git flow feature update user-authentication
```

Update release branch with latest hotfixes:
```bash
git flow release update 1.2.0
```

### Force Rebase Strategy

Update using rebase regardless of configuration:
```bash
git flow feature update my-feature --rebase
```

Update current branch with rebase:
```bash
git flow update --rebase
```

### Typical Workflows

Before finishing a long-running feature:
```bash
# Make sure feature is current before finishing
git flow feature update long-feature
git flow feature finish long-feature
```

Update release with production hotfixes:
```bash
# Hotfix was applied to main, update release
git flow release update 1.2.0
```

Keep hotfix current during development:
```bash
# Update hotfix with any other main changes
git flow hotfix update critical-fix
```

## CONFLICT RESOLUTION

When conflicts occur during update:

```bash
# Start update
git flow feature update my-feature

# Conflicts occur - Git will show conflict markers
# Edit files to resolve conflicts
vim conflicted-file.js

# Stage resolved files
git add conflicted-file.js

# Complete the merge/rebase
git commit  # for merge strategy
# or
git rebase --continue  # for rebase strategy
```

## CONFIGURATION

Update behavior is controlled by these configuration keys:

### Branch-Level Defaults
```bash
git config gitflow.branch.feature.downstreamStrategy rebase
git config gitflow.branch.release.downstreamStrategy merge
git config gitflow.branch.hotfix.downstreamStrategy merge
```

### Command-Level Overrides
```bash
# Always use rebase for feature updates
git config gitflow.feature.downstreamStrategy rebase

# Always use merge for release updates  
git config gitflow.release.downstreamStrategy merge
```

## STRATEGY RECOMMENDATIONS

### Feature Branches
**Rebase** (recommended) - Creates clean, linear history without unnecessary merge commits

### Release Branches  
**Merge** (recommended) - Preserves the integration history of hotfixes and changes

### Hotfix Branches
**Merge** (recommended) - Preserves context of what changes were integrated

## PARENT RELATIONSHIPS

Each topic branch type updates from its configured parent:

- **feature** → updates from **develop**
- **release** → updates from **main** (to get hotfixes)  
- **hotfix** → updates from **main** (rare, but possible)
- **support** → updates from **main**

## EXIT STATUS

**0**
: Successful update operation

**1**
: Topic branch not found

**2**
: Git operation failed (conflicts, etc.)

**3**
: Invalid branch name or configuration

**4**
: Merge conflicts require manual resolution

**5**
: Branch is not a topic branch

## SEE ALSO

**git-flow**(1), **git-flow-start**(1), **git-flow-finish**(1), **git-flow-config**(1), **git-rebase**(1), **git-merge**(1), **gitflow-config**(5)

## NOTES

- The **--rebase** flag always overrides configured strategy
- Update operations preserve your local commits while integrating parent changes
- Use **git flow update** shorthand when on a topic branch
- Conflicts are resolved the same way as standard Git merge/rebase conflicts
- Regular updates keep topic branches easier to merge when finishing
- Update from parent before finishing long-running branches