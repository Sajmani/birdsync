# Adopting River in a repository

Part of River's method. Search `^## ` for the sections; read only the one you need.

## Principles

**You do not have to change any tool you already use.** River is files in your repository, this
skill, and the `river` CLI, run by whatever already runs your tests.

**Adopt one stage at a time.** Requirements alone is worth doing. A stage you have not adopted
is declared **informal** in `bindings.md` rather than left unmentioned, so a partial adoption is
distinguishable from an abandoned one.

One dependency is unavoidable: **Verification depends on Requirements**, because a criterion
with nothing behind it has nothing to check against. The others are free.

**Brownfield is the normal case.** An existing codebase already has requirements, implicit in
its code, tests and documentation. Recover them first and map them to goals afterwards. You will
start with every requirement unmapped and a test suite nobody wrote requirements for; neither is
a failure state, and both are counted.

**A project that already keeps requirements another way** (a spec directory, a requirements
document, a predecessor method) migrates them rather than starting over: every existing
requirement becomes a source item, included, excluded, or composed from, so none drops out
silently, and the old files are retired only after the migration accounts for each. A
requirement the predecessor withdrew becomes a source item with `Status: withdrawn — <why>`,
excluded with the decision that withdrew it ([`requirements.md`](requirements.md), *Status of a
source item*).

**Start a small project with one composed package, split into area files**, one per area of its
sources. Split into more packages only when a second approver, importer, or releasable artifact
appears, or when approving one part keeps re-approving another
([`composition.md`](composition.md), *Splitting into packages*).

## Installing the skill

**This skill is the method.** It is identical in every project that uses River, so it says
nothing about any particular project; what is true of yours goes in `bindings.md`.

Build or install the `river` CLI from a River checkout at the commit you are adopting (`go install
./cmd/river`, until River publishes a Go module path), then run `river install-skill` from your
repository root:

```
river install-skill [--origin <url>]
```

The command writes the embedded skill directory unchanged into `.agents/skills/river/`, the
project-skill location of the Agent Skills convention, and writes `.agents/skills/PROVENANCE.md`
beside it with the CLI's commit pin (`Pin`), origin, retrieval date, and local-amendments pointer:

```markdown
# PROVENANCE

River's skill, copied unchanged.

| | |
| --- | --- |
| Origin | <repository>, `skills/river/` |
| Pin | `<commit>` |
| Retrieved | <YYYY-MM-DD>, by `river install-skill` |
| Local amendments | none; any are listed in `river/bindings.md` |
```

Neither the user nor the agent creates `PROVENANCE.md` by hand. If an agent you use reads skills
only from its own directory (`.claude/skills/`, for one), point that location at
`.agents/skills/river/` rather than keeping a second copy. Keeping every skill file byte-identical
to its origin makes refreshing the skill a reviewable diff: re-run `river install-skill` (or run
`river upgrade`) with a newer `river` binary, which replaces `.agents/skills/river/`, removes any
renamed files, and updates `Pin` while preserving `Origin` and `Local amendments`. Put amendments
in `bindings.md` and say what they change.

**`river check` reads the copy.** From River 2.0 it fails a River module with no copy at
`.agents/skills/river/` in the repository holding it, no `.agents/skills/PROVENANCE.md`, or no pin
there (`skill-copy-missing`, `skill-provenance-missing`, `skill-pin-missing`; notes for 1.x). The
skill states the River version it describes in its `SKILL.md` frontmatter (`metadata:` then
`river: "2.0"`; a skill that states none describes 1.0). When that differs from the module's
declared version, `river check` notes it (`skill-version-mismatch`): run `river upgrade` or
`river install-skill`. `river upgrade --check` shows both versions. River's own repository holds
the skill's origin, `skills/river/`, and keeps no copy of it.

**The method is not a requirements source.** It is a pinned external document, but it
contributes no requirements to your composed set. Its pin belongs in `PROVENANCE.md` and nowhere
else: recording it in `imports.md` too gives two records of one pin, which a single recorded hash
exists to prevent.

## Procedure

1. **Ask which stages the project adopts.** Record each stage as adopted or **informal** in
   `river/bindings.md`.
2. **Install the skill**: run `river install-skill` from the repository root.
3. **Create the River module** at `river/` (or where the project says, recorded in `bindings.md`):
   `README.md` with a `## Module` block, `goals/`, `requirements/README.md`, `requirements/precedence.md`,
   `requirements/decisions.md`, and `requirements/local/<source>/`. The layout is in
   [`module.md`](module.md).
4. **Recover requirements first**, from the existing code, tests, and documentation, one item per
   `## <id>` heading with a bold title and a first paragraph that is the complete rule
   ([`requirements.md`](requirements.md), *Anatomy*). Leave `Goals:` empty where you cannot cite a
   real goal; `river check` counts these and never fails them. Never invent a goal to fill the
   field.
5. **Find the sources** the existing code already depends on: every hosted service, library, and
   dataset ([`discovery.md`](discovery.md)).
6. **Propose goals and their citations** for a person to accept. Cluster the unmapped
   requirements by who benefits and what they gain; each cluster is a candidate goal
   ([`goals.md`](goals.md)). Do not write accepted goals yourself.
7. **Map existing tests to criteria** from the criterion side, in `river/conformance.md`. Never
   annotate source code with requirement IDs; if the project already does, harvest every
   annotation into `conformance.md` first, then remove them. A test no criterion cites is kept and
   counted ([`verification.md`](verification.md)).
8. **Point the agent instructions at River**: add a line to the project's `AGENTS.md` (or
   equivalent) saying requirements live in `river/`, the method is the River skill, composed
   packages are never edited by hand, and `river check` must pass before a change is done.
9. **Wire `river check` in** wherever the project already runs checks: pre-commit, CI, and release.
10. Run `river check` and fix every problem.

## When the method is wrong

It will be. Write the amendment in `bindings.md`, say what it changes and why, and if the change
is general rather than local, send it upstream. A project copy that silently diverges from its
origin is worse than one that never refreshed, because the diff no longer means anything.
