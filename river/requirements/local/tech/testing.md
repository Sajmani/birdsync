package tech

# Tech requirements: Testability

Seams that make birdsync testable, and the rule that tests never touch a live service. One item
per heading, and each item's first paragraph after its fields is the complete rule. Search `^##
` for the list, or `## <name>` for one item. Other tech requirements are in the other files in
this directory, all one package, `tech`.

## clientInterfaces

**birdsync() takes its clients as interfaces**

Formerly: T-013
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

`birdsync()` takes its two clients as the `ebirdClient` and `inatClient`
interfaces, so the sync loop can be driven without a network.

## baseURLSeams

**Both services have base-URL test seams**

Formerly: T-014
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

`inat.NewClient` takes a base URL, and `ebird.DownloadMLAsset` delegates to
an unexported `downloadMLAsset(baseURL, id)`, so an `httptest` server can stand in for
either service. Neither seam may be removed.

## resetFlagsInTests

**Tests call resetFlags first**

Formerly: T-015
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

Flags are package-level variables, so tests share global state. A test calls
`resetFlags()` first.

*Gotcha: `dateTimeFlag.Set("")` returns an error and leaves the previous value in place,
so date flags are zeroed by assignment, not through `Set`.*

## noLiveServicesInTests

**Tests never contact a live service**

Formerly: T-010
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

Tests never contact `api.inaturalist.org`, the Macaulay Library CDN, or any
other live service, and never depend on a real account's contents.
