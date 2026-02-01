# Contributing to git-flow-next

Thank you for your interest in contributing to git-flow-next! This document provides guidelines and instructions for contributing.

## Code of Conduct

By participating in this project, you agree to:
- Be respectful and inclusive of differing viewpoints and experiences
- Use welcoming and inclusive language
- Accept constructive criticism gracefully
- Focus on what is best for the community
- Show empathy towards other community members

## How to Contribute

### Reporting Bugs

1. **Check Existing Issues** - Search the issue tracker to avoid duplicates
2. **Create a Clear Report** Including:
   - A clear, descriptive title
   - Detailed steps to reproduce the bug
   - Expected behavior
   - Actual behavior
   - Your environment (OS, git version, etc.)
   - Any relevant logs or error messages

### Suggesting Enhancements

1. **Check Existing Suggestions** - Search issues and discussions
2. **Provide Context** - Explain why this enhancement would be useful
3. **Consider Scope** - Is it generally useful or specific to your use case?
4. **Describe Implementation** - If possible, outline how it might be implemented

### Pull Requests

1. **Fork the Repository**
2. **Create a Branch**
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/your-bug-fix
   ```

3. **Commit Guidelines**

   Follow the commit message format defined in [COMMIT_GUIDELINES.md](COMMIT_GUIDELINES.md). Key points:
   - Use a `<type>:` prefix (`feat`, `fix`, `refactor`, `test`, `docs`, etc.)
   - Imperative mood ("Add feature" not "Added feature")
   - Subject line ≤50 characters
   - Reference issues in the footer (`Resolves #123`)

   Example:
     ```
     feat: Add support for custom merge strategies

     Implements configurable merge strategies per branch type.

     Resolves #123
     ```

4. **Code Style**
   - Follow existing code style and formatting
   - Add comments for complex logic
   - Update documentation for API changes
   - Include tests for new functionality
   - Ensure all tests pass

5. **Submit Pull Request**
   - Provide a clear title and description
   - Link related issues
   - Include any necessary documentation updates
   - Add notes about testing performed

### Development Setup

1. **Prerequisites**
   - Go 1.19 or later
   - Git 2.25 or later

2. **Local Development**
   ```bash
   # Clone your fork
   git clone https://github.com/YOUR_USERNAME/git-flow-next.git
   cd git-flow-next

   # Add upstream remote
   git remote add upstream https://github.com/gittower/git-flow-next.git

   # Create branch
   git checkout -b your-feature

   # Build
   go build -o git-flow

   # Run tests
   go test ./...
   ```

### Testing

1. **Unit Tests**
   - Add tests for new functionality
   - Update tests for modified code
   - Ensure all tests pass locally

2. **Integration Tests**
   - Test your changes with different Git workflows
   - Verify behavior with merge conflicts
   - Test with different branch configurations

## License and Copyright

- All contributions must be licensed under the BSD 2-Clause license
- You retain copyright on your contributions
- By submitting a pull request, you agree to license your contributions under the same BSD 2-Clause license

## Additional Guidelines

### Documentation

1. **Code Comments**
   - Add comments for complex logic
   - Update function documentation
   - Include examples for new features

2. **User Documentation**
   - Update README.md for user-facing changes
   - Add or update command documentation
   - Include examples for new functionality

### Quality Standards

1. **Code Quality**
   - Follow Go best practices
   - Use meaningful variable and function names
   - Keep functions focused and concise
   - Handle errors appropriately

2. **Testing Requirements**
   - Minimum 80% test coverage for new code
   - Include both positive and negative test cases
   - Test edge cases and error conditions

3. **Performance Considerations**
   - Consider impact on large repositories
   - Avoid unnecessary Git operations
   - Profile code for potential bottlenecks

### Development Workflow

For a detailed overview of our end-to-end development process — from issue analysis through planning, implementation, and PR creation — see [DEV_WORKFLOW.md](DEV_WORKFLOW.md).

### Review Process

All code changes are reviewed according to [REVIEW_GUIDELINES.md](REVIEW_GUIDELINES.md), which covers:
- Code review checklist (architecture, style, testing, documentation)
- Issue categories (blocking, warnings, suggestions)
- Approval requirements
- Merge process

Please review your changes against these guidelines before submitting a pull request.

## Getting Help

- Create an issue for questions
- Join our community discussions
- Read our documentation
- Contact maintainers for guidance

## Recognition

Contributors will be:
- Listed in CONTRIBUTORS.md
- Mentioned in release notes for significant contributions
- Recognized in project documentation where appropriate

Thank you for contributing to git-flow-next! 