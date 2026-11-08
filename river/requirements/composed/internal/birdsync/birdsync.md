package birdsync

# birdsync: composed requirements

**Composed, not authored.** Change it by changing a source and recomposing, never by editing
this file. Composed on 2026-11-05 from every source in `river/requirements/local/` and
`river/requirements/imported/`, under [`precedence.md`](../../../precedence.md) and
[`decisions.md`](../../../decisions.md). The imported part was recomposed from the sources by an
agent with no context, following River's skill, with the local restatements of imported rules
withdrawn ([decisions.restatementsComposed](../../../decisions.md#restatementscomposed)). **Not
yet approved**: `decisions.md` carries no `Approves:` line for this package.

Requirements that apply unchanged are listed by area, linked to their source item. Only what
composition made is written out in full, under *Composed*. Requirements left out are under
*Excluded*, each naming the decision that excluded it. Acceptance criteria are in
[`acceptance.md`](acceptance.md).

Sources at composition: `product` and `tech` (local, at this commit), and `inatapi` (Governing),
`inatterms` and `mlterms` (Mandatory), vendored at the versions recorded in
[`imports.md`](../../../imports.md), the one place versions are recorded.

## Composed from

- inatapi (Governing)
- inatterms (Mandatory)
- mlterms (Mandatory)
- product (Local: product)
- tech (Local: tech)
- decisions
- precedence

## Included from product

### Purpose

- [product.copiesObservations](../../../local/product/scope.md#copiesobservations) — Copies eBird observations, with their media, into iNaturalist
- [product.oneWaySync](../../../local/product/scope.md#onewaysync) — The sync is one-way, eBird to iNaturalist
- [product.csvExportInput](../../../local/product/scope.md#csvexportinput) — The input is the user's eBird CSV export

### Non-goals

- [product.noReverseSync](../../../local/product/scope.md#noreversesync) — No sync from iNaturalist back to eBird
- [product.neverModifyOthersObservations](../../../local/product/scope.md#nevermodifyothersobservations) — Never modifies an observation birdsync did not create
- [product.noTaxonomyReconciliation](../../../local/product/scope.md#notaxonomyreconciliation) — No taxonomy reconciliation between the services
- [product.toolsNotProduct](../../../local/product/scope.md#toolsnotproduct) — Programs in tools/ are not part of the product

### Invocation

- [product.onePositionalArgument](../../../local/product/cli.md#onepositionalargument) — Exactly one positional argument
- [product.flagsBeforeArgument](../../../local/product/cli.md#flagsbeforeargument) — Flags precede the positional argument
- [product.boolFlagEqualsForm](../../../local/product/cli.md#boolflagequalsform) — Boolean flags are disabled with --flag=false
- [product.unopenableCsvExitsEarly](../../../local/product/cli.md#unopenablecsvexitsearly) — An unopenable CSV exits before contacting iNaturalist
- [product.afterLaterThanBeforeExits](../../../local/product/cli.md#afterlaterthanbeforeexits) — --after later than --before is an error
- [product.debugExplainsSkips](../../../local/product/cli.md#debugexplainsskips) — --debug logs why each record was skipped

### Authentication

- [product.credentialsFromEnvironment](../../../local/product/auth.md#credentialsfromenvironment) — Credentials come from INAT_USER_ID and INAT_API_TOKEN
- [product.promptForMissingCredentials](../../../local/product/auth.md#promptformissingcredentials) — Missing credentials are prompted for
- [product.tokenFraming](../../../local/product/auth.md#tokenframing) — Token framing differs between prompt and environment
- [product.unauthorizedSaysRefreshToken](../../../local/product/auth.md#unauthorizedsaysrefreshtoken) — A 401 says the token needs refreshing
- [product.credentialsNeverLogged](../../../local/product/auth.md#credentialsneverlogged) — Credentials are never written to the log

### Identity and idempotence

- [product.syncKey](../../../local/product/identity.md#synckey) — Own observations are identified by a two-field sync key
- [product.idempotentRerun](../../../local/product/identity.md#idempotentrerun) — Re-running creates no duplicates
- [product.taxonNotInSyncKey](../../../local/product/identity.md#taxonnotinsynckey) — The taxon is not part of the sync key
- [product.incompleteKeyNotRecognized](../../../local/product/identity.md#incompletekeynotrecognized) — An incomplete sync key is not recognized

### Downloading existing observations

- [product.downloadExistingObservations](../../../local/product/download.md#downloadexistingobservations) — Existing observations are downloaded before syncing
- [product.dateWindowNarrowsDownload](../../../local/product/download.md#datewindownarrowsdownload) — A date window narrows the download

### Filtering and skip order

- [product.skipOrder](../../../local/product/filtering.md#skiporder) — Records are tested in a fixed skip order
- [product.afterFlag](../../../local/product/filtering.md#afterflag) — --after skips records observed before it
- [product.beforeFlag](../../../local/product/filtering.md#beforeflag) — --before skips records observed after it
- [product.dateFlagFormats](../../../local/product/filtering.md#dateflagformats) — Date flags accept two formats
- [product.verifiableDefaultOn](../../../local/product/filtering.md#verifiabledefaulton) — --verifiable skips medialess records, on by default
- [product.fuzzyMatching](../../../local/product/filtering.md#fuzzymatching) — --fuzzy skips records matching an unkeyed observation
- [product.fuzzyNeverEmptyName](../../../local/product/filtering.md#fuzzyneveremptyname) — Fuzzy matching never uses an empty name
- [product.fuzzyParsedDates](../../../local/product/filtering.md#fuzzyparseddates) — Fuzzy matching compares parsed dates
- [product.fuzzyOffByDefault](../../../local/product/filtering.md#fuzzyoffbydefault) — --fuzzy is off by default and its risk documented

### What birdsync writes

- [product.createdObservationWild](../../../local/product/observations.md#createdobservationwild) — Created observations are wild
- [product.checklistCoordinates](../../../local/product/observations.md#checklistcoordinates) — Coordinates are the checklist's, inexact, with set accuracy
- [product.speciesGuessScientificName](../../../local/product/observations.md#speciesguessscientificname) — The species guess is the eBird scientific name
- [product.observedDateTime](../../../local/product/observations.md#observeddatetime) — Observation date and time come from the eBird columns
- [product.observationFieldMapping](../../../local/product/observations.md#observationfieldmapping) — eBird columns are copied into observation fields
- [product.descriptionContents](../../../local/product/observations.md#descriptioncontents) — What the observation description contains
- [product.observationsPublic](../../../local/product/observations.md#observationspublic) — Created observations are public

### Media

- [product.createBeforeMedia](../../../local/product/media.md#createbeforemedia) — The observation is created before media is attached
- [product.mediaFormats](../../../local/product/media.md#mediaformats) — Photos at 2400px, sounds as MP3
- [product.photoThenSoundFallback](../../../local/product/media.md#photothensoundfallback) — Try the photo URL, then the sound URL
- [product.mediaFilename](../../../local/product/media.md#mediafilename) — Media is uploaded as ML<asset ID> with a canonical extension
- [product.descriptionUpdateKeepsMedia](../../../local/product/media.md#descriptionupdatekeepsmedia) — The description update keeps attached media
- [product.additiveMediaResync](../../../local/product/media.md#additivemediaresync) — Media re-sync is additive
- [product.removedAssetsReported](../../../local/product/media.md#removedassetsreported) — Assets removed from eBird are reported, not removed
- [product.mediaCountMismatchReported](../../../local/product/media.md#mediacountmismatchreported) — A media count mismatch is reported, not corrected
- [product.mediaFailuresTolerated](../../../local/product/media.md#mediafailurestolerated) — A failed media transfer is logged, counted, and survived

### Dry run

- [product.dryRunNoWrites](../../../local/product/dryrun.md#dryrunnowrites) — --dryrun issues no writes
- [product.dryRunLogPrefix](../../../local/product/dryrun.md#dryrunlogprefix) — Skipped actions are logged with DRYRUN:
- [product.dryRunMediaCount](../../../local/product/dryrun.md#dryrunmediacount) — A dry run reports one unclassified media count

### Reporting

- [product.endOfRunSummary](../../../local/product/reporting.md#endofrunsummary) — Every run ends with a summary
- [product.conditionalSkipCounters](../../../local/product/reporting.md#conditionalskipcounters) — A skip counter prints only when its flag is in effect
- [product.conditionalFailureCount](../../../local/product/reporting.md#conditionalfailurecount) — The media failure count prints only when non-zero

### Failure behavior

- [product.writeFailureAborts](../../../local/product/failures.md#writefailureaborts) — A failed create or update aborts the run
- [product.noRollback](../../../local/product/failures.md#norollback) — There is no rollback
- [product.largeAccounts](../../../local/product/download.md#largeaccounts) — Works for accounts over 10,000 observations
- [product.unusableInputExplained](../../../local/product/failures.md#unusableinputexplained) — Unusable input is named before parsing
- [product.permanentUploadFailureRecorded](../../../local/product/failures.md#permanentuploadfailurerecorded) — A permanent upload failure is recorded and not retried
- [product.permanentFailuresReported](../../../local/product/failures.md#permanentfailuresreported) — Permanently failed assets are reported each run
- [product.documentSyncedObservationsIdentifiable](../../../local/product/community.md#documentsyncedobservationsidentifiable) — Documentation says synced observations are identifiable

### Reporting

- [product.dryRunCounterLabels](../../../local/product/reporting.md#dryruncounterlabels) — Dry-run observation counters are labeled hypothetical

### Downloading existing observations

- [product.noTaxonFilter](../../../local/product/download.md#notaxonfilter) — The download is not filtered by taxon

### Failure behavior

- [product.badRecordsSkipped](../../../local/product/failures.md#badrecordsskipped) — Unparseable records are skipped, logged, and counted

## Included from tech

### Module and dependencies

- [tech.singleModule](../../../local/tech/deps.md#singlemodule) — A single Go module with no build system
- [tech.oneExternalDependency](../../../local/tech/deps.md#oneexternaldependency) — Exactly one external dependency
- [tech.goDirectivePolicy](../../../local/tech/deps.md#godirectivepolicy) — go and toolchain directives change only deliberately
- [tech.platformIndependent](../../../local/tech/deps.md#platformindependent) — Builds and runs on every Go platform

### Safety invariants

- [tech.dryRunGateAtCallSite](../../../local/tech/safety.md#dryrungateatcallsite) — Mutations are gated on --dryrun at the call site
- [tech.dryRunPrefixLogged](../../../local/tech/safety.md#dryrunprefixlogged) — A skipped mutation is logged with DRYRUN:
- [tech.honestDryRunCounters](../../../local/tech/safety.md#honestdryruncounters) — Dry-run counters never report work that did not happen
- [tech.clientMutatesUnconditionally](../../../local/tech/safety.md#clientmutatesunconditionally) — inat.Client mutates unconditionally
- [tech.updateKeepsPhotos](../../../local/tech/safety.md#updatekeepsphotos) — UpdateObservation always sets ignore_photos
- [tech.noLiveServicesInTests](../../../local/tech/testing.md#noliveservicesintests) — Tests never contact a live service
- [tech.toolsReadOnly](../../../local/tech/safety.md#toolsreadonly) — No program in tools/ mutates
- [tech.neverRunTools](../../../local/tech/safety.md#neverruntools) — Nothing in the workflow runs a program in tools/
- [tech.credentialsNeverExposed](../../../local/tech/safety.md#credentialsneverexposed) — Credentials are never logged, echoed, or written

### Testability seams

- [tech.clientInterfaces](../../../local/tech/testing.md#clientinterfaces) — birdsync() takes its clients as interfaces
- [tech.baseURLSeams](../../../local/tech/testing.md#baseurlseams) — Both services have base-URL test seams
- [tech.resetFlagsInTests](../../../local/tech/testing.md#resetflagsintests) — Tests call resetFlags first

### External service etiquette

- [tech.errorsIncludeBody](../../../local/tech/etiquette.md#errorsincludebody) — API errors include the response body

### Data format handling

- [tech.csvByHeaderName](../../../local/tech/formats.md#csvbyheadername) — The CSV is read by header name
- [tech.datesViaObserved](../../../local/tech/formats.md#datesviaobserved) — eBird dates are compared via Record.Observed()
- [tech.emptyNameNotIndexed](../../../local/tech/formats.md#emptynamenotindexed) — An empty taxon name never enters the fuzzy index
- [tech.assetTypeNeedsDownload](../../../local/tech/formats.md#assettypeneedsdownload) — Asset type is known only after download

### Resource use

- [tech.memoryScalesWithExport](../../../local/tech/resources.md#memoryscaleswithexport) — Memory scales with the export
- [tech.tempFilesDeleted](../../../local/tech/resources.md#tempfilesdeleted) — Temporary media files are deleted

### Code conventions

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

## Included from imported sources

Imported requirements that apply to birdsync as written, needing nothing birdsync-specific composed from them.

### iNaturalist API

- [inatapi.singleIPAndMediaCaps](../../../imported/inaturalist.org/inatapi/inatapi.md#singleipandmediacaps) — One IP; media download volume caps

### iNaturalist terms

- [inatterms.contributorKeepsRights](../../../imported/inaturalist.org/inatterms/inatterms.md#contributorkeepsrights) — A contributor keeps their rights, and licenses by default as CC BY-NC

### Macaulay Library terms

- [mlterms.contributorKeepsCopyright](../../../imported/ebird.org/mlterms/mlterms.md#contributorkeepscopyright) — A contributor keeps copyright in their own media

## Composed
## requestRate

**Requests to the iNaturalist API stay within iNaturalist's published rate**

Composed from: [inatapi.requestRate](../../../imported/inaturalist.org/inatapi/inatapi.md#requestrate),
[inatterms.rateAndVolume](../../../imported/inaturalist.org/inatterms/inatterms.md#rateandvolume),
[inatterms.dontBehaveLikeMachine](../../../imported/inaturalist.org/inatterms/inatterms.md#dontbehavelikemachine),
[tech.dailyRequestBudget](../../../local/tech/etiquette.md#dailyrequestbudget),
[decisions.dailyRequestBudget](../../../decisions.md#dailyrequestbudget)
Goals: [goals.communityNotBurdened](../../../../goals/goals.md#communitynotburdened), [goals.accountsStaySafe](../../../../goals/goals.md#accountsstaysafe)

birdsync sends requests to `api.inaturalist.org` at no more than about one per second, and no
more than 9,000 in a rolling 24-hour window (within iNaturalist's guideline of around 10,000 a
day), whatever the size of the export and in every mode, including `--dryrun`.

The terms (Mandatory) require posting at a rate that does not hinder other users and warn
against adding a lot of content very quickly; they give no number. The API's recommended
practices (Governing) give the numbers. Composition reads iNaturalist's own published limit as
the measure of the Mandatory rule, so neither the per-second nor the daily half can be given up
locally. `inat.Client` enforces the one-second spacing and a 9,000-request cap per run at the
single choke point every iNaturalist request passes through (`docs/arch.md`), and `birdsync()`
estimates prior 24-hour write requests from the `created_at` timestamps and media of the
downloaded observations so back-to-back runs stay within the same 9,000-request budget, stopping
cleanly at an observation boundary with instructions to re-run in 24 hours
([decisions.dailyRequestBudget](../../../decisions.md#dailyrequestbudget)).
## observationDownloadPageSize

**The existing-observation download asks for 200 observations per page**

Composed from: [inatapi.largestPageSize](../../../imported/inaturalist.org/inatapi/inatapi.md#largestpagesize)
Goals: [goals.communityNotBurdened](../../../../goals/goals.md#communitynotburdened)

Every page request in the download of the user's existing iNaturalist observations
([product.downloadExistingObservations](../../../local/product/download.md#downloadexistingobservations))
sets `per_page=200`, the largest page the API serves for observations.

The value is the source's own ("you can get up to 200 observations in a single request"). This
replaces the withdrawn local [tech.maxPageSize](../../../local/tech/etiquette.md#maxpagesize)
([decisions.restatementsComposed](../../../decisions.md#restatementscomposed)).
## observationDownloadByIDAbove

**The existing-observation download pages by id_above, in ascending id order**

Composed from: [inatapi.idAbovePastTenThousand](../../../imported/inaturalist.org/inatapi/inatapi.md#idabovepasttenthousand),
[decisions.pageByIDAbove](../../../decisions.md#pagebyidabove)
Goals: [goals.birdersKeepOneRecord](../../../../goals/goals.md#birderskeeponerecord), [goals.accountsStaySafe](../../../../goals/goals.md#accountsstaysafe)

The download of the user's existing observations
([product.downloadExistingObservations](../../../local/product/download.md#downloadexistingobservations))
sorts by id ascending (`order_by=id`, `order=asc`), sets `id_above` on each request after the
first to the highest observation ID in the page before, never pages by `page`, ends at the first
page shorter than the page size, and fails with an error if a page does not advance the cursor.

This is how birdsync meets [product.largeAccounts](../../../local/product/download.md#largeaccounts)
within the API's 10,000-result ceiling on page numbers, as decided in
[decisions.pageByIDAbove](../../../decisions.md#pagebyidabove). It replaces the withdrawn local
[tech.idAbovePaging](../../../local/tech/etiquette.md#idabovepaging) ([decisions.restatementsComposed](../../../decisions.md#restatementscomposed)).
## inatUserAgent

**Requests to the iNaturalist API identify birdsync in the User-Agent**

Composed from: [inatapi.identifyUserAgent](../../../imported/inaturalist.org/inatapi/inatapi.md#identifyuseragent),
[tech.userAgent](../../../local/tech/etiquette.md#useragent)
Goals: [goals.communityNotBurdened](../../../../goals/goals.md#communitynotburdened)

Every request birdsync sends to `api.inaturalist.org` carries a `User-Agent` header that names
birdsync and its version, rather than the Go HTTP client's default.

The source asks clients to "consider" a custom User-Agent; [`imports.md`](../../../imports.md)
adopts all eight `inatapi` items, so composition states it as a requirement over every request.
The version is the local addition of [tech.userAgent](../../../local/tech/etiquette.md#useragent). The source allows the app name, an iNaturalist
username, "or anything we might use to differentiate your requests"; the app name is the one
that holds for every user.
## jwtInAuthorizationHeader

**birdsync authenticates with the user's JWT in the Authorization header**

Composed from: [inatapi.jwtAuthorization](../../../imported/inaturalist.org/inatapi/inatapi.md#jwtauthorization)
Goals: [goals.birdersKeepOneRecord](../../../../goals/goals.md#birderskeeponerecord), [goals.accountsStaySafe](../../../../goals/goals.md#accountsstaysafe)

Every authenticated request birdsync sends to `api.inaturalist.org` carries, in its HTTP
`Authorization` header, the JWT the user obtained from `https://www.inaturalist.org/users/api_token`,
whether it came from `INAT_API_TOKEN` as the bare token or from the prompt in its
`{"api_token":"TOKEN"}` framing
([product.credentialsFromEnvironment](../../../local/product/auth.md#credentialsfromenvironment),
[product.promptForMissingCredentials](../../../local/product/auth.md#promptformissingcredentials),
[product.tokenFraming](../../../local/product/auth.md#tokenframing)).

The local sources say where the token comes from and how a 401 from its 24-hour expiry is
reported ([product.unauthorizedSaysRefreshToken](../../../local/product/auth.md#unauthorizedsaysrefreshtoken));
none says how it is sent, which is what the Governing source prescribes.
## ownAccountOnly

**birdsync posts only as the person running it, with their own token**

Composed from: [inatapi.noGroupAccounts](../../../imported/inaturalist.org/inatapi/inatapi.md#nogroupaccounts),
[inatterms.accountHolderResponsible](../../../imported/inaturalist.org/inatterms/inatterms.md#accountholderresponsible),
[decisions.inputIdentityIsTheAccountHolders](../../../decisions.md#inputidentityistheaccountholders)
Goals: [goals.accountsStaySafe](../../../../goals/goals.md#accountsstaysafe)

Every write birdsync makes to iNaturalist
([product.copiesObservations](../../../local/product/scope.md#copiesobservations)) is
authenticated with a token supplied at run time by the person running it
([product.credentialsFromEnvironment](../../../local/product/auth.md#credentialsfromenvironment),
[product.promptForMissingCredentials](../../../local/product/auth.md#promptformissingcredentials));
birdsync contains, ships, and falls back to no account or token of its own.

That keeps birdsync from being a group account: each user posts as themselves. Whether the
person running birdsync is the person whose eBird export it reads is the account holder's
responsibility, not birdsync's
([decisions.inputIdentityIsTheAccountHolders](../../../decisions.md#inputidentityistheaccountholders)).
The source also says to "use one of the various OAuth authentication flows"; whether a token
the user pastes from the api_token page meets that is not settled by precedence
(an open question in `river/proposals.md`).
## authenticateOnlyWhenNecessary

**birdsync sends its token to iNaturalist only on requests that need it**

Composed from: [inatapi.authenticateOnlyWhenNeeded](../../../imported/inaturalist.org/inatapi/inatapi.md#authenticateonlywhenneeded),
[inatterms.duplicatesUndesirable](../../../imported/inaturalist.org/inatterms/inatterms.md#duplicatesundesirable),
[decisions.authenticationIsNecessary](../../../decisions.md#authenticationisnecessary)
Goals: [goals.communityNotBurdened](../../../../goals/goals.md#communitynotburdened), [goals.accountsStaySafe](../../../../goals/goals.md#accountsstaysafe)

birdsync sends its token to `api.inaturalist.org` only on a write, or on a read whose
unauthenticated answer would omit observations birdsync must see to avoid creating duplicates;
every other request to iNaturalist is sent without it.

The owner ruled in [decisions.authenticationIsNecessary](../../../decisions.md#authenticationisnecessary)
that the download of existing observations
([product.downloadExistingObservations](../../../local/product/download.md#downloadexistingobservations))
is such a read: without the token it omits private and geoprivacy-obscured observations, each of
which the next run would create again ([product.idempotentRerun](../../../local/product/identity.md#idempotentrerun)).
So the rule is met today with every request authenticated, and needs no exception from iNaturalist.
Media downloads go to the Macaulay Library and are outside it.
## postsOnlyUsersOwnRecords

**Every observation birdsync posts is a record from the user's own export, and it posts nothing else of its own**

Composed from: [inatterms.noMachineGeneratedContent](../../../imported/inaturalist.org/inatterms/inatterms.md#nomachinegeneratedcontent),
[decisions.notMachineGeneratedContent](../../../decisions.md#notmachinegeneratedcontent)
Goals: [goals.accountsStaySafe](../../../../goals/goals.md#accountsstaysafe), [goals.communityNotBurdened](../../../../goals/goals.md#communitynotburdened)

Every observation birdsync creates
([product.copiesObservations](../../../local/product/scope.md#copiesobservations),
[product.syncKey](../../../local/product/identity.md#synckey)) corresponds to one record in the
eBird export the user supplied ([product.csvExportInput](../../../local/product/scope.md#csvexportinput)),
and birdsync posts no identification, comment, or message of its own.

[decisions.notMachineGeneratedContent](../../../decisions.md#notmachinegeneratedcontent) holds
that birdsync is not machine-generated content because a human decided each observation and the
machine only facilitates posting it, from the user's own curated data. That holds only while
each observation is one the user recorded, and while birdsync generates no other kind of
content; the Mandatory rule, whose penalty is suspension without warning, is met by keeping
both true.
## accountHolderToldOfObligations

**The documentation tells users they answer for what birdsync posts under their account**

Composed from: [inatterms.accountHolderResponsible](../../../imported/inaturalist.org/inatterms/inatterms.md#accountholderresponsible),
[inatterms.dontBehaveLikeMachine](../../../imported/inaturalist.org/inatterms/inatterms.md#dontbehavelikemachine),
[decisions.notMachineGeneratedContent](../../../decisions.md#notmachinegeneratedcontent),
[product.documentCommunityObligations](../../../local/product/community.md#documentcommunityobligations)
Goals: [goals.accountsStaySafe](../../../../goals/goals.md#accountsstaysafe), [goals.communityNotBurdened](../../../../goals/goals.md#communitynotburdened)

birdsync's documentation tells users that iNaturalist holds them responsible for everything
birdsync posts under their account; that they should review the observations birdsync created;
and that an account which adds a lot of content quickly and does not respond to comments and
messages may be suspended, so they should answer comments and messages on synced observations.

birdsync cannot answer comments for the user; the part of
[inatterms.dontBehaveLikeMachine](../../../imported/inaturalist.org/inatterms/inatterms.md#dontbehavelikemachine)
it can reach is pacing ([requestRate](#requestrate)), accuracy
([product.verifiableDefaultOn](../../../local/product/filtering.md#verifiabledefaulton)),
identifiability ([product.documentSyncedObservationsIdentifiable](../../../local/product/community.md#documentsyncedobservationsidentifiable)),
and telling the user. Reviewing what was created matters most: [inatterms.noMachineGeneratedContent](../../../imported/inaturalist.org/inatterms/inatterms.md#nomachinegeneratedcontent)
permits a script to post only where a human decided each observation, so a user who reviews
nothing moves birdsync to the prohibited side of that rule
([product.documentCommunityObligations](../../../local/product/community.md#documentcommunityobligations)).
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

## Excluded

- [product.avesOnlyDownload](../../../local/product/download.md#avesonlydownload) *(excluded)* — The download is restricted to Aves. Excluded by [decisions.noTaxonFilterOnDownload](../../../decisions.md#notaxonfilterondownload)
- [product.malformedRecordAborts](../../../local/product/failures.md#malformedrecordaborts) *(excluded)* — A malformed record aborts the run. Excluded by [decisions.badRecordsSkippedNotFatal](../../../decisions.md#badrecordsskippednotfatal)
- [tech.neverRunMutatingTools](../../../local/tech/safety.md#neverrunmutatingtools) *(excluded)* — Nothing runs the mutating tools. Excluded by [decisions.mutatingToolsDeleted](../../../decisions.md#mutatingtoolsdeleted)
- [tech.maxPageSize](../../../local/tech/etiquette.md#maxpagesize) *(excluded)* — Downloads request the maximum page size. Excluded by [decisions.restatementsComposed](../../../decisions.md#restatementscomposed)
- [tech.requestPacing](../../../local/tech/etiquette.md#requestpacing) *(excluded)* — iNaturalist requests are paced. Excluded by [decisions.restatementsComposed](../../../decisions.md#restatementscomposed)
- [tech.idAbovePaging](../../../local/tech/etiquette.md#idabovepaging) *(excluded)* — The download pages by id_above. Excluded by [decisions.restatementsComposed](../../../decisions.md#restatementscomposed)
