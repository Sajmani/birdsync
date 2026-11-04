package tech

# Tech requirements: Continuous integration

What CI runs. One item per heading, and each item's first paragraph after its fields is the
complete rule. Search `^## ` for the list, or `## <name>` for one item. Other tech requirements
are in the other files in this directory, all one package, `tech`.

## ciStandingChecks

**CI runs the standing checks**

Formerly: T-030
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

CI runs build, vet, gofmt, and the test suite under the race detector, on
every push to `main` and every pull request, using the toolchain declared in `go.mod`.

*The actions it depends on are kept current: GitHub retires the Node runtime older versions
target, and a deprecation warning today is a broken build later. `actions/checkout` and
`actions/setup-go` moved to Node 24 at v5 and v6 respectively; both are pinned at v7.*

## ciLatestStable

**CI also runs latest stable Go, non-blocking**

Formerly: T-031
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

CI additionally runs against the latest stable Go as an early warning, and
that job may fail without blocking a pull request.
