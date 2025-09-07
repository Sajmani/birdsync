# Architecture of birdsync

This document outlines the architecture of the `birdsync` command-line tool.

## Overview

`birdsync` is a Go application that synchronizes bird observation data from eBird to iNaturalist. It reads an eBird data export (CSV file) and creates corresponding observations in iNaturalist, including media attachments.

## High-Level Data Flow

1.  The user runs the `birdsync` command, providing the path to their eBird CSV data export.
2.  The application parses the CSV file to extract individual bird observations.
3.  For each observation, it fetches associated media (photos and sounds) from the Macaulay Library.
4.  It then connects to the iNaturalist API to create a new observation with the data and media.
5.  Flags provided by the user can modify this behavior, for example, to perform a dry run or filter observations by date.

## Package Structure

The codebase is organized into several packages, each with a distinct responsibility:

-   **`main`**: This is the entry point of the application. It contains the `main` function, handles command-line flag parsing, and orchestrates the overall synchronization process by coordinating the other packages.
    -   `birdsync.go`: Contains the core logic for the main application.
    -   `glue.go`: Contains helper functions that connect different parts of the application.

-   **`ebird`**: This package is responsible for all interactions with eBird data.
    -   `ebird/ebird.go`: Contains the logic for parsing the eBird CSV data export file into Go structs that the application can use.

-   **`inat`**: This package provides a client for the iNaturalist API.
    -   `inat/client.go`: An API client for making requests to the iNaturalist API.
    -   `inat/inat.go`: Contains higher-level functions for creating observations and handling other iNaturalist-specific logic.
    -   `inat/types.go`: Defines the Go data structures that map to iNaturalist API objects.
    -   `inat/vars.go`: Holds variables and constants used by the `inat` package.

-   **`media`**: This package handles media processing.
    -   `media.go`: Contains functions for downloading photos and sounds from the Macaulay Library, which are linked in the eBird data.

-   **`tools`**: This directory contains a collection of utility sub-packages for various data manipulation and management tasks. Each sub-package is a self-contained tool that can be used by the main application.
    -   `tools/dedupe`: Logic to prevent creating duplicate observations.
    -   `tools/dump`: Utility to dump data for debugging.
    -   `tools/poke`: A tool to "poke" or check on specific data.
    -   `tools/position`: Helper for handling geospatial positioning data.
    -   `tools/purge`: Tool for purging or deleting data.
    -   `tools/repair`: Tool for fixing or repairing data inconsistencies.

## Testing

-   `birdsync_test.go`: Tests for the main application logic.
-   `media_test.go`: Tests for the media handling logic.

This modular architecture separates concerns, making the codebase easier to understand, maintain, and extend.
