# AGENTS.md

Guidance for AI coding agents working in this repository. Humans should start with
[README.md](README.md) (usage), [spec/arch.md](spec/arch.md) (design), and
[CONTRIBUTING.md](CONTRIBUTING.md) (workflow).

**The requirements are normative and live in [spec/](spec/).** This file summarizes the
rules that prevent damage and points at the rest; where it disagrees with `spec/`, `spec/`
wins. Requirement IDs below (`P-###`, `T-###`) are citations — search `spec/` for the full
statement and rationale.

## What this is

`birdsync` is a Go command-line tool that copies bird observations from an eBird CSV export
into a user's iNaturalist account, including photos and sounds from the Macaulay Library.
Single module, `github.com/Sajmani/birdsync`, one external dependency
(`github.com/google/uuid`).

## Commands

```
go build ./...        # build everything, including tools/
go vet ./...
gofmt -l .            # must print nothing
go test -race ./...   # full test suite; no network, no live accounts
```

Run all four before declaring work finished. There is no Makefile and no lint config beyond
`gofmt` and `go vet`.

## Layout

| Path | Contents |
| --- | --- |
| `birdsync.go` | Flags, `stats`, `main`, and `birdsync()` — the sync loop |
| `glue.go` | `ebirdClient` / `inatClient` interfaces + real impls; `dateTimeFlag` |
| `media.go` | `mlAssetSet`, eBird↔iNaturalist media diffing |
| `ebird/` | CSV parsing, Macaulay Library asset download |
| `inat/` | iNaturalist API v2 client, types, observation-field IDs, credentials |
| `tools/` | `dump`, a read-only observation dumper. Read-only by construction (T-032) |
| `spec/` | Requirements, acceptance criteria, decisions, architecture, process |

## Hard rules

Violating any of these costs a user real data. They are stated in full in
[spec/tech.md](spec/tech.md#safety-invariants).

- **Never make live API calls** (T-010). Tests must not contact `api.inaturalist.org` or the
  Macaulay Library. Use the `ebirdClient` / `inatClient` interfaces in `glue.go`, or the
  base-URL parameter on `inat.NewClient` and `ebird.downloadMLAsset`, to point at an
  `httptest` server (T-013, T-014).
- **Every mutating operation is gated on `--dryrun`** (T-005), at the call site in
  `birdsync()` — not in the client, which is shared with `tools/` (T-008). A `--dryrun` run
  issues no writes at all; reads are fine. Log the skipped action with a `DRYRUN:` prefix
  (T-006), and don't let the counters claim work that didn't happen (T-007). The gates are
  in `birdsync()`, around `UploadMedia`, `UpdateObservation`, and `CreateObservation`.
- **`tools/` stays read-only** (T-032), enforced by `TestToolsAreReadOnly`. Don't add a
  tool that creates, updates, deletes, or uploads. **Don't run one either** (T-033): they
  use the user's real credentials against the live service.
- **Never log or echo credentials** (T-012).
- **Don't add dependencies** (T-002) or touch the `go` / `toolchain` directives (T-003)
  without asking. Both are the maintainer's call, and both are proposed as their own change.

## Working on a change

Changes follow the loop in [spec/process.md](spec/process.md): the spec changes first, then
the acceptance criteria, then the code. Two things follow from that:

- **A behavior change edits `spec/product.md` in the same commit**, and the README too when
  it alters what gets written to iNaturalist or which observations are skipped. The README
  documents flag defaults and skip order, and has drifted from the code before.
- **Never resolve a conflict in code.** If two requirements disagree, stop and take it to
  [spec/decisions.md](spec/decisions.md). Picking one silently is the failure mode the
  process exists to prevent.

Five requirements are currently **not satisfied by the code** — the work arising from
Gate 1, listed in [spec/decisions.md](spec/decisions.md#work-arising). Don't "fix" the code
to match the old behavior when you notice the mismatch.

## Gotchas

Hard-won operational knowledge. The requirement each one protects is cited where there is
one.

- **Flags are package-level globals**, so tests share state. Call
  `resetFlags()` at the top of a new test. It covers every flag except `debug`, which each
  test saves and restores by hand (T-015).
- **`dateTimeFlag.Set("")` fails and leaves the old value.** Zero `after`/`before` by
  assignment (`after = dateTimeFlag{}`), not through `Set`. Existing tests that call
  `after.Set("")` are relying on the previous test's value (T-015).
- **eBird dates come in two formats**, `2006-01-02` and `1/2/2006`, with an optional time.
  Always compare via `ebird.Record.Observed()`, never `Record.Date` directly. Keying on the
  raw field was a real bug: it silently disabled `--fuzzy` for anyone whose export used the
  second format. Regression test: `TestFuzzyMatchDateFormats` (T-019).
- **The eBird CSV has a variable number of columns.** Read fields by header name, never by
  position (T-018).
- **A Macaulay Library asset ID doesn't reveal whether it's a photo or a sound.** The only
  way to know is to download it and see which URL responds. Code that needs the distinction
  before downloading — like `--dryrun` accounting — cannot have it (T-021, P-044).
- **`--verifiable` defaults to `true`**, so the default run skips observations without media
  (P-030).
- **`UpdateObservation` must keep `ignore_photos` set**, or updating a description detaches
  the observation's photos (T-009).
- **An empty taxon name must never enter the fuzzy-match index.** It matches every unnamed
  record on that date and silently drops legitimate observations (T-020, P-032).
- **iNaturalist API tokens expire every 24 hours**, so a token in the environment is likely
  stale. Don't diagnose a 401 as a code bug (P-017).

## Conventions

**American spellings** in prose, comments, and commit messages (T-037) — but **never edit a
quotation or a vendored document to match** (T-038). Quoted terms of service are verbatim
citations whose hashes are recorded; "fixing" a publisher's spelling falsifies the quote and
breaks the record. `TestAmericanSpellings` enforces both halves.


Stated as requirements in [spec/tech.md](spec/tech.md#code-conventions) (T-024 – T-029):
`gofmt`, no import aliasing, standard library grouped first; errors wrapped with the calling
function's name; `log.Fatal` in `main` and `tools/` only; `log.Printf` for progress and
`debugf` for detail; comments explain *why*, especially around the eBird and iNaturalist
quirks the code works around — preserve them.
