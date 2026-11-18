# Recording a decision

Part of River's method. Search `^## ` for the sections; read only the one you need.

## Format

`river/requirements/decisions/` is the decisions package: each file opens with `package
decisions`, and holds entries grouped by topic, usually the requirement area they mainly affect.
Append a new entry at the end of its topic's file, or start a topic file. Each entry is `## <id>`,
a bold line stating the decision, then one field per line:

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

Reaffirms: <id>      (optional: an earlier decision this one confirms still holds)

### Decided

<What was decided. If no reason was given, say so rather than inventing one.>
```

Then `### Decided`, the decision itself, before the reasons; `---` ends the entry. Other fields may
follow these, and `Status:` may carry a note after a semicolon. Every name in `Affects:`,
`Supersedes:`, `Amends:`, and `Reaffirms:` must exist; the checker fails one that doesn't.

## Rules

- **Content is append-only; format is migratable.** Never edit an entry's decision, `Decided by:`,
  `Date:`, or reasons: a decision that turns out badly is superseded or amended by a new entry. By
  hand, only an earlier entry's `Status:` line changes, to name what superseded or amended it. An
  entry may move between topic files; `river links -w` repoints links to it. `river upgrade` may
  migrate formatting across major versions.
- **Accepting a gap is a decision:** `Affects:` names the unmet requirement, and `Expires:` is
  required.
- **A decision records what was decided, not sums of what it was decided against.** Code review
  sees what moved: `river diff --base` lists, for each requirement a change moves, the in-force
  decisions whose `Affects:` names it ([`pull-requests.md`](pull-requests.md)). A `### Sums` block
  left from River 2.x is noted (`decision-sums-ignored`); `river upgrade` removes it.
- You may draft the entry; the person named in `Decided by:` must have actually decided it.

## Interpretation decisions

**An interpretation decision says how to read an imported item, without editing it.** Imported
text is never changed, so what the project concludes about it is a decision whose `Affects:` names
the imported items it reads, linked from the import's entry in `imports/`. It is one of
composition's inputs ([`composition.md`](composition.md), *The three inputs*). Three kinds:

- **The class an item applies to**, when the source doesn't name one: `Applies to: wcag.WebPage`
  for the items listed in `Affects:` ([`composition.md`](composition.md), *Projecting horizontal
  requirements*).
- **The tier of one item**, when a Governing or Mandatory source only recommends: the decision
  says who decided and why, and the tier itself goes in a row of the precedence item table that
  links the decision ([`requirements.md`](requirements.md), *Precedence*).
- **A resolved ambiguity**, such as the checkable reading of "about one per second".

Record it when the source is imported or when composition first needs it, and only with a
person's say.

## When an affected item changes

When a change moves a requirement that a decision affects, the pull request's requirements-impact
section lists the decision. Ask the person whether it still holds. If it doesn't, supersede it. If
it does, nothing need be recorded; a `Reaffirms: <id>` entry is optional, for when the confirmation
itself is worth keeping.

## Seeking

`grep -rn '^Affects:.*offlineMode' river/requirements/decisions/` lists what touches a
requirement, and `grep -rn '^Status: \(amended\|superseded\)' river/requirements/decisions/` lists
what is no longer in force as written.
