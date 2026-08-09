# inat-api — provenance

Vendored copy of iNaturalist's API recommended practices. See
[../../sources.md](../../sources.md) for the adoption record and
[requirements.md](requirements.md) for the transcription.

## What was retrieved

| File | Origin | SHA-256 |
| --- | --- | --- |
| `api-recommended-practices.html` | <https://www.inaturalist.org/pages/api+recommended+practices> | `04bba1540762ab1c80404bea68f133dec9b6b1152d2cf0f3a0e9e33ff7e2aa62` |

- **Retrieved:** 2026-08-09, **by the maintainer, saved from a browser.**
- **Page's own revision line:** "Revised on February 27, 2025 06:23 PM by kueda".

## Why a human retrieved it

`www.inaturalist.org/pages/*` answers an automated request with HTTP 403 and a JavaScript
challenge — confirmed 2026-08-09 for this page, the terms, and the community guidelines. The
`v1/swagger.json` spec is fetchable but carries no rate-limit guidance.

Working around the challenge was rejected: evading a service's bot protection in order to
import its own rules for being a well-behaved client would be self-defeating. Writing the
limits from memory was also rejected — a plausible invented number is indistinguishable from a
real one once it has a requirement ID beside it. See
[process.md](../../process.md#import-style).

## Refreshing

Re-save from a browser, compare the hash, and read the diff. The page carries its own revision
date, which is the quickest way to tell whether anything moved.
