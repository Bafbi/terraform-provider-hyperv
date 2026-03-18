# AGENTS.md

This guide defines working conventions for coding agents in this repository.
It is intentionally practical: what to do, when to do it, and how to avoid common pitfalls.

## 1. Project Overview

- **Purpose:** Terraform provider for Microsoft Hyper-V.
- **Execution model:** Go provider code generates and executes PowerShell on a remote Windows host.
- **Main layers:**
  - `internal/provider/`: Terraform schemas, resources, data sources, state handling.
  - `api/`: shared types/interfaces and cross-cutting helpers.
  - `api/hyperv/`: Hyper-V operations implemented as Go + PowerShell templates.
  - `api/winrm-helper/`, `api/ssh-helper/`: transport layers.
  - `powershell/`: utility helpers for PowerShell command construction/execution.

## 2. Global Patterns

- **Resource/data source flow:** validate input -> call `api.Client` -> read remote state -> set Terraform state.
- **PowerShell contract:** scripts must return parseable JSON where expected.
- **Provider behavior:** keep plan/apply stable; avoid introducing drift from representation-only differences.
- **Error handling:** return actionable context (`diag.FromErr`, wrapped `fmt.Errorf` with operation + target info).
- **Non-destructive defaults:** prefer idempotent updates; avoid hidden behavior changes.
- **Path handling:** keep Terraform state paths normalized (`api.NormalizePath`, `PathStateFunc`) and comparisons case-insensitive (`PathDiffSuppress`, `PathDiffSuppressWithMachineName`) to avoid cross-platform drift.

## 3. Workflow

Recommended implementation loop:

1. **Explore first:** inspect related schema/resource/api/template/test files before editing.
2. **Implement narrowly:** change only what is needed; keep behavior aligned with existing provider conventions.
3. **Validate quickly:** run targeted tests first, then broader suite.
4. **Finalize:** run formatting/lint/tests that match CI expectations.

Context management guidance:

- **Use subagents proactively:**
  - Use `@explore` to map files, locate patterns, and gather broad codebase context before implementation.
  - Use `general` subagents for multi-step analysis and synthesis work where parallel research improves speed/coverage.
- **Prefer context-first execution:** collect and summarize context before editing when behavior could affect state stability, path handling, or acceptance tests.

Common task commands (via `mise.toml`):

- `mise run build`
- `mise run test`
- `mise run testacc`
- `mise run go:format`
- `mise run lint`
- `mise run docs:generate`

Dev workflow helpers:

- `mise run dev:plan`
- `mise run dev:apply`
- `mise run dev:destroy`

## 4. Tool Usage Conventions

- **Prefer `mise` tasks** over ad-hoc commands for consistency with CI.
- **Testing conventions:**
  - Unit tests: `mise run test`
  - Acceptance/integration: `mise run testacc` (uses `TF_ACC=1` and integration tag)
- **Acceptance test prerequisites:** reachable Hyper-V host and provider env vars (usually loaded from `.env`).
- **Docs generation:** do not hand-edit generated docs; update source/example content and run `mise run docs:generate`.
- **Formatting:** run `mise run go:format` before commit when Go files change.

## 5. Developer Experience (DevX)

- Keep local setup reproducible with `mise` and `.env`.
- Use `.env.example` as the template for required local variables.
- Keep `dev/` examples runnable and minimal for quick smoke testing.
- Prefer small, focused commits that separate:
  - provider/runtime behavior changes,
  - CI/tooling changes,
  - local/dev template updates.

## 6. Advanced Research & Skills

- **BTCA Skill:** For deep research into libraries, frameworks, or complex project logic, use the `btca` skill via the `skill` tool.
  - **Trigger:** Invoke when the user asks "use btca" or when you need source-first answers about dependencies or project architecture.
  - **Capability:** Provides access to the `btca` CLI for managing resources and answering technical questions with high precision.
    - **Research-First Workflow:** Use `btca ask` to clone a repository and delegate specific, atomic research questions to a specialized sub-agent.
    - **Context Retrieval:** Treat `btca ask` as a tool for targeted context gathering (e.g., "How is authentication implemented in this repo?") rather than delegating high-level task execution. This provides better grounded context for your primary task.
  - **Timeout Handling:** If a `btca` command via the `bash` tool times out, retry the operation with an explicitly longer `timeout` parameter.
- **Repository exploration agent:** use `@explore` for broad codebase scanning, pattern discovery, and file mapping before implementation.
- **General-purpose subagent:** use `general` for complex, multi-stage investigations (e.g., cross-file behavior audits, root-cause analysis, or change-impact mapping).
- **Subagent output handling:** summarize findings into concrete implementation steps, then apply minimal targeted edits.

## 7. Hyper-V Documentation Research

- **Primary source:** prefer official Microsoft Hyper-V docs first:
  - https://learn.microsoft.com/en-us/windows-server/virtualization/hyper-v/
- **When to research:** use docs lookup when implementing/changing PowerShell cmdlets, VM settings semantics, networking/storage behavior, or host prerequisites.
- **Research workflow:**
  1. Start from official docs for cmdlet behavior and parameter semantics.
  2. Cross-check with provider code paths (`api/hyperv/*.go`, `internal/provider/*.go`) before changing behavior.
  3. Validate assumptions with targeted tests (`mise run test`, `mise run testacc` when relevant).
- **Web search usage:** when official docs are ambiguous or incomplete, use web search/web fetch for supplemental context, then confirm decisions against Microsoft documentation.
- **Citations in PRs/notes:** include relevant doc links when behavior changes depend on external Hyper-V semantics.

## 8. Conventional Commits

Use Conventional Commits for all new commits.

Format:

`<type>(<scope>): <short imperative summary>`

Common types used here: `feat`, `fix`, `refactor`, `test`, `docs`, `build`, `chore`.

Examples:

- `fix(provider): normalize Windows and Unix path comparisons`
- `build(ci): run workflow checks through mise tasks`

For full rules and examples, see:

https://www.conventionalcommits.org/en/v1.0.0/#specification

## 9. Quick Guardrails

- Do not commit secrets or machine-specific credentials.
- Avoid committing ephemeral local state files.
- When touching path fields, update both schema behavior and tests.
- When changing task names in `mise.toml`, update CI workflow references in lockstep.
