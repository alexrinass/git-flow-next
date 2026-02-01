# Git-Flow-AVH Compatibility Analysis

This document provides a comprehensive compatibility analysis between git-flow-avh and git-flow-next, covering both configuration options and command implementations.

## Configuration Format Comparison

### git-flow-avh Format
```bash
# Branch names (legacy)
gitflow.branch.master
gitflow.branch.develop

# Branch prefixes (legacy)
gitflow.prefix.feature
gitflow.prefix.release
gitflow.prefix.hotfix
gitflow.prefix.bugfix
gitflow.prefix.support
gitflow.prefix.versiontag

# Command options (matches our Layer 2)
gitflow.<branchtype>.<command>.<option>
```

### git-flow-next Format
```bash
# System settings
gitflow.origin
gitflow.version
gitflow.initialized

# Branch type configuration (our Layer 1)
gitflow.branch.<branchtype>.type
gitflow.branch.<branchtype>.parent
gitflow.branch.<branchtype>.prefix
gitflow.branch.<branchtype>.tag
gitflow.branch.<branchtype>.tagprefix
gitflow.branch.<branchtype>.upstreamStrategy
gitflow.branch.<branchtype>.downstreamStrategy

# Command-specific overrides (our Layer 2)
gitflow.<branchtype>.<command>.<option>
```

## Supported Options Analysis

### ✅ Fully Compatible Options
These git-flow-avh options are fully supported in git-flow-next:

| AVH Option | git-flow-next Equivalent | Status |
|------------|--------------------------|---------|
| `gitflow.origin` | `gitflow.origin` | ✅ Direct match |
| `gitflow.feature.start.fetch` | `gitflow.feature.start.fetch` | ✅ Direct match |
| `gitflow.feature.finish.fetch` | `gitflow.feature.finish.fetch` | ✅ Direct match |
| `gitflow.release.finish.sign` | `gitflow.release.finish.sign` | ✅ Direct match |
| `gitflow.release.finish.signingkey` | `gitflow.release.finish.signingkey` | ✅ Direct match |
| `gitflow.hotfix.finish.sign` | `gitflow.hotfix.finish.sign` | ✅ Direct match |
| `gitflow.release.finish.keep` | `gitflow.release.finish.keep` | ✅ Direct match |
| `gitflow.feature.finish.keep` | `gitflow.feature.finish.keep` | ✅ Direct match |
| `gitflow.*.finish.no-ff` | `gitflow.*.finish.no-ff` | ✅ Direct match |
| `gitflow.*.finish.preserve-merges` | `gitflow.*.finish.preserve-merges` | ✅ Direct match |

### 🔄 Requires Translation (Auto-Import)
These AVH options need translation but are automatically handled during migration:

| AVH Option | git-flow-next Translation | Notes |
|------------|--------------------------|-------|
| `gitflow.branch.master` | `gitflow.branch.main` | Branch name modernization |
| `gitflow.prefix.feature` | `gitflow.branch.feature.prefix` | New hierarchical format |
| `gitflow.prefix.release` | `gitflow.branch.release.prefix` | New hierarchical format |
| `gitflow.prefix.hotfix` | `gitflow.branch.hotfix.prefix` | New hierarchical format |
| `gitflow.prefix.support` | `gitflow.branch.support.prefix` | New hierarchical format |
| `gitflow.prefix.versiontag` | `gitflow.branch.*.tagprefix` | Applied to all tag-creating types |

### 🔄 Strategy Translation Required
These AVH boolean flags need translation to our strategy-based system. The resolver handles `notag` at runtime (Layer 2), but the import does not auto-convert these to git-flow-next's preferred format:

| AVH Option | git-flow-next Equivalent | Translation Rule |
|------------|--------------------------|------------------|
| `gitflow.feature.finish.rebase=true` | `gitflow.feature.finish.merge=rebase` | Boolean to strategy |
| `gitflow.feature.finish.squash=true` | `gitflow.feature.finish.merge=squash` | Boolean to strategy |
| `gitflow.release.finish.squash=true` | `gitflow.release.finish.merge=squash` | Boolean to strategy |
| `gitflow.release.finish.notag=true` | `gitflow.branch.release.tag=false` | Inverted boolean (works at runtime via resolver, not migrated) |
| `gitflow.hotfix.finish.notag=true` | `gitflow.branch.hotfix.tag=false` | Inverted boolean (works at runtime via resolver, not migrated) |

### ⚠️ Missing Options in git-flow-next
These git-flow-avh options are **not currently supported** in git-flow-next:

| Missing Option | Description | Impact |
|----------------|-------------|--------|
| `gitflow.allowdirty` | Allow operations with dirty working tree | 🔴 High |
| `gitflow.*.finish.push` | Auto-push after finishing | 🔴 High |
| `gitflow.*.finish.pushproduction` | Push to production branch after finish | 🟡 Medium |
| `gitflow.*.finish.pushdevelop` | Push to develop branch after finish | 🟡 Medium |
| `gitflow.*.finish.nobackmerge` | Skip back-merge to develop | 🟡 Medium |
| `gitflow.*.finish.ff-master` | Fast-forward merge to master | 🟡 Medium |

### 🆕 git-flow-next Exclusive Features
These features are **only available** in git-flow-next:

| Feature | Description | Benefit |
|---------|-------------|---------|
| Branch type system | `gitflow.branch.*.type=base\|topic` | Flexible branch taxonomy |
| Parent relationships | `gitflow.branch.develop.parent=main` | Automatic change propagation |
| Auto-update | `gitflow.branch.develop.autoUpdate=true` | Keeps branches synchronized |
| Downstream strategies | `gitflow.*.downstreamStrategy` | Control update merge behavior |
| Universal branch properties | Common properties across all types | Consistent configuration |
| Custom branch types | User-defined branch types | Workflow customization |
| Config management | `git flow config` subcommands | Add, edit, rename, delete branch types |

## Migration Requirements

### Automatic Import (Implemented)
✅ Our `git flow init` automatically handles these conversions:
- `gitflow.branch.master` → `gitflow.branch.main`
- `gitflow.prefix.*` → `gitflow.branch.*.prefix`
- `gitflow.prefix.versiontag` → `gitflow.branch.*.tagprefix`
- Command options are preserved as-is

### Runtime Compatibility (No Migration Needed)
🔄 These AVH options work at runtime without conversion:
- `gitflow.*.finish.notag` — handled by the config resolver at Layer 2

### Manual Translation Required
🔄 These require user intervention or enhanced import logic:
- Boolean strategy flags (`rebase=true`, `squash=true`) → strategy configuration

### Missing Feature Implementation
❌ These require new feature development:

#### 1. Allow Dirty Working Tree
```bash
# AVH
gitflow.allowdirty=true

# Needed in git-flow-next
gitflow.allowdirty=true  # Add system-level support
```

#### 2. Auto-Push on Finish
```bash
# AVH
gitflow.release.finish.push=true

# Needed in git-flow-next
gitflow.release.finish.push=true  # Add to finish command options
```

---

## Command Compatibility Analysis

### Command Status Overview

| Command | AVH | git-flow-next | Status |
|---------|-----|---------------|--------|
| `init` | ✅ | ✅ | ✅ Full parity (`-d/--defaults`, `-f/--force`) |
| `start` | ✅ | ✅ | ✅ Full parity (`-F/--fetch`, `<base>` parameter) |
| `finish` | ✅ | ✅ | 🟡 Most flags supported (see details) |
| `publish` | ✅ | ✅ | ✅ Fully implemented |
| `track` | ✅ | ✅ | ✅ Fully implemented |
| `list` | ✅ | ✅ | 🟡 Missing `-v/--verbose` |
| `checkout` | ✅ | ✅ | ✅ Complete |
| `delete` | ✅ | ✅ | ✅ Complete parity |
| `rename` | ✅ | ✅ | ✅ Complete |
| `diff` | ✅ | ❌ | ❌ Not implemented |
| `pull` | ✅ | ❌ | ❌ Not implemented |
| `config` | ✅ | ✅ | 🟡 Different approach (add/edit/rename/delete vs set/base) |

### General Shorthand Commands
git-flow-next implements shorthand commands that work on the current branch:
- ✅ **`git flow finish`** — Finish current branch
- ✅ **`git flow delete`** — Delete current branch
- ✅ **`git flow rebase`** — Rebase current branch (via `update --rebase`)
- ✅ **`git flow update`** — Update current branch
- ✅ **`git flow rename`** — Rename current branch
- ✅ **`git flow publish`** — Publish current branch

### Command-by-Command Option Analysis

#### ✅ `init` Command
| Flag | AVH | git-flow-next |
|------|-----|---------------|
| `-d/--defaults` | ✅ | ✅ |
| `-f/--force` | ✅ | ✅ |
| `--no-create-branches` | ❌ | ✅ (git-flow-next exclusive) |
| `--local/--global/--system/--file` | ✅ | ❌ |

#### ✅ `start` Command
| Flag | AVH | git-flow-next |
|------|-----|---------------|
| `-F/--fetch` | ✅ | ✅ |
| `<base>` parameter | ✅ | ✅ |
| `--showcommands` | ✅ | ❌ |

#### 🟡 `finish` Command
| Flag | AVH | git-flow-next |
|------|-----|---------------|
| `-F/--fetch` | ✅ | ✅ |
| `-r/--rebase` | ✅ | ✅ (via merge strategy) |
| `-S/--squash` | ✅ | ✅ (via merge strategy) |
| `-k/--keep` | ✅ | ✅ |
| `-D/--force-delete` | ✅ | ✅ |
| `--no-ff` / `--ff` | ✅ | ✅ |
| `-p/--preserve-merges` | ✅ | ✅ |
| `-s/--sign` | ✅ | ✅ |
| `--signingkey` | ✅ | ✅ |
| `-T/--tagname` | ✅ | ✅ |
| `-m/--message` | ✅ | ✅ |
| `--notag` | ✅ | ✅ |
| `--continue/--abort` | ❌ | ✅ (git-flow-next exclusive) |
| `-p/--push` | ✅ | ❌ |
| `--pushproduction` | ✅ | ❌ |
| `--pushdevelop` | ✅ | ❌ |
| `-b/--nobackmerge` | ✅ | ❌ |
| `--ff-master` | ✅ | ❌ |
| `--showcommands` | ✅ | ❌ |

#### Other Commands
| Command | Missing from git-flow-next |
|---------|---------------------------|
| `list` | `-v/--verbose` |
| `checkout` | `--showcommands` (partial — present but not on all commands) |
| `delete` | ✅ Complete parity |
| `rename` | ✅ Complete parity |

---

## Compatibility Scores

### Configuration Compatibility: 85%

- **Core Features**: 100% compatible (branch management, merge strategies)
- **Configuration Translation**: 90% automatic (prefix mapping, branch names)
- **Command Options**: 75% compatible (most finish flags supported, missing push)
- **Advanced Features**: 60% compatible (preserve-merges and no-ff now supported)

### Command Compatibility: 80%

- **Core Commands**: 95% compatible (start, finish, list, delete, checkout, rename)
- **Shorthand Commands**: 100% compatible (all shorthands implemented)
- **Remote Operations**: 70% compatible (publish and track implemented; pull missing)
- **Advanced Operations**: 40% compatible (diff and full config parity still missing)

### Combined Compatibility Score: 82%

- **Core Workflow**: 95% compatible (start → work → finish → delete)
- **Team Collaboration**: 75% compatible (publish and track work; push-on-finish missing)
- **Configuration**: 85% compatible
- **Commands**: 80% compatible

---

## Remaining Implementation Roadmap

### Phase 1: Critical Gaps
**Configuration:**
1. **`gitflow.allowdirty`** — Allow dirty working tree operations
2. **`gitflow.*.finish.push`** — Auto-push after finishing operations

**Commands:**
1. **`diff` command** — Show changes compared to parent branch

### Phase 2: Important Workflow Features
**Configuration:**
1. **Enhanced strategy translation** — Auto-convert boolean flags to strategies during import

**Commands:**
1. **`--verbose` on list** — Detailed branch listing
2. **Missing finish push flags** — `--pushproduction/--pushdevelop`
3. **Universal `--showcommands` flag** — Debug/learning aid across all commands

### Phase 3: Advanced Features
**Commands:**
1. **`pull` command** — Pull from remote
2. **`--nobackmerge` / `--ff-master`** — Specialized finish merge options
3. **Config file location flags** (`--local/--global/--system/--file`) in init

---

*Last updated: 2026-02-01*
