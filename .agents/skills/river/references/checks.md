# Checks

Part of River's method. Search `^## ` for the sections; read only the one you need.

## The river CLI

**River's checks are one Go binary, `river`, the same for every project.** Run `river check`
(or `river check <module>`) wherever the project already runs checks: the agent's edit loop,
pre-commit, CI, review, and release. The shared cases in River's `river/check/testdata/` are its
contract. Record the CLI version you run in `bindings/`, beside the skill's pin.

| Command | What it does |
| --- | --- |
| `river check [module]` | Checks the River module (default `river/`), reporting every finding and exiting 1 on any problem |
| `river compose [--dry-run] [module] [<package>...]` | Writes each named composed package's `composed.sum` (all, with none named), refusing while any other problem remains |
| `river diff [--base <ref>] [--format text\|pr] [module]` | Reports what moved since a git ref (default: the merge base with `main`), by area ([`pull-requests.md`](pull-requests.md)) |
| `river graph [--format text\|mermaid] [--node <q> \| --if-withdrawn <id> \| --if-changed <id>] [module]` | Shows the graph from guardrails through criteria and requirements to imported sources, one node's neighborhood, or the cascade of withdrawing or changing a node |
| `river links [-w] [module] [<extra path>...]` | Lists, or with `-w` rewrites, links whose target file no longer defines the item they name; exits 1 on a dry run that finds any |
| `river sum [-w] [module] [<item>...]` | Prints `requirements/sums.md`, or the sum lines of named inputs; `-w` writes `sums.md`. Never edit it by hand |
| `river install-skill [--origin <text>] [--pin <commit>] [repo]` | Installs or refreshes `.agents/skills/river/` and `.agents/skills/PROVENANCE.md` |
| `river upgrade [--check] [--pin <commit>] [module]` | Migrates a module from any earlier River version to the latest, writing files without committing, and lists what needs a person; `--check` reports only |
| `river version` | Prints the River version the CLI implements and the versions it checks |

Checks of your own, beyond River's, go in your project's own language and test suite.

## What a River module's checks cover

At minimum, a River module should check that:

- every reference resolves: links, `Covers:`, `Affects:`, `Composed from:`, and `Goals:`
- identifiers are well formed and unique within their package
- every source item is accounted for in each composed package it applies to
- each composed package's stamp records its inputs as they are now
- a pinned source still matches its recorded hash, by whatever means its style allows
- the instruction file and every line in the module stay within the limits in
  [`module.md`](module.md)

**Every check must have been watched failing.** Plant a violation, see it reported, remove it.

## Versions

**A River module declares the version of River it follows** in the `## Module` block of `README.md`
at its root, the way `go.mod` declares `go 1.22`:

```
module example.com/app
river 3.0
```

The block holds `module` once, `river` at most once, blank lines and `//` comments; anything else
is `bad-module-block`.

**A River 3.x CLI checks only 3.x modules.** A module declaring 1.x or 2.x, or no version, fails
with `version-too-old`: run `river upgrade`, or check it with River built at the module's pin (see
`.agents/skills/PROVENANCE.md`). A version newer than the CLI knows fails with `unsupported-layout`.

**River's versions are semantic.** A minor version is compatible: a module that passed still
passes, unedited. A major version is not, and `river upgrade` migrates a module across majors: it
makes the edits a tool can, writes them without committing, and lists what needs you.

| Version | What it is |
| --- | --- |
| 1.0 | River before versions; approvals in `decisions.md` |
| 2.0 | `sums.md` and `approvals.md`; `## Composed from`; `*(excluded)*`; per-item precedence rows; strict failures |
| 2.1 | Goal graph, imported goals, and goal versions |
| 3.0 | Core only: source control approves and each composed package carries `composed.sum`; decisions record no sums; goals dormant and gradients set aside; built-in packages as directories; open work out of the module; `river compose`, `river diff --base`, `river links` |

## Codes

Fix a composition problem by recomposing, never by editing the composed package.

| Code | Meaning | Fix |
| --- | --- | --- |
| `unaccounted-requirement` | A source item is not accounted for in a composed package it applies to | Recompose: include, compose from, or exclude it |
| `not-composed` (note) | Sources exist and no composed package yet | Compose when the project is ready |
| `composed-sources` | A package has no `## Composed from` list, or it misses or misnames a source | List its sources: `- product (Local)` |
| `included-and-composed` | A package accounts for one source item in more than one way | Recompose, accounting for each item once |
| `excluded-missing-marker` | A line under `## Excluded` lacks `*(excluded)*` after the item | Write `- [pkg.item](…) *(excluded)* — …` |
| `excluded-without-decision` | An excluded item names no existing decision | Link the decision that excluded it, or record one |
| `withdrawn-not-excluded` | A withdrawn source item is included or composed from | Mark it `*(excluded)*` with the decision that withdrew it |
| `ambiguous-exclusion` (note) | An included line uses exclusion wording | If excluded, add `*(excluded)*` and link the decision |
| `composition-unstamped` | A composed package has no `composed.sum` | Finish composing, then run `river compose <package>` |
| `stale-composition` | The package's files changed since its stamp, an input the stamp records moved, was withdrawn, or was removed, or the package now covers one it doesn't record | Recompose, then run `river compose <package>` |
| `bad-composition-sum` | `composed.sum` is malformed | Run `river compose <package>`; never edit it by hand |
| `stale-sums` | `requirements/sums.md` is missing or out of date | Run `river sum -w` |
| `covers-empty`, `covers-unresolved`, `covers-withdrawn`, `covers-excluded` | A criterion's `Covers:` names nothing, nothing that exists, a withdrawn item, or one the package excludes | Repoint `Covers:`, or drop the criterion on recomposition |
| `withdrawn-criterion` | A criterion or composed item inside `composed/` is marked withdrawn | Recompose without it |
| `guardrail-withdrawn`, `guardrail-unknown` | A conformance entry or a `Verifies:` comment cites something withdrawn, excluded, or nonexistent | Delete or repoint the guardrail |
| `decision-unknown-requirement`, `decision-unknown-ref`, `reaffirms-unknown` | A decision names an item or entry that doesn't exist | Fix the name, or supersede the decision |
| `decision-missing-field`, `decision-bad-status`, `decision-status-mismatch` | A decision lacks a field, or its `Status:` and the entry it names disagree | Fix the entry's format |
| `decision-sums-ignored` (note) | A decision still carries a 2.x `### Sums` block | Leave it; `river upgrade` removes it |
| `precedence-*` | An item row names no source item, a bad tier, or no decision in force; restates its tier in the decision; repeats an item; or the package holds more than one tier table | Fix the row or table |
| `unknown-goal` | A `Goals:` or `Serves:` names a goal that doesn't exist | Fix the name |
| `gradients-set-aside` (note) | The module keeps a `gradients/` directory | Nothing: it is documentation |
| `unlinked-reference`, `unknown-reference`, `broken-reference-link`, `exported-external-link` | A reference is not a link, names nothing, leads nowhere, or leaves an exported package | Link it to its definition; in exported packages write outside names plainly |
| `fields-run-together` | Two field lines sit on adjacent lines outside a code block, so they render as one line | Put a blank line between them |
| `missing-package-decl`, `unexpected-package`, `package-*`, `duplicate-*`, `mixed-packages` | A file's package line is missing or wrong, or a name is declared twice | Fix the declaration |
| `builtin-file-moved` | A 2.x one-file built-in package (such as `requirements/decisions.md`) is left in the module | Run `river upgrade` |
| `approvals-file-removed` | `requirements/approvals.md` is left in the module | Run `river upgrade` |
| `review-rule-missing` (note), `review-rule-duplicate` | Package bindings states no `Review rule:`, or more than one | State the one rule: how code review approves a change to the module |
| `open-work-in-module` (note) | `todo.md`, `todo-closed.md`, or `proposals.md` is in the module | Move it where the project tracks work |
| `misplaced-module` | The `## Module` block is in `requirements/README.md`, or in both READMEs | Keep one, in `README.md` at the module root |
| `version-too-old`, `unsupported-layout` | The module declares a version this CLI doesn't check | Run `river upgrade`, or use the matching CLI |
| `skill-copy-missing`, `skill-provenance-missing`, `skill-pin-missing` | The skill copy, its `PROVENANCE.md`, or its pin is missing | Run `river install-skill` ([`adopting.md`](adopting.md)) |
| `skill-version-mismatch` (note) | The skill copy describes another River version | Run `river upgrade` or `river install-skill` |
| `unrecorded-dependency`, `unrecorded-host` (notes) | A dependency or called host is not in `imports/` | Follow [`discovery.md`](discovery.md); they become problems once `imports/` says `Discovery: complete` |

## In CI

Run `river check` on every pull request and on `main`. Approval is the project's own code review,
under the review rule in `bindings/`; add `river diff --base main --format pr` to the pull request
description ([`pull-requests.md`](pull-requests.md)), or to the CI step summary as a backstop.
`--allow-unapproved` is accepted and ignored, for one major version, so existing CI keeps working.

A green run is not a validated one: what a check cannot tell you is whether a requirement is the
right one.
