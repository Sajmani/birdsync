package product

# Product requirements: Scope

What birdsync is for, and what it deliberately does not do. One item per heading, and each
item's first paragraph after its fields is the complete rule. Search `^## ` for the list, or
`## <name>` for one item. Other product requirements are in the other files in this directory,
all one package, `product`.

The `product` package says what birdsync must do, from a user's point of view. `Formerly:`
gives the identifier an item had under the predecessor method (specman), which older
commits, comments, and issues still cite.

Migrated from `spec/product.md` at
[`dc25e3b`](https://github.com/Sajmani/birdsync/blob/dc25e3b6756bd59c2de931e6d080d47318617da4/spec/product.md),
where it was approved at specman's Gate 1 on 2026-08-09 and amended since. The
requirements were recovered from the implementation, README.md and arch.md, and describe
behavior as it was then. Goals were accepted on 2026-11-04 and are cited in each `Goals:` field.
Split into one file per area on 2026-11-04.

## copiesObservations

**Copies eBird observations, with their media, into iNaturalist**

Formerly: P-001
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

birdsync copies bird observations from a user's eBird data export into that
user's iNaturalist account, including the photos and sounds attached to them in the
Macaulay Library.

## oneWaySync

**The sync is one-way, eBird to iNaturalist**

Formerly: P-002
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

The sync is one-way: eBird → iNaturalist.

## csvExportInput

**The input is the user's eBird CSV export**

Formerly: P-003
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

The input is a `MyEBirdData.csv` file the user downloads from eBird, not a
live eBird API.

*Rationale: the eBird API does not expose personal checklists.*

## noReverseSync

**No sync from iNaturalist back to eBird**

Formerly: P-004
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

birdsync does not sync iNaturalist observations back to eBird.

## neverModifyOthersObservations

**Never modifies an observation birdsync did not create**

Formerly: P-005
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

birdsync never modifies or deletes an iNaturalist observation it did not
create, except to attach media to an observation carrying its own sync key.

## noTaxonomyReconciliation

**No taxonomy reconciliation between the services**

Formerly: P-006
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

birdsync does not reconcile taxonomy between the two services. An eBird
scientific name that iNaturalist cannot resolve produces an observation with an unknown
taxon, which the user fixes by hand.

## toolsNotProduct

**Programs in tools/ are not part of the product**

Formerly: P-007
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

The programs in `tools/` are not part of the product. They are maintenance
utilities, are not installed by `go install` of the root package, and carry no
compatibility promise. They are also read-only ([tech.toolsReadOnly](../tech/tech.md#toolsreadonly)): no tool may alter a user's
account.
