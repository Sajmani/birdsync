package proposals

# Proposals for the owner

**Not normative.** Everything here was drafted by an agent on 2026-11-03, while birdsync
was adopting River, and **none of it has been decided**. River's rule is that an agent drafts
and a person decides. The owner (Sajmani) accepts, changes, or rejects each item. Accepted
goals then go into `goals/`, accepted decisions into `requirements/decisions.md`, and each
item is deleted from here once settled. Search `^## ` for the list.

River's method has no place for undecided drafts, so this file is a local convention. See
[`bindings.md`](bindings.md).

## stagesAdopted

**Adopt the Requirements and Verification stages; leave Implementation and Release informal**

Recorded provisionally in [`bindings.md`](bindings.md). This is the most conservative choice
that keeps what specman already had: requirements, and acceptance criteria with guardrails.
Verification depends on Requirements, so the two go together. specman had no rollout or design
record beyond `docs/arch.md`.

## localPrecedence

**Rank `tech` above `product` within the Local tier**

Specman ranked nothing within Local. The proposal follows the existing statement that the
safety invariants outrank convenience. See [`requirements/precedence.md`](requirements/precedence.md).

## goals

**Accept, reword, or replace four candidate goals**

Every requirement is currently unmapped. Clustering them by who benefits and what they gain
gives four candidates. Each lists the requirements it would be cited by. A requirement may
serve more than one goal, but each is listed once here, under its closest fit. Nothing has
been written to `goals/` or to any `Goals:` field.

### Candidate: birdersKeepOneRecord

Birders who keep their records in eBird get them, with their photos and sounds, into iNaturalist without entering anything twice.

Who benefits: Users who log in eBird. What they gain: one record kept, published in both places.

Would be cited by 51 requirements:

- [product.copiesObservations](requirements/local/product/product.md#copiesobservations)
- [product.oneWaySync](requirements/local/product/product.md#onewaysync)
- [product.csvExportInput](requirements/local/product/product.md#csvexportinput)
- [product.noReverseSync](requirements/local/product/product.md#noreversesync)
- [product.noTaxonomyReconciliation](requirements/local/product/product.md#notaxonomyreconciliation)
- [product.onePositionalArgument](requirements/local/product/product.md#onepositionalargument)
- [product.flagsBeforeArgument](requirements/local/product/product.md#flagsbeforeargument)
- [product.boolFlagEqualsForm](requirements/local/product/product.md#boolflagequalsform)
- [product.unopenableCsvExitsEarly](requirements/local/product/product.md#unopenablecsvexitsearly)
- [product.afterLaterThanBeforeExits](requirements/local/product/product.md#afterlaterthanbeforeexits)
- [product.debugExplainsSkips](requirements/local/product/product.md#debugexplainsskips)
- [product.credentialsFromEnvironment](requirements/local/product/product.md#credentialsfromenvironment)
- [product.promptForMissingCredentials](requirements/local/product/product.md#promptformissingcredentials)
- [product.tokenFraming](requirements/local/product/product.md#tokenframing)
- [product.unauthorizedSaysRefreshToken](requirements/local/product/product.md#unauthorizedsaysrefreshtoken)
- [product.downloadExistingObservations](requirements/local/product/product.md#downloadexistingobservations)
- [product.dateWindowNarrowsDownload](requirements/local/product/product.md#datewindownarrowsdownload)
- [product.skipOrder](requirements/local/product/product.md#skiporder)
- [product.afterFlag](requirements/local/product/product.md#afterflag)
- [product.beforeFlag](requirements/local/product/product.md#beforeflag)
- [product.dateFlagFormats](requirements/local/product/product.md#dateflagformats)
- [product.fuzzyMatching](requirements/local/product/product.md#fuzzymatching)
- [product.fuzzyNeverEmptyName](requirements/local/product/product.md#fuzzyneveremptyname)
- [product.fuzzyParsedDates](requirements/local/product/product.md#fuzzyparseddates)
- [product.fuzzyOffByDefault](requirements/local/product/product.md#fuzzyoffbydefault)
- [product.createdObservationWild](requirements/local/product/product.md#createdobservationwild)
- [product.checklistCoordinates](requirements/local/product/product.md#checklistcoordinates)
- [product.speciesGuessScientificName](requirements/local/product/product.md#speciesguessscientificname)
- [product.observedDateTime](requirements/local/product/product.md#observeddatetime)
- [product.observationFieldMapping](requirements/local/product/product.md#observationfieldmapping)
- [product.descriptionContents](requirements/local/product/product.md#descriptioncontents)
- [product.createBeforeMedia](requirements/local/product/product.md#createbeforemedia)
- [product.mediaFormats](requirements/local/product/product.md#mediaformats)
- [product.photoThenSoundFallback](requirements/local/product/product.md#photothensoundfallback)
- [product.mediaFilename](requirements/local/product/product.md#mediafilename)
- [product.descriptionUpdateKeepsMedia](requirements/local/product/product.md#descriptionupdatekeepsmedia)
- [product.additiveMediaResync](requirements/local/product/product.md#additivemediaresync)
- [product.endOfRunSummary](requirements/local/product/product.md#endofrunsummary)
- [product.conditionalSkipCounters](requirements/local/product/product.md#conditionalskipcounters)
- [product.conditionalFailureCount](requirements/local/product/product.md#conditionalfailurecount)
- [product.writeFailureAborts](requirements/local/product/product.md#writefailureaborts)
- [product.noRollback](requirements/local/product/product.md#norollback)
- [product.badRecordsSkipped](requirements/local/product/product.md#badrecordsskipped)
- [product.largeAccounts](requirements/local/product/product.md#largeaccounts)
- [product.unusableInputExplained](requirements/local/product/product.md#unusableinputexplained)
- [tech.csvByHeaderName](requirements/local/tech/tech.md#csvbyheadername)
- [tech.datesViaObserved](requirements/local/tech/tech.md#datesviaobserved)
- [tech.emptyNameNotIndexed](requirements/local/tech/tech.md#emptynamenotindexed)
- [tech.assetTypeNeedsDownload](requirements/local/tech/tech.md#assettypeneedsdownload)
- [tech.memoryScalesWithExport](requirements/local/tech/tech.md#memoryscaleswithexport)
- [tech.tempFilesDeleted](requirements/local/tech/tech.md#tempfilesdeleted)

### Candidate: accountsStaySafe

Users can run birdsync on their own iNaturalist account without losing data, duplicating observations, or putting the account's standing at risk.

Who benefits: Users. What they gain: an account they can trust birdsync with.

Would be cited by 26 requirements:

- [product.neverModifyOthersObservations](requirements/local/product/product.md#nevermodifyothersobservations)
- [product.toolsNotProduct](requirements/local/product/product.md#toolsnotproduct)
- [product.credentialsNeverLogged](requirements/local/product/product.md#credentialsneverlogged)
- [product.syncKey](requirements/local/product/product.md#synckey)
- [product.idempotentRerun](requirements/local/product/product.md#idempotentrerun)
- [product.taxonNotInSyncKey](requirements/local/product/product.md#taxonnotinsynckey)
- [product.incompleteKeyNotRecognized](requirements/local/product/product.md#incompletekeynotrecognized)
- [product.observationsPublic](requirements/local/product/product.md#observationspublic)
- [product.removedAssetsReported](requirements/local/product/product.md#removedassetsreported)
- [product.mediaCountMismatchReported](requirements/local/product/product.md#mediacountmismatchreported)
- [product.mediaFailuresTolerated](requirements/local/product/product.md#mediafailurestolerated)
- [product.dryRunNoWrites](requirements/local/product/product.md#dryrunnowrites)
- [product.dryRunLogPrefix](requirements/local/product/product.md#dryrunlogprefix)
- [product.dryRunMediaCount](requirements/local/product/product.md#dryrunmediacount)
- [product.dryRunCounterLabels](requirements/local/product/product.md#dryruncounterlabels)
- [product.noTaxonFilter](requirements/local/product/product.md#notaxonfilter)
- [product.permanentUploadFailureRecorded](requirements/local/product/product.md#permanentuploadfailurerecorded)
- [product.permanentFailuresReported](requirements/local/product/product.md#permanentfailuresreported)
- [tech.dryRunGateAtCallSite](requirements/local/tech/tech.md#dryrungateatcallsite)
- [tech.dryRunPrefixLogged](requirements/local/tech/tech.md#dryrunprefixlogged)
- [tech.honestDryRunCounters](requirements/local/tech/tech.md#honestdryruncounters)
- [tech.clientMutatesUnconditionally](requirements/local/tech/tech.md#clientmutatesunconditionally)
- [tech.updateKeepsPhotos](requirements/local/tech/tech.md#updatekeepsphotos)
- [tech.credentialsNeverExposed](requirements/local/tech/tech.md#credentialsneverexposed)
- [tech.toolsReadOnly](requirements/local/tech/tech.md#toolsreadonly)
- [tech.neverRunTools](requirements/local/tech/tech.md#neverruntools)

### Candidate: communityNotBurdened

iNaturalist's identifiers and operators bear as little cost as possible from observations birdsync posts.

Who benefits: The iNaturalist community and its operators. What they gain: synced records worth an identifier's time, and traffic that does not strain the service.

Would be cited by 8 requirements:

- [product.verifiableDefaultOn](requirements/local/product/product.md#verifiabledefaulton)
- [product.documentCommunityObligations](requirements/local/product/product.md#documentcommunityobligations)
- [product.documentSyncedObservationsIdentifiable](requirements/local/product/product.md#documentsyncedobservationsidentifiable)
- [tech.userAgent](requirements/local/tech/tech.md#useragent)
- [tech.maxPageSize](requirements/local/tech/tech.md#maxpagesize)
- [tech.errorsIncludeBody](requirements/local/tech/tech.md#errorsincludebody)
- [tech.requestPacing](requirements/local/tech/tech.md#requestpacing)
- [tech.idAbovePaging](requirements/local/tech/tech.md#idabovepaging)

### Candidate: maintainerChangesSafely

The maintainer and contributors can change birdsync quickly without breaking a user's account.

Who benefits: People who build birdsync. What they gain: cheap, safe changes.

Would be cited by 18 requirements:

- [tech.singleModule](requirements/local/tech/tech.md#singlemodule)
- [tech.oneExternalDependency](requirements/local/tech/tech.md#oneexternaldependency)
- [tech.goDirectivePolicy](requirements/local/tech/tech.md#godirectivepolicy)
- [tech.platformIndependent](requirements/local/tech/tech.md#platformindependent)
- [tech.noLiveServicesInTests](requirements/local/tech/tech.md#noliveservicesintests)
- [tech.clientInterfaces](requirements/local/tech/tech.md#clientinterfaces)
- [tech.baseURLSeams](requirements/local/tech/tech.md#baseurlseams)
- [tech.resetFlagsInTests](requirements/local/tech/tech.md#resetflagsintests)
- [tech.gofmtClean](requirements/local/tech/tech.md#gofmtclean)
- [tech.vetClean](requirements/local/tech/tech.md#vetclean)
- [tech.errorWrapping](requirements/local/tech/tech.md#errorwrapping)
- [tech.noLogFatalInLibraries](requirements/local/tech/tech.md#nologfatalinlibraries)
- [tech.progressLogging](requirements/local/tech/tech.md#progresslogging)
- [tech.commentsExplainWhy](requirements/local/tech/tech.md#commentsexplainwhy)
- [tech.ciStandingChecks](requirements/local/tech/tech.md#cistandingchecks)
- [tech.ciLatestStable](requirements/local/tech/tech.md#cilateststable)
- [tech.americanSpelling](requirements/local/tech/tech.md#americanspelling)
- [tech.quotationsNeverRespelled](requirements/local/tech/tech.md#quotationsneverrespelled)

Missing, by goals.md's checklist: none of the four is a goal for the communities birdsync
depends on *as such*, except `communityNotBurdened`. The Macaulay Library contributors whose
rights `mlterms` protects have no goal of their own, and their imported requirements cite no
goal (imported items carry no `Goals:` field). Consider a fifth goal, *contributors' rights in
their media are respected*.

## mutatingToolsDeleted

**Record the withdrawal of [tech.neverRunMutatingTools](requirements/local/tech/tech.md#neverrunmutatingtools) as a decision**

Drafted entry, for `requirements/decisions.md` if the owner confirms:

```
## mutatingToolsDeleted

**The mutating programs in tools/ were deleted, and the rule that nobody may run them was withdrawn**

Date: <when the tools were deleted>
Decided by: Sajmani
Status: accepted
Affects: tech.neverRunMutatingTools, tech.toolsReadOnly, tech.neverRunTools
Formerly: T-011 withdrawal note

### Decided

The mutating tools were deleted rather than governed by a rule nobody could check;
tech.toolsReadOnly and tech.neverRunTools replace it.
```

Specman recorded the withdrawal only as a `Status:` line on `T-011`, with no decision entry,
so the composed package cannot yet name the decision that excludes it.

## authenticateEveryRequestGap

**Accept as a gap, with a review date, that birdsync authenticates every request**

[inatapi.authenticateOnlyWhenNeeded](requirements/imported/inaturalist.org/inatapi/inatapi.md#authenticateonlywhenneeded)
is Governing, so only iNaturalist can grant an exception, and locally it can only be accepted as
an unmet gap. The reason specman recorded: an unauthenticated download would miss private or
obscured observations, and the next run would duplicate them. Proposed `Expires: 2027-11-03`.
The alternative is to ask iNaturalist for an exception.

## dailyCapGap

**Accept as a gap, with a review date, that the daily half of [birdsync.requestRate](requirements/composed/internal/birdsync/birdsync.md#requestrate) is not enforced**

birdsync keeps no state between runs, so it cannot count today's requests. Proposed
`Expires: 2027-11-03`.

## ownMediaReviewDate

**Give the accepted media-ownership assumption a review date**

[decisions.ownMediaAssumption](requirements/decisions.md#ownmediaassumption) (formerly CR-010) accepts that an export lists only the user's
own assets, with a reopen condition but no date. River requires an accepted gap to carry
`Expires:`. If the owner treats this as an accepted gap against [mlterms.noOthersMediaWithoutPermission](requirements/imported/ebird.org/mlterms/mlterms.md#noothersmediawithoutpermission),
a new entry amending it should add one. Proposed `Expires: 2027-11-03`.

## attributionGuidelines

**Vendor the Macaulay Library attribution guidelines, or accept [mlterms.attributionExpected](requirements/imported/ebird.org/mlterms/mlterms.md#attributionexpected) as unchecked**

The transcription notes that the guidelines page is not vendored, so whether `ML<id>`
filenames and asset URLs satisfy it is unknown. `imports.md` marks `mlterms` partial for this.

## uploadVolumeCap

**Decide whether [inatapi.singleIPAndMediaCaps](requirements/imported/inaturalist.org/inatapi/inatapi.md#singleipandmediacaps) limits uploads**

The source's 5 GB/hour and 24 GB/day caps name downloads. birdsync can upload gigabytes on a first
sync. Specman recorded the question rather than assuming it away.

## confirmDeciders

**Confirm who decided four migrated decisions**

`descriptionListsOnlyUploaded`, `extensionMapFromCDN`, `pageByIDAbove`, and
`notMachineGeneratedContent` (formerly CR-007, CR-009, CR-011, and CR-012) say "Resolved" in
the specman log without naming the decider. They are entered with the owner presumed, and say
so.

## newImports

**Decide on `uuid` and `github-actions` in `imports.md`**

Both were found by discovery and are marked `partial`. For `uuid` the license was read from
the module itself, and adopting it as Advisory is proposed. For `github-actions` no terms were
retrieved: someone needs to read GitHub's terms for Actions and the two actions' licenses, or
the owner records why they impose nothing.

## approveComposedPackage

**Approve `composed/internal/birdsync`, once the items above are settled**

`river hash river/requirements/composed/internal/birdsync` prints the line to record. Only
the owner adds it to `decisions.md`. Recompose first if any decision above changes a source.

## discoveryComplete

**Add `Discovery: complete` to `imports.md` once `newImports` is settled**

`river check` then fails, rather than notes, any unrecorded dependency or host.

## carriedOverOpenQuestions

**Open questions carried over from specman, so none drops out with `spec/`**

- Sound files over 50 MB are rejected by iNaturalist, and this surfaces as a generic upload failure. Should it be detected and reported specifically? (`spec/product.md`, open question 2)
- `inat.Client.UploadMedia` has no test. Its multipart request is exercised only through a mock (`spec/tech.md`, open question 1)
- The memory ceiling of [tech.memoryScalesWithExport](requirements/local/tech/tech.md#memoryscaleswithexport) has grown since [product.noTaxonFilter](requirements/local/product/product.md#notaxonfilter) and is unmeasured: worth a benchmark (`spec/tech.md`, open question 2)
- How much of `main` is worth testing: [product.onePositionalArgument](requirements/local/product/product.md#onepositionalargument) to [product.debugExplainsSkips](requirements/local/product/product.md#debugexplainsskips) live there untested (`spec/acceptance.md`, Gate 2 question 1)
- Whether to land the seven recommended additions in `spec/acceptance.md` now or as they become relevant: a token-in-logs check, error injection, testing `main` via `os/exec`, `UploadMedia` against `httptest`, a dependency count, and a large-export benchmark (Gate 2 question 2)
- Credentials-never-logged ([product.credentialsNeverLogged](requirements/local/product/product.md#credentialsneverlogged), [tech.credentialsNeverExposed](requirements/local/tech/tech.md#credentialsneverexposed)) is a security-relevant invariant with no check at all. A test asserting the token never appears in captured log output is cheap (`spec/acceptance.md`, gaps worth naming 1)
- The create and update failure paths are unexercised: no test sets the mock's `createObsErr` or `updateObsErr` ([product.writeFailureAborts](requirements/local/product/product.md#writefailureaborts); gaps worth naming 2)
- Order-dependent tests: nothing checks that each test calls `resetFlags()` ([tech.resetFlagsInTests](requirements/local/tech/tech.md#resetflagsintests); gaps worth naming 3)
- A successful sound upload has never been observed against the live service (`spec/decisions.md`, open work)
- A download failure is always transient, so a withdrawn asset is retried every run (`spec/decisions.md`, open work)
- Restoring an `Aves` plus `unknown` download filter is a pure bandwidth question, left open on purpose (`spec/decisions.md`)
- Issue #1 is fixed but left open until the reporter confirms on Windows (`spec/decisions.md`)
- [tech.honestDryRunCounters](requirements/local/tech/tech.md#honestdryruncounters) still carries specman's note "not yet satisfied", although `spec/tech.md` says it was implemented. The note was migrated verbatim; the owner may want it removed from the source

The originals are in `spec/` at [`dc25e3b`](https://github.com/Sajmani/birdsync/tree/dc25e3b6756bd59c2de931e6d080d47318617da4/spec).

## inlineRequirementIDs

**Decide what to do with the old IDs that remain in code comments and error messages**

The 43 `// Verifies:` annotations were harvested into `conformance.md` and removed. About 170
other mentions of `P-###`, `T-###`, and `CR-###` remain in comments and test failure messages,
for example `(P-063)`. Each still resolves through a `Formerly:` field. River says not to
annotate code with requirement IDs. Rewriting them is a separate change, proposed for the owner
to schedule.

## talkLinksPinned

**Confirm that the talk's links into `spec/` are pinned to the last specman commit**

`talks/2026-08-10-spec-composition.md` linked to `blob/main/spec/...`, which no longer exists.
Its links now point at `blob/dc25e3b…/spec/...`, where the talk's slides still
resolve. `TestTalkLinksResolve` checks only `blob/main/` links, so the pinned ones are no
longer checked, which is safe because a commit never changes.
