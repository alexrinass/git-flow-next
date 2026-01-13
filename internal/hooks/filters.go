package hooks

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RunVersionFilter executes a version filter and returns the modified version.
// If the filter does not exist or is not executable, the original version is returned.
// If the filter exits with a non-zero status, an error is returned.
func RunVersionFilter(gitDir string, filterType FilterType, version string) (string, error) {
	scriptPath := filepath.Join(gitDir, "hooks", string(filterType))

	// Check if filter exists and is executable
	if !isExecutable(scriptPath) {
		return version, nil
	}

	// Execute the filter with version as argument
	result, err := runFilter(scriptPath, version, nil)
	if err != nil {
		return "", fmt.Errorf("version filter '%s' failed: %w", filterType, err)
	}

	// If the filter returned empty output, use the original version
	if result == "" {
		return version, nil
	}

	return result, nil
}

// RunTagMessageFilter executes a tag message filter and returns the modified message.
// If the filter does not exist or is not executable, the original message is returned.
// If the filter exits with a non-zero status, an error is returned.
func RunTagMessageFilter(gitDir string, filterType FilterType, ctx FilterContext) (string, error) {
	scriptPath := filepath.Join(gitDir, "hooks", string(filterType))

	// Check if filter exists and is executable
	if !isExecutable(scriptPath) {
		return ctx.TagMessage, nil
	}

	// Build environment variables for the filter
	env := buildFilterEnv(ctx)

	// Execute the filter with version and message as arguments
	// The filter receives: $1 = version, $2 = message
	result, err := runFilterWithArgs(scriptPath, []string{ctx.Version, ctx.TagMessage}, env)
	if err != nil {
		return "", fmt.Errorf("tag message filter '%s' failed: %w", filterType, err)
	}

	// If the filter returned empty output, use the original message
	if result == "" {
		return ctx.TagMessage, nil
	}

	return result, nil
}

// isExecutable checks if a file exists and is executable.
func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	// Check if file is executable (any execute bit set)
	return info.Mode()&0111 != 0
}

// buildFilterEnv builds environment variables for filter execution.
func buildFilterEnv(ctx FilterContext) []string {
	env := os.Environ()
	env = append(env,
		fmt.Sprintf("BRANCH_TYPE=%s", ctx.BranchType),
		fmt.Sprintf("BRANCH_NAME=%s", ctx.BranchName),
		fmt.Sprintf("BASE_BRANCH=%s", ctx.BaseBranch),
	)
	if ctx.Version != "" {
		env = append(env, fmt.Sprintf("VERSION=%s", ctx.Version))
	}
	return env
}

// runFilter executes a filter script with input on stdin.
func runFilter(scriptPath string, input string, env []string) (string, error) {
	cmd := exec.Command(scriptPath, input)

	if env != nil {
		cmd.Env = env
	}

	// Set working directory to repository root
	cmd.Dir = filepath.Dir(filepath.Dir(scriptPath))

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("exit code %d: %s", exitErr.ExitCode(), string(exitErr.Stderr))
		}
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

// runFilterWithArgs executes a filter script with arguments.
func runFilterWithArgs(scriptPath string, args []string, env []string) (string, error) {
	cmd := exec.Command(scriptPath, args...)

	if env != nil {
		cmd.Env = env
	}

	// Set working directory to repository root
	cmd.Dir = filepath.Dir(filepath.Dir(scriptPath))

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("exit code %d: %s", exitErr.ExitCode(), string(exitErr.Stderr))
		}
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}
