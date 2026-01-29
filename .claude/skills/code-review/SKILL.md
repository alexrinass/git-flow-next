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

1. **Test Coverage** - Sufficient tests covering success, error, and edge cases
2. **Guidelines & Architecture** - Follows coding guidelines and project patterns
3. **Code Quality** - No ignored errors, unnecessary complexity, or dead code
4. **Security** - No injection vulnerabilities, path traversal, or exposed secrets

Use a concise and direct tone. Put detailed feedback in inline comments with fix prompts.

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
  "body": "Review summary here (general remarks, not file-specific)",
  "event": "COMMENT",
  "commit_id": "SHA_PLACEHOLDER",
  "comments": [
    {"path": "src/file.js", "line": 42, "body": "Inline feedback here"}
  ]
}
EOF

sed -i.bak "s/SHA_PLACEHOLDER/$COMMIT_SHA/" /tmp/review.json && rm /tmp/review.json.bak
gh api repos/$REPO/pulls/$PR_NUMBER/reviews --input /tmp/review.json
```

**Review structure:**
- `body`: General remarks and proposals not tied to specific files
- `comments`: Array of inline comments on specific lines
- `event`: "COMMENT", "APPROVE", or "REQUEST_CHANGES"

**Line numbers:**
- `line` = line number in new file version (right side of diff)
- For deleted lines, add `"side": "LEFT"` to comment object

---

## Action: Respond

Handle @claude mentions that aren't review requests.

### Fetch Context

```bash
# Get review comments (inline on code)
gh api repos/$REPO/pulls/$PR_NUMBER/comments --paginate

# Get issue comments (general PR discussion)
gh api repos/$REPO/issues/$PR_NUMBER/comments --paginate
```

### Responding

```bash
# Reply to inline review comment
gh api repos/$REPO/pulls/$PR_NUMBER/comments/$COMMENT_ID/replies \
  -f body="Response here"

# Reply to general PR comment
gh pr comment $PR_NUMBER --body "Response here"
```

Respond based on the question/request. If project guidelines specify a response format, use it.

---

## Dry Run Mode

When `--dry-run` is passed OR running in local mode (no CI environment detected):
- Analyze normally but DO NOT post to GitHub
- Output what would be posted

```
DRY RUN - No changes will be made

Would submit review to PR #123:
  Event: COMMENT
  Body: "..."

Inline comments (3):
  src/file.js:10 - "..."
  src/file.js:25 - "..."
  src/api.js:100 - "..."
```
