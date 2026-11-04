package product

# Product requirements: Community obligations

What birdsync's documentation must tell users about iNaturalist's community. One item per
heading, and each item's first paragraph after its fields is the complete rule. Search `^## `
for the list, or `## <name>` for one item. Other product requirements are in the other files in
this directory, all one package, `product`.

## documentCommunityObligations

**Documentation states the user's community obligations**

Formerly: P-067
Goals: [goals.communityNotBurdened](../../../goals/goals.md#communitynotburdened)

The documentation tells users what iNaturalist expects of an account that syncs
in bulk: review what was created, and respond to identifiers' comments.

*Rationale: [inatterms.dontBehaveLikeMachine](../../imported/inaturalist.org/inatterms/inatterms.md#dontbehavelikemachine) makes this the account holder's obligation, with suspension named
as the consequence, and birdsync is what puts them in a position to breach it. It also carries
more weight than politeness: [inatterms.noMachineGeneratedContent](../../imported/inaturalist.org/inatterms/inatterms.md#nomachinegeneratedcontent) permits a script to post observations only where a
human decided each one, so a user who reviews nothing is the one who moves birdsync from the
permitted side of that rule to the prohibited side. Recorded from
[decisions.notMachineGeneratedContent](../../decisions.md#notmachinegeneratedcontent).*

## documentSyncedObservationsIdentifiable

**Documentation says synced observations are identifiable**

Formerly: P-068
Goals: [goals.communityNotBurdened](../../../goals/goals.md#communitynotburdened), [goals.communityGainsResearchGrade](../../../goals/goals.md#communitygainsresearchgrade)

The documentation explains that birdsync-created observations are identifiable —
by the attribution line and the two eBird observation fields — so that identifiers who would
rather not work on synced records can recognize them.

*Rationale: raised on the forum thread in [decisions.notMachineGeneratedContent](../../decisions.md#notmachinegeneratedcontent) as a mitigation for identifier burden. It is
already true ([syncKey](identity.md#synckey), [descriptionContents](observations.md#descriptioncontents)); it is simply not written down anywhere a reader would find it.*
