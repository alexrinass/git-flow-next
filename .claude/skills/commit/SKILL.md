---
name: commit
description: Create a commit following COMMIT_GUIDELINES.md
argument-hint: [optional message hint]
allowed-tools: Bash, Read, Grep
---

# Commit Changes

Create a well-formatted commit following the project's COMMIT_GUIDELINES.md.

## Instructions

1. **Read COMMIT_GUIDELINES.md** for all formatting rules, commit types, and conventions

2. **Check Status**
   - Run `git status` to see staged and unstaged changes
   - Run `git diff --cached` to see what will be committed
   - If nothing staged, ask user what to stage or stage all with `git add -A`

3. **Analyze Changes**
   - Review the diff to understand what changed
   - Identify the appropriate commit type and scope per COMMIT_GUIDELINES.md

4. **Write Commit Message** per COMMIT_GUIDELINES.md rules

5. **Create Commit** using HEREDOC format for proper formatting:
   ```bash
   git commit -m "$(cat <<'EOF'
   <type>(<scope>): <subject>

   <body>

   <footer>
   EOF
   )"
   ```

6. **Verify** — run `git log -1` to check the message

## Hard Wrapping Technique

When a guideline specifies a hard line wrap limit, write the text as continuous flowing prose first, then hard wrap at the specified character limit, filling each line as close to the limit as possible.

## Example

```bash
git commit -m "$(cat <<'EOF'
feat(finish): Add squash merge strategy

Implements --squash flag for finish command that performs a squash merge
instead of a regular merge, creating a single commit with all branch
changes on the target branch.

- Add squash option to merge strategy enum
- Update finish command to handle squash flag
- Add configuration support for default squash behavior

Resolves #42
EOF
)"
```
