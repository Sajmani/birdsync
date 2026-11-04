package product

# Product requirements: Dry run

What --dryrun does and reports. One item per heading, and each item's first paragraph after its
fields is the complete rule. Search `^## ` for the list, or `## <name>` for one item. Other
product requirements are in the other files in this directory, all one package, `product`.

## dryRunNoWrites

**--dryrun issues no writes**

Formerly: P-051
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

`--dryrun` issues no write request of any kind. Reads are permitted.

## dryRunLogPrefix

**Skipped actions are logged with DRYRUN:**

Formerly: P-052
Subject: `log.dryrun.prefix`
Value: `"DRYRUN: "`
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

Every action a dry run skips is logged with a `DRYRUN:` prefix.

*Rationale: the README tells users to grep for it.*

## dryRunMediaCount

**A dry run reports one unclassified media count**

Formerly: P-053
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

A dry run reports media as a single "would upload N assets" count rather than
splitting photos from sounds.

*Rationale: it doesn't download the assets, and the asset ID doesn't reveal the type
([photoThenSoundFallback](media.md#photothensoundfallback)).*
