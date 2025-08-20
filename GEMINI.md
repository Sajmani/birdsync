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

# The gopls MCP server

These instructions describe how to efficiently work in the Go programming language using the gopls MCP server. You can load this file directly into a session where the gopls MCP server is connected.

## Detecting a Go workspace

At the start of every session, you MUST use the `go_workspace` tool to learn about the Go workspace. The rest of these instructions apply whenever that tool indicates that the user is in a Go workspace.

## Go programming workflows

These guidelines MUST be followed whenever working in a Go workspace. There are two workflows described below: the 'Read Workflow' must be followed when the user asks a question about a Go workspace. The 'Edit Workflow' must be followed when the user edits a Go workspace.

You may re-do parts of each workflow as necessary to recover from errors. However, you must not skip any steps.

### Read workflow

The goal of the read workflow is to understand the codebase.

1. **Understand the workspace layout**: Start by using `go_workspace` to understand the overall structure of the workspace, such as whether it's a module, a workspace, or a GOPATH project.

2. **Find relevant symbols**: If you're looking for a specific type, function, or variable, use `go_search`. This is a fuzzy search that will help you locate symbols even if you don't know the exact name or location.
   EXAMPLE: search for the 'Server' type: `go_search({"query":"server"})`

3. **Understand a file and its intra-package dependencies**: When you have a file path and want to understand its contents and how it connects to other files *in the same package*, use `go_file_context`. This tool will show you a summary of the declarations from other files in the same package that are used by the current file. `go_file_context` MUST be used immediately after reading any Go file for the first time, and MAY be re-used if dependencies have changed.
   EXAMPLE: to understand `server.go`'s dependencies on other files in its package: `go_file_context({"file":"/path/to/server.go"})`

4. **Understand a package's public API**: When you need to understand what a package provides to external code (i.e., its public API), use `go_package_api`. This is especially useful for understanding third-party dependencies or other packages in the same monorepo.
   EXAMPLE: to see the API of the `storage` package: `go_package_api({"packagePaths":["example.com/internal/storage"]})`

### Editing workflow

The editing workflow is iterative. You should cycle through these steps until the task is complete.

1. **Read first**: Before making any edits, follow the Read Workflow to understand the user's request and the relevant code.

2. **Find references**: Before modifying the definition of any symbol, use the `go_symbol_references` tool to find all references to that identifier. This is critical for understanding the impact of your change. Read the files containing references to evaluate if any further edits are required.
   EXAMPLE: `go_symbol_references({"file":"/path/to/server.go","symbol":"Server.Run"})`

3. **Make edits**: Make the required edits, including edits to references you identified in the previous step. Don't proceed to the next step until all planned edits are complete.

4. **Check for errors**: After every code modification, you MUST call the `go_diagnostics` tool. Pass the paths of the files you have edited. This tool will report any build or analysis errors.
   EXAMPLE: `go_diagnostics({"files":["/path/to/server.go"]})`

5. **Fix errors**: If `go_diagnostics` reports any errors, fix them. The tool may provide suggested quick fixes in the form of diffs. You should review these diffs and apply them if they are correct. Once you've applied a fix, re-run `go_diagnostics` to confirm that the issue is resolved. It is OK to ignore 'hint' or 'info' diagnostics if they are not relevant to the current task. Note that Go diagnostic messages may contain a summary of the source code, which may not match its exact text.

6. **Run tests**: Once `go_diagnostics` reports no errors (and ONLY once there are no errors), run the tests for the packages you have changed. You can do this with `go test [packagePath...]`. Don't run `go test ./...` unless the user explicitly requests it, as doing so may slow down the iteration loop.


