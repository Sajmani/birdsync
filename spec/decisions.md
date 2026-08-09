# Decision log — birdsync

Append-only record of conflicts found while composing the requirements, and how each was
resolved. Written per [process.md](process.md#conflicts-and-resolution).

Entries are never rewritten. A decision that turns out badly is superseded by a later
entry, not edited. Line references in an entry describe the code as it stood on that
entry's date, and are not updated as the file changes — that is what makes them evidence.

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

**Corrected 2026-08-09.** Two claims in this entry were wrong, both because the retrofit never
read the issue tracker — an input [process.md](process.md#gather) names and this analysis
skipped.

1. **The conflict was already known.** [Issue #5](https://github.com/Sajmani/birdsync/issues/5)
   carries the maintainer's own statement of it, from January 2026: *"I tried restricting to
   Aves, but this introduces a subtle bug: if someone changes the taxon for a birdsync-created
   observation to non-Aves, then subsequent birdsync runs won't see it and will just recreate
   it."* The composition analysis below rederived, at considerable cost, a bug that was
   already written down in public.
2. **The filter was not a download-volume optimization.** It came from PR #4, submitted to
   work around issue #5 — a hard failure past 10,000 results. The commit message describing it
   as reducing downloads records the mechanism, not the motive. The cost/benefit below
   therefore understated what removing it gives up: not two requests, but the only mitigation
   that existed for P-065.


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

**Instrument:** a read-only `taxonfilter` tool was written to settle this empirically. It
downloaded the account twice — once with `iconic_taxa[]=Aves`, once without — compared the
two by UUID, and reported how many of the hidden observations carried a birdsync sync key.
It was deleted after producing the evidence below; recover it from commit `fb581a3` if the
question is reopened.

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

## CR-006 — Uploaded file extensions varied by machine

- **Kind:** requirement under-specified, and violated by the implementation
- **Subject:** `media.filename.extension`
- **Involves:** P-045, T-004
- **Found:** 2026-08-09, by CI failing on Linux for the first time

`downloadMLAsset` derived the extension with
`mime.ExtensionsByType(contentType)` and took `extensions[0]`. That call returns
*every* extension registered for the type, sorted, from the machine's mime database, so
the first element is arbitrary:

| | photos (image/jpeg) | sounds (audio/mpeg) |
| --- | --- | --- |
| macOS | `.jpe` | `.m2a` |
| Linux (CI) | `.jfif` | `.m2a` |
| Before `0c7b6d2` | `.jpe` | `.mp3` |

Commit `0c7b6d2` ("Fix crash when uploading sound files") replaced a hardcoded
`ext := ".mp3"` with this lookup, so since June 2026 every sound has been uploaded as
`ML<id>.m2a`. Whether iNaturalist rejects that is not verifiable here (T-010), but the
README already warns users about unexplained media-upload failures.

P-045 said "the extension implied by the response content type", which is ambiguous: several
extensions are implied and it named no way to choose. T-004 is also implicated, since this is
platform-dependent behavior visible to the user.

The test did not catch it because it was written to accept whatever the platform produced —
its allowlist was `{.jpeg, .jpe, .jpg}` — so it passed on macOS while the code picked `.jpe`.
A criterion written around the observed behavior cannot detect that the behavior is wrong.

| Option | Effect |
| --- | --- |
| **A. Map known content types to a canonical extension**, fall back to the mime database | Deterministic on every platform; restores `.mp3`; still works if the CDN serves something new |
| B. Derive from which URL responded (`/2400` → `.jpg`, `/mp3` → `.mp3`) | Simplest, but assumes the photo endpoint is always JPEG and ignores the header |
| C. Prefer a canonical extension if the mime list contains one | Less code, still platform-dependent whenever no canonical form is present |

**Resolved 2026-08-09 (owner): option A.** P-045 amended to name the mapping. The download
tests now assert an exact extension, and `TestFileExtension` covers the mapping directly.

**Confirmed against the live service 2026-08-09.** A narrow real run (`--after 2026-08-05
--before 2026-08-06`) uploaded 8 assets, all logged as `ML<id>.jpg`, all accepted, and all 8
attached: each of the 4 observations shows a media count equal to the number of assets its
description lists. No temp files were left behind (T-023), and no errors occurred.

**Still unverified: the sound path.** All 8 assets were photos, so `audio/mpeg` → `.mp3` has
not been exercised against iNaturalist — and that is the half more likely to have been
rejected as `.m2a`. It cannot easily be tested on purpose, because an asset ID doesn't reveal
whether it is a photo or a sound until it is downloaded (P-044, T-021), so a checklist with
audio can't be picked out of the CSV in advance. Accepted as low residual risk: the mapping is
unit-tested and the photo half is confirmed.

Note the fallback case cannot assert a specific extension — `text/plain` resolves to `.conf`
on macOS — so it asserts only that some extension is returned. Pinning a value there would
rebuild the fragility the mapping removes.

## CR-007 — A failed media upload leaves its URL in the description

- **Kind:** requirement violated by implementation, with an emergent second-order effect
- **Subject:** `observation.description.assets`
- **Involves:** P-040, P-046, P-047, P-050
- **Found:** 2026-08-09, in the output of a `--dryrun` against the maintainer's real account

`addMedia` appends the asset URL to the description at the top of the loop body, before the
download and upload are attempted, and never removes it when either fails:

```go
for _, id := range assetIDs.ids {
    obs.Description += "Macaulay Library Asset: " + mlAssetURL(id) + "\n"   // unconditional
    filename, isPhoto, err := ebirdClient.DownloadMLAsset(id)
    if err != nil { log; s.errors++; continue }                              // already written
    err = inatClient.UploadMedia(...)
    if err != nil { log; s.errors++; continue }                              // still written
}
inatClient.UpdateObservation(obs)
```

Two consequences, both demonstrable from the source:

1. **The description misreports what is attached.** P-040 requires one line "per uploaded
   asset"; the code writes one per *attempted* asset. A direct violation.
2. **The failure becomes permanent.** On the next run `iNatMLAssets` parses the asset ID back
   out of the description, so `mediaChange` computes no difference and never retries the
   upload. P-047's additive re-sync — the mechanism that exists to catch up on missing media —
   is defeated by the record of the failure itself. Neither requirement is wrong alone.

The observation stays unverifiable ("Casual") until the user notices and fixes it by hand,
which is the outcome the README warns about without explaining the cause.

**Evidence:** the dry run reported `iNat description lists 1 ML Asset IDs, but observation has
0 media files` for one observation. The mechanism above is certain; that this specific
observation arose from it is inference — a photo deleted in iNaturalist by hand would look the
same. P-049's mismatch report is what surfaced it, so the detection works even though the
repair doesn't.

| Option | Effect |
| --- | --- |
| **A. Build the description from assets that actually uploaded** | Satisfies P-040; a failed asset stays absent from the description, so P-047 retries it next run |
| B. Keep writing all URLs, but re-attempt anything the observation is missing | Needs the attached-media list, which the update path doesn't currently fetch |
| C. Leave the description alone and report louder | No lie in the data, but the failure is still permanent |

**Recommendation: A.** It is the smaller change and it turns a permanent failure into a
transient one: the next run sees the asset missing from the description and uploads it.

**Origin, checked 2026-08-09.** The append has always been at the top of the loop, added in
`06c72c0`. It was unreachable as a defect until `0c7b6d2`, because a failed download or
upload called `log.Fatalf` and the process died before `UpdateObservation` ran — the
description could not be written with a failed asset in it. `0c7b6d2` ("Fix crash when
uploading sound files") replaced those calls with `log.Printf` and `continue`, which is the
right behavior, but left the append where it was. Making the error path survivable turned a
latent ordering bug into a live one.

That is the second defect from `0c7b6d2`; [CR-006](#cr-006--uploaded-file-extensions-varied-by-machine)
is the first. Both needed either a real failure or a non-macOS machine to become visible.

**Resolved 2026-08-09: option A.** `addMedia` collects the assets that uploaded and builds
the description from those. If none uploaded, the observation is not updated at all and
`updatedObservations` is not incremented — the previous code wrote back an unchanged
description and counted it (T-007).

Verified by `TestFailedUploadIsRetriedNextRun`, which runs the sync twice and feeds the first
run's own output back in as the second run's starting state, since the round trip through the
description is where the defect lived. Asserting only that the description omits the failed
asset would have tested the symptom.

## CR-008 — A permanently rejected asset was retried forever

- **Kind:** consequence of an earlier fix, surfaced by asking what happens next
- **Subject:** `media.upload.permanent_failure`
- **Involves:** P-040, P-047, P-050, P-063, T-034
- **Found:** 2026-08-09, by the maintainer asking what CR-007's fix does to a sound file
  over iNaturalist's 50 MB limit

[CR-007](#cr-007--a-failed-media-upload-leaves-its-url-in-the-description) stopped recording
failed uploads as done, so they are retried. For a transient failure that is the point. For a
permanent one — a file the service will never accept — it means re-downloading the asset from
the Macaulay Library and re-uploading it on every run, indefinitely, with no way to stop it.
birdsync keeps no state between runs except the iNaturalist observation itself, so anything it
must remember has to be recorded there.

CR-007 traded silent permanent loss for loud permanent waste. Both are wrong; the difference
is that the second one is visible.

Compounding it, the user could not tell the two cases apart: `roundTrip` discarded the
response body, so every refusal read `bad HTTP status: 422 Unprocessable Entity` regardless of
cause (T-034).

| Option | Effect |
| --- | --- |
| A. Retry forever, but make the error legible | Simplest; no new state; still re-downloads a large file every run |
| **B. Record permanent failures in the description** | Bounded and truthful; needs the description parsed into two sets so a recorded failure isn't counted as attached |
| C. Skip oversized files before uploading | Avoids the wasted upload, but hardcodes a threshold iNaturalist owns and that the README says is 50 MB only by observation |

**Resolved 2026-08-09 (owner): option B, plus the error-message fix from A.**

**Confirmed end to end against the live service 2026-08-09**, on the observation that
prompted CR-007. Deleting the recorded line from its description asked for a retry; birdsync
downloaded the 61.1 MB sound, iNaturalist answered `413 Request Entity Too Large`, the
failure was classified permanent, recorded in the description with the retry instructions,
counted once, and not retried. Every link in the chain behaved as specified.

That run also showed T-034 overshooting: the 413 came from nginx as a seven-line HTML page,
which went into the log verbatim. An HTML body is now dropped — it says nothing the status
line doesn't — and a genuine message is collapsed onto one line.

**P-064 confirmed on the following run**, which reported `1 ML Asset IDs previously failed to
upload and will not be retried: 637691397`, attempted no download or upload, updated nothing,
and counted no error. The asset is mentioned once per run and costs nothing.

Permanence is decided by the status code: a 4xx means the service rejected the request itself,
except 401 (refresh the token), 408 and 429 (try later). Anything else, including every 5xx and
every network error, stays transient and is retried.

The description is now parsed into two sets — uploaded and permanently failed. Both suppress a
retry; only the uploaded set is compared against the media actually attached, so a recorded
failure does not produce a P-049 count mismatch on every run.

Download failures are out of scope and stay transient: `ebird.DownloadMLAsset` returns an
untyped error, so a 404 for a withdrawn asset is retried like a network blip. Recorded as a
known limit rather than fixed here.

To ask for a retry, delete the `(upload failed)` line from the observation's description.

## CR-009 — No sound file could be downloaded at all

- **Kind:** requirement violated by implementation; a check that passed on invented data
- **Subject:** `media.filename.extension`
- **Involves:** P-043, P-045, T-004
- **Found:** 2026-08-09, while diagnosing the observation that surfaced CR-007

The Macaulay Library CDN serves sounds as `Content-Type: audio/mpeg3`, confirmed against it
for three independent assets (`637691397`, `623187742`, `633968116`); photos come as
`image/jpeg`. `audio/mpeg3` is not a registered media type, so
`mime.ExtensionsByType("audio/mpeg3")` returns an empty list, and `fileExtension` treated
that as an error. `DownloadMLAsset` therefore failed for every sound, before any upload was
attempted.

This dates from `0c7b6d2`, which replaced a hardcoded `ext := ".mp3"` for sounds — which
worked — with detection from the response header, which does not. It is the third defect
from that commit, after [CR-006](#cr-006--uploaded-file-extensions-varied-by-machine) and
[CR-007](#cr-007--a-failed-media-upload-leaves-its-url-in-the-description).

**CR-006's fix did not catch it, and neither did its test.** The canonical map was built from
the value in the existing test fixture, `audio/mpeg`, which was never checked against the
CDN. The test then passed because the fixture and the map agreed with each other. This is
precisely the failure mode added to [process.md](process.md#failure-modes-to-watch-for) hours
earlier — an expected value derived from something other than reality — and repeating it that
quickly says the warning is worth keeping.

**Resolved 2026-08-09:** the map is keyed on content types verified against the CDN, and
`fileExtension` can no longer fail: an unrecognised type falls back to what the endpoint
implies, since the photo URL serves images and the sound URL serves audio. Naming a file
slightly wrong is a far smaller harm than refusing to download it.

The system mime database is no longer consulted at all. It produced both extension bugs: a
platform-dependent answer for JPEG and no answer for MP3.

The test fixture now sends `audio/mpeg3`, and the fallback is observable — the tests assert
which types are mapped rather than only which extension comes out, because the fallback
returns `.mp3` for sounds too and would otherwise mask a missing map entry. Deleting
`audio/mpeg3` from the map now fails the suite; before that assertion it changed nothing any
test could see.

**Confirmed against the live service 2026-08-09.** A run that retried asset `637691397`
logged `Uploading sound as ML637691397.mp3` — the first sound birdsync has downloaded since
`0c7b6d2`. The upload then failed on size, which is a separate matter (below), but the
download and naming worked.

## CR-010 — birdsync cannot tell whose media it is downloading

- **Kind:** conflict between an adopted mandatory source and the implementation
- **Subject:** `media.download.ownership`
- **Involves:** `ml-terms/R2`, `ml-terms/R3`, `ml-terms/R4`, P-042, P-043
- **Found:** 2026-08-09, on vendoring `ml-terms`

`ml-terms/R2` permits a user to download their own Macaulay Library media — the archive is
described as their own backup — and `ml-terms/R1` confirms they keep copyright in it. That is
birdsync's intended case and it is squarely allowed.

`ml-terms/R3` says another contributor's media "may not be downloaded... without permission
from the author". birdsync has no way to honour that distinction: it fetches
`cdn.download.ams.birds.cornell.edu/api/v2/asset/<id>/2400` unauthenticated, and **the CDN
serves any asset to anyone** — verified 2026-08-09 with HEAD requests for two assets belonging
to other contributors, both HTTP 200. Nothing in the pipeline checks, or can check, that an
asset ID belongs to the person running the tool.

So compliance rests entirely on one unverified assumption: that `MyEBirdData.csv` never lists
an asset ID the user didn't upload. If that assumption fails, birdsync copies someone else's
copyrighted photo into the user's iNaturalist account under the user's name.

**Evidence for the assumption.** eBird states that "photos are associated with the eBird
account they are uploaded under" (`ml-terms/R4`), and the export is per-account. Four assets
birdsync has handled for the maintainer were checked against
`search.macaulaylibrary.org/api/v2/search?assetId=`, and all four report
`userDisplayName: "Sameer Ajmani"`; a control asset belonging to another contributor reports
their name, so the field discriminates.

**Evidence against it, or rather the gap.** A shared checklist displays both contributors'
media (`ml-terms/R4`). No page vendored here states whether the recipient's export lists the
other contributor's catalog numbers. A search summary claimed it does not; that claim was not
corroborated by any primary source and is not relied on here.

| Option | Effect |
| --- | --- |
| A. Confirm the export's behavior for shared checklists, then record the assumption as verified | Cheapest if the maintainer or eBird can answer; no code change |
| B. Check ownership before downloading, via the Macaulay search API's `userDisplayName` | Closes the gap mechanically; needs the user's contributor name, and depends on an undocumented API |
| C. Warn in the README and leave it to the user | Honest, but a user cannot audit thousands of asset IDs either |

**Recommendation: A, with B held in reserve.** The whole question disappears if a recipient's
export omits the sender's assets, and B adds a dependency on an unpublished endpoint to guard
against a case that may not exist. B is cheap to add later precisely because the API returns
the contributor.

**Evidence, 2026-08-09.** Tested against the maintainer's current export
(`MyEBirdData 8.csv`, 14,488 rows, 1,904 distinct assets) using
`search.macaulaylibrary.org/api/v2/search`, which reports each asset's `userDisplayName`:

| Sample | Own | Another contributor | Unindexed |
| --- | --- | --- | --- |
| 29 assets on multi-observer checklists (party sizes to 25) | 28 | **0** | 1 |
| 10 assets, one per checklist, from the set below | 1 | **0** | 9 |

A set-difference test — every asset in the export against every asset the API attributes to
the maintainer — left 74 unaccounted for, which looked like leakage. It was not: sampled
across 22 checklists, those assets return NOT FOUND from the search index while still
downloading from the CDN, so the enumeration of contributions was incomplete rather than the
export contaminated. The test is recorded as inconclusive; the sampling is what carries the
weight.

Neither sample found a single asset belonging to another contributor, including on the
checklists most likely to have been shared.

**Recommendation: close as an accepted assumption**, documented rather than enforced. The
sampling supports it, eBird's own statement that media belongs to the uploading account
(`ml-terms/R4`) explains why, and option B stays cheap to add if a counter-example appears —
the API returns the contributor, so a check is a few lines.

Note the residual: this is evidence from one account. A user whose export *does* contain a
co-observer's asset would have birdsync copy it without anyone noticing.

**Resolved 2026-08-09 (owner): closed as an accepted assumption**, documented rather than
enforced. birdsync relies on an eBird export listing only the user's own Macaulay Library
assets. The evidence above supports it and `ml-terms/R4` explains why it should hold.

Reopen on a counter-example: an asset in someone's export whose `userDisplayName` is not
theirs. Option B is then a few lines, since the API returns the contributor.

## CR-011 — The download cannot page past 10,000 results

- **Kind:** requirement violated by implementation; a governing source's limit
- **Subject:** `inat.download.paging`
- **Involves:** P-023, P-061, P-065, `inat-api/R3`
- **Found:** reported by a user as [issue #5](https://github.com/Sajmani/birdsync/issues/5) on
  2026-01-19; missed by the retrofit, which did not read the issue tracker

iNaturalist's recommended practices state it plainly:

> "The `page` and `per_page` parameters can be used to fetch up to (for many endpoints) 10k
> results. An error will be thrown if results beyond 10k are requested."

`inat.Client.DownloadObservations` pages with `page`/`per_page` and stops when it has
`totalResults`. For an account with more than 10,000 observations it never gets there: the
API fails — reported as HTTP 500 in the issue — and the run dies. birdsync is unusable for
those users, and has been since before the retrofit.

Removing the `Aves` filter under [CR-003](#cr-003--aves-only-download-can-defeat-duplicate-detection)
made this worse, and CR-003's analysis did not know it: the filter was PR #4's workaround for
exactly this failure, not the download-volume optimization its commit message described. So
the two conflicts are coupled, and the coupling was invisible while the issue tracker went
unread.

| Option | Effect |
| --- | --- |
| **A. Page with `id_above`**, sorting by id ascending | The method iNaturalist recommends for exactly this case; no ceiling; works whatever the taxon mix, including an account with 10,000 birds |
| B. Restore the `Aves` filter | Postpones the ceiling for users who record non-birds; reinstates CR-003's duplicate bug; no help to a birder with 10,000 bird observations |
| C. Restore the filter as `Aves` + `unknown` | Postpones the ceiling and covers the unresolvable-name case, since `unknown` is a documented `iconic_taxa` value; still fails a disputed ID that resolves to `Animalia`, and still has a ceiling |
| D. Narrow the download by date window | Already available via `--after`; a workaround the user must know to apply, not a fix |

**Recommendation: A.** It removes the ceiling rather than moving it, and it is what the
governing source tells clients to do:

> "One way to use the API and fetch more than 10k records is to sort by id ascending (e.g.
> `&order_by=id&order=asc`) and use the `id_above` parameter set to the ID of the record in
> the last batch."

C is worth keeping in reserve as a bandwidth optimization once A is in place, at which point
it is purely about request volume and can be judged on that alone. It should not be adopted
as a fix for this, because it leaves a ceiling and reintroduces a duplicate risk to buy back
two requests per run.

**Status: escalated — awaiting the owner.**

## Work arising

Phase-3 changes owed by the resolutions above. None may be implemented before
`acceptance.md` exists and passes Gate 2.

All five landed in phase 3 on 2026-08-09, each with a criterion watched failing first.

| From | Change | Requirements | Landed |
| --- | --- | --- | --- |
| CR-001 | Relabel the dry-run summary lines | P-060, T-007 | `53844b9` |
| CR-002 | Return errors from `ebird.Records` and `inat.DownloadObservations` | T-027 | this commit |
| CR-003 | Remove the `iconic_taxa[]=Aves` filter | P-061, P-020 | `084acba` |
| CR-005 | Skip, log, and count records with unparseable fields | P-062 | `4651b5e` |
| — | Delete temporary media files after upload | T-023 | `d5f4416` |

CR-002 landed as option A after all, not the recommended option B: with `tools/` reduced to
`dump`, changing `DownloadObservations`'s signature touched one caller instead of six, so
the reason for deferring half the fix had gone away.

## Deferred items

None.
