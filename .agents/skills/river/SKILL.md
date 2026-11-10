---
name: river
description: >-
  Use this skill in any repository that has a `river/` directory (goals, requirements, and
  acceptance criteria as Markdown), or when the user asks to adopt River, add or change a
  requirement, compose requirements, approve a composed package, settle a conflict between
  requirements, or add a guardrail for an acceptance criterion. It is River's method: it tells the
  agent what to read before editing code, which changes need a person's approval, how to run
  River's checks, and where the full rules are.
metadata:
  river: "2.0"
---

# River

River keeps a project's goals, requirements, and acceptance criteria in `river/` beside the code,
and checks every change against them. **You draft; a person decides.** Whether a requirement is
right, which source wins a conflict, whether a gap may be accepted, and whether a composed version
is approved are never your calls.

## Before you edit code

1. Find the River module: `river/` at the repository root, or wherever the project's `AGENTS.md` (or
   `river/bindings.md`) says.
2. Read the composed packages that govern the code you will touch:
   `river/requirements/composed/{exported,internal}/<name>/`. Search for `## <id>` and read only
   the items you need; each item's first paragraph is its complete rule.
3. If the task would change **what the software must do**, stop and follow *Changing a
   requirement* below before writing code.
4. If the task would **connect to or depend on something the project did not build** — a hosted
   service or API, a dataset, a library, a model, an agent plugin or tool — that is not already in
   `river/requirements/imports.md`, find what governs it **before writing integration code**:
   follow [`references/discovery.md`](references/discovery.md). A person decides whether to
   accept its terms.

## Routine changes

A bug fix, refactor, or feature that every current requirement already allows, and that adds no
new external dependency, needs **no River edits and no approval**. Make the change, then run the
project's tests and `river check`. If a guardrail test fails, the code is wrong: fix the code,
never the test or the requirement.

## Changing a requirement

1. **Never edit `river/requirements/composed/`.** It is generated and approved by hash; a hand edit
   fails `river check` with `stale-approval`.
2. Draft the change in a source under `river/requirements/local/<source>/`. Cite the goals it
   serves (`Goals:`), as links. An imported source (`imported/`) belongs to another team: propose
   the change to its owner instead of editing it.
3. Recompose from composition's three inputs, the sources, `river/requirements/precedence.md`,
   and `river/requirements/decisions.md`, and write the result into the composed package, each
   item naming its sources in `Composed from:`. Follow
   [`references/composition.md`](references/composition.md).
4. If precedence does not settle a conflict, **stop and ask** the named owner. Record what they
   decide as a new entry in `river/requirements/decisions.md` (append-only; see the format in
   [`references/decisions.md`](references/decisions.md)).
5. Run `river sum -w` and `river check --allow-unapproved` so structural, accounting, and reference
   checks pass, then show `river diff` and ask the owner to approve the composed version. **Never
   run `river approve` unless the owner has explicitly approved the current recomposed state.** Once
   the owner approves (in conversation or via a code-review approval hook), `river approve --by
   <owner>` records the approval in `river/requirements/approvals.md`
   ([`references/module.md`](references/module.md), *Keeping the River module's files*), and strict
   `river check` passes.

## Adding a guardrail

1. Write the test, checker, or monitoring rule for one acceptance criterion. Do not annotate
   source code with requirement IDs.
2. **Watch it fail**: break the behavior it guards (or plant a bad input), run it, see it fail, and
   restore. A guardrail never seen failing does not count.
3. Record the mapping and the observation in the criterion's section of `river/conformance.md`,
   never in the criterion: `Watched failing: <date>, by <how>`. The format is in
   [`references/verification.md`](references/verification.md), *Recording conformance*.

## Checking

Run `river check` from the repository root (or `river check <module>` for a River module elsewhere). It exits 1 on any problem,
printing `[code] path: message`; notes (`note [code] ...`) are for a person's attention and never
fail. While drafting a recomposition or running pull-request CI before sign-off, run `river check
--allow-unapproved`, which reports `awaiting-approval` and `stale-approval` as notes; on `main` and
before release, run strict `river check` and fix every problem before calling a change done. Common
codes:

| Code | Meaning | Fix |
| --- | --- | --- |
| `awaiting-approval`, `stale-approval` | A composed package has no approval in `requirements/approvals.md`, or changed (in its files or covered inputs) since its approval; the message names what moved | Recompose if inputs moved, review `river diff`, and once the owner approves run `river approve --by <owner>` |
| `unaccounted-requirement` | A source item is not accounted for: no composed package includes it, composes from it, or excludes it | Recompose, accounting for it |
| `covers-unresolved`, `covers-withdrawn`, `covers-excluded` | A criterion covers an item that does not exist, is withdrawn, or is excluded in that package (`withdrawn-criterion` if the criterion itself is marked withdrawn in `composed/`) | Drop obsolete criteria on recomposition, or repoint `Covers:` to the in-force requirement |
| `guardrail-withdrawn`, `guardrail-unknown` | A `conformance.md` `Criterion:` entry or a `Verifies:` comment cites a withdrawn, excluded, or nonexistent criterion or requirement | Delete or repoint obsolete guardrails and clean up `conformance.md` |
| `stale-sums` | `requirements/sums.md` is out of date (River 2.0) | Run `river sum -w`; never edit it by hand |
| `decision-review` (note) | An item a decision affects changed since it was decided | Ask the person whether it still holds |
| `unknown-goal` | A requirement cites a goal that does not exist | Cite an existing goal, or ask whether a goal is missing |
| `unlinked-reference` | A `package.identifier` reference is not a link | Link it to its definition |
| `missing-package-decl` | A file defining items lacks `package <name>` on line 1 | Add it |
| `unrecorded-dependency` (note) | `go.mod`, `package.json`, or `requirements.txt` names a direct dependency `imports.md` does not | Follow [`references/discovery.md`](references/discovery.md); never just add the name to silence it |
| `unrecorded-host` (note) | Non-test code calls an `http(s)` host `imports.md` does not name | Search for its terms and follow [`references/discovery.md`](references/discovery.md); the project's own hosts go under *Considered and not adopted* |

Both discovery notes become problems once `imports.md` contains the line `Discovery: complete`.
Only a person adds that line, after working through every note.

The `## Module` block in `river/README.md` declares the River version the module follows (`river 2.0`);
`river version` lists the ones the CLI supports. `river sum -w` rewrites `requirements/sums.md` after
any edit to a source, `decisions.md`, or `precedence.md` (River 2.0); `river sum` prints what `-w` would write, or the `### Sums` block for named items. `river graph` (`--node <q>`, `--if-withdrawn <id>`, `--if-changed <id>`, `--format text|mermaid`) shows reachability from guardrails through criteria and requirements to goals before or after a change. `river upgrade` migrates a module and its skill copy to the latest version.

Install the CLI with `go install ./cmd/river` from a River checkout at the skill's pin (until River
publishes a Go module path), and install or refresh the skill copy with `river install-skill`. Every code and command is in [`references/checks.md`](references/checks.md).

## Adopting River in a repository

Follow [`references/adopting.md`](references/adopting.md): install this skill in the repository with `river install-skill`,
then start with requirements recovered from the existing code, tests, and documentation; leave
goals unmapped until a person maps them.

## The method

**This skill is River's method**, identical in every project that uses River; what is true of one
project goes in its `river/bindings.md`. Read only the reference you need, and only the section
(`^## `) within it:

| Reference | What it covers |
| --- | --- |
| [`overview.md`](references/overview.md) | The four stages and four processes, who decides, what is new |
| [`module.md`](references/module.md) | The `river/` layout, writing files to be read in part, what each requirements file holds |
| [`goals.md`](references/goals.md) | Goals versus requirements; gradients, targets, and backstops |
| [`requirements.md`](references/requirements.md) | Sources, precedence tiers, composition rules, conflicts, the anatomy of an item |
| [`composition.md`](references/composition.md) | Composition's three inputs, the procedure, obligations and horizontal requirements, collisions, packages and files |
| [`decisions.md`](references/decisions.md) | The `decisions.md` entry format and its append-only rules |
| [`verification.md`](references/verification.md) | Acceptance criteria, unverifiable requirements, guardrails |
| [`lifecycle.md`](references/lifecycle.md) | Implementation and designs, release and experimentation, optimization, revision |
| [`discovery.md`](references/discovery.md) | Finding what governs a new dependency |
| [`checks.md`](references/checks.md) | The `river` CLI and what a River module's checks must cover |
| [`adopting.md`](references/adopting.md) | Installing this skill, adopting stage by stage, migrating existing requirements |
