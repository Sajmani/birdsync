package birdsync

# birdsync: acceptance criteria

**Composed, not authored**, with [`birdsync.md`](birdsync.md), and changed only by recomposing.
One criterion per independently checkable behavior. Each cites what it covers in `Covers:`.
Which guardrail implements a criterion, and whether it has been watched failing, is
recorded in [`conformance.md`](../../../../conformance.md), not here.

Migrated from specman's `spec/acceptance.md` (`AC-001` to `AC-041`; `Formerly:` gives the
old ID). `AC-039` and `AC-040` check the repository rather than any requirement, so they
are not criteria here: conformance.md lists them as project checks.

## buildsEverywhere

**Everything builds**

Formerly: AC-001
Covers: [tech.singleModule](../../../local/tech/tech.md#singlemodule), [tech.platformIndependent](../../../local/tech/tech.md#platformindependent), [product.toolsNotProduct](../../../local/product/scope.md#toolsnotproduct)

`go build ./...` compiles every package, including `tools/`.

Limit: partial for platform independence: one OS.

## vetPasses

**go vet is clean**

Formerly: AC-002
Covers: [tech.vetClean](../../../local/tech/tech.md#vetclean)

`go vet ./...` reports nothing.

## formatted

**Code is gofmt'd**

Formerly: AC-003
Covers: [tech.gofmtClean](../../../local/tech/tech.md#gofmtclean)

`gofmt -l .` prints nothing.

## suitePassesUnderRace

**The suite passes under the race detector**

Formerly: AC-004
Covers: [tech.ciStandingChecks](../../../local/tech/tech.md#cistandingchecks)

`go test -race ./...` passes with no network and no live accounts.

Limit: runs every criterion below.

## noLiveHostsInTests

**No test names a live service host**

Formerly: AC-005
Covers: [tech.noLiveServicesInTests](../../../local/tech/tech.md#noliveservicesintests)

No `_test.go` file contains the iNaturalist API host or the Macaulay Library CDN host.

Limit: a substring scan; a host assembled from parts escapes it.

## dryRunIssuesNoWrites

**A dry run issues no writes**

Formerly: AC-006
Covers: [product.dryRunNoWrites](../../../local/product/dryrun.md#dryrunnowrites), [tech.dryRunGateAtCallSite](../../../local/tech/tech.md#dryrungateatcallsite), [tech.clientInterfaces](../../../local/tech/tech.md#clientinterfaces)

A `--dryrun` sync through a recording fake makes no create, update, or upload call, while reads still happen.

## createdObservationIdentity

**A created observation carries the sync key and observation fields**

Formerly: AC-007 (part)
Covers: [product.syncKey](../../../local/product/identity.md#synckey), [product.observationFieldMapping](../../../local/product/observations.md#observationfieldmapping), [tech.clientInterfaces](../../../local/tech/tech.md#clientinterfaces)

A sync through a recording fake creates an observation whose sync-key fields and copied eBird observation fields have the required IDs and values.

## createdObservationPlacement

**A created observation is wild, placed, and timed as required**

Formerly: AC-007 (part)
Covers: [product.createdObservationWild](../../../local/product/observations.md#createdobservationwild), [product.checklistCoordinates](../../../local/product/observations.md#checklistcoordinates), [product.speciesGuessScientificName](../../../local/product/observations.md#speciesguessscientificname), [product.observedDateTime](../../../local/product/observations.md#observeddatetime)

The created observation is marked wild, carries the checklist's coordinates marked inexact with the configured accuracy, the scientific name as species guess, and the eBird date and time.

Limit: asserts what is sent, not how iNaturalist interprets it.

## createdObservationMedia

**Media follows creation and the description records it**

Formerly: AC-007 (part)
Covers: [product.descriptionContents](../../../local/product/observations.md#descriptioncontents), [product.createBeforeMedia](../../../local/product/media.md#createbeforemedia), [product.mediaFilename](../../../local/product/media.md#mediafilename), [product.descriptionUpdateKeepsMedia](../../../local/product/media.md#descriptionupdatekeepsmedia)

The created description carries the attribution, details, checklist URL, and protocol but no asset lines; the later update adds the asset URLs, and the upload is keyed by the asset ID.

## updateKeepsPhotos

**Updating an observation keeps its photos**

Formerly: AC-008
Covers: [tech.updateKeepsPhotos](../../../local/tech/tech.md#updatekeepsphotos), [product.descriptionUpdateKeepsMedia](../../../local/product/media.md#descriptionupdatekeepsmedia)

`UpdateObservation` against an `httptest` server sends `ignore_photos` set, and the updated description.

## userAgentSent

**API requests carry the birdsync User-Agent**

Formerly: AC-009
Covers: [tech.userAgent](../../../local/tech/tech.md#useragent)

A request from `inat.Client` carries the birdsync `User-Agent`.

## downloadUsesMaxPageSize

**The download requests the maximum page size**

Formerly: AC-010
Covers: [tech.maxPageSize](../../../local/tech/tech.md#maxpagesize), [product.downloadExistingObservations](../../../local/product/download.md#downloadexistingobservations), [tech.baseURLSeams](../../../local/tech/tech.md#baseurlseams)

`DownloadObservations` against an `httptest` server requests `per_page=200` and returns every observation served.

## downloadHonorsDateWindow

**The download honors the date window**

Formerly: AC-011
Covers: [product.dateWindowNarrowsDownload](../../../local/product/download.md#datewindownarrowsdownload)

With `--after` or `--before` set the download query carries the window, and without them it carries none.

## skipOrderEnforced

**Each skip rule fires in order**

Formerly: AC-012
Covers: [product.idempotentRerun](../../../local/product/identity.md#idempotentrerun), [product.skipOrder](../../../local/product/filtering.md#skiporder), [product.afterFlag](../../../local/product/filtering.md#afterflag), [product.beforeFlag](../../../local/product/filtering.md#beforeflag), [product.verifiableDefaultOn](../../../local/product/filtering.md#verifiabledefaulton), [product.fuzzyMatching](../../../local/product/filtering.md#fuzzymatching)

One set of records exercises every skip rule exactly once, in the order `--after`, `--before`, already synced, `--fuzzy`, `--verifiable`, and the rest are created.

## mediaResyncAdditive

**Media re-sync is additive**

Formerly: AC-013
Covers: [product.additiveMediaResync](../../../local/product/media.md#additivemediaresync)

Assets added to an already-synced checklist are uploaded on the next run.

## fuzzyHandlesBothDateFormats

**Fuzzy matching works with both eBird date formats**

Formerly: AC-014
Covers: [product.fuzzyParsedDates](../../../local/product/filtering.md#fuzzyparseddates), [tech.datesViaObserved](../../../local/tech/tech.md#datesviaobserved), [product.fuzzyMatching](../../../local/product/filtering.md#fuzzymatching)

Fuzzy matching skips a record whose date is written in either eBird format.

## fuzzyIgnoresEmptyNames

**Fuzzy matching ignores empty names**

Formerly: AC-015
Covers: [product.fuzzyNeverEmptyName](../../../local/product/filtering.md#fuzzyneveremptyname), [tech.emptyNameNotIndexed](../../../local/tech/tech.md#emptynamenotindexed)

An existing observation with an empty name never fuzzy-matches a record.

## dryRunMediaCountUnclassified

**A dry run reports one unclassified media count**

Formerly: AC-016
Covers: [product.dryRunMediaCount](../../../local/product/dryrun.md#dryrunmediacount)

Under `--dryrun` the media count is a single unclassified number.

Limit: does not detect a broken dry-run gate; criterion 6 does.

## mediaChangesReported

**Media differences are reported, not applied**

Formerly: AC-017
Covers: [product.removedAssetsReported](../../../local/product/media.md#removedassetsreported), [product.mediaCountMismatchReported](../../../local/product/media.md#mediacountmismatchreported)

Assets removed from eBird, and count mismatches, are reported and not corrected.

## mlAssetDownload

**Macaulay Library assets download as photo or sound**

Formerly: AC-018
Covers: [product.mediaFormats](../../../local/product/media.md#mediaformats), [product.photoThenSoundFallback](../../../local/product/media.md#photothensoundfallback), [tech.assetTypeNeedsDownload](../../../local/tech/tech.md#assettypeneedsdownload), [tech.baseURLSeams](../../../local/tech/tech.md#baseurlseams)

Against an `httptest` server, an asset downloads from the 2400px photo URL, or from the MP3 URL when the photo URL returns 404.

## recordObservedParsing

**eBird dates parse in both formats**

Formerly: AC-019
Covers: [tech.datesViaObserved](../../../local/tech/tech.md#datesviaobserved)

`Record.Observed()` parses both eBird date formats, with and without a time.

## csvReadByHeader

**The CSV is read by header name**

Formerly: AC-020
Covers: [tech.csvByHeaderName](../../../local/tech/tech.md#csvbyheadername), [product.csvExportInput](../../../local/product/scope.md#csvexportinput)

`ebird.Records` reads fields by header name from a file whose columns vary.

Limit: partial for the CSV input requirement.

## syncKeyValidity

**The sync key is valid only when complete**

Formerly: AC-021
Covers: [product.syncKey](../../../local/product/identity.md#synckey), [product.incompleteKeyNotRecognized](../../../local/product/identity.md#incompletekeynotrecognized), [product.taxonNotInSyncKey](../../../local/product/identity.md#taxonnotinsynckey)

`ObservationID.Valid` is true only when both sync-key fields are present.

Limit: partial for the taxon not being in the key.

## ciRunsStandingChecks

**CI runs the standing checks on push and pull request**

Formerly: AC-022
Covers: [tech.ciStandingChecks](../../../local/tech/tech.md#cistandingchecks), [tech.ciLatestStable](../../../local/tech/tech.md#cilateststable)

`.github/workflows/ci.yml` runs build, vet, gofmt, and the race-enabled suite on every push to main and every pull request, plus a non-blocking latest-stable job.

## summaryLabels

**The summary labels and conditional lines are right**

Formerly: AC-023
Covers: [product.dryRunCounterLabels](../../../local/product/reporting.md#dryruncounterlabels), [tech.honestDryRunCounters](../../../local/tech/tech.md#honestdryruncounters), [product.endOfRunSummary](../../../local/product/reporting.md#endofrunsummary), [product.conditionalSkipCounters](../../../local/product/reporting.md#conditionalskipcounters), [product.conditionalFailureCount](../../../local/product/reporting.md#conditionalfailurecount)

The summary says `Would create` and `Would update` under `--dryrun`, and the real labels otherwise. A counter line appears only when its flag is in effect, and the failure line only when non-zero.

## untaxonedObservationRecognized

**A previously synced untaxoned observation is recognized**

Formerly: AC-024
Covers: [product.noTaxonFilter](../../../local/product/download.md#notaxonfilter), [product.idempotentRerun](../../../local/product/identity.md#idempotentrerun)

The download carries no taxon filter, and an observation with no taxon that carries a sync key is recognized as synced and not re-created.

## invalidRecordsSkipped

**Unparseable records are skipped and counted**

Formerly: AC-025
Covers: [product.badRecordsSkipped](../../../local/product/failures.md#badrecordsskipped)

Records with an unparseable date, time, latitude, or longitude are skipped, logged with their line, and counted, and the run continues.

## tempFilesCleanedUp

**Temporary media files are removed after upload**

Formerly: AC-026
Covers: [tech.tempFilesDeleted](../../../local/tech/tech.md#tempfilesdeleted)

After a sync no downloaded media file remains in the temp directory.

## noLogFatalOutsideMain

**Library packages never exit**

Formerly: AC-027
Covers: [tech.noLogFatalInLibraries](../../../local/tech/tech.md#nologfatalinlibraries)

No non-main package in the module calls `log.Fatal` or `os.Exit`.

## toolsAreReadOnly

**tools/ is read-only**

Formerly: AC-028
Covers: [tech.toolsReadOnly](../../../local/tech/tech.md#toolsreadonly), [product.toolsNotProduct](../../../local/product/scope.md#toolsnotproduct)

No program under `tools/` calls a mutating `inat.Client` method.

## canonicalExtensions

**File extensions are canonical on every platform**

Formerly: AC-029
Covers: [product.mediaFormats](../../../local/product/media.md#mediaformats), [product.mediaFilename](../../../local/product/media.md#mediafilename), [tech.platformIndependent](../../../local/tech/tech.md#platformindependent)

Known content types, including the CDN's `audio/mpeg3`, map to `.jpg`, `.png`, `.mp3`, or `.wav`, and an unknown type falls back by endpoint.

## failedUploadRetried

**A failed upload is retried on the next run**

Formerly: AC-030
Covers: [product.descriptionContents](../../../local/product/observations.md#descriptioncontents), [product.additiveMediaResync](../../../local/product/media.md#additivemediaresync), [product.mediaFailuresTolerated](../../../local/product/media.md#mediafailurestolerated), [tech.honestDryRunCounters](../../../local/tech/tech.md#honestdryruncounters)

Across two runs, an asset whose upload failed is absent from the first description and uploaded on the second. If every upload fails, no update is made or counted.

## permanentVersusTransientUpload

**Permanent and transient upload failures are told apart**

Formerly: AC-031
Covers: [product.permanentUploadFailureRecorded](../../../local/product/failures.md#permanentuploadfailurerecorded), [product.permanentFailuresReported](../../../local/product/failures.md#permanentfailuresreported)

A permanently refused upload is recorded and not retried, a transient one is retried, and the status-code classification holds in both directions.

## statusErrorCarriesBody

**API errors carry the response body**

Formerly: AC-032
Covers: [tech.errorsIncludeBody](../../../local/tech/tech.md#errorsincludebody)

A non-200 response's error includes the body, truncated and on one line, and drops an HTML body.

## failedAssetsReported

**Permanently failed assets are reported each run**

Formerly: AC-033
Covers: [product.permanentFailuresReported](../../../local/product/failures.md#permanentfailuresreported)

An asset recorded as permanently failed is reported and not counted as attached or mismatched.

## downloadPagesByID

**The download pages by id_above past 10,000 results**

Formerly: AC-034
Covers: [product.largeAccounts](../../../local/product/download.md#largeaccounts), [tech.idAbovePaging](../../../local/tech/tech.md#idabovepaging)

Against a fake honoring `per_page`, the download pages by `id_above` and errors when the cursor fails to advance.

## requestsPaced

**Requests are paced**

Formerly: AC-035
Covers: [tech.requestPacing](../../../local/tech/tech.md#requestpacing)

Consecutive requests from one client are spaced by at least the minimum interval, and the default interval matches one per second.

Limit: the daily cap is not checked.

## badInputRejected

**Unusable input is named**

Formerly: AC-036
Covers: [product.unusableInputExplained](../../../local/product/failures.md#unusableinputexplained)

A directory, a zip archive, and a CSV without a `Submission ID` column are each rejected with a message naming the problem.

## readmeCommunitySection

**The README states the community obligations**

Formerly: AC-037
Covers: [product.documentCommunityObligations](../../../local/product/community.md#documentcommunityobligations), [product.documentSyncedObservationsIdentifiable](../../../local/product/community.md#documentsyncedobservationsidentifiable)

A person reviewing the README finds what iNaturalist expects of a bulk-syncing account, and how synced observations are identifiable.

Limit: human review.

## downloadCleansUpOnError

**A failed download leaves no temp file**

Formerly: AC-038
Covers: [tech.tempFilesDeleted](../../../local/tech/tech.md#tempfilesdeleted)

A truncated download leaves nothing in the temp directory.

## americanSpellings

**Prose is American English and quotations are untouched**

Formerly: AC-041
Covers: [tech.americanSpelling](../../../local/tech/tech.md#americanspelling), [tech.quotationsNeverRespelled](../../../local/tech/tech.md#quotationsneverrespelled)

No Markdown or Go file outside vendored documents uses a listed British spelling, and blockquotes and vendored documents are never flagged.
