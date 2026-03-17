Write an implementation plan for the requested feature or task and save it to `.plan/`.

## Steps

1. Invoke the `superpowers:writing-plans` skill to structure the plan properly.
2. Derive a short filename from the feature name (kebab-case, no spaces):
   - e.g. `thread-store-interface.md`, `backfill-existing-prs.md`
3. Write the plan to `.plan/<filename>.md`.
4. Confirm the file path to the user when done.

## File location

Always write to `.plan/<filename>.md` — never to the project root, never inline in the conversation.

## When to promote a plan to a GitHub issue

`.plan/` is local and gitignored. If the plan represents a significant architectural decision or work that a collaborator needs visibility into, create a GitHub issue from it using `/create-issue` and reference the plan file in the issue body.

## After the plan is written

Once the plan file is saved, present the user with these three options exactly:

---

**What would you like to do next?**

1. **Clear context and start implementation** — recommended for large plans. Starts a fresh session with the plan file loaded so the context isn't polluted by this planning conversation.
2. **Start implementation now** — begins executing the plan immediately in this session.
3. **Ask me questions** — you have follow-up questions or want to refine the plan before doing anything.

---

Wait for the user to choose before doing anything else.

- If **1**: instruct the user to start a new session and reference the plan with `@.plan/<filename>.md`, then stop.
- If **2**: invoke `superpowers:executing-plans` and begin implementation against the saved plan file.
- If **3**: ask your questions, update `.plan/<filename>.md` with any changes, then present the three options again.

## When to clean up

Delete the plan file once the work it describes is merged. Plans are scratch space, not long-term documentation.
