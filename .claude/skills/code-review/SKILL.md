---
name: code-review
description: PR code review mechanics via GitHub API
allowed-tools: Bash, Read, Grep, Glob
---

# Code Review Skill

PR code review mechanics: create reviews via GitHub API, post inline comments, respond to mentions.

## Input Context

The workflow provides:
- `EVENT`: GitHub event type (`pull_request`, `issue_comment`, `pull_request_review_comment`)
- `REPO`: Repository in `owner/repo` format
- `PR_NUMBER`: Pull request number
- `COMMENT`: Comment body (empty for `pull_request` events)

## Execution Context Detection

```bash
echo "GITHUB_ACTIONS=${GITHUB_ACTIONS:-not_set} CI=${CI:-not_set}"
```

**Execution modes:**
- **CI mode**: `GITHUB_ACTIONS=true` or `CI=true` — Posts reviews to GitHub
- **Local mode**: Neither set — Auto-enables dry-run, outputs review to console

When running locally, inform the user:
> "Running in local mode — dry-run enabled automatically. Review output will be displayed but not posted to GitHub."

---

## Intent Detection

```
If EVENT == "pull_request":
  → INTENT = "review"

If EVENT == "issue_comment" or "pull_request_review_comment":
  If COMMENT contains "@claude review" or "@claude re-review":
    → INTENT = "review"
  Else:
    → INTENT = "respond"
```

---

## Loading Guidelines

Before reviewing or responding, check for project-specific guidelines:

1. **Read CLAUDE.md** (if exists) - may contain review instructions
2. **Follow references** - if CLAUDE.md mentions other files (e.g., "see `REVIEW_GUIDELINES.md`"), read them
3. **Check common locations**: `REVIEW_GUIDELINES.md`, `.github/CONTRIBUTING.md`, `docs/code-standards.md`

Follow `REVIEW_GUIDELINES.md` for review criteria, summary format, and inline comment format.

---

## Action: Review

### Pre-Review Check

Determine what needs reviewing:

```bash
# Get Claude's previous reviews
LAST_REVIEW=$(gh api repos/$REPO/pulls/$PR_NUMBER/reviews \
  --jq '[.[] | select(.user.login == "github-actions[bot]" or .user.login == "claude[bot]")] | sort_by(.submitted_at) | last')

# Get current HEAD commit
CURRENT_HEAD=$(gh pr view $PR_NUMBER --json headRefOid --jq '.headRefOid')

# Get commit SHA from last review
LAST_REVIEW_COMMIT=$(echo "$LAST_REVIEW" | jq -r '.commit_id // empty')
```

**Decision:**
- No previous review → Full review (`gh pr diff $PR_NUMBER`)
- Same commit → Post "No new changes to review"
- New commits in history → Incremental review
- Last commit not in history → Force-push handling

### Incremental Review

```bash
git diff $LAST_REVIEW_COMMIT..$CURRENT_HEAD
```

Only comment on lines changed in new commits.

### Force-Push Handling

```bash
git merge-base --is-ancestor $LAST_REVIEW_COMMIT $CURRENT_HEAD
# Exit 0 = normal push, non-zero = force-push
```

When force-push detected:
1. Full review against base branch
2. Fetch previous comments: `gh api repos/$REPO/pulls/$PR_NUMBER/comments`
3. Skip re-commenting on same file+line+issue

### Submitting the Review

Use `gh api` to submit the review with inline comments:

```bash
# Get HEAD commit SHA
COMMIT_SHA=$(gh pr view $PR_NUMBER --json headRefOid --jq '.headRefOid')

# Submit review (use jq to construct JSON properly)
jq -n \
  --arg body "Review summary here" \
  --arg event "COMMENT" \
  --arg commit_id "$COMMIT_SHA" \
  '{body: $body, event: $event, commit_id: $commit_id, comments: [
    {path: "src/users.js", line: 42, body: "**Issue:** Missing input validation..."},
    {path: "src/api.js", line: 15, body: "**Issue:** Ignored error return value."}
  ]}' | gh api repos/$REPO/pulls/$PR_NUMBER/reviews --input -
```

**Review structure:**
- `body`: The summary (formatted per REVIEW_GUIDELINES.md). If some findings cannot be attached as inline comments (line number uncertain), include them here with file path context.
- `comments`: File-specific findings with confident line numbers. Each comment targets a file and line number so it appears directly on the diff.
- `event`: "COMMENT", "APPROVE", or "REQUEST_CHANGES"

**CRITICAL: Single review per trigger.** Submit exactly ONE review. NEVER post separate issue comments (`gh pr comment`) for initial reviews — all feedback goes in a single review submission. The only exception is responding to re-review requests or @claude mentions, which use issue comments to reply.

**CRITICAL: No fallback comments.** If you cannot determine the correct line number for a finding, include it in the review body with the file path — do NOT post it as a separate issue comment.

### Determining Line Numbers

The `line` field must be the line number in the NEW file version (at the PR's HEAD commit).

**Preferred method — fetch PR branch and read files directly:**

In local mode (or when you have git access), fetch the PR branch and use standard file/git operations:

```bash
# Fetch the PR branch
git fetch origin <branch-name>

# View diff stats
git diff main...FETCH_HEAD --stat

# Read file content directly
git show FETCH_HEAD:path/to/file.go

# Or use the Read tool on fetched files
```

This is faster and more reliable than GitHub API calls, especially for large files.

**Alternative — GitHub API (CI mode or when git fetch unavailable):**

```bash
# Fetch file content at PR HEAD
gh api repos/$REPO/contents/PATH?ref=$COMMIT_SHA --jq '.content' | base64 -d
```

Then search for the specific code pattern to find its line number. This is more reliable than parsing diff hunks.

**Fallback — parse from diff:**

1. Find the `@@` hunk header: `@@ -196,12 +196,41 @@` — the `+196` is the starting line in the new file
2. Count lines from hunk start: only count `+` lines and context lines (no prefix), skip `-` lines
3. For deleted lines, use `"side": "LEFT"` with the old file line number

**When uncertain:** If you cannot confidently determine the line number, do NOT guess. Include the finding in the review body instead:
> **src/users.js**: Missing input validation in `processUser()` function.

---

## Error Handling (CI Mode)

In CI mode, the skill MUST fail explicitly rather than complete silently when it cannot perform its review. Use `exit 1` to indicate failure.

### Permission Errors

After any `gh api` call, check for permission failures:

```bash
# Check if API call failed
if ! RESULT=$(gh api repos/$REPO/pulls/$PR_NUMBER/reviews --input - 2>&1); then
  echo "ERROR: Failed to submit review - $RESULT"
  exit 1
fi
```

Common permission errors to catch:
- `Resource not accessible by integration` — missing workflow permissions
- `Not Found` — repo access denied or PR doesn't exist
- `Validation Failed` — invalid review data (e.g., bad line numbers)

### Review Quality Failures

The skill MUST fail if it cannot produce a meaningful review:

1. **No findings with valid line numbers**: If line detection fails for ALL findings and none can be attached as inline comments, the review adds no value. **Exit with error** rather than posting only a summary.

2. **API submission failure**: If the review submission itself fails for any reason, **exit with error**.

```bash
# Track successful line detections
INLINE_COMMENT_COUNT=<number of comments with confirmed line numbers>

if [ "$INLINE_COMMENT_COUNT" -eq 0 ] && [ "$TOTAL_FINDINGS" -gt 0 ]; then
  echo "ERROR: Found $TOTAL_FINDINGS issues but could not determine line numbers for any. Review cannot proceed."
  exit 1
fi
```

### Success Criteria

A review is considered successful when:
- At least one inline comment has a verified line number, OR
- There are genuinely no findings (clean PR)
- The review was successfully submitted to GitHub

---

## Action: Respond

Handle @claude mentions that aren't review requests.

**IMPORTANT:** In CI mode, you MUST post your response to GitHub. Do not just produce text output — the user cannot see your output, only GitHub comments. Your response is invisible unless you post it.

### Fetch Context

```bash
# Get review comments (inline on code)
gh api repos/$REPO/pulls/$PR_NUMBER/comments --paginate

# Get issue comments (general PR discussion)
gh api repos/$REPO/issues/$PR_NUMBER/comments --paginate
```

### Posting the Response

After formulating your response, you MUST post it to GitHub:

```bash
# Reply to an inline review comment thread
gh api repos/$REPO/pulls/$PR_NUMBER/comments/$COMMENT_ID/replies \
  -f body="Response here"

# Reply to a general PR/issue comment
gh pr comment $PR_NUMBER --body "Response here"
```

**Choose the right method:**
- If the @claude mention is on an inline review comment → use the replies API with the comment ID
- If the @claude mention is on a general PR comment → use `gh pr comment`

Respond based on the question/request. If project guidelines specify a response format, use it.

---

## Dry Run Mode

When `--dry-run` is passed OR running in local mode (no CI environment detected):
- Analyze normally but DO NOT post to GitHub
- Output the exact JSON that would be submitted

````
DRY RUN - Review for PR #123

```json
{
  "body": "<Review summary>",
  "event": "COMMENT",
  "commit_id": "abc123...",
  "comments": [
    {"path": "src/users.js", "line": 42, "body": "**Issue:** ..."},
    {"path": "src/api.js", "line": 15, "body": "**Issue:** ..."}
  ]
}
```
````

For responses, output the comment body:

````
DRY RUN - Response for PR #123

```markdown
Your response text here exactly as it would be posted...
```
````

---

## Exit Codes (CI Mode)

In CI mode, use exit codes to signal workflow success or failure:

| Exit Code | Meaning |
|-----------|---------|
| `0` | Review completed successfully |
| `1` | Review failed (permissions, line detection, API error) |

**Local mode** does not use exit codes — it outputs dry-run results for the user to review.
