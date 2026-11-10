package goals

# Goals — birdsync

Why birdsync exists. Accepted by the owner, Sajmani, on 2026-11-04, from candidates drafted at
adoption. Every local requirement in `requirements/local/` and every composed requirement in
[`requirements/composed/internal/birdsync/`](../requirements/composed/internal/birdsync/README.md)
cites the goals it serves in its `Goals:` field; a goal never lists its requirements, and
included bullets in the composed package link to their source items without repeating `Goals:`.
`communityGainsResearchGrade` and `communityNotBurdened` pull against each other: more
observations are a gain to the community and a cost to its identifiers. Each requirement says
which it serves.

## birdersKeepOneRecord

**Birders keep one record, published in both eBird and iNaturalist**

Birders who keep their records in eBird get them, with their photos and sounds, into iNaturalist without entering anything twice.

## accountsStaySafe

**Users can trust birdsync with their iNaturalist account**

Users can run birdsync on their own iNaturalist account without losing data, duplicating observations, or putting the account's standing at risk.

## communityNotBurdened

**The iNaturalist community bears little cost from synced observations**

iNaturalist's identifiers and operators bear as little cost as possible from observations birdsync posts: synced records worth an identifier's time, and traffic that does not strain the service.

## communityGainsResearchGrade

**The iNaturalist community gains research-grade bird observations**

The iNaturalist community gains bird observations, with photos and sounds, that identifiers can confirm to research grade, and that would otherwise stay only in eBird.

## maintainerChangesSafely

**The maintainer can change birdsync quickly and safely**

The maintainer and contributors can change birdsync quickly without breaking a user's account.
