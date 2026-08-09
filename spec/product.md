# Product requirements — birdsync

What birdsync must do, from a user's point of view. Written per
[process.md](process.md); see it for what the IDs mean and how requirements are amended.

**Status: approved at Gate 1 on 2026-08-09.** This is a retrofit — the requirements below were
reverse-engineered from the implementation, [README.md](../README.md), and
[arch.md](arch.md), and describe *current* behavior. Where current behavior looks
unintended, it is recorded under [Open questions](#open-questions) or as a conflict in
[decisions.md](decisions.md) rather than being quietly restated as if intended.

## Purpose

**P-001** — birdsync copies bird observations from a user's eBird data export into that
user's iNaturalist account, including the photos and sounds attached to them in the
Macaulay Library.

**P-002** — The sync is one-way: eBird → iNaturalist.

**P-003** — The input is a `MyEBirdData.csv` file the user downloads from eBird, not a
live eBird API.
*Rationale: the eBird API does not expose personal checklists.*

## Non-goals

**P-004** — birdsync does not sync iNaturalist observations back to eBird.

**P-005** — birdsync never modifies or deletes an iNaturalist observation it did not
create, except to attach media to an observation carrying its own sync key.

**P-006** — birdsync does not reconcile taxonomy between the two services. An eBird
scientific name that iNaturalist cannot resolve produces an observation with an unknown
taxon, which the user fixes by hand.

**P-007** — The programs in `tools/` are not part of the product. They are maintenance
utilities, are not installed by `go install` of the root package, and carry no
compatibility promise. They are also read-only (T-032): no tool may alter a user's
account.

## Invocation

**P-008** — birdsync takes exactly one positional argument: the path to the eBird CSV
export. Any other number of arguments prints usage and exits non-zero.

**P-009** — Flags must appear before the positional argument.
*Rationale: consequence of the standard Go flag package; documented rather than fixed.*

**P-010** — A boolean flag is disabled with `--flag=false`; `--flag false` is a usage
error.
*Rationale: same as P-009.*

**P-011** — If the CSV file cannot be opened, birdsync exits with an error before
contacting iNaturalist or prompting for credentials.

**P-012** — If `--after` is later than `--before`, birdsync exits with an error rather
than running a sync that cannot match anything.

**P-013** — `--debug` logs the reason each eBird record was skipped, and other detail
not shown in a normal run.

## Authentication

**P-014** — The iNaturalist user ID is read from `INAT_USER_ID`, and the API token from
`INAT_API_TOKEN`.

**P-015** — When either is unset, birdsync prompts for it interactively and explains
where to find it.

**P-016** — At the prompt the token is accepted in its `{"api_token":"TOKEN"}` framing;
from the environment it is the bare token.
*Rationale: matches what the iNaturalist page gives the user, versus what is convenient
to export in a shell.*

**P-017** — A 401 from the API reports that the token needs refreshing, and where to get
a new one, rather than a bare HTTP status.
*Rationale: tokens expire every 24 hours, so this is the most likely failure a user
hits.*

**P-018** — Credentials are never written to the log.

## Identity and idempotence

**P-019** — birdsync identifies its own observations by a sync key: the pair of
iNaturalist observation fields *eBird Checklist* (6033) and *eBird Scientific Name*
(20215).
Subject: `sync.key.fields` · Value: `[6033, 20215]`

**P-020** — Re-running birdsync over the same CSV creates no duplicate observations.

**P-021** — The iNaturalist taxon is not part of the sync key.
*Rationale: the taxon may be changed by the user or the community after upload, and
eBird's "slash" (`Aythya marila/affinis`) and "spuh" (`Melanitta sp.`) names have no
exact iNaturalist equivalent.*

**P-022** — An observation missing either sync-key field is not recognized as
birdsync-created.
*Rationale: observations created by versions that set only the checklist field fall into
this population. A `repair` tool existed to backfill them and was deleted once the
maintainer's account was clean; it is recoverable from the history.*

## Downloading existing observations

**P-023** — Before syncing, birdsync downloads the user's existing iNaturalist
observations into memory.

**P-024** — ~~The download is restricted to the `Aves` iconic taxon.~~
Status: **Withdrawn** by [CR-003](decisions.md#cr-003--aves-only-download-can-defeat-duplicate-detection) — it defeated P-020. Superseded by P-061.

**P-025** — When `--after` or `--before` is set, the download is restricted to that date
window, so duplicate detection and fuzzy matching only consider observations inside it.

## Filtering and skip order

**P-026** — Each eBird record is tested in this order: `--after`, `--before`,
already-synced, `--fuzzy`, `--verifiable`. A record surviving all five becomes a new
iNaturalist observation.
*Rationale: the order is user-visible, because it determines which counter a skipped
record lands in.*

**P-027** — `--after` skips records observed strictly before the given time.

**P-028** — `--before` skips records observed strictly after the given time.

**P-029** — Both accept `2006-01-02 15:04:05` or `2006-01-02`.

**P-030** — `--verifiable` skips records with no Macaulay Library catalog numbers, and
defaults to **on**.
Subject: `cli.verifiable.default` · Value: `true`
*Rationale: iNaturalist will not treat a medialess observation as verifiable, so the default
avoids filling a user's account with "Casual" records.*
*It is also the main thing birdsync does to limit the burden it puts on other people. A
medialess observation gives an identifier nothing to identify, so it is pure cost to the
community — the concern a forum moderator raised about volume
([CR-012](decisions.md#cr-012--is-birdsync-machine-generated-content)). Don't flip this
default for convenience.*

**P-031** — `--fuzzy` skips a record when an existing observation not carrying a sync key
has the same observation date and the same name, matched against either the common name
or the scientific name.

**P-032** — Fuzzy matching never uses an empty name.
*Rationale: an empty name matches every unnamed record on that date and silently drops
legitimate observations.*

**P-033** — Fuzzy matching compares parsed observation dates, not raw CSV date strings.
*Rationale: eBird writes dates in two formats; keying on the raw field silently disabled
`--fuzzy` for anyone whose export used the second one.*

**P-034** — `--fuzzy` is off by default, and its documentation states that it may
suppress a legitimate record when the same bird was seen twice in one day.

## What birdsync writes

**P-035** — A created observation is marked wild, not captive.
*Rationale: eBird checklists record wild birds.*

**P-036** — Latitude and longitude are the checklist's, the location is marked
not-exact, and positional accuracy is `--positional_accuracy_meters`, default 1000.
Subject: `observation.positional_accuracy.default_m` · Value: `1000`
*Rationale: the checklist location approximates a hotspot, not the bird.*

**P-037** — The species guess is the eBird scientific name.

**P-038** — The observation date and time come from the eBird `Date` and `Time` columns.

**P-039** — These eBird columns are copied into iNaturalist observation fields: Count
(1), Common Name (256), Location (157), County (245), State/Province (7739), Number of
Observers (2527), Submission ID (6033), Scientific Name (20215).

**P-040** — The description contains: a line attributing the observation to birdsync, the
eBird observation details when present, the checklist URL, the protocol, the checklist
comments when present, one `Macaulay Library Asset:` line per uploaded asset, and one
`Macaulay Library Asset (upload failed permanently; delete this line from the description to
retry):` line per asset the service permanently refused (P-063).
*The note spells out the remedy because the person who finds the line is reading an
observation, not the source, and the retry mechanism is otherwise undiscoverable.*

**P-041** — Observations birdsync creates are public, like any iNaturalist observation.

## Media

**P-042** — The observation is created first and media attached afterward.
*Rationale: iNaturalist cannot attach media to an observation that does not yet exist.*

**P-043** — Photos are fetched at 2400px and sounds as MP3.
Subject: `media.photo.max_dimension_px` · Value: `2400`

**P-044** — A Macaulay Library asset ID does not say whether it is a photo or a sound, so
birdsync tries the photo URL and falls back to the sound URL on a 404.

**P-045** — Media is uploaded under the filename `ML<asset ID>` plus a canonical extension
for the response content type — `.jpg`, `.png`, `.mp3`, `.wav` — so any file in iNaturalist
can be traced back to its Macaulay Library asset.
Subject: `media.filename.extension` · Value: `{image/jpeg: .jpg, image/png: .png, audio/mpeg: .mp3, audio/mpeg3: .mp3, audio/wav: .wav}`, falling back to `.jpg` or `.mp3` by endpoint
*Rationale: the extension must not depend on the machine running the sync. Amended by
[CR-006](decisions.md#cr-006--uploaded-file-extensions-varied-by-machine); previously it was
whichever extension the system's mime database happened to list first.*

**P-046** — After uploading, the observation description is updated with the asset URLs,
without detaching the media already attached.

**P-047** — Media re-sync is additive. Assets added to an eBird checklist after a sync
are uploaded on the next run.

**P-048** — Assets removed from eBird are reported, not removed from iNaturalist.

**P-049** — A mismatch between the asset count in the description and the media actually
attached is reported, not corrected.

**P-050** — A failed media download or upload is logged and counted, and does not stop
the run.
*The description lists only assets that actually uploaded, so a failure is retried on the
next run rather than recorded as done. See
[CR-007](decisions.md#cr-007--a-failed-media-upload-leaves-its-url-in-the-description).*

## Dry run

**P-051** — `--dryrun` issues no write request of any kind. Reads are permitted.

**P-052** — Every action a dry run skips is logged with a `DRYRUN:` prefix.
Subject: `log.dryrun.prefix` · Value: `"DRYRUN: "`
*Rationale: the README tells users to grep for it.*

**P-053** — A dry run reports media as a single "would upload N assets" count rather than
splitting photos from sounds.
*Rationale: it doesn't download the assets, and the asset ID doesn't reveal the type
(P-044).*

## Reporting

**P-054** — Every run ends with a summary: records processed, records skipped by each
rule, observations created and updated, and media uploaded.

**P-055** — A skip counter is printed only when the flag that produces it is in effect.

**P-056** — A media failure count is printed only when there were failures.

## Failure behavior

**P-057** — ~~A malformed date, time, latitude, or longitude in any record aborts the
whole run, identifying the CSV line.~~
Status: **Withdrawn** by [CR-005](decisions.md#cr-005--record-level-parse-failures-abort-the-run). Superseded by P-062.

**P-058** — A failed observation create or update aborts the run; a failed media transfer
does not.
*Rationale: a create or update failure means a broken token or a broken service, where
continuing would produce one identical error per remaining record. Reaffirmed under
CR-005.*

**P-059** — There is no rollback. Observations created before an abort remain in the
user's account, and re-running is safe because of P-020.

**P-065** — birdsync works for an account with more than 10,000 iNaturalist observations.
*Reported as [issue #5](https://github.com/Sajmani/birdsync/issues/5) in January 2026 and fixed
by [CR-011](decisions.md#cr-011--the-download-cannot-page-past-10000-results): the download
pages by observation id, which has no ceiling, rather than by page number, which the API caps
at 10,000 results.*

**P-066** — Before parsing, birdsync checks the input and names what is wrong, distinguishing
three cases: the path is a directory, the path is a zip archive, and the file has no
`Submission ID` column and so is not an eBird export.
*Rationale: [issue #1](https://github.com/Sajmani/birdsync/issues/1), where a user passed the
extracted download folder and got `Error reading CSV records from ...: Incorrect function.` —
an operating-system message that explains nothing.*
*The check is on the path and the header, never on the text of an operating-system error: that
symptom is what Windows returns for reading a directory, where Unix says "is a directory", so a
fix written against one would not work on the other (T-004).*
*The header case is worse than an unhelpful message. Field lookup is `field[key]`, which yields
0 for an absent column, so a CSV without a `Submission ID` column did not fail — every record
silently took its submission ID from whatever the first column happened to be.*

**P-063** — An upload that fails permanently — the service rejected the file itself, rather
than the attempt — is recorded in the description as failed and is not attempted again. A
failure that may be transient is not recorded, and is retried on the next run.
Subject: `media.upload.permanent_failure` · Value: `HTTP 4xx except 401, 408, 429`
*Rationale: [CR-008](decisions.md#cr-008--a-permanently-rejected-asset-was-retried-forever). Without
the distinction, birdsync either gives up on a network blip or re-downloads a 50 MB sound
file on every run forever. Removing the line by hand is the way to ask for a retry.*

**P-064** — A permanently failed asset is reported on each run, so the user learns about it
without birdsync retrying it.

**P-067** — The documentation tells users what iNaturalist expects of an account that syncs
in bulk: review what was created, and respond to identifiers' comments.
*Rationale: `inat-terms/R5` makes this the account holder's obligation, with suspension named
as the consequence, and birdsync is what puts them in a position to breach it. Recorded from
[CR-012](decisions.md#cr-012--is-birdsync-machine-generated-content).*

**P-068** — The documentation explains that birdsync-created observations are identifiable —
by the attribution line and the two eBird observation fields — so that identifiers who would
rather not work on synced records can recognise them.
*Rationale: raised on the forum thread in CR-012 as a mitigation for identifier burden. It is
already true (P-019, P-040); it is simply not written down anywhere a reader would find it.*

## Amendments from Gate 1

**P-060** — Under `--dryrun`, the observation counters are labelled as hypothetical:
`Would create N` and `Would update N`.
*Rationale: [CR-001](decisions.md#cr-001--a-dry-run-reports-observations-it-did-not-create). The count is
worth reporting; claiming the work happened is not.*

**P-061** — The download of existing iNaturalist observations is not filtered by taxon.
*Rationale: [CR-003](decisions.md#cr-003--aves-only-download-can-defeat-duplicate-detection). An
observation whose eBird name iNaturalist could not resolve has no iconic taxon, so a
taxon-filtered query never returns it, so re-running duplicates it. Idempotence (P-020)
outranks download volume. A consequence is that `--fuzzy` now matches against
observations of any taxon.*

**P-062** — A record whose date, time, latitude, or longitude cannot be parsed is
skipped, logged with its CSV line, and counted. The run continues, and the count appears
in the summary when non-zero.
*Rationale: [CR-005](decisions.md#cr-005--record-level-parse-failures-abort-the-run). This matches how
media failures are already handled (P-050); one bad row in a 12,000-row export should not
end a sync halfway through.*

## Open questions

Resolved at Gate 1: partial-failure policy (CR-005), dry-run counter labels (CR-001),
unknown-taxon re-runs (CR-003), and the `--fuzzy` documentation (CR-004). Still open:

1. ~~**Media ownership.**~~ Answered by vendoring `ml-terms` and closing
   [CR-010](decisions.md#cr-010--birdsync-cannot-tell-whose-media-it-is-downloading): a
   contributor keeps copyright in their own media and may download it freely, and birdsync
   relies on the documented, evidence-supported assumption that an export lists only the
   user's own assets.
2. **Sound files over 50 MB** are rejected by iNaturalist. Currently this surfaces as a
   generic upload failure. Should it be detected and reported specifically?
