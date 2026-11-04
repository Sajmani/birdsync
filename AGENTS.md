# AGENTS.md

Guidance for AI coding agents. Humans: start with [README.md](README.md),
[docs/arch.md](docs/arch.md), and [CONTRIBUTING.md](CONTRIBUTING.md).

**birdsync uses River.** Requirements, acceptance criteria, and decisions live in
[river/](river/). The method is the River skill in
[.agents/skills/river/](.agents/skills/river/SKILL.md): read it before changing what the
software must do. **Never edit `river/requirements/composed/` by hand**; change a source under
`river/requirements/local/` and recompose. **`river check` must pass before a change is done.**
Names below such as `tech.toolsReadOnly` are requirements: search `## toolsReadOnly` under
`river/requirements/`. Old IDs (`P-063`, `T-010`, `CR-007`) resolve through `Formerly:` fields.

## What this is

A Go command-line tool that copies bird observations from an eBird CSV export into a user's
iNaturalist account, with photos and sounds from the Macaulay Library. One module,
`github.com/Sajmani/birdsync`, one external dependency (`github.com/google/uuid`).

## Commands

```
go build ./...
go vet ./...
gofmt -l .            # must print nothing
go test -race ./...   # no network, no live accounts
river check           # from the repo root; build the CLI at the skill's pin
```

Run all five before declaring work finished.

## Layout

| Path | Contents |
| --- | --- |
| `birdsync.go` | Flags, `stats`, `main`, and `birdsync()`, the sync loop |
| `glue.go` | `ebirdClient` / `inatClient` interfaces and real impls; `dateTimeFlag` |
| `media.go` | `mlAssetSet`, eBird↔iNaturalist media diffing |
| `ebird/` | CSV parsing, Macaulay Library asset download |
| `inat/` | iNaturalist API v2 client, types, observation-field IDs, credentials |
| `tools/` | `dump`, a read-only observation dumper |
| `river/` | Requirements, criteria, decisions, conformance; `proposals.md` awaits the owner |

## Hard rules

Breaking these costs a user real data.

- **Never make live API calls** (`tech.noLiveServicesInTests`). Use the client interfaces in
  `glue.go`, or the base-URL parameter on `inat.NewClient` and `ebird.downloadMLAsset`, with
  an `httptest` server.
- **Gate every mutation on `--dryrun` at the call site in `birdsync()`**, not in the client,
  which `tools/` shares (`tech.dryRunGateAtCallSite`). A dry run writes nothing, logs each
  skipped action with `DRYRUN:`, and its counters never claim work that didn't happen.
- **`tools/` stays read-only, and never run a tool** (`tech.toolsReadOnly`,
  `tech.neverRunTools`): they use real credentials against the live service.
- **Never log or echo credentials** (`tech.credentialsNeverExposed`).
- **Don't add dependencies or touch the `go`/`toolchain` directives** without asking
  (`tech.oneExternalDependency`, `tech.goDirectivePolicy`).
- **You draft; the owner decides.** Never accept a goal, approve a composed package (an
  `Approves:` line), settle a conflict precedence cannot settle, or accept a gap. Draft it in
  `river/proposals.md`.

## Working on a change

- A **routine change** that every requirement already allows needs no River edits. If a
  guardrail test fails, the code is wrong.
- A **behavior change** edits a source in `river/requirements/local/` first, then recomposes,
  and updates README.md when it alters what gets written or skipped.
- **Never resolve a conflict in code.** Take it to the owner and record it in
  `river/requirements/decisions.md`.

## Gotchas

- **Flags are package-level globals.** Call `resetFlags()` at the top of a test; save and
  restore `debug` by hand. `dateTimeFlag.Set("")` fails and keeps the old value, so zero
  `after`/`before` by assignment.
- **eBird dates come in two formats.** Compare via `ebird.Record.Observed()`, never
  `Record.Date` (`tech.datesViaObserved`).
- **The CSV has a variable number of columns.** Read fields by header name.
- **An asset ID doesn't say photo or sound** until downloaded (`tech.assetTypeNeedsDownload`).
- **`--verifiable` defaults to true**, deliberately (`product.verifiableDefaultOn`).
- **`UpdateObservation` must keep `ignore_photos` set**, or photos detach.
- **An empty taxon name must never enter the fuzzy index** (`tech.emptyNameNotIndexed`).
- **API tokens expire every 24 hours**; a 401 is rarely a code bug.

## Conventions

American spellings in prose, comments, and commit messages, but **never edit a quotation, a
vendored document, or the skill copy** (`tech.quotationsNeverRespelled`).
`TestAmericanSpellings` enforces both. Code style: `tech.gofmtClean` through
`tech.commentsExplainWhy`. Preserve comments that explain eBird and iNaturalist quirks.
