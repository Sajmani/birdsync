package birdsync

# Identity, formats, and filtering

Sync key, idempotence, CSV and date formats, skip order, and `--fuzzy` matching for [`birdsync`](README.md).

## Included requirements

### Identity and idempotence

- [product.syncKey](../../../local/product/identity.md#synckey) — Own observations are identified by a two-field sync key
- [product.idempotentRerun](../../../local/product/identity.md#idempotentrerun) — Re-running creates no duplicates
- [product.taxonNotInSyncKey](../../../local/product/identity.md#taxonnotinsynckey) — The taxon is not part of the sync key
- [product.incompleteKeyNotRecognized](../../../local/product/identity.md#incompletekeynotrecognized) — An incomplete sync key is not recognized

### Filtering and skip order

- [product.skipOrder](../../../local/product/filtering.md#skiporder) — Records are tested in a fixed skip order
- [product.afterFlag](../../../local/product/filtering.md#afterflag) — --after skips records observed before it
- [product.beforeFlag](../../../local/product/filtering.md#beforeflag) — --before skips records observed after it
- [product.dateFlagFormats](../../../local/product/filtering.md#dateflagformats) — Date flags accept two formats
- [product.verifiableDefaultOn](../../../local/product/filtering.md#verifiabledefaulton) — --verifiable skips medialess records, on by default
- [product.fuzzyMatching](../../../local/product/filtering.md#fuzzymatching) — --fuzzy skips records matching an unkeyed observation
- [product.fuzzyNeverEmptyName](../../../local/product/filtering.md#fuzzyneveremptyname) — Fuzzy matching never uses an empty name
- [product.fuzzyParsedDates](../../../local/product/filtering.md#fuzzyparseddates) — Fuzzy matching compares parsed dates
- [product.fuzzyOffByDefault](../../../local/product/filtering.md#fuzzyoffbydefault) — --fuzzy is off by default and its risk documented

### Data format handling

- [tech.csvByHeaderName](../../../local/tech/formats.md#csvbyheadername) — The CSV is read by header name
- [tech.datesViaObserved](../../../local/tech/formats.md#datesviaobserved) — eBird dates are parsed and compared via Record.Observed()
- [tech.emptyNameNotIndexed](../../../local/tech/formats.md#emptynamenotindexed) — An empty taxon name never enters the fuzzy index

## Acceptance criteria

## skipOrderEnforced

**Each skip rule fires in order**

Formerly: AC-012
Covers: [product.idempotentRerun](../../../local/product/identity.md#idempotentrerun), [product.skipOrder](../../../local/product/filtering.md#skiporder), [product.afterFlag](../../../local/product/filtering.md#afterflag), [product.beforeFlag](../../../local/product/filtering.md#beforeflag), [product.verifiableDefaultOn](../../../local/product/filtering.md#verifiabledefaulton), [product.fuzzyMatching](../../../local/product/filtering.md#fuzzymatching)

One set of records exercises every skip rule exactly once, in the order `--after`, `--before`, already synced, `--fuzzy`, `--verifiable`, and the rest are created.

## fuzzyHandlesBothDateFormats

**Fuzzy matching works with both `2006-01-02` and `1/2/2006` eBird date formats**

Formerly: AC-014
Covers: [product.fuzzyParsedDates](../../../local/product/filtering.md#fuzzyparseddates), [tech.datesViaObserved](../../../local/tech/formats.md#datesviaobserved), [product.fuzzyMatching](../../../local/product/filtering.md#fuzzymatching)

Fuzzy matching skips a record whether eBird wrote its date as `2006-01-02` or as `1/2/2006`.

## fuzzyIgnoresEmptyNames

**Fuzzy matching ignores empty names**

Formerly: AC-015
Covers: [product.fuzzyNeverEmptyName](../../../local/product/filtering.md#fuzzyneveremptyname), [tech.emptyNameNotIndexed](../../../local/tech/formats.md#emptynamenotindexed)

An existing observation with an empty name never fuzzy-matches a record.

## fuzzyMatchesSubspeciesByBinomial

**Fuzzy matching matches subspecies and subspecific groups by binomial species stem**

Covers: [product.fuzzyMatching](../../../local/product/filtering.md#fuzzymatching)

Under `--fuzzy`, an unkeyed observation and an eBird record on the same date match when their
scientific names share the same binomial species stem (`Genus species`), including species-to-subspecies
refinements and eBird subspecific groups, while higher taxa, `sp.` entries, species slashes, and
hybrids do not match across species by stem.

## recordObservedParsing

**eBird dates parse in both `2006-01-02` and `1/2/2006` formats**

Formerly: AC-019
Covers: [tech.datesViaObserved](../../../local/tech/formats.md#datesviaobserved)

`Record.Observed()` parses both eBird CSV date formats — `2006-01-02` (with optional `03:04 PM`
time) and `1/2/2006` (with optional `3:04 PM` time) — with and without a time.

## csvReadByHeader

**The CSV is read by header name**

Formerly: AC-020
Covers: [tech.csvByHeaderName](../../../local/tech/formats.md#csvbyheadername), [product.csvExportInput](../../../local/product/scope.md#csvexportinput)

`ebird.Records` reads fields by header name from a file whose columns vary.

Limit: partial for the CSV input requirement.

## syncKeyValidity

**The sync key is valid only when complete**

Formerly: AC-021
Covers: [product.syncKey](../../../local/product/identity.md#synckey), [product.incompleteKeyNotRecognized](../../../local/product/identity.md#incompletekeynotrecognized), [product.taxonNotInSyncKey](../../../local/product/identity.md#taxonnotinsynckey)

`ObservationID.Valid` is true only when both sync-key fields are present.

Limit: partial for the taxon not being in the key.

## untaxonedObservationRecognized

**A previously synced untaxoned observation is recognized**

Formerly: AC-024
Covers: [product.noTaxonFilter](../../../local/product/download.md#notaxonfilter), [product.idempotentRerun](../../../local/product/identity.md#idempotentrerun)

The download carries no taxon filter, and an observation with no taxon that carries a sync key is recognized as synced and not re-created.
