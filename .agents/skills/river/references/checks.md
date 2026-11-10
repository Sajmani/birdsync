# Checks

Part of River's method. Search `^## ` for the sections; read only the one you need.

## The river CLI

**River's checks are one Go binary, `river`, the same for every project.** Run `river check`
(or `river check <module>`) wherever the project already runs checks: the agent's edit loop,
pre-commit, CI, review, and release. It needs no toolchain in the adopting project and is not
translated into the project's language; the shared cases in River's `river/check/testdata/` are
its contract. Record the CLI version you run in `bindings.md`, beside the skill's pin, so a change
to River's checks is visible.

| Command | What it does |
| --- | --- |
| `river check [module]` | Checks the River module (default `river/`); exits 1 on any problem |
| `river hash <package-dir>` | Prints the `Approves:` line and `### Sums` section to record in `decisions.md` when approving a composed package |
| `river sum [module] [<item>...]` | Prints `requirements/sums.md` to stdout, or, given item names (`product.x`, `decisions.y`, `precedence.local`), prints the `### Sums` section a decision records for them |
| `river sum -w [module]` | Writes `requirements/sums.md` (River 2.0). Never edit that file by hand |
| `river install-skill [--origin <text>] [--pin <commit>] [repo]` | Installs or refreshes `.agents/skills/river/` and `.agents/skills/PROVENANCE.md` in `repo` (default `.`) |
| `river upgrade [--check] [--pin <commit>] [module]` | Migrates the River module and `.agents/skills/river/` to the latest River version, writing files without committing, and lists what needs a person; `--check` reports only |
| `river version` | Prints the River version the CLI implements and the versions it checks |

Checks of your own, beyond River's, go in your project's own language and test suite, which
cannot be skipped and add no toolchain.

**Before River's checks replace yours, list what yours enforced.** River's checker counts some
things a project's own checks may have failed (a requirement citing no goal, say), because a
brownfield project may start that way. If your project was stricter, keep the stricter rule as
your own check; otherwise swapping checkers is how a rule disappears without anyone deciding it
should.

## What a River module's checks cover

Run them before calling a change done. At minimum a River module should check that:

- every requirement's goal citation resolves, and every criterion's requirement citation does
- identifiers are well formed and unique within their package
- a pinned source still matches its recorded hash, by whatever means its style allows
- counts that are claimed are computed, not maintained by hand
- the instruction file and every line in the module stay within the limits in
  [`module.md`](module.md) (*Write every agent-read file to be read in part*)

**Every check must have been watched failing.** Plant a violation, see it reported, remove it.

## Versions

**A River module declares the version of River it follows** in the `## Module` block of `README.md`
at its root, as `river <major>.<minor>` beside `module`, the way `go.mod` declares `go 1.22`:

```
module example.com/app
river 2.0
```

The block holds `module` once, `river` at most once, blank lines and `//` comments; anything else
is `bad-module-block`. A module declaring none is read as 1.0 (`note [layout-undeclared]`); a
version the CLI does not know fails (`unsupported-layout`): upgrade the CLI.

**River's versions are semantic.** A module *passes* when `river check`, which also reads the skill copy's `PROVENANCE.md`,
reports no problems (notes are allowed). A minor version is compatible: a module that passed still passes, unedited, and an
agent following the new skill can carry on without restructuring it; upgrading is re-pinning and
declaring the new minor. A major version is not: a module that passed may fail, or the new skill
asks agents to restructure existing files. So a failure a major introduces is reported as a note
to modules declaring an earlier major, and `river upgrade` migrates a module across majors: it makes
the edits a tool can, writes them without committing, and lists what needs you. `river upgrade
--check` shows where a module stands and what upgrading would change, writing nothing.

| Version | Adds |
| --- | --- |
| 1.0 | River before versions. An approval records the package's hash (`Approves:`) |
| 2.0 | `requirements/sums.md`; `## Composed from`; `*(excluded)*` markers; `### Sums` on decisions and approvals; per-item `precedence.md` rows; `## Module` in root `README.md`; `river install-skill` and `river upgrade`. New failures (notes on 1.x): `unaccounted-requirement`, `covers-*`, `excluded-*`, `withdrawn-not-excluded`, `decision-missing-sums`, `exported-external-link`, `bad-module-block`, `precedence-*`, `included-and-composed`, `reaffirms-unknown`, `misplaced-module` |

**Record your own River version**: the skill's pin in its `PROVENANCE.md` and the layout here.
River keeps no list of adopters.

## Composition checks

Once a River module has composed packages, every source item, local or imported, must be accounted for
in each composed package it applies to: included, composed from, or excluded with `*(excluded)*`
([`composition.md`](composition.md), *Accounted for*). Fix a problem by recomposing, never by
editing the composed package.

| Code | Meaning | Fix |
| --- | --- | --- |
| `unaccounted-requirement` | A source item is not accounted for in a composed package it applies to, or in any composed package | Recompose: include, compose from, or exclude it |
| `not-composed` (note) | Sources exist and no composed package yet | Compose when the project is ready |
| `excluded-missing-marker` | (2.0) A line under `## Excluded` lacks `*(excluded)*` immediately after the item | Write `- [pkg.item](…) *(excluded)* — …` (`river upgrade` migrates 1.x) |
| `excluded-without-decision` | An excluded item line names no existing decision | Link the decision that excluded it, or record one |
| `withdrawn-not-excluded` | (2.0) A withdrawn source item is included unchanged or cited in `Composed from:` | Mark it `*(excluded)*` with the decision that withdrew it |
| `ambiguous-exclusion` (note) | A list line without `*(excluded)*` uses exclusion wording ("excluded", "not applicable", "withdrawn") | If excluded, add `*(excluded)*` and link the decision |
| `exported-external-link` | (2.0) A file in `composed/exported/<pkg>/` links outside `<pkg>/`, which breaks when imported | Write cross-package references in exported packages as plain names |
| `covers-empty` | A criterion's `Covers:` names nothing | Name the requirements it decides |
| `covers-unresolved` | A `Covers:` name resolves to nothing | Fix the name, or the criterion |
| `covers-withdrawn` | A `Covers:` name is a withdrawn item (`Status: withdrawn`) | Repoint it to what replaced the item |
| `decision-unknown-requirement` | An `Affects:` name is no item in its package, local, imported, or composed | Fix the name; if the item is gone, supersede the decision |
| `decision-missing-sums` | (2.0) An in-force decision affecting source items lacks `### Sums` for all of them | Add `### Sums` from `river sum <item>...` (`river upgrade` backfills 1.x) |
| `included-and-composed` | (2.0) A package accounts for one source item in more than one way (included, composed from, or excluded) | Recompose, accounting for each item once |
| `reaffirms-unknown` | (2.0) A decision's `Reaffirms:` names no entry in `decisions.md` | Name the decision it confirms |
| `composed-sources` | (2.0) A package has no `## Composed from` list, names a source that does not exist, or accounts for items from one it does not name | List its sources: `- product (Local)` |
| `stale-sums` | (2.0) `requirements/sums.md` is missing or differs from what `river sum -w` writes | Run `river sum -w` |
| `stale-approval` | The package's files changed since its latest approval, or (2.0) an input the approval's `### Sums` covered (a source item, a decision, a precedence row) changed, was withdrawn, or was removed, or a decision affecting a covered item was recorded since; the message names them | Recompose and ask the owner to approve the record `river hash` prints |
| `misplaced-module` | The `## Module` block is in `requirements/README.md` under River 2.0, or declared in both READMEs | Keep one, in `README.md` at the module root; `river upgrade` moves it |
| `skill-copy-missing` | (2.0) No copy of River's skill at `.agents/skills/river/` beside the module | Run `river install-skill` ([`adopting.md`](adopting.md), *Installing the skill*) |
| `skill-provenance-missing` | (2.0) No `.agents/skills/PROVENANCE.md` for the copy | Run `river install-skill` to write `PROVENANCE.md` with its pin |
| `skill-pin-missing` | (2.0) `PROVENANCE.md` records no pin, or only a placeholder | Run `river install-skill` (with `--pin <commit>` if built outside git) |
| `skill-version-mismatch` (note) | The skill copy describes a River version other than the module declares | Run `river upgrade` or `river install-skill` |
| `precedence-unknown-item` | (2.0) A row of `precedence.md`'s item table names no local or imported source item | Fix the name, or remove the row |
| `precedence-bad-tier` | (2.0) An item row's tier is not Mandatory, Governing, Local, Advisory, or a ranked tier the tier table names | Use one of them |
| `precedence-no-decision` | (2.0) An item row links no decision, or one that does not exist or is superseded | Link the decision in force that set the tier; ask the person if there is none |
| `precedence-tier-restated` | (2.0) A decision behind an item row has a `Tier:` field, or a field naming the row's tier | Remove it: the tier lives only in `precedence.md`; the decision says who and why |
| `precedence-duplicate-item` | (2.0) An item has two rows | Keep one: an item has one tier |
| `decision-review` (note) | A decision's `### Sums` show an item it affects changed, was withdrawn, or was removed since, and no later entry reaffirms it against current sums | Ask the person: reaffirm it in a new entry (`Reaffirms:`), or supersede it |

## Notes, not failures

**Some symptoms are notes, never failures:** a goal cited by one requirement, a goal whose
requirements all cite one other goal too, a goal longer than a paragraph, and a goal that contains
code or names another item. Each can be right. Each is where to look first.

What a check cannot do is tell you whether a requirement genuinely *serves* the goal it cites.
Service is a semantic relation, not a measurement. A green run is not a validated one.
