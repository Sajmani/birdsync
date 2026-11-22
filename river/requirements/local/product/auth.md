package product

# Product requirements: Authentication

Where credentials come from, and how they are handled. One item per heading, and each item's
first paragraph after its fields is the complete rule. Search `^## ` for the list, or
`## <name>` for one item. Other product requirements are in the other files in this directory,
all one package, `product`.

## credentialsFromEnvironment

**API token comes from INAT_API_TOKEN; INAT_USER_ID is optional**

Formerly: P-014

Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

The iNaturalist API token is read from `INAT_API_TOKEN`. If `INAT_USER_ID` is set, birdsync
uses it as the iNaturalist user ID; otherwise it reads the `user_id` claim directly from the
JWT API token without making an extra API request.

*Rationale: iNaturalist JWTs carry the account's numeric `user_id`, which the v2
`/observations` endpoint accepts in `user_id`. Extracting it locally removes the second
credential prompt while staying backward-compatible with explicit `INAT_USER_ID`.*

## promptForMissingCredentials

**Missing credentials are prompted for**

Formerly: P-015

Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

When `INAT_API_TOKEN` is unset, birdsync prompts for the API token interactively and explains
where to find it. It prompts for the iNaturalist user ID only when `INAT_USER_ID` is unset
and the token carries no `user_id` claim. On EOF or a read error at a prompt, credential
resolution returns an error rather than looping.

## tokenFraming

**Both bare tokens and JSON-framed tokens are accepted**

Formerly: P-016

Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

Both `INAT_API_TOKEN` and the interactive prompt accept either a bare token or the
`{"api_token":"TOKEN"}` JSON object served by `https://www.inaturalist.org/users/api_token`,
unwrapping the JSON framing before sending the token in the `Authorization` header.

*Rationale: users copy `{"api_token":"TOKEN"}` directly from the browser page or export the
bare token in shell scripts; accepting either form in both places removes a common
configuration mistake.*

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
