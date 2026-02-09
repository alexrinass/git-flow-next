# Review Guidelines

This document defines the code review criteria for the git-flow-next project. Use these guidelines for both self-reviews before creating PRs and when reviewing others' code.

## Core Review Areas

Every code review must thoroughly evaluate these four areas:

### 1. Test Coverage

Is the test coverage sufficient and does it handle reasonable error and edge cases?

- [ ] New functionality has corresponding tests
- [ ] Tests cover both success and error paths
- [ ] Edge cases are tested (empty inputs, boundary conditions, invalid states)
- [ ] Tests follow [TESTING_GUIDELINES.md](TESTING_GUIDELINES.md)
- [ ] Tests actually verify correct behavior (not just "no error" assertions)
- [ ] No redundant or overly trivial tests that add noise without value

See [Testing Checklist](#testing-checklist) for detailed criteria.

### 2. Coding Guidelines and Architecture

Does the implementation follow our coding guidelines and architecture?

- [ ] Three-layer command pattern followed (Cobra → Wrapper → Execute)
- [ ] Configuration precedence respected (defaults → git config → flags), noting some options intentionally skip defaults (Layer 2 + 3 only)
- [ ] Git operations use `internal/git/` wrappers
- [ ] Error handling uses custom types from `internal/errors`
- [ ] Code follows [CODING_GUIDELINES.md](CODING_GUIDELINES.md)

See [Architecture Checklist](#architecture-checklist) for detailed criteria.

### 3. Code Quality

Flag any obviously low code quality issues:

- [ ] No ignored errors
- [ ] No unnecessary complexity or over-engineering
- [ ] Clear, descriptive naming
- [ ] Consistent code style with the codebase
- [ ] No dead code or commented-out code
- [ ] Functions are focused and reasonably sized

### 4. Security

Thorough review of changes for security concerns:

- [ ] No command injection vulnerabilities (user input in shell commands)
- [ ] No path traversal vulnerabilities
- [ ] Sensitive data not logged or exposed
- [ ] Input validation at system boundaries
- [ ] No hardcoded credentials or secrets
- [ ] Safe handling of file operations

---

## Detailed Checklists

### Architecture Checklist

- [ ] Three-layer command pattern: Cobra handler → Command wrapper → Execute function
- [ ] Config loaded once at start and passed through
- [ ] All git calls go through `internal/git/repo.go` wrappers
- [ ] Custom error types from `internal/errors` with proper exit codes
- [ ] Three-layer config hierarchy: defaults → git config → flags (flags always win)
- [ ] Pointer types (`*bool`) for optional boolean flags
- [ ] Old config entries cleaned up during rename/delete operations

See [CODING_GUIDELINES.md](CODING_GUIDELINES.md) for detailed patterns.

### Testing Checklist

**Structure:**
- [ ] One test case per function (no table-driven for integration tests)
- [ ] Descriptive test names: `TestFunctionNameWithScenario`
- [ ] Test comments with Steps section (mandatory format)

**Implementation:**
- [ ] Using testutil helpers (not direct exec.Command)
- [ ] Proper setup/cleanup with defer
- [ ] Testing both success and error paths
- [ ] Edge cases covered

**Working Directory:**
- [ ] Using `cmd.Dir`, not `os.Chdir()` where possible
- [ ] If `os.Chdir()` used, proper save/restore with defer

**Best Practices:**
- [ ] Long flag variants only (`--force` not `-f`)
- [ ] Using `git flow init --defaults` for test setup
- [ ] Feature branches for general topic branch testing

See [TESTING_GUIDELINES.md](TESTING_GUIDELINES.md) for comprehensive patterns.

### Code Style Checklist

- [ ] Imports organized: stdlib, third-party, local (with blank lines between)
- [ ] Naming conventions: PascalCase exported, camelCase private
- [ ] Exported functions have documentation comments
- [ ] No ignored errors (`_ = someFunc()` is prohibited)
- [ ] No unnecessary abstractions or premature optimization
- [ ] Changes focused on the task at hand

### Documentation Checklist

- [ ] Manpage updated if command/options changed (`docs/`)
- [ ] Config documentation updated if config changed (`docs/gitflow-config.5.md`)
- [ ] Help text updated for CLI changes
- [ ] README updated for user-facing changes

### Commit Message Checklist

- [ ] Subject line ≤50 characters
- [ ] Type matches change: `feat`, `fix`, `refactor`, `test`, `docs`
- [ ] Imperative mood ("Add feature" not "Added feature")
- [ ] Body explains "what" and "why", wrapped at 72 chars
- [ ] Issue referenced in footer (`Resolves #123`)

See [COMMIT_GUIDELINES.md](COMMIT_GUIDELINES.md) for full standards.

---

## Quality Checks

Run these before submitting:

```bash
go build ./...   # Build check
go test ./...    # Test check
go fmt ./...     # Format check
go vet ./...     # Vet check
```

All checks must pass before code review.

---

## Automated PR Review Format

When performing automated AI pull request reviews, use a concise and direct tone.

### Summary Comment

The PR summary comment should follow this structure:

```markdown
- **Code Quality**: <Brief assessment against coding guidelines and architecture>
- **Test Coverage**:
  - `TestName`: <What it tests> — <assessment>
  - `TestOther`: <What it tests> — <assessment>
  - Missing: <Any critical untested code paths>
  - <Flag duplicates or low-value tests if applicable>
- **Security**: <Brief assessment of concerns, if any>
- **Commit Messages**: <Brief assessment against COMMIT_GUIDELINES.md format>

**Action Items:**
- <Item that needs to be addressed>
- <Item that needs to be addressed>
```

The **Test Coverage** section should list every test introduced or changed by name, assess coverage completeness (especially critical missing cases), flag superficial assertions, and note any duplicate or low-value tests.

### Inline Comments

- Put detailed feedback in inline comments on specific lines
- Keep each comment focused on one issue
- Include fix prompts (see below) for actionable issues

### Fix Prompts

For actionable issues, include a copy-pasteable prompt that the developer can give to an AI assistant to fix the issue. Use the same collapsible format as the review body:

````markdown
**Issue:** [Description of the problem]

<details>
<summary>🤖 AI fix prompt</summary>

```prompt
[Clear, specific instructions that can be copy-pasted to an AI to fix this exact issue.
Include the file path, what's wrong, and what the fix should accomplish.]
```

</details>
````

Example:

````markdown
**Issue:** Missing input validation in `processUser()` function.

<details>
<summary>🤖 AI fix prompt</summary>

```prompt
In src/users.js, add input validation to the processUser() function.
Check that the 'user' parameter is not null/undefined and has required
properties 'id' (number) and 'name' (non-empty string). Throw a
TypeError with a descriptive message if validation fails.
```

</details>
````

---

## Review Process

### Issue Categories

Categorize findings by severity:

**Must fix:**
- Security vulnerabilities
- Test failures or missing critical tests
- Build errors
- Missing error handling
- Architecture violations that affect correctness

**Should fix:**
- Missing documentation
- Inconsistent naming
- Suboptimal patterns
- Minor guideline violations

**Nit:**
- Code clarity improvements
- Additional test cases
- Performance optimizations

### Approval Requirements

- At least one maintainer approval
- All CI checks passing
- Documentation updated
- Tests passing

### Merge Process

1. Squash commits if requested
2. Rebase on latest main branch
3. Ensure clean commit history
4. Delete feature branch after merge

---

## References

- [CODING_GUIDELINES.md](CODING_GUIDELINES.md) - Detailed coding standards
- [TESTING_GUIDELINES.md](TESTING_GUIDELINES.md) - Testing methodology
- [COMMIT_GUIDELINES.md](COMMIT_GUIDELINES.md) - Commit message format
- [CONTRIBUTING.md](CONTRIBUTING.md) - Contribution process
