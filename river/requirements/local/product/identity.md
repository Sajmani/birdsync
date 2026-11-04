package product

# Product requirements: Identity and idempotence

How birdsync recognizes observations it created, so reruns create no duplicates. One item per
heading, and each item's first paragraph after its fields is the complete rule. Search `^## `
for the list, or `## <name>` for one item. Other product requirements are in the other files in
this directory, all one package, `product`.

## syncKey

**Own observations are identified by a two-field sync key**

Formerly: P-019
Subject: `sync.key.fields`
Value: `[6033, 20215]`
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

birdsync identifies its own observations by a sync key: the pair of
iNaturalist observation fields *eBird Checklist* (6033) and *eBird Scientific Name*
(20215).

## idempotentRerun

**Re-running creates no duplicates**

Formerly: P-020
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

Re-running birdsync over the same CSV creates no duplicate observations.

## taxonNotInSyncKey

**The taxon is not part of the sync key**

Formerly: P-021
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

The iNaturalist taxon is not part of the sync key.

*Rationale: the taxon may be changed by the user or the community after upload, and
eBird's "slash" (`Aythya marila/affinis`) and "spuh" (`Melanitta sp.`) names have no
exact iNaturalist equivalent.*

## incompleteKeyNotRecognized

**An incomplete sync key is not recognized**

Formerly: P-022
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

An observation missing either sync-key field is not recognized as
birdsync-created.

*Rationale: observations created by versions that set only the checklist field fall into
this population. A `repair` tool existed to backfill them and was deleted once the
maintainer's account was clean; it is recoverable from the history.*
