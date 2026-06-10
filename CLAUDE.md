# go_learn — Project Context

This is a **personal learning sandbox** for practicing Go. The goal is not to ship a product, but to build something small end-to-end while learning idiomatic Go, the Gin web framework, and testing patterns (table-driven tests, `httptest`, mocks/fakes, dependency injection for testability).

## Learning goals

- Idiomatic Go project layout (`cmd/`, `internal/`, `pkg/`)
- Gin: routing, middleware, handler patterns, request validation, error handling
- Layered architecture: **handler → service → repository**
- Testing: unit tests (services), HTTP tests (`net/http/httptest` + Gin), table-driven tests, integration tests
- Error handling and proper HTTP status codes
- Configuration (env vars, flags)
- Possibly: SQLite/Postgres later, simple auth, structured logging

## Chosen project: Personal Bookmark Manager

- **Surface:** CRUD on bookmarks (URL, title, tags, notes) + search/filter by tag or text
- **Why:** Real query/filter use case, natural fit for handler/service/repository layering, plenty of room for table-driven tests.
- **Roadmap:**
  1. ✅ In-memory repo, model, Create + List endpoints, tests for each layer (current step)
  2. Get by ID, Update, Delete
  3. Filter by tag, simple text search
  4. Swap in-memory repo for SQLite (same interface)
  5. Optional: structured logging, env-based config, simple auth

## Working style — IMPORTANT

**Claude does NOT write the implementation.** This is a learning project. The user writes the code; Claude **explains concepts, then reviews and corrects** what the user wrote.

- Explain Go idioms, package responsibilities, and design tradeoffs in prose.
- Tiny illustrative snippets (3–5 lines) are OK to clarify a concept.
- Full files / full functions / boilerplate the user could write themselves → **do not write them**. Wait for the user.
- Scaffolding commands the user explicitly delegates (`go mod init`, `mkdir`, `go get`) are fine — those aren't learning material.
- After the user writes code, review for: correctness, idiomatic style, naming, error handling, missing tests, layering violations.
- Build incrementally — small steps, each one testable. Tests alongside code, not as an afterthought.

## Tech baseline (proposed, not locked in)

- Go (latest stable)
- Gin for HTTP
- Standard library `testing` + `net/http/httptest`
- `testify/assert` is fine if useful, but try to learn the stdlib idioms first
- Storage: in-memory map → SQLite (`modernc.org/sqlite` for cgo-free) → optionally Postgres
- Config via env vars

## Next session — start here

1. Confirm which of the 3 projects to build (recommendation: bookmark manager).
2. Scaffold the directory structure (`cmd/server/main.go`, `internal/handler`, `internal/service`, `internal/repository`, `internal/model`).
3. `go mod init` and pull in Gin.
4. First feature end-to-end: a single endpoint with handler + service + repo + tests for each layer.

## Git

This directory is **not currently a git repo**. Before the first commit, follow the global rule in `~/.claude/CLAUDE.md`: validate `git config user.name`, `user.email`, and remotes, and confirm with the user which identity (personal vs. work) to use. Personal project → personal identity expected, but always confirm.