# Create Technical Specification

You are tasked with creating comprehensive technical specifications by synthesizing outputs from research and implementation planning into a high-level architectural document.

## Initial Response

When this command is invoked:

1. **Check if parameters were provided**:
   - If file paths to research and/or plan documents were provided, skip the default message
   - Immediately read any provided files FULLY
   - Begin the spec creation process

2. **If no parameters provided**, respond with:

```
I'll help you create a technical specification. This spec will synthesize your research findings and implementation plan into a comprehensive architectural document.

Please provide:
1. Path to research document (e.g., `thoughts/shared/research/YYYY-MM-DD-topic.md`)
2. Path to implementation plan (e.g., `thoughts/shared/plans/YYYY-MM-DD-topic.md`)
3. Any additional context or constraints
4. Link to product requirements (Linear ticket, PRD, etc.)

I'll analyze these inputs and create a detailed tech spec that covers problem statement, architecture, implementation approach, and roll-out strategy.

Tip: You can invoke this command with documents directly:
`/create_spec thoughts/shared/research/2025-01-08-feature.md thoughts/shared/plans/2025-01-08-feature.md`
```

Then wait for the user's input.

## Process Steps

### Step 1: Context Gathering & Analysis

1. **Read all mentioned files immediately and FULLY**:
   - Research documents
   - Implementation plans
   - Product requirement documents
   - Related tech specs
   - Any Linear tickets or design documents
   - **IMPORTANT**: Use the Read tool WITHOUT limit/offset parameters to read entire files
   - **CRITICAL**: DO NOT spawn sub-tasks before reading these files yourself in the main context

2. **Verify understanding with the user**:
   ```
   Based on the research and plan documents, I understand we're building:

   **Problem**: [Summary of user problem from research]
   **Solution**: [High-level approach from plan]
   **Scope**: [Key phases from implementation plan]

   Before I create the spec, I need clarification on:
   - [Business context question]
   - [Stakeholder identification]
   - [Success metrics question]
   - [Rollout approach question]

   Are there any additional constraints or requirements I should consider?
   ```

### Step 2: Structure Development

Present a proposed structure to the user:

```
Here's the spec outline I'll create:

**Table of Contents**

# Problem Statement
[Summarize user problem and customer value]
- Link to PRD

## Open Questions
[Remaining questions the team needs answered]

## Assumptions
[Assumptions being made in this doc]

## Glossary [optional]
[Terms and definitions]

## Requirements
[Functional and non-functional requirements]

## Out of scope
[What we're specifically NOT solving]

## Prior References
[Links to previous tech specs with relevance notes]

# Proposed Solution (Architecture Design)
[Solution description with complexity callouts]

## Alternate Solution(s)
[Other solutions considered and why abandoned]

## Risk and mitigations
[Potential risks and mitigation strategies]

# Roll-out Plan
[How this should be rolled out]

## Permissions / Auditing / Monitoring
[Permissions needed, audit events, monitoring metrics]

# References
[External documents and API docs]

Does this structure capture everything needed?
```

### Step 4: Write the Technical Specification

1. **Determine the filename**:
   ```
   thoughts/shared/specs/YYYY-MM-DD-PRO-XXXX-brief-description.md
   ```
   - Format: `YYYY-MM-DD-PRO-XXXX-description.md` where:
     - YYYY-MM-DD is today's date
     - PRO-XXXX is the ticket number (omit if no ticket)
     - description is a brief kebab-case description
   - Examples:
     - With ticket: `thoughts/shared/specs/2025-01-15-PRO-1234-scheduled-config-changes.md`
     - Without ticket: `thoughts/shared/specs/2025-01-15-role-management.md`

2. **Use this template structure**:

````markdown
# [Feature Name] Technical Specification

**Table of Contents**

# Problem Statement

[Summarize the user problem and the customer value we'd like to create. e.g. for summative bonus, it'd be adding policy-level reporting, chargebacks for churned policies. Link to PRD]

Clearly articulate the problem from the user's perspective:
- What pain point are users experiencing?
- What can't users do today that they need to do?
- What inefficiency or risk does this create?

## Open Questions

[List remaining questions the team needs answered]

## Assumptions

[What assumptions are being made in this doc]

## Glossary [optional]

[List terms and their definitions used within the doc or links to our internal glossary]

## Requirements

[]

## Out of scope

[What are you specifically not trying to solve, link to any use cases & examples from PRD]

## Prior References

[Link to previous tech specs. If possible, add something that summarizes and clarifies previous tech specs.

e.g. for tech specs that you're referencing give a one liner of each one's relevance and takeaways and/or add a diagram or table]

# Proposed Solution (Architecture Design)

[Describe your solution and requirements. Add call outs for specific points of complexity]

## Alternate Solution(s)

[What other solutions were considered? Why were they abandoned?]

## Risk and mitigations

[]

# Roll-out Plan

[How should this be rolled out]

## Permissions / Auditing / Monitoring

[What are the permissions we need to support]

[Audit log events]

[What monitoring metrics should be tracked?]

# References

[Other external documents and references e.g. external API docs]
````

### Step 5: Review and Iterate

1. **Present the draft spec location**:
   ```
   I've created the technical specification at:
   `thoughts/specs/YYYY-MM-DD-PRO-XXXX-description.md`

   Please review and let me know:
   - Is the problem statement clear?
   - Are the requirements complete?
   - Are assumptions documented?
   - Should any alternate solutions be added?
   - Are open questions listed?
   - Are risks and mitigations comprehensive?
   ```

2. **Iterate based on feedback**:
   - Add missing sections
   - Clarify technical details
   - Expand on alternatives
   - Add or refine requirements
   - Resolve or add open questions
   - Update risk analysis

3. **Continue refining** until the user is satisfied

## Important Guidelines

1. **Synthesize, Don't Copy**:
   - Extract key insights from research and plan
   - Write for a broader audience (technical and non-technical stakeholders)
   - Elevate implementation details to architectural concepts
   - Connect technical decisions to business outcomes

2. **Business Focus First**:
   - Start with problem and impact (why we're doing this)
   - Then describe solution approach (what we're building)
   - Save implementation details for references to the plan
   - Make it readable by product, operations, and executive stakeholders

3. **Completeness**:
   - Read all input documents FULLY before writing
   - Don't leave sections incomplete or with TODOs
   - Resolve open questions or explicitly call them out
   - Include all necessary references

4. **Clarity**:
   - Use terminology consistently
   - Define domain-specific terms
   - Provide concrete examples
   - Include diagrams from prior references when relevant

5. **Actionability**:
   - Success metrics must be measurable
   - Rollout plan must be specific
   - Open questions must identify who can answer

6. **No Placeholders**:
   - Never write the spec with placeholder values
   - If information is missing, ask the user
   - If you're uncertain, research or clarify
   - All sections must be complete before presenting

7. **Minimize code snippets**
   - The code snippets will be part of the detailed implementation plan
	 - The tech spec will be focused on high-level summaries
	 - Only include code snippets that are essential to concisely summarize a concept

## Integration with Other Commands

This command builds on:

- **`/research_codebase`**: Provides detailed technical findings and current state analysis
- **`/create_plan`**: Provides implementation approach, phases, and success criteria

The tech spec synthesizes both into a comprehensive architectural document that:
- Adds business context and stakeholder perspective
- Defines success metrics and monitoring
- Documents simple and concise rollout plan using feature flags in necessary
- Records alternatives considered
- Tracks open questions and decisions

## Example Interaction Flow

```
User: /create_spec thoughts/shared/research/2025-01-15-pay-cycle.md thoughts/shared/plans/2025-01-15-pay-cycle.md
