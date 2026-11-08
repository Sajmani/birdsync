package birdsync

# Dry run and end-of-run reporting

`--dryrun` call-site gates, logging, and end-of-run summary reporting for [`birdsync`](README.md).

## Included requirements

### Dry-run safety

- [product.dryRunNoWrites](../../../local/product/dryrun.md#dryrunnowrites) — --dryrun issues no writes
- [product.dryRunLogPrefix](../../../local/product/dryrun.md#dryrunlogprefix) — Skipped actions are logged with DRYRUN:
- [product.dryRunMediaCount](../../../local/product/dryrun.md#dryrunmediacount) — A dry run reports one unclassified media count
- [tech.dryRunGateAtCallSite](../../../local/tech/safety.md#dryrungateatcallsite) — Mutations are gated on --dryrun at the call site
- [tech.dryRunPrefixLogged](../../../local/tech/safety.md#dryrunprefixlogged) — A skipped mutation is logged with DRYRUN:
- [tech.honestDryRunCounters](../../../local/tech/safety.md#honestdryruncounters) — Dry-run counters never report work that did not happen
- [tech.clientMutatesUnconditionally](../../../local/tech/safety.md#clientmutatesunconditionally) — inat.Client mutates unconditionally

### Summary reporting

- [product.endOfRunSummary](../../../local/product/reporting.md#endofrunsummary) — Every run ends with a summary
- [product.conditionalSkipCounters](../../../local/product/reporting.md#conditionalskipcounters) — A skip counter prints only when its flag is in effect
- [product.conditionalFailureCount](../../../local/product/reporting.md#conditionalfailurecount) — The media failure count prints only when non-zero
- [product.dryRunCounterLabels](../../../local/product/reporting.md#dryruncounterlabels) — Dry-run observation counters are labeled hypothetical

## Acceptance criteria

## dryRunIssuesNoWrites

**A dry run issues no writes**

Formerly: AC-006
Covers: [product.dryRunNoWrites](../../../local/product/dryrun.md#dryrunnowrites), [tech.dryRunGateAtCallSite](../../../local/tech/safety.md#dryrungateatcallsite), [tech.clientInterfaces](../../../local/tech/testing.md#clientinterfaces)

A `--dryrun` sync through a recording fake makes no create, update, or upload call, while reads still happen.

## dryRunMediaCountUnclassified

**A dry run reports one unclassified media count**

Formerly: AC-016
Covers: [product.dryRunMediaCount](../../../local/product/dryrun.md#dryrunmediacount)

Under `--dryrun` the media count is a single unclassified number.

Limit: does not detect a broken dry-run gate; [dryRunIssuesNoWrites](#dryrunissuesnowrites) does.

## summaryLabels

**The summary labels and conditional lines are right**

Formerly: AC-023
Covers: [product.dryRunCounterLabels](../../../local/product/reporting.md#dryruncounterlabels), [tech.honestDryRunCounters](../../../local/tech/safety.md#honestdryruncounters), [product.endOfRunSummary](../../../local/product/reporting.md#endofrunsummary), [product.conditionalSkipCounters](../../../local/product/reporting.md#conditionalskipcounters), [product.conditionalFailureCount](../../../local/product/reporting.md#conditionalfailurecount)

The summary says `Would create` and `Would update` under `--dryrun`, and the real labels otherwise. A counter line appears only when its flag is in effect, and the failure line only when non-zero.
