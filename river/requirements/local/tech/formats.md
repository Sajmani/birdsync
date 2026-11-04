package tech

# Tech requirements: Data format handling

How birdsync reads eBird and Macaulay Library data. One item per heading, and each item's first
paragraph after its fields is the complete rule. Search `^## ` for the list, or `## <name>` for
one item. Other tech requirements are in the other files in this directory, all one package,
`tech`.

## csvByHeaderName

**The CSV is read by header name**

Formerly: T-018
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

The eBird CSV is read by header name, never by column position, with
`FieldsPerRecord = -1`.

*Rationale: the export has a variable number of columns between users.*

## datesViaObserved

**eBird dates are compared via Record.Observed()**

Formerly: T-019
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

Any comparison of an eBird observation date goes through `Record.Observed()`,
never `Record.Date` directly.

*Rationale: eBird writes `2006-01-02` or `1/2/2006`, with an optional time. Keying on the
raw field was a real bug; regression test: `TestFuzzyMatchDateFormats`.*

## emptyNameNotIndexed

**An empty taxon name never enters the fuzzy index**

Formerly: T-020
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

An empty taxon name never enters the fuzzy-match index (implements [product.fuzzyNeverEmptyName](../product/filtering.md#fuzzyneveremptyname)).

## assetTypeNeedsDownload

**Asset type is known only after download**

Formerly: T-021
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

Code that needs to know whether a Macaulay Library asset is a photo or a
sound must download it first; the ID does not encode it.
