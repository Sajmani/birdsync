package decisions

# Decisions

How each conflict and judgment about birdsync's requirements was settled, oldest first. Content
is append-only: entries are appended at the end and their substance is never rewritten. By hand,
only `Status:` changes when an entry is superseded, amended, or reaffirmed; `river upgrade` may
migrate fields across major River versions. Search `^## ` for the list, `## <name>` for one
entry, and `^Affects:.*<name>` for what touches a requirement.

The first twelve entries were migrated on 2026-11-03 from specman's decision log, where they
were `CR-001` to `CR-012` (`Formerly:` gives the old ID). Each is condensed here and links to
the full original at commit `dc25e3b`. Where the original entry names no decider, `Decided
by:` says so. Decisions drafted for the owner but not yet made are not entered here. They
wait in [`../proposals.md`](../proposals.md) until the owner decides them.

---

## dryRunCountersRelabeled

**Under --dryrun the summary labels the observation counters as hypothetical**

Date: 2026-08-09
Decided by: Sajmani
Status: accepted
Affects: [product.dryRunCounterLabels](local/product/reporting.md#dryruncounterlabels), [tech.honestDryRunCounters](local/tech/safety.md#honestdryruncounters)
Formerly: CR-001

### Decided

Option A of three: under `--dryrun` the summary reads `Would create N` and `Would update N`.
The count stays, because it is useful, but the label no longer claims the work happened.
Recorded as [product.dryRunCounterLabels](local/product/reporting.md#dryruncounterlabels). [tech.honestDryRunCounters](local/tech/safety.md#honestdryruncounters) stands unamended. The alternatives were to narrow [tech.honestDryRunCounters](local/tech/safety.md#honestdryruncounters) to the
media counters, or to document the false summary.

Migrated from specman's decision log; the full entry, with the evidence and the options considered, is [CR-001 at `dc25e3b`](https://github.com/Sajmani/birdsync/blob/dc25e3b6756bd59c2de931e6d080d47318617da4/spec/decisions.md#cr-001--a-dry-run-reports-observations-it-did-not-create).

### Sums

Recorded by `river upgrade` on 2026-11-08, not when decided.

```
product.dryRunCounterLabels sha256:8396b2d7996de5398da2b0b32718790aed3a47c7cd0fdbd22a263fec1eee3147
tech.honestDryRunCounters sha256:453abe22ea5b6963457f7dff179adf302f75a3a0f574a6a63fd5c308c08ce747
```

---

## libraryPackagesReturnErrors

**The ebird and inat packages return errors instead of calling log.Fatal**

Date: 2026-08-09
Decided by: Sajmani
Status: accepted
Affects: [tech.noLogFatalInLibraries](local/tech/style.md#nologfatalinlibraries)
Formerly: CR-002

### Decided

Option A: fix both `ebird.Records` and `inat.DownloadObservations`. The second gains an
`error` return and its callers in `tools/` are updated. They are verified by compilation and
review only, since the tools may not be run. Scoping [tech.noLogFatalInLibraries](local/tech/style.md#nologfatalinlibraries) to new code was rejected, because a
check cannot tell new code from old.

Migrated from specman's decision log; the full entry, with the evidence and the options considered, is [CR-002 at `dc25e3b`](https://github.com/Sajmani/birdsync/blob/dc25e3b6756bd59c2de931e6d080d47318617da4/spec/decisions.md#cr-002--logfatal-in-packages-required-to-return-errors).

### Sums

Recorded by `river upgrade` on 2026-11-08, not when decided.

```
tech.noLogFatalInLibraries sha256:82214729800048320e3d5fe758c6c8d6e07f10ff8dddb1f5133f77b29a6f643f
```

---

## noTaxonFilterOnDownload

**The download of existing observations is not filtered by taxon**

Date: 2026-08-09
Decided by: Sajmani
Status: accepted
Affects: [product.avesOnlyDownload](local/product/download.md#avesonlydownload), [product.noTaxonFilter](local/product/download.md#notaxonfilter), [product.idempotentRerun](local/product/identity.md#idempotentrerun)
Formerly: CR-003

### Decided

Option B: drop the `iconic_taxa[]=Aves` filter and download all of the user's observations.
An observation whose eBird name iNaturalist cannot resolve has no iconic taxon, so a filtered
download misses it and the next run duplicates it ([product.idempotentRerun](local/product/identity.md#idempotentrerun)). [product.avesOnlyDownload](local/product/download.md#avesonlydownload) is withdrawn and replaced
by [product.noTaxonFilter](local/product/download.md#notaxonfilter). The owner rejected a second query for untaxoned observations, because an
unfiltered query re-fetches every `Aves` observation anyway. The extra download volume is
knowingly re-accepted. A side effect is that `--fuzzy` now matches observations of any taxon.

Migrated from specman's decision log; the full entry, with the evidence and the options considered, is [CR-003 at `dc25e3b`](https://github.com/Sajmani/birdsync/blob/dc25e3b6756bd59c2de931e6d080d47318617da4/spec/decisions.md#cr-003--aves-only-download-can-defeat-duplicate-detection).

### Sums

Recorded by `river upgrade` on 2026-11-08, not when decided.

```
product.avesOnlyDownload sha256:fc62b838e469c9c80272492bde864fb098c32a4566b6f04b3b2a0725c3881652
product.idempotentRerun sha256:cd83c688a4b9644398e6265321116165b0a871adcf3d03f38816440867e8ce5c
product.noTaxonFilter sha256:7c63f4369c2b7295650d431c22404bffe8ae1784301638e1faffda53b0189c98
```

---

## fuzzyDocsCorrected

**The README and flag help describe --fuzzy as matching observations without a complete sync key**

Date: 2026-08-09
Decided by: Sajmani
Status: accepted
Affects: [product.fuzzyMatching](local/product/filtering.md#fuzzymatching), [product.incompleteKeyNotRecognized](local/product/identity.md#incompletekeynotrecognized)
Formerly: CR-004

### Decided

Option A: correct the documentation to match the code. The alternative was to exclude
partially keyed observations from the fuzzy index, which would have created duplicates.
Documentation only.

Migrated from specman's decision log; the full entry, with the evidence and the options considered, is [CR-004 at `dc25e3b`](https://github.com/Sajmani/birdsync/blob/dc25e3b6756bd59c2de931e6d080d47318617da4/spec/decisions.md#cr-004----fuzzy-is-documented-as-matching-non-birdsync-observations).

### Sums

Recorded by `river upgrade` on 2026-11-08, not when decided.

```
product.fuzzyMatching sha256:ea3c63d6ab6372f09fef4d4bb89e317eddd805c4da73400a4d9a3c85a3ba011f
product.incompleteKeyNotRecognized sha256:c8faa9397239e0d1556012fd49eb3e24d62629186f2cd7e7a2f8e69cca3bb6cb
```

---

## badRecordsSkippedNotFatal

**A record with an unparseable date, time, latitude, or longitude is skipped and counted, not fatal**

Date: 2026-08-09
Decided by: Sajmani
Status: accepted
Affects: [product.malformedRecordAborts](local/product/failures.md#malformedrecordaborts), [product.badRecordsSkipped](local/product/failures.md#badrecordsskipped), [product.writeFailureAborts](local/product/failures.md#writefailureaborts)
Formerly: CR-005

### Decided

Option A: skip, log with its CSV line, and count the record, reporting the count in the
summary when non-zero. Recorded as [product.badRecordsSkipped](local/product/failures.md#badrecordsskipped). [product.malformedRecordAborts](local/product/failures.md#malformedrecordaborts) is withdrawn. [product.writeFailureAborts](local/product/failures.md#writefailureaborts) is unchanged: a failed
create or update still aborts the run.

Migrated from specman's decision log; the full entry, with the evidence and the options considered, is [CR-005 at `dc25e3b`](https://github.com/Sajmani/birdsync/blob/dc25e3b6756bd59c2de931e6d080d47318617da4/spec/decisions.md#cr-005--record-level-parse-failures-abort-the-run).

### Sums

Recorded by `river upgrade` on 2026-11-08, not when decided.

```
product.badRecordsSkipped sha256:5f9e08c98113bd6d5525276c470625cf0e2099ebbe2a4d56e440b7bb47b08cdb
product.malformedRecordAborts sha256:36d1f322e123e53ec0dcf619edcf58dd0681757553583b6e972d3744d71f6a5b
product.writeFailureAborts sha256:6e897691284dfe52e808a3110b9b08894c17b846a891ddda0c61166938c89c44
```

---

## canonicalFileExtensions

**Uploaded file extensions come from a fixed map of content types, not the system mime database**

Date: 2026-08-09
Decided by: Sajmani
Status: accepted
Affects: [product.mediaFilename](local/product/media.md#mediafilename), [tech.platformIndependent](local/tech/deps.md#platformindependent)
Formerly: CR-006

### Decided

Option A: map known content types to canonical extensions. [product.mediaFilename](local/product/media.md#mediafilename) was amended to name the
mapping, which makes the filename the same on every platform ([tech.platformIndependent](local/tech/deps.md#platformindependent)). Confirmed against the
live service for photos. The sound path was accepted as a low residual risk.

Migrated from specman's decision log; the full entry, with the evidence and the options considered, is [CR-006 at `dc25e3b`](https://github.com/Sajmani/birdsync/blob/dc25e3b6756bd59c2de931e6d080d47318617da4/spec/decisions.md#cr-006--uploaded-file-extensions-varied-by-machine).

### Sums

Recorded by `river upgrade` on 2026-11-08, not when decided.

```
product.mediaFilename sha256:6d1d2fe0862bdb300c0b205daf8325798606ede11a0cd71714389bf9e5def5ff
tech.platformIndependent sha256:5d30d7a04c46b775bd2dd6c52b04d4418660a74837ef6a864af426b898a1cd70
```

---

## descriptionListsOnlyUploaded

**The description lists only assets that actually uploaded**

Date: 2026-08-09
Decided by: Sajmani (presumed: the original entry does not name who decided; to be confirmed)
Status: accepted; decider confirmed by decidersConfirmed
Affects: [product.descriptionContents](local/product/observations.md#descriptioncontents), [product.additiveMediaResync](local/product/media.md#additivemediaresync), [product.mediaFailuresTolerated](local/product/media.md#mediafailurestolerated), [tech.honestDryRunCounters](local/tech/safety.md#honestdryruncounters)
Formerly: CR-007

### Decided

Option A: build the description from the assets that uploaded, so a failed asset stays
absent and the next run retries it ([product.additiveMediaResync](local/product/media.md#additivemediaresync)). If none uploaded, the observation is not updated
and not counted as updated ([tech.honestDryRunCounters](local/tech/safety.md#honestdryruncounters)).

Migrated from specman's decision log; the full entry, with the evidence and the options considered, is [CR-007 at `dc25e3b`](https://github.com/Sajmani/birdsync/blob/dc25e3b6756bd59c2de931e6d080d47318617da4/spec/decisions.md#cr-007--a-failed-media-upload-leaves-its-url-in-the-description).

### Sums

Recorded by `river upgrade` on 2026-11-08, not when decided.

```
product.additiveMediaResync sha256:9193833bcc33f33e045f7664fd595b7c382f3be0b7f06a7e8ce225c26b0524dd
product.descriptionContents sha256:1fde0426f1551fedc4fd7c2a70109415dea8bd2defde45c8ea3cc41c6634377e
product.mediaFailuresTolerated sha256:686e85c77589e9d6d19a2e0e74aab1c684c8f52f5847b5e5f0b37d06b1ca7daa
tech.honestDryRunCounters sha256:453abe22ea5b6963457f7dff179adf302f75a3a0f574a6a63fd5c308c08ce747
```

---

## permanentFailuresRecorded

**A permanent upload failure is recorded in the description and not retried; API errors carry the response body**

Date: 2026-08-09
Decided by: Sajmani
Status: accepted
Affects: [product.permanentUploadFailureRecorded](local/product/failures.md#permanentuploadfailurerecorded), [product.permanentFailuresReported](local/product/failures.md#permanentfailuresreported), [tech.errorsIncludeBody](local/tech/etiquette.md#errorsincludebody)
Formerly: CR-008

### Decided

Option B, plus the error-message fix from option A. A 4xx response other than 401, 408, or
429 is permanent: it is recorded in the description with retry instructions and reported each
run ([product.permanentUploadFailureRecorded](local/product/failures.md#permanentuploadfailurerecorded), [product.permanentFailuresReported](local/product/failures.md#permanentfailuresreported)). Anything else is transient. Error messages include the response body,
minus HTML ([tech.errorsIncludeBody](local/tech/etiquette.md#errorsincludebody)). Download failures stay transient, which is a known limit.

Migrated from specman's decision log; the full entry, with the evidence and the options considered, is [CR-008 at `dc25e3b`](https://github.com/Sajmani/birdsync/blob/dc25e3b6756bd59c2de931e6d080d47318617da4/spec/decisions.md#cr-008--a-permanently-rejected-asset-was-retried-forever).

### Sums

Recorded by `river upgrade` on 2026-11-08, not when decided.

```
product.permanentFailuresReported sha256:c63206700cdfcc88ae615b51746d3f403c7b1a8404f6110c89ffecf3c80d6b40
product.permanentUploadFailureRecorded sha256:03d0d3bcf65a06522722dd2b32b91e71135839598c5fc10b036d7bc0ac58bda0
tech.errorsIncludeBody sha256:32b7d9c033b89fbee756aaee5bc1d8c6fa1cf95967d0e673d878d2e5ff6c4b60
```

---

## extensionMapFromCDN

**The extension map is keyed on content types verified against the CDN, and an unrecognized type falls back by endpoint**

Date: 2026-08-09
Decided by: Sajmani (presumed: the original entry does not name who decided; to be confirmed)
Status: accepted; decider confirmed by decidersConfirmed
Affects: [product.mediaFormats](local/product/media.md#mediaformats), [product.mediaFilename](local/product/media.md#mediafilename), [tech.platformIndependent](local/tech/deps.md#platformindependent)
Formerly: CR-009

### Decided

The Macaulay Library CDN serves sounds as `audio/mpeg3`, which the system mime database does
not know, so every sound download failed. The map is now keyed on content types verified
against the CDN. An unknown type falls back to what the endpoint implies, and the system mime
database is no longer consulted.

Migrated from specman's decision log; the full entry, with the evidence and the options considered, is [CR-009 at `dc25e3b`](https://github.com/Sajmani/birdsync/blob/dc25e3b6756bd59c2de931e6d080d47318617da4/spec/decisions.md#cr-009--no-sound-file-could-be-downloaded-at-all).

### Sums

Recorded by `river upgrade` on 2026-11-08, not when decided.

```
product.mediaFilename sha256:6d1d2fe0862bdb300c0b205daf8325798606ede11a0cd71714389bf9e5def5ff
product.mediaFormats sha256:09d1c7c82b23160077392704e9f5b47eb07ffc34c925ed75e3d9590a78ae24e1
tech.platformIndependent sha256:5d30d7a04c46b775bd2dd6c52b04d4418660a74837ef6a864af426b898a1cd70
```

---

## ownMediaAssumption

**birdsync relies on an eBird export listing only the user's own Macaulay Library assets**

Date: 2026-08-09
Decided by: Sajmani
Status: amended by inputIdentityIsTheAccountHolders
Affects: [mlterms.noOthersMediaWithoutPermission](imported/ebird.org/mlterms/mlterms.md#noothersmediawithoutpermission), [mlterms.downloadOwnMedia](imported/ebird.org/mlterms/mlterms.md#downloadownmedia)
Formerly: CR-010
Reopen: on a counter-example; no review date was set

### Decided

Closed as an accepted assumption, documented rather than enforced. Sampling the
maintainer's export found no asset belonging to another contributor, and [mlterms.mediaOwnedByUploader](imported/ebird.org/mlterms/mlterms.md#mediaownedbyuploader)
explains why. Reopen on a counter-example: an asset in someone's export whose
`userDisplayName` is not theirs. Checking ownership through the Macaulay search API stays in
reserve.

Migrated from specman's decision log; the full entry, with the evidence and the options considered, is [CR-010 at `dc25e3b`](https://github.com/Sajmani/birdsync/blob/dc25e3b6756bd59c2de931e6d080d47318617da4/spec/decisions.md#cr-010--birdsync-cannot-tell-whose-media-it-is-downloading).

### Sums

Recorded by `river upgrade` on 2026-11-08, not when decided.

```
mlterms.downloadOwnMedia sha256:dae5e29525094d7ae153ecfd4dc84321f37bb4d3efe5b2fa8aed3f0e664980a1
mlterms.noOthersMediaWithoutPermission sha256:b99f0a9a41a7926c5fb7e585444020c0c72ff37d3bde00d404c08bd63b229a8b
```

---

## pageByIDAbove

**The observation download pages by id_above with an ascending id sort**

Date: 2026-08-09
Decided by: Sajmani (presumed: the original entry does not name who decided; to be confirmed)
Status: accepted; decider confirmed by decidersConfirmed
Affects: [product.largeAccounts](local/product/download.md#largeaccounts), [tech.idAbovePaging](local/tech/etiquette.md#idabovepaging), [product.downloadExistingObservations](local/product/download.md#downloadexistingobservations), [product.noTaxonFilter](local/product/download.md#notaxonfilter)
Formerly: CR-011

### Decided

Option A: page by `id_above`, as [inatapi.idAbovePastTenThousand](imported/inaturalist.org/inatapi/inatapi.md#idabovepasttenthousand) recommends, and stop on a short page. The
client requires the cursor to advance on every page. This removes the 10,000-result ceiling
behind issue #5 rather than moving it. Restoring an `Aves` plus `unknown` filter stays open as
a pure bandwidth question.

Migrated from specman's decision log; the full entry, with the evidence and the options considered, is [CR-011 at `dc25e3b`](https://github.com/Sajmani/birdsync/blob/dc25e3b6756bd59c2de931e6d080d47318617da4/spec/decisions.md#cr-011--the-download-cannot-page-past-10000-results).

### Sums

Recorded by `river upgrade` on 2026-11-08, not when decided.

```
product.downloadExistingObservations sha256:b02bb364de5671513974db0c6f73b75e652117681d277ae7e8a660f1b587c7c5
product.largeAccounts sha256:d0eb854ab9a2c850c442872f52fbae976f225b528e1322b4bf18b2f0eb9a020e
product.noTaxonFilter sha256:7c63f4369c2b7295650d431c22404bffe8ae1784301638e1faffda53b0189c98
tech.idAbovePaging sha256:8c497a5540785c6c3a41bf908c1350c56bccea5931cc6fe5622e75839da6d0c2
```

---

## notMachineGeneratedContent

**birdsync is not machine generated content under iNaturalist's guidelines**

Date: 2026-08-09
Decided by: Sajmani (presumed: the original entry does not name who decided; to be confirmed)
Status: accepted; decider confirmed by decidersConfirmed
Affects: [inatterms.noMachineGeneratedContent](imported/inaturalist.org/inatterms/inatterms.md#nomachinegeneratedcontent), [inatterms.dontBehaveLikeMachine](imported/inaturalist.org/inatterms/inatterms.md#dontbehavelikemachine), [product.documentCommunityObligations](local/product/community.md#documentcommunityobligations), [product.documentSyncedObservationsIdentifiable](local/product/community.md#documentsyncedobservationsidentifiable)
Formerly: CR-012

### Decided

Resolved on the rule's own text: the vendored machine-generated-content definition permits a
machine "to facilitate posting this content" where a human decided each observation, and
gives a script posting from the user's own curated data as an acceptable example. The volume
concern a forum moderator raised stands, and became [product.documentCommunityObligations](local/product/community.md#documentcommunityobligations) and [product.documentSyncedObservationsIdentifiable](local/product/community.md#documentsyncedobservationsidentifiable). It is also recorded
in [product.verifiableDefaultOn](local/product/filtering.md#verifiabledefaulton)'s rationale, so that the `--verifiable` default is not flipped for convenience.

Migrated from specman's decision log; the full entry, with the evidence and the options considered, is [CR-012 at `dc25e3b`](https://github.com/Sajmani/birdsync/blob/dc25e3b6756bd59c2de931e6d080d47318617da4/spec/decisions.md#cr-012--is-birdsync-machine-generated-content).

### Sums

Recorded by `river upgrade` on 2026-11-08, not when decided.

```
inatterms.dontBehaveLikeMachine sha256:23a4d147281025f155f8be5e9845937558ac5093c056d745dd9439dd561c4238
inatterms.noMachineGeneratedContent sha256:078a5847c3d950a4e8c153ffce4dd3f4e81afdd4841e2c7878ab76e65aab2c1d
product.documentCommunityObligations sha256:d0b474e0dce11942608cff676d16f3c5d433bc84c7191007f339eeee0ec38391
product.documentSyncedObservationsIdentifiable sha256:e8c80a6aa660c0aa42c4126060136d80ae4e3b273cbc0c434e4d3fe2fa8796ff
```

---

## authenticationIsNecessary

**Every request birdsync makes to iNaturalist needs authentication, so authenticating each one meets iNaturalist's rule to authenticate only when necessary**

Date: 2026-11-04
Decided by: Sajmani
Status: accepted
Affects: [inatapi.authenticateOnlyWhenNeeded](imported/inaturalist.org/inatapi/inatapi.md#authenticateonlywhenneeded), [product.downloadExistingObservations](local/product/download.md#downloadexistingobservations), [product.idempotentRerun](local/product/identity.md#idempotentrerun)

### Decided

No exception from iNaturalist is needed, and no gap is accepted. birdsync's writes need the
token by definition. Its one read, the download of the user's existing observations, needs it
too: an unauthenticated download omits private and geoprivacy-obscured observations, and every
one it omits is created again on the next run. Preventing those duplicates is necessary to meet
[product.idempotentRerun](local/product/identity.md#idempotentrerun) and
[inatterms.duplicatesUndesirable](imported/inaturalist.org/inatterms/inatterms.md#duplicatesundesirable),
so authentication is required, which is the case the rule allows. The owner's words: "Preventing
creation of duplicate observations, even private or obscured, seems necessary to meet other
requirements, so I don't see why this requires an exception."

Media downloads go to the Macaulay Library, not iNaturalist, and are outside this rule. A future
iNaturalist request that does not need the token should be sent without it.

### Sums

Recorded by `river upgrade` on 2026-11-08, not when decided.

```
inatapi.authenticateOnlyWhenNeeded sha256:6f4320a719979c0da2ec19a82b71fba6e2ef8e1902ae3e65aa6a3675a196a081
product.downloadExistingObservations sha256:b02bb364de5671513974db0c6f73b75e652117681d277ae7e8a660f1b587c7c5
product.idempotentRerun sha256:cd83c688a4b9644398e6265321116165b0a871adcf3d03f38816440867e8ce5c
```

---

## inputIdentityIsTheAccountHolders

**birdsync cannot tell whether the export and the iNaturalist account belong to the same person; that is the account holder's responsibility, not a gap**

Date: 2026-11-04
Decided by: Sajmani
Status: accepted
Amends: ownMediaAssumption
Affects: [mlterms.noOthersMediaWithoutPermission](imported/ebird.org/mlterms/mlterms.md#noothersmediawithoutpermission), [inatterms.accountHolderResponsible](imported/inaturalist.org/inatterms/inatterms.md#accountholderresponsible), [inatterms.contributorKeepsRights](imported/inaturalist.org/inatterms/inatterms.md#contributorkeepsrights)

### Decided

Nothing stops a user from running birdsync on someone else's eBird export and uploading that
person's observations and media to their own iNaturalist account. birdsync automates uploads but
has no way to check that the input CSV and the iNaturalist account represent the same human, and
checking is technically infeasible. A user can already do the same by hand without birdsync.

So [ownMediaAssumption](#ownmediaassumption) is not an accepted gap against
[mlterms.noOthersMediaWithoutPermission](imported/ebird.org/mlterms/mlterms.md#noothersmediawithoutpermission)
and needs no review date. The obligation rests with the person running birdsync, as
[inatterms.accountHolderResponsible](imported/inaturalist.org/inatterms/inatterms.md#accountholderresponsible)
already says of everything done under an account. The assumption stands as recorded, extended
from "the export lists only the user's own assets" to "the export is the user's own", and still
reopens on a counter-example in which birdsync itself could have told.

### Sums

Recorded by `river upgrade` on 2026-11-08, not when decided.

```
inatterms.accountHolderResponsible sha256:d63590df5a6c3c996211234a70454cf5bbb6bbfa775ea520d3fc58744a935831
inatterms.contributorKeepsRights sha256:6a8d22140b7f977b0a435e4341f0c9feaa01168986ee46f8a517d522a900e160
mlterms.noOthersMediaWithoutPermission sha256:b99f0a9a41a7926c5fb7e585444020c0c72ff37d3bde00d404c08bd63b229a8b
```

---

## mutatingToolsDeleted

**The mutating programs in tools/ were deleted, and the rule that nobody may run them was withdrawn**

Date: 2026-08-09
Decided by: Sajmani
Status: accepted
Affects: [tech.neverRunMutatingTools](local/tech/safety.md#neverrunmutatingtools), [tech.toolsReadOnly](local/tech/safety.md#toolsreadonly), [tech.neverRunTools](local/tech/safety.md#neverruntools)
Formerly: T-011 withdrawal note

### Decided

The mutating tools were deleted rather than governed by a rule nobody could check, in commit
[`cca3255`](https://github.com/Sajmani/birdsync/commit/cca3255) ("Delete the mutating tools; enforce that tools/ stays read-only").
[tech.toolsReadOnly](local/tech/safety.md#toolsreadonly) and [tech.neverRunTools](local/tech/safety.md#neverruntools) replace [tech.neverRunMutatingTools](local/tech/safety.md#neverrunmutatingtools). specman recorded
the withdrawal only as a `Status:` line on `T-011`; this entry was written on 2026-11-04, when
the owner confirmed it, and is dated by the commit that made the change.

### Sums

Recorded by `river upgrade` on 2026-11-08, not when decided.

```
tech.neverRunMutatingTools sha256:257f148b05b0f337bb997452bf11c40541d0725cdacede405526cd0a9a2ff99c
tech.neverRunTools sha256:993c0629cf1682d31c4ffcdc88b124f5986ce5456d4baaa5ff6b62b6d31338c0
tech.toolsReadOnly sha256:71a2f130e5ee913c0006699546c02cd3a95e10ca275124924946d9a90646da81
```

---

## localTechAboveProduct

**Within the Local tier, tech outranks product**

Date: 2026-11-04
Decided by: Sajmani
Status: accepted
Affects: [precedence](precedence.md)

### Decided

Confirmed as drafted at adoption: `tech` ranks above `product` in `precedence.md`, because
birdsync's safety invariants outrank convenience. specman's manifest ranked nothing within
Local.

---

## decidersConfirmed

**The owner decided the four migrated decisions whose specman entries named no decider**

Date: 2026-11-04
Decided by: Sajmani
Status: accepted
Affects: none; it confirms the record of [descriptionListsOnlyUploaded](#descriptionlistsonlyuploaded), [extensionMapFromCDN](#extensionmapfromcdn), [pageByIDAbove](#pagebyidabove), and [notMachineGeneratedContent](#notmachinegeneratedcontent)

### Decided

Their `Decided by:` lines said "presumed"; the owner confirmed on 2026-11-04 that he decided
each. Those entries are otherwise unchanged, as `decisions.md` is append-only.

---

## uuidAdopted

**github.com/google/uuid is adopted as an Advisory source for its license**

Date: 2026-11-04
Decided by: Sajmani
Status: accepted
Affects: [tech.oneExternalDependency](local/tech/deps.md#oneexternaldependency)

### Decided

Adopted as drafted at discovery: the module's three-clause BSD-style license, read from the
module at v1.6.0, is recorded in `imports.md` as Advisory. birdsync redistributes no binaries,
so the notice obligation does not arise today.

### Sums

Recorded by `river upgrade` on 2026-11-08, not when decided.

```
tech.oneExternalDependency sha256:4bcba12881b3d5fdff440eca9b56b1121d34c3b5ce8c9a6e564d9d5c55d5926e
```

---

## restatementsComposed

**Local requirements that only restated an imported rule are withdrawn; composition derives them instead**

Date: 2026-11-05
Decided by: Sajmani
Status: accepted
Affects: [tech.maxPageSize](local/tech/etiquette.md#maxpagesize), [tech.requestPacing](local/tech/etiquette.md#requestpacing), [tech.idAbovePaging](local/tech/etiquette.md#idabovepaging), [tech.userAgent](local/tech/etiquette.md#useragent), [product.documentCommunityObligations](local/product/community.md#documentcommunityobligations)

### Decided

specman's local requirements included hand-written restatements of rules from the imported
sources, so composition never had to derive them. On 2026-11-04 an agent with no context,
following only River's skill, recomposed the imported part of the package with five such items
hidden. It reproduced [tech.maxPageSize](local/tech/etiquette.md#maxpagesize), [tech.idAbovePaging](local/tech/etiquette.md#idabovepaging), and [tech.requestPacing](local/tech/etiquette.md#requestpacing) in
substance, the paging details with help from [pageByIDAbove](#pagebyidabove). It reproduced
[tech.userAgent](local/tech/etiquette.md#useragent) except for the version. It reproduced part of
[product.documentCommunityObligations](local/product/community.md#documentcommunityobligations), missing that users should review what was created.

So the three fully reproduced items are withdrawn, and the composed items replace them.
[tech.userAgent](local/tech/etiquette.md#useragent) is narrowed to the version, its local addition. [product.documentCommunityObligations](local/product/community.md#documentcommunityobligations)
stays, and is composed into the composed requirement on documenting the user's obligations. The
single choke point that enforces pacing, and why it replaced an argument from slowness, are
design and move to `docs/arch.md`.

The same recomposition wrote out six imported obligations as birdsync rules, where the
migration had included them unchanged: the JWT header, posting only as the account holder,
authenticating only when necessary, posting only the user's own records, using only the media
the export lists, and attributing Macaulay assets. Their criteria are new, and uncovered until
guardrails exist.

### Sums

Recorded by `river upgrade` on 2026-11-08, not when decided.

```
product.documentCommunityObligations sha256:d0b474e0dce11942608cff676d16f3c5d433bc84c7191007f339eeee0ec38391
tech.idAbovePaging sha256:8c497a5540785c6c3a41bf908c1350c56bccea5931cc6fe5622e75839da6d0c2
tech.maxPageSize sha256:3e70df34eefab7db73ee2cc9d83490f9d78e4bc584a86940068b4bd0d6e2dc19
tech.requestPacing sha256:6510a97f49d15ccce60477361dc74f296458044aae560b1e82da2443417fdcb7
tech.userAgent sha256:55fa2e311e0751625bf89a00ca14f4cc1031b291523bb12a74326363449eec9e
```

---

## dailyRequestBudget

**birdsync enforces a 9,000-request daily budget in the client and across runs via downloaded created_at timestamps**

Date: 2026-11-08
Decided by: Sajmani
Status: accepted
Affects: [inatapi.requestRate](imported/inaturalist.org/inatapi/inatapi.md#requestrate), [inatterms.rateAndVolume](imported/inaturalist.org/inatterms/inatterms.md#rateandvolume), [inatterms.dontBehaveLikeMachine](imported/inaturalist.org/inatterms/inatterms.md#dontbehavelikemachine), [tech.dailyRequestBudget](local/tech/etiquette.md#dailyrequestbudget)

### Decided

birdsync meets the daily half of [inatapi.requestRate](imported/inaturalist.org/inatapi/inatapi.md#requestrate)
("around 10k API requests a day") without external state by enforcing a conservative 9,000-request
budget in two places ([tech.dailyRequestBudget](local/tech/etiquette.md#dailyrequestbudget)):
`inat.Client` refuses requests past 9,000 in a single run (covering downloads and `--dryrun`), and
`birdsync()` requests `created_at` in the existing-observation download so it can estimate write
requests already made on the account in the past 24 hours (one create per recent observation, one
upload per attached or permanently failed asset, and one description update when assets were
handled), add the download and write requests of the current run, and stop cleanly before starting
an observation that would push the 24-hour total past 9,000. Because re-running birdsync is
idempotent ([product.idempotentRerun](local/product/identity.md#idempotentrerun)), the user can
resume 24 hours later with the same command.

Two small residuals are accepted as negligible: concurrent birdsync processes started before
either finishes writing, and runs narrowed with `--after` or `--before` that omit observations
observed outside that date window even if synced earlier the same day.

### Sums

```
inatapi.requestRate sha256:d74d7fa37dba61d3f71c184b00150432dbc73ae88cc099044e37f8bef901f976
inatterms.dontBehaveLikeMachine sha256:23a4d147281025f155f8be5e9845937558ac5093c056d745dd9439dd561c4238
inatterms.rateAndVolume sha256:31b383a584ace928b6292b3ff22e3096288c7350911191148f3d309e585c9278
tech.dailyRequestBudget sha256:407b28a3dd428824dc860c2b03a9f27758432f910441d553dfd589dbd84e310c
```
---

## birdsyncApproved

**Approve composed/internal/birdsync**

Date: 2026-11-10
Decided by: Sajmani
Status: accepted
Affects: none
Approves: birdsync, sha256:a4fc4272df0f979da68b752d0e4cf611a7f027748e63b9adc8d40673253c30c2

### Decided

Approved `composed/internal/birdsync/` on PR #8 (`river-migration`).

### Sums

```
decisions.authenticationIsNecessary sha256:d76e21999da442db0a26813eecedfb51552950ba80bff6bd8f48d4d2f884d350
decisions.badRecordsSkippedNotFatal sha256:1934a3f18b89d88bf3c469e81e301e70197fcfd5aba848f9345272b05087be44
decisions.canonicalFileExtensions sha256:ea18f968f8d2ae1d54c0ce0dae2d8de2b31a03d62a848f4c28d7d63630a48450
decisions.dailyRequestBudget sha256:539aee901c0dc74bf99312cb68936c103fb339ff2a45510bb78c9f5b681cc423
decisions.descriptionListsOnlyUploaded sha256:10038789afb12631c85ced8865d04b747e9b19c6e9912c58df07acf87c23e30f
decisions.dryRunCountersRelabeled sha256:f2cfb3e3b4023da80781b002f262a8130275d69a1339e5ef87c2ef826d30f3f1
decisions.extensionMapFromCDN sha256:411ab77554227401ec5438d2b49a30ba858a8f89c8e7637f0226ac8ca5fa86eb
decisions.fuzzyDocsCorrected sha256:3733edbea16c65c88b2a5dfa9cdb44793f380d6516e207513aec6b82a9003570
decisions.inputIdentityIsTheAccountHolders sha256:c8f9c5436abe6079bf5e5686ef90add7a6ffc9e41dd9e5ecd7e3faf753a8f357
decisions.libraryPackagesReturnErrors sha256:8bb97f43e97a41e34486c3e32ea253a4ed5102932e7991d587ad05a692338e0e
decisions.mutatingToolsDeleted sha256:ca97c3ecfe7c0bd2b1d086af5158f596624180236c5013a733f0504c37fde349
decisions.noTaxonFilterOnDownload sha256:6acd2746e41dc8b104f303ba036007e78e70cdda2d30f5565838ca008692d023
decisions.notMachineGeneratedContent sha256:258a8fcb7f974ae2cf9cb9a0481485a6a89809eaf3020fb3f72b814177b56557
decisions.ownMediaAssumption sha256:d5771e264b75acd4906454aac293dd99fea5c5f48c15dd55fb61cd8380fa5f14
decisions.pageByIDAbove sha256:9af11ce2d061e57917a974b425d7874010b7d34d9fd58569fc6be712bde30acd
decisions.permanentFailuresRecorded sha256:036b0f3663a1bc668701629c15297ed4677fec71c915b6875c88d70d2ddf830e
decisions.restatementsComposed sha256:e34aa3981ab02aab15d3f00ffcc52a8dad9d6c2aef51e451c03c23bb43ecba18
decisions.uuidAdopted sha256:465558ac9e00ae9feaf21e9a1fb9dc163d7279b12d9ec96549c16f9c5f264914
inatapi.authenticateOnlyWhenNeeded sha256:6f4320a719979c0da2ec19a82b71fba6e2ef8e1902ae3e65aa6a3675a196a081
inatapi.idAbovePastTenThousand sha256:e704186a5f2e384b3f67400adbfc1578b1f76445695da95035acd799920d9cd3
inatapi.identifyUserAgent sha256:71c0db5218fded5551c38b994570d3aee12ae5962dc25a28166830e5f6418eef
inatapi.jwtAuthorization sha256:b2651b4dbd4eabb2a36b4b87548aa97b8bd1234ac0446ec35c238b2e48bec843
inatapi.largestPageSize sha256:b3d51caec96c1c76e449c7ba10fa2718c114e8ccc2519999864b893d68f94f06
inatapi.noGroupAccounts sha256:c012426c607a284c1e40ddab954bf2bba4a33d4809dddd018ae9c736d69debe5
inatapi.requestRate sha256:d74d7fa37dba61d3f71c184b00150432dbc73ae88cc099044e37f8bef901f976
inatapi.singleIPAndMediaCaps sha256:007ac2cf8f2a5ed03860205cd8eb7d19c83eac4f8e386f70915252a8b560eafe
inatterms.accountHolderResponsible sha256:d63590df5a6c3c996211234a70454cf5bbb6bbfa775ea520d3fc58744a935831
inatterms.contributorKeepsRights sha256:6a8d22140b7f977b0a435e4341f0c9feaa01168986ee46f8a517d522a900e160
inatterms.dontBehaveLikeMachine sha256:23a4d147281025f155f8be5e9845937558ac5093c056d745dd9439dd561c4238
inatterms.duplicatesUndesirable sha256:95c693bd2539d676c1e3e1774c19f0511cc77f66e93ebe7989b81b7b6d5793b4
inatterms.noMachineGeneratedContent sha256:078a5847c3d950a4e8c153ffce4dd3f4e81afdd4841e2c7878ab76e65aab2c1d
inatterms.rateAndVolume sha256:31b383a584ace928b6292b3ff22e3096288c7350911191148f3d309e585c9278
mlterms.attributionExpected sha256:c36c3a440f121e7dec84e081ac6bc71beefd4536ecc59e9cc407b66a3a8350a3
mlterms.contributorKeepsCopyright sha256:7ffd0fd6ad0af908c17baeb56bb92245cc7f832be60ecaff708d5a2cae7063e2
mlterms.downloadOwnMedia sha256:dae5e29525094d7ae153ecfd4dc84321f37bb4d3efe5b2fa8aed3f0e664980a1
mlterms.mediaOwnedByUploader sha256:ab74d7c15e829896c2a5355b53462183e9d1bab59ae4627cf21035f028086bd7
mlterms.noOthersMediaWithoutPermission sha256:b99f0a9a41a7926c5fb7e585444020c0c72ff37d3bde00d404c08bd63b229a8b
precedence.advisory sha256:a5076f38d24fe5e65a0d4e833af9990f641176132838806cccec1c707c1584f1
precedence.governing sha256:4d8a3e3d4853d9d8e136b597e69086fcfd478e6c622e51b3e1105f2aea6aee16
precedence.local-product sha256:369440908eab7afccccdea5c95a5ff3ec1239c113bba68a50a32dcde410ce35a
precedence.local-tech sha256:7a7e7c4efffac5dda038508bdeb9b5d22943847308d3ece88639293a6988ce0b
precedence.mandatory sha256:ef52942de0ee503115e3720878e8a06739ac2d6b6512c752d60413c0dedb6a05
product.additiveMediaResync sha256:9193833bcc33f33e045f7664fd595b7c382f3be0b7f06a7e8ce225c26b0524dd
product.afterFlag sha256:9a49f9ee67481a66bf3324e1d45000c332800fa0029a2a31c374ebef8341fda1
product.afterLaterThanBeforeExits sha256:a2ac3cfd6a20d11170340c8b941fa4ee79064fa8a62e72f49167d3f1aa188124
product.avesOnlyDownload sha256:fc62b838e469c9c80272492bde864fb098c32a4566b6f04b3b2a0725c3881652
product.badRecordsSkipped sha256:5f9e08c98113bd6d5525276c470625cf0e2099ebbe2a4d56e440b7bb47b08cdb
product.beforeFlag sha256:a03138469eb9c9fb91c7d152c4cc215937c6e31d79d89c1a579a3284963e0c0a
product.boolFlagEqualsForm sha256:a815d07187ce7ecefeddfed2b5c7554dbba7cfa5353dd6b05447534ebc13bb05
product.checklistCoordinates sha256:57d8869042463a30cdd92cd140d1fd30a02e58b93f20a5b42c9fa21eefabc11d
product.conditionalFailureCount sha256:f495a4be0619af85df071023575611108c9453a404a05d0481b5738068032f76
product.conditionalSkipCounters sha256:aef8c37b1f207b34d9b08ac224c301e99e387f974532b77e7533fc12f28c5037
product.copiesObservations sha256:c3958422f217bbf88f782bf43fe943e132c38b1b2530ddf528618f1d02275448
product.createBeforeMedia sha256:94eb342a028772a67f6c5682cb6640d550eced0dd47619c0aaa6aacb67eab092
product.createdObservationWild sha256:791403c6c324db26d0c026c75fbd22db2046e33752c0e7109caf43b916d44682
product.credentialsFromEnvironment sha256:003590cedfaba8fa6dc4396e2135ace1cf35496ab4f7c5038a5466dff41f936a
product.credentialsNeverLogged sha256:d8eb84c0707a16a283079d1e020b6e85ef70bc5455561dfcad02d8fe575b153f
product.csvExportInput sha256:74ade4c2f27b33ab383dc5a51d70f73e8c7ea62ce047da9295da3f132a93b27b
product.dateFlagFormats sha256:658b73de15ae862c9e4b2a70c65a970d15e6e01d5aa3f025a77f719c396d04ca
product.dateWindowNarrowsDownload sha256:f4b41c74c98075ddd69cb89ede75925436e28775cf9be2e88b1ad68f63a1555d
product.debugExplainsSkips sha256:7c944a6c64261c751f480d3657a9c41c5516eabefac84a89978fbe90702b48fc
product.descriptionContents sha256:1fde0426f1551fedc4fd7c2a70109415dea8bd2defde45c8ea3cc41c6634377e
product.descriptionUpdateKeepsMedia sha256:e34fcca98c7cbe3a4011bad965934cb9b4bcde2e4d351a6a3d154a074652b84b
product.documentCommunityObligations sha256:d0b474e0dce11942608cff676d16f3c5d433bc84c7191007f339eeee0ec38391
product.documentSyncedObservationsIdentifiable sha256:e8c80a6aa660c0aa42c4126060136d80ae4e3b273cbc0c434e4d3fe2fa8796ff
product.downloadExistingObservations sha256:b02bb364de5671513974db0c6f73b75e652117681d277ae7e8a660f1b587c7c5
product.dryRunCounterLabels sha256:8396b2d7996de5398da2b0b32718790aed3a47c7cd0fdbd22a263fec1eee3147
product.dryRunLogPrefix sha256:bab185b92a86028f3636596e02536e83fa920e433b324b102325ba6f2ae85869
product.dryRunMediaCount sha256:23f84e8f443dbb37de74da4cc33c98c88cada9eba9178cb9cc63328c03974e86
product.dryRunNoWrites sha256:6233a64cac81a48337be4b73c645982b7bfae8b4f392a0a057d90073f5b52e91
product.endOfRunSummary sha256:76671371e0f703410814c085acce3a050c5db705246c8949ce5bd8cb12465aec
product.flagsBeforeArgument sha256:0df71b0a68ad463b92d07978758c01d31c8554b94fdf3af95899ff5e1af44fb0
product.fuzzyMatching sha256:ea3c63d6ab6372f09fef4d4bb89e317eddd805c4da73400a4d9a3c85a3ba011f
product.fuzzyNeverEmptyName sha256:886c1dd6bf64a7487b3a8ac46e6fc8a43f180553f5e177e1671a472f8540547f
product.fuzzyOffByDefault sha256:d683a5e1524cb1b1ec9556953e0a1b9840e8b2fdde2274e97e9cae8f915d45d9
product.fuzzyParsedDates sha256:dc7b34f961ea5bb7986b3995fc1272e96faef677fb9f5e0463cbc44f29ee458a
product.idempotentRerun sha256:cd83c688a4b9644398e6265321116165b0a871adcf3d03f38816440867e8ce5c
product.incompleteKeyNotRecognized sha256:c8faa9397239e0d1556012fd49eb3e24d62629186f2cd7e7a2f8e69cca3bb6cb
product.largeAccounts sha256:d0eb854ab9a2c850c442872f52fbae976f225b528e1322b4bf18b2f0eb9a020e
product.malformedRecordAborts sha256:36d1f322e123e53ec0dcf619edcf58dd0681757553583b6e972d3744d71f6a5b
product.mediaCountMismatchReported sha256:26c61b21f1ea8f8d9ba78b8ffe4ba7951a9dec6ed6584c03212aa866042c1296
product.mediaFailuresTolerated sha256:686e85c77589e9d6d19a2e0e74aab1c684c8f52f5847b5e5f0b37d06b1ca7daa
product.mediaFilename sha256:6d1d2fe0862bdb300c0b205daf8325798606ede11a0cd71714389bf9e5def5ff
product.mediaFormats sha256:09d1c7c82b23160077392704e9f5b47eb07ffc34c925ed75e3d9590a78ae24e1
product.neverModifyOthersObservations sha256:fd48f09b5730790da143114eefe7abc0a8329c8e502d3d671c92b5aab626b54a
product.noReverseSync sha256:3cd36ba5cecc4d8d8019d8b8bc9fcff7e8b2ab35bff64711eb845dbf88e797a5
product.noRollback sha256:53620c932063f6d33bda075d9382d2dd069d4f52d6903924be040956d2d94fd1
product.noTaxonFilter sha256:7c63f4369c2b7295650d431c22404bffe8ae1784301638e1faffda53b0189c98
product.noTaxonomyReconciliation sha256:b0de609650fe0ca8fc890fd8656e94734ab7c49c4647c1318cd0c63915af1572
product.observationFieldMapping sha256:ff0077b38bc5a6d1f29e98168badc01977d0e6112d7faa184a1807e7a2910543
product.observationsPublic sha256:e4165695e6da85d8776f66c3db14a12534994a9dd49e76199b0b0a1c39b2c291
product.observedDateTime sha256:603f7846b3d99209834bd47edadbeccc98e82c1f75c6e2b156db431156f4d3b6
product.onePositionalArgument sha256:ccff4f46399078018cd877fb5cbefaeea7296a7672e190b351dad781a448a099
product.oneWaySync sha256:efcd7f3fb00332a1e9d6fd351b797a6221a0b543bb48787e7f7371eab8dbef1f
product.permanentFailuresReported sha256:c63206700cdfcc88ae615b51746d3f403c7b1a8404f6110c89ffecf3c80d6b40
product.permanentUploadFailureRecorded sha256:03d0d3bcf65a06522722dd2b32b91e71135839598c5fc10b036d7bc0ac58bda0
product.photoThenSoundFallback sha256:2b7d9d849e79cddf233df06bac10716cdd90938d51b7dd040cd62ac3b9605058
product.promptForMissingCredentials sha256:3e4e30fab003e189a0ddc5356ab260c6db84c562c7e7ebf403e21f1e57756da2
product.removedAssetsReported sha256:7fbbc93c19ea0f603fa3c4bb58b05d4f75fbbf2f7626293d14b5c60bbd0b1d51
product.skipOrder sha256:3754b19c432a390e82b8ae7012ce3d4495af0e2798c25bae27d290361cdbd342
product.speciesGuessScientificName sha256:7fec544b3df7171e979618a0bdccb7bddd3ea27235fc68fde90baa0fcf508d90
product.syncKey sha256:f7a7b5816dbbb811bcb476a0adbf7f7f2bf9f508d8a5f301b68d375faefac33f
product.taxonNotInSyncKey sha256:adea3071e6ceb295367cf63e1d2499416883172ce70de29743a2893d1767fb06
product.tokenFraming sha256:4634fed24d50690766d01f04bf7596f68e94bc3ac7d6ee8ae70c7b1fff0e14c9
product.toolsNotProduct sha256:b57ae82ce460d3e2c21a7c78c422aa2cc9ede95b0ac5eea5072036d4a455cfff
product.unauthorizedSaysRefreshToken sha256:f0a82334ab613552c78a8ef26922b8704c8056e9f7693068ece97c78a7f6ebc6
product.unopenableCsvExitsEarly sha256:29dccffb1e49626347195c08ec05424ebe8d4be60425bda44f205be6184f20fc
product.unusableInputExplained sha256:2617377cdf6bbaeb6e4dd06dccd6e0d79e2db14ff48ad835b29909cc40525395
product.verifiableDefaultOn sha256:025a9317552d7efd8a87c85ff89de2d8d2deb1d974473a05ecabb7d0c6e385d5
product.writeFailureAborts sha256:6e897691284dfe52e808a3110b9b08894c17b846a891ddda0c61166938c89c44
tech.americanSpelling sha256:9b4e5166032049f8f87eafc74f71992a70132da10233c647452a789901b33361
tech.assetTypeNeedsDownload sha256:f4b3ac140638b08ce5320603f191a086437430a444648dec921f5327a492c889
tech.baseURLSeams sha256:0babe7b6cc845dd79a6f2a416aaa4bee11c74d57fe0f3d8b29bf4e2fb63eac2e
tech.ciLatestStable sha256:ab86108ef3677568effcf5ea75027efeba6e65d15e2abac00445ec4e9821b51b
tech.ciStandingChecks sha256:3688c1a5710576f331f64474d0469c0b929b8efba1f9b0ad46a31a36c63156c1
tech.clientInterfaces sha256:312f1884f6cd259e005649f08f8bcb10478264db22380dabe9600ab91b5f3fff
tech.clientMutatesUnconditionally sha256:3c594bce5015da5390e6f803af700db70421c67b6e0b4a594fc417f9fb57e209
tech.commentsExplainWhy sha256:c74fcb1d3b14400a6ddafbae72c02e47740b4dbc288c2aabcafe0d3f34d21873
tech.credentialsNeverExposed sha256:e2dcb4a3bd6e4271788f9bff5a9b2165a46579241b83d552eb2b16ce4c92e5b6
tech.csvByHeaderName sha256:d8286787f46d4184123ba631c87caf478039c23639bd479e3aa1543351514081
tech.dailyRequestBudget sha256:407b28a3dd428824dc860c2b03a9f27758432f910441d553dfd589dbd84e310c
tech.datesViaObserved sha256:90018726da5e6e419cd21f29ef8a127af34d55ced5067cd007bb3e421479ff10
tech.dryRunGateAtCallSite sha256:0ca213c1d817609fd1ccadb2987c12c809358ac804f96db104ef221ce25de8b6
tech.dryRunPrefixLogged sha256:3d3406eaaf5e889b36db70458217d1d2e4296b9b0b6c7ed218e7854078edd7b0
tech.emptyNameNotIndexed sha256:86534ee9f9269cb5ebd133e06638ae50d097b3a42ef44b374e89626da5c82c7d
tech.errorWrapping sha256:002ada239d2d7e4e746c6b494846a40c756f5818d89ceab4baed2829a9f378ba
tech.errorsIncludeBody sha256:32b7d9c033b89fbee756aaee5bc1d8c6fa1cf95967d0e673d878d2e5ff6c4b60
tech.goDirectivePolicy sha256:5d65f2e5d853130e60d64d2563a47f541965902c1157ba600b2662f5e4bb7887
tech.gofmtClean sha256:f89469279a6e60aa90504515240337dfaf52cbcf72e033a1c98f4283cf16b805
tech.honestDryRunCounters sha256:453abe22ea5b6963457f7dff179adf302f75a3a0f574a6a63fd5c308c08ce747
tech.idAbovePaging sha256:8c497a5540785c6c3a41bf908c1350c56bccea5931cc6fe5622e75839da6d0c2
tech.maxPageSize sha256:3e70df34eefab7db73ee2cc9d83490f9d78e4bc584a86940068b4bd0d6e2dc19
tech.memoryScalesWithExport sha256:738b4a15cfa76bfa9f7902b7347c9f90b8d519d046d608120ae69906f29dca19
tech.neverRunMutatingTools sha256:257f148b05b0f337bb997452bf11c40541d0725cdacede405526cd0a9a2ff99c
tech.neverRunTools sha256:993c0629cf1682d31c4ffcdc88b124f5986ce5456d4baaa5ff6b62b6d31338c0
tech.noLiveServicesInTests sha256:7ade179dfbef7aafa41e9bfee7c006f9ae56ace2a5c351708ae6949a06da15e5
tech.noLogFatalInLibraries sha256:82214729800048320e3d5fe758c6c8d6e07f10ff8dddb1f5133f77b29a6f643f
tech.oneExternalDependency sha256:4bcba12881b3d5fdff440eca9b56b1121d34c3b5ce8c9a6e564d9d5c55d5926e
tech.platformIndependent sha256:5d30d7a04c46b775bd2dd6c52b04d4418660a74837ef6a864af426b898a1cd70
tech.progressLogging sha256:efddec1e0048f666a8ab2b668080324c5c4e201f9465f04fd3695f9321e3959f
tech.quotationsNeverRespelled sha256:773967c902b5c74ea0fc0c07cbcc326c8037b81125f6f8b1cfc4cab2e803cd1e
tech.requestPacing sha256:6510a97f49d15ccce60477361dc74f296458044aae560b1e82da2443417fdcb7
tech.resetFlagsInTests sha256:dbf132305439da0d272baaf965558155aa3203549aa90cc2057670c25b886a1b
tech.singleModule sha256:4efefcf898af23d513d68b1d7ba72c35bf2675517fb414a118303bf1fe7b22e6
tech.tempFilesDeleted sha256:2d20469e3eb967a2888d06098a0eb237f2626f0d46efac9d63ed0cf46486c9df
tech.toolsReadOnly sha256:71a2f130e5ee913c0006699546c02cd3a95e10ca275124924946d9a90646da81
tech.updateKeepsPhotos sha256:28cb54444c3731e1ecffbaa3024d110ad0675b977e0336903a378e684fd77afa
tech.userAgent sha256:55fa2e311e0751625bf89a00ca14f4cc1031b291523bb12a74326363449eec9e
tech.vetClean sha256:586358c0fd40f90f408a05826db8a8be08014fd4007e53b224f49073bbbb0908
```
