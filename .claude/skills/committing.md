---
name: committing
description: Use when about to create a git commit in this repository — enforces conventional commit format and project-specific staging rules
type: skill
---

## Commit message format

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<optional scope>): <short description>

<optional body — wrap at 72 chars>
```

### Types
| Type | When to use |
|------|-------------|
| `feat` | New behaviour or webhook event support |
| `fix` | Bug fix |
| `refactor` | Code restructure with no behaviour change |
| `test` | Adding or updating tests or fixtures |
| `docs` | CLAUDE.md, FIXTURES.md, or other documentation |
| `chore` | Dependency updates, go.mod changes, config |

### Scopes (optional but encouraged)
Use the package or concept being changed:
- `webhook` — event parsing / routing
- `slack` — Slack API interaction
- `thread` — ThreadStore interface or usage
- `config` — ConfigStore interface or usage
- `logger` — Logger interface or usage
- `testdata` — fixture files

### Examples
```
feat(webhook): add support for pullrequest:comment_updated event
fix(thread): return error when ThreadStore TTL write fails
test(testdata): add fixture for pullrequest with no reviewers
docs: update FIXTURES.md with commit_status resolution notes
```

## What to stage

- Stage specific files by name — never `git add .` or `git add -A`
- Always verify the diff before committing: `git diff --staged`
- Do not commit:
  - `.env` or any file containing secrets
  - Generated files that are not checked in by convention
  - Unrelated changes bundled with the intended change

## Checklist before committing

- [ ] `go test ./...` passes
- [ ] `go vet ./...` passes
- [ ] If interfaces changed → `CLAUDE.md` is updated
- [ ] If fixtures changed → `testdata/webhooks/FIXTURES.md` is updated
