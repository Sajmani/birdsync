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
   200 per page, restricted only to the `--after`/`--before` date window when set
   (`inat.Client.DownloadObservations`). Deliberately not filtered by taxon:
   see CR-003 in [decisions.md](decisions.md).
3. **Build two indexes** over those observations, at the top of `birdsync()`:
   - `previouslySynced`, keyed by `ebird.ObservationID` — the pair of eBird observation-field
     values that identifies a birdsync-created observation.
   - `fuzzyMatch`, keyed by observation date plus name, for every observation whose
     `ObservationID` is *not* valid. Each is indexed twice, under its common name and under its
     scientific name. Used only when `--fuzzy` is set.

   "Not valid" means *either* eBird field is missing (`ebird.ObservationID.Valid`), so an observation
   created by an old version of birdsync that set the checklist ID but not the scientific name
   lands in the fuzzy index instead of `previouslySynced`. A `repair` tool used to backfill
   that population; it was deleted once the maintainer's account was clean.
4. **Read the CSV.** `ebird.Records` returns an `iter.Seq[ebird.Record]` over the export.
   Despite the iterator, this isn't streaming: the whole file is read with
   `csv.Reader.ReadAll` and the iterator walks the resulting slice.
5. **Filter and act** on each record, in this order: `--after`, `--before`, already-synced,
   `--fuzzy`, `--verifiable`. Records that survive all five become new iNaturalist observations.
   A record whose date or coordinates won't parse is skipped and counted rather than ending
   the run: the date is checked before the date filters, the coordinates just before the
   observation is built.
   The already-synced branch isn't a pure skip: when eBird has assets the iNaturalist
   description doesn't list, it uploads them and updates the observation (the `addMedia`
   closure).
6. **Create, then attach media.** The observation is created first, and each Macaulay Library
   asset is downloaded and uploaded to the now-existing observation afterward. The observation
   description is then updated with the asset URLs. Media cannot be attached to an observation
   that doesn't exist yet, which is why the order matters. Each downloaded file is removed
   straight after its upload attempt, so at most one asset sits on disk at a time.
7. **Report.** `stats.summary()` builds the end-of-run lines and `main` logs them. It is a
   separate function so it can be tested: the labels change under `--dryrun`, and a summary
   that contradicts what happened is invisible until a user acts on it.

## The sync key

Birdsync recognizes its own observations using two iNaturalist
[observation fields](https://www.inaturalist.org/observation_fields):

- [eBird Checklist](https://www.inaturalist.org/observation_fields/6033) (`inat.EBirdField`,
  ID 6033) — the eBird submission ID, e.g. `S193523301`
- [eBird Scientific Name](https://www.inaturalist.org/observation_fields/20215)
  (`inat.EBirdScientificNameField`, ID 20215)

Together these form `ebird.ObservationID`, which makes the sync
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

One standalone `main` package, `dump`, which downloads the user's observations and prints
them as JSON. It is *not* imported by birdsync, and `go install
github.com/Sajmani/birdsync@latest` doesn't install it (though `go install ./...` from a
clone does). Run it with `go run ./tools/dump`.

`tools/` is read-only by construction, enforced by `TestToolsAreReadOnly` (T-032). It used
to hold six more programs that created, updated, and deleted real observations, guarded
only by a `debug` constant at the top of each file — `purge` shipped with its guard off,
deleting for real. Four of them (`dedupe`, `repair`, `position`, `purge`) were one-time
cleanups for defects that have since been fixed, `poke` was a token smoke test, and
`taxonfilter` answered CR-003. All were deleted; the history has them if one is needed
again.

Deleting them turned a rule people had to remember into a property a check can enforce,
which is the trade this repository prefers wherever it is available.

## Testing

`go test ./...` runs everything and touches neither the network nor a real iNaturalist account.

| File | Covers |
| --- | --- |
| `birdsync_test.go` | The sync loop and `stats.summary()`, via `mockEBirdClient` and `mockINatClient` |
| `guard_test.go` | Static analysis over the repository itself: no live hostnames in tests, no writes under `tools/`, no `log.Fatal` in library packages |
| `media_test.go` | `mediaChange`; the `mlAssetSet` helpers only indirectly |
| `ebird/ebird_test.go` | CSV parsing (temp file), `Record.Observed` date formats, `ObservationID.Valid`, and `downloadMLAsset` against an `httptest` server |
| `inat/inat_test.go` | `DownloadObservations`: pagination, query parameters, and the error path |
| `inat/client_test.go` | `CreateObservation`, `UpdateObservation` (including `ignore_photos`), `DeleteObservation` |

The fakes in `birdsync_test.go` record every mutating call they receive. That is what makes
`--dryrun` checkable directly — asserting that nothing was written, rather than inferring it
from counters, which is how the bug in CR-001 survived.

`inat.Client.UploadMedia` has no test — the multipart request it builds is only exercised
through the mock in `birdsync_test.go`. It's the largest untested surface here, and the only
place binary data is written to a user's account.

`spec/acceptance.md` maps each requirement to the check that verifies it, and lists the ones
with no check.

Two things make this work: the client interfaces in `glue.go`, and base-URL parameters that let
a test server stand in for the real service — `inat.NewClient` takes one, and the exported
`ebird.DownloadMLAsset` delegates to an unexported `downloadMLAsset(baseURL, id)` that the
package's own tests call directly.

The flags are package-level variables, so tests mutate global state. `resetFlags()` in
`birdsync_test.go` restores all of them except `debug`, which each test saves and restores
itself. Call it at the top of any new test. Note that `dateTimeFlag.Set("")` returns an error
and leaves the previous value in place, so the date flags must be zeroed by assignment rather
than through `Set`.
