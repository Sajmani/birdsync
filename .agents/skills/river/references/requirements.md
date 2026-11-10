# Requirements

Part of River's method. Search `^## ` for the sections; read only the one you need.

## Sources

A **source** provides requirements. A **local** source is authored in your repository. An
**imported** source is written by someone else — a compliance standard, a terms of service, a
contract — and is **pinned and integrity-checked**, so a change to it produces a reviewable
diff. Record the hash in exactly one place.

**Vendoring is one way to pin a source and is often not available.** A license may forbid
redistribution, terms may forbid automated access, or the source may be an API with nothing to
copy. Record an **import style** per source — vendored, fetched at build, queried, quoted,
absorbed, consulted — and pin it by whatever the style allows. What must hold is that a change
is detectable, whatever the mechanism.

**A source you cannot retrieve is a recorded gap, never a reconstruction.** Terms pages refuse
automated clients; publishers go offline. Record what could not be read and leave it unread.
Never supply its contents from memory: a plausible rate limit or license term is
indistinguishable from a real one, and an agent will produce one readily.

**A source has a lifecycle, and refusal is part of it.** A source you evaluated and
**rejected** is kept, with the reason: the refusal is a durable decision, and without it the
next person evaluates the source again. A source you adopted and later **retired** is kept too,
with what depended on it. Neither is the same as a source nobody considered.

**Look for sources before integrating anything new.** Before a change adds a hosted service, an
API, a dataset, a library, a model, or an agent plugin, find what governs its use: terms of
service, license, usage limits, attribution, and data-use terms. Record each in `imports.md`,
carry its obligations into a source, and let a person decide whether to accept them before the
integration lands. A source nobody considered is the one an auditor finds first.

**One source is one package.** A source may span several files in its directory without becoming
several packages; files organize, packages namespace. If a grouping deserves a prefix it
deserves to be a package, and if it does not, it is a file.

## Precedence

`precedence.md` puts each source in a tier, highest first. Composition uses the tiers to settle
a conflict on its own; anything they can't settle goes to a person. There are four:

| Tier | When it conflicts with another tier | Who can grant an exception | Typical sources |
| --- | --- | --- | --- |
| **Mandatory** | Always wins | Nobody locally | Law, regulation, accessibility standards, company security policy |
| **Governing** | Beats your own requirements | The source's owner, when you ask; record it in `decisions.md` | Brand guidelines, company engineering standards, app store rules |
| **Local** | Your own requirements | You | Everything you decide for yourself |
| **Advisory** | Loses to your own; applies where they say nothing | You, with a reason on record | Best practices, style guides, internal cookbooks |

The tier belongs to how you adopt a source, not to the source: the same brand guide can be
Governing for a shipping product and Advisory for a prototype. Within a tier you may rank
sources further, as in *Local: tech* above *Local: product*.

**One item may take a different tier from its package**, as when a Governing source only
recommends something. Record it in `precedence.md`'s item table, one row per item, with a link to
the decision that is the authority for it:

```markdown
| Item | Tier | Decision |
| --- | --- | --- |
| [inatapi.identifyUserAgent](imported/inaturalist.org/inatapi/inatapi.md#identifyuseragent) | Advisory | [decisions.userAgentRecommended](decisions.md#useragentrecommended) |
```

An item's tier is its row's if it has one, and its package's otherwise. The tier is written only
here; the decision says who decided and why, and doesn't repeat it. Nothing else records an
override: a local requirement that beats an Advisory item is ordinary precedence. From River 2.0,
`river check` fails a row naming no source item, a tier that isn't one of the four (or a ranked
tier the tier table names), a row without a live decision, two rows for one item, and a decision
whose fields restate the tier.

## Composition

**Composition** applies precedence and recorded decisions to your sources and produces one or
more **composed requirements packages** — explicit, versioned artifacts rather than an implied
state. Its three inputs, sources, precedence, and decisions, count equally; a fact in none of them
is asked for and recorded before composing on it. Each package names the sources it came from;
`imports.md` holds their versions, and the CLI records their sums in the River module's `sums.md`. The
procedure is in [`composition.md`](composition.md).

A composed package is named for its **audience or scope**, the way a Go package is, never for a
version: a version is a commit, changelist or revision, so a directory called `v1/` is always
wrong. Once published, a composed version is never edited; a later composition is a new version,
and simultaneous milestones live on separate lines of development. Identify a version by what
your version-control system never reassigns (a commit hash, changelist or revision number) and
check it by its content hash. A tag or label is a readable alias, worth protecting where your
system allows, but no system makes one immutable by default, so never trust it alone.

Downstream repositories import your **exported composed packages**, never your raw sources, so
every cross-repository obligation has been through composition with explicit precedence and
human sign-off.

**A composed package names everything that applies, and writes out only what composition made.**
Requirements that apply unchanged are a list grouped by area, one line each, linking the source
item and giving its title:

```markdown
### Diagram

- [product.scaleToggle](../../../local/product/diagram.md#scaletoggle) — Diagrams scale per section or globally, never per tile

## Excluded

- [tech.localNowPublicLater](../../../local/tech/posture.md#localnowpubliclater) *(excluded)* — Local now, public later ([decisions.publicDeploymentDeferred](../../../decisions.md#publicdeploymentdeferred))
```

What composition wrote — a project rule from an obligation, a rule over a domain, an
intermediary's contract, a requirement a higher tier changed — follows in full, each item naming
its sources in `Composed from:`, every entry a link to the item it names (plain names in an
exported package, which links only within its own directory). **Every source requirement is
included, excluded (marked `*(excluded)*` after its link), or composed from, in every composed
package it applies to** ([`composition.md`](composition.md), *Accounted for*); one that is none of
these has silently dropped out. Included requirements keep their source names. Do not commit a
flattened copy with everything inlined beside the sources: every search would find each
requirement twice. If an importer or auditor needs one, generate it when you tag.

**Include an item unchanged only if it already says what this project must do.** An obligation
phrased for an operator or a user ("use OAuth", "don't machine-generate content") is composed into
a rule for this project ([`composition.md`](composition.md), *Obligations become rules*).

**A requirement over every instance of something is composed once**, over its domain, never per
instance. A source names the domain a requirement applies to in an `Applies to:` field, in its own
term — `Applies to: wcag.WebPage`; where it doesn't, an interpretation decision beside the import
names it, and the imported text stays as it is. You define what counts as an instance in a local
requirement whose `Binds:` field names that domain, and decide the edge cases there: is an error
page a web page? The composed package states the rule once, *every web page meets
`wcag.ContrastEnhanced`*, and records any excluded instance as an accepted gap. Its criteria test
the domain exhaustively over what varies and by sampling the instances. The local definition
makes the source's term concrete and **never narrows it**: if the source's own definition makes
something an instance, leaving it out is an exception, recorded instance by instance. **Nothing in
the module lists the instances.** Verification discovers them from the system, so a new page
inherits its obligations by existing and cannot escape by being left off a list. Discover from
what the system declares — its routes, its templates — as well as by following links: an instance
nothing links to is still an instance. The procedure is [`composition.md`](composition.md),
*Projecting horizontal requirements*.

**An intermediary carries a contract.** A stylesheet, a framework or a middleware that makes its
users meet a requirement states the terms: given an assumption about how it is used, it
guarantees the requirement, over a stated scope, and says whether its users may override it. The
stylesheet guarantees 7:1 text, given that a page takes its colors only from it; the framework
guarantees scrubbed request logs, given that nothing logs around it. A contract from another
repository is imported at a pin like any other requirement. Layered intermediaries each carry
their own, and the order they run in is part of the contract, checked where they are assembled.
Where an intermediary carries a project's obligation, the project's composed rule cites the
contract, and the project's own requirement is the contract's assumption.

## Status of a source item

**A source item is in force unless its `Status:` line says otherwise.** The one other status is
`Status: withdrawn — <why>`, with the decision that withdrew it linked in the reason:

```markdown
## noRollback

**Releases never roll back schema changes**

Status: withdrawn — replaced by a composed rule; see [decisions.rollbackComposed](../../decisions.md#rollbackcomposed)

Releases never roll back a schema change; a bad change is fixed forward.
```

A withdrawn item **keeps its text**, so a reader can see what was withdrawn, and **cites no
goals**, since it serves none any more. It stays in its source and is marked `*(excluded)*` in a
composed package (`- [pkg.item](…) *(excluded)* — …`), with the deciding decision, so accounting
still finds it and `river check` rejects including or composing from it (`withdrawn-not-excluded`).
Nothing may cite it in `Covers:`; repoint a criterion that did
([`composition.md`](composition.md), *Replacing a local item*). Withdraw rather than delete: a
deleted item takes its history with it, and a link to it breaks without saying why.

## Conflicts

**Never resolve a conflict silently.** Two requirements that disagree go to `decisions.md`:
what conflicted, which sources, what was decided, by whom, and on what authority. Decision content
is append-only: a decision that turns out badly is superseded by a new one, never edited, because
the record of having been wrong is worth keeping. Only its `Status:` line changes by hand, while
`river upgrade` may migrate format across River major versions ([`decisions.md`](decisions.md)).
How to tell a conflict from an overlap, and which overlaps precedence settles, is in
[`composition.md`](composition.md), *Overlaps and collisions*.

## Anatomy

```markdown
## noSilentFailures

**Errors surface to the caller**

Subject: `api.errors`
Goals: `goals.trustworthyApi`

Every error reaches the caller. None is logged and dropped, and none is replaced by a default
value.

A swallowed error is indistinguishable from success, and the caller acts on a lie.
```

**The first paragraph after the fields is the complete statement**: what must be true, readable
on its own. Rationale, history and examples follow it, so a reader who needs only the rule can
stop there. Keep the rationale, after the rule: it tells a future reader a constraint is
deliberate, and it holds the incident a good failing test is built from.

Identifiers are descriptive CamelCase — `noSilentFailures`, not `P-014`. Inside a package its
own names are bare; everywhere else they are qualified, `product.noSilentFailures`. Package
names are unique within a River module, which is what lets a bare prefix resolve without an import
block.

**Every file that defines identifiers opens with `package <name>`** (`tech.declaredPackages`),
and that name is the prefix. A file in a package directory takes the directory's name; a file
standing alone takes its own, up to any hyphen: `decisions.md` is `package decisions`, and
`todo.md` and `todo-closed.md` are both `package todo`, so an item keeps its name when it closes.
Examples in this skill are the one exception.

**The heading is the identifier and nothing else; the title is the bold line beneath it.** A
heading's anchor is derived from its text, so this makes every item's anchor its identifier
lowercased — `product.md#nosilentfailures` — and lets the title be reworded without breaking
a single link to it. Write link fragments in lowercase; anchors are generated that way and
matched exactly. Decisions follow the same form.

Tell your agents to find one requirement by searching for `## noSilentFailures` and to read only
that section; a source file holding twenty requirements costs twenty times as much to read whole.
