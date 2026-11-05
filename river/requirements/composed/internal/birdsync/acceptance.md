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

Limit: runs every criterion below.

## noLiveHostsInTests

**No test names a live service host**

Formerly: AC-005
Covers: [tech.noLiveServicesInTests](../../../local/tech/testing.md#noliveservicesintests)

No `_test.go` file contains the iNaturalist API host or the Macaulay Library CDN host.

Limit: a substring scan; a host assembled from parts escapes it.

## dryRunIssuesNoWrites

**A dry run issues no writes**

Formerly: AC-006
Covers: [product.dryRunNoWrites](../../../local/product/dryrun.md#dryrunnowrites), [tech.dryRunGateAtCallSite](../../../local/tech/safety.md#dryrungateatcallsite), [tech.clientInterfaces](../../../local/tech/testing.md#clientinterfaces)

A `--dryrun` sync through a recording fake makes no create, update, or upload call, while reads still happen.

## createdObservationIdentity

**A created observation carries the sync key and observation fields**

Formerly: AC-007 (part)
Covers: [product.syncKey](../../../local/product/identity.md#synckey), [product.observationFieldMapping](../../../local/product/observations.md#observationfieldmapping), [tech.clientInterfaces](../../../local/tech/testing.md#clientinterfaces)

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
Covers: [tech.updateKeepsPhotos](../../../local/tech/safety.md#updatekeepsphotos), [product.descriptionUpdateKeepsMedia](../../../local/product/media.md#descriptionupdatekeepsmedia)

`UpdateObservation` against an `httptest` server sends `ignore_photos` set, and the updated description.

## userAgentSent

**API requests carry the birdsync User-Agent**

Formerly: AC-009
Covers: [inatUserAgent](birdsync.md#inatuseragent), [tech.userAgent](../../../local/tech/etiquette.md#useragent)

A request from `inat.Client` carries the birdsync `User-Agent`.

## downloadUsesMaxPageSize

**The download requests the maximum page size**

Formerly: AC-010
Covers: [observationDownloadPageSize](birdsync.md#observationdownloadpagesize), [product.downloadExistingObservations](../../../local/product/download.md#downloadexistingobservations), [tech.baseURLSeams](../../../local/tech/testing.md#baseurlseams)

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
Covers: [product.fuzzyParsedDates](../../../local/product/filtering.md#fuzzyparseddates), [tech.datesViaObserved](../../../local/tech/formats.md#datesviaobserved), [product.fuzzyMatching](../../../local/product/filtering.md#fuzzymatching)

Fuzzy matching skips a record whose date is written in either eBird format.

## fuzzyIgnoresEmptyNames

**Fuzzy matching ignores empty names**

Formerly: AC-015
Covers: [product.fuzzyNeverEmptyName](../../../local/product/filtering.md#fuzzyneveremptyname), [tech.emptyNameNotIndexed](../../../local/tech/formats.md#emptynamenotindexed)

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
Covers: [product.mediaFormats](../../../local/product/media.md#mediaformats), [product.photoThenSoundFallback](../../../local/product/media.md#photothensoundfallback), [tech.assetTypeNeedsDownload](../../../local/tech/formats.md#assettypeneedsdownload), [tech.baseURLSeams](../../../local/tech/testing.md#baseurlseams)

Against an `httptest` server, an asset downloads from the 2400px photo URL, or from the MP3 URL when the photo URL returns 404.

## recordObservedParsing

**eBird dates parse in both formats**

Formerly: AC-019
Covers: [tech.datesViaObserved](../../../local/tech/formats.md#datesviaobserved)

`Record.Observed()` parses both eBird date formats, with and without a time.

## csvReadByHeader

**The CSV is read by header name**

Formerly: AC-020
Covers: [tech.csvByHeaderName](../../../local/tech/formats.md#csvbyheadername), [product.csvExportInput](../../../local/product/scope.md#csvexportinput)

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
Covers: [tech.ciStandingChecks](../../../local/tech/ci.md#cistandingchecks), [tech.ciLatestStable](../../../local/tech/ci.md#cilateststable)

`.github/workflows/ci.yml` runs build, vet, gofmt, and the race-enabled suite on every push to main and every pull request, plus a non-blocking latest-stable job.

## summaryLabels

**The summary labels and conditional lines are right**

Formerly: AC-023
Covers: [product.dryRunCounterLabels](../../../local/product/reporting.md#dryruncounterlabels), [tech.honestDryRunCounters](../../../local/tech/safety.md#honestdryruncounters), [product.endOfRunSummary](../../../local/product/reporting.md#endofrunsummary), [product.conditionalSkipCounters](../../../local/product/reporting.md#conditionalskipcounters), [product.conditionalFailureCount](../../../local/product/reporting.md#conditionalfailurecount)

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
Covers: [tech.tempFilesDeleted](../../../local/tech/resources.md#tempfilesdeleted)

After a sync no downloaded media file remains in the temp directory.

## noLogFatalOutsideMain

**Library packages never exit**

Formerly: AC-027
Covers: [tech.noLogFatalInLibraries](../../../local/tech/style.md#nologfatalinlibraries)

No non-main package in the module calls `log.Fatal` or `os.Exit`.

## toolsAreReadOnly

**tools/ is read-only**

Formerly: AC-028
Covers: [tech.toolsReadOnly](../../../local/tech/safety.md#toolsreadonly), [product.toolsNotProduct](../../../local/product/scope.md#toolsnotproduct)

No program under `tools/` calls a mutating `inat.Client` method.

## canonicalExtensions

**File extensions are canonical on every platform**

Formerly: AC-029
Covers: [product.mediaFormats](../../../local/product/media.md#mediaformats), [product.mediaFilename](../../../local/product/media.md#mediafilename), [tech.platformIndependent](../../../local/tech/deps.md#platformindependent)

Known content types, including the CDN's `audio/mpeg3`, map to `.jpg`, `.png`, `.mp3`, or `.wav`, and an unknown type falls back by endpoint.

## failedUploadRetried

**A failed upload is retried on the next run**

Formerly: AC-030
Covers: [product.descriptionContents](../../../local/product/observations.md#descriptioncontents), [product.additiveMediaResync](../../../local/product/media.md#additivemediaresync), [product.mediaFailuresTolerated](../../../local/product/media.md#mediafailurestolerated), [tech.honestDryRunCounters](../../../local/tech/safety.md#honestdryruncounters)

Across two runs, an asset whose upload failed is absent from the first description and uploaded on the second. If every upload fails, no update is made or counted.

## permanentVersusTransientUpload

**Permanent and transient upload failures are told apart**

Formerly: AC-031
Covers: [product.permanentUploadFailureRecorded](../../../local/product/failures.md#permanentuploadfailurerecorded), [product.permanentFailuresReported](../../../local/product/failures.md#permanentfailuresreported)

A permanently refused upload is recorded and not retried, a transient one is retried, and the status-code classification holds in both directions.

## statusErrorCarriesBody

**API errors carry the response body**

Formerly: AC-032
Covers: [tech.errorsIncludeBody](../../../local/tech/etiquette.md#errorsincludebody)

A non-200 response's error includes the body, truncated and on one line, and drops an HTML body.

## failedAssetsReported

**Permanently failed assets are reported each run**

Formerly: AC-033
Covers: [product.permanentFailuresReported](../../../local/product/failures.md#permanentfailuresreported)

An asset recorded as permanently failed is reported and not counted as attached or mismatched.

## downloadPagesByID

**The download pages by id_above past 10,000 results**

Formerly: AC-034
Covers: [product.largeAccounts](../../../local/product/download.md#largeaccounts), [observationDownloadByIDAbove](birdsync.md#observationdownloadbyidabove)

Against a fake honoring `per_page`, the download pages by `id_above` and errors when the cursor fails to advance.

## requestsPaced

**Requests are paced**

Formerly: AC-035
Covers: [requestRate](birdsync.md#requestrate)

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
Covers: [accountHolderToldOfObligations](birdsync.md#accountholdertoldofobligations), [product.documentSyncedObservationsIdentifiable](../../../local/product/community.md#documentsyncedobservationsidentifiable)

A person reviewing the README finds what iNaturalist expects of a bulk-syncing account (that the user answers for what is posted, reviews what was created, and responds to comments and messages), and how synced observations are identifiable.

Limit: human review.

## downloadCleansUpOnError

**A failed download leaves no temp file**

Formerly: AC-038
Covers: [tech.tempFilesDeleted](../../../local/tech/resources.md#tempfilesdeleted)

A truncated download leaves nothing in the temp directory.

## americanSpellings

**Prose is American English and quotations are untouched**

Formerly: AC-041
Covers: [tech.americanSpelling](../../../local/tech/style.md#americanspelling), [tech.quotationsNeverRespelled](../../../local/tech/style.md#quotationsneverrespelled)

No Markdown or Go file outside vendored documents uses a listed British spelling, and blockquotes and vendored documents are never flagged.

## pacingPerDay

**No more than 10,000 iNaturalist requests in a day**

Covers: [requestRate](birdsync.md#requestrate)

Given an export whose sync would need more than 10,000 requests, with a clock birdsync takes
from the test, the test server receives no more than 10,000 requests in any 24-hour span.

Currently expected to fail: [`imports.md`](../../../imports.md) records the daily half as
unenforced. Whether the span is within a run or across runs is open; the criterion is written
for within a run, and does not cover requests across separate runs until the owner decides.

## downloadOrderedByID

**The download sorts by id ascending and never uses page**

Covers: [observationDownloadByIDAbove](birdsync.md#observationdownloadbyidabove)

Every request of the existing-observation download carries `order_by=id` and `order=asc`, and
none carries a `page` parameter.

## downloadStopsOnShortPage

**The download ends at the first short page**

Covers: [observationDownloadByIDAbove](birdsync.md#observationdownloadbyidabove)

When the test server answers a download request with fewer observations than the page size,
birdsync sends no further download request.

## jwtFromEnvironmentInHeader

**A token from INAT_API_TOKEN is sent in the Authorization header**

Covers: [jwtInAuthorizationHeader](birdsync.md#jwtinauthorizationheader)

With `INAT_API_TOKEN` set to a token, every authenticated request the test server receives
carries that token in its `Authorization` header.

## jwtFromPromptUnwrapped

**A token entered at the prompt is sent unwrapped in the Authorization header**

Covers: [jwtInAuthorizationHeader](birdsync.md#jwtinauthorizationheader)

With `INAT_API_TOKEN` unset and `{"api_token":"TOKEN"}` entered at the prompt, every
authenticated request the test server receives carries `TOKEN`, without the JSON framing, in its
`Authorization` header.

## noBuiltInCredentials

**Without the user's token, birdsync writes nothing**

Covers: [ownAccountOnly](birdsync.md#ownaccountonly)

With `INAT_USER_ID` and `INAT_API_TOKEN` unset and the prompt given no token, the test server
receives no write request and no request carrying an `Authorization` header.

## writesCarrySuppliedToken

**Every write carries the token the user supplied**

Covers: [ownAccountOnly](birdsync.md#ownaccountonly)

Every write request the test server receives during a run carries, in its `Authorization`
header, exactly the token supplied for that run, and no other.

## tokenOnlyWhereNeeded

**Only writes and the existing-observation download carry the token**

Covers: [authenticateOnlyWhenNecessary](birdsync.md#authenticateonlywhennecessary)

Of the requests the test server receives during a run, every one carrying an `Authorization`
header is a write or a request of the existing-observation download; any other request carries
none. A new iNaturalist request added to birdsync falls under this check by being made, not by
being listed.

## authenticatedDownloadSeesPrivate

**The authenticated download sees private and obscured observations, so a rerun duplicates none**

Covers: [authenticateOnlyWhenNecessary](birdsync.md#authenticateonlywhennecessary)

Given a test server that returns the user's private and geoprivacy-obscured birdsync-created
observations only to an authenticated request, a second run over the same export creates no
observation.

## createsMatchExportRecords

**Every created observation corresponds to one export record**

Covers: [postsOnlyUsersOwnRecords](birdsync.md#postsonlyusersownrecords)

Over a run, every observation create the test server receives carries a sync key
([product.syncKey](../../../local/product/identity.md#synckey)) that matches one record of the
input export, and no two creates match the same record.

## noIdentificationsCommentsOrMessages

**birdsync posts no identification, comment, or message**

Covers: [postsOnlyUsersOwnRecords](birdsync.md#postsonlyusersownrecords)

Over a run that creates, uploads, and updates, every write request the test server receives
targets an observation or its media; none targets identifications, comments, or messages.

## mediaRequestsOnlyListedAssets

**Every Macaulay Library request is for an asset the export lists**

Covers: [mediaOnlyFromUsersExport](birdsync.md#mediaonlyfromusersexport)

Over a run, every request a test server standing in for the Macaulay Library CDN receives is
for an asset whose catalog number appears in a record of the input export that the run syncs,
and no request is a search.

Does not cover whether the listed assets are the user's own; that rests with the account holder
([decisions.inputIdentityIsTheAccountHolders](../../../decisions.md#inputidentityistheaccountholders)).

## mediaAttachedToListingRecord

**Each asset is attached to the observation of the record that lists it**

Covers: [mediaOnlyFromUsersExport](birdsync.md#mediaonlyfromusersexport)

Given an export in which two records list different assets, each upload the test server
receives is attached to the observation created or found for the record that lists that asset.

## assetNamedInDescription

**Each uploaded asset is named on its observation**

Covers: [macaulayAssetsAttributed](birdsync.md#macaulayassetsattributed)

After a run, the description of every observation to which an asset was uploaded contains a
`Macaulay Library Asset:` line linking `https://macaulaylibrary.org/asset/<id>` for that asset,
and the file was uploaded under the name `ML<id>` with its extension.

Does not decide [macaulayAssetsAttributed](birdsync.md#macaulayassetsattributed): it checks what
birdsync does, not whether that meets the Macaulay Library's attribution guidelines.

## attributionMeetsGuidelines

**birdsync's attribution meets the Macaulay Library's guidelines**

Covers: [macaulayAssetsAttributed](birdsync.md#macaulayassetsattributed)

A reviewer, the owner, compares the attribution birdsync writes with the Macaulay Library's media
attribution guidelines, as vendored in `mlterms`, and confirms it meets them; repeated at every
refresh of `mlterms`.

Cannot be performed yet: the guidelines page has not been retrieved. Until it is, the
requirement counts as uncovered.
