# Catena

Open source modern git server. My attempt at making a GitHub-like system. Catena is experimental for now. The goal is to build a fast, modern, self-hostable social Git platform with a great developer experience and a clean web UI.

## Roadmap

in rough order of priority:
- [x] Authentication
- [x] Create repo
- [x] Clone / push to repo (git-http-backend)
- [x] Git authentication
- [ ] Web repo browser
- [ ] Manage account & tokens
- [ ] Issues <- at this point a production deployment will be created
- [ ] Pull requests
  - [ ] Code review
  - [ ] Review approvals and requested changes
  - [ ] Merge strategies
- [ ] Notifications
  - [ ] Mentions
- [ ] Project/activity feeds <- at this point the project is considered minimally viable
- [ ] Webhooks
- [ ] Repository events
- [ ] Post-receive processing
- [ ] CI/CD integration model
- [ ] First-party CI/pipeline runner
- [ ] External CI integrations
- [ ] Status checks
- [ ] Release management <- at this point the project is considered at v1.0
- [ ] `catena` CLI tool
  - [ ] CLI login
  - [ ] Git credential helper integration
  - [ ] CLI repository management
  - [ ] CLI token management
- [ ] Proper docs
- [ ] Production volume mounting
- [ ] Advanced object storage
  - [ ] LFS support
  - [ ] Backups
- [ ] Repository sharding strategy
- [ ] Multi-replica Git serving strategy
- [ ] Background job runner
- [ ] Metrics and analytics
- [ ] Rate limiting
- [ ] Secret scanning
- [ ] Abuse prevention
- [ ] Error tracking
