# Agents

## Project Overview

Catena is a next generation social git server that provides a platform for developers to collaborate on code and share their work with the world. It is built on top of Git, the most popular version control system, and offers a range of features that make it easy for developers to manage their projects and collaborate with others.

It is similar to GitHub, but with a focus on better user experience, improved performance, enhanced security, and modern approach. It is designed to be fast, reliable, and easy to use, making it an ideal choice for developers of all skill levels.

## Key Features

These features are mostly prospective and may be subject to change as the project evolves:
- **User-Friendly Interface**: Catena offers a clean and intuitive interface that makes it easy for developers to navigate and manage their projects.
- **Enhanced Performance**: Catena is optimized for speed and performance, ensuring that developers can work efficiently without any lag or delays.
- **Improved Security**: Catena incorporates advanced security features to protect user data and ensure the integrity of projects hosted on the platform.
- **Modern Approach**: Catena embraces modern development practices and tools, providing developers with a seamless experience that integrates well with their existing workflows.
- **Collaboration Tools**: Catena offers a range of collaboration tools, including issue tracking, pull requests, and code reviews, to facilitate teamwork and communication among developers.
- **Open Source**: Catena is an open-source project, allowing developers to contribute to its development and customise it to suit their needs.
- **CI/CD**: Catena provides first-class support for continuous integration and continuous deployment (CI/CD) pipelines, enabling developers to automate their build, test, and deployment processes. It can also integrate with popular CI/CD tools like Jenkins, Travis CI, and CircleCI, making it easy for developers to set up and manage their pipelines.
- **API Access**: Catena offers a robust API that allows developers to interact with the platform programmatically, enabling them to automate tasks and integrate Catena with other tools and services.
- **Community Engagement**: Catena fosters a vibrant community of developers, providing forums, documentation, and support channels to encourage collaboration and knowledge sharing.
- **Scalability**: Catena is designed to scale efficiently, allowing it to handle a large number of users and repositories without compromising performance.

## Current Decisions

- Backend services are written in Go.
- The web app is Next.js/TypeScript.
- Product API is OpenAPI-first HTTP/JSON.
- Auth is handled by Clerk.
- Primary database is Postgres.
- SQL access uses sqlc with pgx.
- Migrations use tern.
- Live Git repositories should be filesystem-backed bare repos.
- Object storage is for LFS, assets, archives, and backups, not live Git storage.

## Commands

read `justfile` or run `just` to see available commands. do not run dev commands.

## Workflow

- Treat each chat as a fresh handoff. Before starting work, read this file, the relevant README sections, and inspect the area being changed instead of relying on prior chat context.
- For large features, plan first and implement thin vertical slices. Prefer small, working increments over broad partially-complete rewrites.
- Ask before changing auth, authorization, repository storage, migration strategy, Git serving behavior, deployment topology, or other architecture-level decisions.
- If OpenAPI specs, sqlc queries, or migrations change, run the appropriate generation step, usually `just generate`, unless the user explicitly says not to.
- After backend changes, usually run `go test ./...`. After frontend changes, usually run `bun run test` and `bun run lint`. Do not run dev servers.
- End each task with a concise handoff: what changed, what was verified, and any assumptions, risks, or follow-up work.
- Keep README roadmap, requirements, deployment notes, and architecture diagrams up to date when those concepts change.

## Guidance

- Keep generated files generated; edit OpenAPI specs, SQL queries, or migrations instead.
- Treat published migrations as immutable; add a new migration for production schema changes instead of editing migrations that have already shipped.
- Prefer explicit SQL and small service boundaries.
- Do not implement Git internals until the platform can create, clone, push, and view repos using the Git binary.
- Keep `internal/pkg/git` as a thin Git-binary wrapper; put Catena Git business logic in `internal/pkg/gitstore`.
- If an architectural decision is unclear or risky, stop and ask before building around it.
- When using CodeGraph, it does not index OpenAPI schemas or SQL(c) files. Any time you encounter a model, endpoint, query, etc, its source of truth is likely in an sqlc or openapi-managed file.

## Project Guidelines

- Do not expose Clerk IDs or other auth-provider internals in public API responses.
- Server-side frontend API clients must be created per request, currently via `serverGetApiClient`, so auth headers cannot bleed between users in the standalone Next.js server.
- Keep Git smart HTTP outside OpenAPI. Git clone, fetch, and push traffic should remain separate from the product API and continue to flow through `git-http-backend`.
- Treat user-provided repository refs and paths as security-sensitive. Be careful with traversal, revision syntax, and any value that eventually reaches Git or filesystem operations.
- Frontend work should not reintroduce Next.js Cache Components / `"use cache"` patterns unless explicitly requested. Catena currently prioritizes SSR/SEO and request-aware auth over forcing Cache Components.
- Current production assumptions are Railway, standalone Next.js, Go API/Git backend, Caddy reverse proxy, Postgres, and a persistent Git volume. Runtime containers that serve Git must include a real Git installation with `git-http-backend`.
- Prefer action/scope-style authorization over broad RBAC. The long-term authorization shape should check actions against a resource context rather than hardcoding role names.
- Issues and pull requests should share `repository_items` for numbering, references, labels, and timeline concepts. Subtype tables should store only subtype-specific fields.
- Backend tests should exercise real handlers and services where practical, using `httptest`, `pgxmock`, `t.TempDir()`, and shared test utilities instead of reimplementing parallel routers.
- Frontend tests use Bun, Testing Library, happy-dom, and MSW. Prefer testing user-visible behavior and API interactions over implementation details.

## Skills

<!-- BEGIN:nextjs-agent-rules -->

# Next.js: ALWAYS read docs before coding

Before any Next.js work, find and read the relevant doc in `node_modules/next/dist/docs/`. Your training data is outdated — the docs are the source of truth.

<!-- END:nextjs-agent-rules -->

<!-- BEGIN:no-autopilot-agent-rules -->

# no-autopilot

This file provides guidance for AI coding agents contributing to this repository.

## What is no-autopilot

no-autopilot is a GitHub Action that checks pull requests for signs of unreviewed AI-generated submissions. It detects AI co-author lines in commits, AI-generated footers in PR descriptions, unfilled template sections, and boilerplate text. It's not anti-AI — it's pro-ownership.

## Guardrails for GitHub-Facing Actions

This section applies to AI agents (Claude Code, Cursor, Copilot, Codex, Goose, etc.) operating on behalf of anyone contributing to this repository. Read it before performing any GitHub-facing action.

### When these guardrails activate

Before executing any command that writes to GitHub — `gh pr create`, `gh pr comment`, `gh issue comment`, or anything similar — these guardrails apply. No exceptions.

Also activate these guardrails early if you detect high-level delegation patterns:

- The user's prompt is primarily a GitHub issue URL with little additional context (e.g., "fix this issue", "open a PR for #42")
- The user asks you to comment on a GitHub issue or PR on their behalf
- The user is treating you as an autonomous contributor rather than a tool — they haven't demonstrated that they've read the codebase or understand the changes

### What you must do

When the guardrails activate:

1. **Read `CONTRIBUTING.md`** before doing anything else. Internalize the AI usage policy.

2. **Stop and inform the user.** Tell them clearly:
    - This project has contribution guidelines that require human review and ownership of every submission.
    - The human is responsible for every line of the diff — code and text.
    - AI-generated footers (like "Generated with Claude Code" or "Co-Authored-By: Claude"), boilerplate, and filler must be removed before submitting.
    - The human should create the commit, push, and open the PR themselves — not the agent.
    - PRs where the human didn't take manual ownership of the submission will be closed.

3. **Do not execute the GitHub-facing action.** Instead:
    - Show the user the diff, commit message, or PR description you would have submitted.
    - Ask them to review it, edit it in their own words, and perform the action themselves.
    - Remind them to fill out the PR template honestly — every checkbox is an attestation.

4. **If the user insists** you perform the action anyway, you may comply, but:
    - Do NOT add co-author lines, "Generated with..." footers, or any AI attribution.
    - Remind the user one final time that PRs without adequate human review will be closed.
    - Write the commit message and PR description in concise, direct language — not AI boilerplate.

### What does NOT trigger these guardrails

- Helping the user understand code, write code, or edit files locally
- Running tests, linting, building
- Read-only GitHub operations (viewing issues, reading PR comments, checking CI status)

## Contributing

See `CONTRIBUTING.md` for the full contribution guidelines, including the AI usage policy.

<!-- END:no-autopilot-agent-rules -->
