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
[cut list](#cut-list-in-order) is at the end.

---

## 0:00 – 0:45 · The problem

**Say:** You write requirements for your product. But your product is also governed by rules
you did not write — brand guidelines, platform terms of service, API usage policies,
accessibility law. Nobody writes those down as requirements, so they live in whoever read them
last, and you find out you violated one when a user's account gets suspended.

Composing requirements from several authors raises three questions no single document answers:

1. **Which sources apply to us?**
2. **When they disagree, who wins?**
3. **What breaks only in the combination?**

That third one is the interesting one, and it's what the next ninety seconds is about.

---

## 0:45 – 2:15 · Why composition needs a process

**Slide:** `talks/sdd-example.png` — three stacked boxes, each a rendering of the same page.

**Say:** Suppose I'm building a web page. I write two product requirements: the background is
blue, and the text is light grey. That's the top box. It's legible — **7.41:1** contrast,
which passes WCAG AA *and* the stricter AAA. Good design, and I checked.

Now two things arrive that I didn't write.

**Brand guidelines** say the background must be the lighter brand blue. Perfectly reasonable
requirement. Fine on its own — that's the middle box, and notice **the text hasn't changed**.
Only the background has.

**Accessibility governance** says text must meet 4.5:1 against its background. Also fine on
its own. Also not negotiable.

Middle box: **1.56:1**. You can barely read it. And here's the thing — **no two of those
requirements contradict each other.** Read any pair and you'd approve it. The defect exists
only in the combination, which is why composition has to be a step someone performs rather
than something you hope falls out of code review.

**Resolution needs a rule for who yields.** That's what tiers are for:

| Tier | Example here | On conflict |
| --- | --- | --- |
| **Mandatory** | Accessibility | Always wins; nobody local may override it |
| **Governing** | Brand guidelines | Beats local requirements; the brand owner may grant an exception |
| **Local** | My two requirements | Mine to change |
| **Advisory** | House style guide | Applies wherever local is silent |

Work it through: accessibility cannot yield, brand outranks me, so the only thing that *may*
move is my text colour. Bottom box — dark grey on the brand blue, **5.48:1**. Passes.

**The punchline:** the tiers didn't just detect the conflict, they determined *which
requirement had to change*. And that dark grey on my *original* blue would have been
**1.15:1** — so no single edit fixes it either. The resolution only exists in the composition
too.

**Show:** [the four tiers, as specified](https://github.com/Sajmani/birdsync/blob/main/spec/process.md#sources-and-composition)

---

## 2:15 – 3:00 · birdsync, and formalising what it already did

**Say:** birdsync copies bird sightings from my eBird account into iNaturalist, with the
photos and sound recordings. About 2,000 lines of Go. **I wrote it by hand**, and I did the
homework — years ago, when I started, I read the API guidance, the terms of service, the rules
about media copyright and about automated posting. I made decisions based on them.

What I never did was write any of it down.

The first step wasn't to write new requirements. It was to **recover the ones already implicit**
in the code, the README, and the command-line help, and write them down with identifiers.

**Show:** [spec/product.md](https://github.com/Sajmani/birdsync/blob/main/spec/product.md)
— 68 product requirements, reverse-engineered from a tool that already worked.

**Say:** Each one is individually testable, and each says *why*. That "why" earns its keep
later: it's the difference between a deliberate constraint and an accident nobody remembers
making.

---

## 3:00 – 3:45 · Acceptance criteria, and the first thing they found

**Say:** Then every requirement gets a check that would fail if it were violated. And a check
you have not watched fail is not yet a check — so I broke the code deliberately, one
requirement at a time, to confirm each one goes red.

birdsync's most important promise is `--dryrun`: it reads, it never writes. That's the
guarantee you rely on before pointing a tool at your own data. It had a test.

I changed the code so `--dryrun` wrote to the live account. **The test stayed green.**

**Show:** [Criteria that do not bite](https://github.com/Sajmani/birdsync/blob/main/spec/acceptance.md#criteria-that-do-not-bite)

**Say:** It only ever inspected counters — and the counters were incremented outside the
safety gate. The fake clients discarded everything handed to them, so nothing *could* observe
whether a write happened. Fixing that meant recording the calls; only then was the guarantee
checkable at all. That's `AC-006`, and it's why I now believe `--dryrun` rather than hope.

---

## 3:45 – 4:45 · Discovering the sources

**Say:** Only then did we go looking for the rules I hadn't written. Four sources, each with a
tier, an owner, and a scope.

**Show:** [spec/sources.md](https://github.com/Sajmani/birdsync/blob/main/spec/sources.md)

- **eBird / Macaulay Library** — *mandatory*. birdsync downloads my photos from the Macaulay
  Library and re-uploads them to iNaturalist. Is that allowed? I had never asked it in those
  terms.
- **iNaturalist terms and community guidelines** — *mandatory*.
- **iNaturalist API recommended practices** — *governing*. Rate limits, paging, client
  identification.
- **Go practice** — *advisory*, so my local rules beat it.

**Say:** Two details worth pausing on. **Rejected sources are recorded too** —
[considered and not adopted](https://github.com/Sajmani/birdsync/blob/main/spec/sources.md#considered-and-not-adopted),
so "we're not subject to GDPR, and here's why" isn't re-litigated every year. And each source
is **vendored and hashed** — because until now there was no way to notice that a document I'd
read once had since changed.

---

## 4:45 – 6:15 · Resolving the sources against what I'd built

**Say:** Here is what that omission actually costs, in one file.

**Show:** [inat/inat.go, the comment at line 31](https://github.com/Sajmani/birdsync/blob/main/inat/inat.go#L31-L34)

**Say:** I wrote that in July 2025. It quotes iNaturalist's API recommended practices by URL,
and right below it the code does what that page asks about page size. So I read the page, I
agreed with it, I implemented it, and I left the citation for the next person.

**The rate limit is on the same page.** One request per second. There was no rate limiter in
the code until yesterday — and I want to be precise about why, because "he didn't read it" is
the wrong lesson.

I *had* read it. I'd reasoned that birdsync stays under a request per second on its own: it
fetches observations in large pages, and between each write it downloads a photo. Slow
operations, naturally spaced. That reasoning was correct — a real run measures about 0.2
requests per second.

But look at what kind of guarantee that is. It was never written down, so no reviewer could
check it. It's **conditional** — turn off `--verifiable` and there's no photo to download
between creates, and the spacing evaporates. And it's **load-bearing on the code being slow**,
so the first person to add concurrency deletes the guarantee without touching anything that
mentions it.

That's the failure mode, and it isn't carelessness. Reading a document and acting on it leaves
no record of *which parts you acted on*, *why you decided what you decided*, or *whether it's
still true*. The knowledge was real; it lived in one head.

So composition didn't discover that I should care about these rules. It did four things I had
never done:

**1. It turned conclusions into requirements that can be checked.** I had decided years ago
that re-uploading my own Macaulay photos was fine. The terms confirm it — a contributor keeps
copyright and may download their own media. But composing that against the code surfaced
something I hadn't considered: birdsync fetches by asset ID from an unauthenticated CDN that
will serve *anyone's* asset, so the whole thing rests on an assumption about what an eBird
export contains. I sampled 39 assets from the checklists most likely to be shared. None
belonged to anyone else. That assumption is now written down **with its evidence**, where
before it was a conclusion I'd reached and forgotten reaching.

**2. It nearly stopped the project.** The community guidelines, under a heading marked as
grounds for immediate suspension:

> "We do not allow machines to generate and post content on iNat with no human oversight
> curating each piece of content."

birdsync is a machine that posts observations. The penalty falls on **the user's account**,
not on mine.

**Show:** [CR-012](https://github.com/Sajmani/birdsync/blob/main/spec/decisions.md#cr-012--is-birdsync-machine-generated-content)

**Say:** I had thought about this one too, back at the start. What I couldn't produce on my
own was the answer's *durability*. We resolved it by fetching the definition page that rule
links to, which lists *acceptable* examples — one nearly a description of birdsync: "writing a
script to create observations from a manually curated local folder of your images and metadata
on your desktop." Permitted, and now permitted **on the record, with the text quoted**.

And it **added two requirements**, `P-067` and `P-068`, because the guidelines place an
obligation on the *user* — review what was created, answer identifiers — that I had understood
and never passed on to anyone using the tool.

**3. It gave an existing decision a second reason.** `--verifiable` skips observations with
no photo — a deliberate choice I made for data quality. A moderator's comment reframed it: an
observation with no media gives an identifier nothing to identify, so it is pure cost to other
people. Same default, now with both reasons attached, so the next person to find it
inconvenient knows what they'd be trading away.

**4. It produced two technical requirements I simply did not have.** `T-035`: about one
request per second. `T-036`: page with `id_above`, because page numbers stop working past
10,000 results.

That second one is the whole story in miniature, so it's worth the thirty seconds:

- **January.** A user hits a hard failure — more than 10,000 observations and the download
  dies. He files it, *and* sends a pull request that filters the download to birds only, which
  makes it fail later rather than never. He says so explicitly in the issue.
- **June.** I merge it. Reasonable: it's an improvement, from a contributor, for a real problem.
- **August.** Composition finds that the filter causes birdsync to lose track of its own
  observations and silently duplicate them.

And here's the part I want to land: **the rollback couldn't stand on its own.** Reverting the
filter would have handed that user his original bug back. Removing the mitigation obliged us
to remove the limit — so `id_above` went in, the ceiling disappeared, and the January report
was finally fixed. His reply on the issue was "Good catch."

A well-intentioned workaround, merged in good faith, quietly traded one failure for a subtler
one. That's the class of thing composition is *for*.

---

## 6:15 – 7:45 · Criteria, then code

**Say:** Those requirements generated criteria, and the criteria drove the code. Three
concrete changes, each traceable to a rule I had read once and half-remembered.

| New requirement | New criterion | Code change |
| --- | --- | --- |
| `T-035` pacing | `AC-035` — three requests must take at least two intervals | A minimum interval in `roundTrip`, the one place any request is made. Compliance no longer depends on the code staying slow — which is what makes it safe to parallelise later |
| `T-036` `id_above` paging | `AC-034` — no `page` parameter; the cursor must advance | Rewrote the download loop; birdsync now works for accounts over 10,000 observations |
| `P-067`, `P-068` user obligations | `AC-037` — human review | A new [README section](https://github.com/Sajmani/birdsync/blob/main/README.md#being-a-good-inaturalist-citizen) telling users to review what was created and to answer identifiers |

**Show:** [the traceability table](https://github.com/Sajmani/birdsync/blob/main/spec/acceptance.md#traceability)

**Say:** Every requirement and what verifies it — including the ones with nothing verifying
them, listed as gaps rather than quietly omitted.

**And the criteria caught the spec drifting too.** The transcriptions are the highest-stakes
artifacts here — legal text, transcribed by an agent — so a test checks that every quoted
passage still appears in the vendored document.

**Show:** [`TestTranscribedQuotesAppearInSources`](https://github.com/Sajmani/birdsync/blob/main/guard_test.go#L238-L262)

**Say:** All 29 passages. It immediately caught two places where the transcription had
*improved* the text it was quoting — swapped quotation marks, and two bullet points joined
with a full stop the source does not contain. Small infidelities, in exactly the documents
where fidelity is the entire point.

---

## 7:45 – 9:00 · What went wrong with the process

**Say:** Four honest ones.

**1. It skipped an input its own method names.** Phase 1 lists "issues" as a source. It never
opened the issue tracker — where a user had reported the paging limit seven months earlier,
and where the reason an existing filter existed was written down. It read the commit message
instead, and so recorded that filter as a bandwidth optimisation when it was a workaround for
that very bug. I had to point at the issue myself before any of it surfaced.

**2. Inference dressed as evidence — twice, in the same entry.** A conflict record cited a
cleanup tool's existence as proof of a bug; the tool predated the bug by five months, and one
`git log` would have settled it. Then the *correction* to that claim asserted I had diagnosed
the conflict myself "months earlier" — read off a comment whose timestamp nobody checked. The
comment was written during this work, hours after the analysis found the conflict. Same
failure, one level up, in the paragraph about that failure.

*(If you want one slide from this talk, this is the one: an agent that cites evidence it has
not verified will do it again inside the correction.)*

**3. Checks that pass for the wrong reason, repeatedly.** A content-type map built from a test
fixture nobody had checked against the real service. A paging probe with a threshold below
every record in the account. An error-message test that matched macOS's own wording and would
have failed on the platform that reported the bug. The same shape every time: the expected
value came from something other than reality.

**4. Vendoring published someone else's secret.** Pages saved from a browser carry inline
scripts, and those carried iNaturalist's Google Maps API key into a public repo, where secret
scanning found it.

**Show:** [the revision log](https://github.com/Sajmani/birdsync/blob/main/spec/process.md#revising-this-process)

**Say:** Ten entries, all from this work. Every one is a real incident that changed the
method. That's the honest status: it is being written by being used, and it is not finished.

**Close:** Twelve conflicts found and resolved, and a seven-month-old user-reported bug
closed as a side effect of resolving one of them properly.

But I want to end on the thing that surprised me. I had read these documents. I had made
defensible decisions from them. What I had no way to do was **prove which decisions I'd made,
show why, notice when the rules changed, or tell my users what the rules asked of them.** That
is what composition added — not the reading, the *residue* of the reading.

And be careful about credit: the specs found some of these; **CI found one, a `curl` against
the real service found another, and the issue tracker had a third all along.** Composition
made the decisions findable and durable. Running the thing against reality is what found them.
You need both.

---

## Backup material, if asked

| Question | Where |
| --- | --- |
| What does a resolution record contain? | [CR-011](https://github.com/Sajmani/birdsync/blob/main/spec/decisions.md#cr-011--the-download-cannot-page-past-10000-results) |
| What is still open? | [Open work](https://github.com/Sajmani/birdsync/blob/main/spec/decisions.md#open-work-as-of-2026-08-09) |
| How complete is coverage? | ~60 of 100 requirements have an automated check; gaps are [named](https://github.com/Sajmani/birdsync/blob/main/spec/acceptance.md#gaps-worth-naming) |
| Is the method reusable? | [process.md](https://github.com/Sajmani/birdsync/blob/main/spec/process.md) is project-independent; per-project detail is confined to [bindings](https://github.com/Sajmani/birdsync/blob/main/spec/tech.md#project-bindings) |
| Why not just ask a lawyer? | The transcriptions say plainly they are not legal advice. The value is knowing *which* questions need one |
| What did resolving CR-012 cost? | Two web fetches, one forum thread, one human with a browser |
| Where did the colour numbers come from? | Sampled from `talks/sdd-example.png` (`#EEEEEE` on `#0000FF`, then on `#9FC5E8`, then `#434343` on `#9FC5E8`) and computed with the WCAG relative-luminance formula. Recompute if the slide is edited |

## Cut list, in order

1. Backup material
2. Failure 4, the leaked key — the best story, but 1–3 are the substantive ones
3. Resolution point 3, the `--verifiable` reframing
4. The `TestTranscribedQuotesAppearInSources` beat at 7:15 — keep the traceability table
