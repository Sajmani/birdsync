# The River module

A **River module** is the `river/` directory: goals, gradients, requirements, and the records
kept about them, declared once in `river/README.md` with a `module` path and the River version.
It is not a Go module, though it borrows go.mod's form.

Part of River's method. Search `^## ` for the sections; read only the one you need.

## Layout

```
river/
  README.md                  what is here, and the module declaration (River 2.0)
  bindings.md                how THIS project works, and which stages it adopts
  conformance.md             what is built, against which version, and which criteria are met
  proposals.md               drafts awaiting a person, and open questions
  goals/                     package goals
  gradients/                 package gradients
  requirements/
    README.md                what requirements/ holds
    imports.md               every imported source, with its pinned version
    precedence.md            standing rules: which source outranks which
    decisions.md             log of how each conflict and judgment was settled
    approvals.md             approved version and covered sums per composed package, written by `river approve` (2.0)
    sums.md                  a sum per input to composition, written by `river sum -w`, never by hand (2.0)
    local/<name>/            a source written here: one source, one package
    imported/<import path>/  another module's source, pinned, with PROVENANCE.md
    composed/
      exported/<name>/       a composed package others may import, and its criteria
      internal/<name>/       a composed package private to this repository
```

A composed package holds a `README.md` (what it is composed from, the precedence applied, and
what is excluded) and one file per area, each with that area's requirements and their criteria
([`composition.md`](composition.md), *Files within a package*). `conformance.md`'s format is in
[`verification.md`](verification.md), *Recording conformance*.

Each composed file says it is composed and from what, and is changed by recomposing, never by
hand. Where a composed package sits is the record of who may import it: moving it into
`exported/` exports it, and no other file lists exports. These names were chosen by testing
candidate layouts on zero-context agents acting as new engineers.

The method is not in the River module: it is this skill, kept beside the module at
`.agents/skills/river/` ([`adopting.md`](adopting.md)). The checks are not in the module either:
the `river` CLI runs them ([`checks.md`](checks.md)).

**Put the module in `river/` at the repository root**, unless you have a reason not to; if you
put it elsewhere, say where in `bindings.md`. Paths inside it are relative to its root, and your
agent instruction file (`AGENTS.md` or equivalent) says that the project uses River and where the
module is, so an agent that has not loaded the skill still learns to.

**Everything is plain Markdown.** No database, no server, no format that needs a tool to read.

**Write every agent-read file to be read in part** (`product.seekableArtifacts`). An agent
searches and opens ranges; it should not have to read a file whole to find one item.

- **Keep the instruction file short** — under 6,000 characters. It is loaded in every session,
  so it holds commands, the rules whose breach costs a commit, and a map. No history.
- **One item per heading, the identifier as the heading**, its title in bold beneath. Then
  `grep -n -A1 '^## '` is the table of contents and `## <id>` finds one item. Keep no
  hand-maintained index.
- **No line over 500 characters** outside a code block. A paragraph in a table row is one search
  hit that cannot be read by range; give the item a heading instead.
- **Fixed field names** (`Status:`, `Owner:`) so a search can filter items without opening them.
- **Separate by lifecycle.** Closed items and background notes go in their own files, out of the
  path of the open ones.
- **Say how to search a long file** in one line at its top, and link rather than copy between
  files.
- **Link every reference to its definition.** A name such as `product.noSilentFailures` is written
  as a Markdown link to the heading that defines it, in prose and in `Goals:` fields alike, so a
  reader follows it in one step and a renamed item breaks a link the checker reports. Exported and
  imported packages, which move between modules, link only within their own directory.

## Keeping the River module's files

What each file under `requirements/` must contain, so that anyone adding to it keeps it
consistent. **Read this before writing to any of them.** An existing entry is a good example of
the format, but it cannot show the rules below that only apply when something changes.

**`decisions.md`** — the format and its rules are in [`decisions.md`](decisions.md).

**`imports.md`** — one entry per imported source, `` ## `<name>` ``, then a field table: Name,
Origin, Version, Integrity, Tier, Import style (`vendored`, `fetched at build`, `queried`,
`quoted`, `absorbed`, `consulted`), Scope, Adopted parts, Owner, Status. The version is recorded
here and nowhere else. **Status** is `adopted` when every part named in Adopted parts is in force,
`partial` when some named part is not yet (not transcribed, not pinned, not checked), `rejected`
when evaluated and declined, and `retired` when adopted and later removed; adopting one criterion
of a large standard is a narrow Adopted parts, not `partial`.

**`imported/<import path>/`** — one directory per imported source, at its import path: the host
that publishes it, then the path naming the publisher's project, ending in the package name, the
way a Go import path does: `github.com/OWASP/asvs`, `w3.org/wcag`. A source known by a DOI uses
the publisher behind it (`figshare.com/avonet`), not the resolver. The directory holds
`<name>.md` (first line `package <name>`, then the transcribed items, each `## <Identifier>` with
a bold title) and `PROVENANCE.md` (each file's path, origin and SHA-256, and how and when it was
retrieved). Mark anything not yet retrieved as such; never fill it in from memory. The source's
tier goes in `precedence.md`.

**`precedence.md`** — one table of tiers, highest first, naming the sources in each. A source is
added to a tier when it is imported. Below it, optionally, an item table (`| Item | Tier |
Decision |`): one row per source item whose tier differs from its package's, linking the item and
the decision behind it ([`requirements.md`](requirements.md), *Precedence*). Tiers are recorded
here and nowhere else: a decision says who and why, not the tier.

**Composed packages** — each opens by saying what it is composed from and that it is changed by
recomposing. **Without a composition tool, recompose by hand:** rewrite the file from its sources,
`precedence.md` and `decisions.md`, in a commit that says it recomposed and why. A requirement
accepted as unmet stays listed, with the decision beside it, as an excluded one is
([`requirements.md`](requirements.md), *Composition*):
`— <title>. Accepted unmet by decisions.<id>`, the decision linked.

**`approvals.md` — a person approves each composed version with its hash** (`tech.approvalsCarryHashes`).
From River 2.0, approvals live in `requirements/approvals.md` (`package approvals`), written by
`river approve --by <owner>` and never by hand, with one `## <package>` section per composed package
updated in place rather than appended to `decisions.md` (git history preserves prior approvals, and
`approvals.md` is not an input in `sums.md`). Each section carries `Date:`, `Approved by:`,
`Approves: <package>, sha256:<hash>` (the SHA-256 of the package's `.md` files, each contributing
its path relative to the package directory, a newline, and its bytes, in path order), and a
`### Sums` block holding the `sums.md` lines for **every input the package was composed from** — the
source items it includes, composes from, or excludes, every decision it cites or whose `Affects:`
names one of them, and every row of `precedence.md`'s tables:

    ## app

    **Approved version of composed/internal/app**

    Date: 2026-11-10
    Approved by: Sam
    Approves: app, sha256:9be64dc7…

    ### Sums

    ```
    decisions.dropGamma sha256:ecca6c31…
    precedence.governing sha256:76feefa0…
    product.alphaRule sha256:db2ea2c9…
    terms.RateLimit sha256:1a7a8758…
    ```

The commit that updates the entry identifies the approved version. A composed package is never
edited by hand; a change is a recomposition and a new approval. Under River 2.0, strict `river
check` fails both a never-approved package (`awaiting-approval`) and one whose files or covered
inputs changed since approval (`stale-approval`, naming each); `river check --allow-unapproved`
reports both as notes so draft branches and pull requests can check everything else before the owner
signs off ([`checks.md`](checks.md), *Draft branches, PR CI, and review hooks*). `river diff`
summarizes what changed since the recorded approval. Repositories on GitHub may mark
`river/requirements/sums.md` and `river/requirements/approvals.md` `linguist-generated=true` in
`.gitattributes` so pull requests collapse machine-written hash blocks by default. `PROVENANCE.md`
hashes stay out of the approval: they vouch for a transcription, not for what was approved. A River
1.x module records `Approves:` in `decisions.md`, and `river upgrade` moves those entries into
`approvals.md`.

**An agent that needs a new or changed requirement** adds it to a source under `local/`,
recomposes, runs `river sum -w` and `river check --allow-unapproved`, shows `river diff`, and waits
for the package's owner to approve before running `river approve --by <owner>`. It never edits a
composed file by hand or runs `river approve` without the owner's approval of that state.

## Proposals

**Drafts awaiting a person, and open questions, go in `river/proposals.md`**, and nowhere else.
A file in `goals/` defines accepted goals, an entry in `decisions.md` claims a decision was made,
and a composed package is approved by hash, so a draft in any of them says something false or
changes what was approved. Proposed goals, decisions, imports, withdrawals, and stage choices wait
here, each as `## <id>` with a bold title, what is proposed, and what the person must decide. The
file opens with `package proposals`.

**An item leaves when it is decided:** recorded where it belongs (a goal file, `decisions.md`,
`imports.md`, `bindings.md`), or dropped, with the outcome, if the person declines it. Never put a
question or a draft inside a composed package; the CLI counts open proposals as a note.
