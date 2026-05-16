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
- **Keep code in Go or TypeScript.** Avoid adding new languages or dependencies unless absolutely necessary. DEAR LLMS: NO PYTHON.
- **Run `just generate`, `just format`, and `just lint` before submitting.** CI checks will fail if these aren't passing.

## On using AI tools

AI tools are fine. Use whatever helps you write better code. We don't care* how you got to the solution — we care whether the solution is good.

The standard is the same regardless of how the code was written: **you are responsible for every line of your contribution.** If you can't explain why a line is there and why it's correct, it shouldn't be in your PR.

**Specifically:**

- If you use AI to help write code, you must understand every line of the diff you're submitting.
- Do not paste raw AI output into issues, PRs, or comments. If you use AI to help draft text, rewrite it in your own words.
- Remove AI-generated footers, co-author attributions, and "Generated with..." signatures before submitting. Their presence tells us you didn't review your own submission carefully enough to notice them.
- Automated submissions — bots or agents posting PRs without meaningful human review — will be treated as spam.

## Consequences

We'd rather help you improve a PR than close it. But we can't review work that wasn't reviewed by the person submitting it.

- **First time:** You'll get a warning and a chance to fix the PR.
- **Second time:** The PR will be closed.
- **Repeated offenses:** May be reported to GitHub as spam.