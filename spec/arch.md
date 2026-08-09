# Architecture of birdsync

This document describes how the `birdsync` command-line tool is put together, for people
working on the code. For instructions on running it, see [README.md](../README.md).

## Overview

`birdsync` is a Go program that copies bird observations from an eBird data export into a
user's iNaturalist account, along with the photos and sounds attached to those observations
in the [Macaulay Library](https://www.macaulaylibrary.org/).

The sync is one-way (eBird → iNaturalist) and file-driven: the eBird API doesn't expose
personal checklists, so the input is the `MyEBirdData.csv` file the user downloads from eBird.

## Data flow

One run does the following:

1. **Authenticate.** Read the iNaturalist user ID and API token from `INAT_USER_ID` and
   `INAT_API_TOKEN`, prompting interactively if they're unset (`inat/vars.go`).
2. **Download existing observations.** Fetch the user's iNaturalist observations into memory,
   200 per page, restricted to the `Aves` iconic taxon and to the `--after`/`--before` date
   window when set (`inat.Client.DownloadObservations`, inat/inat.go:19).
3. **Build two indexes** over those observations (birdsync.go:136-167):
   - `previouslySynced`, keyed by `ebird.ObservationID` — the pair of eBird observation-field
     values that identifies a birdsync-created observation.
   - `fuzzyMatch`, keyed by observation date plus name, for every observation whose
     `ObservationID` is *not* valid. Each is indexed twice, under its common name and under its
     scientific name. Used only when `--fuzzy` is set.

   "Not valid" means *either* eBird field is missing (ebird/ebird.go:164), so an observation
   created by an old version of birdsync that set the checklist ID but not the scientific name
   lands in the fuzzy index instead of `previouslySynced`. That is the population `tools/repair`
   exists to fix.
4. **Read the CSV.** `ebird.Records` returns an `iter.Seq[ebird.Record]` over the export
   (ebird/ebird.go:80). Despite the iterator, this isn't streaming: the whole file is read with
   `csv.Reader.ReadAll` and the iterator walks the resulting slice.
5. **Filter and act** on each record, in this order: `--after`, `--before`, already-synced,
   `--fuzzy`, `--verifiable`. Records that survive all five become new iNaturalist observations.
   The already-synced branch isn't a pure skip: when eBird has assets the iNaturalist
   description doesn't list, it uploads them and updates the observation (birdsync.go:254-263).
6. **Create, then attach media.** The observation is created first, and each Macaulay Library
   asset is downloaded and uploaded to the now-existing observation afterward. The observation
   description is then updated with the asset URLs. Media cannot be attached to an observation
   that doesn't exist yet, which is why the order matters.

## The sync key

Birdsync recognizes its own observations using two iNaturalist
[observation fields](https://www.inaturalist.org/observation_fields):

- [eBird Checklist](https://www.inaturalist.org/observation_fields/6033) (`inat.EBirdField`,
  ID 6033) — the eBird submission ID, e.g. `S193523301`
- [eBird Scientific Name](https://www.inaturalist.org/observation_fields/20215)
  (`inat.EBirdScientificNameField`, ID 20215)

Together these form `ebird.ObservationID` (ebird/ebird.go:149), which makes the sync
idempotent: re-running birdsync over the same CSV skips everything it already uploaded.

The iNaturalist taxon is deliberately *not* used as the key, because the taxon may be changed
by the user or the community after upload, and because eBird names don't always map cleanly
onto iNaturalist taxa — eBird uses forms like `Aythya marila/affinis` ("slashes") and
`Melanitta sp.` ("spuhs") that have no exact iNaturalist equivalent.

## Package structure

The repository is a single Go module, `github.com/Sajmani/birdsync`, with one external
dependency (`github.com/google/uuid`).

### `main` (repository root)

- **`birdsync.go`** — flags, the `stats` counters, `main`, and `birdsync()`, which holds the
  whole sync loop. `birdsync()` takes its two clients as interfaces, so tests drive it without
  a network.
- **`glue.go`** — the seam that makes the above testable. Defines the `ebirdClient` and
  `inatClient` interfaces plus the real implementations that forward to the `ebird` and `inat`
  packages. Also defines `dateTimeFlag`, the `flag.Value` behind `--after` and `--before`.
- **`media.go`** — reconciling media between the two services. `mlAssetSet` is an ordered set
  of Macaulay Library asset IDs; `eBirdMLAssets` parses them from the CSV column and
  `iNatMLAssets` parses them back out of an iNaturalist observation's description text.
  `mediaChange` diffs the two and reports what changed.

### `ebird`

Everything that knows about eBird's data formats.

- `Record` mirrors a row of `MyEBirdData.csv`. Fields are read by *header name*, not position,
  and the CSV reader is set to `FieldsPerRecord = -1`, because eBird's export has a variable
  number of columns.
- `Record.Observed` parses the date and time. eBird writes dates as either `2006-01-02` or
  `1/2/2006` and the time may be absent, so this function handles four combinations. Anything
  comparing dates should go through it rather than reading `Record.Date` directly.
- `DownloadMLAsset` fetches an asset from the Macaulay Library CDN. An asset ID doesn't say
  whether it's a photo or a sound, so this tries the photo URL (`/asset/<id>/2400`) and falls
  back to the sound URL (`/asset/<id>/mp3`) on a 404. It returns a temp-file path, an
  `isPhoto` flag, and derives the file extension from the response `Content-Type`.

### `inat`

A hand-written client for the [iNaturalist API v2](https://api.inaturalist.org/v2/docs/).

- `client.go` — `Client` and its `roundTrip` helper, which sets the `Authorization` and
  `User-Agent` headers and turns a 401 into a "refresh your token" message.
  `CreateObservation`, `UpdateObservation`, `DeleteObservation`, and `UploadMedia`.
  `UpdateObservation` always sets `ignore_photos` so that updating a description can't clobber
  attached media.
- `inat.go` — `DownloadObservations`, which handles pagination and the `fields` parameter that
  selects which parts of each observation the API returns.
- `types.go` — the API's JSON shapes, and the observation-field ID constants.
- `vars.go` — `GetUserID` and `GetAPIToken`, including the interactive prompts.

Note the asymmetry between `Observation` (what birdsync sends) and `Result` (what the API
returns); they are not the same shape.

### `tools/`

Seven standalone `main` packages for maintaining observations after the fact. They are *not*
imported by birdsync, and `go install github.com/Sajmani/birdsync@latest` doesn't install them
(though `go install ./...` from a clone does). Run them with `go run ./tools/<name>`.

`dedupe` and `repair` exist for a historical reason worth knowing. Before the two eBird
observation fields were added (`3d036ba` and `7478a09`, July 2025), birdsync could not
reliably recognize its own observations, so re-runs created duplicates. `dedupe` cleans up
that era's damage and `repair` backfills the missing field on what survived. **Neither is
evidence of a current defect** — treating `dedupe`'s existence as a symptom of a live bug
is a mistake already made once and retracted, in [decisions.md](decisions.md) CR-003.

`taxonfilter` is the exception to the "no tests" rule below: because T-011 forbids running
anything here, its comparison logic and query construction are covered by
`tools/taxonfilter/taxonfilter_test.go` instead.

Each builds its own `inat.Client`, except `taxonfilter`, which issues its own GET requests so
it can vary the `iconic_taxa[]` parameter. `dedupe`, `purge`, and `position` find birdsync-created
observations using the same `ebird.ObservationID` key as the sync loop; `repair` reads the two
observation fields directly; `dump` and `poke` don't filter at all.

| Tool | Purpose | Guard |
| --- | --- | --- |
| `dump` | Print all downloaded observations as JSON | read-only |
| `poke` | Create one test observation, or upload one asset | none |
| `dedupe` | Delete duplicate observations sharing a sync key | `debug` const |
| `purge` | Delete observations with no photos and no sounds | `debug` const |
| `position` | Reset positional accuracy to `ebird.PositionalAccuracy` | `debug` const |
| `repair` | Backfill the eBird scientific name field | none |
| `taxonfilter` | Compare a filtered and unfiltered download; diagnose CR-003 | read-only |

`dedupe`, `purge`, and `position` are gated by a `debug` constant at the top of the file:
`true` logs what would happen, `false` performs it. There is no `--dryrun` flag. As checked in,
`purge` deletes for real while `dedupe` and `position` are log-only. Check the constant before
running any of them.

## Testing

`go test ./...` runs everything and touches neither the network nor a real iNaturalist account.

| File | Covers |
| --- | --- |
| `birdsync_test.go` | The sync loop, via `mockEBirdClient` and `mockINatClient` |
| `media_test.go` | `mediaChange`; the `mlAssetSet` helpers only indirectly |
| `ebird/ebird_test.go` | CSV parsing (temp file), `Record.Observed` date formats, `ObservationID.Valid`, and `downloadMLAsset` against an `httptest` server |
| `inat/inat_test.go` | `DownloadObservations`, against an `httptest` server |
| `inat/client_test.go` | `CreateObservation`, `UpdateObservation`, `DeleteObservation` |

`inat.Client.UploadMedia` has no test — the multipart request it builds (inat/client.go:125) is
only exercised through the mock in `birdsync_test.go`. It's the largest untested surface here.

Two things make this work: the client interfaces in `glue.go`, and base-URL parameters that let
a test server stand in for the real service — `inat.NewClient` takes one, and the exported
`ebird.DownloadMLAsset` delegates to an unexported `downloadMLAsset(baseURL, id)` that the
package's own tests call directly.

The flags are package-level variables, so tests mutate global state. `resetFlags()` in
`birdsync_test.go` restores all of them except `debug`, which each test saves and restores
itself. Call it at the top of any new test. Note that `dateTimeFlag.Set("")` returns an error
and leaves the previous value in place, so the date flags must be zeroed by assignment rather
than through `Set`.
