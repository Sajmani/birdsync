package product

What birdsync must do, from a user's point of view: one item per heading, and each item's
first paragraph after its fields is the complete rule. Search `^## ` for the list, or
`## <name>` for one item; `Formerly:` gives the identifier this item had under the
predecessor method (specman), which older commits, comments, and issues still cite.

Migrated from `spec/product.md` at
[`dc25e3b`](https://github.com/Sajmani/birdsync/blob/dc25e3b6756bd59c2de931e6d080d47318617da4/spec/product.md),
where it was approved at specman's Gate 1 on 2026-08-09 and amended since. The
requirements were recovered from the implementation, README.md and arch.md, and describe
behavior as it was then. No goals have been accepted yet, so every `Goals:` field is empty.

## copiesObservations

**Copies eBird observations, with their media, into iNaturalist**

Formerly: P-001
Area: Purpose
Goals:

birdsync copies bird observations from a user's eBird data export into that
user's iNaturalist account, including the photos and sounds attached to them in the
Macaulay Library.

## oneWaySync

**The sync is one-way, eBird to iNaturalist**

Formerly: P-002
Area: Purpose
Goals:

The sync is one-way: eBird → iNaturalist.

## csvExportInput

**The input is the user's eBird CSV export**

Formerly: P-003
Area: Purpose
Goals:

The input is a `MyEBirdData.csv` file the user downloads from eBird, not a
live eBird API.

*Rationale: the eBird API does not expose personal checklists.*

## noReverseSync

**No sync from iNaturalist back to eBird**

Formerly: P-004
Area: Non-goals
Goals:

birdsync does not sync iNaturalist observations back to eBird.

## neverModifyOthersObservations

**Never modifies an observation birdsync did not create**

Formerly: P-005
Area: Non-goals
Goals:

birdsync never modifies or deletes an iNaturalist observation it did not
create, except to attach media to an observation carrying its own sync key.

## noTaxonomyReconciliation

**No taxonomy reconciliation between the services**

Formerly: P-006
Area: Non-goals
Goals:

birdsync does not reconcile taxonomy between the two services. An eBird
scientific name that iNaturalist cannot resolve produces an observation with an unknown
taxon, which the user fixes by hand.

## toolsNotProduct

**Programs in tools/ are not part of the product**

Formerly: P-007
Area: Non-goals
Goals:

The programs in `tools/` are not part of the product. They are maintenance
utilities, are not installed by `go install` of the root package, and carry no
compatibility promise. They are also read-only ([tech.toolsReadOnly](../tech/tech.md#toolsreadonly)): no tool may alter a user's
account.

## onePositionalArgument

**Exactly one positional argument**

Formerly: P-008
Area: Invocation
Goals:

birdsync takes exactly one positional argument: the path to the eBird CSV
export. Any other number of arguments prints usage and exits non-zero.

## flagsBeforeArgument

**Flags precede the positional argument**

Formerly: P-009
Area: Invocation
Goals:

Flags must appear before the positional argument.

*Rationale: consequence of the standard Go flag package; documented rather than fixed.*

## boolFlagEqualsForm

**Boolean flags are disabled with --flag=false**

Formerly: P-010
Area: Invocation
Goals:

A boolean flag is disabled with `--flag=false`; `--flag false` is a usage
error.

*Rationale: same as [flagsBeforeArgument](#flagsbeforeargument).*

## unopenableCsvExitsEarly

**An unopenable CSV exits before contacting iNaturalist**

Formerly: P-011
Area: Invocation
Goals:

If the CSV file cannot be opened, birdsync exits with an error before
contacting iNaturalist or prompting for credentials.

## afterLaterThanBeforeExits

**--after later than --before is an error**

Formerly: P-012
Area: Invocation
Goals:

If `--after` is later than `--before`, birdsync exits with an error rather
than running a sync that cannot match anything.

## debugExplainsSkips

**--debug logs why each record was skipped**

Formerly: P-013
Area: Invocation
Goals:

`--debug` logs the reason each eBird record was skipped, and other detail
not shown in a normal run.

## credentialsFromEnvironment

**Credentials come from INAT_USER_ID and INAT_API_TOKEN**

Formerly: P-014
Area: Authentication
Goals:

The iNaturalist user ID is read from `INAT_USER_ID`, and the API token from
`INAT_API_TOKEN`.

## promptForMissingCredentials

**Missing credentials are prompted for**

Formerly: P-015
Area: Authentication
Goals:

When either is unset, birdsync prompts for it interactively and explains
where to find it.

## tokenFraming

**Token framing differs between prompt and environment**

Formerly: P-016
Area: Authentication
Goals:

At the prompt the token is accepted in its `{"api_token":"TOKEN"}` framing;
from the environment it is the bare token.

*Rationale: matches what the iNaturalist page gives the user, versus what is convenient
to export in a shell.*

## unauthorizedSaysRefreshToken

**A 401 says the token needs refreshing**

Formerly: P-017
Area: Authentication
Goals:

A 401 from the API reports that the token needs refreshing, and where to get
a new one, rather than a bare HTTP status.

*Rationale: tokens expire every 24 hours, so this is the most likely failure a user
hits.*

## credentialsNeverLogged

**Credentials are never written to the log**

Formerly: P-018
Area: Authentication
Goals:

Credentials are never written to the log.

## syncKey

**Own observations are identified by a two-field sync key**

Formerly: P-019
Area: Identity and idempotence
Subject: `sync.key.fields`
Value: `[6033, 20215]`
Goals:

birdsync identifies its own observations by a sync key: the pair of
iNaturalist observation fields *eBird Checklist* (6033) and *eBird Scientific Name*
(20215).

## idempotentRerun

**Re-running creates no duplicates**

Formerly: P-020
Area: Identity and idempotence
Goals:

Re-running birdsync over the same CSV creates no duplicate observations.

## taxonNotInSyncKey

**The taxon is not part of the sync key**

Formerly: P-021
Area: Identity and idempotence
Goals:

The iNaturalist taxon is not part of the sync key.

*Rationale: the taxon may be changed by the user or the community after upload, and
eBird's "slash" (`Aythya marila/affinis`) and "spuh" (`Melanitta sp.`) names have no
exact iNaturalist equivalent.*

## incompleteKeyNotRecognized

**An incomplete sync key is not recognized**

Formerly: P-022
Area: Identity and idempotence
Goals:

An observation missing either sync-key field is not recognized as
birdsync-created.

*Rationale: observations created by versions that set only the checklist field fall into
this population. A `repair` tool existed to backfill them and was deleted once the
maintainer's account was clean; it is recoverable from the history.*

## downloadExistingObservations

**Existing observations are downloaded before syncing**

Formerly: P-023
Area: Downloading existing observations
Goals:

Before syncing, birdsync downloads the user's existing iNaturalist
observations into memory.

## avesOnlyDownload

**The download is restricted to Aves**

Formerly: P-024
Area: Downloading existing observations
Status: withdrawn by [decisions.noTaxonFilterOnDownload](../../decisions.md#notaxonfilterondownload) — it defeated [idempotentRerun](#idempotentrerun). Superseded by [noTaxonFilter](#notaxonfilter).
Goals:

The download is restricted to the `Aves` iconic taxon.

## dateWindowNarrowsDownload

**A date window narrows the download**

Formerly: P-025
Area: Downloading existing observations
Goals:

When `--after` or `--before` is set, the download is restricted to that date
window, so duplicate detection and fuzzy matching only consider observations inside it.

## skipOrder

**Records are tested in a fixed skip order**

Formerly: P-026
Area: Filtering and skip order
Goals:

Each eBird record is tested in this order: `--after`, `--before`,
already-synced, `--fuzzy`, `--verifiable`. A record surviving all five becomes a new
iNaturalist observation.

*Rationale: the order is user-visible, because it determines which counter a skipped
record lands in.*

## afterFlag

**--after skips records observed before it**

Formerly: P-027
Area: Filtering and skip order
Goals:

`--after` skips records observed strictly before the given time.

## beforeFlag

**--before skips records observed after it**

Formerly: P-028
Area: Filtering and skip order
Goals:

`--before` skips records observed strictly after the given time.

## dateFlagFormats

**Date flags accept two formats**

Formerly: P-029
Area: Filtering and skip order
Goals:

Both accept `2006-01-02 15:04:05` or `2006-01-02`.

## verifiableDefaultOn

**--verifiable skips medialess records, on by default**

Formerly: P-030
Area: Filtering and skip order
Subject: `cli.verifiable.default`
Value: `true`
Goals:

`--verifiable` skips records with no Macaulay Library catalog numbers, and
defaults to **on**.

*Rationale: iNaturalist will not treat a medialess observation as verifiable, so the default
avoids filling a user's account with "Casual" records.*

*It is also the main thing birdsync does to limit the burden it puts on other people. A
medialess observation gives an identifier nothing to identify, so it is pure cost to the
community — the concern a forum moderator raised about volume
([decisions.notMachineGeneratedContent](../../decisions.md#notmachinegeneratedcontent)). Don't flip this
default for convenience.*

## fuzzyMatching

**--fuzzy skips records matching an unkeyed observation**

Formerly: P-031
Area: Filtering and skip order
Goals:

`--fuzzy` skips a record when an existing observation not carrying a sync key
has the same observation date and the same name, matched against either the common name
or the scientific name.

## fuzzyNeverEmptyName

**Fuzzy matching never uses an empty name**

Formerly: P-032
Area: Filtering and skip order
Goals:

Fuzzy matching never uses an empty name.

*Rationale: an empty name matches every unnamed record on that date and silently drops
legitimate observations.*

## fuzzyParsedDates

**Fuzzy matching compares parsed dates**

Formerly: P-033
Area: Filtering and skip order
Goals:

Fuzzy matching compares parsed observation dates, not raw CSV date strings.

*Rationale: eBird writes dates in two formats; keying on the raw field silently disabled
`--fuzzy` for anyone whose export used the second one.*

## fuzzyOffByDefault

**--fuzzy is off by default and its risk documented**

Formerly: P-034
Area: Filtering and skip order
Goals:

`--fuzzy` is off by default, and its documentation states that it may
suppress a legitimate record when the same bird was seen twice in one day.

## createdObservationWild

**Created observations are wild**

Formerly: P-035
Area: What birdsync writes
Goals:

A created observation is marked wild, not captive.

*Rationale: eBird checklists record wild birds.*

## checklistCoordinates

**Coordinates are the checklist's, inexact, with set accuracy**

Formerly: P-036
Area: What birdsync writes
Subject: `observation.positional_accuracy.default_m`
Value: `1000`
Goals:

Latitude and longitude are the checklist's, the location is marked
not-exact, and positional accuracy is `--positional_accuracy_meters`, default 1000.

*Rationale: the checklist location approximates a hotspot, not the bird.*

## speciesGuessScientificName

**The species guess is the eBird scientific name**

Formerly: P-037
Area: What birdsync writes
Goals:

The species guess is the eBird scientific name.

## observedDateTime

**Observation date and time come from the eBird columns**

Formerly: P-038
Area: What birdsync writes
Goals:

The observation date and time come from the eBird `Date` and `Time` columns.

## observationFieldMapping

**eBird columns are copied into observation fields**

Formerly: P-039
Area: What birdsync writes
Goals:

These eBird columns are copied into iNaturalist observation fields: Count
(1), Common Name (256), Location (157), County (245), State/Province (7739), Number of
Observers (2527), Submission ID (6033), Scientific Name (20215).

## descriptionContents

**What the observation description contains**

Formerly: P-040
Area: What birdsync writes
Goals:

The description contains: a line attributing the observation to birdsync, the
eBird observation details when present, the checklist URL, the protocol, the checklist
comments when present, one `Macaulay Library Asset:` line per uploaded asset, and one
`Macaulay Library Asset (upload failed permanently; delete this line from the description to
retry):` line per asset the service permanently refused ([permanentUploadFailureRecorded](#permanentuploadfailurerecorded)).

*The note spells out the remedy because the person who finds the line is reading an
observation, not the source, and the retry mechanism is otherwise undiscoverable.*

## observationsPublic

**Created observations are public**

Formerly: P-041
Area: What birdsync writes
Goals:

Observations birdsync creates are public, like any iNaturalist observation.

## createBeforeMedia

**The observation is created before media is attached**

Formerly: P-042
Area: Media
Goals:

The observation is created first and media attached afterward.

*Rationale: iNaturalist cannot attach media to an observation that does not yet exist.*

## mediaFormats

**Photos at 2400px, sounds as MP3**

Formerly: P-043
Area: Media
Subject: `media.photo.max_dimension_px`
Value: `2400`
Goals:

Photos are fetched at 2400px and sounds as MP3.

## photoThenSoundFallback

**Try the photo URL, then the sound URL**

Formerly: P-044
Area: Media
Goals:

A Macaulay Library asset ID does not say whether it is a photo or a sound, so
birdsync tries the photo URL and falls back to the sound URL on a 404.

## mediaFilename

**Media is uploaded as ML<asset ID> with a canonical extension**

Formerly: P-045
Area: Media
Subject: `media.filename.extension`
Value: `{image/jpeg: .jpg, image/png: .png, audio/mpeg: .mp3, audio/mpeg3: .mp3, audio/wav: .wav}`, falling back to `.jpg` or `.mp3` by endpoint
Goals:

Media is uploaded under the filename `ML<asset ID>` plus a canonical extension
for the response content type — `.jpg`, `.png`, `.mp3`, `.wav` — so any file in iNaturalist
can be traced back to its Macaulay Library asset.

*Rationale: the extension must not depend on the machine running the sync. Amended by
[decisions.canonicalFileExtensions](../../decisions.md#canonicalfileextensions); previously it was
whichever extension the system's mime database happened to list first.*

## descriptionUpdateKeepsMedia

**The description update keeps attached media**

Formerly: P-046
Area: Media
Goals:

After uploading, the observation description is updated with the asset URLs,
without detaching the media already attached.

## additiveMediaResync

**Media re-sync is additive**

Formerly: P-047
Area: Media
Goals:

Media re-sync is additive. Assets added to an eBird checklist after a sync
are uploaded on the next run.

## removedAssetsReported

**Assets removed from eBird are reported, not removed**

Formerly: P-048
Area: Media
Goals:

Assets removed from eBird are reported, not removed from iNaturalist.

## mediaCountMismatchReported

**A media count mismatch is reported, not corrected**

Formerly: P-049
Area: Media
Goals:

A mismatch between the asset count in the description and the media actually
attached is reported, not corrected.

## mediaFailuresTolerated

**A failed media transfer is logged, counted, and survived**

Formerly: P-050
Area: Media
Goals:

A failed media download or upload is logged and counted, and does not stop
the run.

*The description lists only assets that actually uploaded, so a failure is retried on the
next run rather than recorded as done. See
[decisions.descriptionListsOnlyUploaded](../../decisions.md#descriptionlistsonlyuploaded).*

## dryRunNoWrites

**--dryrun issues no writes**

Formerly: P-051
Area: Dry run
Goals:

`--dryrun` issues no write request of any kind. Reads are permitted.

## dryRunLogPrefix

**Skipped actions are logged with DRYRUN:**

Formerly: P-052
Area: Dry run
Subject: `log.dryrun.prefix`
Value: `"DRYRUN: "`
Goals:

Every action a dry run skips is logged with a `DRYRUN:` prefix.

*Rationale: the README tells users to grep for it.*

## dryRunMediaCount

**A dry run reports one unclassified media count**

Formerly: P-053
Area: Dry run
Goals:

A dry run reports media as a single "would upload N assets" count rather than
splitting photos from sounds.

*Rationale: it doesn't download the assets, and the asset ID doesn't reveal the type
([photoThenSoundFallback](#photothensoundfallback)).*

## endOfRunSummary

**Every run ends with a summary**

Formerly: P-054
Area: Reporting
Goals:

Every run ends with a summary: records processed, records skipped by each
rule, observations created and updated, and media uploaded.

## conditionalSkipCounters

**A skip counter prints only when its flag is in effect**

Formerly: P-055
Area: Reporting
Goals:

A skip counter is printed only when the flag that produces it is in effect.

## conditionalFailureCount

**The media failure count prints only when non-zero**

Formerly: P-056
Area: Reporting
Goals:

A media failure count is printed only when there were failures.

## malformedRecordAborts

**A malformed record aborts the run**

Formerly: P-057
Area: Failure behavior
Status: withdrawn by [decisions.badRecordsSkippedNotFatal](../../decisions.md#badrecordsskippednotfatal). Superseded by [badRecordsSkipped](#badrecordsskipped).
Goals:

A malformed date, time, latitude, or longitude in any record aborts the
whole run, identifying the CSV line.

## writeFailureAborts

**A failed create or update aborts the run**

Formerly: P-058
Area: Failure behavior
Goals:

A failed observation create or update aborts the run; a failed media transfer
does not.

*Rationale: a create or update failure means a broken token or a broken service, where
continuing would produce one identical error per remaining record. Reaffirmed under
[decisions.badRecordsSkippedNotFatal](../../decisions.md#badrecordsskippednotfatal).*

## noRollback

**There is no rollback**

Formerly: P-059
Area: Failure behavior
Goals:

There is no rollback. Observations created before an abort remain in the
user's account, and re-running is safe because of [idempotentRerun](#idempotentrerun).

## largeAccounts

**Works for accounts over 10,000 observations**

Formerly: P-065
Area: Failure behavior
Goals:

birdsync works for an account with more than 10,000 iNaturalist observations.

*Reported as [issue #5](https://github.com/Sajmani/birdsync/issues/5) in January 2026 and fixed
by [decisions.pageByIDAbove](../../decisions.md#pagebyidabove): the download
pages by observation id, which has no ceiling, rather than by page number, which the API caps
at 10,000 results.*

## unusableInputExplained

**Unusable input is named before parsing**

Formerly: P-066
Area: Failure behavior
Goals:

Before parsing, birdsync checks the input and names what is wrong, distinguishing
three cases: the path is a directory, the path is a zip archive, and the file has no
`Submission ID` column and so is not an eBird export.

*Rationale: [issue #1](https://github.com/Sajmani/birdsync/issues/1), where a user passed the
extracted download folder and got `Error reading CSV records from ...: Incorrect function.` —
an operating-system message that explains nothing.*

*The check is on the path and the header, never on the text of an operating-system error: that
symptom is what Windows returns for reading a directory, where Unix says "is a directory", so a
fix written against one would not work on the other ([tech.platformIndependent](../tech/tech.md#platformindependent)).*

*The header case is worse than an unhelpful message. Field lookup is `field[key]`, which yields
0 for an absent column, so a CSV without a `Submission ID` column did not fail — every record
silently took its submission ID from whatever the first column happened to be.*

## permanentUploadFailureRecorded

**A permanent upload failure is recorded and not retried**

Formerly: P-063
Area: Failure behavior
Subject: `media.upload.permanent_failure`
Value: `HTTP 4xx except 401, 408, 429`
Goals:

An upload that fails permanently — the service rejected the file itself, rather
than the attempt — is recorded in the description as failed and is not attempted again. A
failure that may be transient is not recorded, and is retried on the next run.

*Rationale: [decisions.permanentFailuresRecorded](../../decisions.md#permanentfailuresrecorded). Without
the distinction, birdsync either gives up on a network blip or re-downloads a 50 MB sound
file on every run forever. Removing the line by hand is the way to ask for a retry.*

## permanentFailuresReported

**Permanently failed assets are reported each run**

Formerly: P-064
Area: Failure behavior
Goals:

A permanently failed asset is reported on each run, so the user learns about it
without birdsync retrying it.

## documentCommunityObligations

**Documentation states the user's community obligations**

Formerly: P-067
Area: Failure behavior
Goals:

The documentation tells users what iNaturalist expects of an account that syncs
in bulk: review what was created, and respond to identifiers' comments.

*Rationale: [inatterms.dontBehaveLikeMachine](../../imported/inaturalist.org/inatterms/inatterms.md#dontbehavelikemachine) makes this the account holder's obligation, with suspension named
as the consequence, and birdsync is what puts them in a position to breach it. It also carries
more weight than politeness: [inatterms.noMachineGeneratedContent](../../imported/inaturalist.org/inatterms/inatterms.md#nomachinegeneratedcontent) permits a script to post observations only where a
human decided each one, so a user who reviews nothing is the one who moves birdsync from the
permitted side of that rule to the prohibited side. Recorded from
[decisions.notMachineGeneratedContent](../../decisions.md#notmachinegeneratedcontent).*

## documentSyncedObservationsIdentifiable

**Documentation says synced observations are identifiable**

Formerly: P-068
Area: Failure behavior
Goals:

The documentation explains that birdsync-created observations are identifiable —
by the attribution line and the two eBird observation fields — so that identifiers who would
rather not work on synced records can recognize them.

*Rationale: raised on the forum thread in [decisions.notMachineGeneratedContent](../../decisions.md#notmachinegeneratedcontent) as a mitigation for identifier burden. It is
already true ([syncKey](#synckey), [descriptionContents](#descriptioncontents)); it is simply not written down anywhere a reader would find it.*

## dryRunCounterLabels

**Dry-run observation counters are labeled hypothetical**

Formerly: P-060
Area: Reporting
Goals:

Under `--dryrun`, the observation counters are labeled as hypothetical:
`Would create N` and `Would update N`.

*Rationale: [decisions.dryRunCountersRelabeled](../../decisions.md#dryruncountersrelabeled). The count is
worth reporting; claiming the work happened is not.*

## noTaxonFilter

**The download is not filtered by taxon**

Formerly: P-061
Area: Downloading existing observations
Goals:

The download of existing iNaturalist observations is not filtered by taxon.

*Rationale: [decisions.noTaxonFilterOnDownload](../../decisions.md#notaxonfilterondownload). An
observation whose eBird name iNaturalist could not resolve has no iconic taxon, so a
taxon-filtered query never returns it, so re-running duplicates it. Idempotence ([idempotentRerun](#idempotentrerun))
outranks download volume. A consequence is that `--fuzzy` now matches against
observations of any taxon.*

## badRecordsSkipped

**Unparseable records are skipped, logged, and counted**

Formerly: P-062
Area: Failure behavior
Goals:

A record whose date, time, latitude, or longitude cannot be parsed is
skipped, logged with its CSV line, and counted. The run continues, and the count appears
in the summary when non-zero.

*Rationale: [decisions.badRecordsSkippedNotFatal](../../decisions.md#badrecordsskippednotfatal). This matches how
media failures are already handled ([mediaFailuresTolerated](#mediafailurestolerated)); one bad row in a 12,000-row export should not
end a sync halfway through.*
