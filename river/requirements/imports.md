# Imported sources

Every source birdsync imports, every dependency found when River was adopted, and every source
considered and turned down, with the reason. One entry per source: search `` ^## ` `` for the
list. Each version is recorded here and nowhere else. Per-file hashes live in each source's
`PROVENANCE.md`, and the `Integrity` field points there rather than copying them. Tiers are
in [`precedence.md`](precedence.md).

The first five entries were carried over from specman's source manifest (`spec/sources.md` at
[`dc25e3b`](https://github.com/Sajmani/birdsync/blob/dc25e3b6756bd59c2de931e6d080d47318617da4/spec/sources.md)),
where the owner adopted or retired them. The entries for `uuid` and `github-actions` were found
by the discovery step of adopting River, on 2026-11-03. The owner adopted `uuid` on 2026-11-04;
`github-actions` stays partial until someone reads GitHub's terms for Actions.

## `inatapi`

| Field | Value |
| --- | --- |
| Name | iNaturalist API recommended practices |
| Origin | <https://www.inaturalist.org/pages/api+recommended+practices> |
| Version | Retrieved 2026-08-09 by the maintainer from a browser |
| Integrity | SHA-256 per file in [`PROVENANCE.md`](imported/inaturalist.org/inatapi/PROVENANCE.md) |
| Tier | Governing |
| Import style | vendored |
| Scope | Every request to the iNaturalist API, `api.inaturalist.org` |
| Adopted parts | All eight transcribed items in [`inatapi.md`](imported/inaturalist.org/inatapi/inatapi.md) |
| Owner | iNaturalist |
| Status | partial: the daily half of `requestRate` is not enforced, and no accepted-gap decision covers it yet. `authenticateOnlyWhenNeeded` is met ([decisions.authenticationIsNecessary](decisions.md#authenticationisnecessary)) |

Bulk access routes, such as observation exports and the GBIF dataset, serve consumers of
iNaturalist data rather than a user syncing their own, and are not adopted.

## `inatterms`

| Field | Value |
| --- | --- |
| Name | iNaturalist terms of use, community guidelines, and machine-generated-content definition |
| Origin | <https://www.inaturalist.org/pages/terms>, <https://www.inaturalist.org/pages/community+guidelines>, <https://www.inaturalist.org/pages/machine_generated_content> |
| Version | Retrieved 2026-08-09 by the maintainer from a browser |
| Integrity | SHA-256 per file in [`PROVENANCE.md`](imported/inaturalist.org/inatterms/PROVENANCE.md) |
| Tier | Mandatory |
| Import style | vendored |
| Scope | Everything birdsync writes to a user's account on `inaturalist.org` and `www.inaturalist.org`, and the manner of writing it |
| Adopted parts | All six transcribed items in [`inatterms.md`](imported/inaturalist.org/inatterms/inatterms.md) |
| Owner | iNaturalist |
| Status | adopted |

The conduct rules birdsync cannot reach (harassment, sockpuppets, explicit content, commercial
AI training) bind the user but are not adopted.

## `mlterms`

| Field | Value |
| --- | --- |
| Name | Macaulay Library and eBird media terms |
| Origin | Three eBird Help Center pages, listed in [`PROVENANCE.md`](imported/ebird.org/mlterms/PROVENANCE.md) |
| Version | Retrieved 2026-08-09 by `curl`, as served |
| Integrity | SHA-256 per file in [`PROVENANCE.md`](imported/ebird.org/mlterms/PROVENANCE.md) |
| Tier | Mandatory |
| Import style | vendored |
| Scope | Downloading assets from the Macaulay Library CDN, `cdn.download.ams.birds.cornell.edu`, and re-uploading them |
| Adopted parts | All five transcribed items in [`mlterms.md`](imported/ebird.org/mlterms/mlterms.md) |
| Owner | Cornell Lab of Ornithology |
| Status | partial: whether birdsync meets `attributionExpected` is unchecked, because the attribution guidelines page is not vendored |

The media request and licensing process is not adopted: birdsync never requests assets it
lacks IDs for, and makes no commercial use.

## `go-practice`

| Field | Value |
| --- | --- |
| Name | Go language and toolchain best practices |
| Origin | Effective Go, Go Code Review Comments, the `gofmt` and `go vet` defaults |
| Version | Tracks the toolchain in `go.mod` |
| Integrity | none; absorbed |
| Tier | Advisory |
| Import style | absorbed |
| Scope | All Go source in this repository |
| Adopted parts | Formatting, error wrapping, naming, comment style, written into `tech` as [tech.gofmtClean](local/tech/tech.md#gofmtclean) through [tech.commentsExplainWhy](local/tech/tech.md#commentsexplainwhy) |
| Owner | The Go team; locally, the repository owner |
| Status | adopted |

## `specman`

| Field | Value |
| --- | --- |
| Name | specman, the spec-driven development method |
| Origin | <https://github.com/Sajmani/specman>, file `process.md` |
| Version | Base pin `sha256:63bfdecb99a256d2819aa6d0e55bfb0bccf19b96db40e28f72572315e81cbdbc`, taken 2026-09-23 |
| Integrity | the pin above |
| Tier | Advisory |
| Import style | quoted (copied to `spec/process.md`) |
| Scope | How work on this repository was done |
| Adopted parts | In full, until River replaced it |
| Owner | The repository owner |
| Status | retired, 2026-11-03, when birdsync adopted River. River's method is the skill at `.agents/skills/river/`, pinned in its `PROVENANCE.md` rather than here. What depended on specman: every former `P-###`, `T-###`, `AC-###` and `CR-###` identifier, now kept as `Formerly:` fields, and the `spec/` directory, now retired. The copy survives at [`dc25e3b`](https://github.com/Sajmani/birdsync/blob/dc25e3b6756bd59c2de931e6d080d47318617da4/spec/process.md) |

## `uuid`

| Field | Value |
| --- | --- |
| Name | github.com/google/uuid |
| Origin | <https://github.com/google/uuid> |
| Version | v1.6.0 (`go.mod`) |
| Integrity | `go.sum`: `h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=` |
| Tier | Advisory |
| Import style | fetched at build |
| Scope | Generating observation UUIDs (`github.com/google/uuid` is the module's only external dependency; see [tech.oneExternalDependency](local/tech/tech.md#oneexternaldependency)) |
| Adopted parts | Its license |
| Owner | Google, per its copyright line |
| Status | adopted ([decisions.uuidAdopted](decisions.md#uuidadopted)) |

Discovery, 2026-11-03: read the `LICENSE` file shipped in the module at v1.6.0, from the Go
module cache (SHA-256 `0a8d61ed3cbfd5312326e8126c31ce9c627a283adc99131b56896d29ada04b2d`). It
is a three-clause BSD-style license, "Copyright (c) 2009,2014 Google Inc.", that permits
redistribution in source and binary forms if the copyright notice, conditions, and disclaimer
are kept. The module ships no `NOTICE` file. birdsync redistributes no binaries, since users
build with `go install`, so the notice obligation does not arise today.

## `github-actions`

| Field | Value |
| --- | --- |
| Name | GitHub Actions, with `actions/checkout@v7` and `actions/setup-go@v7` |
| Origin | <https://github.com/features/actions>; the two actions' repositories |
| Version | `@v7` tags, as used in `.github/workflows/ci.yml` (a tag is not a pin) |
| Integrity | none recorded |
| Tier | Advisory |
| Import style | consulted |
| Scope | Continuous integration ([tech.ciStandingChecks](local/tech/tech.md#cistandingchecks)) |
| Adopted parts | none yet |
| Owner | GitHub |
| Status | partial: found during discovery and awaiting the owner's decision. Terms not retrieved |

Discovery, 2026-11-03: one web search, "GitHub Actions terms for additional products and
features usage limits". It returned a search summary, not the owner's pages. That summary
pointed to GitHub's *Terms for Additional Products and Features* and its Actions usage-limits
documentation, saying that runners in public repositories are free and that some uses are
prohibited. **Neither page was retrieved or read**, so none of their terms are recorded here.
The licenses of the two actions were not retrieved either. Recorded as a gap.

## Considered and not adopted

| Source | Status | Why not |
| --- | --- | --- |
| Accessibility standards (WCAG, Section 508) | rejected | birdsync has no graphical interface. Its output is uncolored plain text on stdout and stderr. Revisit if a UI or colored output is added |
| Brand or design guidelines | rejected | None exist; this is a personal open-source project |
| Privacy regulation (GDPR, CCPA) | rejected | birdsync collects, stores, and transmits nothing to its author. It moves a user's own data between two services they have accounts with, on their own machine |
| Semantic versioning | rejected | There are no release tags; users install `@latest`. Revisit if tagged releases begin |
| eBird API terms | rejected | birdsync does not use the eBird API. The CSV export exists because personal checklists are not exposed ([product.csvExportInput](local/product/product.md#csvexportinput)) |
| `ebird.org` | rejected | Never called. birdsync writes the checklist URL `https://ebird.org/checklist/<id>` into a description as text. The input file is downloaded by the user, not by birdsync |
| `macaulaylibrary.org` | rejected | Never called. birdsync writes `https://macaulaylibrary.org/asset/<id>` into a description as text. Attribution for Macaulay assets is governed by `mlterms` |

The first five rows were carried over from specman. `ebird.org` and `macaulaylibrary.org` were
added by discovery on 2026-11-03, because `river check` reported them as hosts.
`www.inaturalist.org` and `inaturalist.org` are covered by `inatterms` above.
