# Technical requirements — birdsync

Constraints the implementation must satisfy, independent of what the tool does for a
user. Written per [process.md](process.md).

**Status: approved at Gate 1 on 2026-08-09.** Retrofitted from the code,
[AGENTS.md](../AGENTS.md), and [CONTRIBUTING.md](../CONTRIBUTING.md). Those two files now
keep only a terse summary of the hard safety rules, with citations back to the `T-###`
requirements here.

The five requirements that were unsatisfied at Gate 1 have all been implemented; see
[decisions.md](decisions.md#work-arising).

## Module and dependencies

**T-001** — birdsync is a single Go module, `github.com/Sajmani/birdsync`, with no build
system, no code generation, and no linter configuration beyond `gofmt` and `go vet`.

**T-002** — The module has exactly one external dependency,
`github.com/google/uuid`. Adding another requires the owner's approval.
Subject: `build.external_dependencies.count` · Value: `1`
*Rationale: the small footprint is deliberate; it keeps `go install` fast and the
supply-chain surface near zero.*

**T-003** — The `go` and `toolchain` directives are raised only as a deliberate,
standalone change, never as a side effect of other work.
Subject: `build.go_directive` · Value: `1.23.0`
*Rationale: users on the default `GOTOOLCHAIN=auto` silently download a newer toolchain,
but anyone pinned to `GOTOOLCHAIN=local` on an older Go is blocked outright.*

**T-004** — The tool must build and run on any platform the Go toolchain supports, using
only the standard library plus T-002. No platform-specific code.

## Safety invariants

These are the rules whose violation costs a user real data. They outrank convenience.

**T-005** — Every operation that creates, updates, deletes, or uploads is gated on
`--dryrun` **at the call site in `birdsync()`**, not inside the client.
*Rationale: `inat.Client` is shared with `tools/`, which has no such flag. The gate
belongs where the decision is made, not in a client that other callers reuse.*

**T-006** — A skipped mutation is logged with a `DRYRUN:` prefix (implements P-052).

**T-007** — A dry run's counters must not report work that did not happen. Where a dry
run cannot know something, it counts it as unknown rather than guessing.
*Status: **not yet satisfied** — the created and updated counters are labeled as though
the work happened. [CR-001](decisions.md#cr-001--a-dry-run-reports-observations-it-did-not-create) resolved this
in favour of relabeling; implements P-060.*

**T-008** — `inat.Client`'s methods mutate unconditionally and know nothing about
`--dryrun`. This is deliberate, not an oversight.

**T-009** — `UpdateObservation` always sets `ignore_photos`.
*Rationale: without it, updating a description detaches the observation's photos.*

**T-010** — Tests never contact `api.inaturalist.org`, the Macaulay Library CDN, or any
other live service, and never depend on a real account's contents.

**T-011** — ~~Nothing in the build, test, or development workflow runs a program in
`tools/`. They create, update, and delete observations in whatever account the
environment's credentials point at, and `purge` is checked in with its guard off.~~
Status: **Withdrawn** — the mutating tools were deleted rather than governed by a rule
nobody could check. Superseded by T-032 and T-033.

**T-032** — No program in `tools/` performs a mutating operation. Nothing there may
create, update, delete, or upload.
*Rationale: `tools/` held six programs that modified real accounts, guarded by a `debug`
constant that one of them shipped with turned off and by an instruction not to run them.
A structural property that a check can enforce is worth more than a warning people have to
read and remember. Deleted tools remain in the history if one is needed again.*

**T-033** — Nothing in the build, test, or development workflow runs a program in
`tools/`. They contact the live service with the user's real credentials, which no
automated process should do on the user's behalf.

**T-012** — Credentials are never logged, echoed, or written to a file.

## Testability seams

**T-013** — `birdsync()` takes its two clients as the `ebirdClient` and `inatClient`
interfaces, so the sync loop can be driven without a network.

**T-014** — `inat.NewClient` takes a base URL, and `ebird.DownloadMLAsset` delegates to
an unexported `downloadMLAsset(baseURL, id)`, so an `httptest` server can stand in for
either service. Neither seam may be removed.

**T-015** — Flags are package-level variables, so tests share global state. A test calls
`resetFlags()` first.
*Gotcha: `dateTimeFlag.Set("")` returns an error and leaves the previous value in place,
so date flags are zeroed by assignment, not through `Set`.*

## External service etiquette

**T-016** — Every request to the iNaturalist API carries a `User-Agent` identifying
birdsync and its version.
Subject: `http.user_agent` · Value: `birdsync/0.1`

**T-017** — Observation downloads request the maximum supported page size.
Subject: `inat.page_size` · Value: `200`
*Rationale: required by `inat/api-recommended-practices` (see
[sources.md](sources.md)).*

**T-034** — An error from a non-200 API response includes the response body, truncated, not
just the status line.
*Rationale: iNaturalist explains its refusals in the body. Discarding it left every failure
looking identical — `bad HTTP status: 422 Unprocessable Entity` whether the file was too
large, the format unsupported, or the asset withdrawn. Classification under P-063 depends on
the status code, but the human reading the log needs the reason.*
*An HTML body is dropped rather than included: it comes from a proxy rather than the API, and
a 413 arrives as a seven-line nginx page. A message that is kept is collapsed onto one line.*

**T-035** — Requests to the iNaturalist API are paced to the rate its recommended practices
ask for: about one per second, and about 10,000 per day.
Subject: `inat.request_rate` · Value: `1/second, 10000/day`
*Rationale: `inat-api/R1`, a governing source. Exceeding it returns HTTP 429, and the page
warns that persistent offenders may be IP-blocked. A first sync of a media-heavy account is
thousands of writes: reads are not the exposure.*
*Enforced in `Client.roundTrip`, which every request in the package passes through — the only
reason a single choke point suffices. The daily cap is not enforced: birdsync keeps no state
between runs, so it cannot know how many requests today has already seen.*

*This limit was already being met before the limiter existed, by an argument the maintainer
had reasoned through and not written down: observations are fetched in large pages, and
between writes birdsync downloads a media file, so the operations are slow enough to stay
under one per second on their own. Measurement agrees — a real run averaged about five seconds
per upload and three per download page, roughly 0.2 requests per second.*

*The argument is nonetheless worth replacing with a mechanism, for two reasons. It is
**conditional**: with `--verifiable=false` there is no media to fetch between creates, so the
spacing that made it true disappears. And it is **load-bearing on slowness** — the first
change that adds concurrency, or that makes uploads faster, silently removes the guarantee
with nothing to notice. An explicit limiter inverts that: compliance stops depending on the
code staying slow, which is what makes it safe to parallelize later.*

**T-036** — The observation download pages with `id_above` and an ascending id sort, never
with page numbers, and stops when a page comes back short.
*Rationale: `inat-api/R3` and [CR-011](decisions.md#cr-011--the-download-cannot-page-past-10000-results).
`total_results` shrinks as the cursor advances, so it cannot decide when to stop. The client
requires the cursor to advance on every page: without that check, a server ignoring `id_above`
makes the loop re-fetch one page forever.*

## Data format handling

**T-018** — The eBird CSV is read by header name, never by column position, with
`FieldsPerRecord = -1`.
*Rationale: the export has a variable number of columns between users.*

**T-019** — Any comparison of an eBird observation date goes through `Record.Observed()`,
never `Record.Date` directly.
*Rationale: eBird writes `2006-01-02` or `1/2/2006`, with an optional time. Keying on the
raw field was a real bug; regression test: `TestFuzzyMatchDateFormats`.*

**T-020** — An empty taxon name never enters the fuzzy-match index (implements P-032).

**T-021** — Code that needs to know whether a Macaulay Library asset is a photo or a
sound must download it first; the ID does not encode it.

## Resource use

**T-022** — Memory scales with the export: the CSV is read whole with `ReadAll`, and all
downloaded iNaturalist observations are held in memory. This is accepted at the current
scale (tens of thousands of records) and is a known ceiling, not a design goal.

**T-023** — Temporary files created while downloading media are deleted before the run
ends.
*The file belongs to the caller once `DownloadMLAsset` returns it, so `birdsync()` removes
it after the upload attempt, successful or not — at most one asset is on disk at a time.*
*The download's error paths honor it too: `downloadMLAsset` closes and removes the temp file
on every path but success. Leaving it behind — and on Windows leaving the handle open — is what
[issue #1](https://github.com/Sajmani/birdsync/issues/1) reported as "The process cannot access
the file because it is being used by another process" on a later attempt.*

## Code conventions

**T-024** — `gofmt -l .` prints nothing.

**T-025** — `go vet ./...` is clean.

**T-026** — Errors are wrapped with the calling function's name:
`fmt.Errorf("DownloadMLAsset(%s): %w", id, err)`. `inat.Client.roundTrip` is the
exception, wrapping with a phrase describing the step.

**T-027** — `log.Fatal` is acceptable in `main` and in `tools/`; the `ebird` and `inat`
packages return errors to their caller.
*Enforced by static analysis rather than convention, since a test can show that one
function returns an error but only analysis can show that no function exits. `os.Exit` is
checked too, being the same thing under another name.*

**T-028** — User-visible progress goes through `log.Printf`; verbose detail goes through
`debugf`, gated on `--debug`.

**T-029** — Comments explain *why*. The explanations of eBird and iNaturalist quirks are
the most valuable comments in the repository and are preserved through refactors.

## Continuous integration

**T-030** — CI runs build, vet, gofmt, and the test suite under the race detector, on
every push to `main` and every pull request, using the toolchain declared in `go.mod`.
*The actions it depends on are kept current: GitHub retires the Node runtime older versions
target, and a deprecation warning today is a broken build later. `actions/checkout` and
`actions/setup-go` moved to Node 24 at v5 and v6 respectively; both are pinned at v7.*

**T-031** — CI additionally runs against the latest stable Go as an early warning, and
that job may fail without blocking a pull request.

## Project bindings

Per [process.md](process.md#project-bindings):

| Binding | Value |
| --- | --- |
| Standing checks | `go build ./...`, `go vet ./...`, `gofmt -l .` (must print nothing), `go test -race ./...` |
| ID citation syntax | Go line comment: `// Verifies: P-020.` |
| Artifact locations | `spec/`; vendored sources in `spec/sources/<name>-<pin>/` |
| Subject vocabulary | Dotted paths rooted at `build`, `cli`, `http`, `inat`, `log`, `media`, `observation`, `sync`. The owner adds roots. |
| Hard safety constraints | T-005, T-010, T-012, T-032, T-033 |
| External systems | `api.inaturalist.org` and `cdn.download.ams.birds.cornell.edu` are never contacted from a test; use `httptest` via T-014 |
| Approval authority | The repository owner (Sajmani) for every gate, waiver, amendment, and risk acceptance |
| Source refresh cadence | Reviewed when a `spec/` change is made, and at minimum annually |
| Existing docs | `README.md` owns usage; `CONTRIBUTING.md` owns workflow mechanics; `AGENTS.md` owns agent-specific operating notes; `spec/` owns requirements, criteria, and design |

## Open questions

Resolved at Gate 1: temp-file cleanup is a defect (T-023), and the agent and contributor
guides keep a terse summary with citations rather than being reduced to bare links.
Still open:

1. **`inat.Client.UploadMedia` has no test.** The multipart request it builds
   is only exercised through a mock. It is the largest untested surface, and it is also the
   one that writes binary data to a user's account.
2. **T-022's memory ceiling grows** now that P-061 downloads every observation rather
   than birds only. Nobody has measured what a large account costs. Worth a benchmark
   criterion in phase 2 rather than a guess here.
3. **Rate limiting is unexpressed.** `inat-api` is adopted as governing, but only its
   `User-Agent` and page-size guidance became requirements (T-016, T-017). Its request-rate
   guidance has no `T-###` and no check. See [sources.md](sources.md#adopted).
