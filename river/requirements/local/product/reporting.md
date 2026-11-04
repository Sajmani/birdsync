package product

# Product requirements: Reporting

What a run reports when it ends. One item per heading, and each item's first paragraph after its
fields is the complete rule. Search `^## ` for the list, or `## <name>` for one item. Other
product requirements are in the other files in this directory, all one package, `product`.

## endOfRunSummary

**Every run ends with a summary**

Formerly: P-054
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

Every run ends with a summary: records processed, records skipped by each
rule, observations created and updated, and media uploaded.

## conditionalSkipCounters

**A skip counter prints only when its flag is in effect**

Formerly: P-055
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

A skip counter is printed only when the flag that produces it is in effect.

## conditionalFailureCount

**The media failure count prints only when non-zero**

Formerly: P-056
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

A media failure count is printed only when there were failures.

## dryRunCounterLabels

**Dry-run observation counters are labeled hypothetical**

Formerly: P-060
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

Under `--dryrun`, the observation counters are labeled as hypothetical:
`Would create N` and `Would update N`.

*Rationale: [decisions.dryRunCountersRelabeled](../../decisions.md#dryruncountersrelabeled). The count is
worth reporting; claiming the work happened is not.*
