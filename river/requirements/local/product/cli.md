package product

# Product requirements: Invocation

How birdsync is invoked: arguments, flags, and checks before any network request. One item per
heading, and each item's first paragraph after its fields is the complete rule. Search `^## `
for the list, or `## <name>` for one item. Other product requirements are in the other files in
this directory, all one package, `product`.

## onePositionalArgument

**Exactly one positional argument**

Formerly: P-008
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

birdsync takes exactly one positional argument: the path to the eBird CSV
export. Any other number of arguments prints usage and exits non-zero.

## flagsBeforeArgument

**Flags precede the positional argument**

Formerly: P-009
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

Flags must appear before the positional argument.

*Rationale: consequence of the standard Go flag package; documented rather than fixed.*

## boolFlagEqualsForm

**Boolean flags are disabled with --flag=false**

Formerly: P-010
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

A boolean flag is disabled with `--flag=false`; `--flag false` is a usage
error.

*Rationale: same as [flagsBeforeArgument](#flagsbeforeargument).*

## unopenableCsvExitsEarly

**An unopenable CSV exits before contacting iNaturalist**

Formerly: P-011
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

If the CSV file cannot be opened, birdsync exits with an error before
contacting iNaturalist or prompting for credentials.

## afterLaterThanBeforeExits

**--after later than --before is an error**

Formerly: P-012
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

If `--after` is later than `--before`, birdsync exits with an error rather
than running a sync that cannot match anything.

## debugExplainsSkips

**--debug logs why each record was skipped**

Formerly: P-013
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

`--debug` logs the reason each eBird record was skipped, and other detail
not shown in a normal run.
