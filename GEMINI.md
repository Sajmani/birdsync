# GEMINI.md

## Project Overview

This project, `birdsync`, is a command-line tool written in Go. Its primary purpose is to synchronize bird observation data from eBird to iNaturalist. It reads a CSV file exported from eBird, processes the observations, and creates corresponding entries in iNaturalist, including photos and sounds from the Macaulay Library.

The project is structured into several packages:
- `main`: The main application package.
- `ebird`: Contains logic for parsing eBird data.
- `inat`: Contains the client for interacting with the iNaturalist API.
- `tools`: Contains various utility tools for data manipulation.

## Building and Running

### Prerequisites

- Go 1.23.0 or higher.
- An eBird data export file (`MyEBirdData.csv`).

### Installation

To install the `birdsync` tool, run the following command:

```bash
go install github.com/Sajmani/birdsync@latest
```

This will install the binary in your `$GOPATH/bin` directory.

### Running

To run the tool, you need to provide the path to your eBird CSV file:

```bash
birdsync [flags] MyEBirdData.csv
```

**Common Flags:**

- `--dryrun`: Perform a dry run without creating any iNaturalist observations.
- `--verifiable`: Only sync observations with associated media (photos or sounds).
- `--fuzzy`: Avoid creating duplicate observations for the same bird on the same day.
- `--after YYYY-MM-DD`: Only sync observations after a specific date.
- `--before YYYY-MM-DD`: Only sync observations before a specific date.

For more details on the available flags, refer to the `README.md` file or run `birdsync --help`.

## Development Conventions

### Code Style

The codebase follows standard Go formatting and conventions. It is recommended to use `gofmt` to format your code before committing.

### Dependencies

Dependencies are managed using Go modules. The `go.mod` file lists the direct and indirect dependencies. To add a new dependency, use `go get`.

### Testing

The project does not appear to have a dedicated test suite in the provided file structure. If you add new features, it is recommended to also add corresponding tests.

### Contribution Guidelines

There are no explicit contribution guidelines in the repository. However, it is recommended to follow the existing code style and to open an issue to discuss any major changes before submitting a pull request.
