# Contributing to birdsync

Thanks for your interest in improving birdsync. Bug reports and pull requests are welcome.

This file covers the mechanics: setup, the checks to run, and how to test safely. What
birdsync must do, and the constraints the code has to satisfy, are written down in
[spec/](spec/) — those are normative, and where this file disagrees with them, they win.

| If you want to know | Read |
| --- | --- |
| How to run the tool | [README.md](README.md) |
| How the code is put together | [spec/arch.md](spec/arch.md) |
| What it must do, and why | [spec/product.md](spec/product.md) |
| Constraints on the implementation | [spec/tech.md](spec/tech.md) |
| How a change travels from requirement to code | [spec/process.md](spec/process.md) |

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
go test -race ./...
```

CI runs the same four on every pull request, plus a non-blocking job against the latest
stable Go as an early warning.

If your change alters behavior, it also needs a matching edit to
[spec/product.md](spec/product.md), and to [README.md](README.md) when it changes what gets
written to iNaturalist or which observations are skipped. The README documents flag defaults
and skip order, and has drifted from the code before.

## Testing

**Tests must never make live API calls** (`spec/tech.md` T-010). They must not contact
`api.inaturalist.org` or the Macaulay Library, and they must not depend on anyone's real
observations. Two seams exist for this, and one of them will fit whatever you're testing:

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
account, and it only works if it's airtight. Anything that creates, updates, deletes, or
uploads sits behind the flag:

```go
if dryRun {
    log.Printf("DRYRUN: %s ...", whatWouldHappen)
} else {
    err = inatClient.CreateObservation(obs)
    // ...
}
```

The rules are `spec/tech.md` T-005 through T-008; in short, gate at the call site inside
`birdsync()` rather than in the client (`tools/` shares it), prefix the log line with
`DRYRUN:` because the README tells users to grep for it, and don't let the counters report
work that didn't happen. The gates are in `birdsync()`, around `UploadMedia`,
`UpdateObservation`, and `CreateObservation`.

This applies to birdsync itself. Nothing in `tools/` may mutate at all, so the question
doesn't arise there.

## Code style

Standard Go, written out as T-024 through T-029 in [spec/tech.md](spec/tech.md). The short
version: `gofmt` decides formatting; wrap errors with the calling function's name
(`fmt.Errorf("DownloadMLAsset(%s): %w", id, err)`); keep `log.Fatal` out of `ebird` and
`inat`; send progress to `log.Printf` and detail to `debugf`.

Comments should explain *why*. Much of this code works around undocumented quirks in the
eBird export format and the iNaturalist API, and those explanations are the most valuable
comments in the repository. Please don't strip them.

Please ask before adding a dependency — the single-dependency footprint is deliberate — and
raise the `go` directive in `go.mod` only as a deliberate, standalone change. Most users are
on the default `GOTOOLCHAIN=auto` and will download a newer toolchain without noticing, but
anyone pinned to `GOTOOLCHAIN=local` on an older Go will be blocked outright.

## Testing against a real account

The automated tests can't cover everything, and at some point you'll want to see a change work
end to end. Please be careful:

- Use `--dryrun` first. It prints what it would create and upload without touching anything.
- Narrow the run with `--after` and `--before` so you're working with a handful of
  observations rather than your whole life list.
- Consider a throwaway iNaturalist account rather than your own.
- Remember that observations you create during testing are public, and that the API token from
  <https://www.inaturalist.org/users/api_token> expires every 24 hours.

`tools/` holds only `dump`, which is read-only, and a check enforces that anything added
there stays that way (T-032 in [spec/tech.md](spec/tech.md)). Please don't add a tool that
writes: the ones that used to live there were guarded by a `debug` constant that was easy
to get wrong, and one of them shipped deleting for real.

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
