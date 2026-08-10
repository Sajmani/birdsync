# inat-api — provenance

Vendored copy of iNaturalist's API recommended practices. See
[../../sources.md](../../sources.md) for the adoption record and
[requirements.md](requirements.md) for the transcription.

## What was retrieved

| File | Origin | SHA-256 |
| --- | --- | --- |
| `api-recommended-practices.html` | <https://www.inaturalist.org/pages/api+recommended+practices> | `ffb53ae11e1cce9c8bc2cef7743d478318ae7f9525eccde68d56d9eab4f7dbd3` |

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

## Active content removed

The `.html` files here are the documents **with `<script>`, `<style>` and `<link>` elements
removed**. Nothing else was altered: text, structure and attributes are as served.

This is not cosmetic. A browser-saved iNaturalist page embeds the site's Google Maps
JavaScript API key in a `<script src=...>` tag, and committing it here published a third
party's credential and tripped GitHub secret scanning (alert #1). The key belongs to
iNaturalist and appears in every page they serve, so nothing became public that was not
already, but redistributing someone else's key is not this repository's business — and a
terms-of-service snapshot has no use for scripts.

Both hashes are recorded so the stripping can be reproduced and audited:

| File | SHA-256 as served | SHA-256 as stored |
| --- | --- | --- |
| `api-recommended-practices.html` | `04bba1540762ab1c80404bea68f133dec9b6b1152d2cf0f3a0e9e33ff7e2aa62` | `ffb53ae11e1cce9c8bc2cef7743d478318ae7f9525eccde68d56d9eab4f7dbd3` |
| `community-guidelines.html` | `dc84d9b271dabf10dc2c52512e50a5ac2bb4ff23711a36cb73ac26282eb5ce84` | `2f32109e84002e0208ecb86aa18331ed6b62d849c1729a9c4e816690626044b8` |
| `machine-generated-content.html` | `b2453ade68999662bb893f8b1b534b040545821a860f053670b8d88d8ae38cbc` | `57ebd02cf03f0eebc5823ddea52e7299b086acf0ff54d6b719a666f8fc642160` |
| `terms.html` | `3ce69d6eac87fcc93ccca8b67a4b97f93991291ca2509467781ab8adc593122c` | `de9737c4108f4058f0f6173ed8ea95894a6d56a224a5e245ae411538a0154afb` |

`TestTranscribedQuotesAppearInSources` checks that every passage quoted in
`requirements.md` still appears in these files, so the stripping cannot quietly remove
something a requirement depends on.
