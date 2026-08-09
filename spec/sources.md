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
| Version | Retrieved 2026-08-09 by the maintainer from a browser; hash in [PROVENANCE.md](sources/inat-api-2026-08-09/PROVENANCE.md) |
| Tier | **Governing** |
| Import style | By reference, vendored at `sources/inat-api-2026-08-09/` |
| Scope | All requests to the iNaturalist API |
| Adopted parts | `inat-api/R1`–`R8`, in [requirements.md](sources/inat-api-2026-08-09/requirements.md) |
| Exclusions | Bulk access routes — observation exports, the GBIF dataset — which serve consumers of iNaturalist data, not a user syncing their own |
| Owner | iNaturalist |

The service operator can rate-limit or block a client that ignores these, so they override
local convenience. Eight requirements transcribed, of which six are satisfied: T-016
(`User-Agent`), T-017 (`per_page=200`), T-035 (pacing), T-036 (`id_above` paging), and two
that hold by construction. `inat-api/R7` — authenticate only when necessary — is a **knowing
departure**, recorded with its reason rather than left as an oversight.
### `inat-terms` — iNaturalist terms of service and community guidelines

| Field | Value |
| --- | --- |
| Origin | <https://www.inaturalist.org/pages/terms>, <https://www.inaturalist.org/pages/community+guidelines> |
| Version | Retrieved 2026-08-09 by the maintainer from a browser; hashes in [PROVENANCE.md](sources/inat-terms-2026-08-09/PROVENANCE.md) |
| Tier | **Mandatory** |
| Import style | By reference, vendored at `sources/inat-terms-2026-08-09/` |
| Scope | Everything birdsync writes to a user's account, and the manner of writing it |
| Adopted parts | `inat-terms/R1`–`R6`, in [requirements.md](sources/inat-terms-2026-08-09/requirements.md) |
| Exclusions | Conduct rules birdsync cannot reach: harassment, sockpuppets, explicit content, commercial AI training |
| Owner | iNaturalist |
| Also vendored | The definition of machine generated content, which the guidelines link to and CR-012 turns on |

**Vendored and transcribed 2026-08-09.** Raised the sharpest question in this specification
and then answered it: the guidelines prohibit machines posting content "with no human oversight
curating each piece", on pain of account suspension, but the linked definition lists writing a
script to post observations from your own curated data among its examples of *acceptable*
behavior. [CR-012](decisions.md#cr-012--is-birdsync-machine-generated-content), closed.

Also `inat-terms/R3`, a rate-and-volume clause in the terms themselves, which binds harder than
the API guidance behind T-035, and `inat-terms/R5`, which puts an obligation on the user that
became P-067.

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

All three by-reference sources are now vendored and pinned. The two iNaturalist ones could
not be retrieved by an agent and were saved by the maintainer from a browser: `www.inaturalist.org/pages/*` answers an automated request with HTTP
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
