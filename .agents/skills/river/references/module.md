# The River module

A **River module** is the `river/` directory: requirements, the records kept about them, and,
optionally, goals as documentation. It is declared once in `river/README.md` with a `module` path
and the River version (`river 3.0`). It is not a Go module, though it borrows go.mod's form.

Part of River's method. Search `^## ` for the sections; read only the one you need.

## Layout

```
river/
  README.md                  what is here, and the module declaration
  bindings/                  package bindings: how THIS project works, and its review rule
  conformance/               package conformance: which guardrails implement which criteria
  goals/                     optional documentation; River checks only that links into it resolve
  requirements/
    README.md                what requirements/ holds
    imports/                 package imports: every imported source, with its pinned version
    precedence/              package precedence: which source outranks which
    decisions/               package decisions: how each conflict and judgment was settled
    sums.md                  a sum per input to composition, written by `river sum -w`, never by hand
    local/<name>/            a source written here: one source, one package
    imported/<import path>/  another module's source, pinned, with PROVENANCE.md
    composed/
      exported/<name>/       a composed package others may import, its criteria, and composed.sum
      internal/<name>/       a composed package private to this repository, and composed.sum
```

**Every authored package is a directory named after it**, built-in ones included: each file in it
declares the package on its first line (`package decisions`), an item can live in any file, and
identifiers are unique across the package. Split a long package into topic files a reader can read
whole; `river links -w` repoints links after an item moves. `sums.md` and each package's
`composed.sum` are written by the CLI.

A composed package holds one file per area, each with that area's requirements and their criteria
([`composition.md`](composition.md), *Files within a package*). Where a composed package sits is
the record of who may import it: moving it into `exported/` exports it.

**Open work is not in the module.** Track it where the project tracks work: an issue tracker, or a
`TODO.md` at the repository root. Link to requirements by qualified name if useful.

The method is not in the River module: it is this skill, kept beside the module at
`.agents/skills/river/` ([`adopting.md`](adopting.md)). The checks are not in the module either:
the `river` CLI runs them ([`checks.md`](checks.md)).

**Put the module in `river/` at the repository root**, unless you have a reason not to; if you
put it elsewhere, say where in `bindings/`. Your agent instruction file (`AGENTS.md` or equivalent)
says that the project uses River and where the module is, so an agent that has not loaded the skill
still learns to.

**Everything is plain Markdown.** No database, no server, no format that needs a tool to read.

**Write every agent-read file to be read in part.** An agent searches and opens ranges; it should
not have to read a file whole to find one item.

- **Keep the instruction file short**, under 6,000 characters: commands, the rules whose breach
  costs a commit, and a map. No history.
- **One item per heading, the identifier as the heading**, its title in bold beneath. Then
  `grep -n -A1 '^## '` is the table of contents and `## <id>` finds one item.
- **No line over 500 characters** outside a code block. Give a long item a heading instead.
- **Fixed field names** (`Status:`, `Owner:`) so a search can filter items without opening them.
  **Put a blank line between two fields**, so each renders on its own line; Markdown joins
  adjacent lines into one paragraph (`fields-run-together`).
- **Link every reference to its definition.** A name such as `product.noSilentFailures` is written
  as a Markdown link to the heading that defines it, so a reader follows it in one step and a moved
  item breaks a link the checker reports. Exported and imported packages, which move between
  modules, link only within their own directory.

## Keeping the River module's files

What each built-in package must contain. **Read this before writing to any of them.**

**`decisions/`**: the entry format and its rules are in [`decisions.md`](decisions.md).

**`imports/`**: one entry per imported source, `` ## `<name>` ``, then a field table: Name,
Origin, Version, Integrity, Tier, Import style (`vendored`, `fetched at build`, `queried`,
`quoted`, `absorbed`, `consulted`), Scope, Adopted parts, Owner, Status. The version is recorded
here and nowhere else. **Status** is `adopted`, `partial` (a named part not yet in force),
`rejected`, or `retired`.

**`imported/<import path>/`**: one directory per imported source, at its import path, the way a Go
import path names a package: `github.com/OWASP/asvs`, `w3.org/wcag`. It holds `<name>.md` (first
line `package <name>`, then the transcribed items) and `PROVENANCE.md` (each file's path, origin and
SHA-256, and how and when it was retrieved). Never fill in a source from memory.

**`precedence/`**: exactly one table of tiers, highest first, naming the sources in each, and at
most one item table (`| Item | Tier | Decision |`) for source items whose tier differs from their
package's ([`requirements.md`](requirements.md), *Precedence*). Tiers are recorded here and nowhere
else.

**`bindings/`**: what is true of this project only: where the module is, which stages it adopts,
local amendments to the method, and **the review rule**: how a change to the module is approved,
in exactly one `Review rule:` field across the package's files:

```markdown
Review rule: commits on main, merged through owner-reviewed pull requests
```

A module stating no rule is noted (`review-rule-missing`); two fail (`review-rule-duplicate`).

**Composed packages** open by saying what they are composed from and that they are changed by
recomposing, never by hand. A requirement accepted as unmet stays listed with the decision beside
it, as an excluded one is ([`requirements.md`](requirements.md), *Composition*).

## Composition stamps and review

**Source control approves.** A change to the module, a recomposition included, is approved the way
the project approves code, under the review rule in `bindings/`. River keeps no approvals of its
own.

**Each composed package carries a stamp, `composed.sum`**, written by `river compose <package>` at
the end of every composition, never by hand. It records a hash of the package's own files (its
`package` line) and the sum of every input the package was composed from: the source items it
includes, composes from, or excludes; the decisions it cites or whose `Affects:` names one of them;
and the rows of the precedence tables. When the package's files are edited by hand, or an input
moves afterward, `river check` fails the package with `stale-composition`, saying which; the fix is
to recompose and run `river compose` again. `river sum -w` never rewrites a stamp: `sums.md` is the
present, a stamp is the state at composition.

Every pull request that changes the module carries a requirements-impact section built from
`river diff --base` ([`pull-requests.md`](pull-requests.md)).
