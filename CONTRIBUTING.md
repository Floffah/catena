# Contributing to Catena

Contributions are welcome, here's how to make them count.

## Running in Development

You can do the following to get a development instance running:
- Install (Go)[https://go.dev/dl/], [Bun](https://github.com/oven-sh/bun), and [Just](https://github.com/casey/just)
- Start a postgres instance (e.g. with `docker compose -f deployments/dev.docker-compose.yml up -d db`)
- Create a [Clerk](https://clerk.com) account and create a new applicatipn
- Set up .env and web/.env (templates are provided) with all credentials and configuration values from the above
- Run `just dev` to start both the backend (restarts on change) and next app (hot reloads on change)

## Contributing Code

- **Read the code before changing it.** Understand the existing patterns and match the style.
- **Small PRs, incremental improvements.** A series of focused, reviewable PRs is better than one large change.
- **Discuss before building big things.** For major features or architectural changes, open an issue first. Once there's agreement, break the work into smaller pieces.
- **Bug fixes and small improvements can go straight to PR.** Not everything needs a discussion.
- **Tests are required for new features and bug fixes.** If you add functionality, please add tests to cover it.
- **Documentation updates are also contributions.** If you see something in the docs that could be clearer, feel free to open a PR to improve it.
- **Keep code in Go or TypeScript.** Avoid adding new languages or dependencies unless absolutely necessary.
- **Run `just generate`, `just format`, and `just lint` before submitting.** CI checks will fail if these aren't passing.

## AI-assisted contributions

Read the [AI contribution policy](AI_POLICY.md) for the author and peer review requirements. The same requirements apply regardless of which tools you use.

## Submitting a pull request

Use an Angular-style Conventional Commit title with an imperative summary. The scope is optional, as in `feat(repos): add branch browsing` or `fix: preserve file permissions`.

Complete the [pull request template](.github/pull_request_template.md) and confirm your own review in the contributor checklist. Explain what changed, why it matters, and how you approached it. Include the checks you ran and disclose known failures, limitations, and any breaking changes. Short answers are welcome for small changes.

Draft pull requests are welcome for early feedback. For open, ready-for-review pull requests from human contributors, GitHub Actions posts a reminder when **What**, **Why**, or **How** is empty, or when the **Contributor checklist** is missing, empty, or has unchecked items. Editing the description refreshes the reminder and completing it removes the comment.
