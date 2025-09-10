# GIT-FLOW-FINISH(1)

## NAME

git-flow-finish - Complete and merge topic branches

## SYNOPSIS

**git-flow** *topic* **finish** [*name*] [*options*]

**git-flow finish** [*options*]

## DESCRIPTION

Complete a topic branch by merging it to its parent branch according to the configured merge strategy. This command works with any topic branch type (feature, release, hotfix, support, or custom types).

The finish operation performs several steps:
1. Merges the topic branch to its parent branch using the configured upstream strategy
2. Optionally creates and signs tags
3. Updates any child branches that depend on the parent
4. Optionally deletes the topic branch (local and/or remote)

If conflicts occur during merging, the operation can be continued with **--continue** or aborted with **--abort**.

## ARGUMENTS

*topic*
: The topic branch type (feature, release, hotfix, support, or any configured custom type)

*name*
: Name of the topic branch to finish. If omitted when using the shorthand **git-flow finish**, the current branch is used.

## OPTIONS

### Operation Control

**--continue**, **-c**
: Continue the finish operation after resolving merge conflicts

**--abort**, **-a**
: Abort the finish operation and return to the original state

**--force**, **-f**
: Force finish a non-standard branch using this branch type's strategy

### Tag Creation

**--tag**
: Create a tag for the finished branch (overrides configuration)

**--notag**
: Don't create a tag for the finished branch (overrides configuration)

**--sign**
: Sign the tag cryptographically with GPG

**--no-sign**
: Don't sign the tag cryptographically

**--signingkey** *keyid*
: Use the given GPG key for the digital signature

**--message**, **-m** *message*
: Use the given message for the tag

**--messagefile** *file*
: Use contents of the given file as tag message

**--tagname** *name*
: Use the given tag name instead of the default

### Branch Retention

**--keep**
: Keep the topic branch after finishing (don't delete)

**--no-keep**
: Delete the topic branch after finishing (default behavior)

**--keepremote**
: Keep the remote tracking branch after finishing

**--no-keepremote**
: Delete the remote tracking branch after finishing

**--keeplocal**
: Keep the local branch after finishing

**--no-keeplocal**
: Delete the local branch after finishing

**--force-delete**
: Force delete the branch even if not fully merged

**--no-force-delete**
: Don't force delete the branch (default)

## MERGE STRATEGIES

The merge strategy used when finishing is determined by configuration, not command-line flags:

**merge**
: Standard Git merge creating a merge commit (preserves branch history)

**rebase**
: Rebase topic branch onto parent creating linear history

**squash**
: Combine all topic branch commits into a single commit

The strategy is configured via:
1. **Branch defaults**: `gitflow.branch.<type>.upstreamStrategy`
2. **Command overrides**: `gitflow.<type>.finish.merge`

## EXAMPLES

### Basic Usage

Finish current topic branch (shorthand):
```bash
git flow finish
```

Finish specific feature branch:
```bash
git flow feature finish user-authentication
```

Finish release branch with tag:
```bash
git flow release finish 1.2.0 --tag
```

### Handling Conflicts

When conflicts occur during finish:
```bash
# Start finish operation
git flow feature finish my-feature

# Conflicts occur, resolve them manually
git add conflicted-files
git commit

# Continue the finish operation
git flow feature finish my-feature --continue
```

Or abort if needed:
```bash
git flow feature finish my-feature --abort
```

### Tag Management

Create signed tag with custom message:
```bash
git flow release finish 1.2.0 --sign --message "Release version 1.2.0"
```

Use specific GPG key:
```bash
git flow release finish 1.2.0 --sign --signingkey ABC123DEF
```

### Branch Retention

Keep branch for backporting:
```bash
git flow hotfix finish 1.1.1 --keep
```

Clean up both local and remote:
```bash
git flow feature finish my-feature --no-keeplocal --no-keepremote
```

## WORKFLOW BEHAVIOR

### Dual Merge Pattern

For release and hotfix branches, finish typically performs dual merges:
1. **Topic → Parent**: Merge the topic branch to its parent (e.g., release → main)
2. **Parent → Develop**: Merge parent back to develop to synchronize branches

### Child Branch Updates

After merging to the parent, any child branches configured to auto-update will be updated with the new changes.

## CONFIGURATION

Finish behavior is controlled by these configuration keys:

### Branch-Level Defaults
```bash
git config gitflow.branch.<type>.upstreamStrategy merge
git config gitflow.branch.<type>.tag true
```

### Command-Level Overrides
```bash
git config gitflow.<type>.finish.merge rebase
git config gitflow.<type>.finish.sign true
git config gitflow.<type>.finish.signingkey ABC123DEF
```

## EXIT STATUS

**0**
: Successful finish operation

**1**
: Topic branch not found

**2**
: Git operation failed (conflicts, etc.)

**3**
: Invalid branch name or configuration

**4**
: Merge conflicts require manual resolution

**5**
: Tag creation failed

**6**
: GPG signing failed

## SEE ALSO

**git-flow**(1), **git-flow-start**(1), **git-flow-config**(1), **git-flow-update**(1), **gitflow-config**(5)

## NOTES

- Merge strategy is determined by configuration, not command-line flags
- Use **--continue** and **--abort** for conflict resolution
- Tag creation behavior varies by topic branch type configuration
- The **git-flow finish** shorthand automatically detects current topic branch type
- Child branches are automatically updated when their parent changes
- Some topic branch types (like releases and hotfixes) may create tags by default