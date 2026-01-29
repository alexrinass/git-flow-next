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

**Use project guidelines for:**
- What to focus on in reviews
- Review summary format
- Inline comment format
- Response format/tone

**If no guidelines found, use defaults:**

1. **Test Coverage** - Evaluate test additions and changes thoroughly (see Test Analysis below)
2. **Guidelines & Architecture** - Follows coding guidelines and project patterns
3. **Code Quality** - No ignored errors, unnecessary complexity, or dead code
4. **Security** - No injection vulnerabilities, path traversal, or exposed secrets

Use a concise and direct tone. Put detailed feedback in inline comments with fix prompts.

---

## Test Analysis

When reviewing tests, the **Test Coverage** section in the review summary must go beyond a one-line assessment. Include:

1. **Overview** — Briefly list each test introduced or changed and what it covers.
2. **Coverage completeness** — Are new/changed code paths covered? Call out any critical missing cases (error paths, edge cases, boundary conditions).
3. **Test quality** — Do tests follow project patterns (naming, structure, testutil usage, one test per function)?
4. **Test correctness** — Do tests actually verify the right behavior, or are they superficial (e.g., only asserting "no error" without checking output or state)?
5. **Duplicates and low-value tests** — Flag tests that are redundant or test something too trivial to warrant their own test function.

Format the Test Coverage section in the review summary like this:

```
- **Test Coverage**:
  - `TestFoo`: Tests X when Y — good coverage of error path
  - `TestBar`: Tests basic Z — missing: no test for Z with invalid input
  - Missing: No test for the <scenario> code path, which is critical because <reason>
  - `TestBaz` and `TestQux` overlap significantly — consider consolidating
```

Keep it concise but specific. Every test should be mentioned by name with a brief note on what it covers and any gaps.

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

Use GitHub's Reviews API to submit review with inline comments:

```bash
COMMIT_SHA=$(gh pr view $PR_NUMBER --json headRefOid --jq '.headRefOid')

cat > /tmp/review.json << 'EOF'
{
  "body": "- **Code Quality**: Brief assessment...\n- **Test Coverage**:\n  - `TestProcessUser`: Tests happy path — good\n  - `TestProcessUserError`: Tests DB failure — missing: no test for invalid email input\n  - Missing: No test for concurrent access, which is critical for the new batch endpoint\n- **Security**: No concerns.\n- **Commit Messages**: Follow guidelines.\n\n**Action Items:**\n- Fix ignored error in processUser()\n- Add input validation for email field\n- Add test for invalid email input",
  "event": "COMMENT",
  "commit_id": "SHA_PLACEHOLDER",
  "comments": [
    {"path": "src/users.js", "line": 42, "body": "**Issue:** Missing input validation...\n\n<details>\n<summary>Fix prompt</summary>\n\n```prompt\nIn src/users.js, add input validation to processUser()...\n```\n\n</details>"},
    {"path": "src/api.js", "line": 15, "body": "**Issue:** Ignored error return value from `db.Save()`."},
    {"path": "src/old.js", "line": 8, "side": "LEFT", "body": "**Issue:** This deleted validation was still needed."}
  ]
}
EOF

sed -i.bak "s/SHA_PLACEHOLDER/$COMMIT_SHA/" /tmp/review.json && rm /tmp/review.json.bak
gh api repos/$REPO/pulls/$PR_NUMBER/reviews --input /tmp/review.json
```

**Review structure:**
- `body`: ONLY the summary (formatted per REVIEW_GUIDELINES.md or default format). NEVER put file-specific feedback in the body.
- `comments`: ALL file-specific findings MUST go here as inline comments on specific lines. Each comment targets a file and line number so it appears directly on the diff.
- `event`: "COMMENT", "APPROVE", or "REQUEST_CHANGES"

### Determining Line Numbers

The `line` field in each comment must be the line number in the file (not the diff position). To find the correct line number from `gh pr diff` output:

1. Find the relevant `@@` hunk header for your finding:
   ```
   @@ -196,12 +196,41 @@ func executeFinish(...)
   ```
   The `+196` means the new file version starts at line 196 in this hunk.

2. Count lines from the hunk start. Only count `+` lines and context lines (lines without `+` or `-` prefix). Skip `-` lines entirely — they don't exist in the new file.

3. **Worked example:**
   ```diff
   @@ -10,7 +10,9 @@ func process() {
        existing line       // line 10 (context)
        another line        // line 11 (context)
   -    old code            // SKIP (deleted, not in new file)
   +    new code            // line 12
   +    added line          // line 13
        unchanged           // line 14
   ```
   To comment on `added line`, use `"line": 13`.

4. For deleted lines (prefixed with `-`), add `"side": "LEFT"` and use the line number from the **old** file (the `-10` side of the hunk header).

**Important:** If you cannot confidently determine the exact line number, read the file at the PR's HEAD to verify:
```bash
gh api repos/$REPO/contents/PATH?ref=$COMMIT_SHA --jq '.content' | base64 -d | head -n 50
```

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
- Output the exact JSON that would be submitted, so the user sees exactly what would be posted

For reviews, output the review JSON:

````
DRY RUN - Review for PR #123

```json
{
  "body": "- **Code Quality**: ...\n- **Test Coverage**:\n  - `TestFoo`: ...\n  - Missing: ...\n- **Security**: ...\n...",
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
