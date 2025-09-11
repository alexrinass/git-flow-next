#!/bin/bash

# Run tests in the repository
set -e

if [ $# -eq 0 ]; then
    # No arguments - run all tests
    echo "Running all tests..."
    go test -v ./...
elif [ $# -eq 1 ]; then
    # Single argument - run all tests in that package/file
    echo "Running tests in: $1"
    go test -v "$1"
else
    # Multiple arguments - first is package, rest are test names
    package="$1"
    shift
    
    # Build regex pattern for multiple test functions
    pattern=""
    for test in "$@"; do
        if [ -z "$pattern" ]; then
            pattern="$test"
        else
            pattern="$pattern|$test"
        fi
    done
    
    echo "Running tests: $pattern in $package"
    go test -v "$package" -run "$pattern"
fi

echo "Tests completed!"