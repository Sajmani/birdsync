package tech

# Tech requirements: Resource use

Memory and disk. One item per heading, and each item's first paragraph after its fields is the
complete rule. Search `^## ` for the list, or `## <name>` for one item. Other tech requirements
are in the other files in this directory, all one package, `tech`.

## memoryScalesWithExport

**Memory scales with the export**

Formerly: T-022
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

Memory scales with the export: the CSV is read whole with `ReadAll`, and all
downloaded iNaturalist observations are held in memory. This is accepted at the current
scale (tens of thousands of records) and is a known ceiling, not a design goal.

## tempFilesDeleted

**Temporary media files are deleted**

Formerly: T-023
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

Temporary files created while downloading media are deleted before the run
ends.

*The file belongs to the caller once `DownloadMLAsset` returns it, so `birdsync()` removes
it after the upload attempt, successful or not — at most one asset is on disk at a time.*

*The download's error paths honor it too: `downloadMLAsset` closes and removes the temp file
on every path but success. Leaving it behind — and on Windows leaving the handle open — is what
[issue #1](https://github.com/Sajmani/birdsync/issues/1) reported as "The process cannot access
the file because it is being used by another process" on a later attempt.*
