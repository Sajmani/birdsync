package tech

Constraints the implementation must satisfy, independent of what birdsync does for a
user: one item per heading, and each item's first paragraph after its fields is the
complete rule. Search `^## ` for the list, or `## <name>` for one item; `Formerly:` gives
the identifier this item had under the predecessor method (specman).

Migrated from `spec/tech.md` at
[`dc25e3b`](https://github.com/Sajmani/birdsync/blob/dc25e3b6756bd59c2de931e6d080d47318617da4/spec/tech.md),
where it was approved at specman's Gate 1 on 2026-08-09 and amended since. Goals were
accepted on 2026-11-04 and are cited in each `Goals:` field.

## singleModule

**A single Go module with no build system**

Formerly: T-001
Area: Module and dependencies
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

birdsync is a single Go module, `github.com/Sajmani/birdsync`, with no build
system, no code generation, and no linter configuration beyond `gofmt` and `go vet`.

## oneExternalDependency

**Exactly one external dependency**

Formerly: T-002
Area: Module and dependencies
Subject: `build.external_dependencies.count`
Value: `1`
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

The module has exactly one external dependency,
`github.com/google/uuid`. Adding another requires the owner's approval.

*Rationale: the small footprint is deliberate; it keeps `go install` fast and the
supply-chain surface near zero.*

## goDirectivePolicy

**go and toolchain directives change only deliberately**

Formerly: T-003
Area: Module and dependencies
Subject: `build.go_directive`
Value: `1.23.0`
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

The `go` and `toolchain` directives are raised only as a deliberate,
standalone change, never as a side effect of other work.

*Rationale: users on the default `GOTOOLCHAIN=auto` silently download a newer toolchain,
but anyone pinned to `GOTOOLCHAIN=local` on an older Go is blocked outright.*

## platformIndependent

**Builds and runs on every Go platform**

Formerly: T-004
Area: Module and dependencies
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

The tool must build and run on any platform the Go toolchain supports, using
only the standard library plus [oneExternalDependency](#oneexternaldependency). No platform-specific code.

## dryRunGateAtCallSite

**Mutations are gated on --dryrun at the call site**

Formerly: T-005
Area: Safety invariants
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

Every operation that creates, updates, deletes, or uploads is gated on
`--dryrun` **at the call site in `birdsync()`**, not inside the client.

*Rationale: `inat.Client` is shared with `tools/`, which has no such flag. The gate
belongs where the decision is made, not in a client that other callers reuse.*

## dryRunPrefixLogged

**A skipped mutation is logged with DRYRUN:**

Formerly: T-006
Area: Safety invariants
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

A skipped mutation is logged with a `DRYRUN:` prefix (implements [product.dryRunLogPrefix](../product/dryrun.md#dryrunlogprefix)).

## honestDryRunCounters

**Dry-run counters never report work that did not happen**

Formerly: T-007
Area: Safety invariants
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

A dry run's counters must not report work that did not happen. Where a dry
run cannot know something, it counts it as unknown rather than guessing.

*Status: **not yet satisfied** — the created and updated counters are labeled as though
the work happened. [decisions.dryRunCountersRelabeled](../../decisions.md#dryruncountersrelabeled) resolved this
in favour of relabeling; implements [product.dryRunCounterLabels](../product/reporting.md#dryruncounterlabels).*

## clientMutatesUnconditionally

**inat.Client mutates unconditionally**

Formerly: T-008
Area: Safety invariants
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

`inat.Client`'s methods mutate unconditionally and know nothing about
`--dryrun`. This is deliberate, not an oversight.

## updateKeepsPhotos

**UpdateObservation always sets ignore_photos**

Formerly: T-009
Area: Safety invariants
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

`UpdateObservation` always sets `ignore_photos`.

*Rationale: without it, updating a description detaches the observation's photos.*

## noLiveServicesInTests

**Tests never contact a live service**

Formerly: T-010
Area: Safety invariants
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

Tests never contact `api.inaturalist.org`, the Macaulay Library CDN, or any
other live service, and never depend on a real account's contents.

## neverRunMutatingTools

**Nothing runs the mutating tools**

Formerly: T-011
Area: Safety invariants
Status: withdrawn — the mutating tools were deleted rather than governed by a rule
Goals:

Nothing in the build, test, or development workflow runs a program in
`tools/`. They create, update, and delete observations in whatever account the
environment's credentials point at, and `purge` is checked in with its guard off.
nobody could check. Superseded by [toolsReadOnly](#toolsreadonly) and [neverRunTools](#neverruntools).

## toolsReadOnly

**No program in tools/ mutates**

Formerly: T-032
Area: Safety invariants
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

No program in `tools/` performs a mutating operation. Nothing there may
create, update, delete, or upload.

*Rationale: `tools/` held six programs that modified real accounts, guarded by a `debug`
constant that one of them shipped with turned off and by an instruction not to run them.
A structural property that a check can enforce is worth more than a warning people have to
read and remember. Deleted tools remain in the history if one is needed again.*

## neverRunTools

**Nothing in the workflow runs a program in tools/**

Formerly: T-033
Area: Safety invariants
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

Nothing in the build, test, or development workflow runs a program in
`tools/`. They contact the live service with the user's real credentials, which no
automated process should do on the user's behalf.

## credentialsNeverExposed

**Credentials are never logged, echoed, or written**

Formerly: T-012
Area: Safety invariants
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

Credentials are never logged, echoed, or written to a file.

## clientInterfaces

**birdsync() takes its clients as interfaces**

Formerly: T-013
Area: Testability seams
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

`birdsync()` takes its two clients as the `ebirdClient` and `inatClient`
interfaces, so the sync loop can be driven without a network.

## baseURLSeams

**Both services have base-URL test seams**

Formerly: T-014
Area: Testability seams
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

`inat.NewClient` takes a base URL, and `ebird.DownloadMLAsset` delegates to
an unexported `downloadMLAsset(baseURL, id)`, so an `httptest` server can stand in for
either service. Neither seam may be removed.

## resetFlagsInTests

**Tests call resetFlags first**

Formerly: T-015
Area: Testability seams
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

Flags are package-level variables, so tests share global state. A test calls
`resetFlags()` first.

*Gotcha: `dateTimeFlag.Set("")` returns an error and leaves the previous value in place,
so date flags are zeroed by assignment, not through `Set`.*

## userAgent

**Every iNaturalist request carries a birdsync User-Agent**

Formerly: T-016
Area: External service etiquette
Subject: `http.user_agent`
Value: `birdsync/0.1`
Goals: [goals.communityNotBurdened](../../../goals/goals.md#communitynotburdened)

Every request to the iNaturalist API carries a `User-Agent` identifying
birdsync and its version.

## maxPageSize

**Downloads request the maximum page size**

Formerly: T-017
Area: External service etiquette
Subject: `inat.page_size`
Value: `200`
Goals: [goals.communityNotBurdened](../../../goals/goals.md#communitynotburdened)

Observation downloads request the maximum supported page size.

*Rationale: required by `inat/api-recommended-practices` (see
[imports.md](../../imports.md)).*

## errorsIncludeBody

**API errors include the response body**

Formerly: T-034
Area: External service etiquette
Goals: [goals.communityNotBurdened](../../../goals/goals.md#communitynotburdened)

An error from a non-200 API response includes the response body, truncated, not
just the status line.

*Rationale: iNaturalist explains its refusals in the body. Discarding it left every failure
looking identical — `bad HTTP status: 422 Unprocessable Entity` whether the file was too
large, the format unsupported, or the asset withdrawn. Classification under [product.permanentUploadFailureRecorded](../product/failures.md#permanentuploadfailurerecorded) depends on
the status code, but the human reading the log needs the reason.*

*An HTML body is dropped rather than included: it comes from a proxy rather than the API, and
a 413 arrives as a seven-line nginx page. A message that is kept is collapsed onto one line.*

## requestPacing

**iNaturalist requests are paced**

Formerly: T-035
Area: External service etiquette
Subject: `inat.request_rate`
Value: `1/second, 10000/day`
Goals: [goals.communityNotBurdened](../../../goals/goals.md#communitynotburdened)

Requests to the iNaturalist API are paced to the rate its recommended practices
ask for: about one per second, and about 10,000 per day.

*Rationale: [inatapi.requestRate](../../imported/inaturalist.org/inatapi/inatapi.md#requestrate), a governing source. Exceeding it returns HTTP 429, and the page
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

## idAbovePaging

**The download pages by id_above**

Formerly: T-036
Area: External service etiquette
Goals: [goals.communityNotBurdened](../../../goals/goals.md#communitynotburdened)

The observation download pages with `id_above` and an ascending id sort, never
with page numbers, and stops when a page comes back short.

*Rationale: [inatapi.idAbovePastTenThousand](../../imported/inaturalist.org/inatapi/inatapi.md#idabovepasttenthousand) and [decisions.pageByIDAbove](../../decisions.md#pagebyidabove).
`total_results` shrinks as the cursor advances, so it cannot decide when to stop. The client
requires the cursor to advance on every page: without that check, a server ignoring `id_above`
makes the loop re-fetch one page forever.*

## csvByHeaderName

**The CSV is read by header name**

Formerly: T-018
Area: Data format handling
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

The eBird CSV is read by header name, never by column position, with
`FieldsPerRecord = -1`.

*Rationale: the export has a variable number of columns between users.*

## datesViaObserved

**eBird dates are compared via Record.Observed()**

Formerly: T-019
Area: Data format handling
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

Any comparison of an eBird observation date goes through `Record.Observed()`,
never `Record.Date` directly.

*Rationale: eBird writes `2006-01-02` or `1/2/2006`, with an optional time. Keying on the
raw field was a real bug; regression test: `TestFuzzyMatchDateFormats`.*

## emptyNameNotIndexed

**An empty taxon name never enters the fuzzy index**

Formerly: T-020
Area: Data format handling
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

An empty taxon name never enters the fuzzy-match index (implements [product.fuzzyNeverEmptyName](../product/filtering.md#fuzzyneveremptyname)).

## assetTypeNeedsDownload

**Asset type is known only after download**

Formerly: T-021
Area: Data format handling
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

Code that needs to know whether a Macaulay Library asset is a photo or a
sound must download it first; the ID does not encode it.

## memoryScalesWithExport

**Memory scales with the export**

Formerly: T-022
Area: Resource use
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

Memory scales with the export: the CSV is read whole with `ReadAll`, and all
downloaded iNaturalist observations are held in memory. This is accepted at the current
scale (tens of thousands of records) and is a known ceiling, not a design goal.

## tempFilesDeleted

**Temporary media files are deleted**

Formerly: T-023
Area: Resource use
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

Temporary files created while downloading media are deleted before the run
ends.

*The file belongs to the caller once `DownloadMLAsset` returns it, so `birdsync()` removes
it after the upload attempt, successful or not — at most one asset is on disk at a time.*

*The download's error paths honor it too: `downloadMLAsset` closes and removes the temp file
on every path but success. Leaving it behind — and on Windows leaving the handle open — is what
[issue #1](https://github.com/Sajmani/birdsync/issues/1) reported as "The process cannot access
the file because it is being used by another process" on a later attempt.*

## gofmtClean

**gofmt -l . prints nothing**

Formerly: T-024
Area: Code conventions
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

`gofmt -l .` prints nothing.

## vetClean

**go vet is clean**

Formerly: T-025
Area: Code conventions
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

`go vet ./...` is clean.

## errorWrapping

**Errors are wrapped with the calling function's name**

Formerly: T-026
Area: Code conventions
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

Errors are wrapped with the calling function's name:
`fmt.Errorf("DownloadMLAsset(%s): %w", id, err)`. `inat.Client.roundTrip` is the
exception, wrapping with a phrase describing the step.

## noLogFatalInLibraries

**log.Fatal only in main and tools/**

Formerly: T-027
Area: Code conventions
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

`log.Fatal` is acceptable in `main` and in `tools/`; the `ebird` and `inat`
packages return errors to their caller.

*Enforced by static analysis rather than convention, since a test can show that one
function returns an error but only analysis can show that no function exits. `os.Exit` is
checked too, being the same thing under another name.*

## progressLogging

**Progress via log.Printf, detail via debugf**

Formerly: T-028
Area: Code conventions
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

User-visible progress goes through `log.Printf`; verbose detail goes through
`debugf`, gated on `--debug`.

## americanSpelling

**Prose uses American spelling**

Formerly: T-037
Area: Code conventions
Subject: `docs.spelling`
Value: `en-US`
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

Prose in this repository — documentation, comments, commit messages — uses
American spellings.

*Rationale: consistency, and the maintainer's preference. It is the weakest requirement here
on its own; it is written down for the exception below.*

## quotationsNeverRespelled

**Spelling changes never touch quotations or vendored documents**

Formerly: T-038
Area: Code conventions
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

A spelling change never touches a quotation or a vendored document.

*Rationale: this is the part with consequences. A bulk find-and-replace across `river/` would
silently rewrite text quoted from someone else's terms of service, turning a verbatim citation
into a paraphrase that still looks verbatim, and would break the hashes recorded in each
source's `PROVENANCE.md`. The rule is not "be careful"; the check skips blockquote lines and
never opens `river/requirements/imported/**` or the skill copy in `.agents/skills/` (under
specman these were `spec/` and `spec/sources/**`).*

## commentsExplainWhy

**Comments explain why**

Formerly: T-029
Area: Code conventions
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

Comments explain *why*. The explanations of eBird and iNaturalist quirks are
the most valuable comments in the repository and are preserved through refactors.

## ciStandingChecks

**CI runs the standing checks**

Formerly: T-030
Area: Continuous integration
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

CI runs build, vet, gofmt, and the test suite under the race detector, on
every push to `main` and every pull request, using the toolchain declared in `go.mod`.

*The actions it depends on are kept current: GitHub retires the Node runtime older versions
target, and a deprecation warning today is a broken build later. `actions/checkout` and
`actions/setup-go` moved to Node 24 at v5 and v6 respectively; both are pinned at v7.*

## ciLatestStable

**CI also runs latest stable Go, non-blocking**

Formerly: T-031
Area: Continuous integration
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

CI additionally runs against the latest stable Go as an early warning, and
that job may fail without blocking a pull request.
