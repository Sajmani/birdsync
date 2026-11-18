---
name: river
description: >-
  Use this skill in any repository that has a `river/` directory (requirements and acceptance
  criteria as Markdown), or when the user asks to adopt River, add or change a requirement, compose
  requirements, settle a conflict between requirements, review a pull request that changes them,
  or add a guardrail for an acceptance criterion. It is River's method: it tells the agent what to
  read before editing code, which changes need a person's decision, how to run River's checks, and
  where the full rules are.
metadata:
  river: "3.0"
---

# River

River keeps a project's requirements and acceptance criteria in `river/` beside the code, and
checks every change against them. **You draft; a person decides.** Whether a requirement is right,
which source wins a conflict, whether a gap may be accepted, and whether a change to the
requirements merges are never your calls.

## Before you edit code

1. Find the River module: `river/` at the repository root, or wherever the project's `AGENTS.md` (or
   `river/bindings/`) says.
2. Read the composed packages that govern the code you will touch:
   `river/requirements/composed/{exported,internal}/<name>/`. Search for `## <id>` and read only
   the items you need; each item's first paragraph is its complete rule.
3. If the task would change **what the software must do**, stop and follow *Changing a
   requirement* below before writing code.
4. If the task would **connect to or depend on something the project did not build** — a hosted
   service or API, a dataset, a library, a model, an agent plugin or tool — that is not already in
   `river/requirements/imports/`, find what governs it **before writing integration code**:
   follow [`references/discovery.md`](references/discovery.md). A person decides whether to
   accept its terms.

## Routine changes

A bug fix, refactor, or feature that every current requirement already allows, and that adds no
new external dependency, needs **no River edits**. Make the change, then run the project's tests
and `river check`. If a guardrail test fails, the code is wrong: fix the code, never the test or
the requirement.

## Changing a requirement

1. **Never edit `river/requirements/composed/` by hand.** It changes only by recomposing; its stamp
   records the hash of its files, so a hand edit fails `river check` with `stale-composition`.
2. Draft the change in a source under `river/requirements/local/<source>/`. An imported source
   (`imported/`) belongs to another team: propose the change to its owner instead of editing it.
3. Recompose from composition's three inputs, the sources, `river/requirements/precedence/`, and
   `river/requirements/decisions/`, and write the result into the composed package, each item
   naming its sources in `Composed from:`. Follow
   [`references/composition.md`](references/composition.md).
4. If precedence does not settle a conflict, **stop and ask** the named owner. Record what they
   decide as a new entry in `river/requirements/decisions/` (append-only; see
   [`references/decisions.md`](references/decisions.md)).
5. Run `river sum -w`, then `river compose <package>`, which checks the package and writes its
   `composed.sum` stamp.
6. Open a pull request whose description has a **Requirements impact** section from
   `river diff --base main --format pr`, triaged into what needs the reviewer's call and what
   doesn't ([`references/pull-requests.md`](references/pull-requests.md)). The review approves the
   change, under the review rule in `river/bindings/`.

## Adding a guardrail

1. Write the test, checker, or monitoring rule for one acceptance criterion. Do not annotate
   source code with requirement IDs.
2. **Watch it fail**: break the behavior it guards (or plant a bad input), run it, see it fail, and
   restore. A guardrail never seen failing does not count.
3. Record the mapping and the observation in the criterion's section in `river/conformance/`,
   never in the criterion: `Watched failing: <date>, by <how>`. The format is in
   [`references/verification.md`](references/verification.md), *Recording conformance*.

## Checking

Run `river check` from the repository root (or `river check <module>` for a River module
elsewhere). It exits 1 on any problem, printing `[code] path: message`; notes
(`note [code] ...`) are for a person's attention and never fail. Fix every problem before a change
is done.

| Code | Meaning | Fix |
| --- | --- | --- |
| `stale-composition`, `composition-unstamped` | A composed package's files or inputs changed since its stamp, or it has none | Recompose, then `river compose <package>` |
| `unaccounted-requirement` | A source item is not included, composed from, or excluded in a package it applies to | Recompose, accounting for it |
| `stale-sums` | `requirements/sums.md` is out of date | Run `river sum -w`; never edit it by hand |
| `covers-*`, `guardrail-*` | A criterion or guardrail cites something withdrawn, excluded, or missing | Repoint it, or drop it on recomposition |
| `unlinked-reference`, `broken-reference-link` | A qualified name is not a link, or its link leads nowhere | Link it; `river links -w` repoints moved items |
| `version-too-old` | The module declares River 1.x or 2.x | Run `river upgrade`, or use River built at the module's pin |
| `unrecorded-dependency`, `unrecorded-host` (notes) | A dependency or host is not in `imports/` | Follow [`references/discovery.md`](references/discovery.md); never just add the name |

The `## Module` block in `river/README.md` declares the River version the module follows
(`river 3.0`). `river graph` shows what reaches what, from guardrails to sources, before or after a
change; `river upgrade` migrates an older module and its skill copy. Install the CLI with
`go install ./cmd/river` from a River checkout at the skill's pin, and refresh the skill copy with
`river install-skill`. Every code and command is in [`references/checks.md`](references/checks.md).

## The method

**This skill is River's method**, identical in every project that uses River; what is true of one
project goes in its `river/bindings/`. To adopt River, follow
[`references/adopting.md`](references/adopting.md). Read only the reference you need, and only the
section (`^## `) within it:

| Reference | What it covers |
| --- | --- |
| [`overview.md`](references/overview.md) | River's core phases, who decides, what is new |
| [`module.md`](references/module.md) | The `river/` layout, built-in packages, composition stamps and review |
| [`requirements.md`](references/requirements.md) | Sources, precedence tiers, composition rules, conflicts, the anatomy of an item |
| [`composition.md`](references/composition.md) | Composition's three inputs, the procedure, horizontal requirements, collisions, packages and files |
| [`decisions.md`](references/decisions.md) | The decision entry format and its append-only rules |
| [`pull-requests.md`](references/pull-requests.md) | The requirements-impact section of a pull request |
| [`verification.md`](references/verification.md) | Acceptance criteria, unverifiable requirements, guardrails |
| [`lifecycle.md`](references/lifecycle.md) | Implementation and designs, revision, future phases |
| [`goals.md`](references/goals.md) | Goals, dormant in River 3.0 |
| [`discovery.md`](references/discovery.md) | Finding what governs a new dependency |
| [`checks.md`](references/checks.md) | The `river` CLI and what a River module's checks must cover |
| [`adopting.md`](references/adopting.md) | Installing this skill, adopting stage by stage, migrating existing requirements |

Phases River does not yet define (goals, gradients, release guardrails, experimentation) are
described in `ROADMAP.md` at the root of River's own repository.
