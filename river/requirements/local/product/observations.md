package product

# Product requirements: What birdsync writes

The contents of each observation birdsync creates. One item per heading, and each item's first
paragraph after its fields is the complete rule. Search `^## ` for the list, or `## <name>` for
one item. Other product requirements are in the other files in this directory, all one package,
`product`.

## createdObservationWild

**Created observations are wild**

Formerly: P-035
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord), [goals.communityGainsResearchGrade](../../../goals/goals.md#communitygainsresearchgrade)

A created observation is marked wild, not captive.

*Rationale: eBird checklists record wild birds.*

## checklistCoordinates

**Coordinates are the checklist's, inexact, with set accuracy**

Formerly: P-036
Subject: `observation.positional_accuracy.default_m`
Value: `1000`
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord), [goals.communityGainsResearchGrade](../../../goals/goals.md#communitygainsresearchgrade)

Latitude and longitude are the checklist's, the location is marked
not-exact, and positional accuracy is `--positional_accuracy_meters`, default 1000.

*Rationale: the checklist location approximates a hotspot, not the bird.*

## speciesGuessScientificName

**The species guess is the eBird scientific name**

Formerly: P-037
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord), [goals.communityGainsResearchGrade](../../../goals/goals.md#communitygainsresearchgrade)

The species guess is the eBird scientific name.

## observedDateTime

**Observation date and time come from the eBird columns**

Formerly: P-038
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord), [goals.communityGainsResearchGrade](../../../goals/goals.md#communitygainsresearchgrade)

The observation date and time come from the eBird `Date` and `Time` columns.

## observationFieldMapping

**eBird columns are copied into observation fields**

Formerly: P-039
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

These eBird columns are copied into iNaturalist observation fields: Count
(1), Common Name (256), Location (157), County (245), State/Province (7739), Number of
Observers (2527), Submission ID (6033), Scientific Name (20215).

## descriptionContents

**What the observation description contains**

Formerly: P-040
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

The description contains: a line attributing the observation to birdsync, the
eBird observation details when present, the checklist URL, the protocol, the checklist
comments when present, one `Macaulay Library Asset:` line per uploaded asset, and one
`Macaulay Library Asset (upload failed permanently; delete this line from the description to
retry):` line per asset the service permanently refused ([permanentUploadFailureRecorded](failures.md#permanentuploadfailurerecorded)).

*The note spells out the remedy because the person who finds the line is reading an
observation, not the source, and the retry mechanism is otherwise undiscoverable.*

## observationsPublic

**Created observations are public**

Formerly: P-041
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe), [goals.communityGainsResearchGrade](../../../goals/goals.md#communitygainsresearchgrade)

Observations birdsync creates are public, like any iNaturalist observation.
