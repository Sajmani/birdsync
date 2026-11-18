# Pull requests

Part of River's method. Search `^## ` for the sections; read only the one you need.

## Requirements impact

**Every pull request that changes the River module carries a "Requirements impact" section.**
Source control approves River changes, so the reviewer must see what the change moves. Build the
section from:

```
river diff --base main --format pr
```

(`--base` defaults to the merge base with `main`.) It reports, by area: requirements added,
changed, withdrawn, or removed; imported sources re-pinned and what changed upstream; composed
packages recomposed and how, and those now stale; precedence rows; criteria, and requirements
newly without one; guardrails, criteria newly without one, and new watched-failing lines; and
decisions added, superseded, or amended, with the existing decisions affecting each moved
requirement. Long lists are capped; `river diff --base main` gives them in full.

**Put your triage first**, above the report, in three parts:

1. **Needs your call:** one line each, saying why. Typical entries: a requirement whose meaning
   changed; a new exclusion or precedence change; a criterion weakened or withdrawn; a guardrail
   removed; coverage lost (a requirement with no criterion, or a criterion with no guardrail); an
   existing decision the change may contradict.
2. **Unaffected:** one summary line per area, such as "13 decisions affect
   `tech.layoutVersion`; none is contradicted: the change moves a list, not the rule".
3. **Amended here:** what this pull request itself supersedes or amends, and why.

Then the full report, as `river diff` writes it.

**Stop and ask before opening the pull request** when more than a handful of entries need the
owner's call: a change that needs many calls is probably several changes.

**For a direct push**, with no pull request, show the owner the same output and get their answer
before pushing.

## Keeping pull requests small

- **Keep requirements changes in their own pull requests**, apart from the code that implements
  them, so the review of what must be true is not buried in a diff of how.
- **Split a heavily cited item.** When `river diff` hints that a changed item is affected by many
  decisions, splitting it into smaller items keeps later reviews short.
- Sums ignore reflowed text, so rewrapping a paragraph moves nothing.
