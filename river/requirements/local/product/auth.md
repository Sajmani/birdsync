package product

# Product requirements: Authentication

Where credentials come from, and how they are handled. One item per heading, and each item's
first paragraph after its fields is the complete rule. Search `^## ` for the list, or
`## <name>` for one item. Other product requirements are in the other files in this directory,
all one package, `product`.

## credentialsFromEnvironment

**Credentials come from INAT_USER_ID and INAT_API_TOKEN**

Formerly: P-014
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

The iNaturalist user ID is read from `INAT_USER_ID`, and the API token from
`INAT_API_TOKEN`.

## promptForMissingCredentials

**Missing credentials are prompted for**

Formerly: P-015
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

When either is unset, birdsync prompts for it interactively and explains
where to find it.

## tokenFraming

**Token framing differs between prompt and environment**

Formerly: P-016
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

At the prompt the token is accepted in its `{"api_token":"TOKEN"}` framing;
from the environment it is the bare token.

*Rationale: matches what the iNaturalist page gives the user, versus what is convenient
to export in a shell.*

## unauthorizedSaysRefreshToken

**A 401 says the token needs refreshing**

Formerly: P-017
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

A 401 from the API reports that the token needs refreshing, and where to get
a new one, rather than a bare HTTP status.

*Rationale: tokens expire every 24 hours, so this is the most likely failure a user
hits.*

## credentialsNeverLogged

**Credentials are never written to the log**

Formerly: P-018
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

Credentials are never written to the log.
