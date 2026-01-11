package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/gittower/git-flow-next/internal/config"
	"github.com/gittower/git-flow-next/internal/errors"
	"github.com/gittower/git-flow-next/internal/git"
)

// TrackCommand is the implementation of the track command for topic branches
func TrackCommand(branchType string, name string) {
	if err := track(branchType, name); err != nil {
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

// track performs the actual tracking branch creation logic
func track(branchType string, name string) error {
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
		return &errors.EmptyBranchNameError{}
	}

	// Get configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return &errors.GitError{Operation: "load configuration", Err: err}
	}

	// Get branch configuration
	branchConfig, ok := cfg.Branches[branchType]
	if !ok {
		return &errors.InvalidBranchTypeError{BranchType: branchType}
	}

	// Construct full branch name
	fullBranchName := name
	if branchConfig.Prefix != "" && !strings.HasPrefix(name, branchConfig.Prefix) {
		fullBranchName = branchConfig.Prefix + name
	}

	// Check if branch already exists locally
	if err := git.BranchExists(fullBranchName); err == nil {
		return &errors.BranchExistsError{BranchName: fullBranchName}
	}

	// Determine remote (from config or default to "origin")
	remote := cfg.Remote
	if remote == "" {
		remote = "origin"
	}

	// Fetch from remote to ensure we have latest refs
	fmt.Printf("Fetching from '%s'...\n", remote)
	if err := git.Fetch(remote); err != nil {
		return &errors.GitError{
			Operation: fmt.Sprintf("fetch from remote '%s'", remote),
			Err:       err,
		}
	}

	// Check if branch exists on remote
	if !git.RemoteBranchExists(remote, fullBranchName) {
		return &errors.RemoteBranchNotFoundError{
			Remote:     remote,
			BranchName: fullBranchName,
		}
	}

	// Create local tracking branch
	fmt.Printf("Setting up tracking branch for '%s'...\n", fullBranchName)
	if err := git.CreateTrackingBranch(fullBranchName, remote, fullBranchName); err != nil {
		return &errors.GitError{
			Operation: fmt.Sprintf("create tracking branch '%s'", fullBranchName),
			Err:       err,
		}
	}

	fmt.Printf("Successfully created tracking branch '%s' from '%s/%s'\n",
		fullBranchName, remote, fullBranchName)
	return nil
}
