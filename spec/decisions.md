# Decision log — birdsync

Append-only record of conflicts found while composing the requirements, and how each was
resolved. Written per [process.md](process.md#conflicts-and-resolution).

Entries are never rewritten. A decision that turns out badly is superseded by a later
entry, not edited.

**Status: all conflicts resolved at Gate 1 on 2026-08-09 by the repository owner.** Four
of the five resolutions require code changes, which are phase-3 work and are listed in
[Work arising](#work-arising).

---

## CR-001 — A dry run reports observations it did not create

- **Kind:** requirement violated by implementation
- **Subject:** `stats.created`, `stats.updated`
- **Involves:** T-007 (a dry run's counters must not lie), P-051, P-054
- **Found:** retrofit, reading `birdsync.go:246` and `birdsync.go:362`

`s.updatedObservations++` and `s.createdObservations++` sit *outside* the `if dryRun`
branch, so a dry run ends with `Created 57 new iNaturalist observations` and
`Updated 57 iNaturalist observations` having created and updated nothing. The media
counter was carefully built to avoid exactly this (`stats.pendingMedia`, and the comment
above it), which suggests the create/update counters were an oversight rather than a
decision.

The behavior is arguably useful — a user does want to know how many observations a real
run would create — so the defect is in the *labels*, not the counting.

| Option | Effect |
| --- | --- |
| **A. Relabel under `--dryrun`** — print `Would create N` / `Would update N` | Satisfies T-007, keeps the useful number, small change, user-visible |
| B. Amend T-007 to cover media counters only | No code change; weakens the invariant that makes `--dryrun` trustworthy |
| C. Leave as is, document in the README | Cheapest; the summary still says something untrue |

**Resolved 2026-08-09 (owner): option A.** Under `--dryrun` the summary reads
`Would create N` / `Would update N`. Recorded as P-060; T-007 stands unamended. Code
change pending.

---

## CR-002 — `log.Fatal` in packages required to return errors

- **Kind:** requirement violated by implementation
- **Involves:** T-027
- **Found:** retrofit; already acknowledged in `CONTRIBUTING.md` and `arch.md`

`ebird.Records` (`ebird/ebird.go:83, 94, 97`) and `inat.DownloadObservations`
(`inat/inat.go:41, 65, 71`) call `log.Fatal` instead of returning to their caller. Both
files already warn readers not to copy the pattern, which is a documented deviation
rather than a resolved one — the process has no such category.

Note the asymmetry: `ebird.Records` already returns `(iter.Seq[Record], error)`, so
fixing it changes no signature. `inat.DownloadObservations` returns only `[]Result`, so
fixing it changes a signature that `tools/` also calls.

| Option | Effect |
| --- | --- |
| **A. Fix both** | Satisfies T-027; touches six call sites in `tools/`, which cannot be run to verify (T-011) |
| B. Fix `ebird.Records` now, defer the other | Removes the free half of the violation; leaves the deviation open with an owner and a date |
| C. Scope T-027 to new code | Honest about intent, but makes the convention unenforceable — a check can't tell new code from old |

**Resolved 2026-08-09 (owner): option A — fix both.** `inat.DownloadObservations` gains
an `error` return and its callers in `tools/` are updated. Since T-011 forbids running
those programs, they will be verified by compilation and review only, and that limit is
recorded here rather than discovered later. Code change pending.

---

## CR-003 — Aves-only download can defeat duplicate detection

- **Kind:** derived conflict between two requirements
- **Subject:** `sync.idempotence`
- **Involves:** P-020 (re-running creates no duplicates), P-024 (download is restricted
  to the `Aves` iconic taxon), P-006 (an unresolvable eBird name yields an unknown taxon)
- **Found:** retrofit, reading `inat/inat.go:47` against `birdsync.go:136-167`

No pair of these contradicts. Together they do:

1. birdsync creates an observation whose species guess is an eBird name iNaturalist
   cannot resolve — a "slash", a "spuh", a domestic form (P-006, P-021).
2. That observation has no taxon, therefore no iconic taxon.
3. The next run's download filters on `iconic_taxa[]=Aves`, so the observation is not
   returned.
4. It is therefore absent from `previouslySynced`, and the same eBird record is created
   again — violating P-020.

~~Circumstantial support: `tools/dedupe` exists specifically to delete observations sharing
a sync key, and the `Aves` filter was added later as a download-volume optimization
(commit `6f3b826`, "Only download Aves to reduce number of downloads") — after the
duplicate-detection design.~~

**Retracted 2026-08-09**, confirmed by the maintainer and by the commit history:

| Date | Commit | Event |
| --- | --- | --- |
| 2025-07-29 | `3d036ba` | eBird Checklist field (6033) introduced |
| 2025-07-30 | `7478a09` | eBird Scientific Name field added "to improve duplicate detection" |
| 2025-08-01 | `568cf08` | `dedupe` added, with `ebird.ObservationID.Valid` |
| 2026-01-18 | `6f3b826` | `Aves` download filter added |

`dedupe` was built to clear duplicates created *before* the birdsync observation fields
existed to detect them — it arrived two days after the sync key was completed, and five
months before the `Aves` filter. It is evidence about P-022, the same population
`tools/repair` serves, and says nothing about this conflict.

The claim should never have been written: the chronology was one `git log` away, and it
was asserted as "circumstantial support" without being checked. See
[process.md](process.md#conflicts-and-resolution) on separating evidence from inference,
added as a result.

What stands is the narrower point: the filter postdates the duplicate-detection design, so
it was not written with the sync key in mind.

**This has not been confirmed against the live API**, because T-010 forbids it. Confirming
it needs either the owner's knowledge of the API's filter semantics, or a test against a
fake that encodes the assumed behavior — which would only prove birdsync's logic, not
iNaturalist's.

**Instrument:** `tools/taxonfilter` was written to settle this empirically. It downloads
the account twice — once with `iconic_taxa[]=Aves`, once without — compares the two by
UUID, and reports how many of the observations the filter hides carry a birdsync sync key.
It is read-only. Run it with `go run ./tools/taxonfilter` and paste the verdict below.

**Evidence, 2026-08-09**, from a run against the maintainer's account (1478 observations,
983 of them birdsync-created):

```
Returned by iconic_taxa[]=Aves:  1026
  with no taxon at all:             0
  with a non-Aves taxon:          452
Present unfiltered but MISSING from the filtered query: 452
  of those, created by birdsync:    0
```

**This is inconclusive on the question that matters, not a refutation.** The account
contains no unidentified observations, so the filter was never handed one to drop. All 452
it dropped have a non-Aves iconic taxon — the filter working as designed. Step 2 of the
mechanism, whether the filter hides an observation with no iconic taxon, remains untested
because there was nothing to test it with.

Two corrections came out of this run:

- **The tool's first verdict overstated the result**, reporting "CR-003's mechanism is
  real" when the data showed only that non-bird observations are filtered out. Fixed: it
  now reports `INCONCLUSIVE` when the account has no unidentified observations, and
  distinguishes that from `REFUTED`.
- **The detector was wrong.** It classified on the taxon being absent, but iNaturalist may
  represent an unidentified observation as a placeholder such as "Life" or "Unknown"
  rather than as a null taxon. The tool now classifies on the *iconic* taxon being absent,
  which is the property the filter actually tests, and is correct under either
  representation. Both shapes are covered in `taxonfilter_test.go`.

**Why the decision stands regardless.** Option B is correct whether or not the mechanism
is real, and the measured cost of choosing it is smaller than assumed: 1478 observations
instead of 1026, which at 200 per page is **8 requests instead of 6**. The optimization
being protected is worth two HTTP requests per run. Against that, the downside of being
wrong is silent duplicate creation in a user's account. The asymmetry decides it without
needing the answer.

This also disposes of a refinement that the evidence briefly made attractive — asking for
`iconic_taxa[]=Aves` *and* unidentified observations in one query, keeping most of the
saving. It would only be safe if iNaturalist's parameter accepts that combination *and*
"unknown" captures the placeholder representation, and getting either wrong reintroduces
the bug silently. Two requests is not worth that risk.

**Residual uncertainty, accepted:** whether a birdsync observation can reach the
unidentified state at all is still unverified. The README says it can ("the observation
species name will say Unknown"), and eBird slashes and spuhs are the obvious candidates,
but this account has none — either they all resolved, or the maintainer fixed them by hand
as the README advises. Settling it definitively would require creating an observation with
an unresolvable name in a throwaway account. Not worth doing, given that the fix is going
in anyway.

| Option | Effect |
| --- | --- |
| A. Second query for observations with no iconic taxon, merged into the same index | Preserves the download-volume win; depends on the API supporting such a query |
| **B. Drop the `Aves` filter** | Certainly correct; reverts a deliberate optimization and downloads every observation the user has |
| C. Filter after download instead of in the query | Same volume cost as B, no benefit |
| D. Accept, and rely on `tools/dedupe` | Status quo; leaves users to discover duplicates and run a destructive tool with no dry-run flag |

**Resolved 2026-08-09 (owner): option B — confirmed as a real bug; drop the filter and
download all of the user's observations.** The owner rejected option A on the grounds
that a second, unfiltered query would re-fetch every `Aves` observation anyway, making
the first query pure overhead.

Consequences to carry into the amended requirements:

- P-024 is withdrawn and replaced by P-061.
- Download volume grows for users who record more than birds. This is the cost the
  `Aves` filter was introduced to avoid, knowingly re-accepted to preserve P-020.
- `--fuzzy` now matches against observations of any taxon, not just birds. The
  documented limitation that it "only compares against bird observations" becomes false
  once the code changes, so the README edit ships in the same commit as the fix.

Code change pending.

---

## CR-004 — `--fuzzy` is documented as matching non-birdsync observations

- **Kind:** documentation contradicts implementation
- **Involves:** P-031, P-022
- **Found:** retrofit, comparing `README.md:55` with `birdsync.go:147-166`

The README and the flag's help text both say `--fuzzy` compares against
"non-birdsync observations". The index is actually built from every observation whose
sync key is incomplete, which includes birdsync-created observations from versions that
set the checklist field but not the scientific name — the population `tools/repair`
exists to fix. For a user with such observations, `--fuzzy` suppresses records that a
correct reading of the documentation says it would sync.

`arch.md` already describes the code correctly, so the three documents disagree with each
other as well as with the code.

| Option | Effect |
| --- | --- |
| **A. Correct the README and the flag help** to say "an observation not carrying a complete birdsync sync key" | Documentation matches behavior; no code change |
| B. Change the code to exclude partially-keyed observations from the fuzzy index | Matches the documentation's promise, but those observations really are duplicates, so this would create them |

**Resolved 2026-08-09 (owner): option A.** The README and the flag help text are
corrected to describe what the code does. This is a documentation-only fix describing
*current* behavior, so it lands now rather than waiting for phase 3.

---

## CR-005 — Record-level parse failures abort the run

- **Kind:** undocumented requirement, surfaced by the retrofit
- **Involves:** P-057 (now withdrawn), P-050, P-058
- **Found:** retrofit, reading `birdsync.go:180` and `birdsync.go:299`

A malformed date, time, latitude, or longitude in any single CSV row calls `log.Fatalf`
and ends the run, potentially after thousands of observations have already been created.
A failed media transfer, by contrast, is logged, counted, and tolerated (P-050). Nothing
recorded why the two classes of failure are treated differently, which is the signature
of an accident rather than a decision.

| Option | Effect |
| --- | --- |
| **A. Skip and count the bad record** | Matches the media-failure policy; one bad row in a large export no longer ends the run |
| B. Keep aborting | A malformed export is a reason to stop and investigate |
| C. Validate the whole CSV up front, then abort before writing anything | All-or-nothing runs; a second full pass over the file |

**Resolved 2026-08-09 (owner): option A.** Recorded as P-062. The failure is logged with
its CSV line, counted, and reported in the summary when non-zero. P-057 is withdrawn.

Note this does **not** change P-058: a failed observation create or update still aborts
the run. Those failures indicate a broken token or a broken service, where continuing
would produce thousands of identical errors.

Code change pending.

---

## Work arising

Phase-3 changes owed by the resolutions above. None may be implemented before
`acceptance.md` exists and passes Gate 2.

| From | Change | Requirements |
| --- | --- | --- |
| CR-001 | Relabel the dry-run summary lines | P-060, T-007 |
| CR-002 | Return errors from `ebird.Records` and `inat.DownloadObservations`; update `tools/` callers | T-027 |
| CR-003 | Remove the `iconic_taxa[]=Aves` filter; update the README's fuzzy-matching limitation in the same commit | P-061, P-020 |
| CR-005 | Skip, log, and count records with unparseable fields | P-062 |
| — | Delete temporary media files after upload | T-023 |

## Deferred items

None.
