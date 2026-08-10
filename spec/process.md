# Spec-driven development

This document describes *how to work on a software project*, not what any particular project
does. It defines a three-phase loop — **context2spec → spec2test → spec2code** — that turns
requirements into verified code, and names the artifacts each phase produces.

It is deliberately project-independent, so the same copy can seed a new project or be retrofitted
into an existing one. Everything specific to a given project — build commands, safety
constraints, language conventions — is confined to a short set of **[project
bindings](#project-bindings)** recorded elsewhere, leaving this file identical across projects
and cheap to update from a canonical copy.

It is also a living document. When a step stops earning its keep, change this file first and then
follow the changed version. See [Revising this process](#revising-this-process).

## Why

Requirements that exist only in someone's head can't be checked against the code, so they drift,
and the drift stays invisible until a user hits it. Documentation drifts the same way, in the
other direction: it describes what the code used to do.

Spec-driven development makes requirements a reviewable, referenceable artifact, and then makes
"is this requirement actually enforced?" a question with a mechanical answer. The payoff is
largest where mistakes are expensive or hard to see: anything that writes to a user's data,
handles money or credentials, is hard to roll back, must satisfy an outside authority, or is
maintained by people — or agents — who weren't there when the decisions were made.

## The loop

```
  context2spec ─ gather ─► compose ─► resolve
       │                                 │
       │        product.md, tech.md, sources.md, decisions.md
       │                                 │
       │                                 ▼
       │                    GATE 1: human approves the requirements
       ▼
  spec2test    ──►  acceptance.md + checks ──►  GATE 2: human approves the criteria
       │
       ▼
  spec2code    ──►  code, tests, arch.md
       │
       └───────────►  next change re-enters at context2spec
```

## The artifacts

| File | Answers | Produced by | Nature |
| --- | --- | --- | --- |
| `process.md` | How do we work? | — (this file) | Normative |
| `product.md` | What must the system do, from a user's point of view? | context2spec | Normative |
| `tech.md` | What constraints must the implementation satisfy? | context2spec | Normative |
| `sources.md` | Whose requirements apply here, and who outranks whom? | context2spec | Normative |
| `decisions.md` | How was each conflict resolved, by whom? | context2spec | Historical |
| `acceptance.md` | How do we know each requirement holds? | spec2test | Normative |
| `arch.md` | How is the code actually put together? | spec2code | Descriptive |

The normative/descriptive split decides who wins an argument:

- **A normative file disagrees with the code → the code is wrong**, or the requirement was, and
  needs an explicit, acknowledged amendment.
- **`arch.md` disagrees with the code → `arch.md` is wrong.** It documents reality, including
  reality nobody is proud of.
- **`decisions.md` is append-only history.** Entries are never rewritten; a decision that turned
  out badly is superseded by a new entry, not edited.

These sit alongside whatever docs the project already has — a README, a contributor guide, an
agent instruction file. Those keep their jobs. One rule prevents the set from rotting: **a fact
lives in exactly one file and the others link to it.** `product.md` is normative ("the system
must…"); a README is instructional ("run it like this"). When they state the same fact twice,
they will eventually disagree, and nobody will know which one to believe.

Size these to the project. A small tool with no external sources may leave `sources.md` and
`decisions.md` nearly empty, and should. The gates matter more than the volume; manufacturing
ceremony is its own failure mode.

### Where they live

The whole bundle goes in a `spec/` directory at the repository root, keeping the root for the
files people and tools expect to find there:

```
README.md            entry point for users
CONTRIBUTING.md      entry point for contributors
AGENTS.md            entry point for agents (tooling expects it at the root)
spec/
  process.md         this file
  product.md
  tech.md
  sources.md         the adoption manifest
  decisions.md       the CR-### log
  acceptance.md
  arch.md
  sources/           vendored copies of by-reference external sources
    <name>/
      PROVENANCE.md    where each file came from, its hash, who fetched it
      requirements.md  the transcription, as <name>/R1, R2, …
      *.html           the documents, exactly as served
```

### Inside a vendored source

One directory per **adopted source** — a named body of rules with one tier, one owner and one
scope — not one per document. A source often spans several pages, and how many says nothing
about it: a publisher may put its rules on one page or split them across three and a linked
definition.

| File | Contents |
| --- | --- |
| `PROVENANCE.md` | Every file's origin URL, its SHA-256, the retrieval date, **who or what retrieved it and how**, any upstream revision date, and what is deliberately *not* vendored |
| `requirements.md` | The transcription: `<name>/R1`, `R2`, … each quoting the source verbatim, then saying what it means for this project. Ends with what was considered and not adopted |
| the documents | As served, except that **active content is removed** — `<script>`, `<style>`, `<link>`. Never reformatted, never reworded, never partially quoted |

**Strip active content, and record that you did.** A page saved from a browser carries the
publisher's scripts, and scripts carry the publisher's secrets: vendoring four iNaturalist
pages committed their Google Maps API key to a public repository and tripped secret scanning.
Nothing there was private — the key is in every page they serve — but redistributing another
party's credential is nobody's idea of good citizenship, and a rules snapshot has no use for
JavaScript. Record both hashes, as served and as stored, so the removal can be reproduced and
audited.

**Check that quotations survive.** Once documents are edited at all, a transcription can quote
something the stored copy no longer contains. An automated check that every quoted passage
appears in its source costs little and is worth more than it looks: run against the
transcriptions here it found the stripping to be clean, and found two places where an earlier
transcription had quietly improved the text it was quoting — swapping quotation marks, and
joining two bullet points with a full stop that the source does not contain.

**The directory name carries no version.** Put the pin in `PROVENANCE.md`, as a hash and a
date. Naming the directory for its vintage seems tidy and quietly defeats the purpose:
refreshing then adds one directory and deletes another, so the diff — the thing vendoring
exists to produce — becomes unreadable, and stale copies pile up. An unversioned directory
makes a refresh an in-place change whose diff is exactly what the publisher altered. Git keeps
the old bytes.

Four reasons to group rather than scatter these at the root:

- **The bundle already needs a directory.** Vendored sources can't live at the root sanely, so
  half the bundle would end up nested anyway.
- **One directory is one review rule.** A CODEOWNERS entry on `spec/` is how Gate 1 gets enforced
  mechanically: changes to requirements require the owner's approval, changes to code don't.
- **Seeding a new project is a directory copy**, and `process.md` sits at the same path in every
  project, so improvements to it propagate with a one-line command.
- **It draws the normative/instructional boundary in the filesystem**, which makes the
  no-duplication rule easier to keep.

The cost is discoverability: someone browsing the repository sees the root and may never open
`spec/`. Link to it from the README and the contributor guide.

Two variations are fine, and both are [project bindings](#project-bindings): nesting under an
existing documentation tree (`docs/spec/`), and leaving a pre-existing `arch.md` at the root
rather than churning links during a retrofit. What matters is that the location is recorded and
consistent, not which location it is.

## Anatomy of a requirement

A requirement is a heading with an ID, a one-sentence statement, and whatever of the following
apply:

```
### P-020 — Body text color
Subject: page.text.color
Value:   #D3D3D3
Status:  Superseded by brand/P-009 (CR-007)

Body text on content pages is light grey.

Rationale: reduces glare against large light areas in the reading view.
```

- **Statement** — one testable assertion. If it needs an "and," split it.
- **Rationale** — why, not how. This is what lets a future reader tell an intentional constraint
  from an accident.
- **Status** — `Active` by default; otherwise `Superseded by …`, `Withdrawn (reason)`, or
  `Deferred (owner)`.
- **Subject and Value** — required *only* for requirements that fix a value on a shared property,
  optional everywhere else. This is the structure that makes conflict detection mechanical, and
  it's deliberately not demanded of prose requirements: the small fraction of requirements that
  set values is where nearly all conflicts live.

A **subject** is a dotted path from a vocabulary the project maintains (`page.text.color`,
`api.timeout.seconds`, `build.minimum_runtime_version`). The vocabulary is a project binding.
Two requirements sharing a subject are talking about the same thing, which is exactly what a
conflict check needs to know.

## Requirement IDs and provenance

Every requirement and criterion has a stable identifier, and every requirement has a source.

| Form | Meaning | Example |
| --- | --- | --- |
| `P-###` | Local product requirement | `P-012` |
| `T-###` | Local technical requirement | `T-007` |
| `AC-###` | Local acceptance criterion | `AC-031` |
| `CR-###` | Conflict resolution record | `CR-007` |
| `<source>/<id>` | Requirement from an adopted external source | `brand/P-009`, `wcag/1.4.6` |

Rules:

- **IDs are allocated monotonically and never reused.** `P-014` means one thing forever.
- **Requirements are never renumbered or deleted.** One that no longer applies is marked
  `Withdrawn` or `Superseded` in place. Deleting it silently orphans every reference to it in the
  code, the tests, the decision log, and the history.
- **External sources keep their native numbering.** Import WCAG's `1.4.6` as `wcag/1.4.6`, not as
  a renumbered local requirement. Renumbering destroys the traceability back to the authority,
  which is the main reason to cite an outside source at all.
- **Every requirement is testable as written.** "The system should be fast" is not a requirement;
  "a request at the 99th percentile completes within 300 ms under the standard load profile" is.
- **Requirements say what and why, not how**, unless the how *is* the constraint — which is what
  `tech.md` is for.

Code and tests cite IDs in comments, so the trace from requirement to enforcement is a search:

```
// Verifies: P-012, brand/P-009.
```

Use the host language's comment syntax, but keep the literal token `Verifies:` so one search
covers the whole tree regardless of language. Cite IDs on tests as a matter of course; cite them
in non-test code only where the reason for the code is non-obvious without one — a safety gate, a
workaround for a third-party quirk, a branch that exists solely to satisfy a stated constraint.

## Sources and composition

Requirements arrive from more than one place: the project's own owner, company-wide brand or
engineering policy, language and platform best practice, and outside authorities such as
regulation. Composition is the act of assembling them into one effective requirement set.

### Precedence tiers

Each adopted source is declared at one of four tiers. The tier determines what happens on
conflict and who may grant an exception.

| Tier | On conflict with another tier | Exception granted by | Typical sources |
| --- | --- | --- | --- |
| **Mandatory** | Always wins; never overridden locally | Nobody locally — see [unsatisfiable mandatory requirements](#when-a-mandatory-requirement-cannot-be-satisfied) | Law, regulation, accessibility standards, org security policy |
| **Governing** | Overrides local requirements | The source's owner, on request, recorded in `decisions.md` | Brand guidelines, company engineering standards, platform store rules |
| **Local** | The project's own `product.md` / `tech.md` | The project owner | Everything the project decides for itself |
| **Advisory** | Loses to local; applies wherever local is silent | The project owner, with a recorded rationale | Language and framework best practices, style guides, internal cookbooks |

Two observations about this table:

- **Advisory sources shape the spec rather than constrain it.** Adopting a language best-practice
  guide as advisory means its requirements become the default for everything the project hasn't
  decided, and the project may depart from any of them by saying so. That "by saying so" is the
  point: a recorded departure is a decision, an unrecorded one is a mistake.
- **The tier is a property of the adoption, not of the source.** The same brand guide may be
  governing for a shipping product and advisory for a prototype. Declare it per project.

### The source manifest

`sources.md` records every adopted source. For each one:

| Field | Meaning |
| --- | --- |
| Name | The ID prefix used in citations (`brand`, `wcag`, `go`) |
| Origin | Repository, URL, or document, precise enough to re-fetch |
| Version | Tag, commit, edition, or retrieval date — a **pin**, not "latest" |
| Integrity | Content hash of the vendored copy, where the format allows one |
| Tier | mandatory / governing / local / advisory |
| Import style | By reference or absorbed (below) |
| Scope | Where it applies — the whole project, or a named subset |
| Adopted parts | Which sections or conformance level, when adoption is partial |
| Exclusions | Requirements deliberately not adopted, each with a reason |
| Owner | Who to ask for an exception |

Also record the sources you considered and rejected. "We are not subject to X" is a decision that
someone will otherwise re-litigate every year.

### Import style

**By reference** — vendor a pinned copy into the repository (`spec/sources/<name>/`) and cite
its requirements with their native IDs. Upstream changes arrive only when
someone deliberately refreshes the pin, and they arrive as a reviewable diff. This is the default
for mandatory and governing sources, where you want to know about updates.

**Absorbed** — copy the guidance in at adoption time as local requirements, with the origin noted
in each one's rationale. The link to upstream is cut: later changes there won't reach you. This is
appropriate for advisory sources where stability beats currency, and for prose guidance that has
to be transcribed by hand anyway.

Vendor rather than fetch at check time. The spec stays self-contained and diffable, checks don't
need network access, and an upstream edit can never silently change what your project requires.

**Transcription is an artifact.** Regulations and design guidelines rarely come machine-readable,
so importing them means a human or agent turning prose into numbered requirements. That
transcription can be wrong. Record who produced it, from which document and revision, on what
date, and treat it as reviewable material rather than as ground truth.

**Some sources refuse to be fetched.** A publisher may block automated retrieval — a bot
challenge, a login wall, a paywall, a PDF behind a click-through. Then a human retrieves the
document and the provenance records that, including who did it and how. Do not work around
the block: a service that has asked not to be scraped has asked, and evading it to import its
own rules is self-defeating. Above all, do not fill the gap from memory. An agent asked for a
rate limit it cannot fetch will produce a plausible number, and a plausible number in a
requirements document is indistinguishable from a real one. Leave the requirement unwritten
and the source marked unpinned until someone can supply the text.

### Scope

A source, or an individual requirement, may declare a scope: the components, surfaces, or
configurations it governs. Precedence only applies where scopes overlap. A more specific
requirement does **not** automatically beat a more general one from a higher tier — a governing
rule yields to a local one only where the governing source's own declared scope excludes the case.
Otherwise, "but ours is more specific" becomes a way to launder any override.

## Conflicts and resolution

A conflict exists when the composed requirement set admits no implementation. Kinds differ, and
so does how you find them:

| Kind | Example | Detected by |
| --- | --- | --- |
| Direct contradiction | Two local requirements set the same subject to different values | Grouping by subject |
| Cross-source override | Local value versus a governing source's value | Grouping by subject; resolution is mechanical |
| **Derived / emergent** | Text color and background color each fine alone, but their contrast violates an accessibility rule | **Evaluating a rule over the resolved set** — no pair contradicts |
| Conditional | The conflict appears only in dark mode, a locale, or a deployment tier | Enumerating the configurations the spec declares |
| Over-constraint | No value satisfies every applicable constraint | Attempting resolution and failing |

The emergent row is why conflict detection can't be careful reading alone. "Text is light grey"
does not contradict "background is teal"; the violation exists only in a rule computed over both.
Catching that class requires [spec-level criteria](#spec-level-criteria).

### Resolution order

1. **Group by subject.** Every subject with more than one active value is a candidate conflict.
2. **Apply tier precedence.** Different tiers: the higher tier wins, the loser's status becomes
   `Superseded by <id> (CR-###)`. Same tier: precedence decides nothing — escalate.
3. **Report every mechanical resolution anyway.** An override that silently voids a stated
   requirement is the most dangerous outcome in this whole process, because the requirement's
   author still believes it holds. Two local requirements contradicting each other stay worth
   reporting even when a governing source makes both moot: it means somebody misunderstood
   something, and that misunderstanding is probably not confined to one line.
4. **Evaluate spec-level criteria** over the resolved values to surface emergent conflicts.
5. **Escalate what's left to a human, with candidate resolutions.** Never fewer than two, each
   showing which requirements it satisfies, which it violates, and what it costs. The arithmetic
   is mechanical; the choice between "ask the brand owner for a deviation" and "adopt the lower
   conformance level" is a judgment with organizational consequences, and belongs to a person.
6. **Record the outcome** as a `CR-###` entry in `decisions.md`, and update the status of every
   requirement it touched.

An escalation frequently reveals that the *adoption* was underspecified rather than that the
requirements are wrong — "we follow the accessibility standard" without naming a conformance
level, for instance. Fixing `sources.md` is a legitimate resolution.

### Resolution records

Each `CR-###` in `decisions.md` records: the subject, the conflicting requirement IDs and their
sources, the kind of conflict, whether it was resolved mechanically or by escalation, the options
considered, the decision, the person who made it, the date, the resulting status of every
requirement involved, and any follow-up owed to a source owner.

Sources are never edited to resolve a conflict. Vendored copies stay byte-identical to upstream;
local requirements are annotated in place, never rewritten to pretend the conflict didn't happen.

**Evidence and inference are labelled separately, and evidence is checked before it is cited.**
A plausible causal story — *this cleanup tool exists, so the bug it cleans up must be
happening* — is a hypothesis. It belongs in the record as one, next to the check that would
settle it. Most such stories are settled in one command, because version control knows when
each thing was written and in what order: a tool that predates the code it supposedly
compensates for cannot be evidence for it. A guess promoted to "supporting evidence" is worse
than no evidence, because the next reader cannot tell which is which and will build on it.

### When a mandatory requirement cannot be satisfied

Work stops. There is no local mechanism to override a mandatory source, and "note the gap and
carry on" is how non-compliance becomes routine bookkeeping.

Two ways forward: change the other constraints so the requirement can be met (usually by asking a
governing source's owner for a deviation), or record an explicit **risk acceptance** in
`decisions.md` naming the person who accepted the risk, their authority to do so, the exposure
they accepted, an expiry date, and a review date. An unsigned, unexpiring gap row in a table is
not a risk acceptance.

### Worked example

```
Sources:  brand/    governing   pinned v4.2
          wcag/     mandatory   2.2, Level AAA, adopted in full
Local:    P-012  page.background.color = blue
          P-013  page.background.color = green
          P-020  page.text.color       = #D3D3D3   (light grey)
```

- Subject `page.background.color` has three active values: `P-012`, `P-013`, and `brand/P-009`
  (teal `#008080`). Brand is governing, so teal wins and both local requirements are superseded —
  **and the `P-012` / `P-013` contradiction is still reported**, because two local requirements
  disagreeing means a shared misunderstanding.
- A spec-level criterion for `wcag/1.4.6` computes contrast over the resolved values: light grey
  on teal is roughly 3.2:1 against a required 7:1. Mandatory tier, so the standard cannot yield.
- No tier ordering resolves the remainder, so it escalates with candidates:

  | Option | Contrast | Cost |
  | --- | --- | --- |
  | White text on brand teal | ≈4.8:1 — still fails | Non-compliant unless the real target is Level AA |
  | White text on a darkened teal `#004D4D` | ≈9.7:1 ✓ | Needs a brand deviation or a brand-approved dark variant |
  | Near-black text on brand teal | ≈4.4:1 — fails | — |
  | Adopt Level AA (4.5:1) instead of AAA | White on teal passes | Changes the compliance posture — a real decision, not a workaround |

- The human chooses; `CR-007` records the choice, the options rejected, and the fate of `P-012`,
  `P-013`, and `P-020`.

Note that the last option is a change to `sources.md`, not to any requirement. Escalations often
land there.

## Two entry modes

**Greenfield.** Requirements come from the person who wants the system. Phase 1 is an interview,
and the main risk is inventing requirements nobody asked for. Write down non-goals as
deliberately as goals, and identify the mandatory and governing sources early — discovering a
governing source late invalidates work.

**Retrofit.** The system already exists, works, and has users. Phase 1 reverse-engineers the
requirements from the code and docs, then confirms them with the maintainer; phase 2 backfills
criteria that lock in the behavior users already depend on. The retrofit will surface behavior
nobody intended — finding it is the point. Record those as open questions rather than silently
"fixing" them: in a system with users, undocumented behavior may still be load-bearing. Expect the
first composition against external sources to produce a burst of conflicts representing years of
undetected drift; resolve them in priority order rather than all at once.

In both modes, from the first approved spec onward, every behavior change starts with a spec edit.
Code that traces to no requirement is either an undocumented requirement (write it down) or an
accident (remove it).

## Phase 1 — context2spec

**Goal:** a complete, agreed, conflict-free statement of what the implementation must satisfy.

**Inputs:** the existing code; existing docs; version-control history, issues, and post-mortems
(bugs are requirements that were never written down); the specifications and quirks of external
systems the project integrates with; the external sources the project is subject to; and the
humans who own the project.

### Gather

1. Survey the code and docs; extract the behavior actually implemented and the constraints
   actually honored.
2. Identify every source with a claim on the project — regulatory, organizational, platform,
   best-practice. Ask the owner what the project is subject to; this is not inferable from code.
3. Sort each finding into product (observable by a user) or technical (a constraint on how it's
   built).
4. Write them as numbered, individually-testable requirements, with subjects and values on the
   ones that fix a value.
5. **Ask the owner** about everything the code can't tell you: which behaviors are intentional
   versus incidental, what's out of scope, which quirks are load-bearing, what the system must
   never do. Batch the questions rather than dripping them out one at a time.
6. Record what you could not resolve under "Open questions" in the relevant file. An unresolved
   question is a legitimate output of this phase; a guess dressed up as a requirement is not.
7. **Record which inputs were actually consulted**, by name, in the spec. Listing inputs in a
   method does not cause anyone to read them, and a skipped one is invisible afterwards — the
   spec looks equally complete either way. The issue tracker is the one most often missed and
   the most likely to hold a defect someone has already reported, and the context that
   explains why the code looks the way it does.

### Compose

7. Adopt each source in `sources.md`: tier, pin, import style, scope, adopted parts, exclusions,
   owner. Vendor by-reference sources; transcribe absorbed ones with attribution.
8. Assemble the effective requirement set and group it by subject.

### Resolve

9. Work the [resolution order](#resolution-order). Escalate what precedence can't decide, with
   candidate options and their costs.
10. Record every resolution as a `CR-###`, and update the status of each requirement involved.

**Outputs:**

`product.md` — the user-visible contract: what a user can ask the system to do; its interface
surface and compatibility promises; what it produces, writes, or sends; its filtering, ordering,
and precedence rules; idempotence and retry semantics; observable failure behavior and error
codes; and explicit **non-goals**.

`tech.md` — the constraints on the implementation: safety invariants that must never be violated;
security and privacy obligations; dependency and supply-chain policy; supported platforms,
runtimes, and versions; testability seams the design must preserve; performance, memory, and cost
budgets; error-handling and observability conventions; backward-compatibility obligations; and the
[project bindings](#project-bindings). Any pre-existing project rules — the ones in an agent
instruction file or contributor guide — are promoted here into numbered requirements so they
become checkable rather than folkloric.

`sources.md` — the adoption manifest, including rejected sources.

`decisions.md` — the `CR-###` log.

**→ GATE 1: the owner reviews and approves `product.md`, `tech.md`, and `sources.md` before phase
2 begins.** Precondition: no unresolved conflicts. A conflict may pass the gate only as a
`Deferred` entry naming an owner and a date — never by being left unmentioned.

## Phase 2 — spec2test

**Goal:** for each requirement, a specific check that would fail if the requirement were violated.

**Method:** for each requirement, choose the *cheapest checker that would actually catch the
violation*:

| Method | Good for |
| --- | --- |
| Unit test | Pure logic: parsing, transformation, key construction, precedence rules |
| Integration / contract test | Request and response shapes, wiring, end-to-end flows against fakes |
| Property-based test / fuzzing | Invariants over generated inputs; parsers and serializers |
| Static analysis | Structural invariants a compiler or linter can enforce; forbidden imports, calls, or strings |
| Type system / schema | Making an illegal state unrepresentable rather than checking for it |
| Dynamic analysis | Races, leaks, memory errors, undefined behavior |
| Benchmark / budget | Throughput, latency, allocation, and cost ceilings |
| Model checking | Ordering and state invariants too broad to enumerate by example |
| **Spec-level check** | Rules computed over the requirement set itself (below) |
| Agentic code review | Judgment calls that resist encoding: "is every mutation behind the safety gate?" |
| Agentic product testing | Spec-driven exploration of the running system against a fake backend |
| Human review | Taste, product intent, and anything with irreversible real-world consequences |

### Spec-level criteria

Some criteria have the **resolved requirement set** as their subject rather than the code. A
contrast ratio computed from two declared color values can be checked before a line of code
exists; so can a latency budget that must be the sum of its parts, a version-support matrix that
must be non-empty, or a rule that every user-visible string has a localization requirement.

Every criterion therefore declares a **level**: `spec` or `code`. Spec-level criteria run during
context2spec's resolve step and again whenever a requirement or a source pin changes. They are the
only mechanism that catches emergent conflicts, and requirements imported from mandatory sources
are usually best expressed as them.

### Guidance

- **Prefer executable checks.** Human review is the fallback, not the default; it doesn't run in
  CI, and it doesn't run at 2am.
- **A criterion that can be promoted, should be.** When someone works out how to automate a
  human-review criterion, change its method and note the change.
- **Name the command.** Every criterion states exactly how to run it. "Reviewed manually," without
  a procedure, is not a criterion.
- **A criterion must be able to fail.** Before trusting a new check, break the code — or the spec
  — and watch it go red. A check that passes against a deliberately broken subject is worse than
  no check, because it is believed.
- **Test the requirement, not the implementation.** A check that asserts today's internals will
  block tomorrow's refactor while catching none of the violations that matter.
- **Derive the assertion from the requirement, not from a run.** Writing a check by executing
  the code and recording what came back produces something that passes by construction and
  certifies whatever the code does, bugs included. It is the standard way a retrofit goes
  wrong, because the behavior is right there and the requirement has to be thought about. The
  tell is an assertion shaped like the output — a set of accepted values, a golden file nobody
  read, a tolerance wide enough to fit the observed answer.
- **Cover imported requirements too.** An adopted mandatory or governing requirement without a
  criterion is a compliance claim nobody is checking.
- **Coverage is tracked, not assumed.** `acceptance.md` carries a traceability table mapping every
  requirement — local and imported — to the criteria that verify it. Requirements with no
  criterion are listed explicitly as gaps. An honest gap beats a criterion that doesn't bite.

Add meta-criteria for composition itself: every source pinned and vendored; every vendored copy
matching its recorded hash; no requirement in `Deferred` past its date; every `CR-###` referencing
requirements that exist; every risk acceptance unexpired.

**Outputs:** `acceptance.md`, plus the checks themselves. Each criterion records: its ID, the
requirement(s) it verifies, its level (`spec` or `code`), its method, how to run it, and its
status (`verified`, `partial`, `gap`, `waived` — a waiver names who waived it and why).

`acceptance.md` also names the project's **standing checks**: the small set every change must pass
regardless of what it touched (typically build, format, lint/type-check, and the test suite under
the strictest analysis available).

**→ GATE 2: the owner reviews and approves `acceptance.md` before implementation begins.**

## Phase 3 — spec2code

**Goal:** code that satisfies the requirements and passes the criteria.

**The invariant: spec2code never resolves a conflict.** If implementation reaches a subject with
more than one active value, it stops and re-enters context2spec for that subject. A code generator
that picks one — whichever requirement it read last, or whichever it found most convincing — has
made a product decision silently, left no record, and produced an implementation that will
contradict the spec forever. This is the single most important rule in this document. Discovering
a conflict during implementation is normal and fine; deciding it there is not.

**Method:**

1. Write or change the code to satisfy the target requirements.
2. Run the **full** acceptance suite, not just the new checks. A change that satisfies `P-020` by
   breaking `P-003` is not done.
3. Run the non-automated criteria that apply — agentic review, human review — rather than quietly
   skipping the inconvenient ones.
4. Update `arch.md` to describe the resulting structure, and any user-facing docs whose subject
   matter changed.
5. Report which criteria ran and which didn't. **Never claim a criterion passed without running
   it.** A false green is the most expensive output this process can produce.

**Outputs:** code, checks, an updated `arch.md`.

`arch.md` describes the code as it is: data flow, module structure, the seams that make it
testable, where state lives, and the known rough edges. It is written for someone about to change
the code, so it should be honest about what is untested or fragile.

## The loop for ongoing changes

Not every change needs all three phases. Classify first:

| Change | Process |
| --- | --- |
| New feature or changed behavior | Full loop: spec edit, compose, resolve, then criteria, then code. |
| Bug fix | The spec already forbids it. Add a criterion that catches it under the existing requirement, then fix. If no requirement covers it, that's a spec gap — add the requirement first. |
| Refactor | No spec change. The acceptance suite must stay green; if the work requires editing a check's *assertions*, it isn't a refactor. |
| **Adopting a new source** | `sources.md` entry, then a full compose-and-resolve pass. Expect conflicts; budget for them rather than discovering them under deadline. |
| **Refreshing a source pin** | Re-fetch, diff against the vendored copy, re-resolve any conflict the diff touches, re-run spec-level and affected code-level criteria. Never bump a pin without reading the diff. |
| Dependency or toolchain bump | No spec change unless it moves a supported-version requirement. Standing checks must pass. |
| Docs only | No spec change, unless writing the doc reveals the spec was wrong. |

**Never edit a requirement to match code that violates it** without saying so explicitly and
getting agreement. That converts a bug into a feature by fiat, and it is the main way this kind of
process rots.

## Drift control

Specs decay, and composed specs decay from two directions at once. Four defenses:

- **Same-commit rule.** A change that touches behavior touches the spec in the same commit.
- **Periodic reconciliation.** Re-run context2spec against the current code and diff the result
  against the normative files. Each disagreement is a defect in the code or in the spec; decide
  which, and record the decision.
- **Upstream review.** On a schedule, check each pinned source for a newer version. Staying behind
  is a legitimate choice; not knowing you're behind is not. For mandatory sources this is a
  compliance obligation, not hygiene.
- **Orphan and gap sweeps.** Search for cited IDs that no longer exist, requirements no check
  references, criteria whose status has been `gap` long enough to count as a decision, `Deferred`
  items past their date, and expired risk acceptances.

## Human gates

| Gate | When | What's approved |
| --- | --- | --- |
| Gate 1 | End of context2spec | `product.md`, `tech.md`, `sources.md`; no unresolved conflicts |
| Gate 2 | End of spec2test | `acceptance.md` and the checks it names |

An agent stops at each gate and asks. Approving requirements is cheap; discovering after
implementation that they were the wrong requirements is not.

Gates are about *direction*, not correctness review — the owner is confirming that these are the
right things to build and the right way to check them, not proofreading. The one thing the owner
must actually read line by line is the list of requirements that lost a conflict, because that is
where their intent was overruled.

Conflict escalations are a third, unscheduled kind of gate: they can arise in any phase, and they
always return to a human.

## Project bindings

This file is project-independent. Each project supplies the following, recorded in `tech.md`
(under a "Project bindings" heading) rather than here, so that this file stays identical across
projects:

| Binding | Example of what to record |
| --- | --- |
| Standing checks | The exact commands every change must pass, in order |
| Comment syntax for ID citations | The language's line-comment form, preserving the `Verifies:` token |
| Artifact locations | Where the bundle lives, if not `spec/`; any file deliberately left elsewhere |
| Subject vocabulary | The dotted-path namespace, and who may add to it |
| Hard safety constraints | What must never happen, even in a test or a scratch script |
| External systems | Which real services must never be contacted from tests, and the approved fakes |
| Approval authority | Who can pass a gate, waive a criterion, amend a requirement, grant a source deviation, or accept a risk |
| Source refresh cadence | How often pinned sources are checked for updates |
| Existing docs | Which project docs own which facts, so the no-duplication rule is enforceable |

Where a project already has an agent instruction file or contributor guide, its rules outrank this
process; this document never grants an exception to a project's own safety constraints.

## Failure modes to watch for

- **Spec theater.** Requirements written to look thorough, verified by criteria that can't fail.
- **Retro-fitted requirements.** Amending the spec to describe whatever the code does, which makes
  the process a rubber stamp.
- **Silent override.** A governing source quietly voids a local requirement and nobody tells its
  author, who goes on believing it holds.
- **Tier laundering.** Declaring an inconvenient source "advisory," or claiming a local
  requirement is "more specific," to escape an override that should have been negotiated.
- **Adoption by gesture.** "We follow the accessibility standard" with no conformance level, no
  pin, and no criteria — a compliance claim with nothing behind it.
- **Inference dressed as evidence.** A causal story cited in a decision record without being
  checked, when the repository's own history would have settled it. It survives because it
  sounds right and nobody re-derives it.
- **Stale pins.** Vendored sources years behind upstream, discovered during an audit.
- **Conflict deferral as a habit.** `Deferred` entries accumulating past their dates until the
  register is a to-do list nobody reads.
- **Gate skipping under deadline.** The gates cost minutes; being wrong about the requirements
  costs the whole implementation.
- **Doc sprawl.** The same fact in four files, three of them stale.
- **Over-specification.** Requirements that pin down implementation detail, converting every
  future refactor into a spec negotiation.
- **Unfalsifiable coverage claims.** A traceability table where every row says "verified" and
  nobody has checked that the checks bite.
- **Characterization sold as specification.** A check whose expected value was copied from a
  run, so it can only ever confirm that the code still does what it did. It will even fail
  honestly when the code changes — which is what makes it look like a real check.

## Revising this process

Change this file when the process changes, in the same commit as the change. Record what changed
and why below, so an improvement can be told apart from a drift. When a project's experience
suggests an improvement, make it in the canonical copy and propagate it, rather than letting each
project's copy diverge.

| Date | Change |
| --- | --- |
| 2026-08-09 | Initial version. Three phases; `P-###`/`T-###`/`AC-###` IDs cited in code comments; criteria in a separate `acceptance.md`; human gates after context2spec and spec2test; per-project details isolated as bindings. |
| 2026-08-09 | Added composition and conflict resolution: source provenance and qualified IDs; four precedence tiers; `sources.md` manifest with pinned, vendored imports; `decisions.md` conflict-resolution log; optional subject/value structure on value-setting requirements; spec-level acceptance criteria for emergent conflicts; the rule that spec2code never resolves a conflict; blocking treatment of unsatisfiable mandatory requirements. |
| 2026-08-09 | Required active content to be stripped from vendored documents, after browser-saved pages committed a publisher's API key to a public repo; and required a check that quoted passages still appear in the stored copy, which immediately found two transcriptions that had tidied the text they quoted. |
| 2026-08-09 | Documented what a vendored source directory contains, and dropped the version from its name: dating the directory turned a refresh into an add-plus-delete and hid the upstream diff that vendoring exists to produce. |
| 2026-08-09 | Required phase 1 to record which inputs it actually consulted, after a retrofit skipped the issue tracker: it missed an open bug report seven months old, and mistook a workaround for an optimization by trusting a commit message over the issue that prompted it. *(Amended 2026-08-10: this entry originally also claimed the retrofit had rediscovered a conflict already diagnosed in that tracker. It had not — the comment cited was written during the work, not before it. See CR-003.)* |
| 2026-08-09 | Added guidance for a source that cannot be fetched: a human retrieves it, the block is not worked around, and the requirement is left unwritten rather than filled in from memory. Prompted by iNaturalist answering automated requests for its own API guidance with a bot challenge. |
| 2026-08-09 | Warned against deriving a check's expected value from a run rather than from the requirement, after a two-month-old defect survived a test whose allowlist had been built from the buggy output. Distinct from "a criterion must be able to fail": that test would have failed if the code broke, but could never report that the code was already wrong. |
| 2026-08-09 | Required decision records to separate checked evidence from inference, after a retrofit cited a cleanup tool's existence as support for a bug the tool predated by five months. Version control would have settled it in one command. |
| 2026-08-09 | Collected the bundle under `spec/`, leaving the repository root for entry-point docs. Motivated by vendored sources needing a directory regardless, by CODEOWNERS on one directory being a mechanical enforcement of Gate 1, and by making `process.md` sit at an identical path in every project so it can be propagated. |
