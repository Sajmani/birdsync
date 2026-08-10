<!--
Kept in the repo so the links stay correct as the files move. Links use heading
anchors rather than line numbers wherever the target is markdown, because line
numbers drift silently on every commit above them. TestTalkLinksResolve in
guard_test.go verifies every link in this directory against the working tree.
-->

# Spec composition and conflict resolution
## A 9-minute talk, presented live from the birdsync repo

**Repo:** <https://github.com/Sajmani/birdsync>
Links point at `main` and are checked by `TestTalkLinksResolve`. Timings are cumulative; the
last beat is the one to cut if you are running long.

---

## 0:00 – 1:00 · The problem

**Say:** Every project imports rules it did not write — brand guidelines, language style
guides, platform policies, law. Nobody writes those down as requirements, so they live in
people's heads, and you find out you violated one when a user account gets suspended.

Three things break when requirements come from more than one author:

1. **You don't know what applies.** Nobody enumerated the sources.
2. **You can't tell who wins.** Your requirement says one thing, the platform's says another.
3. **The conflicts aren't visible in any single document** — they only exist in the
   combination.

birdsync copies bird sightings from eBird into iNaturalist. Small tool, ~2000 lines. It turns
out to be governed by four external rule sets, and I found that out the hard way.

**Show:** <https://github.com/Sajmani/birdsync/blob/main/spec/process.md#sources-and-composition>
— *Sources and composition*, and the four precedence tiers.

> The tier is a property of the adoption, not of the source.

---

## 1:00 – 2:00 · (1) Discovering the sources

**Say:** The agent's first pass produced a manifest of what governs this project — with a
tier, an owner, and a scope for each. The interesting entries are the ones nobody would have
thought to ask about.

**Show:** <https://github.com/Sajmani/birdsync/blob/main/spec/sources.md>

Three points to land while it's on screen:

- **`ml-terms` came from a question, not a checklist.** birdsync downloads photos from the
  Macaulay Library and re-uploads them to iNaturalist. Is that allowed? Nobody had asked.
- **Rejected sources are recorded too** —
  [Considered and not adopted](https://github.com/Sajmani/birdsync/blob/main/spec/sources.md#considered-and-not-adopted).
  "We are not subject to GDPR, and here's why" is a decision that otherwise gets re-litigated
  annually.
- **Tiers do real work.** Go best practice is *advisory* — local rules beat it. iNaturalist's
  terms are *mandatory* — nothing local overrides them.

---

## 2:00 – 3:00 · (2) Extracting requirements

**Say:** Each source becomes numbered requirements that quote the original verbatim, then say
what it means *for this project*. The transcription is treated as an artifact that can be
wrong — it records who produced it, from what, on what date, and that it is not legal advice.

**Show:** <https://github.com/Sajmani/birdsync/blob/main/spec/sources/inat-terms/requirements.md#inat-termsr4--machine-generated-content-is-prohibited-hard-rule>
— `inat-terms/R4`, the one that mattered.

**Say:** Note the shape: the quote, then the consequence. That separation is what lets someone
check my reading against the source without trusting me.

**Optional, if time:**
<https://github.com/Sajmani/birdsync/blob/main/spec/sources/ml-terms/PROVENANCE.md>
— hashes, retrieval method, and an explicit "this is not legal advice and has not been
reviewed by a lawyer."

---

## 3:00 – 4:45 · (3) Discovering conflicts — two kinds

### 3:00 – 3:50 · The emergent kind

**Show:** <https://github.com/Sajmani/birdsync/blob/main/spec/decisions.md#cr-003--aves-only-download-can-defeat-duplicate-detection> — CR-003.

**Say:** Three requirements, no two of which contradict:

- Download only bird observations (added to work around an unrelated API limit).
- Never create a duplicate.
- An eBird name iNaturalist can't resolve produces an observation with **no taxon**.

Each is fine. Together: the untaxoned observation isn't returned by a bird-filtered query, so
birdsync can't see its own work and creates it again. **Silent, permanent duplicates.**

This is the argument for composition as a distinct step. No amount of careful reading of any
one requirement finds it.

### 3:50 – 4:45 · The source-versus-purpose kind

**Show:** <https://github.com/Sajmani/birdsync/blob/main/spec/decisions.md#cr-012--is-birdsync-machine-generated-content> — CR-012.

**Say:** iNaturalist's community guidelines, marked `(!)` — grounds for immediate suspension:

> "We do not allow machines to generate and post content on iNat with no human oversight
> curating each piece of content."

birdsync is a machine that posts observations. The penalty falls on **the user's account**,
not on the tool. A person could lose their observations for using software the README
recommends — and this had been shipping for a year without anyone asking.

---

## 4:45 – 6:00 · (4) Resolving

**Say:** Three resolution mechanisms, in order of cost.

**Mechanically, by tier.** Most conflicts never reach a human: the higher tier wins and the
loser is marked superseded. But the override is always *reported* — silently voiding
someone's requirement is the worst outcome, because its author still believes it holds.

**By escalation with options.** The agent computes; the human decides. Never fewer than two
options, each with what it satisfies and what it costs.

**By going and finding out.** CR-012 was resolved by fetching the definition page the rule
links to, which turned out to list *acceptable* examples — one of which is nearly a
description of birdsync.

**Show:** <https://github.com/Sajmani/birdsync/blob/main/spec/decisions.md#cr-012--is-birdsync-machine-generated-content>

> "Writing a script to create observations from a manually curated local folder of your
> images and metadata on your desktop."

**The rule that makes this work — say it explicitly:**
[spec2code never resolves a conflict](https://github.com/Sajmani/birdsync/blob/main/spec/process.md#phase-3--spec2code).
If code generation hits a contested value it stops and goes back. Otherwise the decision gets
made by whichever requirement the generator read last, and no record survives.

---

## 6:00 – 7:15 · (5) Criteria that stop regressions

**Say:** A conflict resolution is worthless if the next change quietly undoes it. Every
resolution ships with a check — and the check has to be *watched failing first*.

**Show:** <https://github.com/Sajmani/birdsync/blob/main/spec/acceptance.md#criteria-that-do-not-bite>
— *Criteria that do not bite.*

**Say:** This is the finding I'd most like people to take away. The tool's central safety
guarantee — `--dryrun` writes nothing — had a test. Mutating the code so `--dryrun` wrote to
the live account left that test **green**. It only ever checked counters. Every criterion in
this table was mutation-tested, and the table records the one that failed to fail.

**Then show:** <https://github.com/Sajmani/birdsync/blob/main/guard_test.go#L238-L262>
— `TestTranscribedQuotesAppearInSources`.

**Say:** Because the transcriptions are the highest-stakes artifacts here and an agent wrote
them, a test checks that all 29 quoted passages actually appear in the vendored documents. It
caught two places where the agent had "improved" the text it was quoting: swapped quotation
marks, and joined two bullet points with a full stop the source doesn't contain.

---

## 7:15 – 8:45 · (6) What went wrong with the process

**Say:** Six things, and they're the most useful part.

**1. The agent skipped an input its own method names.** Phase 1 lists "issues" as a source.
It never opened the issue tracker — where the maintainer had *already diagnosed* CR-003, in
public, months earlier. A day of analysis rediscovered a known bug.
→ <https://github.com/Sajmani/birdsync/blob/main/spec/decisions.md#cr-003--aves-only-download-can-defeat-duplicate-detection>

**2. Inference presented as evidence.** CR-003 cited a cleanup tool's existence as support
for a bug — a tool that predated the bug by five months. One `git log` would have settled it.
The retraction is in the same entry.

**3. Checks that pass for the wrong reason — repeatedly.** A content-type map built from a
test fixture nobody had verified against the real service. A paging probe with a threshold
below every record in the account. A test whose expected error text matched macOS's own
message. Same failure each time: the expected value came from something other than reality.

**4. Sources that refuse to be fetched.** Three iNaturalist pages return 403 with a
JavaScript challenge. The agent can't vendor them, shouldn't work around the block, and must
not fill the gap from memory — a plausible invented rate limit is indistinguishable from a
real one.

**5. Vendoring leaked someone else's secret.** Browser-saved pages carry inline scripts;
those carried iNaturalist's Google Maps API key straight into a public repo and tripped
secret scanning.

**6. The process changed nine times in one session.**
→ <https://github.com/Sajmani/birdsync/blob/main/spec/process.md#revising-this-process>
Every row is a real incident. That's the honest state of the methodology: it is being written
by being used, and it isn't finished.

---

## 8:45 – 9:00 · Close

**Say:** Twelve conflicts found and resolved. Five defects in one subsystem, three of them
from a single well-intentioned commit.

But be honest about attribution: the specs found CR-003 and CR-007. **CI found the extension
bug. A `curl` against the real CDN found the broken sound downloads. The issue tracker had
the paging limit all along.**

The composition process made things *findable and durable* — every decision has evidence, a
date, and a check. Running the thing against reality is what *found* them. You need both, and
the process should say so.

---

## Backup material, if asked

| Question | Where |
| --- | --- |
| What does a resolution record contain? | [decisions.md CR-011](https://github.com/Sajmani/birdsync/blob/main/spec/decisions.md#cr-011--the-download-cannot-page-past-10000-results) |
| How complete is the coverage? | [acceptance.md traceability](https://github.com/Sajmani/birdsync/blob/main/spec/acceptance.md#traceability) — ~60 of 100, gaps listed |
| What's still open? | [Open work](https://github.com/Sajmani/birdsync/blob/main/spec/decisions.md#open-work-as-of-2026-08-09) |
| Is this reusable? | [process.md](https://github.com/Sajmani/birdsync/blob/main/spec/process.md) is project-independent; per-project detail is confined to [bindings](https://github.com/Sajmani/birdsync/blob/main/spec/tech.md#project-bindings) |
| What did a conflict cost to resolve? | CR-012: two web fetches, one forum thread, one human with a browser |
| Why not just ask a lawyer? | The transcriptions say plainly they aren't legal advice; the point is knowing *which* questions need one |

## Cut list, in order

1. The optional PROVENANCE link at 2:00
2. The second conflict kind (3:50) — CR-003 alone carries the argument
3. Backup material
4. Limitation 5 (the leaked key) — the best story, but limitations 1–3 are the substantive ones
