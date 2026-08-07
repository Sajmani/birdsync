# birdsync
Birdsync syncs eBird observations, photos, and sounds to iNaturalist.

# Requirements

You must download your data from eBird using
https://ebird.org/downloadMyData.
Save the zip file and unzip it to get the `MyEBirdData.csv` file.

To run birdsync, you'll need the Go language toolchain.
Download it from http://go.dev.

Birdsync is a command line program.
On Macs you can run commands using Terminal.

# Install and run birdsync

Install or update `birdsync` to the latest version using
```
go install github.com/Sajmani/birdsync@latest
```
By default, Go installs binaries in the directory `$HOME/go/bin`.
[Customize this location](https://pkg.go.dev/cmd/go#hdr-Compile_and_install_packages_and_dependencies) by setting the `GOBIN` environment variable.

Run `birdsync` from the command line, specifying the path to your `MyEBirdData.csv` file:
```
$HOME/go/bin/birdsync MyEBirdData.csv
```
Consider running a "dry run" to test what birdsync would do without actually touching your iNaturalist observations:
```
$HOME/go/bin/birdsync --dryrun MyEBirdData.csv
```
Look for log lines starting with "DRYRUN" to see what observations birdsync will create and which media files it will copy. The dry run prints out the full observation data structure and so can be quite verbose.

Birdsync will prompt you to enter your iNaturalist user name and [API token](https://www.inaturalist.org/users/api_token), which allow the tool to read and write your personal iNaturalist observations.
Copy the full string from the web page, including both curly braces: `{"api_token":"TOKEN"}`

To skip these interactive steps, you can provide your iNaturalist user name and API token as environment variables, but remember that you need to refresh your [API token](https://www.inaturalist.org/users/api_token) every 24 hours:
```
export INAT_USER_ID=(your iNaturalist user name)
export INAT_API_TOKEN=(just the TOKEN part of {"api_token":"TOKEN"})
```
Birdsync provides command-line flags to customize its behavior:
*  `-after 2006-01-02`
        Sync only observations observed after the provided date and time (formatted as "2006-01-02 15:04:05"). The time can be omitted (2006-01-02).
* `-before 2006-01-02`
        Sync only observations observed before the provided date and time (formatted as "2006-01-02 15:04:05"). The time can be omitted (2006-01-02).
* `-dryrun`
        Don't actually sync any observations, just log what birdsync would do
* `-verifiable` (**default `true`**)
        Sync only observations that include Macaulay Catalog Numbers (photos or sound), as iNaturalist requres media to consider an observation verifiable.
        This is on by default, so by default birdsync will _not_ sync observations
        that have no photos or sounds. To sync those observations too, pass `--verifiable=false`.
* `-fuzzy`
        Don't create a birdsync observation if a non-birdsync observation already exists for the same bird on the same date. This fuzzy matching is useful when you've entered the same observation manually into both eBird and iNaturalist, but it may skip legitimate uploads if you saw the same bird twice on the same day.
* `-positional_accuracy_meters` (default `1000`)
        Positional accuracy in meters of the iNaturalist observations created by birdsync.
        Since the latitude and longitude of birdsync observations is set to the checklist location,
        this may be distant from the actual location where individual birds were observed.
        Birdsync uses default positional accuracy of 1000 meters; use this flag to adjust it.
* `-debug`
        Log verbosely. Useful for seeing exactly why each eBird observation was skipped.

Boolean flags must be turned off using `=`: `--verifiable=false` works, but `--verifiable false`
fails with a usage error, because `false` is read as a positional argument rather than as the
flag's value.

On the command line, flags must be listed _before_ your `MyEBirdData.csv` file:
```
$HOME/go/bin/birdsync --fuzzy --after 2025-07-01 MyEBirdData.csv
```

Birdsync exits with an error if `--after` is later than `--before`, since that combination
can't match any records.

## What birdsync prints when it finishes

Birdsync ends each run with a summary of what it did, for example:
```
Finished processing 12043 eBird observations
Skipped 11890 previously uploaded by birdsync
Skipped 96 unverifiable eBird observations
Created 57 new iNaturalist observations
Updated 57 iNaturalist observations
Uploaded 61 photos to iNaturalist
Uploaded 4 sounds to iNaturalist
```
The skip counts for `--fuzzy`, `--after`, `--before`, and `--verifiable` are only printed when
those flags are in effect. A final "Failed to upload N media assets" line appears if any
media downloads or uploads failed; those failures are logged but don't stop the run.

Under `--dryrun`, the last two lines are replaced by a single
"Would upload N media assets to iNaturalist". A dry run doesn't download anything, and a
Macaulay Library asset ID doesn't say whether it refers to a photo or a sound, so a dry run
can't split the count into photos and sounds.

## Checking the results

Once birdsync has finished running, you should check the observations it created:
- If iNaturalist doesn't recognize the scientific name provided by eBird, the observation species name will say "Unknown". Fix this by editing the observation in iNaturalist.
- If the iNaturalist observation has no photos or sounds, either because none were in eBird or because birdsync failed to copy them, then the observation will be marked "Casual". Fix this by uploading media for these observations or deleting them. Note that birdsync skips observations without media by default (`--verifiable` defaults to true), so this mostly happens when a media upload failed. iNaturalist rejects sound files larger than 50 MB; in these cases you will need to add a smaller file to the observation.

# How birdsync works

Given (`iNaturalist user name`, `eBird CSV file`):
- Download that user's existing iNaturalist bird observations into memory.
  Only observations in the `Aves` iconic taxon are downloaded, and if `--after` or `--before`
  are set, only observations in that date range are downloaded.
- Index these iNaturalist observations by ([eBird submission ID](https://www.inaturalist.org/observation_fields/6033), [eBird scientific name](https://www.inaturalist.org/observation_fields/20215))
- Index any non-birdsync observations by date for fuzzy matching, under both their
  common name and their scientific name
- For each eBird observation in `eBird CSV file`, in this order:
  - If `--after` is set, skip any eBird observations before that date
  - If `--before` is set, skip any eBird observations after that date
  - Skip any eBird observations that have already been uploaded
    - If photos or sounds have been added to eBird since the last sync, upload them to iNaturalist
      and append their URLs to the observation description
    - If photos or sounds have been _removed_ from eBird since the last sync, log the difference
      but leave the iNaturalist observation alone
  - If `--fuzzy` is set, skip any eBird observations for the same bird and day as a non-birdsync observation
  - If `--verifiable` is set (the default), skip any eBird observations lacking photos or sounds
  - Create a new iNaturalist observation from the eBird observation
  - For each [Macaulay Library](https://www.macaulaylibrary.org/) catalog ID for this eBird observation:
    - Download the photo or sound from the Macaulay Library.
      Photos are fetched at 2400px; sounds are fetched as MP3.
      Asset IDs don't say whether they're a photo or a sound, so birdsync tries the photo
      URL first and falls back to the sound URL.
    - Upload the photo or sound to iNaturalist, associated with the new observation.
      Media is uploaded under the filename `ML<asset ID>`, so you can trace any file in
      iNaturalist back to its Macaulay Library asset.

## What birdsync writes to iNaturalist

Each observation birdsync creates is marked as wild (not captive), with the checklist's
latitude and longitude, a non-exact location, and the positional accuracy set by
`--positional_accuracy_meters`. The species guess is the eBird scientific name, and the
observation date/time is taken from the eBird `Date` and `Time` columns.

Birdsync copies these eBird columns into iNaturalist [observation fields](https://www.inaturalist.org/observation_fields):

| eBird column | iNaturalist observation field |
| --- | --- |
| Count | [Count](https://www.inaturalist.org/observation_fields/1) |
| Common Name | [Common Name](https://www.inaturalist.org/observation_fields/256) |
| Location | [Location](https://www.inaturalist.org/observation_fields/157) |
| County | [County](https://www.inaturalist.org/observation_fields/245) |
| State/Province | [State or Province](https://www.inaturalist.org/observation_fields/7739) |
| Number of Observers | [Number of Observers](https://www.inaturalist.org/observation_fields/2527) |
| Submission ID | [eBird Checklist](https://www.inaturalist.org/observation_fields/6033) |
| Scientific Name | [eBird Scientific Name](https://www.inaturalist.org/observation_fields/20215) |

The last two fields are what birdsync uses to recognize its own observations on later runs,
so don't remove them if you want re-syncing to work. Birdsync deliberately doesn't rely on the
iNaturalist taxon for this, because the taxon may be corrected by you or the community
after upload.

The observation description contains a note that birdsync created it, the eBird observation
details and checklist comments (when present), the checklist URL, the eBird protocol, and one
`Macaulay Library Asset:` line per uploaded photo or sound.

# Limitations

Birdsync only works in the eBird → iNaturalist direction because (as far as I can tell) the [eBird API](https://support.ebird.org/en/support/solutions/articles/48000838205-download-ebird-data#API) doesn't support reading or writing personal checklists, only reading "limited, recent and summary outputs of eBird data".

Birdsync cannot detect whether iNaturalist observations that you've created manually are duplicates of those in your eBird checklists unless you mark your existing iNaturalist observations with the [eBird submission ID](https://www.inaturalist.org/observation_fields/6033) and [eBird scientific name](https://www.inaturalist.org/observation_fields/20215) observation fields. The `--fuzzy` matching feature provides a convenient way to avoid creating duplicates, but it may also suppress creating legitimate observations if you happened to see the same bird twice on the same day and entered it once into each tool.

`--fuzzy` only compares against bird (`Aves`) observations, since those are the only ones
birdsync downloads. It also ignores iNaturalist observations that have no taxon name, so an
unidentified observation won't suppress your eBird records for that day.

Media re-syncing is one-way and additive. If you add photos or sounds to an eBird checklist
after a sync, the next run uploads them. If you remove media from eBird, or if the assets listed
in the iNaturalist description don't match the media actually attached to the observation,
birdsync reports the discrepancy but does not fix it.

Restricting a run with `--after` or `--before` also restricts which existing iNaturalist
observations get downloaded. Duplicate detection and fuzzy matching therefore only consider
observations inside that date window.

# Tools

The `tools` directory contains small standalone programs for maintaining observations that
birdsync has already created. They are separate `main` packages, so
`go install github.com/Sajmani/birdsync@latest` does not install them; run them from a clone
with `go run ./tools/<name>`. They read `INAT_USER_ID` and `INAT_API_TOKEN` (or prompt) the
same way birdsync does.

- `dump` — download your observations and print them as JSON. Read-only; useful for debugging.
- `poke` — create a single test observation, or upload one Macaulay Library asset to an
  observation you name. Useful for checking that your API token works.
- `dedupe` — find observations that share an (eBird checklist, eBird scientific name) pair and
  delete all but one, preferring the one with the most identifications, then the oldest.
- `purge` — delete birdsync observations that have neither photos nor sounds.
- `position` — reset the positional accuracy of birdsync observations to 1000 meters.
- `repair` — backfill the eBird scientific name field on observations created before birdsync
  started setting it. Takes `MyEBirdData.csv` as its only argument.

**These tools delete and modify observations, and none of them has a `--dryrun` flag.**
`dedupe`, `purge`, and `position` are instead gated by a `debug` constant at the top of their
source file: when `debug` is `true` the tool only logs what it would do, and when `debug` is
`false` it makes the changes for real. As checked in, `purge` deletes for real, while `dedupe`
and `position` are in log-only mode. `poke` and `repair` have no such guard and always write.
Read the source before running any of these against your account.

# Development

Run the tests with:
```
go test ./...
```
The tests use fake eBird and iNaturalist clients and local HTTP test servers, so they don't
touch the network or your real observations.
