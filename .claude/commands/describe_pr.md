# Generate PR Description

You are tasked with generating a comprehensive pull request description following the repository's standard template.

## Steps to follow:

1. **Read the PR description template:**

   - First, check if `pull_request_template.md` exists in the root of the repository
   - Read the template carefully to understand all sections and requirements

1. **Identify the PR to describe:**

   - Check if the current branch has an associated PR: `gt branch info`

1. **Check for existing description:**

   - Check if `thoughts/shared/prs/{number}_description.md` already exists
   - If it exists, read it and inform the user you'll be updating it
   - Consider what has changed since the last description was written

1. **Gather comprehensive PR information:**

   - Get the changed files: `git diff $(gt parent) --name-only`
     - Identify non-generated code files to inspect and inform the PR description
     - Pay close attention to any changes in the `db/`, `proto/`, `api/`, `workflows/`, `pkg/` or `frontend/apps` directories
     - Ignore generated files in `pkg/proto`, `frontend/packages/proto`

1. **Analyze the changes thoroughly:** (ultrathink about the code changes, their architectural implications, and potential impacts)

   - Read through the entire diff carefully
   - For context, read any files that are referenced but not shown in the diff
   - Understand the purpose and impact of each change
   - Identify user-facing changes vs internal implementation details
   - Look for breaking changes or migration requirements

1. **IMPORTANT: Ask the engineer questions of how they have manually verified the changes**

   - If anything, please include a summary of the manual verification the engineer has done

1. **Get metadata**

   - run `git config user.name` to get the git username for the pr description directory

1. **Generate the description:**

   - Fill out each section from the template thoroughly:
     - Answer each question/section based on your analysis
     - Be specific about problems solved and changes made
     - Focus on user impact where relevant
     - Include technical details in appropriate sections
     - Write a concise changelog entry
   - Ensure all checklist items are addressed (checked or explained)

1. **Save and sync the description:**

   - Write the completed description to `thoughts/{git_username}/prs/{number}_description.md`
   - Show the user the generated description

1. **Update the PR:**

   - Update the PR description using the GitHub CLI:
     ```bash
     gh pr edit {number} --body-file /full/path/to/thoughts/{git_username}/prs/{number}_description.md
     ```
   - Note: Use `--body-file` instead of `--body` with command substitution, as `$(cat ...)` may not expand properly in all shells
   - If the update fails, provide the PR URL for manual update: https://github.com/premiumlabs/monorepo/pull/{number}

## Important notes:

- Be thorough but concise - descriptions should be scannable
- Focus on the "why" as much as the "what"
  - Ask the user questions if anything is unclear
- Include any breaking changes or migration notes prominently
- If the PR touches multiple components, organize the description accordingly
