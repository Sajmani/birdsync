# Composing requirements

Part of River's method. Search `^## ` for the sections; read only the one you need.

**Composition has three inputs, and they count equally: sources, precedence, and decisions.**
Sources say what is required; precedence says which source wins a conflict; decisions record what
a person settled that neither could. A composition that ignores any one of them is wrong.

The rules are in [`requirements.md`](requirements.md), sections *Precedence*, *Composition*, and
*Conflicts*, and in [`verification.md`](verification.md), *Acceptance criteria*. This file is the
procedure, and the judgment each step needs.

## The three inputs

1. **Sources:** every package under `river/requirements/local/` and
   `river/requirements/imported/`. A withdrawn source item is still read: it must be listed as
   excluded ([`requirements.md`](requirements.md), *Status of a source item*).
2. **Precedence:** `river/requirements/precedence.md`. **Mandatory** always wins. **Governing**
   beats local requirements unless its owner grants an exception. **Local** is the project's own.
   **Advisory** loses to local and fills in where local says nothing.
3. **Decisions:** `river/requirements/decisions.md`. Besides rulings on conflicts, it holds
   **interpretation decisions**, which say how to read an imported item without editing it: the
   class of things it applies to ([*Projecting horizontal requirements*](#projecting-horizontal-requirements)),
   whether one recommendation binds ([*Soft wording and vague numbers*](#soft-wording-and-vague-numbers)),
   and what an ambiguous number means. Each interpretation decision names the imported items it
   reads in `Affects:` and sits, in effect, beside the import: `imports.md` may link it from the
   import's entry, and anyone composing from that source applies it.

A composed package names its sources in a `## Composed from` list. The version of each import is
in `imports.md`, and the CLI records the content sums of what composition read in the River module's
`sums.md`, so a later change to a source is visible. Neither is repeated in the package.

**If a fact composition needs is in none of the three, stop and ask.** Record the answer as a
decision before composing on it. That is what makes a composition reproducible: anyone holding
the same three inputs gets the same result ([*Blind recomposition*](#blind-recomposition)).

## Procedure

1. **List the inputs** above.
2. **Review the decisions first.** List every decision whose affected items changed, were
   withdrawn, or were removed since it was decided, and ask the person whether each still holds
   ([*Reviewing decisions at recomposition*](#reviewing-decisions-at-recomposition)). Never
   apply one silently.
3. **Account for every source item.** Each is **included** unchanged, **composed from**, or
   **excluded**, in every composed package it applies to ([*Accounted for*](#accounted-for)), and
   that package names the item's source in its `## Composed from` list (required from River 2.0;
   step 4). Included unchanged is only for an item that already says what this project must do
   ([*Obligations become rules*](#obligations-become-rules)). An item over every instance of
   something is projected ([*Projecting horizontal requirements*](#projecting-horizontal-requirements)).
4. **Write the package** (named for its audience or scope, never a version; see
   [*Splitting into packages*](#splitting-into-packages)):
   - name every source it accounts for items from, and `decisions` where it rests on one, under
     `## Composed from`, one per line with its tier (`- product (Local)`); River 2.0 requires
     the list ([`module.md`](module.md), *Keeping the River module's files*);
   - list what applies unchanged, one line each, as a link to the source item plus its title,
     grouped by area (without repeating `Goals:`, which stays on each source item);
   - write out in full only what composition made: a project rule from an obligation, a rule
     over a domain, a rule narrowed by a higher tier, an intermediary's contract. Each such item
     cites the goals it serves in `Goals:` and has `Composed from:` linking every source item and
     decision it rests on (plain names in an exported package, which links only within its own
     directory);
   - mark each excluded item with `*(excluded)*` immediately after its link (`- [pkg.item](…) *(excluded)* — …`),
     usually under `## Excluded`, and link the decision that excluded it. A withdrawn item is
     excluded the same way, linking the decision that withdrew it.
5. **Check that every source item is accounted for**
   ([*Accounted for*](#accounted-for)). None may be omitted from a package it applies to, or from
   every package in the River module; one that is has silently dropped out. The CLI's accounting check
   reports unaccounted items (`unaccounted-requirement`), but count them yourself while composing.
6. **Write one acceptance criterion per independently checkable behavior**, each with `Covers:`
   naming what it covers, beside the requirements it covers
   ([*Files within a package*](#files-within-a-package)). Drop any criterion whose covered
   requirements were withdrawn or excluded (`withdrawn-criterion`, `covers-withdrawn`,
   `covers-excluded`), and clean up its `conformance.md` section and obsolete guardrails
   (`guardrail-withdrawn`). Where a composed item replaces a local one, keep the local item's
   criteria ([*Replacing a local item*](#replacing-a-local-item)).
7. **Where precedence does not decide a conflict, stop.** Ask the owner and record the ruling in
   `decisions.md` before composing on it ([*Overlaps and collisions*](#overlaps-and-collisions)).
8. **Questions and drafts go in `river/proposals.md`**, never into the composed package
   ([`module.md`](module.md), *Proposals*).
9. Run `river sum -w` (River 2.0), so `sums.md` records every input composition read (sources,
   decisions, precedence), then `river check --allow-unapproved`; review `river diff` with the
   owner, and once the owner approves, run `river approve --by <owner>` and `river check`
   ([`module.md`](module.md), *Keeping the River module's files*).

## Accounted for

**An input item is accounted for in a composed package when the package does one of three things
with it:** lists it unchanged, as a link with its title; cites it in a composed rule's
`Composed from:`; or marks it `*(excluded)*` after the item link and links the decision that
excluded it (`- [pkg.item](…) *(excluded)* — …`). A mention in prose, or in a criterion's
`Covers:`, does not account for it, and a list line without `*(excluded)*` counts as an inclusion
even if its prose says "not applicable" (`note [ambiguous-exclusion]`).

**Which packages owe accounting depends on what each package governs:**
- **One composed package:** it represents the whole project and must account for every local and
  imported source item in the River module.
- **Several composed packages** (such as `app` and `dataset`, or exported `framework` and internal
  `build`): an item with `Applies to:` (on the local item, or in an interpretation decision for an
  imported item) must be accounted for in **every** composed package that represents or `Binds:`
  that target; an item without `Applies to:` must be accounted for in **every** composed package
  whose `## Composed from` names its source; and every source item in the module must be accounted
  for in at least one composed package. `river check` reports any gap as `unaccounted-requirement`.
## Obligations become rules

**Include a source item unchanged only if it already says what this project must do.** An
obligation phrased for someone else, such as an operator, a user, or "the application", is
composed into a rule for this project, and the rule cites the obligation in `Composed from:`.
"Use OAuth" becomes *the client authenticates with OAuth 2.0 authorization-code flow and never
asks for a password*; "don't machine-generate observations" becomes a rule about what the tool
lets a user upload. Composition writes that rule into the composed package; a local source item
is needed only for what the project adds beyond it.

**A horizontal obligation, one over every web page or every service, is composed once, over its
domain, and never restated per instance.** Twenty pages under six accessibility rules is six
rules, not 120. Its criteria test the domain **exhaustively or by sampling**: every kind of
instance and every input that varies the property, and a sample of the instances
([`verification.md`](verification.md), *Verification and guardrails*). The procedure is
[*Projecting horizontal requirements*](#projecting-horizontal-requirements).

**Where an intermediary carries the rule, the composed rule cites its contract.** A stylesheet,
a framework, or a middleware that makes its users meet the obligation states a contract
([`requirements.md`](requirements.md), *Composition*). The project's composed rule says the
obligation holds through that contract, and the project's own requirement is the contract's
assumption, such as *every request goes through the framework*.

## Projecting horizontal requirements

**A horizontal requirement is mapped over classes of things, never over a list of instances.**
Six steps:

1. **Recognize it.** "All", "any", "every", "each", or a rule stated for a kind of thing (a page,
   a request, a record, a dataset) marks a horizontal requirement.
2. **Name the classes it applies to**, in the source's own term: `Applies to: wcag.WebPage`. If
   the source says, its transcription carries the field. If it doesn't, record the class in an
   **interpretation decision** whose `Affects:` names the imported items, and link it from the
   import's entry in `imports.md`. Never edit imported text to add the field.
3. **Classify the project's entities.** A local requirement with a `Binds:` field naming the
   class says what counts as an instance here, and settles the edge cases: is an error page a
   web page? Is a health check a request? The binding makes the term concrete and **never narrows
   it**: leaving out something the source's definition includes is an exception, recorded
   instance by instance under the source's tier.
4. **Map requirements over entities by class.** The composed package states each rule once, over
   its class, citing the requirement, the binding, and any interpretation decision. **Nothing
   lists the instances.** Verification discovers them from the system, from what it declares
   (routes, templates, handlers) and by following links, so a new instance inherits its
   obligations by existing.
5. **Several horizontal rules on one class:** where they agree, each applies. Where they
   conflict, precedence settles it; where precedence can't, a person does, and the ruling is a
   decision ([*Overlaps and collisions*](#overlaps-and-collisions)).
6. **Intermediaries carry contracts.** Where a stylesheet or framework makes every instance meet
   a rule, its contract states the guarantee, the assumption, and the scope, and the project
   checks the assumption for each kind of user, not the rule for each instance
   ([`verification.md`](verification.md), *Verification and guardrails*).

**A worked case.** Acme's brand group publishes two policies, `accessibility` (WCAG 2.2 contrast,
Mandatory) and `brandstyle` (the palette, Governing), and a stylesheet whose contract, `brandkit`,
meets both for a page that takes its colors only from the stylesheet. Acme's policy group publishes
`servicelogs`, over every request a service serves, and the server framework's contract,
`framework`, meets it for every request routed through it. Greeter, a small web service, imports
all five:

- **Recognize:** "every piece of text on an Acme page" and "every request is logged" are
  horizontal.
- **Classes:** the policies don't name a class, so greeter records an interpretation decision:
  the accessibility and brand items apply to an Acme page, and the logging items to a request.
- **Bind:** greeter's local requirements bind those classes: every HTML response greeter serves
  is an Acme page, including its error pages; every HTTP request it receives is a request.
- **Map:** greeter's composed package states that each policy holds over its class, through the
  contract that carries it, and lists no page and no route.
- **Intermediaries:** greeter checks the contracts' assumptions (no color of its own; nothing
  logged outside the framework), plus the contrast of the page as served, end to end. The
  contracts' guarantees are checked once, in the stylesheet's and the framework's own
  repositories.

## Soft wording and vague numbers

**Settle how strongly an imported item binds, and what its numbers mean, at import, by
decision.** Precedence ranks sources; it doesn't say how firmly "please consider" binds, or how
to check "about one request per second".

- **A recommendation in a Governing or Mandatory source** ("should", "please consider") may be
  made **Advisory for that item**: followed by default, losing to a local requirement where the
  two disagree. A person decides it, in a decision saying who and why, and the tier itself goes in
  a row of `precedence.md`'s item table that links that decision
  ([`requirements.md`](requirements.md), *Precedence*). The rest of the source keeps its tier. A
  local requirement that then beats the item is ordinary precedence: it needs no field naming the
  item, and the item table is where a reader sees every item whose tier was changed.
- **An ambiguous number** ("about one per second": 1.5 per second at peak? bursts averaging one
  per second over ten minutes?) is resolved by an interpretation decision that states the
  checkable reading, which the criterion then tests.
- **What import left open, composition resolves**, also as a decision, and asks a person when the
  reading isn't obvious. Never pick a reading silently: the strictest one may be unbuildable, and
  the loosest one may break the source's terms.

## Restatements

**A local item that restates an imported or composed rule is redundant, not wrong.** A project
may keep one for its own reasons: a page that was accessible before it adopted an accessible
framework still has its own requirement and its own tests. Composition notes each restatement it
finds, the CLI reports it as a note, and a person decides whether to withdraw the local item
([*Replacing a local item*](#replacing-a-local-item)). Keep a local item for what it adds beyond
the rule it restates, such as the version string a User-Agent must carry.

## Overlaps and collisions

**When two sources state rules on the same subject, treat it as a name collision, and resolve
it the way programming languages that have proven these rules do.** This is guidance; the
complete rule set is still open.

- **Different tiers:** the higher tier wins, as the shallower field wins under Go embedding.
  Where the lower tier's rule is stricter and compatible, both apply.
- **Same tier, same subject, no decision:** an error until a person resolves it, as an ambiguous
  selector at equal depth is in Go, and as a conflicting method from two traits or mixins must be
  resolved explicitly. Ranking within a tier (*Local: tech* above *Local: product*) settles it if
  the ranking is recorded in `precedence.md`.
- **Horizontal rules** resemble aspects: the class is the pointcut, and precedence orders the
  advice. Two rules over one class are composed in precedence order, and a conflict between them
  is settled as above.
- **One source item, accounted for twice in one package:** don't. If composition wrote a rule from
  an item, or excluded it with `*(excluded)*`, don't also list it as included unchanged: doing two
  of {included, composed from, excluded} in one package leaves a reader unable to tell whether or
  how it governs. From River 2.0, `river check` fails a package that does (`included-and-composed`).

## Reviewing decisions at recomposition

**A decision is an input only while it still makes sense.** At each recomposition, list every
decision whose `Affects:` names an item that changed, was withdrawn, or was removed since the
decision was made, and ask the person whether each still holds. The CLI notes such decisions
(`decision-review`), using the per-item sums in `sums.md`. A decision that no longer holds is
superseded by a new one, never edited ([`decisions.md`](decisions.md)). A decision whose
`Affects:` names an item that doesn't exist fails the check.

## Blind recomposition

**Check that composition reproduces what people wrote by hand.** Hide the local items you suspect
restate imported rules, recompose from the sources, precedence, and decisions alone, and compare.
What reappears was a restatement. What doesn't reappear was either something the project adds,
so keep it, or a fact missing from the inputs, so record it as a decision. Run it after a large
migration, or when the composed package has grown by hand-composition over many changes.

## Splitting into packages

**A composed package is the unit of approval and of import; split along those lines and nowhere
else.** Split when parts have:

- a **different importer or audience**: the exported versus internal line;
- a **different approver**, such as compliance-owned rules beside product rules;
- a **separate releasable artifact**, such as a library and the service that uses it;
- a **different change cadence**: approving one part keeps re-approving a stable other part.

**Don't split by source** (composition merges sources), **by version** (a version is a commit),
**or by size alone** (use files for that).

**Splitting has costs.** Each composed package owes accounting for every source item that applies
to the entity it represents (via `Applies to:` and `Binds:`, or via `## Composed from`), so shared
rules that govern two packages must be accounted for in both ([*Accounted for*](#accounted-for)).
Criteria live with the package whose requirement they
cover. A requirement composed in one package and needed by another is cited across packages, and
the method doesn't yet define how one composed package declares that it imports another, or
whether a cross-package citation pins a sum. Until it does, prefer one package wherever the four
reasons above don't apply.

## Files within a package

**Files organize; the package is still approved as one.** Splitting a package into files never
changes what is approved, only how it reads.

- **Split by area**, mirroring the source files: a source split into `posting.md` and `auth.md`
  gives the package `posting.md` and `auth.md`.
- **Each area file holds that area's included list, its composed items, and its criteria**, so a
  reader of one area sees each rule and its check together. Included bullets link to their source
  items, where their `Goals:` fields live; composed items carry their own `Goals:` fields.
- **A package `README.md`** carries the `## Composed from` list, the precedence applied, and the
  `## Excluded` list.
- **Split a file** when it outgrows a few hundred lines, or mixes areas a reader wouldn't load
  together.

**For a small project with one artifact, use one package, split into area files**, until it gains
a second approver or one part keeps being re-approved for changes elsewhere.

## Replacing a local item

**When a composed item replaces a local one, withdraw the local item and keep its criteria.**

1. Mark the local item `Status: withdrawn — <why>`, citing the decision, and exclude it in the
   composed package with `*(excluded)*` (`- [pkg.item](…) *(excluded)* — …`), linking the decision
   ([`requirements.md`](requirements.md), *Status of a source item*).
2. Repoint each criterion's `Covers:` from the local item to the composed one, and repoint the
   criterion's rows in `river/conformance.md`.
3. **Keep the existing criteria rather than writing near-duplicates.** Their guardrails and their
   watched-failing evidence survive only if the criteria do. Add a criterion only for behavior the
   composed rule requires that no existing criterion checks.
