# PR Review Output Format

## Overview

The PR review produces a structured summary designed for human reviewers. It provides two quick-glance signals (verdict and impact), a brief assessment, action items grouped by severity, a collapsible test case overview, and a collapsible AI fix prompt for automated remediation.

## Structure

### Header (no heading, starts the review)

Two inline status indicators followed by a 1-3 sentence assessment:

- **Verdict**: One of three values:
  - `✅ **Approved**` — Good to merge, no issues.
  - `✅ **Approved with notes**` — Good to merge, but review the notes.
  - `🚫 **Changes requested**` — Do not merge until issues are resolved.
- **Impact**: One of three values:
  - `🔴 **High impact**` — Touches critical paths, wide scope of impact, data loss or security implications if wrong.
  - `🟡 **Medium impact**` — Meaningful logic changes but scoped impact, moderate scope of impact.
  - `🟢 **Low impact**` — Config, docs, refactors, isolated additions, minimal scope of impact.

The assessment summary briefly explains why those two states were chosen. It should describe what the changeset does, how confident the review is, and what the scope of impact is if something goes wrong.

### Action Items (h3)

All findings are grouped under a single `### Action Items` heading, then organized by severity as h4 subsections. Each item is a concise one-liner describing what is wrong without file or line references — inline diff comments carry that detail. This includes test issues, missing tests, and any other findings from all review areas.

- `#### Must fix` — Must fix before merge. Issues that would cause failures, data loss, security vulnerabilities, or violate critical project guidelines.
- `#### Should fix` — Expected to be addressed but won't block merge. Convention violations, gaps in documentation, maintainability concerns.
- `#### Nit` — Optional improvements at the author's discretion. Style, additional tests, minor optimizations.

If a severity subsection has no items, omit it entirely. If there are no action items at all, omit the entire section.

### Test Cases (collapsible, collapsed by default)

A `<details>` block (collapsed by default) providing an overview of tests in the changeset. Contains a table with columns: Test, Coverage, Verdict.

- **Test**: The test function name (e.g. `TestProcessBasic`). Include the file name only when needed to disambiguate (e.g. when multiple test files define similarly-named tests).
- **Coverage**: A brief explanation of the scenario — what setup, action, and expected outcome the test verifies.
- **Verdict**: `✅` for solid tests, `⚠️` with a short note for tests with concerns.

This section is purely an overview of what tests exist in the changeset. Do not include test issues, missing tests, or a summary sentence here — those belong in Action Items.

If the changeset contains no tests, omit this section entirely.

### AI Fix Prompt (collapsible, at the bottom)

A single collapsible `<details>` block at the very end. **Never** place a horizontal rule (`---`) or any other separator before it. The block contains a fenced code block (` ```prompt `) with numbered fix instructions organized by severity category — including nits. Every severity level that has findings in the review must have a corresponding section in the AI fix prompt. Each item provides enough context for an AI agent to locate and fix the issue: file path, what is wrong, and what the fix should look like. Fix instructions should reference the project's conventions and patterns where applicable.

The content must always be inside a fenced code block so the user can copy it in one click. The numbering is continuous across categories so items are easy to reference. The user can copy the prompt, remove items they don't want addressed, and paste it into an AI coding agent.

**If there are no findings (clean PR), omit the AI fix prompt section entirely.** Do not include an empty collapsible block.

## Rules

1. The header has no heading. Action Items uses h3 (`###`). Severity groups within Action Items use h4 (`####`).
2. The action items stay concise and do not include file/line details — inline diff comments provide that context.
3. If a severity subsection has no items, omit it. If the entire Action Items section is empty, omit it.
4. If an area defined in the project's review guidelines was evaluated and had no findings, state what was evaluated and that no issues were found in the assessment summary. The reviewer should know which areas were checked.
5. The verdict and impact are independent dimensions: verdict reflects quality of the change, impact reflects consequence of getting it wrong.
   Verdict is determined by the highest severity present: **Must fix** → Changes requested, **Should fix** (with or without Nit) → Approved with notes, **Nit** only or no findings → Approved.
6. Severity definitions:
   - **Must fix**: Issues that would cause failures, data loss, security vulnerabilities, or violate critical project guidelines. The PR must not merge until these are resolved.
   - **Should fix**: Issues that affect maintainability, violate project conventions, or have gaps that should be addressed. Won't block merge but expected to be resolved.
   - **Nit**: Optional improvements. Style preferences, additional test ideas, minor optimizations. Author's discretion.

## Review Scope

The review evaluates the changeset against the project's review guidelines. The project must provide its own review guidelines that define what to check and what standards to enforce. This format does not prescribe what to look for — it prescribes how to present findings.

Typical review areas include security, code quality & architecture, test coverage, performance, and formal checks (documentation, commit messages, style). The project's review guidelines define the specifics for each area.

All findings from all areas are merged into the Action Items severity subsections. If an area was evaluated and had no findings, mention it in the assessment summary so the reviewer knows it was checked.

## Concrete Examples

### Example 1: Changes requested with must-fix issues

```markdown
🚫 **Changes requested** · 🔴 **High impact**

Modifies the core data processing pipeline and introduces a new
caching layer that multiple modules depend on. Two must-fix issues
prevent merge: a silently ignored error in a critical path and a
flaky test. The scope of impact is high — failures here could cause
silent data loss. Security and performance were reviewed with no
concerns.

### Action Items

#### Must fix
- Silently ignored error in directory creation could fail without indication on restricted systems.
- Test mutates shared state without cleanup, causing ordering-dependent flaky failures.
- No test for concurrent access to the shared cache.

#### Should fix
- Business logic mixed into the request handler instead of the service layer, violating project architecture guidelines.
- No documentation update for the new `--force` flag.
- Error path when upstream service returns 429 is untested.

#### Nit
- Configuration value used directly instead of going through the resolution chain defined in project guidelines.
- Imports not grouped according to project style conventions.

<details><summary>Test Cases</summary>

| Test | Coverage | Verdict |
|------|----------|---------|
| `TestProcessBasic` | Processes standard input with default config, verifies output matches expected transformation | ✅ |
| `TestProcessEmpty` | Processes empty input, checks output is empty without errors | ⚠️ Shallow — only asserts no error, doesn't verify output state |

</details>

<details><summary>🤖 AI fix prompt</summary>

```prompt
## Must fix

1. In `internal/processor/setup.go`, the error returned by `os.MkdirAll`
   is not checked. Wrap the call in proper error handling and return a
   descriptive error using the project's error handling conventions.

2. In `processor_test.go::TestProcess`, the test modifies shared state
   but has no cleanup. Add cleanup that resets the state after each test,
   or refactor to use an isolated fixture per test case.

3. Add a test for concurrent access to the shared cache — spawn multiple
   goroutines and verify no data races or corruption occur.

## Should fix

4. In `handlers/process.go`, move the business logic out of the HTTP
   handler into the service layer, following the project's layered
   architecture pattern.

5. Add documentation for the new `--force` flag. Include a description
   of its behavior, default value, and any interaction with existing flags.

6. Add a test for the error path when the upstream service returns 429.

## Nit

7. In `internal/processor/config.go`, update the configuration resolution
   to follow the project's config precedence chain instead of reading the
   flag value directly.

8. Reorganize imports in affected files to follow the project's import
   grouping conventions.
```

</details>
```

### Example 2: Approved with notes, no must-fix issues

```markdown
✅ **Approved with notes** · 🟡 **Medium impact**

Adds a new machine-readable output mode for the status endpoint. The
implementation follows existing patterns well and test coverage is
solid. A few minor items to consider around documentation and code
structure. Security, performance, and architecture were reviewed with
no concerns.

### Action Items

#### Should fix
- No documentation for the new `--format` flag.

#### Nit
- The output formatter could share an interface with the existing pretty-printer to reduce duplication.
- Commit message on `a1b2c3d` exceeds the project's subject line length limit.
- No test for extremely long entity names in machine-readable output.

<details><summary>Test Cases</summary>

| Test | Coverage | Verdict |
|------|----------|---------|
| `TestMachineOutput` | Runs status with mixed entity states, verifies machine-readable output matches expected format | ✅ |
| `TestMachineOutputEmpty` | Runs status with no entities present, checks for correct empty output | ✅ |
| `TestMachineOutputNested` | Creates nested entity structures, verifies paths appear with correct prefixes | ✅ |

</details>

<details><summary>🤖 AI fix prompt</summary>

```prompt
## Should fix

1. Add documentation for the new `--format` flag. Include a description
   of the output format, its intended use in automation, and how it
   differs from the default human-readable output.

## Nit

2. Extract a shared formatter interface that both the pretty-printer
   and machine-readable formatter implement, reducing duplication in
   header and entry rendering logic.

3. Amend commit `a1b2c3d` to shorten the subject line per the project's
   commit message conventions.

4. Add a test for extremely long entity names to verify the
   machine-readable output handles them correctly.
```

</details>
```

## Follow-up Reviews

When new commits are pushed to a branch that has already been reviewed, the AI reviews only the new commits and produces a follow-up review. The follow-up review must not repeat findings from the previous review — it builds on top of it.

### Follow-up Header

Same verdict and impact format as an initial review, but the assessment summary must state:

- The commit range reviewed (e.g. `a1b2c3d`..`f4e5d6a`)
- How many commits are covered
- A link or reference to the previous review it builds on
- An updated overall assessment reflecting the current state of the PR

The verdict and impact reflect the **current state of the entire PR**, not just the new commits. If the previous review had must-fix issues and the new commits resolve them, the verdict should update accordingly.

### Follow-up Sections

The follow-up review uses three sections to track state changes:

- `### Resolved from previous review` — Items from the prior review that have been addressed. Use strikethrough on the original finding with a brief note on how it was resolved.
- `### Still open from previous review` — Items from the prior review that remain unaddressed. Listed without re-explaining — the previous review has the detail.
- `### New findings` — New issues introduced by the new commits. Uses the same severity subsections (`#### Must fix`, `#### Should fix`, `#### Nit`). If a severity level has no new findings, omit it.

If a section would be empty, omit it entirely. For example, if all previous items are resolved and there are no new findings, only the "Resolved" section appears.

### Follow-up Test Cases

Incremental — only covers tests added or changed in the new commits. Does not repeat the full test table from the previous review. If no tests were added or changed, omit this section.

### Follow-up AI Fix Prompt

Combines still-open items from the previous review and new findings into a single actionable prompt, organized by category (`## Still open`, `## New findings` with severity subsections).

### Example 3: Follow-up review after fixes

```markdown
✅ **Approved** · 🔴 **High impact**

Follow-up review covering commits `a1b2c3d`..`f4e5d6a` (3 commits),
building on [previous review](#link-to-previous-review).

The must-fix issues from the prior review have been resolved. Error
handling is now in place and the flaky test has been fixed. One new
nit introduced by the new commits. All areas reviewed with no
remaining concerns.

### Resolved from previous review
- ~~Silently ignored error in directory creation~~ — Fixed, proper error handling added.
- ~~Test mutates shared state without cleanup~~ — Fixed, uses isolated fixtures now.
- ~~No documentation update for the new `--force` flag~~ — Fixed, documentation added.

### Still open from previous review
- Business logic mixed into the request handler instead of the service layer.

### New findings

#### Nit
- New helper function in `utils/parse.go` duplicates existing logic in `utils/format.go`.

<details><summary>Test Cases</summary>

| Test | Coverage | Verdict |
|------|----------|---------|
| `TestProcessConcurrent` | Spawns multiple goroutines accessing the shared cache simultaneously, verifies no data races or corruption | ✅ |

</details>

<details><summary>🤖 AI fix prompt</summary>

```prompt
## Still open

1. In `handlers/process.go`, move the business logic out of the HTTP
   handler into the service layer, following the project's layered
   architecture pattern.

## New nit

2. In `utils/parse.go`, the new `extractFields` function duplicates
   logic already present in `utils/format.go::parseFields`. Reuse the
   existing function or extract a shared helper.
```

</details>
```
