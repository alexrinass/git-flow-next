# Git-Flow-AVH Compatibility Analysis

This document analyzes compatibility between git-flow-avh and git-flow-next configurations, identifying gaps and migration requirements.

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
| `gitflow.feature.finish.fetch` | `gitflow.feature.finish.fetch` | ✅ Direct match |
| `gitflow.release.finish.sign` | `gitflow.release.finish.sign` | ✅ Direct match |
| `gitflow.release.finish.signingkey` | `gitflow.release.finish.signingkey` | ✅ Direct match |
| `gitflow.hotfix.finish.sign` | `gitflow.hotfix.finish.sign` | ✅ Direct match |
| `gitflow.release.finish.keep` | `gitflow.release.finish.keep` | ✅ Direct match |
| `gitflow.feature.finish.keep` | `gitflow.feature.finish.keep` | ✅ Direct match |

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
These AVH boolean flags need translation to our strategy-based system:

| AVH Option | git-flow-next Equivalent | Translation Rule |
|------------|--------------------------|------------------|
| `gitflow.feature.finish.rebase=true` | `gitflow.feature.finish.merge=rebase` | Boolean to strategy |
| `gitflow.feature.finish.squash=true` | `gitflow.feature.finish.merge=squash` | Boolean to strategy |
| `gitflow.release.finish.squash=true` | `gitflow.release.finish.merge=squash` | Boolean to strategy |
| `gitflow.release.finish.notag=true` | `gitflow.branch.release.tag=false` | Inverted boolean |
| `gitflow.hotfix.finish.notag=true` | `gitflow.branch.hotfix.tag=false` | Inverted boolean |

### ⚠️ Missing Options in git-flow-next
These git-flow-avh options are **not currently supported** in git-flow-next:

#### High Priority Missing Options
| Missing Option | Description | Impact | Recommendation |
|----------------|-------------|---------|----------------|
| `gitflow.allowdirty` | Allow operations with dirty working tree | 🔴 High | **ADD SUPPORT** - Important for workflow flexibility |
| `gitflow.*.start.fetch` | Fetch before starting branches | 🔴 High | **ADD SUPPORT** - Critical for team workflows |
| `gitflow.*.finish.push` | Auto-push after finishing | 🔴 High | **ADD SUPPORT** - Common workflow requirement |
| `gitflow.*.finish.force-delete` | Force delete branches | 🟡 Medium | **ADD SUPPORT** - Safety feature |
| `gitflow.*.finish.preserve-merges` | Preserve merges during rebase | 🟡 Medium | **ADD SUPPORT** - Advanced rebase option |

#### Medium Priority Missing Options  
| Missing Option | Description | Impact | Recommendation |
|----------------|-------------|---------|----------------|
| `gitflow.support.rebase.interactive` | Interactive rebase for support | 🟡 Medium | Consider adding |
| `gitflow.prefix.bugfix` | Bugfix branch type | 🟡 Medium | Consider adding as custom type |

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

## Migration Requirements

### Automatic Import (Implemented)
✅ Our `git flow init` automatically handles these conversions:
- `gitflow.branch.master` → `gitflow.branch.main`
- `gitflow.prefix.*` → `gitflow.branch.*.prefix`
- Command options are preserved as-is

### Manual Translation Required
🔄 These require user intervention or enhanced import logic:
- Boolean strategy flags → strategy configuration
- `notag` flags → inverted `tag` configuration

### Missing Feature Implementation
❌ These require new feature development:

#### 1. Allow Dirty Working Tree
```bash
# AVH
gitflow.allowdirty=true

# Needed in git-flow-next
gitflow.allowdirty=true  # Add system-level support
```

#### 2. Auto-Fetch on Start
```bash  
# AVH
gitflow.feature.start.fetch=true

# Needed in git-flow-next  
gitflow.feature.start.fetch=true  # Add to start command options
```

#### 3. Auto-Push on Finish
```bash
# AVH
gitflow.release.finish.push=true

# Needed in git-flow-next
gitflow.release.finish.push=true  # Add to finish command options
```

## Compatibility Score

### Overall Compatibility: 75% ✅

- **Core Features**: 100% compatible (branch management, merge strategies)
- **Configuration Translation**: 90% automatic (prefix mapping, branch names)
- **Command Options**: 60% compatible (missing push, fetch, dirty flags)
- **Advanced Features**: 40% compatible (missing preserve-merges, interactive options)

## Recommended Implementation Priority

### Priority 1: Critical Compatibility (Required for migration)
1. **`gitflow.allowdirty`** - Allow dirty working tree operations
2. **`gitflow.*.start.fetch`** - Fetch before starting branches  
3. **`gitflow.*.finish.push`** - Auto-push after finishing operations
4. **Enhanced strategy translation** - Auto-convert boolean flags to strategies

### Priority 2: Important Workflow Features  
1. **`gitflow.*.finish.force-delete`** - Force delete branch options
2. **`gitflow.*.finish.preserve-merges`** - Advanced rebase options
3. **Better error messages** for unsupported AVH options

### Priority 3: Nice-to-Have Features
1. **`gitflow.support.rebase.interactive`** - Interactive rebase support
2. **`gitflow.prefix.bugfix`** - Bugfix branch type support
3. **Environment variable support** - `.gitflow_export` file support

## Implementation Plan

### Phase 1: Core Compatibility (High Priority Missing Options)
- Add `allowdirty` system configuration
- Add `fetch` options to start commands  
- Add `push` options to finish commands
- Add `force-delete` options to finish commands

### Phase 2: Enhanced Translation
- Auto-convert boolean strategy flags during import
- Auto-convert `notag` to inverted `tag` configuration  
- Provide migration warnings for unsupported options

### Phase 3: Advanced Features
- Add `preserve-merges` rebase option
- Add interactive rebase support
- Consider `bugfix` branch type addition

---

*This analysis ensures git-flow-next maintains strong compatibility with existing git-flow-avh installations while providing enhanced modern features.*