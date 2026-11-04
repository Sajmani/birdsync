package product

# Product requirements: Downloading existing observations

Downloading the user's existing iNaturalist observations before syncing. One item per heading,
and each item's first paragraph after its fields is the complete rule. Search `^## ` for the
list, or `## <name>` for one item. Other product requirements are in the other files in this
directory, all one package, `product`.

## downloadExistingObservations

**Existing observations are downloaded before syncing**

Formerly: P-023
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

Before syncing, birdsync downloads the user's existing iNaturalist
observations into memory.

## avesOnlyDownload

**The download is restricted to Aves**

Formerly: P-024
Status: withdrawn by [decisions.noTaxonFilterOnDownload](../../decisions.md#notaxonfilterondownload) — it defeated [idempotentRerun](identity.md#idempotentrerun). Superseded by [noTaxonFilter](#notaxonfilter).
Goals:

The download is restricted to the `Aves` iconic taxon.

## dateWindowNarrowsDownload

**A date window narrows the download**

Formerly: P-025
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

When `--after` or `--before` is set, the download is restricted to that date
window, so duplicate detection and fuzzy matching only consider observations inside it.

## noTaxonFilter

**The download is not filtered by taxon**

Formerly: P-061
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

The download of existing iNaturalist observations is not filtered by taxon.

*Rationale: [decisions.noTaxonFilterOnDownload](../../decisions.md#notaxonfilterondownload). An
observation whose eBird name iNaturalist could not resolve has no iconic taxon, so a
taxon-filtered query never returns it, so re-running duplicates it. Idempotence ([idempotentRerun](identity.md#idempotentrerun))
outranks download volume. A consequence is that `--fuzzy` now matches against
observations of any taxon.*

## largeAccounts

**Works for accounts over 10,000 observations**

Formerly: P-065
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

birdsync works for an account with more than 10,000 iNaturalist observations.

*Reported as [issue #5](https://github.com/Sajmani/birdsync/issues/5) in January 2026 and fixed
by [decisions.pageByIDAbove](../../decisions.md#pagebyidabove): the download
pages by observation id, which has no ceiling, rather than by page number, which the API caps
at 10,000 results.*
