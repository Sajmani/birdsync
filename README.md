# birdsync
Sync eBird observations and photos to iNaturalist

As a prerequisite, you must download your data from eBird using
https://ebird.org/downloadMyData — save the zip file and unzip it
to get the MyEBirdData.csv file.

To run birdsync, you'll need the Go language toolchain.
Download it from http://go.dev.
Then you can run the tool as
`go run github.com/Sajmani/birdsync MyEBirdData.csv`.
The tool will prompt you to enter your iNaturalist user name and API token.
Alternatively, you can provide these as environment variables.

Birdsync works as follows:
- Download all iNaturalist observations for INAT_USER_ID into memory
- Index these observations by (eBird checklist ID, species name)
- Read eBird observations from the CSV file provided as a command line argument
- For each eBird observation:
  - Skip any eBird observations that have already been uploaded
  - Create the iNaturalist observation
  - For each Macaulay Library ID for this eBird observation:
    - Download the image from the Macaulay Library
    - Upload the image to iNaturalist, associated with the new observation
