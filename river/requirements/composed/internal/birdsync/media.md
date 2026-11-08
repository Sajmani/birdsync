package birdsync

# Media downloads, uploads, attribution, and resource cleanup

Macaulay Library asset downloads, iNaturalist media uploads, attribution, and temporary-file and
memory bounds for [`birdsync`](README.md).

## Included requirements

### Media handling

- [product.createBeforeMedia](../../../local/product/media.md#createbeforemedia) — The observation is created before media is attached
- [product.mediaFormats](../../../local/product/media.md#mediaformats) — Photos at 2400px, sounds as MP3
- [product.photoThenSoundFallback](../../../local/product/media.md#photothensoundfallback) — Try the photo URL, then the sound URL
- [product.mediaFilename](../../../local/product/media.md#mediafilename) — Media is uploaded as ML<asset ID> with a canonical extension
- [product.descriptionUpdateKeepsMedia](../../../local/product/media.md#descriptionupdatekeepsmedia) — The description update keeps attached media
- [product.additiveMediaResync](../../../local/product/media.md#additivemediaresync) — Media re-sync is additive
- [product.removedAssetsReported](../../../local/product/media.md#removedassetsreported) — Assets removed from eBird are reported, not removed
- [product.mediaCountMismatchReported](../../../local/product/media.md#mediacountmismatchreported) — A media count mismatch is reported, not corrected
- [product.mediaFailuresTolerated](../../../local/product/media.md#mediafailurestolerated) — A failed media transfer is logged, counted, and survived
- [tech.updateKeepsPhotos](../../../local/tech/safety.md#updatekeepsphotos) — UpdateObservation always sets ignore_photos
- [tech.assetTypeNeedsDownload](../../../local/tech/formats.md#assettypeneedsdownload) — Asset type is known only after download
- [mlterms.contributorKeepsCopyright](../../../imported/ebird.org/mlterms/mlterms.md#contributorkeepscopyright) — A contributor keeps copyright in their own media

### Resource use

- [tech.memoryScalesWithExport](../../../local/tech/resources.md#memoryscaleswithexport) — Memory scales with the export
- [tech.tempFilesDeleted](../../../local/tech/resources.md#tempfilesdeleted) — Temporary media files are deleted

## Composed requirements

## mediaOnlyFromUsersExport

**birdsync downloads and re-uploads only the Macaulay Library assets the user's export lists**

Composed from: [mlterms.downloadOwnMedia](../../../imported/ebird.org/mlterms/mlterms.md#downloadownmedia),
[mlterms.noOthersMediaWithoutPermission](../../../imported/ebird.org/mlterms/mlterms.md#noothersmediawithoutpermission),
[mlterms.mediaOwnedByUploader](../../../imported/ebird.org/mlterms/mlterms.md#mediaownedbyuploader),
[inatterms.accountHolderResponsible](../../../imported/inaturalist.org/inatterms/inatterms.md#accountholderresponsible),
[decisions.ownMediaAssumption](../../../decisions.md#ownmediaassumption),
[decisions.inputIdentityIsTheAccountHolders](../../../decisions.md#inputidentityistheaccountholders)
Goals: [goals.accountsStaySafe](../../../../goals/goals.md#accountsstaysafe), [goals.birdersKeepOneRecord](../../../../goals/goals.md#birderskeeponerecord)

birdsync downloads from the Macaulay Library only assets whose catalog numbers appear in a
record of the eBird export the user supplied
([product.copiesObservations](../../../local/product/scope.md#copiesobservations)), attaches
each only to the observation made from the record that lists it, and never searches for,
requests, or downloads any other asset.

A user may download their own media and no one else's without permission; media belong to the
account that uploaded them. birdsync cannot tell whose an asset is, so it relies on the export
listing only the user's own assets
([decisions.ownMediaAssumption](../../../decisions.md#ownmediaassumption)), and whether the
export is the user's own rests with the account holder
([decisions.inputIdentityIsTheAccountHolders](../../../decisions.md#inputidentityistheaccountholders)).
What birdsync can control is that it reaches for nothing beyond what the export lists. Both
decisions reopen on a counter-example birdsync itself could have detected.

## macaulayAssetsAttributed

**Each Macaulay Library asset birdsync uploads is attributed on its observation**

Composed from: [mlterms.attributionExpected](../../../imported/ebird.org/mlterms/mlterms.md#attributionexpected)
Goals: [goals.birdersKeepOneRecord](../../../../goals/goals.md#birderskeeponerecord)

Every Macaulay Library asset birdsync uploads is attributed as the Macaulay Library's media
attribution guidelines require.

What birdsync does today is name each uploaded asset in its observation's description with a
`Macaulay Library Asset:` line linking `https://macaulaylibrary.org/asset/<id>`
([product.descriptionContents](../../../local/product/observations.md#descriptioncontents)), and
upload it under the filename `ML<asset ID>`
([product.mediaFilename](../../../local/product/media.md#mediafilename)). Whether that meets the
guidelines is unknown: the guidelines page was not retrieved ([`imports.md`](../../../imports.md),
`mlterms`, Status), and its contents are not supplied from memory. Until it is read, this
requirement is uncovered.

## Acceptance criteria

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

## mediaResyncAdditive

**Media re-sync is additive**

Formerly: AC-013
Covers: [product.additiveMediaResync](../../../local/product/media.md#additivemediaresync)

Assets added to an already-synced checklist are uploaded on the next run.

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

## tempFilesCleanedUp

**Temporary media files are removed after upload**

Formerly: AC-026
Covers: [tech.tempFilesDeleted](../../../local/tech/resources.md#tempfilesdeleted)

After a sync no downloaded media file remains in the temp directory.

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

## downloadCleansUpOnError

**A failed download leaves no temp file**

Formerly: AC-038
Covers: [tech.tempFilesDeleted](../../../local/tech/resources.md#tempfilesdeleted)

A truncated download leaves nothing in the temp directory.

## mediaRequestsOnlyListedAssets

**Every Macaulay Library request is for an asset the export lists**

Covers: [mediaOnlyFromUsersExport](#mediaonlyfromusersexport)

Over a run, every request a test server standing in for the Macaulay Library CDN receives is
for an asset whose catalog number appears in a record of the input export that the run syncs,
and no request is a search.

Does not cover whether the listed assets are the user's own; that rests with the account holder
([decisions.inputIdentityIsTheAccountHolders](../../../decisions.md#inputidentityistheaccountholders)).

## mediaAttachedToListingRecord

**Each asset is attached to the observation of the record that lists it**

Covers: [mediaOnlyFromUsersExport](#mediaonlyfromusersexport)

Given an export in which two records list different assets, each upload the test server
receives is attached to the observation created or found for the record that lists that asset.

## assetNamedInDescription

**Each uploaded asset is named on its observation**

Covers: [macaulayAssetsAttributed](#macaulayassetsattributed)

After a run, the description of every observation to which an asset was uploaded contains a
`Macaulay Library Asset:` line linking `https://macaulaylibrary.org/asset/<id>` for that asset,
and the file was uploaded under the name `ML<id>` with its extension.

Does not decide [macaulayAssetsAttributed](#macaulayassetsattributed): it checks what
birdsync does, not whether that meets the Macaulay Library's attribution guidelines.

## attributionMeetsGuidelines

**birdsync's attribution meets the Macaulay Library's guidelines**

Covers: [macaulayAssetsAttributed](#macaulayassetsattributed)

A reviewer, the owner, compares the attribution birdsync writes with the Macaulay Library's media
attribution guidelines, as vendored in `mlterms`, and confirms it meets them; repeated at every
refresh of `mlterms`.

Cannot be performed yet: the guidelines page has not been retrieved. Until it is, the
requirement counts as uncovered.
