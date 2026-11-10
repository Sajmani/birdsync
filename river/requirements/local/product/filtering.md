package product

# Product requirements: Filtering and skip order

Which eBird records are skipped, and in what order the tests apply. One item per heading, and
each item's first paragraph after its fields is the complete rule. Search `^## ` for the list,
or `## <name>` for one item. Other product requirements are in the other files in this
directory, all one package, `product`.

## skipOrder

**Records are tested in a fixed skip order**

Formerly: P-026
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

Each eBird record is tested in this order: `--after`, `--before`,
already-synced, `--fuzzy`, `--verifiable`. A record surviving all five becomes a new
iNaturalist observation.

*Rationale: the order is user-visible, because it determines which counter a skipped
record lands in.*

## afterFlag

**--after skips records observed before it**

Formerly: P-027
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

`--after` skips records observed strictly before the given time.

## beforeFlag

**--before skips records observed after it**

Formerly: P-028
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

`--before` skips records observed strictly after the given time.

## dateFlagFormats

**Date flags accept two formats**

Formerly: P-029
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

Both accept `2006-01-02 15:04:05` or `2006-01-02`.

## verifiableDefaultOn

**--verifiable skips medialess records, on by default**

Formerly: P-030
Subject: `cli.verifiable.default`
Value: `true`
Goals: [goals.communityNotBurdened](../../../goals/goals.md#communitynotburdened), [goals.communityGainsResearchGrade](../../../goals/goals.md#communitygainsresearchgrade)

`--verifiable` skips records with no Macaulay Library catalog numbers, and
defaults to **on**.

*Rationale: iNaturalist will not treat a medialess observation as verifiable, so the default
avoids filling a user's account with "Casual" records.*

*It is also the main thing birdsync does to limit the burden it puts on other people. A
medialess observation gives an identifier nothing to identify, so it is pure cost to the
community — the concern a forum moderator raised about volume
([decisions.notMachineGeneratedContent](../../decisions.md#notmachinegeneratedcontent)). Don't flip this
default for convenience.*

## fuzzyMatching

**--fuzzy skips records matching an unkeyed observation**

Formerly: P-031
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

`--fuzzy` skips a record when an existing observation not carrying a sync key
has the same observation date and the same name, matched against either the common name
or the scientific name.

## fuzzyNeverEmptyName

**Fuzzy matching never uses an empty name**

Formerly: P-032
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

Fuzzy matching never uses an empty name.

*Rationale: an empty name matches every unnamed record on that date and silently drops
legitimate observations.*

## fuzzyParsedDates

**Fuzzy matching compares parsed dates**

Formerly: P-033
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

Fuzzy matching compares parsed observation dates, not raw CSV date strings, across both eBird CSV
date formats (`2006-01-02` and `1/2/2006`).

*Rationale: keying on the raw field silently disabled `--fuzzy` for anyone whose export used the
second format.*

## fuzzyOffByDefault

**--fuzzy is off by default and its risk documented**

Formerly: P-034
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

`--fuzzy` is off by default, and its documentation states that it may
suppress a legitimate record when the same bird was seen twice in one day.
