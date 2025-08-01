# birdsync
Birdsync syncs eBird observations and photos to iNaturalist.

# Requirements

You must download your data from eBird using
https://ebird.org/downloadMyData.
Save the zip file and unzip it to get the `MyEBirdData.csv` file.

To run birdsync, you'll need the Go language toolchain.
Download it from http://go.dev.


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
Birdsync will prompt you to enter your iNaturalist user name and [API token](https://www.inaturalist.org/users/api_token), which allow the tool to read and write your personal iNaturalist observations.

To skip these interactive steps, you can provide your iNaturalist user name and API token as environment variables, but remember that you need to refresh your [API token](https://www.inaturalist.org/users/api_token) every 24 hours:
```
export INAT_USER_ID=sameerajmani
export INAT_API_TOKEN=(the TOKEN from {"api_token":"TOKEN"})
```
Birdsync provides command-line flags to customize its behavior:
*  `-after 2006-01-02`
        Sync only observations observed after the provided date and time (formatted as "2006-01-02 15:04:05"). The time can be omitted (2006-01-02).
* `-before 2006-01-02`
        Sync only observations observed before the provided date and time (formatted as "2006-01-02 15:04:05"). The time can be omitted (2006-01-02).
* `-dryrun`
        Don't actually sync any observations, just log what birdsync would do
* `-verifiable`
        Sync only observations that include Macaulay Catalog Numbers (photos or sound)
* `-fuzzy`
        Don't create a birdsync observation if a non-birdsync observation already exists for the same bird on the same date. This fuzzy matching is useful when you've entered the same observation manually into both eBird and iNaturalist, but it may skip legitimate uploads if you saw the same bird twice on the same day.

On the command line, flags must be listed _before_ your `MyEBirdData.csv` file:
```
$HOME/go/bin/birdsync --verifiable --after 2025-07-01 MyEBirdData.csv
```

# How birdsync works

Given (`iNaturalist user name`, `eBird CSV file`):
- Download all iNaturalist observations for `iNaturalist user name` into memory
- Index these iNaturalist observations by ([eBird submission ID](https://www.inaturalist.org/observation_fields/6033), [eBird scientific name](https://www.inaturalist.org/observation_fields/20215))
- Index any non-birdsync observations by date and common name for fuzzy matching
- For each eBird observation in `eBird CSV file`:
  - Skip any eBird observations that have already been uploaded
  - If `--after` is set, skip any eBird observations before that date
  - If `--before` is set, skip any eBird observations after that date
  - If `--verifiable` is set, skip any eBird observations lacking photos
  - If `--fuzzy` is set, skip any eBird observations for the same bird and day as a non-birdsync observation
  - Create a new iNaturalist observation from the eBird observation
  - For each [Macaulay Library](https://www.macaulaylibrary.org/) catalog ID for this eBird observation:
    - Download the image from the Macaulay Library
    - Upload the image to iNaturalist, associated with the new observation

# Limitations

Birdsync only works in the eBird → iNaturalist direction because (as far as I can tell) the [eBird API](https://support.ebird.org/en/support/solutions/articles/48000838205-download-ebird-data#API) doesn't support reading or writing personal checklists, only reading "limited, recent and summary outputs of eBird data".

Birdsync cannot detect whether iNaturalist observations that you've created manually are duplicates of those in your eBird checklists unless you mark your existing iNaturalist observations with the [eBird submission ID](https://www.inaturalist.org/observation_fields/6033) and [eBird scientific name](https://www.inaturalist.org/observation_fields/20215) observation fields. The `--fuzzy` matching feature provides a convenient way to avoid creating duplicates, but it may also suppress creating legitimate observations if you happened to see the same bird twice on the same day and entered it once into each tool.
