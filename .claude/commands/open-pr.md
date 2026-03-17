Open a pull request using the standard project template.

## Pre-flight checks
Before opening the PR:
- All tests pass: `go test ./...`
- No vet issues: `go vet ./...`
- The branch is pushed to remote: `git push -u origin <branch>`
- There is a related issue to close (if applicable)

## Steps

1. Gather:
   - A short PR title following the convention below
   - The issue number being closed (if any)
   - A bullet list of what changed and why
   - Whether any adapter interfaces changed and if the change is breaking

2. Open the PR:

```bash
gh pr create \
  --title "<title>" \
  --body "$(cat <<'EOF'
<filled pull_request_template.md body>
EOF
)"
```

## Title conventions
Use the same prefix as the related commit(s):
- `feat: <what was added>`
- `fix: <what was fixed>`
- `refactor: <what was restructured>`
- `test: <what test coverage was added>`
- `docs: <what was documented>`

Keep titles under 72 characters. Do not end with a period.

## Notes
- Never open a PR directly to `main` for a breaking adapter interface change without prior discussion in an issue.
- The PR checklist in `.github/pull_request_template.md` must be completed before requesting review.
