# Conformance

What birdsync's code is checked against, and how. This is the mutable record of which guardrail
implements each acceptance criterion in [`composed/internal/birdsync/`](requirements/composed/internal/birdsync/README.md), and whether it has been
watched failing. A criterion has been watched failing only if its section here carries a
`Watched failing:` line. Search `^## ` for the criteria, or `## <name>` for one.

## Built against

- Package: [`composed/internal/birdsync`](requirements/composed/internal/birdsync/README.md)
- Version: **none approved**. The package has no `Approves:` line in `decisions.md`, so this
  record maps guardrails to a draft. Recheck it when the owner approves a version.
- Checked on: 2026-11-08, with `go test -race ./...` passing and `river check` passing.

## Counts

Computed on 2026-11-08 by script, not maintained by hand. Recount rather than editing these
numbers.

- Requirements in sources: 126 (107 local, 19 imported). 101 local are live; 6 are withdrawn and excluded
- Live local requirements covered by a criterion, directly or through a composed requirement
  they are composed into: 73. Uncovered: 28, counted, not failed
- Composed requirements: 11. Criteria: 56. Watched failing on record: 9. Watched not failing: 1.
  Not recorded: 32. With no guardrail yet: 14
- Test functions: 50. Cited by a criterion: 46. Project checks: 2. Cited by none: 2
- Annotations harvested: 43, from 6 files. Annotated tests no criterion cites: 2

## buildsEverywhere

Criterion: [birdsync.buildsEverywhere](requirements/composed/internal/birdsync/engineering.md#buildseverywhere)
Guardrails: `go build ./...`, in CI
Gates: yes, in CI
Specman status: verified
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## vetPasses

Criterion: [birdsync.vetPasses](requirements/composed/internal/birdsync/engineering.md#vetpasses)
Guardrails: `go vet ./...`, in CI
Gates: yes, in CI
Specman status: verified
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## formatted

Criterion: [birdsync.formatted](requirements/composed/internal/birdsync/engineering.md#formatted)
Guardrails: `gofmt -l .`, in CI
Gates: yes, in CI
Specman status: verified
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## suitePassesUnderRace

Criterion: [birdsync.suitePassesUnderRace](requirements/composed/internal/birdsync/engineering.md#suitepassesunderrace)
Guardrails: `go test -race ./...`, in CI
Gates: yes, in CI
Specman status: verified
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## noLiveHostsInTests

Criterion: [birdsync.noLiveHostsInTests](requirements/composed/internal/birdsync/engineering.md#nolivehostsintests)
Guardrails: `TestNoLiveServiceHostsInTests` (`guard_test.go`)
Gates: yes, in CI
Specman status: verified
Watched failing: 2026-08-09, by adding a scratch `_test.go` naming the live API host; the test failed (specman acceptance.md, commit `defc5e3`)

## dryRunIssuesNoWrites

Criterion: [birdsync.dryRunIssuesNoWrites](requirements/composed/internal/birdsync/dryrun.md#dryrunissuesnowrites)
Guardrails: `TestDryRunIssuesNoWrites` (`birdsync_test.go`)
Gates: yes, in CI
Specman status: verified
Watched failing: 2026-08-09, by changing `if dryRun` to `if false` at the `CreateObservation` gate; the test failed (commit `defc5e3`)

## createdObservationIdentity

Criterion: [birdsync.createdObservationIdentity](requirements/composed/internal/birdsync/observations.md#createdobservationidentity)
Guardrails: `TestCreatedObservationContent` (`birdsync_test.go`)
Gates: yes, in CI
Specman status: verified
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## createdObservationPlacement

Criterion: [birdsync.createdObservationPlacement](requirements/composed/internal/birdsync/observations.md#createdobservationplacement)
Guardrails: `TestCreatedObservationContent` (`birdsync_test.go`)
Gates: yes, in CI
Specman status: verified
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## createdObservationMedia

Criterion: [birdsync.createdObservationMedia](requirements/composed/internal/birdsync/media.md#createdobservationmedia)
Guardrails: `TestCreatedObservationContent` (`birdsync_test.go`)
Gates: yes, in CI
Specman status: verified
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## updateKeepsPhotos

Criterion: [birdsync.updateKeepsPhotos](requirements/composed/internal/birdsync/media.md#updatekeepsphotos)
Guardrails: `TestClient_UpdateObservation` (`inat/client_test.go`)
Gates: yes, in CI
Specman status: verified
Watched failing: 2026-08-09, by changing `IgnorePhotos: true` to `false`; the test failed (commit `defc5e3`)

## userAgentSent

Criterion: [birdsync.userAgentSent](requirements/composed/internal/birdsync/download.md#useragentsent)
Guardrails: `TestClient_CreateObservation` (`inat/client_test.go`)
Gates: yes, in CI
Specman status: verified
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## downloadUsesMaxPageSize

Criterion: [birdsync.downloadUsesMaxPageSize](requirements/composed/internal/birdsync/download.md#downloadusesmaxpagesize)
Guardrails: `TestDownloadObservations` (`inat/inat_test.go`)
Gates: yes, in CI
Specman status: verified
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## downloadHonorsDateWindow

Criterion: [birdsync.downloadHonorsDateWindow](requirements/composed/internal/birdsync/download.md#downloadhonorsdatewindow)
Guardrails: `TestDownloadObservationsDateWindow` (`inat/inat_test.go`), `TestDownloadObservationsNoDateWindow` (`inat/inat_test.go`)
Gates: yes, in CI
Specman status: verified
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## skipOrderEnforced

Criterion: [birdsync.skipOrderEnforced](requirements/composed/internal/birdsync/filtering.md#skiporderenforced)
Guardrails: `TestBirdsync` (`birdsync_test.go`)
Gates: yes, in CI
Specman status: verified
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## mediaResyncAdditive

Criterion: [birdsync.mediaResyncAdditive](requirements/composed/internal/birdsync/media.md#mediaresyncadditive)
Guardrails: `TestUpdateMedia` (`birdsync_test.go`)
Gates: yes, in CI
Specman status: verified
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## fuzzyHandlesBothDateFormats

Criterion: [birdsync.fuzzyHandlesBothDateFormats](requirements/composed/internal/birdsync/filtering.md#fuzzyhandlesbothdateformats)
Guardrails: `TestFuzzyMatchDateFormats` (`birdsync_test.go`)
Gates: yes, in CI
Specman status: verified
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## fuzzyIgnoresEmptyNames

Criterion: [birdsync.fuzzyIgnoresEmptyNames](requirements/composed/internal/birdsync/filtering.md#fuzzyignoresemptynames)
Guardrails: `TestFuzzyMatchIgnoresEmptyNames` (`birdsync_test.go`)
Gates: yes, in CI
Specman status: verified
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## dryRunMediaCountUnclassified

Criterion: [birdsync.dryRunMediaCountUnclassified](requirements/composed/internal/birdsync/dryrun.md#dryrunmediacountunclassified)
Guardrails: `TestDryRunMediaCount` (`birdsync_test.go`)
Gates: yes, in CI
Specman status: partial — see [below](#criteria-that-do-not-bite)
Watched not failing: 2026-08-09, the `if dryRun` → `if false` mutation at the `CreateObservation` gate left `TestDryRunMediaCount` green: it does not catch a broken dry-run gate (commit `defc5e3`). It still decides the unclassified count; that it fails when that count is wrong has not been recorded

## mediaChangesReported

Criterion: [birdsync.mediaChangesReported](requirements/composed/internal/birdsync/media.md#mediachangesreported)
Guardrails: `TestMediaChange` (`media_test.go`)
Gates: yes, in CI
Specman status: verified
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## mlAssetDownload

Criterion: [birdsync.mlAssetDownload](requirements/composed/internal/birdsync/media.md#mlassetdownload)
Guardrails: `TestDownloadMLAsset_Photo` (`ebird/ebird_test.go`), `TestDownloadMLAsset_Sound` (`ebird/ebird_test.go`)
Gates: yes, in CI
Specman status: verified
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## recordObservedParsing

Criterion: [birdsync.recordObservedParsing](requirements/composed/internal/birdsync/filtering.md#recordobservedparsing)
Guardrails: `TestRecord_Observed` (`ebird/ebird_test.go`)
Gates: yes, in CI
Specman status: verified
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## csvReadByHeader

Criterion: [birdsync.csvReadByHeader](requirements/composed/internal/birdsync/filtering.md#csvreadbyheader)
Guardrails: `TestRecords` (`ebird/ebird_test.go`)
Gates: yes, in CI
Specman status: verified
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## syncKeyValidity

Criterion: [birdsync.syncKeyValidity](requirements/composed/internal/birdsync/filtering.md#synckeyvalidity)
Guardrails: `TestObservationID_Valid` (`ebird/ebird_test.go`)
Gates: yes, in CI
Specman status: verified
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## ciRunsStandingChecks

Criterion: [birdsync.ciRunsStandingChecks](requirements/composed/internal/birdsync/engineering.md#cirunsstandingchecks)
Guardrails: `.github/workflows/ci.yml`
Gates: yes, in CI
Specman status: verified
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## summaryLabels

Criterion: [birdsync.summaryLabels](requirements/composed/internal/birdsync/dryrun.md#summarylabels)
Guardrails: `TestSummaryDryRunLabels` (`birdsync_test.go`), `TestSummaryRealRunLabels` (`birdsync_test.go`), `TestSummaryConditionalLines` (`birdsync_test.go`)
Gates: yes, in CI
Specman status: verified (landed as a pending criterion)
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## untaxonedObservationRecognized

Criterion: [birdsync.untaxonedObservationRecognized](requirements/composed/internal/birdsync/filtering.md#untaxonedobservationrecognized)
Guardrails: `TestDownloadObservationsNoTaxonFilter` (`inat/inat_test.go`), `TestUntaxonedObservationIsRecognized` (`birdsync_test.go`)
Gates: yes, in CI
Specman status: verified (landed as a pending criterion)
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## invalidRecordsSkipped

Criterion: [birdsync.invalidRecordsSkipped](requirements/composed/internal/birdsync/failures.md#invalidrecordsskipped)
Guardrails: `TestInvalidRecordsAreSkipped` (`birdsync_test.go`)
Gates: yes, in CI
Specman status: verified (landed as a pending criterion)
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## tempFilesCleanedUp

Criterion: [birdsync.tempFilesCleanedUp](requirements/composed/internal/birdsync/media.md#tempfilescleanedup)
Guardrails: `TestTempFilesAreCleanedUp` (`birdsync_test.go`)
Gates: yes, in CI
Specman status: verified (landed as a pending criterion)
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## noLogFatalOutsideMain

Criterion: [birdsync.noLogFatalOutsideMain](requirements/composed/internal/birdsync/engineering.md#nologfataloutsidemain)
Guardrails: `TestNoLogFatalInLibraryPackages` (`guard_test.go`)
Gates: yes, in CI
Specman status: verified (landed as a pending criterion)
Watched failing: 2026-08-09, by restoring a `log.Fatal` in `inat`; the test failed (commit `cecd667`)

## toolsAreReadOnly

Criterion: [birdsync.toolsAreReadOnly](requirements/composed/internal/birdsync/scope.md#toolsarereadonly)
Guardrails: `TestToolsAreReadOnly` (`guard_test.go`)
Gates: yes, in CI
Specman status: verified
Watched failing: 2026-08-09, by adding a scratch tool calling `DeleteObservation`; the test failed (commit `cecd667`)

## canonicalExtensions

Criterion: [birdsync.canonicalExtensions](requirements/composed/internal/birdsync/media.md#canonicalextensions)
Guardrails: `TestFileExtension` (`ebird/ebird_test.go`), `TestDownloadMLAsset_Photo` (`ebird/ebird_test.go`), `TestDownloadMLAsset_Sound` (`ebird/ebird_test.go`)
Gates: yes, in CI
Specman status: verified
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## failedUploadRetried

Criterion: [birdsync.failedUploadRetried](requirements/composed/internal/birdsync/media.md#faileduploadretried)
Guardrails: `TestFailedUploadIsRetriedNextRun` (`birdsync_test.go`), `TestAllUploadsFailingLeavesNoUpdate` (`birdsync_test.go`)
Gates: yes, in CI
Specman status: verified
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## permanentVersusTransientUpload

Criterion: [birdsync.permanentVersusTransientUpload](requirements/composed/internal/birdsync/failures.md#permanentversustransientupload)
Guardrails: `TestPermanentUploadFailureIsNotRetried` (`birdsync_test.go`), `TestTransientUploadFailureIsRetried` (`birdsync_test.go`), `TestStatusErrorPermanence` (`inat/client_test.go`)
Gates: yes, in CI
Specman status: verified
Watched failing: 2026-08-09, by forcing `Permanent()` to `true`, then to `false`; the tests failed both ways (commit `27ecd40`)

## statusErrorCarriesBody

Criterion: [birdsync.statusErrorCarriesBody](requirements/composed/internal/birdsync/download.md#statuserrorcarriesbody)
Guardrails: `TestStatusErrorIncludesBody` (`inat/client_test.go`), `TestStatusErrorDropsHTMLBody` (`inat/client_test.go`), `TestStatusErrorCollapsesWhitespace` (`inat/client_test.go`)
Gates: yes, in CI
Specman status: verified
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## failedAssetsReported

Criterion: [birdsync.failedAssetsReported](requirements/composed/internal/birdsync/failures.md#failedassetsreported)
Guardrails: `TestMediaChange` (`media_test.go`)
Gates: yes, in CI
Specman status: verified
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## downloadPagesByID

Criterion: [birdsync.downloadPagesByID](requirements/composed/internal/birdsync/download.md#downloadpagesbyid)
Guardrails: `TestDownloadObservationsPagesByID` (`inat/inat_test.go`), `TestDownloadObservationsCursorMustAdvance` (`inat/inat_test.go`)
Gates: yes, in CI
Specman status: verified
Watched failing: 2026-08-09, by removing `id_above` from the request; it first hung, and fails in two requests since the cursor-advance guard (commit `03999b5`)

## requestsPaced

Criterion: [birdsync.requestsPaced](requirements/composed/internal/birdsync/download.md#requestspaced)
Guardrails: `TestClientPacesRequests` (`inat/client_test.go`), `TestDefaultMinRequestIntervalMatchesGuidance` (`inat/client_test.go`)
Gates: yes, in CI
Specman status: verified
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## badInputRejected

Criterion: [birdsync.badInputRejected](requirements/composed/internal/birdsync/failures.md#badinputrejected)
Guardrails: `TestRecordsRejectsBadInput` (`ebird/ebird_test.go`)
Gates: yes, in CI
Specman status: verified
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## readmeCommunitySection

Criterion: [birdsync.readmeCommunitySection](requirements/composed/internal/birdsync/observations.md#readmecommunitysection)
Guardrails: Human review of README.md
Gates: no; informs, by human review
Specman status: verified
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## downloadCleansUpOnError

Criterion: [birdsync.downloadCleansUpOnError](requirements/composed/internal/birdsync/media.md#downloadcleansuponerror)
Guardrails: `TestDownloadMLAssetCleansUpOnError` (`ebird/ebird_test.go`)
Gates: yes, in CI
Specman status: verified
Not yet watched failing under River. Specman's `verified` status meant the check had been watched failing, but nothing records how or when, so this criterion has no `Watched failing:` line until someone breaks the behavior, sees the guardrail fail, and records it here.

## americanSpellings

Criterion: [birdsync.americanSpellings](requirements/composed/internal/birdsync/engineering.md#americanspellings)
Guardrails: `TestAmericanSpellings` (`guard_test.go`)
Gates: yes, in CI
Specman status: verified — five behaviors mutation-tested
Watched failing: 2026-08-10, by planting a British word in prose, a generated inflection, and an irregular; each failed, and the same word in a quotation and in a vendored transcription correctly passed (commit `2328744`)

## pacingPerDay

Criterion: [birdsync.pacingPerDay](requirements/composed/internal/birdsync/download.md#pacingperday)
Guardrails: `TestClientEnforcesMaxRequests` (`inat/client_test.go`), `TestDailyRequestBudget` (`birdsync_test.go`)
Gates: yes, in CI
Watched failing: 2026-11-08, by disabling the `c.requestCount >= c.maxRequests` guard in `inat/client.go` and forcing `estimateRecentWriteRequests` to return 0 in `birdsync.go`; both tests failed (`4th request error = <nil>, want ErrRequestBudgetExceeded` and `got 2 creates, 2 uploads, 2 updates; want 1, 1, 1`)

## downloadOrderedByID

Criterion: [birdsync.downloadOrderedByID](requirements/composed/internal/birdsync/download.md#downloadorderedbyid)
Guardrails: none mapped yet. Written by composition on 2026-11-05; an existing test may already
check it, and mapping one is the next step
Gates: no
Not watched failing.

## downloadStopsOnShortPage

Criterion: [birdsync.downloadStopsOnShortPage](requirements/composed/internal/birdsync/download.md#downloadstopsonshortpage)
Guardrails: none mapped yet. Written by composition on 2026-11-05; an existing test may already
check it, and mapping one is the next step
Gates: no
Not watched failing.

## jwtFromEnvironmentInHeader

Criterion: [birdsync.jwtFromEnvironmentInHeader](requirements/composed/internal/birdsync/auth.md#jwtfromenvironmentinheader)
Guardrails: none mapped yet. Written by composition on 2026-11-05; an existing test may already
check it, and mapping one is the next step
Gates: no
Not watched failing.

## jwtFromPromptUnwrapped

Criterion: [birdsync.jwtFromPromptUnwrapped](requirements/composed/internal/birdsync/auth.md#jwtfrompromptunwrapped)
Guardrails: none mapped yet. Written by composition on 2026-11-05; an existing test may already
check it, and mapping one is the next step
Gates: no
Not watched failing.

## noBuiltInCredentials

Criterion: [birdsync.noBuiltInCredentials](requirements/composed/internal/birdsync/auth.md#nobuiltincredentials)
Guardrails: none mapped yet. Written by composition on 2026-11-05; an existing test may already
check it, and mapping one is the next step
Gates: no
Not watched failing.

## writesCarrySuppliedToken

Criterion: [birdsync.writesCarrySuppliedToken](requirements/composed/internal/birdsync/auth.md#writescarrysuppliedtoken)
Guardrails: none mapped yet. Written by composition on 2026-11-05; an existing test may already
check it, and mapping one is the next step
Gates: no
Not watched failing.

## tokenOnlyWhereNeeded

Criterion: [birdsync.tokenOnlyWhereNeeded](requirements/composed/internal/birdsync/auth.md#tokenonlywhereneeded)
Guardrails: none mapped yet. Written by composition on 2026-11-05; an existing test may already
check it, and mapping one is the next step
Gates: no
Not watched failing.

## authenticatedDownloadSeesPrivate

Criterion: [birdsync.authenticatedDownloadSeesPrivate](requirements/composed/internal/birdsync/auth.md#authenticateddownloadseesprivate)
Guardrails: none mapped yet. Written by composition on 2026-11-05; an existing test may already
check it, and mapping one is the next step
Gates: no
Not watched failing.

## createsMatchExportRecords

Criterion: [birdsync.createsMatchExportRecords](requirements/composed/internal/birdsync/scope.md#createsmatchexportrecords)
Guardrails: none mapped yet. Written by composition on 2026-11-05; an existing test may already
check it, and mapping one is the next step
Gates: no
Not watched failing.

## noIdentificationsCommentsOrMessages

Criterion: [birdsync.noIdentificationsCommentsOrMessages](requirements/composed/internal/birdsync/scope.md#noidentificationscommentsormessages)
Guardrails: none mapped yet. Written by composition on 2026-11-05; an existing test may already
check it, and mapping one is the next step
Gates: no
Not watched failing.

## mediaRequestsOnlyListedAssets

Criterion: [birdsync.mediaRequestsOnlyListedAssets](requirements/composed/internal/birdsync/media.md#mediarequestsonlylistedassets)
Guardrails: none mapped yet. Written by composition on 2026-11-05; an existing test may already
check it, and mapping one is the next step
Gates: no
Not watched failing.

## mediaAttachedToListingRecord

Criterion: [birdsync.mediaAttachedToListingRecord](requirements/composed/internal/birdsync/media.md#mediaattachedtolistingrecord)
Guardrails: none mapped yet. Written by composition on 2026-11-05; an existing test may already
check it, and mapping one is the next step
Gates: no
Not watched failing.

## assetNamedInDescription

Criterion: [birdsync.assetNamedInDescription](requirements/composed/internal/birdsync/media.md#assetnamedindescription)
Guardrails: none mapped yet. Written by composition on 2026-11-05; an existing test may already
check it, and mapping one is the next step
Gates: no
Not watched failing.

## attributionMeetsGuidelines

Criterion: [birdsync.attributionMeetsGuidelines](requirements/composed/internal/birdsync/media.md#attributionmeetsguidelines)
Guardrails: none mapped yet. Written by composition on 2026-11-05; an existing test may already
check it, and mapping one is the next step
Gates: no
Not watched failing.

## Project checks

Tests that check the repository rather than a requirement. These were specman criteria
`AC-039` and `AC-040`. River makes criteria only from requirements, so these are kept as the
project's own checks (River's checks.md: *Checks of your own*), not as criteria.

- `TestTranscribedQuotesAppearInSources` (`guard_test.go`), formerly AC-039: every quoted passage in an imported transcription appears in its vendored documents. 29 passages
- `TestTalkLinksResolve` (`guard_test.go`), formerly AC-040: every link from `talks/` into this repository resolves

## Harvested annotations

Until 2026-11-03, test functions cited the requirements they cover in a `// Verifies:` comment
above the function. River keeps requirement IDs out of source code, so every annotation was
copied here verbatim, checked against the source by script, and then removed. `Formerly:` in
each source item resolves the old IDs. No `// Criterion:` annotations existed.

| Test | File | Annotation | Criteria citing the test |
| --- | --- | --- | --- |
| `TestBirdsync` | `birdsync_test.go` | `Verifies: P-020, P-026, P-027, P-028, P-030, P-031.` | `skipOrderEnforced` |
| `TestUpdateMedia` | `birdsync_test.go` | `Verifies: P-047.` | `mediaResyncAdditive` |
| `TestDryRunIssuesNoWrites` | `birdsync_test.go` | `Verifies: P-051, T-005.` | `dryRunIssuesNoWrites` |
| `TestInvalidRecordsAreSkipped` | `birdsync_test.go` | `Verifies: P-062.` | `invalidRecordsSkipped` |
| `TestUntaxonedObservationIsRecognized` | `birdsync_test.go` | `Verifies: P-020, P-061.` | `untaxonedObservationRecognized` |
| `TestCreatedObservationContent` | `birdsync_test.go` | `Verifies: P-019, P-035, P-036, P-037, P-038, P-039, P-040.` | `createdObservationIdentity`, `createdObservationPlacement`, `createdObservationMedia` |
| `TestTempFilesAreCleanedUp` | `birdsync_test.go` | `Verifies: T-023.` | `tempFilesCleanedUp` |
| `TestSummaryDryRunLabels` | `birdsync_test.go` | `Verifies: P-060, T-007.` | `summaryLabels` |
| `TestSummaryRealRunLabels` | `birdsync_test.go` | `Verifies: P-054.` | `summaryLabels` |
| `TestSummaryConditionalLines` | `birdsync_test.go` | `Verifies: P-055, P-056.` | `summaryLabels` |
| `TestFailedUploadIsRetriedNextRun` | `birdsync_test.go` | `Verifies: P-040, P-047, P-050.` | `failedUploadRetried` |
| `TestAllUploadsFailingLeavesNoUpdate` | `birdsync_test.go` | `Verifies: P-040, T-007.` | `failedUploadRetried` |
| `TestPermanentUploadFailureIsNotRetried` | `birdsync_test.go` | `Verifies: P-063, P-064.` | `permanentVersusTransientUpload` |
| `TestTransientUploadFailureIsRetried` | `birdsync_test.go` | `Verifies: P-063.` | `permanentVersusTransientUpload` |
| `TestDownloadMLAsset_Sound` | `ebird/ebird_test.go` | `Verifies: P-043, P-044.` | `mlAssetDownload`, `canonicalExtensions` |
| `TestDownloadMLAsset_Photo` | `ebird/ebird_test.go` | `Verifies: P-043, P-044.` | `mlAssetDownload`, `canonicalExtensions` |
| `TestRecord_Observed` | `ebird/ebird_test.go` | `Verifies: T-019.` | `recordObservedParsing` |
| `TestRecords` | `ebird/ebird_test.go` | `Verifies: T-018.` | `csvReadByHeader` |
| `TestObservationID_Valid` | `ebird/ebird_test.go` | `Verifies: P-019, P-022.` | `syncKeyValidity` |
| `TestFileExtension` | `ebird/ebird_test.go` | `Verifies: P-045, T-004.` | `canonicalExtensions` |
| `TestRecordsRejectsBadInput` | `ebird/ebird_test.go` | `Verifies: P-066.` | `badInputRejected` |
| `TestDownloadMLAssetCleansUpOnError` | `ebird/ebird_test.go` | `Verifies: T-023.` | `downloadCleansUpOnError` |
| `TestNoLiveServiceHostsInTests` | `guard_test.go` | `Verifies: T-010.` | `noLiveHostsInTests` |
| `TestToolsAreReadOnly` | `guard_test.go` | `Verifies: T-032.` | `toolsAreReadOnly` |
| `TestNoLogFatalInLibraryPackages` | `guard_test.go` | `Verifies: T-027.` | `noLogFatalOutsideMain` |
| `TestTranscribedQuotesAppearInSources` | `guard_test.go` | `Verifies: the accuracy of every `<source>/R#` transcription.` | none |
| `TestAmericanSpellings` | `guard_test.go` | `Verifies: T-037, T-038.` | `americanSpellings` |
| `TestClient_CreateObservation` | `inat/client_test.go` | `Verifies: T-016.` | `userAgentSent` |
| `TestClient_UpdateObservation` | `inat/client_test.go` | `Verifies: T-009, P-046.` | `updateKeepsPhotos` |
| `TestStatusErrorIncludesBody` | `inat/client_test.go` | `Verifies: T-034.` | `statusErrorCarriesBody` |
| `TestStatusErrorPermanence` | `inat/client_test.go` | `Verifies: P-063.` | `permanentVersusTransientUpload` |
| `TestStatusErrorDropsHTMLBody` | `inat/client_test.go` | `Verifies: T-034.` | `statusErrorCarriesBody` |
| `TestStatusErrorCollapsesWhitespace` | `inat/client_test.go` | `Verifies: T-034.` | `statusErrorCarriesBody` |
| `TestClientPacesRequests` | `inat/client_test.go` | `Verifies: T-035.` | `requestsPaced` |
| `TestDefaultMinRequestIntervalMatchesGuidance` | `inat/client_test.go` | `Verifies: T-035.` | `requestsPaced` |
| `TestDownloadObservations` | `inat/inat_test.go` | `Verifies: T-017, P-023.` | `downloadUsesMaxPageSize` |
| `TestDownloadObservationsNoTaxonFilter` | `inat/inat_test.go` | `Verifies: P-061.` | `untaxonedObservationRecognized` |
| `TestDownloadObservationsDateWindow` | `inat/inat_test.go` | `Verifies: P-025.` | `downloadHonorsDateWindow` |
| `TestDownloadObservationsNoDateWindow` | `inat/inat_test.go` | `Verifies: P-025.` | `downloadHonorsDateWindow` |
| `TestDownloadObservationsError` | `inat/inat_test.go` | `Verifies: T-027.` | none |
| `TestDownloadObservationsPagesByID` | `inat/inat_test.go` | `Verifies: P-065, T-036.` | `downloadPagesByID` |
| `TestDownloadObservationsCursorMustAdvance` | `inat/inat_test.go` | `Verifies: T-036.` | `downloadPagesByID` |
| `TestMediaChange` | `media_test.go` | `Verifies: P-048, P-049.` | `mediaChangesReported`, `failedAssetsReported` |

## Tests no criterion cites

Kept and counted, never removed for lacking a criterion. Each may be an undocumented
requirement, a test of an implementation detail, or a test of behavior nobody wants any more.

- `TestClient_DeleteObservation` (`inat/client_test.go`).
- `TestDownloadObservationsError` (`inat/inat_test.go`). Annotated `Verifies: T-027.`.

## Requirements not met

Requirements the code is known not to satisfy, none of them yet an accepted gap. Each has a
drafted decision for the owner in [`proposals.md`](proposals.md).

- [birdsync.macaulayAssetsAttributed](requirements/composed/internal/birdsync/media.md#macaulayassetsattributed): not known to be met; the attribution guidelines it depends on were never retrieved

[inatapi.authenticateOnlyWhenNeeded](requirements/imported/inaturalist.org/inatapi/inatapi.md#authenticateonlywhenneeded) was listed here until 2026-11-04, when the owner decided that
every request birdsync makes needs authentication, so the rule is met
([decisions.authenticationIsNecessary](requirements/decisions.md#authenticationisnecessary)).
[birdsync.requestRate](requirements/composed/internal/birdsync/download.md#requestrate) was listed here until 2026-11-08, when the 9,000-request 24-hour budget was implemented
([decisions.dailyRequestBudget](requirements/decisions.md#dailyrequestbudget)).
