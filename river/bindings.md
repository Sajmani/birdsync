# Bindings

How birdsync applies River. The method itself is the skill at
[`.agents/skills/river/`](../.agents/skills/river/SKILL.md), copied unchanged and pinned in its
[`PROVENANCE.md`](../.agents/skills/PROVENANCE.md). Everything true of this project rather than
of River goes here.

## Stages

Drafted at adoption on 2026-11-03 and confirmed by the owner, Sajmani, on 2026-11-04.

| Stage | Status | Notes |
| --- | --- | --- |
| Requirements | adopted | Sources in `requirements/local/` and `requirements/imported/`, composed into `requirements/composed/internal/birdsync/` |
| Verification | adopted | Criteria in the composed package's `acceptance.md`. Guardrails are Go tests and CI, mapped in [`conformance.md`](conformance.md) |
| Implementation | informal | The design record is [`docs/arch.md`](../docs/arch.md), kept in step with the code by convention, not by River |
| Release | informal | No release process: users install `@latest` |

## Tools

| | |
| --- | --- |
| River skill | Pinned in [`.agents/skills/PROVENANCE.md`](../.agents/skills/PROVENANCE.md), and only there |
| River CLI | Built from River at the skill's pin; CI reads the pin from `PROVENANCE.md`. The CLI reports no version of its own |
| Bundle | `river/` at the repository root |
| Where `river check` runs | CI (`.github/workflows/ci.yml`, job `river`), which is skipped until the repository variable `RIVER_REPO` names a clonable River repository, because River is not published yet. Also before calling any change done (`AGENTS.md`). There is no pre-commit hook and no release step |
| `goals/` | [`goals/goals.md`](goals/goals.md), five goals accepted by the owner on 2026-11-04 |

## Project bindings

Carried over from specman's project bindings (`spec/tech.md`, *Project bindings*).

| Binding | Value |
| --- | --- |
| Owner | The repository owner, Sajmani, approves every composed package, accepted gap, and decision |
| Standing checks | `go build ./...`, `go vet ./...`, `gofmt -l .` (must print nothing), `go test -race ./...`, `river check` |
| Hard safety constraints | [tech.dryRunGateAtCallSite](requirements/local/tech/tech.md#dryrungateatcallsite), [tech.noLiveServicesInTests](requirements/local/tech/tech.md#noliveservicesintests), [tech.credentialsNeverExposed](requirements/local/tech/tech.md#credentialsneverexposed), [tech.toolsReadOnly](requirements/local/tech/tech.md#toolsreadonly), [tech.neverRunTools](requirements/local/tech/tech.md#neverruntools) |
| External systems | `api.inaturalist.org` and `cdn.download.ams.birds.cornell.edu` are never contacted from a test. Use `httptest` through the base-URL seams |
| Subject vocabulary | `Subject:` fields are dotted paths rooted at `build`, `cli`, `http`, `inat`, `log`, `media`, `observation`, or `sync`. The owner adds roots |
| Source refresh | Imported sources are reviewed when a source changes, and at least once a year. `www.inaturalist.org/pages/*` refuses automated clients, so refreshing those needs a person with a browser |
| Existing docs | `README.md` owns usage, `CONTRIBUTING.md` owns workflow mechanics, `AGENTS.md` owns agent operating notes, `docs/arch.md` owns design, and `river/` owns requirements and criteria |

## Local conventions

Additions to the method, none of them changing it:

- **`Formerly:`** on a source item, criterion, or decision gives its identifier under specman
  (`P-###`, `T-###`, `AC-###`, `CR-###`, `<source>/R#`). Commits, issues, comments, and test
  messages from before 2026-11-03 still cite those identifiers. Search `Formerly: P-063` to
  resolve one.
- **`Area:`** on a local source item gives the heading it sat under in specman. The composed
  package groups by it.
- **`Subject:` and `Value:`** carry specman's machine-checkable subject and value, where an
  item had one.
- **[`proposals.md`](proposals.md)** holds drafts awaiting the owner: candidate goals,
  decisions, and approvals. River has no place for undecided drafts. Keeping them here keeps
  them out of `goals/` and `decisions.md`, which hold only what a person has decided.
- **Vendored HTML is scanned before commit.** Strip active content and credential-shaped
  strings, such as presigned `X-Amz-*` query parameters, from a vendored page before committing
  it. Record each redaction, with the hash as served and as stored, in that source's
  `PROVENANCE.md`. The `mlterms` pages predate this policy and are kept as served until their
  next re-fetch.
- **Hand composition:** with no composition tool, the composed package is rewritten by hand
  from its sources, in a commit that says it recomposed and why.

## Amendments to the method

None. The skill copy is byte-identical to its origin.
