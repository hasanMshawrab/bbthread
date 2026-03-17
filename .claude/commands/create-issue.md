Create a GitHub issue using the appropriate template.

## Steps

1. Determine the issue type:
   - **Bug** → use `.github/ISSUE_TEMPLATE/bug_report.md`
   - **Feature / improvement** → use `.github/ISSUE_TEMPLATE/feature_request.md`

2. Collect the required information by asking the user for any missing fields from the relevant template.

3. Create the issue:

```bash
gh issue create \
  --title "<concise title>" \
  --body "$(cat <<'EOF'
<filled template body>
EOF
)" \
  --label "<bug|enhancement>"
```

## Title conventions
- Bug: `fix: <what is broken>` — e.g. `fix: build status event not threaded when PR predates integration`
- Feature: `feat: <what is being added>` — e.g. `feat: support pullrequest:comment_updated event`

## Notes
- Do not create an issue for work that already has an open issue. Search first: `gh issue list --search "<keywords>"`
- If the issue touches an adapter interface, make sure the `Adapter interface changes` section in the template is filled in.
