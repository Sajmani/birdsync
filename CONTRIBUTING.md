# Contributing to birdsync

Thanks for your interest in improving birdsync. Bug reports and pull requests are welcome.

## Getting set up

You need the [Go toolchain](https://go.dev/dl/). `go.mod` declares `go 1.23.0` as the language
version, along with `toolchain go1.24.5`, so by default the `go` command will fetch 1.24.5 if
your local Go is older.

```
git clone https://github.com/Sajmani/birdsync
cd birdsync
go build ./...
go test ./...
```

That's the whole setup. There is no Makefile, no code generation, and no linter beyond
`gofmt` and `go vet`. The project has exactly one external dependency,
`github.com/google/uuid`.

## Before you send a pull request

Run all four of these and make sure they're clean:

```
go build ./...
go vet ./...
gofmt -l .          # must print nothing
go test ./...
```

CI runs these same checks on every pull request, except that it runs the tests with the race
detector (`go test -race ./...`). Running that locally too will save you a round trip.

## Testing

**Tests must never make live API calls.** They must not contact `api.inaturalist.org` or the
Macaulay Library, and they must not depend on anyone's real observations. There are two
mechanisms for this, and one of them will fit whatever you're testing:

- **Fake clients.** `birdsync()` takes `ebirdClient` and `inatClient` interfaces (defined in
  `glue.go`). `birdsync_test.go` has `mockEBirdClient` and `mockINatClient` implementations.
  Use these to test sync-loop behavior — which observations get created, skipped, or updated.
- **Local HTTP servers.** `inat.NewClient` takes a base URL, and `ebird.DownloadMLAsset`
  delegates to an unexported `downloadMLAsset(baseURL, id)` that tests inside the `ebird`
  package call directly. Either way an `httptest.Server` can stand in for the real service.
  Use this to test request formatting and response parsing.

Two things to know about the test suite:

- The command-line flags are package-level variables, so tests share global state. Call
  `resetFlags()` at the start of any new test.
- `dateTimeFlag.Set("")` returns an error and leaves the previous value in place. Zero the
  `after` and `before` flags by assignment, not through `Set`.

If you're fixing a bug, please add a test that fails without your fix. It's worth confirming
that it does fail — revert your change, run the test, and check that it goes red before you
put it back.

## Every mutating operation must honor `--dryrun`

`--dryrun` is the safety net users rely on before letting this tool loose on their iNaturalist
account, and it only works if it's airtight. If you add anything that creates, updates,
deletes, or uploads, it must sit behind the flag:

```go
if dryRun {
    log.Printf("DRYRUN: %s ...", whatWouldHappen)
} else {
    err = inatClient.CreateObservation(obs)
    // ...
}
```

Three rules for getting this right:

- **Gate at the call site**, inside `birdsync()`. The `inat.Client` methods mutate
  unconditionally and don't know the flag exists; that's deliberate, since `tools/` uses the
  same client without it. The existing gates are at birdsync.go:212, :236, and :350.
- **Prefix the log line with `DRYRUN:`.** The README tells users to grep for it.
- **Don't let the counters lie.** A dry run must not report work it didn't do. When it can't
  know something — a Macaulay Library asset ID doesn't reveal whether it's a photo or a sound
  without downloading it — count it as unknown rather than guessing. That's what
  `stats.pendingMedia` is for.

A `--dryrun` run may issue reads, but it must issue no writes at all.

Note that the programs in `tools/` predate this and have no `--dryrun`; they use a `debug`
constant instead. If you add a tool, a real flag is preferable.

## Code style

Standard Go. `gofmt` decides formatting; beyond that:

- Wrap errors with the calling function's name: `fmt.Errorf("DownloadMLAsset(%s): %w", id, err)`.
- `log.Fatal` is fine in `main` and in `tools/`, but the `ebird` and `inat` packages should
  return errors to their caller instead. Two existing functions break this rule —
  `ebird.Records` and `inat.DownloadObservations` — so don't take them as the example to follow.
- User-facing progress goes through `log.Printf`; verbose detail goes through `debugf`, which
  only prints under `--debug`.
- Comments should explain *why*. Much of this code works around undocumented quirks in the
  eBird export format and the iNaturalist API, and those explanations are the most valuable
  comments in the repository. Please don't strip them.

Please ask before adding a dependency — the single-dependency footprint is deliberate.

Please also raise the `go` directive in `go.mod` only as a deliberate, standalone change rather
than as a side effect of something else. Most users are on the default `GOTOOLCHAIN=auto` and
will download a newer toolchain without noticing, but anyone pinned to `GOTOOLCHAIN=local` on
an older Go will be blocked outright.

If your change affects what gets written to iNaturalist, or which observations are skipped,
update [README.md](README.md) in the same pull request. It documents the flag defaults and the
order in which skip rules are applied, and it has drifted from the code before.

## Testing against a real account

The automated tests can't cover everything, and at some point you'll want to see a change work
end to end. Please be careful:

- Use `--dryrun` first. It prints what it would create and upload without touching anything.
- Narrow the run with `--after` and `--before` so you're working with a handful of
  observations rather than your whole life list.
- Consider a throwaway iNaturalist account rather than your own.
- Remember that observations you create during testing are public, and that the API token from
  <https://www.inaturalist.org/users/api_token> expires every 24 hours.

**The programs in `tools/` are destructive and have no dry-run flag.** `dedupe`, `purge`, and
`position` are gated by a `debug` constant at the top of the source file — `true` means log
only, `false` means do it. As checked in, `purge` deletes for real. Read the source before
running any of them, and don't point them at an account whose data you care about.

## Reporting bugs

Open an issue at <https://github.com/Sajmani/birdsync/issues>. Helpful details:

- The exact command line you ran.
- Relevant log output, ideally from a run with `--debug`.
- If it's a parsing problem, the header row of your `MyEBirdData.csv` and a sample data row
  with anything private removed. eBird's export varies between users — the number of columns
  differs, and dates appear as either `2006-01-02` or `1/2/2006` — so the raw shape of your
  file is often the thing that identifies the bug.

Please don't paste your iNaturalist API token into an issue.

## License

By contributing, you agree that your contributions will be licensed under the
[MIT License](LICENSE) that covers this project.
