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
