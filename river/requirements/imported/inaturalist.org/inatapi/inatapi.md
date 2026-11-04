package inatapi

# inat-api — transcribed requirements

Tier: **governing**. Scope: every request birdsync makes to the iNaturalist API.

IDs are assigned by this transcription; the page has no numbering. Each entry quotes the
source so the reading can be checked. See [PROVENANCE.md](PROVENANCE.md).

Identifiers were renamed from the transcription's `R#` numbers to descriptive names when
birdsync adopted River; `Formerly:` gives the old one. The quotations are unchanged.

## requestRate

**About one request per second, about 10,000 a day**

Formerly: inat-api/R1

> "Please keep requests to about 1 per second, and around 10k API requests a day"
>
> "Requests exceeding this limit might be throttled, and will return an HTTP 429 exception
> “Too Many Requests”... Please add delays into your code to keep under these limits...
> We may block IPs that consistently exceed these limits"

**Implemented as `requestPacing`.** `Client.roundTrip` waits out a minimum interval. The daily cap is
not enforced: birdsync keeps no state between runs and cannot know what today has spent.

## largestPageSize

**Use the largest supported page size**

Formerly: inat-api/R2

> "If using the API to fetch a lot of results, please use the highest supported per_page
> value. For example you can get up to 200 observations in a single request, which would be
> faster and more efficient than fetching the default 30 results at a time"

**Implemented as `maxPageSize`** (`per_page=200`).

## idAbovePastTenThousand

**Page numbers cap at 10,000 results; use id_above beyond that**

Formerly: inat-api/R3

> "The page and per_page parameters can be used to fetch up to (for many endpoints) 10k
> results. An error will be thrown if results beyond 10k are requested."
>
> "One way to use the API and fetch more than 10k records is to sort by id ascending (e.g.
> &order_by=id&order=asc) and use the id_above parameter set to the ID of the record in the
> last batch."

**Implemented as `idAbovePaging`**, after [issue #5](https://github.com/Sajmani/birdsync/issues/5) and
the decision `pageByIDAbove`.

## identifyUserAgent

**Identify the client with a User-Agent**

Formerly: inat-api/R4

> "please consider using a custom User Agent to identify your application... this could be
> set to the name of your app, an iNaturalist username, or anything we might use to
> differentiate your requests"

**Implemented as `userAgent`** (`birdsync/0.1`).

## jwtAuthorization

**JWTs, expiring in 24 hours, in an Authorization header**

Formerly: inat-api/R5

> "The preferred way to make an authenticated request in the newer API is to use a JSON Web
> Token (JWT), available at https://www.inaturalist.org/users/api_token... These tokens expire
> in 24 hours... Include the JWT in an HTTP Authorization header"

**Matches current behavior** and explains `unauthorizedSaysRefreshToken`: a 401 usually means the token aged out, not
that anything is broken.

## noGroupAccounts

**Users log in as themselves; no group accounts**

Formerly: inat-api/R6

> "If you are creating an application that will allow users to post data to iNaturalist,
> please do not do it using a group account (i.e. many individuals submitting data under a
> single username). Use one of the various OAuth authentication flows to allow users to log
> in to iNaturalist as themselves"

**Satisfied by construction.** birdsync runs on the user's machine with the user's own token
and posts under their own account. There is no birdsync service account.

## authenticateOnlyWhenNeeded

**Authenticate only when necessary**

Formerly: inat-api/R7

> "Authenticate API requests only when required. Since the data returned from authenticated
> requests might contain private information, these requests are not cached, therefore they
> can require more server time to process"

**Not satisfied, deliberately.** birdsync sends the token on every request, including the
observation download, which costs iNaturalist uncached work. Reading the user's own
observations unauthenticated would miss anything private or geoprivacy-obscured, and missing
an existing observation means creating a duplicate — the failure mode of
the decision `noTaxonFilterOnDownload`.
Recorded as a knowing departure rather than an oversight; revisit if iNaturalist objects.

## singleIPAndMediaCaps

**One IP; media download volume caps**

Formerly: inat-api/R8

> "Please use a single IP address for fetching data. If we think multiple IPs are being used
> in coordination to bypass rate limits, we may block those IPs regardless of query rate"
>
> "Downloading over 5 GB of media per hour or 24 GB of media per day may result in a
> permanent block"

**Not applicable as written.** birdsync runs on one machine, and downloads media from the
Macaulay Library rather than from iNaturalist. Whether the volume cap also governs *uploads*
is not stated; birdsync can upload gigabytes on a first sync, so the question is real.
Recorded rather than assumed away.

---

## Not transcribed

Bulk access routes — observation exports, the weekly GBIF dataset — are for consumers of
iNaturalist data. birdsync writes a user's own records and reads only their own account.
