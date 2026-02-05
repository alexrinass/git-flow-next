---
name: pr-summary
description: Generate a PR summary and write to .ai/ pr_summary.md
allowed-tools: Read, Grep, Glob, Write, Bash
---

# Create PR Summary

Generate a pull request summary based on branch changes.

## Instructions

1. **Gather Context**
   - Get current branch name
   - Find workflow folder for this branch (`.ai/issue-<number>-*` or `.ai/<feature>/`)
   - Read analysis.md and plan.md if they exist
   - Get associated issue number from branch name or commits

2. **Analyze Changes**
   ```bash
   # All commits on this branch
   git log main...HEAD --oneline

   # Full diff
   git diff main...HEAD --stat

   # Changed files
   git diff main...HEAD --name-only
   ```

3. **Generate Summary**

   Follow the format from `.github/PULL_REQUEST_TEMPLATE.md`.

   Write to `.ai/<folder>/pr_summary.md`:

   ```markdown
   <Summary prose — no header. Describe what changed and why in 1-3 sentences.
   Mention key areas touched and important decisions.
   Link to resolved issues with "Resolves #ISSUE".>

   ## Notes

   <Optional. Call out risks, edge cases, breaking changes, or scope clarifications.
   Remove this section if not applicable.>
   ```

   **Guidelines:**
   - Keep it concise — prose, not bullet lists or tables
   - Focus on what and why, not line-by-line details
   - Do NOT include checklists — those are for author verification only

4. **Create PR Command**

   Output the `gh pr create` command:

   ```bash
   gh pr create --title "<title>" --body "$(cat .ai/<folder>/pr_summary.md)"
   ```

5. **Report Completion**
   - Show path to pr_summary.md
   - Show the PR creation command
   - Remind to push branch first if not already pushed

## PR Title Guidelines

- Use imperative mood: "Add feature" not "Added feature"
- Be specific but concise
- Include issue number if applicable: "Add squash merge support (#42)"
