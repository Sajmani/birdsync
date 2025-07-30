# birdsync
Sync eBird observations and photos to iNaturalist

As a prerequisite, you must download your data from eBird using
https://ebird.org/downloadMyData — save the zip file and unzip it
to get the MyEBirdData.csv file.

To run birdsync, you'll need the Go language toolchain.
Download it from go.dev.
Then you can run the tool as
`go run github.com/Sajmani/birdsync MyEBirdData.csv`.
The tool will prompt you to enter your iNaturalist user name and API token.
Alternatively, you can provide these as environment variables.

Birdsync works as follows:
- download all iNaturalist observations for INAT_USER_ID into memory
- index these observations by (eBird checklist ID, species name)
- read eBird observations from the CSV file provided as a command line argument
- For each eBird observation:
  - skip* any eBird observations that have already been uploaded
  - create the iNaturalist observation
  - For each Macaulay Library ID for this eBird observation:
    - Download the image from the Macaulay Library
    - Upload the image to iNaturalist, associated with the new observation

**Known limitation**: Since we detect previously synced observations using
(eBird checklist ID, species name), we will reupload an observation if
the species name is changed in iNaturalist. This may happen based on
iNaturalist community idenfitications, resulting in duplicates.
As far as I can tell, there are no other fields in the eBird CSV
export that we can use to detect duplicate observations more reliably.
