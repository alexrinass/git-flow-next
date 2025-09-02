# git-flow-avh Command Reference

This directory contains the complete command reference for git-flow-avh, organized by command type for easy lookup during compatibility comparison with git-flow-next.

## Command Files

### Core Commands
- [**init**](init.md) - Initialize repository for git-flow usage
- [**config**](config.md) - Configuration management
- [**general**](general.md) - General shorthand commands (work on current branch)

### Branch Type Commands
- [**bugfix**](bugfix.md) - Bugfix branch operations
- [**feature**](feature.md) - Feature branch operations  
- [**hotfix**](hotfix.md) - Hotfix branch operations
- [**release**](release.md) - Release branch operations
- [**support**](support.md) - Support branch operations

## Common Operations

Each branch type supports these standard operations (where applicable):

| Operation | Description |
|-----------|-------------|
| `list` | List existing branches of this type |
| `start` | Start a new branch |
| `finish` | Finish/merge a branch |
| `publish` | Push branch to remote origin |
| `track` | Start tracking a remote branch |
| `diff` | Show changes compared to parent |
| `rebase` | Rebase branch on parent |
| `checkout` | Switch to branch |
| `pull` | Pull branch from remote |
| `delete` | Delete branch |
| `rename` | Rename branch |

## Common Flags

### Universal Flags
- `-h,--[no]help` - Show help
- `--showcommands` - Show git commands while executing

### Finish Operation Flags
- `-F,--[no]fetch` - Fetch from origin before operation
- `-k,--[no]keep` - Keep branch after finishing
- `--[no]keeplocal` - Keep local branch
- `--[no]keepremote` - Keep remote branch
- `-D,--[no]force_delete` - Force delete branch

### Merge Strategy Flags
- `-r,--[no]rebase` - Rebase before merging
- `-S,--[no]squash` - Squash commits during merge
- `--no-ff` - Never fast-forward during merge
- `-p,--[no]preserve-merges` - Preserve merges while rebasing

### Tag Creation Flags (release/hotfix)
- `-s,--[no]sign` - Sign tag cryptographically
- `-u,--signingkey` - Use specific GPG key for signing
- `-m,--message` - Use custom tag message
- `-f,--messagefile` - Use tag message from file
- `-n,--[no]notag` - Don't create tag
- `-T,--tagname` - Use custom tag name

### Remote Operations
- `-p,--[no]push` - Push to origin after finish
- `--[no]pushproduction` - Push production branch
- `--[no]pushdevelop` - Push develop branch
- `-b,--[no]nobackmerge` - Don't back-merge into develop

## Usage Examples

```bash
# Initialize with defaults
git flow init -d

# Start a feature
git flow feature start my-feature

# Finish feature with rebase and keep branch
git flow feature finish -r -k my-feature

# Start hotfix and finish with tag
git flow hotfix start 1.0.1
git flow hotfix finish -s 1.0.1

# Release with custom tag message
git flow release finish -m "Release v2.0.0" 2.0.0
```

## Compatibility Notes

This reference is maintained for compatibility comparison with git-flow-next. Key differences:

1. **Unified Architecture**: git-flow-next uses unified topic branch commands vs. separate feature/hotfix/release commands
2. **Configuration-Driven**: Behavior determined by git config rather than hardcoded branch types  
3. **Modern Workflow Support**: Flexible branch hierarchies and merge strategies

Refer to individual command files for complete option details and usage patterns.