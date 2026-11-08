package birdsync

# Module, testability, conventions, and CI

Go module and dependencies, testability seams, code and spelling conventions, and continuous
integration for [`birdsync`](README.md).

## Included requirements

### Module and dependencies

- [tech.singleModule](../../../local/tech/deps.md#singlemodule) — A single Go module with no build system
- [tech.oneExternalDependency](../../../local/tech/deps.md#oneexternaldependency) — Exactly one external dependency
- [tech.goDirectivePolicy](../../../local/tech/deps.md#godirectivepolicy) — go and toolchain directives change only deliberately
- [tech.platformIndependent](../../../local/tech/deps.md#platformindependent) — Builds and runs on every Go platform

### Testability and live-service isolation

- [tech.noLiveServicesInTests](../../../local/tech/testing.md#noliveservicesintests) — Tests never contact a live service
- [tech.clientInterfaces](../../../local/tech/testing.md#clientinterfaces) — birdsync() takes its clients as interfaces
- [tech.baseURLSeams](../../../local/tech/testing.md#baseurlseams) — Both services have base-URL test seams
- [tech.resetFlagsInTests](../../../local/tech/testing.md#resetflagsintests) — Tests call resetFlags first

### Code and prose conventions

- [tech.gofmtClean](../../../local/tech/style.md#gofmtclean) — gofmt -l . prints nothing
- [tech.vetClean](../../../local/tech/style.md#vetclean) — go vet is clean
- [tech.errorWrapping](../../../local/tech/style.md#errorwrapping) — Errors are wrapped with the calling function's name
- [tech.noLogFatalInLibraries](../../../local/tech/style.md#nologfatalinlibraries) — log.Fatal only in main and tools/
- [tech.progressLogging](../../../local/tech/style.md#progresslogging) — Progress via log.Printf, detail via debugf
- [tech.americanSpelling](../../../local/tech/style.md#americanspelling) — Prose uses American spelling
- [tech.quotationsNeverRespelled](../../../local/tech/style.md#quotationsneverrespelled) — Spelling changes never touch quotations or vendored documents
- [tech.commentsExplainWhy](../../../local/tech/style.md#commentsexplainwhy) — Comments explain why

### Continuous integration

- [tech.ciStandingChecks](../../../local/tech/ci.md#cistandingchecks) — CI runs the standing checks
- [tech.ciLatestStable](../../../local/tech/ci.md#cilateststable) — CI also runs latest stable Go, non-blocking

## Acceptance criteria

## buildsEverywhere

**Everything builds**

Formerly: AC-001
Covers: [tech.singleModule](../../../local/tech/deps.md#singlemodule), [tech.platformIndependent](../../../local/tech/deps.md#platformindependent), [product.toolsNotProduct](../../../local/product/scope.md#toolsnotproduct)

`go build ./...` compiles every package, including `tools/`.

Limit: partial for platform independence: one OS.

## vetPasses

**go vet is clean**

Formerly: AC-002
Covers: [tech.vetClean](../../../local/tech/style.md#vetclean)

`go vet ./...` reports nothing.

## formatted

**Code is gofmt'd**

Formerly: AC-003
Covers: [tech.gofmtClean](../../../local/tech/style.md#gofmtclean)

`gofmt -l .` prints nothing.

## suitePassesUnderRace

**The suite passes under the race detector**

Formerly: AC-004
Covers: [tech.ciStandingChecks](../../../local/tech/ci.md#cistandingchecks)

`go test -race ./...` passes with no network and no live accounts.

Limit: runs every criterion in this package.

## noLiveHostsInTests

**No test names a live service host**

Formerly: AC-005
Covers: [tech.noLiveServicesInTests](../../../local/tech/testing.md#noliveservicesintests)

No `_test.go` file contains the iNaturalist API host or the Macaulay Library CDN host.

Limit: a substring scan; a host assembled from parts escapes it.

## ciRunsStandingChecks

**CI runs the standing checks on push and pull request**

Formerly: AC-022
Covers: [tech.ciStandingChecks](../../../local/tech/ci.md#cistandingchecks), [tech.ciLatestStable](../../../local/tech/ci.md#cilateststable)

`.github/workflows/ci.yml` runs build, vet, gofmt, and the race-enabled suite on every push to main and every pull request, plus a non-blocking latest-stable job.

## noLogFatalOutsideMain

**Library packages never exit**

Formerly: AC-027
Covers: [tech.noLogFatalInLibraries](../../../local/tech/style.md#nologfatalinlibraries)

No non-main package in the module calls `log.Fatal` or `os.Exit`.

## americanSpellings

**Prose is American English and quotations are untouched**

Formerly: AC-041
Covers: [tech.americanSpelling](../../../local/tech/style.md#americanspelling), [tech.quotationsNeverRespelled](../../../local/tech/style.md#quotationsneverrespelled)

No Markdown or Go file outside vendored documents uses a listed British spelling, and blockquotes and vendored documents are never flagged.
