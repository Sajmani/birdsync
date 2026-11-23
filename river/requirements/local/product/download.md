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

Status: withdrawn by [decisions.noTaxonFilterOnDownload](../../decisions/decisions.md#notaxonfilterondownload) — it defeated [idempotentRerun](identity.md#idempotentrerun). Superseded by [noTaxonFilter](#notaxonfilter).

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

*Rationale: [decisions.noTaxonFilterOnDownload](../../decisions/decisions.md#notaxonfilterondownload). An
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
by [decisions.pageByIDAbove](../../decisions/decisions.md#pagebyidabove): the download
pages by observation id, which has no ceiling, rather than by page number, which the API caps
at 10,000 results.*

## incrementalStateCache

**Routine syncs use a local state cache and updated_since; --full forces a complete download**

Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord), [goals.communityNotBurdened](../../../goals/goals.md#communitynotburdened), [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

By default birdsync keeps a per-account local state cache at
`$XDG_STATE_HOME/birdsync/state.json` (falling back to `os.UserCacheDir()/birdsync/state.json`,
overridable via `--state`, and disabled when `--state=""`). After an initial full download,
subsequent unwindowed runs query iNaturalist with `updated_since` (using a 1-minute overlap
before the cached watermark) and merge changed observations into the cache. `--full` ignores
the watermark and replaces the cached observations with a complete download, reconciling
observations deleted on iNaturalist. Date-windowed runs (`--after` or `--before`) never
advance the full-history watermark, `--dryrun` never creates or modifies the state file, a
corrupt or unreadable cache falls back to a full download, and no credentials are stored in
the state file.

*Rationale: [decisions.localStateCache](../../decisions/decisions.md#localstatecache). Avoids
re-downloading thousands of unchanged observations on every run while preserving exact
duplicate detection, `--fuzzy` matching, and `--dryrun` safety.*
