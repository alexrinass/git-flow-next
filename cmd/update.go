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
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update [branch]",
	Short: "Update a branch with changes from its parent branch",
	Long: `Update a branch with changes from its parent branch.
This command will update the specified branch (or current branch if none specified)
with changes from its parent branch using the configured downstream strategy (merge or rebase).
If merge conflicts occur, they will be handled according to the configured merge state handling.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var branchName string
		if len(args) > 0 {
			branchName = args[0]
		}
		useRebase, _ := cmd.Flags().GetBool("rebase")
		if err := executeUpdate("", branchName, useRebase); err != nil {
			var exitCode errors.ExitCode
			if flowErr, ok := err.(errors.Error); ok {
				exitCode = flowErr.ExitCode()
			} else {
				exitCode = errors.ExitCodeGitError
			}
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(int(exitCode))
		}
		return nil
	},
}

// createUpdateCommand creates a new update command for a specific branch type
func createUpdateCommand(branchType string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [name]",
		Short: fmt.Sprintf("Update a %s branch with changes from its parent branch", branchType),
		Long: fmt.Sprintf(`Update a %s branch with changes from its parent branch.
This command will update the specified %s branch (or current branch if none specified)
with changes from its parent branch using the configured downstream strategy (merge or rebase).
If merge conflicts occur, they will be handled according to the configured merge state handling.`, branchType, branchType),
		RunE: func(cmd *cobra.Command, args []string) error {
			var name string
			if len(args) > 0 {
				name = args[0]
			}
			useRebase, _ := cmd.Flags().GetBool("rebase")
			if err := executeUpdate(branchType, name, useRebase); err != nil {
				var exitCode errors.ExitCode
				if flowErr, ok := err.(errors.Error); ok {
					exitCode = flowErr.ExitCode()
				} else {
					exitCode = errors.ExitCodeGitError
				}
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(int(exitCode))
			}
			return nil
		},
	}
	
	// Add --rebase flag to the command
	cmd.Flags().Bool("rebase", false, "Force rebase strategy instead of configured strategy")
	
	return cmd
}

func init() {
	// Add --rebase flag to the root update command
	updateCmd.Flags().Bool("rebase", false, "Force rebase strategy instead of configured strategy")
	rootCmd.AddCommand(updateCmd)
}

// AddUpdateCommand adds the update command to a topic branch command
func AddUpdateCommand(parentCmd *cobra.Command) {
	parentCmd.AddCommand(createUpdateCommand(parentCmd.Name()))
}

// executeUpdate updates a branch with changes from its parent branch
func executeUpdate(branchType string, name string, useRebase bool) error {
	// Validate that git-flow is initialized
	initialized, err := config.IsInitialized()
	if err != nil {
		return &errors.GitError{Operation: "check if git-flow is initialized", Err: err}
	}
	if !initialized {
		return &errors.NotInitializedError{}
	}

	// Get configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return &errors.GitError{Operation: "load configuration", Err: err}
	}

	var branchName string
	if branchType != "" {
		// If branch type is specified, construct full branch name
		if name == "" {
			// If no name provided, try to get current branch and verify it's of the correct type
			currentBranch, err := git.GetCurrentBranch()
			if err != nil {
				return &errors.GitError{Operation: "get current branch", Err: err}
			}
			branchConfig, ok := cfg.Branches[branchType]
			if !ok {
				return &errors.InvalidBranchTypeError{BranchType: branchType}
			}
			if !strings.HasPrefix(currentBranch, branchConfig.Prefix) {
				return &errors.GitError{Operation: "validate current branch", Err: fmt.Errorf("current branch is not a %s branch", branchType)}
			}
			branchName = currentBranch
		} else {
			// Construct full branch name from type and name
			branchConfig, ok := cfg.Branches[branchType]
			if !ok {
				return &errors.InvalidBranchTypeError{BranchType: branchType}
			}
			branchName = branchConfig.Prefix + name
		}
	} else {
		// No branch type specified, use provided branch name or current branch
		if name == "" {
			currentBranch, err := git.GetCurrentBranch()
			if err != nil {
				return &errors.GitError{Operation: "get current branch", Err: err}
			}
			branchName = currentBranch
		} else {
			branchName = name
		}
	}

	// Check if branch exists
	if err := git.BranchExists(branchName); err != nil {
		return &errors.BranchNotFoundError{BranchName: branchName}
	}

	// Get parent branch
	parentBranch, err := update.GetParentBranch(cfg, branchName)
	if err != nil {
		return err
	}

	// Check if parent branch exists
	if err := git.BranchExists(parentBranch); err != nil {
		return &errors.BranchNotFoundError{BranchName: parentBranch}
	}

	// Get branch configuration for merge strategy
	var strategy string
	for branchKey, bc := range cfg.Branches {
		if bc.Type == string(config.BranchTypeBase) && branchKey == branchName {
			strategy = bc.DownstreamStrategy
			break
		}
		if bc.Type == string(config.BranchTypeTopic) && bc.Prefix != "" && strings.HasPrefix(branchName, bc.Prefix) {
			strategy = bc.DownstreamStrategy
			break
		}
	}

	if strategy == "" {
		strategy = "merge" // Default to merge if no strategy configured
	}

	// Override strategy if --rebase flag is set
	if useRebase {
		strategy = "rebase"
	}

	// Create merge state
	state := &mergestate.MergeState{
		Action:         "update",
		BranchName:     branchName,
		ParentBranch:   parentBranch,
		MergeStrategy:  strategy,
		CurrentStep:    "merge",
		FullBranchName: branchName,
	}

	// Update the branch using shared logic
	return update.UpdateBranchFromParent(branchName, parentBranch, strategy, true, state)
}
