package tech

# Tech requirements: Module and dependencies

The module, its dependencies, and the platforms it builds on. One item per heading, and each
item's first paragraph after its fields is the complete rule. Search `^## ` for the list, or
`## <name>` for one item. Other tech requirements are in the other files in this directory, all
one package, `tech`.

The `tech` package holds the constraints the implementation must satisfy, independent of what
birdsync does for a user. `Formerly:` gives the identifier an item had under the predecessor
method (specman).

Migrated from `spec/tech.md` at
[`dc25e3b`](https://github.com/Sajmani/birdsync/blob/dc25e3b6756bd59c2de931e6d080d47318617da4/spec/tech.md),
where it was approved at specman's Gate 1 on 2026-08-09 and amended since. Goals were
accepted on 2026-11-04 and are cited in each `Goals:` field. Split into one file per area on
2026-11-04.

## singleModule

**A single Go module with no build system**

Formerly: T-001
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

birdsync is a single Go module, `github.com/Sajmani/birdsync`, with no build
system, no code generation, and no linter configuration beyond `gofmt` and `go vet`.

## oneExternalDependency

**Exactly one external dependency**

Formerly: T-002
Subject: `build.external_dependencies.count`
Value: `1`
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

The module has exactly one external dependency,
`github.com/google/uuid`. Adding another requires the owner's approval.

*Rationale: the small footprint is deliberate; it keeps `go install` fast and the
supply-chain surface near zero.*

## goDirectivePolicy

**go and toolchain directives change only deliberately**

Formerly: T-003
Subject: `build.go_directive`
Value: `1.23.0`
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

The `go` and `toolchain` directives are raised only as a deliberate,
standalone change, never as a side effect of other work.

*Rationale: users on the default `GOTOOLCHAIN=auto` silently download a newer toolchain,
but anyone pinned to `GOTOOLCHAIN=local` on an older Go is blocked outright.*

## platformIndependent

**Builds and runs on every Go platform**

Formerly: T-004
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

The tool must build and run on any platform the Go toolchain supports, using
only the standard library plus [oneExternalDependency](#oneexternaldependency). No platform-specific code.
