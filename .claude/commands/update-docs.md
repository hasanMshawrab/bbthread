Review recent changes and update all affected documentation.

## When to update CLAUDE.md
- A new adapter interface is added or an existing one changes signature
- A new webhook event type is supported
- The core flow changes (e.g. new steps in the thread-lookup or backfill logic)
- The Slack API usage changes (e.g. new API methods used)
- A new top-level directory or architectural concept is introduced

## When to update testdata/webhooks/FIXTURES.md
- A new fixture file is added to `testdata/webhooks/`
- An existing fixture is modified in a way that affects what it tests
- A new intentional design decision is embedded in the fixtures (shared IDs, edge-case fields, etc.)

## Steps

1. Read the current `CLAUDE.md` and identify any sections affected by recent changes.
2. Read `testdata/webhooks/FIXTURES.md` and check if any fixture files were added or changed.
3. Edit only the sections that are out of date — do not rewrite unaffected sections.
4. Keep descriptions factual and concise. Do not add implementation advice that belongs in code comments.
