# Checks

Part of River's method. Search `^## ` for the sections; read only the one you need.

## The river CLI

**River's checks are one Go binary, `river`, the same for every project.** Run `river check`
(or `river check <bundle>`) wherever the project already runs checks: the agent's edit loop,
pre-commit, CI, review, and release. It needs no toolchain in the adopting project and is not
translated into the project's language; the shared cases in River's `river/check/testdata/` are
its contract. Record the CLI version you run in `bindings.md`, beside the skill's pin, so a change
to River's checks is visible.

Checks of your own, beyond River's, go in your project's own language and test suite, which
cannot be skipped and add no toolchain.

**Before River's checks replace yours, list what yours enforced.** River's checker counts some
things a project's own checks may have failed (a requirement citing no goal, say), because a
brownfield project may start that way. If your project was stricter, keep the stricter rule as
your own check; otherwise swapping checkers is how a rule disappears without anyone deciding it
should.

## What a bundle's checks cover

Run them before calling a change done. At minimum a bundle should check that:

- every requirement's goal citation resolves, and every criterion's requirement citation does
- identifiers are well formed and unique within their package
- a pinned source still matches its recorded hash, by whatever means its style allows
- counts that are claimed are computed, not maintained by hand
- the instruction file and every line in the bundle stay within the limits in
  [`bundle.md`](bundle.md) (*Write every agent-read file to be read in part*)

**Every check must have been watched failing.** Plant a violation, see it reported, remove it.

## Notes, not failures

**Some symptoms are notes, never failures:** a goal cited by one requirement, a goal whose
requirements all cite one other goal too, a goal longer than a paragraph, and a goal that contains
code or names another item. Each can be right. Each is where to look first.

What a check cannot do is tell you whether a requirement genuinely *serves* the goal it cites.
Service is a semantic relation, not a measurement. A green run is not a validated one.
