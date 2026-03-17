# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Go module: `github.com/hasanMshawrab/bitslack`
Go version: 1.24.2

A Go **library** that receives Bitbucket webhook events and forwards them to Slack as **threaded messages** — all events for a given PR appear as replies under the original message rather than as new top-level messages.

Consumers embed this library into their own server and wire it up by providing adapter implementations. The library defines interfaces; callers supply the concrete backends.

## Common Commands

```bash
make check                            # Full check: build + vet + lint + arch-lint + test
make build                            # Build all packages
make test                             # Run all tests (excluding e2e)
make test-unit                        # Run unit tests only (internal/)
make test-integration                 # Run integration tests (handler)
make lint                             # Run golangci-lint
make lint-fix                         # Run golangci-lint with auto-fix
make arch-lint                        # Run architecture dependency linter
make fmt                              # Format code (goimports + golines)
make tools                            # Install pinned dev tools
go test ./internal/... -run TestName  # Run a single test by name
```

## File System Structure

```
bitslack/
├── bitslack.go          # Public API: Config struct, New() constructor, Handler()
├── adapter.go           # Interface definitions: ThreadStore, ConfigStore, Logger
├── handler.go           # http.Handler — receives webhooks and drives the core flow
│
├── internal/
│   ├── bitbucket/       # Bitbucket REST API client (PR resolution from commit hash)
│   ├── slack/           # Slack API client (chat.postMessage, chat.update)
│   ├── event/           # Webhook event types, JSON parsing, routing by event key
│   └── format/          # Slack message formatting (opening message, reply text)
│
├── examples/
│   └── server/
│       ├── main.go          # Reference server wired with concrete adapters
│       └── docker-compose.yml  # Full stack for local E2E testing
│
├── testdata/
│   └── webhooks/
│       ├── FIXTURES.md          # Explains fixture design decisions
│       ├── pullrequest/         # One JSON file per pullrequest:* event
│       └── commit_status/       # One JSON file per repo:commit_status_* event
│
├── .claude/
│   ├── commands/        # Custom slash commands: /plan, /create-issue, /open-pr, /update-docs
│   └── skills/          # Superpowers skills: committing
│
├── .github/
│   ├── ISSUE_TEMPLATE/  # bug_report.md, feature_request.md
│   └── pull_request_template.md
│
├── .plan/               # Local planning scratch space — gitignored
├── .gitignore
├── go.mod
└── CLAUDE.md
```

### Key boundaries

- **Public surface** (`bitslack.go`, `adapter.go`, `handler.go`) — everything a consumer needs to import and wire up. Keep this minimal and stable.
- **`internal/`** — all implementation details. Nothing in `internal/` is importable by consumers. Each sub-package has a single clear responsibility.
- **`examples/`** — the only place that may use concrete third-party adapter implementations. The core library never depends on them.
- **`testdata/`** — Go test convention; files here are accessible via `os.Open("testdata/...")` in tests without any path manipulation.

## Architecture

### Adapter / Plugin Model

The library is backend-agnostic. Consumers construct the core engine by injecting adapters that satisfy these interfaces:

- **ConfigStore** — provides repo→channel mapping and username→Slack user ID lookup. The library only calls lookup methods (`GetChannel`, `GetSlackUserID`); it never loads or caches data itself. The consumer controls their own data lifecycle — preloading, caching, or fetching on demand is entirely up to them.
- **ThreadStore** — stores and retrieves the PR→Slack thread `ts` mapping. Needs TTL support (30-day expiry per PR). Could be backed by Redis, Memcached, an in-process map, etc.
- **Logger** — structured logging with three methods: `Info(message string)`, `Warn(message string)`, `Error(message string)`. If none is provided, the library defaults to a no-op logger.

The library ships no concrete adapter implementations — those live in the caller's codebase or in separate companion packages.

### Opening Message Format

The first message posted for a PR (either on `pullrequest:created` or backfilled) must display:
- **PR title** (bold)
- **Author** — Slack @mention if a username mapping exists, otherwise plain Bitbucket username
- **Repository** name
- **Reviewers** — each as a Slack @mention if mapped, otherwise plain username

### Opening Message Updates

The opening message is a live document — it is edited (via `chat.update`) to stay in sync with PR state changes:

- `pullrequest:updated` — if the title or reviewer list changed, update the opening message in place
- **Adding a reviewer** — edit the message to add their @mention; Slack will automatically notify them (no separate notification needed)
- **Removing a reviewer** — edit the message to remove their @mention; Slack will not notify them of the removal. If they have not yet engaged with the thread (no reply, no click-through), they will stop receiving future thread notifications. If they have already engaged, Slack marks them as a thread follower and they will continue to receive updates regardless — this is a known Slack limitation.

### Core Flow

1. Caller's HTTP server receives a Bitbucket webhook and passes the raw event to the library.
2. The library identifies the PR (see "Build Status Events" below).
3. Look up the thread `ts` for that PR via `ThreadStore`.
4. If no `ts` exists (new PR **or** an existing PR that predates the integration):
   - Call the Bitbucket API to fetch full PR details (`GET /repositories/{workspace}/{repo}/pullrequests/{id}`)
   - Post a synthetic opening message to Slack → store the returned `ts` via `ThreadStore`
   - If either step fails, log the error and drop the event gracefully (no panic, no partial state)
5. Post the triggering event as a reply using `thread_ts`.

### Build Status Events

Bitbucket `repo:commit_status_created` / `repo:commit_status_updated` events do **not** include a PR ID — only a commit hash. To resolve the PR, call the Bitbucket API:

```
GET /repositories/{workspace}/{repo}/commit/{hash}/pullrequests
```

### Slack Integration

Uses `chat.postMessage` with the `thread_ts` field to post replies into an existing thread.

The caller provides a Slack Bot Token (`xoxb-...`) at construction time — the library has no opinion on how it is stored or retrieved:

```go
client := bitslack.New(bitslack.Config{
    SlackToken: "xoxb-...",
    // adapters...
})
```

Required OAuth scopes: `chat:write`. Add `chat:write.public` if the bot needs to post to channels it hasn't been explicitly invited to.

### Testing Strategy

- **Unit** — test event parsing, message formatting, and thread-lookup decision logic using mock adapters. This is the bulk of the test suite.
- **Integration** — test concrete adapter implementations against real infrastructure spun up via docker-compose (e.g. a Redis-backed `ThreadStore`).
- **E2E** — lives in `examples/`, fires real Bitbucket-shaped webhook payloads at a running instance and asserts against a Slack API stub.

An `examples/` directory provides a reference server wired with concrete adapters and a `docker-compose.yml` for running the full stack locally. This doubles as the E2E test harness.

### Supported Webhook Events

**Pull Request**
- `pullrequest:created`
- `pullrequest:updated`
- `pullrequest:approved`
- `pullrequest:unapproved`
- `pullrequest:fulfilled` (merged)
- `pullrequest:rejected` (declined)
- `pullrequest:comment_created`

**Build Status**
- `repo:commit_status_created`
- `repo:commit_status_updated`
