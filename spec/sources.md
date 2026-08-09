# Requirement sources — birdsync

Whose requirements apply to this project, and who outranks whom. Written per
[process.md](process.md#sources-and-composition); see it for what the tiers mean.

**Status: draft, pending Gate 1.** No source has been vendored yet — see
[Vendoring status](#vendoring-status).

## Adopted

### `inat-api` — iNaturalist API recommended practices

| Field | Value |
| --- | --- |
| Origin | <https://www.inaturalist.org/pages/api+recommended+practices> |
| Version | **Not pinned — retrieval blocked** (HTTP 403 with a JavaScript challenge, 2026-08-09) |
| Tier | **Governing** |
| Import style | By reference |
| Scope | All requests to the iNaturalist API |
| Adopted parts | Client identification, page size, request rate |
| Owner | iNaturalist |

The service operator can rate-limit or block a client that ignores these, so they
override local convenience. Already reflected in T-016 (`User-Agent`) and T-017
(`per_page=200`); the rate-limit guidance is **not yet expressed as a requirement**,
which is a gap to close in phase 2.

### `inat-terms` — iNaturalist terms of service and community guidelines

| Field | Value |
| --- | --- |
| Origin | <https://www.inaturalist.org/pages/terms>, <https://www.inaturalist.org/pages/community+guidelines> |
| Version | **Not pinned — retrieval blocked** (HTTP 403 with a JavaScript challenge, 2026-08-09) |
| Tier | **Mandatory** |
| Import style | By reference |
| Scope | Everything birdsync writes to a user's account |
| Adopted parts | Content ownership and licensing; automated posting |
| Owner | iNaturalist |

Governs what may be uploaded and under whose authority. Relevant to
[product.md open question 4](product.md#open-questions): birdsync uploads media
originating in the Macaulay Library into a user's iNaturalist account.

### `ml-terms` — Macaulay Library / eBird media terms

| Field | Value |
| --- | --- |
| Origin | Three eBird Help Center pages; see [PROVENANCE.md](sources/ml-terms-2026-08-09/PROVENANCE.md) |
| Version | Retrieved 2026-08-09; SHA-256 of each page recorded in PROVENANCE.md |
| Tier | **Mandatory** |
| Import style | By reference, vendored at `sources/ml-terms-2026-08-09/` |
| Scope | Downloading assets from the Macaulay Library CDN and re-uploading them |
| Adopted parts | `ml-terms/R1`–`R5`, in [requirements.md](sources/ml-terms-2026-08-09/requirements.md) |
| Exclusions | The media request and licensing process: birdsync never requests assets it lacks IDs for, and makes no commercial use |
| Owner | Cornell Lab of Ornithology |

**Vendored and transcribed 2026-08-09.** The working assumption — that a user's export
references only their own assets, making the copy a personal data transfer — is documented
and evidence-supported rather than enforced (CR-010, closed). A contributor keeps copyright in their own media (`ml-terms/R1`) and may
download it freely (`R2`), so birdsync's intended case is squarely permitted, which answers
[product.md](product.md#open-questions)'s first open question. But another contributor's media
may not be downloaded without permission (`R3`), and birdsync cannot tell the difference: it
fetches by asset ID from an unauthenticated CDN that serves anyone's asset. See
[CR-010](decisions.md#cr-010--birdsync-cannot-tell-whose-media-it-is-downloading).

The transcription is an interpretation by an AI agent, unreviewed by counsel, and says so at
the top of the file.

### `go-practice` — Go language and toolchain best practices

| Field | Value |
| --- | --- |
| Origin | Effective Go, Go Code Review Comments, the `gofmt` and `go vet` defaults |
| Version | Tracks the toolchain in `go.mod` |
| Tier | **Advisory** |
| Import style | Absorbed |
| Scope | All Go source in this repository |
| Adopted parts | Formatting, error wrapping, naming, comment style |
| Owner | The Go team; locally, the repository owner |

Absorbed rather than referenced: the relevant guidance is already written into T-024
through T-029 as local requirements. Advisory tier means a local requirement wins where
the two differ — for example T-027's rule about `log.Fatal`, which is stricter than
general Go practice, and T-002's dependency ceiling, which is far stricter.

## Considered and not adopted

| Source | Why not |
| --- | --- |
| Accessibility standards (WCAG, Section 508) | birdsync has no graphical interface. Its output is uncolored plain text on stdout/stderr, which is already screen-reader compatible. Revisit if a UI or colored output is ever added. |
| Brand or design guidelines | None exist; this is a personal open-source project. |
| Privacy regulation (GDPR, CCPA) | birdsync collects, stores, and transmits nothing to its author. It moves a user's own data between two services the user has accounts with, running entirely on the user's machine. No controller or processor relationship arises. |
| Semantic versioning | There are no release tags; users install `@latest`. Worth revisiting if tagged releases begin. |
| eBird API terms | Not adopted because birdsync does not use the eBird API — the CSV export exists precisely because personal checklists are not exposed (P-003). |

Recording these matters as much as recording the adopted ones: each is a question that
will otherwise be re-asked every year.

## Vendoring status

`ml-terms` is vendored and pinned. The two iNaturalist sources are **not**, and cannot be
retrieved by an agent: `www.inaturalist.org/pages/*` answers an automated request with HTTP
403 and a JavaScript challenge (checked 2026-08-09 for the terms, the community guidelines,
and the API recommended practices). The `v1/swagger.json` spec is fetchable but carries no
rate-limit guidance.

Working around the challenge is not on the table — evading a service's bot protection in
order to import its own rules for being a well-behaved client would be self-defeating. Nor is
writing the limits from memory: a plausible invented number is indistinguishable from a real
one once it is in a requirements file, which is how CR-009 happened. These stay unpinned
until a human saves the pages.

This is a deliberate deferral, not an oversight: fetching and transcribing three sets of
terms of service is a substantial task with legal consequence, and the transcription is
itself a reviewable artifact needing the owner's sign-off. Doing it uninstructed would
produce exactly the kind of unreviewed authority the process is meant to prevent.

Proposed order, highest risk first:

1. ~~`ml-terms`~~ — done 2026-08-09.
2. `inat-terms` — same category, lower exposure.
3. `inat-api` — mostly already implemented; vendoring mainly adds the rate limits.

`go-practice` needs no vendoring, being absorbed.
