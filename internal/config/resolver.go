package config

import (
	"fmt"
)

// ResolvedFinishOptions contains all resolved configuration options for the finish command
type ResolvedFinishOptions struct {
	// Tag options
	ShouldTag    bool
	TagName      string
	ShouldSign   bool
	SigningKey   string
	TagMessage   string
	MessageFile  string

	// Branch retention options
	Keep         bool
	KeepRemote   bool
	KeepLocal    bool
	ForceDelete  bool
}

// TagOptions represents command-line tag options
// Note: This should match the TagOptions type in cmd package
type TagOptions struct {
	ShouldTag   *bool
	ShouldSign  *bool
	SigningKey  string
	Message     string
	MessageFile string
	TagName     string
}

// BranchRetentionOptions represents command-line retention options  
// Note: This should match the BranchRetentionOptions type in cmd package
type BranchRetentionOptions struct {
	Keep        *bool
	KeepRemote  *bool
	KeepLocal   *bool
	ForceDelete *bool
}

// ResolveFinishOptions resolves all finish command options using three-layer precedence:
// Layer 1: Branch configuration defaults
// Layer 2: Command-specific git config (gitflow.<branchtype>.finish.*)
// Layer 3: Command-line arguments (highest priority)
func ResolveFinishOptions(cfg *Config, branchType string, branchName string, tagOpts *TagOptions, retentionOpts *BranchRetentionOptions) *ResolvedFinishOptions {
	branchConfig := cfg.Branches[branchType]

	return &ResolvedFinishOptions{
		// Tag resolution
		ShouldTag:    resolveFinishShouldTag(cfg, branchConfig, branchType, tagOpts),
		TagName:      resolveFinishTagName(branchConfig, branchType, branchName, tagOpts),
		ShouldSign:   resolveFinishShouldSign(cfg, branchType, tagOpts),
		SigningKey:   resolveFinishSigningKey(cfg, branchType, tagOpts),
		TagMessage:   resolveFinishTagMessage(branchName, tagOpts),
		MessageFile:  resolveFinishMessageFile(cfg, branchType, tagOpts),

		// Retention resolution
		Keep:         resolveFinishKeep(cfg, branchType, retentionOpts),
		KeepRemote:   resolveFinishKeepRemote(cfg, branchType, retentionOpts),
		KeepLocal:    resolveFinishKeepLocal(cfg, branchType, retentionOpts),
		ForceDelete:  resolveFinishForceDelete(cfg, branchType, retentionOpts),
	}
}

// resolveFinishShouldTag resolves whether to create a tag
func resolveFinishShouldTag(cfg *Config, branchConfig BranchConfig, branchType string, tagOpts *TagOptions) bool {
	// Layer 1: Branch configuration default
	shouldTag := branchConfig.Tag

	// Layer 2: Command-specific config (notag inverts the default)
	if notag := getCommandConfigBool(cfg, fmt.Sprintf("gitflow.%s.finish.notag", branchType)); notag {
		shouldTag = false
	}

	// Layer 3: Command-line flags override config
	if tagOpts != nil && tagOpts.ShouldTag != nil {
		shouldTag = *tagOpts.ShouldTag
	}

	return shouldTag
}

// resolveFinishTagName resolves the tag name to use
func resolveFinishTagName(branchConfig BranchConfig, branchType string, branchName string, tagOpts *TagOptions) string {
	// Layer 1: Default is branch name with prefix from branch config
	tagName := branchName
	if branchConfig.TagPrefix != "" {
		tagName = branchConfig.TagPrefix + branchName
	}

	// Layer 2: No command-specific config for tag name

	// Layer 3: Command-line custom tag name overrides default
	if tagOpts != nil && tagOpts.TagName != "" {
		tagName = tagOpts.TagName
	}

	return tagName
}

// resolveFinishShouldSign resolves whether to sign the tag
func resolveFinishShouldSign(cfg *Config, branchType string, tagOpts *TagOptions) bool {
	// Layer 1: Default is not signing
	shouldSign := false

	// Layer 2: Check command-specific signing config
	if sign := getCommandConfigBool(cfg, fmt.Sprintf("gitflow.%s.finish.sign", branchType)); sign {
		shouldSign = true
	}

	// Layer 3: Command-line signing flags override config
	if tagOpts != nil && tagOpts.ShouldSign != nil {
		shouldSign = *tagOpts.ShouldSign
	}

	return shouldSign
}

// resolveFinishSigningKey resolves the signing key to use
func resolveFinishSigningKey(cfg *Config, branchType string, tagOpts *TagOptions) string {
	// Layer 1: Default is empty
	signingKey := ""

	// Layer 2: Check command-specific signing key
	if key := getCommandConfigString(cfg, fmt.Sprintf("gitflow.%s.finish.signingkey", branchType)); key != "" {
		signingKey = key
	}

	// Layer 3: Command-line signing key overrides config
	if tagOpts != nil && tagOpts.SigningKey != "" {
		signingKey = tagOpts.SigningKey
	}

	return signingKey
}

// resolveFinishTagMessage resolves the tag message
func resolveFinishTagMessage(branchName string, tagOpts *TagOptions) string {
	// Layer 1: Default message
	message := fmt.Sprintf("Tagging version %s", branchName)

	// Layer 2: No command-specific config for tag message

	// Layer 3: Command-line message overrides default
	if tagOpts != nil && tagOpts.Message != "" {
		message = tagOpts.Message
	}

	return message
}

// resolveFinishMessageFile resolves the message file path
func resolveFinishMessageFile(cfg *Config, branchType string, tagOpts *TagOptions) string {
	// Layer 1: Default is empty
	messageFile := ""

	// Layer 2: Check command-specific message file config
	if file := getCommandConfigString(cfg, fmt.Sprintf("gitflow.%s.finish.messagefile", branchType)); file != "" {
		messageFile = file
	}

	// Layer 3: Command-line message file overrides config
	if tagOpts != nil && tagOpts.MessageFile != "" {
		messageFile = tagOpts.MessageFile
	}

	return messageFile
}

// resolveFinishKeep resolves whether to keep the branch
func resolveFinishKeep(cfg *Config, branchType string, retentionOpts *BranchRetentionOptions) bool {
	// Layer 1: Default is to delete (not keep)
	keep := false

	// Layer 2: Check command-specific config
	if keepConfig := getCommandConfigBool(cfg, fmt.Sprintf("gitflow.%s.finish.keep", branchType)); keepConfig {
		keep = true
	}

	// Layer 3: Command-line flags override config
	if retentionOpts != nil && retentionOpts.Keep != nil {
		keep = *retentionOpts.Keep
	}

	return keep
}

// resolveFinishKeepRemote resolves whether to keep the remote branch
func resolveFinishKeepRemote(cfg *Config, branchType string, retentionOpts *BranchRetentionOptions) bool {
	// Layer 1: Default is to delete remote
	keepRemote := false

	// Layer 2: Check command-specific config
	if keepRemoteConfig := getCommandConfigBool(cfg, fmt.Sprintf("gitflow.%s.finish.keepremote", branchType)); keepRemoteConfig {
		keepRemote = true
	}

	// Layer 3: Command-line flags override config
	if retentionOpts != nil && retentionOpts.KeepRemote != nil {
		keepRemote = *retentionOpts.KeepRemote
	}

	return keepRemote
}

// resolveFinishKeepLocal resolves whether to keep the local branch
func resolveFinishKeepLocal(cfg *Config, branchType string, retentionOpts *BranchRetentionOptions) bool {
	// Layer 1: Default is to delete local
	keepLocal := false

	// Layer 2: Check command-specific config
	if keepLocalConfig := getCommandConfigBool(cfg, fmt.Sprintf("gitflow.%s.finish.keeplocal", branchType)); keepLocalConfig {
		keepLocal = true
	}

	// Layer 3: Command-line flags override config
	if retentionOpts != nil && retentionOpts.KeepLocal != nil {
		keepLocal = *retentionOpts.KeepLocal
	}

	return keepLocal
}

// resolveFinishForceDelete resolves whether to force delete the branch
func resolveFinishForceDelete(cfg *Config, branchType string, retentionOpts *BranchRetentionOptions) bool {
	// Layer 1: Default is not to force delete
	forceDelete := false

	// Layer 2: Check command-specific config
	if forceDeleteConfig := getCommandConfigBool(cfg, fmt.Sprintf("gitflow.%s.finish.force-delete", branchType)); forceDeleteConfig {
		forceDelete = true
	}

	// Layer 3: Command-line flags override config
	if retentionOpts != nil && retentionOpts.ForceDelete != nil {
		forceDelete = *retentionOpts.ForceDelete
	}

	return forceDelete
}

// getCommandConfigBool gets a boolean config value from preloaded config
func getCommandConfigBool(cfg *Config, configKey string) bool {
	value, exists := cfg.CommandConfig[configKey]
	return exists && value == "true"
}

// getCommandConfigString gets a string config value from preloaded config  
func getCommandConfigString(cfg *Config, configKey string) string {
	value, exists := cfg.CommandConfig[configKey]
	if !exists {
		return ""
	}
	return value
}