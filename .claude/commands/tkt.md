# Create Worktree for Linear Ticket

I need you to help me set up a new worktree and begin working on a Linear ticket.

Arguments provided: $ARGUMENTS

Parse the arguments as follows:
- First argument: Linear ticket ID (e.g., PRO-7111)
- Remaining arguments: Additional context or specific instructions

Example usage:
- `/tkt PRO-7111` - Just the ticket ID
- `/tkt PRO-7111 "focus on the frontend implementation only"`
- `/tkt PRO-7111 "this connects to the payment processor API"`

Please follow these steps:

1. **Fetch the Linear ticket details** using the Linear MCP:
   - Extract the ticket ID from the first argument
   - Use `mcp__linear-server__get_issue` with the extracted ticket ID
   - Note the `gitBranchName` field from the response
   - Review the title, description, and any acceptance criteria
   - Check if there's a `parentId` and fetch parent ticket details
   - Format parent ticket info as markdown for the CLAUDE-TKT-ENV.md file

2. **Create a worktree** for this ticket:
   - Pass the entire ticket JSON response as the first parameter
   - If parent ticket exists, pass the parent ticket JSON as the second parameter
   - Include any additional context from the command arguments as the third parameter
   - Run: `just claude-worktree-create '<ticket-json>' '<parent-json>' "<additional-context>"`
   - Note: Both JSON parameters must be passed as single-quoted strings to preserve formatting
   - Note the port number that will be assigned for frontend development

3. **Navigate to the new worktree**:
   - Change to the worktree directory: `../<ticket-part-of-branch-name>`
   - (e.g., if branch is `cc/pro-7111-feat`, directory will be `../pro-7111-feat`)

4. **Create a tech spec** (if needed):
   - Create a tech spec in `docs/tech-specs/<ticket-id>-<descriptive-name>.md`
   - Base it on the ticket requirements and acceptance criteria
   - Use Graphite to create a spec branch: `gt create --message "spec: <ticket title>"`
   - Commit the tech spec

5. **Create implementation branch**:
   - After tech spec is committed, create implementation branch
   - Use: `gt create --message "feat: <ticket title>"`

6. **Begin implementation**:
   - Start implementing based on the ticket requirements
   - Remember to use the assigned port for frontend development
   - Do NOT run database migrations in worktree checkouts

**Important**: Take into account any additional context provided after the ticket ID when implementing the solution.

Please start by parsing the arguments and fetching the ticket details.