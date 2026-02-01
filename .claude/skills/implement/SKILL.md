---
name: implement
description: Execute an implementation plan, making changes and committing properly
allowed-tools: Read, Grep, Glob, Edit, Write, Bash
---

# Implement Plan

Execute the implementation plan, making code changes and committing according to guidelines.

## Instructions

1. **Find and Load Plan**
   - Detect workflow folder from current git branch
   - Read `.ai/<folder>/plan.md`
   - If no plan exists, suggest running `/create-plan` first

2. **Verify Prerequisites**
   - Confirm on correct feature branch
   - Check working directory is clean (`git status`)
   - Ensure tests pass before starting: `go test ./...`

3. **Load Guidelines**
   - Read CODING_GUIDELINES.md for implementation standards
   - Read COMMIT_GUIDELINES.md for commit format
   - Read GIT_TEST_SCENARIOS.md when implementing tests (required for Git scenario setup)

4. **Execute Tasks**

   For each task in the plan:

   ### Before Each Task
   - Review task requirements and identify all files to modify
   - Read CODING_GUIDELINES.md rules relevant to the change

   ### Code Changes
   - Make focused, atomic changes
   - Follow existing patterns in the codebase
   - Keep changes minimal — don't over-engineer

   ### After Each Task
   - Verify build: `go build ./...`
   - Run tests: `go test ./...`
   - Review changes: `git diff`

5. **Commit Strategy**

   ### When to Commit
   - After completing a logical unit of work
   - After each checkpoint in the plan
   - Keep commits atomic and focused
   - Use `/commit` skill for proper formatting

6. **Checkpoint Verification**

   At each checkpoint in the plan:
   - [ ] Build succeeds: `go build ./...`
   - [ ] Tests pass: `go test ./...`
   - [ ] Expected behavior works
   - [ ] Changes committed

7. **Track Progress**

   Update plan.md checkboxes as tasks complete:
   ```markdown
   - [x] Completed task
   - [ ] Pending task
   ```

8. **Handle Issues**

   If problems arise:
   - Document the issue
   - Check if it affects the plan
   - Adjust approach if needed
   - Ask for clarification if blocked

## Completion

When all tasks are done:
1. Verify all tests pass
2. Check all checkboxes in plan.md are complete
3. Suggest running `/local-review` before PR
