# inat-terms — provenance

Vendored copy of iNaturalist's terms of use and community guidelines. See
[../../sources.md](../../sources.md) for the adoption record and
[requirements.md](requirements.md) for the transcription.

## What was retrieved

| File | Origin | SHA-256 |
| --- | --- | --- |
| `terms.html` | <https://www.inaturalist.org/pages/terms> | `de9737c4108f4058f0f6173ed8ea95894a6d56a224a5e245ae411538a0154afb` |
| `community-guidelines.html` | <https://www.inaturalist.org/pages/community+guidelines> | `2f32109e84002e0208ecb86aa18331ed6b62d849c1729a9c4e816690626044b8` |
| `machine-generated-content.html` | <https://www.inaturalist.org/pages/machine_generated_content> | `57ebd02cf03f0eebc5823ddea52e7299b086acf0ff54d6b719a666f8fc642160` |

- **Retrieved:** 2026-08-09, **by the maintainer, saved from a browser** — the pages answer an
  automated request with HTTP 403 and a JavaScript challenge.
- **Retrieved by:** transcribed by an AI agent, unreviewed by counsel.

## Page revision dates

`machine-generated-content.html` carries "Revised on May 20, 2026 08:20 PM by tiwane". It is
the page [CR-012](../../decisions.md#cr-012--is-birdsync-machine-generated-content) turns on,
and the one most worth diffing on a refresh: it defines a rule whose stated penalty is account
suspension, and its scope has evidently been revised before.

## What this is not

**Not legal advice, and not reviewed by a lawyer.** It is a snapshot of published guidance,
kept so a future reader can see what the terms said when the requirements were written.

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
