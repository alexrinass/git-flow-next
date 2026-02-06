package util

import "strings"

// ExpandMessagePlaceholders expands placeholders in merge commit messages.
// Supported placeholders:
//
//	%b - branch name (e.g., feature/my-feature)
//	%B - full refname (e.g., refs/heads/feature/my-feature)
//	%p - parent branch name (e.g., develop)
//	%P - full parent refname (e.g., refs/heads/develop)
//	%% - literal percent sign
func ExpandMessagePlaceholders(message, branch, parent string) string {
	// Use strings.NewReplacer for efficient multi-replacement
	replacer := strings.NewReplacer(
		"%%", "\x00", // Temporarily escape %%
		"%b", branch,
		"%B", "refs/heads/"+branch,
		"%p", parent,
		"%P", "refs/heads/"+parent,
	)
	result := replacer.Replace(message)
	return strings.ReplaceAll(result, "\x00", "%")
}
