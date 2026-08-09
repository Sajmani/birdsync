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
| Version | Not yet pinned |
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
| Version | Not yet pinned |
| Tier | **Mandatory** |
| Import style | By reference |
| Scope | Everything birdsync writes to a user's account |
| Adopted parts | Content ownership and licensing; automated posting |
| Owner | iNaturalist |

Governs what may be uploaded and under whose authority. Relevant to
[product.md open question 4](product.md#open-questions): birdsync uploads media
originating in the Macaulay Library into a user's iNaturalist account.

### `ml-terms` — Macaulay Library / eBird terms of use

| Field | Value |
| --- | --- |
| Origin | <https://www.birds.cornell.edu/home/terms-of-use/>, eBird data-access policy |
| Version | Not yet pinned |
| Tier | **Mandatory** |
| Import style | By reference |
| Scope | Downloading assets from the Macaulay Library CDN and re-uploading them |
| Adopted parts | Permitted use of downloaded assets; attribution |
| Owner | Cornell Lab of Ornithology |

The working assumption is that a user's own export references their own assets, making
the copy a personal data transfer. That assumption has not been checked against the
actual terms, and it is the one place where birdsync could put a user in the wrong. This
is the highest-value source to vendor first.

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

**No source is vendored.** Every by-reference source above is declared but unpinned,
which does not satisfy [process.md](process.md#import-style)'s requirement that a
by-reference source be vendored and pinned.

This is a deliberate deferral, not an oversight: fetching and transcribing three sets of
terms of service is a substantial task with legal consequence, and the transcription is
itself a reviewable artifact needing the owner's sign-off. Doing it uninstructed would
produce exactly the kind of unreviewed authority the process is meant to prevent.

Proposed order, highest risk first:

1. `ml-terms` — the only source that could put a user in the wrong.
2. `inat-terms` — same category, lower exposure.
3. `inat-api` — mostly already implemented; vendoring mainly adds the rate limits.

`go-practice` needs no vendoring, being absorbed.
