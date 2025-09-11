package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/gittower/git-flow-next/internal/config"
	"github.com/gittower/git-flow-next/internal/errors"
	"github.com/gittower/git-flow-next/internal/git"
	"github.com/gittower/git-flow-next/internal/mergestate"
	"github.com/gittower/git-flow-next/internal/update"
)

// Step constants
const (
	stepMerge          = "merge"
	stepCreateTag      = "create_tag"
	stepUpdateChildren = "update_children"
	stepDeleteBranch   = "delete_branch"
)

// Strategy constants
const (
	strategyRebase = "rebase"
	strategySquash = "squash"
	strategyMerge  = "merge"
)

// =============================================================================
// PUBLIC ENTRY POINTS
// =============================================================================

// FinishCommand is the implementation of the finish command for topic branches
func FinishCommand(branchType string, name string, continueOp bool, abortOp bool, force bool, tagOptions *config.TagOptions, retentionOptions *config.BranchRetentionOptions, mergeOptions *config.MergeStrategyOptions) {
	if err := executeFinish(branchType, name, continueOp, abortOp, force, tagOptions, retentionOptions, mergeOptions); err != nil {
		var exitCode errors.ExitCode
		if flowErr, ok := err.(errors.Error); ok {
			exitCode = flowErr.ExitCode()
		} else {
			exitCode = errors.ExitCodeGitError
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(int(exitCode))
	}
}

// =============================================================================
// MAIN EXECUTION FLOW
// =============================================================================

// executeFinish performs the actual branch finishing logic and returns any errors
func executeFinish(branchType string, name string, continueOp bool, abortOp bool, force bool, tagOptions *config.TagOptions, retentionOptions *config.BranchRetentionOptions, mergeOptions *config.MergeStrategyOptions) error {
	// Get configuration early
	cfg, err := config.LoadConfig()
	if err != nil {
		return &errors.GitError{Operation: "load configuration", Err: err}
	}

	// Get branch configuration
	branchConfig, ok := cfg.Branches[branchType]
	if !ok {
		return &errors.InvalidBranchTypeError{BranchType: branchType}
	}

	// Check if there's a merge in progress
	if mergestate.IsMergeInProgress() {
		state, err := mergestate.LoadMergeState()
		if err != nil {
			return &errors.GitError{Operation: "load merge state", Err: err}
		}

		// Get the branch config for the state's branch type
		stateBranchConfig, ok := cfg.Branches[state.BranchType]
		if !ok {
			return &errors.InvalidBranchTypeError{BranchType: state.BranchType}
		}

		if abortOp {
			return handleAbort(state)
		}

		if continueOp {
			// Resolve options for continue operation
			resolvedOptions := config.ResolveFinishOptions(cfg, state.BranchType, state.BranchName, tagOptions, retentionOptions, mergeOptions)
			return handleContinue(cfg, state, stateBranchConfig, resolvedOptions)
		}

		return &errors.MergeInProgressError{BranchName: state.FullBranchName}
	}

	// Don't allow continue or abort if no merge is in progress
	if continueOp || abortOp {
		return &errors.NoMergeInProgressError{}
	}

	// Resolve branch name (try with and without prefix)
	resolvedName, err := resolveBranchName(name, branchConfig)
	if err != nil {
		return err
	}
	name = resolvedName

	// If the branch exists but doesn't have the expected prefix
	if !strings.HasPrefix(name, branchConfig.Prefix) {
		if !force {
			// Get the short name for tag creation
			shortName := name
			if strings.Contains(name, "/") {
				parts := strings.Split(name, "/")
				shortName = parts[len(parts)-1]
			}

			// Prompt user for confirmation
			fmt.Printf("Warning: Branch '%s' is not a standard %s branch (missing prefix '%s').\n", name, branchType, branchConfig.Prefix)
			fmt.Printf("Finishing this branch will:\n")
			fmt.Printf("1. Merge it into '%s' using the %s strategy\n", branchConfig.Parent, branchConfig.UpstreamStrategy)

			// Resolve options early for confirmation dialog
			resolvedOptions := config.ResolveFinishOptions(cfg, branchType, shortName, tagOptions, retentionOptions, mergeOptions)

			if resolvedOptions.ShouldTag {
				fmt.Printf("2. Create a tag '%s'\n", resolvedOptions.TagName)
			}

			fmt.Printf("3. Delete the branch after successful merge\n\n")
			fmt.Printf("Do you want to continue? [y/N]: ")

			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" {
				return fmt.Errorf("operation cancelled by user")
			}
		}
	}

	// Regular finish command flow
	return finishBranch(cfg, branchType, name, branchConfig, tagOptions, retentionOptions, mergeOptions)
}

func finishBranch(cfg *config.Config, branchType string, name string, branchConfig config.BranchConfig, tagOptions *config.TagOptions, retentionOptions *config.BranchRetentionOptions, mergeOptions *config.MergeStrategyOptions) error {
	// Validate that git-flow is initialized
	initialized, err := config.IsInitialized()
	if err != nil {
		return &errors.GitError{Operation: "check if git-flow is initialized", Err: err}
	}
	if !initialized {
		return &errors.NotInitializedError{}
	}

	// Validate inputs
	if name == "" {
		return &errors.InvalidBranchNameError{BranchName: name}
	}

	// Get the short name by removing the prefix if it exists
	shortName := name
	if strings.HasPrefix(name, branchConfig.Prefix) {
		shortName = strings.TrimPrefix(name, branchConfig.Prefix)
	} else if strings.Contains(name, "/") {
		// For non-standard branches, use the last part after the slash
		parts := strings.Split(name, "/")
		shortName = parts[len(parts)-1]
	}

	// Check if branch exists
	if err := git.BranchExists(name); err != nil {
		return &errors.BranchNotFoundError{BranchName: name}
	}

	// Get target branch (always the parent branch for finishing)
	targetBranch := branchConfig.Parent

	// Check if target branch exists
	if err := git.BranchExists(targetBranch); err != nil {
		return &errors.BranchNotFoundError{BranchName: targetBranch}
	}

	// Find child base branches that need to be updated
	childBranches := []string{}
	for branchName, branch := range cfg.Branches {
		if branch.Type == string(config.BranchTypeBase) && branch.Parent == targetBranch {
			fmt.Printf("Found child base branch '%s' to update\n", branchName)
			childBranches = append(childBranches, branchName)
		}
	}

	// Save merge state before starting
	state := &mergestate.MergeState{
		Action:          "finish",
		BranchType:      branchType,
		BranchName:      shortName,
		CurrentStep:     stepMerge,
		ParentBranch:    targetBranch,
		MergeStrategy:   branchConfig.UpstreamStrategy,
		FullBranchName:  name,
		ChildBranches:   childBranches,
		UpdatedBranches: []string{},
	}
	if err := mergestate.SaveMergeState(state); err != nil {
		return &errors.GitError{Operation: "save merge state", Err: err}
	}

	// Resolve all options once at the beginning
	resolvedOptions := config.ResolveFinishOptions(cfg, branchType, shortName, tagOptions, retentionOptions, mergeOptions)

	return executeSteps(cfg, state, branchConfig, resolvedOptions)
}

// =============================================================================
// STATE MACHINE AND CONTROL FLOW
// =============================================================================

// executeSteps runs the state machine for the finish operation
func executeSteps(cfg *config.Config, state *mergestate.MergeState, branchConfig config.BranchConfig, resolvedOptions *config.ResolvedFinishOptions) error {
	for {
		var err error
		switch state.CurrentStep {
		case stepMerge:
			err = handleMergeStep(cfg, state, branchConfig, resolvedOptions)
		case stepCreateTag:
			err = handleCreateTagStep(state, resolvedOptions)
		case stepUpdateChildren:
			err = handleUpdateChildrenStep(cfg, state, branchConfig, resolvedOptions)
		case stepDeleteBranch:
			return handleDeleteBranchStep(state, resolvedOptions) // Final step
		default:
			return &errors.GitError{Operation: fmt.Sprintf("unknown step '%s'", state.CurrentStep), Err: nil}
		}

		if err != nil {
			return err
		}
	}
}

func handleContinue(cfg *config.Config, state *mergestate.MergeState, branchConfig config.BranchConfig, resolvedOptions *config.ResolvedFinishOptions) error {
	// For merge step continuation, check if conflicts are resolved
	if state.CurrentStep == stepMerge {
		if git.HasConflicts() {
			return &errors.UnresolvedConflictsError{}
		}
		// Move to next step since merge conflicts are resolved
		state.CurrentStep = stepCreateTag
		if err := mergestate.SaveMergeState(state); err != nil {
			return &errors.GitError{Operation: "save merge state", Err: err}
		}
	}

	return executeSteps(cfg, state, branchConfig, resolvedOptions)
}

func handleAbort(state *mergestate.MergeState) error {
	// Abort the merge based on strategy
	var err error
	switch state.MergeStrategy {
	case strategyMerge:
		err = git.MergeAbort()
	case strategyRebase:
		err = git.RebaseAbort()
	default:
		err = git.MergeAbort() // Default to merge abort
	}

	if err != nil {
		return &errors.GitError{Operation: "abort merge", Err: err}
	}

	// Checkout the original branch
	if err := git.Checkout(state.FullBranchName); err != nil {
		return &errors.GitError{Operation: fmt.Sprintf("checkout original branch '%s'", state.FullBranchName), Err: err}
	}

	// Clear the merge state
	if err := mergestate.ClearMergeState(); err != nil {
		return &errors.GitError{Operation: "clear merge state", Err: err}
	}

	return nil
}

// =============================================================================
// STEP HANDLERS (Called by state machine)
// =============================================================================

// handleMergeStep handles the merge step of the finish operation
func handleMergeStep(cfg *config.Config, state *mergestate.MergeState, branchConfig config.BranchConfig, resolvedOptions *config.ResolvedFinishOptions) error {
	// Checkout target branch
	err := git.Checkout(state.ParentBranch)
	if err != nil {
		return &errors.GitError{Operation: fmt.Sprintf("checkout target branch '%s'", state.ParentBranch), Err: err}
	}
	fmt.Printf("Switched to branch '%s'\n", state.ParentBranch)

	// Perform merge based on resolved strategy
	fmt.Printf("Merging using strategy: %v\n", resolvedOptions.MergeStrategy)
	var mergeErr error
	switch resolvedOptions.MergeStrategy {
	case strategyRebase:
		fmt.Printf("Rebase strategy selected\n")
		// For rebase, we need to:
		// 1. Stay on feature branch
		err = git.Checkout(state.FullBranchName)
		if err != nil {
			return &errors.GitError{Operation: "checkout feature branch for rebase", Err: err}
		}
		// 2. Rebase onto target branch with options
		mergeErr = git.RebaseWithOptions(state.ParentBranch, resolvedOptions.PreserveMerges)
		if mergeErr == nil {
			// 3. If rebase succeeds, checkout target and merge
			err = git.Checkout(state.ParentBranch)
			if err != nil {
				return &errors.GitError{Operation: "checkout target branch after rebase", Err: err}
			}
			// Use NoFastForward option for the final merge
			mergeErr = git.MergeWithOptions(state.FullBranchName, resolvedOptions.NoFastForward)
		}
	case strategySquash:
		squashMessage := fmt.Sprintf("Squashed commit of branch '%s'", state.FullBranchName)
		mergeErr = git.MergeSquashWithMessage(state.FullBranchName, squashMessage)
	case strategyMerge:
		mergeErr = git.MergeWithOptions(state.FullBranchName, resolvedOptions.NoFastForward)
	default:
		return &errors.GitError{Operation: fmt.Sprintf("unknown merge strategy: %s", resolvedOptions.MergeStrategy), Err: nil}
	}

	if mergeErr != nil {
		if strings.Contains(mergeErr.Error(), "conflict") {
			// Save state before returning conflict error
			state.CurrentStep = stepMerge
			if err := mergestate.SaveMergeState(state); err != nil {
				return &errors.GitError{Operation: "save merge state", Err: err}
			}

			msg := fmt.Sprintf("Merge conflicts detected. Resolve conflicts and run 'git flow %s finish --continue %s'\n", state.BranchType, state.BranchName)
			msg += fmt.Sprintf("To abort the merge, run 'git flow %s finish --abort %s'", state.BranchType, state.BranchName)
			fmt.Println(msg)
			return &errors.UnresolvedConflictsError{}
		}
		return &errors.GitError{Operation: "merge branch", Err: mergeErr}
	}

	// Move to next step (tag creation)
	state.CurrentStep = stepCreateTag
	if err := mergestate.SaveMergeState(state); err != nil {
		return &errors.GitError{Operation: "save merge state", Err: err}
	}

	return nil
}

// handleCreateTagStep handles the tag creation step
func handleCreateTagStep(state *mergestate.MergeState, resolvedOptions *config.ResolvedFinishOptions) error {
	if resolvedOptions.ShouldTag {
		if err := createTagForBranchResolved(state, resolvedOptions); err != nil {
			return err
		}
	}

	// Move to next step
	state.CurrentStep = stepUpdateChildren
	if err := mergestate.SaveMergeState(state); err != nil {
		return &errors.GitError{Operation: "save merge state", Err: err}
	}
	return nil
}

// handleUpdateChildrenStep handles updating child base branches
func handleUpdateChildrenStep(cfg *config.Config, state *mergestate.MergeState, branchConfig config.BranchConfig, resolvedOptions *config.ResolvedFinishOptions) error {
	// Find next child branch to update
	nextBranch := findNextBranchToUpdate(state)

	// If no more branches to update, move to final step
	if nextBranch == "" {
		state.CurrentStep = stepDeleteBranch
		if err := mergestate.SaveMergeState(state); err != nil {
			return &errors.GitError{Operation: "save merge state", Err: err}
		}
		return nil
	}

	// Update the next child branch
	if err := updateChildBranch(cfg, nextBranch, state); err != nil {
		return err
	}

	// Mark this branch as updated
	state.UpdatedBranches = append(state.UpdatedBranches, nextBranch)
	if err := mergestate.SaveMergeState(state); err != nil {
		return &errors.GitError{Operation: "save merge state", Err: err}
	}

	// Continue with next branch
	return nil
}

// handleDeleteBranchStep handles branch deletion
func handleDeleteBranchStep(state *mergestate.MergeState, resolvedOptions *config.ResolvedFinishOptions) error {
	// Ensure we're on the parent branch before deletion
	if err := git.Checkout(state.ParentBranch); err != nil {
		return &errors.GitError{Operation: fmt.Sprintf("checkout parent branch '%s'", state.ParentBranch), Err: err}
	}

	// Apply keep logic: if keep is set, it overrides individual settings
	keepRemote := resolvedOptions.KeepRemote
	keepLocal := resolvedOptions.KeepLocal
	if resolvedOptions.Keep {
		keepRemote = true
		keepLocal = true
	}

	// Delete branches based on settings
	if err := deleteBranchesIfNeeded(state, keepRemote, keepLocal, resolvedOptions.ForceDelete); err != nil {
		return err
	}

	// Clean up base branch configuration if branch was deleted
	if !keepLocal {
		configKey := fmt.Sprintf("gitflow.branch.%s.base", state.FullBranchName)
		if err := git.UnsetConfig(configKey); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to clean up base config: %v\n", err)
		}
	}

	// Clear the merge state
	if err := mergestate.ClearMergeState(); err != nil {
		return &errors.GitError{Operation: "clear merge state", Err: err}
	}

	fmt.Printf("Successfully finished branch '%s' and updated %d child base branches\n", state.FullBranchName, len(state.UpdatedBranches))
	return nil
}

// =============================================================================
// HELPER FUNCTIONS (Called by step handlers and main flow)
// =============================================================================

// resolveBranchName tries to find the branch name with and without prefix
func resolveBranchName(name string, branchConfig config.BranchConfig) (string, error) {
	// Try name as-is first
	if err := git.BranchExists(name); err == nil {
		return name, nil
	}

	// If not found as-is, try with prefix
	if !strings.HasPrefix(name, branchConfig.Prefix) {
		fullName := branchConfig.Prefix + name
		if err := git.BranchExists(fullName); err == nil {
			return fullName, nil
		}
	}

	return "", &errors.BranchNotFoundError{BranchName: name}
}

// createTagForBranchResolved creates a tag using resolved options
func createTagForBranchResolved(state *mergestate.MergeState, options *config.ResolvedFinishOptions) error {
	// Determine if we should use message file
	useMessageFile := options.MessageFile != ""

	// Create the tag using the git module
	gitTagOptions := &git.TagOptions{
		Message:     options.TagMessage,
		MessageFile: options.MessageFile,
		Sign:        options.ShouldSign,
		SigningKey:  options.SigningKey,
	}

	// Use MessageFile if specified, otherwise use Message
	if useMessageFile {
		gitTagOptions.Message = "" // Clear message since we're using file
	} else {
		gitTagOptions.MessageFile = "" // Clear file since we're using message
	}

	if err := git.CreateTag(options.TagName, gitTagOptions); err != nil {
		return &errors.GitError{Operation: fmt.Sprintf("create tag '%s'", options.TagName), Err: err}
	}
	fmt.Printf("Created tag '%s'\n", options.TagName)
	return nil
}

// findNextBranchToUpdate finds the next child branch that needs updating
func findNextBranchToUpdate(state *mergestate.MergeState) string {
	for _, branch := range state.ChildBranches {
		alreadyUpdated := false
		for _, updated := range state.UpdatedBranches {
			if branch == updated {
				alreadyUpdated = true
				break
			}
		}
		if !alreadyUpdated {
			return branch
		}
	}
	return ""
}

// updateChildBranch updates a single child branch
func updateChildBranch(cfg *config.Config, branchName string, state *mergestate.MergeState) error {
	fmt.Printf("Updating child base branch '%s' from '%s'...\n", branchName, state.ParentBranch)

	// Get merge strategy for this child branch from provided config
	childBranchConfig, ok := cfg.Branches[branchName]
	if !ok {
		return &errors.GitError{Operation: fmt.Sprintf("get config for branch '%s'", branchName), Err: fmt.Errorf("branch config not found")}
	}

	// Use the shared update logic
	err := update.UpdateBranchFromParent(branchName, state.ParentBranch, childBranchConfig.DownstreamStrategy, true, state)
	if err != nil {
		if _, ok := err.(*errors.UnresolvedConflictsError); ok {
			msg := fmt.Sprintf("Merge conflicts detected while updating base branch '%s'. Resolve conflicts and run 'git flow %s finish --continue %s'\n", branchName, state.BranchType, state.BranchName)
			msg += fmt.Sprintf("To abort the merge, run 'git flow %s finish --abort %s'", state.BranchType, state.BranchName)
			fmt.Println(msg)
			return err
		}
		return err
	}

	return nil
}

// deleteBranchesIfNeeded deletes branches based on retention settings
func deleteBranchesIfNeeded(state *mergestate.MergeState, keepRemote, keepLocal, forceDelete bool) error {
	// Delete remote branch if not keeping it and if remote branch exists
	if !keepRemote {
		// Only attempt to delete if the remote branch actually exists
		if git.RemoteBranchExists("origin", state.FullBranchName) {
			remoteBranch := fmt.Sprintf("origin/%s", state.FullBranchName)
			if err := git.DeleteRemoteBranch("origin", state.FullBranchName); err != nil {
				return &errors.GitError{Operation: fmt.Sprintf("delete remote branch '%s'", remoteBranch), Err: err}
			}
		}
	}

	// Delete local branch if not keeping it
	if !keepLocal {
		if err := git.DeleteBranch(state.FullBranchName, forceDelete); err != nil {
			return &errors.GitError{Operation: fmt.Sprintf("delete branch '%s'", state.FullBranchName), Err: err}
		}
	}

	return nil
}
