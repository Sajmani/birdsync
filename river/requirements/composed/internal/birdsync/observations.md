package birdsync

# Observations, community obligations, and rights

What [`birdsync`](README.md) writes into each iNaturalist observation, what it documents about
community obligations, and contributor rights.

## Included requirements

- [product.createdObservationWild](../../../local/product/observations.md#createdobservationwild) — Created observations are wild
- [product.checklistCoordinates](../../../local/product/observations.md#checklistcoordinates) — Coordinates are the checklist's, inexact, with set accuracy
- [product.speciesGuessScientificName](../../../local/product/observations.md#speciesguessscientificname) — The species guess is the eBird scientific name
- [product.observedDateTime](../../../local/product/observations.md#observeddatetime) — Observation date and time come from the eBird columns
- [product.observationFieldMapping](../../../local/product/observations.md#observationfieldmapping) — eBird columns are copied into observation fields
- [product.descriptionContents](../../../local/product/observations.md#descriptioncontents) — What the observation description contains
- [product.observationsPublic](../../../local/product/observations.md#observationspublic) — Created observations are public
- [product.documentSyncedObservationsIdentifiable](../../../local/product/community.md#documentsyncedobservationsidentifiable) — Documentation says synced observations are identifiable
- [inatterms.contributorKeepsRights](../../../imported/inaturalist.org/inatterms/inatterms.md#contributorkeepsrights) — A contributor keeps their rights, and licenses by default as CC BY-NC

## Composed requirements

## accountHolderToldOfObligations

**The documentation tells users they answer for what birdsync posts under their account**

Composed from: [inatterms.accountHolderResponsible](../../../imported/inaturalist.org/inatterms/inatterms.md#accountholderresponsible),
[inatterms.dontBehaveLikeMachine](../../../imported/inaturalist.org/inatterms/inatterms.md#dontbehavelikemachine),
[decisions.notMachineGeneratedContent](../../../decisions.md#notmachinegeneratedcontent),
[product.documentCommunityObligations](../../../local/product/community.md#documentcommunityobligations)
Goals: [goals.accountsStaySafe](../../../../goals/goals.md#accountsstaysafe), [goals.communityNotBurdened](../../../../goals/goals.md#communitynotburdened)

birdsync's documentation tells users that iNaturalist holds them responsible for everything
birdsync posts under their account; that they should review the observations birdsync created;
and that an account which adds a lot of content quickly and does not respond to comments and
messages may be suspended, so they should answer comments and messages on synced observations.

birdsync cannot answer comments for the user; the part of
[inatterms.dontBehaveLikeMachine](../../../imported/inaturalist.org/inatterms/inatterms.md#dontbehavelikemachine)
it can reach is pacing ([requestRate](download.md#requestrate)), accuracy
([product.verifiableDefaultOn](../../../local/product/filtering.md#verifiabledefaulton)),
identifiability ([product.documentSyncedObservationsIdentifiable](../../../local/product/community.md#documentsyncedobservationsidentifiable)),
and telling the user. Reviewing what was created matters most: [inatterms.noMachineGeneratedContent](../../../imported/inaturalist.org/inatterms/inatterms.md#nomachinegeneratedcontent)
permits a script to post only where a human decided each observation, so a user who reviews
nothing moves birdsync to the prohibited side of that rule
([product.documentCommunityObligations](../../../local/product/community.md#documentcommunityobligations)).

## Acceptance criteria

## createdObservationIdentity

**A created observation carries the sync key and observation fields**

Formerly: AC-007 (part)
Covers: [product.syncKey](../../../local/product/identity.md#synckey), [product.observationFieldMapping](../../../local/product/observations.md#observationfieldmapping), [tech.clientInterfaces](../../../local/tech/testing.md#clientinterfaces)

A sync through a recording fake creates an observation whose sync-key fields and copied eBird observation fields have the required IDs and values.

## createdObservationPlacement

**A created observation is wild, placed, and timed as required**

Formerly: AC-007 (part)
Covers: [product.createdObservationWild](../../../local/product/observations.md#createdobservationwild), [product.checklistCoordinates](../../../local/product/observations.md#checklistcoordinates), [product.speciesGuessScientificName](../../../local/product/observations.md#speciesguessscientificname), [product.observedDateTime](../../../local/product/observations.md#observeddatetime)

The created observation is marked wild, carries the checklist's coordinates marked inexact with the configured accuracy, the scientific name as species guess, and the eBird date and time.

Limit: asserts what is sent, not how iNaturalist interprets it.

## readmeCommunitySection

**The README states the community obligations**

Formerly: AC-037
Covers: [accountHolderToldOfObligations](#accountholdertoldofobligations), [product.documentSyncedObservationsIdentifiable](../../../local/product/community.md#documentsyncedobservationsidentifiable)

A person reviewing the README finds what iNaturalist expects of a bulk-syncing account (that the user answers for what is posted, reviews what was created, and responds to comments and messages), and how synced observations are identifiable.

Limit: human review.
