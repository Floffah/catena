# Catena

Open source modern git server. My attempt at making a GitHub-like system. Catena is experimental for now. The goal is to build a fast, modern, self-hostable social Git platform with a great developer experience and a clean web UI.

Links:
- [Production Instance](https://www.oncatena.com)
- [API Docs](https://api.oncatena.com/docs)

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
## Contents

- [Roadmap](#roadmap)
- [Goals](#goals)
- [Deployment](#deployment)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Roadmap

in rough order of priority:
- [x] Authentication [FR-003]
- [x] Create repo [FR-001]
- [x] Clone / fetch / push to repo with git-http-backend [FR-002]
- [x] Git authentication [FR-003]
- [x] Web repo browser [FR-006, FR-007]
- [x] Product API foundation [FR-005]
- [x] Database foundation [FR-004]
- [ ] Manage account & tokens [FR-003]
- [ ] Issues <- at this point a production deployment will be created [FR-008]
- [ ] External project management integrations [FR-009]
- [ ] Pull requests [FR-010]
  - [ ] Diff rendering [FR-011]
  - [ ] Code review UI [FR-012]
  - [ ] Review approvals and requested changes [FR-012]
  - [ ] Merge strategies [FR-013]
- [ ] Notifications [FR-034]
  - [ ] Mentions [FR-034]
- [ ] Project/activity feeds <- at this point the project is considered minimally viable [FR-034]
- [ ] Webhooks [FR-017]
- [ ] Repository events [FR-017]
- [ ] Post-receive processing [FR-017]
- [ ] Third-party API and automation surface [FR-019]
- [ ] CI/CD integration model [FR-014]
- [ ] First-party CI/pipeline runner [FR-014]
- [ ] External CI integrations [FR-015]
- [ ] External deployment integrations [FR-015]
- [ ] Status checks [FR-016]
- [ ] Release management <- at this point the project is considered at v1.0 [FR-035]
  - [ ] Release notes [FR-035]
  - [ ] Changelog generation [FR-036]
  - [ ] Milestone management [FR-036]
  - [ ] Release artifact management [FR-037]
- [ ] `catena` CLI tool [FR-018]
  - [ ] CLI login [FR-018]
  - [ ] Git credential helper integration [FR-018]
  - [ ] CLI repository management [FR-018]
  - [ ] CLI token management [FR-018]
- [ ] Proper docs [FR-043]
- [ ] Production volume mounting [FR-044]
- [ ] Advanced object storage [FR-028]
  - [ ] LFS support [FR-028]
  - [ ] Backups [FR-028]
- [ ] Repository sharding strategy [FR-027]
- [ ] Multi-replica Git serving strategy [FR-027]
- [ ] Background job runner [FR-045]
- [ ] Metrics and analytics [FR-046]
- [ ] Rate limiting [FR-047]
- [ ] Secret scanning [FR-030]
- [ ] Abuse prevention [FR-030]
- [ ] Error tracking [FR-048]
- [ ] `catenarc` repository configuration [FR-020]
- [ ] Push-based repository creation [FR-021]
- [ ] Fork support [FR-022]
- [ ] AI pull request review [FR-023]
- [ ] Agent task offload [FR-024]
- [ ] Agent-powered triage [FR-025]
- [ ] code.storage snippet sharing [FR-026]
- [ ] Organisation management [FR-029]
- [ ] Team management [FR-029]
- [ ] Repository collaborators and access control [FR-029]
- [ ] Branch rules [FR-030]
- [ ] Pull request requirements [FR-030]
- [ ] Repository visibility controls [FR-030]
- [ ] System-wide code search [FR-031]
- [ ] Trending explore page [FR-032]
- [ ] Algorithmic explore page [FR-032]
- [ ] User profiles [FR-033]
- [ ] Following users [FR-033]
- [ ] Starring repositories [FR-033]
- [ ] Audit logs [FR-038]
- [ ] Dependency management [FR-039]
- [ ] Dependency graphs [FR-039]
- [ ] Vulnerability scanning [FR-040]
- [ ] Repository import [FR-041]
- [ ] Repository export [FR-041]
- [ ] Repository mirroring [FR-041]
- [ ] Distributed Catena backend instances [FR-042]
- [ ] Global repository network [FR-042]

## Goals

In no order. Features marked with \* are differentiating features. All of these features are core to the vision of Catena
- **FR-001** Host filesystem-backed bare Git repositories.
- **FR-002** Provide Git-over-HTTP clone, fetch, and push.
- **FR-003** Authenticate users with Clerk and manage Catena Git access tokens.
- **FR-004** Store product data in Postgres.
- **FR-005** Expose an OpenAPI HTTP/JSON product API.
- **FR-006** Provide a Next.js web UI for repository creation and browsing.
- **FR-007** Support repository README, tree, branch, commit, and path resolution views.
- **FR-008** Have an intentionally limited issue tracker for triage NOT project management.
- **FR-009** \* Provide first-class integration with external project management tools like grapharc or Linear.
- **FR-010** Provide a pull request system.
- **FR-011** \* Use diffs.com for diff rendering.
- **FR-012** \* Use trees.software for code review UI.
- **FR-013** Support squash, rebase, and merge commit merge strategies.
- **FR-014** Provide native CI/CD pipelines with provided and self-hostable runners.
- **FR-015** \* Provide first-class integrations with external CI/CD and deployment tools like Depot, Blacksmith, Render, Fly.io, and Vercel.
- **FR-016** Support status checks and a GitHub-like commit status API for CI/CD integration.
- **FR-017** Support webhooks, repository events, and post-receive processing for custom workflows.
- **FR-018** Provide a `catena` CLI tool for repository management, token management, and Git credential helper integration.
- **FR-019** Provide a robust API for third-party integrations and automation.
- **FR-020** \* Allow `catenarc` configuration files for file-based repository configuration in addition to UI configuration.
- **FR-021** \* Support push-based repository creation in addition to UI-based repository creation.
- **FR-022** Support repository forks.
- **FR-023** Provide first-class AI review for pull requests.
- **FR-024** Provide first-class agent task offload.
- **FR-025** \* Provide agent-powered triage and issue management.
- **FR-026** \* Provide agentic support features like code.storage for quick iteration and snippet sharing.
- **FR-027** Support repository sharding and multi-replica Git serving strategies for scalability.
- **FR-028** Support advanced object storage strategies for LFS, backups, and large repositories.
- **FR-029** Provide organisation and team management features for access control and collaboration.
- **FR-030** Provide repository security features including branch rules, pull request requirements, secret scanning, visibility, abuse prevention, and other security features.
- **FR-031** Provide system-wide code search.
- **FR-032** Provide trending and algorithmic explore pages.
- **FR-033** Provide user profiles and social features like following, starring, and forking.
- **FR-034** Provide notifications and activity feeds for PRs, issues, reviews, comments, repository activity, releases, and followed users.
- **FR-035** Provide release management features including release notes.
- **FR-036** Provide changelog generation and milestone management.
- **FR-037** Provide release artifact management.
- **FR-038** Provide audit logs.
- **FR-039** Provide dependency management and dependency graphs.
- **FR-040** Provide vulnerability scanning.
- **FR-041** \* Support import, export, and mirroring to and from other Git hosting platforms like GitHub, GitLab, and Bitbucket.
- **FR-042** \* Allow distribution of Git servers so users can run their own Catena backend instance and participate in the same global network of repositories, users, and organisations.
- **FR-043** Provide proper user, API, operator, and contributor documentation.
- **FR-044** Support production storage deployment with mounted Git volumes.
- **FR-045** Provide a background job runner for async platform work.
- **FR-046** Provide metrics, analytics, and observability.
- **FR-047** Provide rate limiting.
- **FR-048** Provide error tracking.

## Deployment

Catena is not currently production ready. See [CONTRIBUTING.md](CONTRIBUTING.md) if you're interested in contributing or running a development instance.