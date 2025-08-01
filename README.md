# birdsync
Birdsync syncs eBird observations and photos to iNaturalist.

# Requirements

You must download your data from eBird using
https://ebird.org/downloadMyData.
Save the zip file and unzip it to get the `MyEBirdData.csv` file.

To run birdsync, you'll need the Go language toolchain.
Download it from http://go.dev.


# Install and run birdsync

Install or update `birdsync` to the latest version using `go install github.com/Sajmani/birdsync@latest`, then run it from the command line, specifying the path to your `MyEBirdData.csv` file:
```
birdsync MyEBirdData.csv
```
Birdsync tool will prompt you to enter your iNaturalist user name and [API token](https://www.inaturalist.org/users/api_token), which allows the tool to read and write your personal iNaturalist observations.

To skip these interactive steps, you can provide your iNaturalist user name and API token as environment variables, but remember that you need to refresh your API token every 24 hours:
```
export INAT_USER_ID=sameerajmani
export INAT_API_TOKEN=(the TOKEN from {"api_token":"TOKEN"})
```
Birdsync provides command-line flags to customize its behavior:
*  `-after value`
        Sync only observations observed after the provided DateTime (2006-01-02 15:04:05). The time can be omitted (2006-01-02).
* `-before value`
        Sync only observations observed before the provided DateTime (2006-01-02 15:04:05). The time can be omitted (2006-01-02).
* `-dryrun`
        Don't actually sync any observations, just log what birdsync would do
* `-verifiable`
        Sync only observations that include Macaulay Catalog Numbers (photos or sound)
Command-line flags must be listed before your `MyEBirdData.csv` file:
```
birdsync --verifiable --after 2025-07-01 MyEBirdData.csv
```

# How it works

Given (`iNaturalist user name`, `eBird CSV file`):
- Download all iNaturalist observations for `iNaturalist user name` into memory
- Index these iNaturalist observations by ([eBird submission ID](https://www.inaturalist.org/observation_fields/6033), [eBird scientific name](https://www.inaturalist.org/observation_fields/20215))
- For each eBird observation in `eBird CSV file`:
  - Skip any eBird observations that have already been uploaded
  - If `--after` is set, skip any eBird observations before that date
  - If `--before` is set, skip any eBird observations after that date
  - If `--verifiable` is set, skip any eBird observations lacking photos
  - Create a new iNaturalist observation from the eBird observation
  - For each Macaulay Library ID for this eBird observation:
    - Download the image from the Macaulay Library
    - Upload the image to iNaturalist, associated with the new observation

# Limitations

Birdsync only works in the eBird → iNaturalist direction because (as far as I can tell) the [eBird API](https://support.ebird.org/en/support/solutions/articles/48000838205-download-ebird-data#API) doesn't support reading or writing personal checklists, only reading "limited, recent and summary outputs of eBird data".