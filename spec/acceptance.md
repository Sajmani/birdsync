# Acceptance criteria — birdsync

How we know each requirement holds. Written per
[process.md](process.md#phase-2--spec2test).

**Status: approved at Gate 2 on 2026-08-09.**

Two of the questions below were deferred rather than answered, and are carried as open
items rather than settled silently: how much of `main` is worth testing, and whether the
[recommended additions](#recommended-additions) land now or as they become relevant.
Neither blocks the phase-3 fixes. The third, AC-024's fidelity, was settled by the
evidence in [CR-003](decisions.md).

Each criterion names its level (`spec` checks the requirement set, `code` checks the
implementation), its method, and the exact command to run it. Statuses are `verified`
(the check exists and has been watched failing), `partial`, `gap` (no check), and
`pending` (specified, but lands with a phase-3 fix).

Tests cite the requirements they cover in a `// Verifies:` comment, so the trace from a
requirement to its check is:

```
grep -rn 'Verifies:.*P-036' --include='*.go' .
```

## Standing checks

Every change must pass all four, regardless of what it touched.

| ID | Command | Method | Covers |
| --- | --- | --- | --- |
| AC-001 | `go build ./...` | Compilation | T-001, T-004 |
| AC-002 | `go vet ./...` | Static analysis | T-025 |
| AC-003 | `gofmt -l .` (must print nothing) | Static analysis | T-024 |
| AC-004 | `go test -race ./...` | Unit + integration + dynamic analysis | Everything below |

## Criteria

All are level `code` unless stated otherwise. "Command" is the `-run` pattern under
`go test`.

| ID | Check | Method | Verifies | Status |
| --- | --- | --- | --- | --- |
| AC-005 | `TestNoLiveServiceHostsInTests` | Static analysis over the repo's own source | T-010 | verified |
| AC-006 | `TestDryRunIssuesNoWrites` | Integration, recording fake | P-051, T-005 | verified |
| AC-007 | `TestCreatedObservationContent` | Integration, recording fake | P-019, P-035–P-040, P-042, P-045, P-046 | verified |
| AC-008 | `TestClient_UpdateObservation` | Integration, `httptest` | T-009, P-046 | verified |
| AC-009 | `TestClient_CreateObservation` | Integration, `httptest` | T-016 | verified |
| AC-010 | `TestDownloadObservations` | Integration, `httptest` | T-017, P-023 | verified |
| AC-011 | `TestDownloadObservationsDateWindow`, `…NoDateWindow` | Integration, `httptest` | P-025 | verified |
| AC-012 | `TestBirdsync` | Integration, fakes | P-020, P-026, P-027, P-028, P-030, P-031 | verified |
| AC-013 | `TestUpdateMedia` | Integration, fakes | P-047 | verified |
| AC-014 | `TestFuzzyMatchDateFormats` | Integration, table-driven | P-033, T-019 | verified |
| AC-015 | `TestFuzzyMatchIgnoresEmptyNames` | Integration, fakes | P-032, T-020 | verified |
| AC-016 | `TestDryRunMediaCount` | Integration, fakes | P-053 | partial — see [below](#criteria-that-do-not-bite) |
| AC-017 | `TestMediaChange` | Unit, table-driven | P-048, P-049 | verified |
| AC-018 | `TestDownloadMLAsset_Photo`, `_Sound` | Integration, `httptest` | P-043, P-044 | verified |
| AC-019 | `TestRecord_Observed` | Unit, table-driven | T-019 | verified |
| AC-020 | `TestRecords` | Unit, temp file | T-018 | verified |
| AC-021 | `TestObservationID_Valid` | Unit, table-driven | P-019, P-022 | verified |
| AC-022 | `.github/workflows/ci.yml` runs AC-001–AC-004 on push and PR | CI configuration | T-030, T-031 | verified |
| AC-028 | `TestToolsAreReadOnly` | Static analysis (AST) over `tools/` | T-032 | verified |

### Criteria that do not bite

**AC-016 does not catch a broken `--dryrun` gate.** Mutating `birdsync.go:350` from
`if dryRun` to `if false` leaves `TestDryRunMediaCount` green, because it only inspects
counters — and the counters increment outside the gate ([CR-001](decisions.md)). AC-006
was written for this reason and does fail under the same mutation. AC-016 is kept for
what it does check (P-053, the unclassified media count), not as a dry-run guarantee.

This is the concrete case for the process's rule that a criterion must be watched
failing. Three checks were mutation-tested while writing this document:

| Criterion | Mutation | Result |
| --- | --- | --- |
| AC-006 | `if dryRun` → `if false` at birdsync.go:350 | fails, as intended |
| AC-016 | same mutation | **passes** — does not bite |
| AC-008 | `IgnorePhotos: true` → `false` | fails, as intended |
| AC-005 | added a scratch `_test.go` naming the live API host | fails, as intended |

### Known limits

- **AC-005 is a substring scan.** It catches a real hostname pasted into a test, not a
  hostname assembled from parts or read from the environment. Parsing every expression
  that could produce a host would cost more than it catches; the residual risk is that a
  determined mistake gets through, and human review is the backstop.
- **AC-007 asserts against a fake, not iNaturalist.** It pins down what birdsync *sends*.
  Whether iNaturalist interprets those fields as intended is unverifiable here by T-010,
  and is covered only by the maintainer's manual testing against a real account.

## Pending criteria

Specified now, landing with the phase-3 fixes in
[decisions.md](decisions.md#work-arising). Each must be watched failing against today's
code before its fix is written; none can be committed before then, because a red suite
would break AC-004 for everyone.

| ID | Check | Verifies | From |
| --- | --- | --- | --- |
| AC-023 | A `--dryrun` run's summary says `Would create` / `Would update`, never `Created` / `Updated` | P-060, T-007 | CR-001 |
| ~~AC-024~~ | *landed* — `TestDownloadObservationsNoTaxonFilter`, `TestUntaxonedObservationIsRecognized` | P-061, P-020 | CR-003 |
| AC-025 | A record with an unparseable date, time, latitude, or longitude is skipped and counted, and the run completes and processes later records | P-062 | CR-005 |
| AC-026 | `downloadMLAsset` leaves no file behind after the caller finishes with it | T-023 | Gate 1 |
| AC-027 | Static analysis: no `log.Fatal` outside `main` and `tools/` | T-027 | CR-002 |

AC-024 is worth two checks rather than one. Asserting the absent parameter tests the
mechanism; the empty-taxon round trip tests the requirement, and would survive a future
change of mechanism.

**AC-024 has a prerequisite that no automated check can supply.** Whether
`iconic_taxa[]=Aves` actually drops untaxoned observations is a fact about iNaturalist,
not about birdsync, and T-010 forbids asking it from the test suite. A read-only
`taxonfilter` tool was written to answer it empirically; its run is recorded in
[decisions.md](decisions.md) CR-003, and the tool was deleted afterwards (recoverable at
commit `fb581a3`).

The run was inconclusive — the account held no unidentified observations — so **AC-024 is
written against the requirement rather than against iNaturalist's filter semantics**: a
previously synced observation must be recognized whatever its taxon. That check is
independent of a fact we could not establish, which makes it the more durable one anyway.

AC-027 is a static check because the convention is unenforceable by example — a test can
only show that one function returns an error, never that no function calls `log.Fatal`.

## Traceability

Every requirement, and what verifies it. This table is the honest accounting: 33 of 93
requirements have an automated check.

| Requirement | Criteria | Status |
| --- | --- | --- |
| P-001 purpose | — | gap (narrative) |
| P-002 one-way sync | — | gap (architectural; human review) |
| P-003 CSV input | AC-020 | partial |
| P-004 no reverse sync | — | gap (non-goal; human review) |
| P-005 never modifies others' observations | — | **gap — see [Recommended additions](#recommended-additions)** |
| P-006 no taxonomy reconciliation | — | gap (non-goal) |
| P-007 `tools/` not part of the product, read-only | AC-001, AC-028 | verified |
| P-008 one positional argument | — | gap (`main`, untested) |
| P-009 flags precede the argument | — | gap (Go flag behavior) |
| P-010 `--flag=false` form | — | gap (Go flag behavior) |
| P-011 unopenable CSV exits early | — | gap (`main`, untested) |
| P-012 `--after` later than `--before` exits | — | gap (`main`, untested) |
| P-013 `--debug` explains skips | — | gap |
| P-014 credentials from environment | — | gap |
| P-015 interactive prompts | — | gap (interactive; human review) |
| P-016 token framing | — | gap |
| P-017 401 says refresh the token | — | gap |
| P-018 credentials never logged | — | **gap — security-relevant** |
| P-019 sync key | AC-007, AC-021 | verified |
| P-020 idempotence | AC-012, AC-024 | verified |
| P-021 taxon not part of the key | AC-021 | partial |
| P-022 incomplete key not recognized | AC-021 | verified |
| P-023 downloads existing observations | AC-010 | verified |
| P-024 *withdrawn (CR-003)* | — | n/a |
| P-025 date window narrows the download | AC-011 | verified |
| P-026 skip order | AC-012 | verified |
| P-027 `--after` | AC-012 | verified |
| P-028 `--before` | AC-012 | verified |
| P-029 flag date formats | — | gap (`dateTimeFlag.Set` untested) |
| P-030 `--verifiable` default on | AC-012 | verified |
| P-031 `--fuzzy` matching | AC-012, AC-014 | verified |
| P-032 empty names excluded | AC-015 | verified |
| P-033 parsed dates, not raw | AC-014 | verified |
| P-034 `--fuzzy` off by default, documented | — | gap (documentation; human review) |
| P-035 wild, not captive | AC-007 | verified |
| P-036 coordinates, inexact, accuracy | AC-007 | verified |
| P-037 species guess | AC-007 | verified |
| P-038 observed date/time | AC-007 | verified |
| P-039 observation-field mapping | AC-007 | verified |
| P-040 description contents | AC-007 | verified |
| P-041 observations are public | — | gap (property of iNaturalist) |
| P-042 create before attaching media | AC-007 | verified |
| P-043 2400px photos, MP3 sounds | AC-018 | verified |
| P-044 photo-then-sound fallback | AC-018 | verified |
| P-045 `ML<id>` filename | AC-007 | partial — asset ID checked, filename not |
| P-046 description updated, media kept | AC-007, AC-008 | verified |
| P-047 additive media re-sync | AC-013 | verified |
| P-048 removals reported, not applied | AC-017 | verified |
| P-049 count mismatch reported | AC-017 | verified |
| P-050 media failures tolerated | — | **gap — error paths untested** |
| P-051 dry run issues no writes | AC-006 | verified |
| P-052 `DRYRUN:` prefix | — | gap (log output unasserted) |
| P-053 unclassified media count | AC-016 | verified |
| P-054 end-of-run summary | — | gap (`main`, untested) |
| P-055 conditional counters | — | gap (`main`, untested) |
| P-056 failure count printed | — | gap (`main`, untested) |
| P-057 *withdrawn (CR-005)* | — | n/a |
| P-058 create/update failure aborts | — | gap |
| P-059 no rollback | — | gap (consequence of P-058) |
| P-060 dry-run labels | AC-023 | pending |
| P-061 no taxon filter | AC-024 | verified |
| P-062 bad records skipped and counted | AC-025 | pending |
| T-001 single module | AC-001 | verified |
| T-002 one dependency | — | **gap — see [Recommended additions](#recommended-additions)** |
| T-003 `go`/`toolchain` policy | — | gap (human review) |
| T-004 platform independence | AC-001 | partial (one OS in CI) |
| T-005 `--dryrun` gates at the call site | AC-006 | verified |
| T-006 `DRYRUN:` prefix | — | gap |
| T-007 honest counters | AC-023 | pending |
| T-008 client mutates unconditionally | — | gap (design note; human review) |
| T-009 `ignore_photos` always set | AC-008 | verified |
| T-010 no live API calls in tests | AC-005 | verified |
| T-011 *withdrawn* | — | n/a |
| T-032 `tools/` is read-only | AC-028 | verified |
| T-033 never run `tools/` | — | gap (process rule; human review) |
| T-012 credentials never logged | — | **gap — security-relevant** |
| T-013 client interfaces | AC-006, AC-007, AC-012 | verified (used, so preserved) |
| T-014 base-URL seams | AC-010, AC-018 | verified (used, so preserved) |
| T-015 `resetFlags` in every test | — | **gap — two tests violated this until now** |
| T-016 `User-Agent` header | AC-009 | verified |
| T-017 page size 200 | AC-010 | verified |
| T-018 CSV read by header name | AC-020 | verified |
| T-019 dates via `Observed()` | AC-014, AC-019 | verified |
| T-020 empty name excluded | AC-015 | verified |
| T-021 asset type unknown before download | AC-018 | verified |
| T-022 memory ceiling | — | gap (unmeasured; benchmark recommended) |
| T-023 temp files deleted | AC-026 | pending |
| T-024 `gofmt` | AC-003 | verified |
| T-025 `go vet` | AC-002 | verified |
| T-026 error wrapping | — | gap (human review) |
| T-027 `log.Fatal` placement | AC-027 | pending |
| T-028 `log.Printf` vs `debugf` | — | gap (human review) |
| T-029 comments explain why | — | gap (human review) |
| T-030 CI runs the standing checks | AC-022 | verified |
| T-031 CI early-warning job | AC-022 | verified |

## Gaps worth naming

Not all gaps are equal. These are the ones where the absence of a check is itself a risk,
rather than a requirement that simply isn't mechanically checkable:

1. **P-018 and T-012 — credentials never logged.** A security-relevant invariant with no
   check at all. A static scan for the token variable reaching a logging call is
   feasible, and cheaper than the incident.
2. **P-050 — media failures tolerated.** The error paths are entirely untested: the mock
   fields `createObsErr`, `updateObsErr`, and `uploadMediaErr` exist but no test sets
   them. A run that dies on the first failed upload would pass the whole suite.
3. **T-015 — `resetFlags`.** `TestBirdsync` and `TestUpdateMedia` did not call it.
   `TestUpdateMedia` used `after.Set("")`, which fails silently, so it was passing only
   because the previous test happened to leave a window containing its record's date.
   Fixed in this phase; the underlying hazard — order-dependent tests — has no check.
4. **`main` is untested.** P-008 through P-013 and P-054 through P-056 all live in
   `main`, which no test exercises. Flag parsing, the argument check, the
   `--after`/`--before` sanity check, and the entire summary block are unverified.
5. **`inat.Client.UploadMedia` has no test.** The multipart request it builds is the one
   place binary data is written to a user's account, and it is exercised only through a
   mock that ignores its arguments.

## Recommended additions

Cheap, and each closes a gap above. Not implemented pending Gate 2, because each is a
judgment call about how much checking this project wants:

| Proposal | Method | Would cover | Cost |
| --- | --- | --- | --- |
| Assert the token never appears in captured log output | Unit, `log.SetOutput` to a buffer | P-018, T-012 | small |
| Table-driven error-injection through the existing mock fields | Integration | P-050, P-058 | small |
| Extract `main`'s summary block into a testable function | Refactor + unit | P-054–P-056 | medium |
| Test `main`'s argument and flag validation via `os/exec` on the built binary | Integration | P-008, P-011, P-012 | medium |
| `UploadMedia` against an `httptest` server, asserting the multipart filename | Integration | P-045, and the untested surface | small |
| Count `go list -m all` and fail above one dependency | Static analysis | T-002 | small |
| Benchmark a large synthetic export | Benchmark | T-022 | medium |

## Open questions for Gate 2

1. **How much of `main` is worth testing?** It holds seven requirements and no tests.
   Extracting the summary block is easy; testing flag parsing means either exporting more
   or driving the built binary.
2. **Should the recommended additions land now or as they become relevant?** Doing all
   seven roughly doubles the test suite.
3. **AC-024's second half needs a decision on fidelity.** Testing that an empty-taxon
   observation round-trips requires the fake to model iNaturalist's filter semantics — the
   very thing we could not confirm in CR-003. A fake that encodes a guess proves only that
   birdsync matches the guess.
