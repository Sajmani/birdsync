# The bundle

Part of River's method. Search `^## ` for the sections; read only the one you need.

## Layout

```
river/
  bindings.md                how THIS project works, and which stages it adopts
  conformance.md             what is built, against which version, and which criteria are met
  goals/                     package goals
  gradients/                 package gradients
  requirements/
    README.md                what is here, and the module declaration
    imports.md               every imported source, with its pinned version
    precedence.md            standing rules: which source outranks which
    decisions.md             log of how each conflict and judgment was settled
    local/<name>/            a source written here: one source, one package
    imported/<import path>/  another module's source, pinned, with PROVENANCE.md
    composed/
      exported/<name>/       a composed package others may import, and its criteria
      internal/<name>/       a composed package private to this repository
```

Each composed file says it is composed and from what, and is changed by recomposing, never by
hand. Where a composed package sits is the record of who may import it: moving it into
`exported/` exports it, and no other file lists exports. These names were chosen by testing
candidate layouts on zero-context agents acting as new engineers.

The method is not in the bundle: it is this skill, kept beside the bundle at
`.agents/skills/river/` ([`adopting.md`](adopting.md)). The checks are not in the bundle either:
the `river` CLI runs them ([`checks.md`](checks.md)).

**Put the bundle in `river/` at the repository root**, unless you have a reason not to; if you
put it elsewhere, say where in `bindings.md`. Paths inside it are relative to its root, and your
agent instruction file (`AGENTS.md` or equivalent) says that the project uses River and where the
bundle is, so an agent that has not loaded the skill still learns to.

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

## Keeping the bundle's files

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
added to a tier when it is imported.

**Composed packages** — each opens by saying what it is composed from and that it is changed by
recomposing. **Without a composition tool, recompose by hand:** rewrite the file from its sources,
`precedence.md` and `decisions.md`, in a commit that says it recomposed and why. A requirement
accepted as unmet stays listed, with the decision beside it, as an excluded one is
([`requirements.md`](requirements.md), *Composition*):
`— <title>. Accepted unmet by decisions.<id>`, the decision linked.

**A person approves each composed version, in `decisions.md`, with its hash** (`tech.approvalsCarryHashes`).
The entry carries `Approves: <package>, sha256:<hash>`: the SHA-256 of the package's `.md` files,
each contributing its path relative to the package directory, a newline and its bytes, in path
order. The commit that adds the entry identifies the published version. A published package is
never edited; a change is a recomposition and a new approval. The checker fails a package whose
latest approval names a different hash (`stale-approval`), so a hand edit, by a person or an agent,
is caught at the next run; a package never approved is reported as awaiting approval, not failed.

**An agent that needs a new or changed requirement** adds it to a source under `local/`,
recomposes, and asks the package's owner to approve the result. It never edits a composed file.
