// Package hooks provides filter and hook execution for git-flow operations.
//
// Filters transform input/output values (e.g., version numbers, tag messages).
// Hooks execute scripts before/after operations (e.g., pre-flow-feature-start).
//
// All scripts are located in .git/hooks/ following these patterns:
//   - Filters: filter-flow-{type}-{action}-{target}
//   - Hooks: {pre,post}-flow-{type}-{action}
package hooks

// FilterType represents the type of filter script.
type FilterType string

const (
	// FilterVersionReleaseStart filters the version for release start.
	FilterVersionReleaseStart FilterType = "filter-flow-release-start-version"

	// FilterVersionHotfixStart filters the version for hotfix start.
	FilterVersionHotfixStart FilterType = "filter-flow-hotfix-start-version"

	// FilterTagMessageReleaseFinish filters the tag message for release finish.
	FilterTagMessageReleaseFinish FilterType = "filter-flow-release-finish-tag-message"

	// FilterTagMessageHotfixFinish filters the tag message for hotfix finish.
	FilterTagMessageHotfixFinish FilterType = "filter-flow-hotfix-finish-tag-message"
)

// FilterContext contains data passed to filters.
type FilterContext struct {
	BranchType string // The type of branch (e.g., "release", "hotfix")
	BranchName string // The short name of the branch
	Version    string // The version number (for release/hotfix)
	TagMessage string // The tag message to filter
	BaseBranch string // The base/parent branch
}

// HookPhase represents pre or post execution.
type HookPhase string

const (
	// HookPre runs before the operation.
	HookPre HookPhase = "pre"

	// HookPost runs after the operation.
	HookPost HookPhase = "post"
)

// HookAction represents the git-flow action being performed.
type HookAction string

const (
	// HookActionStart is the start action.
	HookActionStart HookAction = "start"

	// HookActionFinish is the finish action.
	HookActionFinish HookAction = "finish"

	// HookActionPublish is the publish action.
	HookActionPublish HookAction = "publish"

	// HookActionTrack is the track action.
	HookActionTrack HookAction = "track"

	// HookActionDelete is the delete action.
	HookActionDelete HookAction = "delete"
)

// HookContext contains data passed to hooks via environment variables.
type HookContext struct {
	BranchType string // The type of branch (e.g., "feature", "release")
	BranchName string // The short name of the branch
	FullBranch string // The full branch name with prefix
	BaseBranch string // The base/parent branch
	Origin     string // The remote name
	Version    string // The version (for release/hotfix)
	ExitCode   int    // For post-hooks: exit code of the operation
}

// HookResult contains the result of hook execution.
type HookResult struct {
	Executed bool   // Whether the hook was found and executed
	ExitCode int    // Exit code of the hook script
	Output   string // Combined stdout/stderr output
	Error    error  // Any error that occurred during execution
}
