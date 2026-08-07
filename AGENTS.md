# AGENTS.md

Guidance for AI coding agents working in this repository. Humans should start with
[README.md](README.md) (usage), [arch.md](arch.md) (design), and
[CONTRIBUTING.md](CONTRIBUTING.md) (workflow).

## What this is

`birdsync` is a Go command-line tool that copies bird observations from an eBird CSV export
into a user's iNaturalist account, including photos and sounds from the Macaulay Library.
Single module, `github.com/Sajmani/birdsync`, one external dependency
(`github.com/google/uuid`).

## Commands

```
go build ./...        # build everything, including tools/
go test ./...         # full test suite; no network, no live accounts
go vet ./...
gofmt -l .            # must print nothing
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
| `tools/` | Six standalone `main` packages for account maintenance |

## Rules

- **Never make live API calls.** Tests must not contact `api.inaturalist.org` or the Macaulay
  Library. Use the `ebirdClient` / `inatClient` interfaces in `glue.go`, or the base-URL
  parameter on `inat.NewClient` and `ebird.downloadMLAsset`, to point at an `httptest` server.
- **Every mutating operation must be disabled by `--dryrun`.** Anything that creates, updates,
  deletes, or uploads goes behind `if dryRun { log what would happen } else { do it }`. A
  `--dryrun` run must issue no write requests at all; reads are fine. The gate belongs at the
  call site in `birdsync()` — the `inat.Client` methods mutate unconditionally and know nothing
  about the flag. Existing gates: birdsync.go:212 (`UploadMedia`), :236 (`UpdateObservation`),
  :350 (`CreateObservation`).
  - Log the skipped action with a `DRYRUN:` prefix; users grep for it.
  - Keep the counters honest. A dry run can't know whether an asset is a photo or a sound, so
    it counts into `stats.pendingMedia` rather than guessing. Don't report work as done that
    didn't happen.
- **Never run anything in `tools/`.** They create, update, and delete real observations in
  whoever's account the environment's credentials point at. `purge` is checked in with its
  `debug` guard off, so it deletes for real. Reading them is fine.
- **Don't add dependencies** without asking. The single-dependency footprint is deliberate.
- **Don't touch the `go` or `toolchain` directives as a side effect of other work.** Raising the
  `go` line is the maintainer's call: users on the default `GOTOOLCHAIN=auto` transparently
  download a newer toolchain, but anyone pinned to `GOTOOLCHAIN=local` on an older Go is
  blocked. Propose it as its own change.
- Changes that alter what gets written to iNaturalist, or which observations are skipped,
  need a corresponding README update. The README documents flag defaults and skip order, and
  has drifted from the code before.

## Gotchas

- **Flags are package-level globals** (`birdsync.go:23`), so tests share state. Call
  `resetFlags()` at the top of a new test. It covers every flag except `debug`, which each test
  saves and restores by hand.
- **`dateTimeFlag.Set("")` fails and leaves the old value.** Zero `after`/`before` by
  assignment (`after = dateTimeFlag{}`), not through `Set`. Existing tests that call
  `after.Set("")` are relying on the previous test's value.
- **eBird dates come in two formats**, `2006-01-02` and `1/2/2006`, with an optional time.
  Always compare via `ebird.Record.Observed()`, never `Record.Date` directly. Keying on the raw
  field was a real bug: it silently disabled `--fuzzy` for anyone whose export used the second
  format. Regression test at `birdsync_test.go:291`.
- **A Macaulay Library asset ID doesn't reveal whether it's a photo or a sound.** The only way
  to know is to download it and see which URL responds. Code that needs the distinction before
  downloading (like `--dryrun` accounting) cannot have it.
- **`--verifiable` defaults to `true`**, so the default run skips observations without media.
- **`UpdateObservation` must keep `ignore_photos` set**, or updating a description detaches the
  observation's photos.
- **An empty taxon name must never enter the fuzzy-match index.** It matches every unnamed
  record on that date and silently drops legitimate observations.
- **iNaturalist API tokens expire every 24 hours**, so a token in the environment is likely
  stale. Don't diagnose a 401 as a code bug.

## Conventions

- Standard `gofmt`; no aliasing of imports; standard library grouped first.
- Errors are wrapped with `fmt.Errorf("Caller: %w", err)` using the calling function's name.
  `inat.Client.roundTrip` is the exception, wrapping with a phrase describing the step
  (`"making HTTP request: %w"`).
- `log.Fatal` is acceptable in `main` and the tools, but not in `ebird` or `inat`, which return
  errors. (`ebird.Records` and `inat.DownloadObservations` violate this; don't copy them.)
- User-visible progress goes through `log.Printf`. Verbose detail goes through `debugf`, which
  is gated on `--debug`.
- Comments explain *why*, especially around the eBird and iNaturalist quirks the code works
  around. Preserve them.
