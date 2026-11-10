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
Reaffirms: <id>      (only when confirming that an earlier decision under review still holds)

### Decided

<What was decided. If no reason was given, say so rather than inventing one.>
```

Then `### Decided`, the decision itself, before the reasons; `---` ends the entry. If no reason
was given, say so rather than supplying one. Other fields may follow these, and `Status:` may
carry a note after a semicolon.

## Rules

- **Content is append-only; format is migratable.** Never edit an entry's decision, `Decided by:`,
  `Date:`, or reasons: a decision that turns out badly is superseded or amended by a new entry,
  because the record of having been wrong is worth keeping. By hand, only an earlier entry's
  `Status:` line changes, to name what superseded, amended, or reaffirmed it. Across River major
  versions, `river upgrade` may migrate entry formatting, repoint links, and backfill required
  sections, recording in the commit that it did.
- **Accepting a gap is a decision:** `Affects:` names the unmet requirement, and `Expires:` is
  required. An accepted gap with no review date is indistinguishable from a forgotten one.
- **Composed-package approvals live in `requirements/approvals.md` from River 2.0**, written by
  `river approve --by <owner>` after the owner approves the recomposed state, rather than in a
  River module's `decisions.md` (`misplaced-approval` if `Approves:` appears here under 2.0; `river
  upgrade` moves 1.x `Approves:` entries into `approvals.md`; see [`module.md`](module.md), *Keeping
  the River module's files*). **Never run `river approve` unless the owner has approved that
  state.**
- **From River 2.0, every in-force decision affecting source items records what it was decided
  against:** a `### Sums` section at the end of the entry, from `river sum <item>...` for the
  source items in `Affects:` (`decision-missing-sums` if missing; `river upgrade` backfills
  baseline sums onto 1.x decisions with `Recorded by \`river upgrade\` on YYYY-MM-DD, not when
  decided.`). When an affected item later changes, is withdrawn, or is removed, `river check`
  notes the decision (`decision-review`) for a person to confirm or supersede. Every name in
  `Affects:` must exist.
- You may draft the entry; the person named in `Decided by:` must have actually decided it.

## Interpretation decisions

**An interpretation decision says how to read an imported item, without editing it.** Imported
text is never changed, so what the project concludes about it is a decision whose `Affects:` names
the imported items it reads, linked from the import's entry in `imports.md`. It is one of
composition's inputs ([`composition.md`](composition.md), *The three inputs*). Three kinds:

- **The class an item applies to**, when the source doesn't name one: `Applies to: wcag.WebPage`
  for the items listed in `Affects:` ([`composition.md`](composition.md), *Projecting horizontal
  requirements*).
- **The tier of one item**, when a Governing or Mandatory source only recommends: the decision
  says who decided and why, and the tier itself goes in a row of `precedence.md`'s item table that
  links the decision ([`requirements.md`](requirements.md), *Precedence*). Don't repeat the tier in
  the decision: it is recorded once.
- **A resolved ambiguity**, such as the checkable reading of "about one per second".

Record it when the source is imported or when composition first needs it, before composing on it,
and, like any decision, only with a person's say.

## Reviewing a decision

**A decision stays in force only while what it affects is unchanged.** When an item in its
`Affects:` changes, is withdrawn, or is removed, the decision is due for review: the CLI notes it
(`decision-review`), and the next recomposition asks the person whether it still holds
([`composition.md`](composition.md), *Reviewing decisions at recomposition*). If it doesn't hold,
supersede it.

**If it holds, reaffirm it in a new entry**, never by editing the old one: `Reaffirms: <id>`, its
own `Decided by:` naming the person who confirmed it, an `Affects:` naming the same items, and a
fresh `### Sums` of them from `river sum <item>...`. Change the old entry's `Status:` to add
`; reaffirmed by <new id>`, and nothing else in it. The CLI stops noting the old decision while
the new entry's sums match the items, and notes both again when they move. From River 2.0,
`Reaffirms:` naming no entry fails (`reaffirms-unknown`).

## Seeking

`awk '/^## /{id=$2} /^Affects:.*\.offlineMode(,|$)/{print id}' decisions.md` lists what touches a
requirement, and `awk '/^## /{id=$2} /^Status: (amended|superseded)/{print id": "$0}'` lists what
is no longer in force as written.
