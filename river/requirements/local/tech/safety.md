package tech

# Tech requirements: Safety invariants

Invariants that keep birdsync from changing data it should not. One item per heading, and each
item's first paragraph after its fields is the complete rule. Search `^## ` for the list, or
`## <name>` for one item. Other tech requirements are in the other files in this directory, all
one package, `tech`.

## dryRunGateAtCallSite

**Mutations are gated on --dryrun at the call site**

Formerly: T-005
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

Every operation that creates, updates, deletes, or uploads is gated on
`--dryrun` **at the call site in `birdsync()`**, not inside the client.

*Rationale: `inat.Client` is shared with `tools/`, which has no such flag. The gate
belongs where the decision is made, not in a client that other callers reuse.*

## dryRunPrefixLogged

**A skipped mutation is logged with DRYRUN:**

Formerly: T-006
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

A skipped mutation is logged with a `DRYRUN:` prefix (implements [product.dryRunLogPrefix](../product/dryrun.md#dryrunlogprefix)).

## honestDryRunCounters

**Dry-run counters never report work that did not happen**

Formerly: T-007
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

A dry run's counters must not report work that did not happen. Where a dry
run cannot know something, it counts it as unknown rather than guessing.

*Status: **not yet satisfied** — the created and updated counters are labeled as though
the work happened. [decisions.dryRunCountersRelabeled](../../decisions.md#dryruncountersrelabeled) resolved this
in favour of relabeling; implements [product.dryRunCounterLabels](../product/reporting.md#dryruncounterlabels).*

## clientMutatesUnconditionally

**inat.Client mutates unconditionally**

Formerly: T-008
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

`inat.Client`'s methods mutate unconditionally and know nothing about
`--dryrun`. This is deliberate, not an oversight.

## updateKeepsPhotos

**UpdateObservation always sets ignore_photos**

Formerly: T-009
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

`UpdateObservation` always sets `ignore_photos`.

*Rationale: without it, updating a description detaches the observation's photos.*

## neverRunMutatingTools

**Nothing runs the mutating tools**

Formerly: T-011
Status: withdrawn — the mutating tools were deleted rather than governed by a rule
Goals:

Nothing in the build, test, or development workflow runs a program in
`tools/`. They create, update, and delete observations in whatever account the
environment's credentials point at, and `purge` is checked in with its guard off.
nobody could check. Superseded by [toolsReadOnly](#toolsreadonly) and [neverRunTools](#neverruntools).

## toolsReadOnly

**No program in tools/ mutates**

Formerly: T-032
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

No program in `tools/` performs a mutating operation. Nothing there may
create, update, delete, or upload.

*Rationale: `tools/` held six programs that modified real accounts, guarded by a `debug`
constant that one of them shipped with turned off and by an instruction not to run them.
A structural property that a check can enforce is worth more than a warning people have to
read and remember. Deleted tools remain in the history if one is needed again.*

## neverRunTools

**Nothing in the workflow runs a program in tools/**

Formerly: T-033
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

Nothing in the build, test, or development workflow runs a program in
`tools/`. They contact the live service with the user's real credentials, which no
automated process should do on the user's behalf.

## credentialsNeverExposed

**Credentials are never logged, echoed, or written**

Formerly: T-012
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

Credentials are never logged, echoed, or written to a file.
