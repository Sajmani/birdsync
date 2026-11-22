package birdsync

# Authentication and account ownership

Credentials, JWT authorization, account ownership, and authentication scope for [`birdsync`](README.md).

## Included requirements

- [product.credentialsFromEnvironment](../../../local/product/auth.md#credentialsfromenvironment) — API token comes from INAT_API_TOKEN; INAT_USER_ID is optional
- [product.promptForMissingCredentials](../../../local/product/auth.md#promptformissingcredentials) — Missing credentials are prompted for
- [product.tokenFraming](../../../local/product/auth.md#tokenframing) — Both bare tokens and JSON-framed tokens are accepted
- [product.unauthorizedSaysRefreshToken](../../../local/product/auth.md#unauthorizedsaysrefreshtoken) — A 401 says the token needs refreshing
- [product.credentialsNeverLogged](../../../local/product/auth.md#credentialsneverlogged) — Credentials are never written to the log
- [tech.credentialsNeverExposed](../../../local/tech/safety.md#credentialsneverexposed) — Credentials are never logged, echoed, or written

## Composed requirements

## jwtInAuthorizationHeader

**birdsync authenticates with the user's JWT in the Authorization header**

Composed from: [inatapi.jwtAuthorization](../../../imported/inaturalist.org/inatapi/inatapi.md#jwtauthorization)

Goals: [goals.birdersKeepOneRecord](../../../../goals/goals.md#birderskeeponerecord), [goals.accountsStaySafe](../../../../goals/goals.md#accountsstaysafe)

Every authenticated request birdsync sends to `api.inaturalist.org` carries, in its HTTP
`Authorization` header, the bare JWT the user obtained from `https://www.inaturalist.org/users/api_token`,
whether it came from `INAT_API_TOKEN` or from the interactive prompt and whether it was
supplied as a bare token or in its `{"api_token":"TOKEN"}` JSON framing
([product.credentialsFromEnvironment](../../../local/product/auth.md#credentialsfromenvironment),
[product.promptForMissingCredentials](../../../local/product/auth.md#promptformissingcredentials),
[product.tokenFraming](../../../local/product/auth.md#tokenframing)).

The local sources say where the token comes from, how its framing is unwrapped, and how a 401
from its 24-hour expiry is reported
([product.unauthorizedSaysRefreshToken](../../../local/product/auth.md#unauthorizedsaysrefreshtoken));
none says how it is sent, which is what the Governing source prescribes.

## ownAccountOnly

**birdsync posts only as the person running it, with their own token**

Composed from: [inatapi.noGroupAccounts](../../../imported/inaturalist.org/inatapi/inatapi.md#nogroupaccounts),
[inatterms.accountHolderResponsible](../../../imported/inaturalist.org/inatterms/inatterms.md#accountholderresponsible),
[decisions.inputIdentityIsTheAccountHolders](../../../decisions/decisions.md#inputidentityistheaccountholders)
Goals: [goals.accountsStaySafe](../../../../goals/goals.md#accountsstaysafe)

Every write birdsync makes to iNaturalist
([product.copiesObservations](../../../local/product/scope.md#copiesobservations)) is
authenticated with a token supplied at run time by the person running it
([product.credentialsFromEnvironment](../../../local/product/auth.md#credentialsfromenvironment),
[product.promptForMissingCredentials](../../../local/product/auth.md#promptformissingcredentials));
birdsync contains, ships, and falls back to no account or token of its own.

That keeps birdsync from being a group account: each user posts as themselves. Whether the
person running birdsync is the person whose eBird export it reads is the account holder's
responsibility, not birdsync's
([decisions.inputIdentityIsTheAccountHolders](../../../decisions/decisions.md#inputidentityistheaccountholders)).
The source also says to "use one of the various OAuth authentication flows"; whether a token
the user pastes from the api_token page meets that is not settled by precedence
(an open question in `river/proposals.md`).

## authenticateOnlyWhenNecessary

**birdsync sends its token to iNaturalist only on requests that need it**

Composed from: [inatapi.authenticateOnlyWhenNeeded](../../../imported/inaturalist.org/inatapi/inatapi.md#authenticateonlywhenneeded),
[inatterms.duplicatesUndesirable](../../../imported/inaturalist.org/inatterms/inatterms.md#duplicatesundesirable),
[decisions.authenticationIsNecessary](../../../decisions/decisions.md#authenticationisnecessary)
Goals: [goals.communityNotBurdened](../../../../goals/goals.md#communitynotburdened), [goals.accountsStaySafe](../../../../goals/goals.md#accountsstaysafe)

birdsync sends its token to `api.inaturalist.org` only on a write, or on a read whose
unauthenticated answer would omit observations birdsync must see to avoid creating duplicates;
every other request to iNaturalist is sent without it.

The owner ruled in [decisions.authenticationIsNecessary](../../../decisions/decisions.md#authenticationisnecessary)
that the download of existing observations
([product.downloadExistingObservations](../../../local/product/download.md#downloadexistingobservations))
is such a read: without the token it omits private and geoprivacy-obscured observations, each of
which the next run would create again ([product.idempotentRerun](../../../local/product/identity.md#idempotentrerun)).
So the rule is met today with every request authenticated, and needs no exception from iNaturalist.
Media downloads go to the Macaulay Library and are outside it.

## Acceptance criteria

## jwtFromEnvironmentInHeader

**A token from INAT_API_TOKEN is sent unwrapped in the Authorization header**

Covers: [jwtInAuthorizationHeader](#jwtinauthorizationheader), [product.tokenFraming](../../../local/product/auth.md#tokenframing)

With `INAT_API_TOKEN` set to either a bare token or `{"api_token":"TOKEN"}`, every
authenticated request the test server receives carries `TOKEN`, without JSON framing, in its
`Authorization` header.

## jwtFromPromptUnwrapped

**A token entered at the prompt is sent unwrapped in the Authorization header**

Covers: [jwtInAuthorizationHeader](#jwtinauthorizationheader), [product.tokenFraming](../../../local/product/auth.md#tokenframing), [product.promptForMissingCredentials](../../../local/product/auth.md#promptformissingcredentials)

With `INAT_API_TOKEN` unset and either `{"api_token":"TOKEN"}` or a bare `TOKEN` entered at
the prompt, every authenticated request the test server receives carries `TOKEN`, without
JSON framing, in its `Authorization` header.

## userIDFromJWT

**User ID is read from JWT claims when INAT_USER_ID is unset**

Covers: [product.credentialsFromEnvironment](../../../local/product/auth.md#credentialsfromenvironment), [product.promptForMissingCredentials](../../../local/product/auth.md#promptformissingcredentials)

With `INAT_USER_ID` unset and a JWT carrying a `user_id` claim, birdsync extracts that user
ID locally without reading stdin or contacting iNaturalist. With `INAT_USER_ID` set,
`INAT_USER_ID` takes precedence. Only when `INAT_USER_ID` is unset and the token carries no
`user_id` claim does birdsync prompt on stdin, returning an error on EOF.

## credentialsNotLeaked

**Credentials do not appear in logs or errors**

Covers: [product.credentialsNeverLogged](../../../local/product/auth.md#credentialsneverlogged), [tech.credentialsNeverExposed](../../../local/tech/safety.md#credentialsneverexposed)

Across valid and malformed token inputs, interactive prompt retries, and authenticated HTTP
requests (including 401 and error responses), the supplied API token never appears in
captured log output or returned errors. This checks the paths birdsync exercises; it cannot
prove absence across unreachable code paths.

## noBuiltInCredentials

**Without the user's token, birdsync writes nothing**

Covers: [ownAccountOnly](#ownaccountonly)

With `INAT_USER_ID` and `INAT_API_TOKEN` unset and the prompt given no token, the test server
receives no write request and no request carrying an `Authorization` header.

## writesCarrySuppliedToken

**Every write carries the token the user supplied**

Covers: [ownAccountOnly](#ownaccountonly)

Every write request the test server receives during a run carries, in its `Authorization`
header, exactly the token supplied for that run, and no other.

## tokenOnlyWhereNeeded

**Only writes and the existing-observation download carry the token**

Covers: [authenticateOnlyWhenNecessary](#authenticateonlywhennecessary)

Of the requests the test server receives during a run, every one carrying an `Authorization`
header is a write or a request of the existing-observation download; any other request carries
none. A new iNaturalist request added to birdsync falls under this check by being made, not by
being listed.

## authenticatedDownloadSeesPrivate

**The authenticated download sees private and obscured observations, so a rerun duplicates none**

Covers: [authenticateOnlyWhenNecessary](#authenticateonlywhennecessary)

Given a test server that returns the user's private and geoprivacy-obscured birdsync-created
observations only to an authenticated request, a second run over the same export creates no
observation.
