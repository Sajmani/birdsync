package decisions

# Decisions

How each conflict and judgment about birdsync's requirements was settled, oldest first. Entries
are appended at the end and never rewritten. When an entry is superseded or amended, only its
`Status:` line changes. Search `^## ` for the list, `## <name>` for one entry, and
`^Affects:.*<name>` for what touches a requirement.

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
