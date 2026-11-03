# Recording a decision

Part of River's method. Search `^## ` for the sections; read only the one you need.

## Format

`river/requirements/decisions.md` opens with `package decisions`, then holds one entry per
decision, oldest first, **appended at the end**. Each is `## <id>`, a bold line stating the
decision, then one field per line:

```
---

## <descriptiveCamelCaseId>

**<The decision, in one sentence>**

Date: 2026-10-07
Decided by: <person>
Status: accepted | amended by <id>, <id> | superseded by <id>
Supersedes: <id>     (only when it replaces an earlier decision in whole)
Amends: <id>, <id>   (only when it changes earlier decisions in part)
Affects: product.x   (the requirements it changes or settles, not ones it merely cites)
Expires: 2027-01-01  (only for an accepted gap: when it must be looked at again)
Approves: app, sha256:<hash>  (only for approving a composed version)

### Decided

<What was decided. If no reason was given, say so rather than inventing one.>
```

Then `### Decided`, the decision itself, before the reasons; `---` ends the entry. If no reason
was given, say so rather than supplying one. Other fields may follow these, and `Status:` may
carry a note after a semicolon.

## Rules

- **Append only, with one exception:** when a new entry supersedes or amends an earlier one,
  change the earlier entry's `Status:` line to say so, and nothing else in it, so that a reader of
  the old entry learns it is no longer the whole story. A decision that turns out badly is
  superseded, never edited: the record of having been wrong is worth keeping.
- **Accepting a gap is a decision:** `Affects:` names the unmet requirement, and `Expires:` is
  required. An accepted gap with no review date is indistinguishable from a forgotten one.
- **An approval of a composed version** adds `Approves: <name>, sha256:<hash>`, from
  `river hash <package-dir>` ([`bundle.md`](bundle.md), *Keeping the bundle's files*). **Only the
  approving person adds this line.**
- You may draft the entry; the person named in `Decided by:` must have actually decided it.

## Seeking

`awk '/^## /{id=$2} /^Affects:.*\.offlineMode(,|$)/{print id}' decisions.md` lists what touches a
requirement, and `awk '/^## /{id=$2} /^Status: (amended|superseded)/{print id": "$0}'` lists what
is no longer in force as written.
