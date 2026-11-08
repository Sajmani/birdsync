package birdsync

# Observation download and iNaturalist API etiquette

Existing-observation downloads, `id_above` paging, User-Agent identification, and API request
pacing and daily budgeting for [`birdsync`](README.md).

## Included requirements

- [product.downloadExistingObservations](../../../local/product/download.md#downloadexistingobservations) — Existing observations are downloaded before syncing
- [product.dateWindowNarrowsDownload](../../../local/product/download.md#datewindownarrowsdownload) — A date window narrows the download
- [product.largeAccounts](../../../local/product/download.md#largeaccounts) — Works for accounts over 10,000 observations
- [product.noTaxonFilter](../../../local/product/download.md#notaxonfilter) — The download is not filtered by taxon
- [tech.errorsIncludeBody](../../../local/tech/etiquette.md#errorsincludebody) — API errors include the response body
- [inatapi.singleIPAndMediaCaps](../../../imported/inaturalist.org/inatapi/inatapi.md#singleipandmediacaps) — One IP; media download volume caps

## Composed requirements

## requestRate

**Requests to the iNaturalist API stay within iNaturalist's published rate**

Composed from: [inatapi.requestRate](../../../imported/inaturalist.org/inatapi/inatapi.md#requestrate),
[inatterms.rateAndVolume](../../../imported/inaturalist.org/inatterms/inatterms.md#rateandvolume),
[inatterms.dontBehaveLikeMachine](../../../imported/inaturalist.org/inatterms/inatterms.md#dontbehavelikemachine),
[tech.dailyRequestBudget](../../../local/tech/etiquette.md#dailyrequestbudget),
[decisions.dailyRequestBudget](../../../decisions.md#dailyrequestbudget)
Goals: [goals.communityNotBurdened](../../../../goals/goals.md#communitynotburdened), [goals.accountsStaySafe](../../../../goals/goals.md#accountsstaysafe)

birdsync sends requests to `api.inaturalist.org` at no more than about one per second, and no
more than 9,000 in a rolling 24-hour window (within iNaturalist's guideline of around 10,000 a
day), whatever the size of the export and in every mode, including `--dryrun`.

The terms (Mandatory) require posting at a rate that does not hinder other users and warn
against adding a lot of content very quickly; they give no number. The API's recommended
practices (Governing) give the numbers. Composition reads iNaturalist's own published limit as
the measure of the Mandatory rule, so neither the per-second nor the daily half can be given up
locally. `inat.Client` enforces the one-second spacing and a 9,000-request cap per run at the
single choke point every iNaturalist request passes through (`docs/arch.md`), and `birdsync()`
estimates prior 24-hour write requests from the `created_at` timestamps and media of the
downloaded observations so back-to-back runs stay within the same 9,000-request budget, stopping
cleanly at an observation boundary with instructions to re-run in 24 hours
([decisions.dailyRequestBudget](../../../decisions.md#dailyrequestbudget)).

## observationDownloadPageSize

**The existing-observation download asks for 200 observations per page**

Composed from: [inatapi.largestPageSize](../../../imported/inaturalist.org/inatapi/inatapi.md#largestpagesize)
Goals: [goals.communityNotBurdened](../../../../goals/goals.md#communitynotburdened)

Every page request in the download of the user's existing iNaturalist observations
([product.downloadExistingObservations](../../../local/product/download.md#downloadexistingobservations))
sets `per_page=200`, the largest page the API serves for observations.

The value is the source's own ("you can get up to 200 observations in a single request"). This
replaces the withdrawn local [tech.maxPageSize](../../../local/tech/etiquette.md#maxpagesize)
([decisions.restatementsComposed](../../../decisions.md#restatementscomposed)).

## observationDownloadByIDAbove

**The existing-observation download pages by id_above, in ascending id order**

Composed from: [inatapi.idAbovePastTenThousand](../../../imported/inaturalist.org/inatapi/inatapi.md#idabovepasttenthousand),
[decisions.pageByIDAbove](../../../decisions.md#pagebyidabove)
Goals: [goals.birdersKeepOneRecord](../../../../goals/goals.md#birderskeeponerecord), [goals.accountsStaySafe](../../../../goals/goals.md#accountsstaysafe)

The download of the user's existing observations
([product.downloadExistingObservations](../../../local/product/download.md#downloadexistingobservations))
sorts by id ascending (`order_by=id`, `order=asc`), sets `id_above` on each request after the
first to the highest observation ID in the page before, never pages by `page`, ends at the first
page shorter than the page size, and fails with an error if a page does not advance the cursor.

This is how birdsync meets [product.largeAccounts](../../../local/product/download.md#largeaccounts)
within the API's 10,000-result ceiling on page numbers, as decided in
[decisions.pageByIDAbove](../../../decisions.md#pagebyidabove). It replaces the withdrawn local
[tech.idAbovePaging](../../../local/tech/etiquette.md#idabovepaging) ([decisions.restatementsComposed](../../../decisions.md#restatementscomposed)).

## inatUserAgent

**Requests to the iNaturalist API identify birdsync in the User-Agent**

Composed from: [inatapi.identifyUserAgent](../../../imported/inaturalist.org/inatapi/inatapi.md#identifyuseragent),
[tech.userAgent](../../../local/tech/etiquette.md#useragent)
Goals: [goals.communityNotBurdened](../../../../goals/goals.md#communitynotburdened)

Every request birdsync sends to `api.inaturalist.org` carries a `User-Agent` header that names
birdsync and its version, rather than the Go HTTP client's default.

The source asks clients to "consider" a custom User-Agent; [`imports.md`](../../../imports.md)
adopts all eight `inatapi` items, so composition states it as a requirement over every request.
The version is the local addition of [tech.userAgent](../../../local/tech/etiquette.md#useragent). The source allows the app name, an iNaturalist
username, "or anything we might use to differentiate your requests"; the app name is the one
that holds for every user.

## Acceptance criteria

## userAgentSent

**API requests carry the birdsync User-Agent**

Formerly: AC-009
Covers: [inatUserAgent](#inatuseragent), [tech.userAgent](../../../local/tech/etiquette.md#useragent)

A request from `inat.Client` carries the birdsync `User-Agent`.

## downloadUsesMaxPageSize

**The download requests the maximum page size**

Formerly: AC-010
Covers: [observationDownloadPageSize](#observationdownloadpagesize), [product.downloadExistingObservations](../../../local/product/download.md#downloadexistingobservations), [tech.baseURLSeams](../../../local/tech/testing.md#baseurlseams)

`DownloadObservations` against an `httptest` server requests `per_page=200` and returns every observation served.

## downloadHonorsDateWindow

**The download honors the date window**

Formerly: AC-011
Covers: [product.dateWindowNarrowsDownload](../../../local/product/download.md#datewindownarrowsdownload)

With `--after` or `--before` set the download query carries the window, and without them it carries none.

## statusErrorCarriesBody

**API errors carry the response body**

Formerly: AC-032
Covers: [tech.errorsIncludeBody](../../../local/tech/etiquette.md#errorsincludebody)

A non-200 response's error includes the body, truncated and on one line, and drops an HTML body.

## downloadPagesByID

**The download pages by id_above past 10,000 results**

Formerly: AC-034
Covers: [product.largeAccounts](../../../local/product/download.md#largeaccounts), [observationDownloadByIDAbove](#observationdownloadbyidabove)

Against a fake honoring `per_page`, the download pages by `id_above` and errors when the cursor fails to advance.

## requestsPaced

**Requests are paced**

Formerly: AC-035
Covers: [requestRate](#requestrate)

Consecutive requests from one client are spaced by at least the minimum interval, and the default interval matches one per second.

Limit: covers per-second spacing; the 24-hour budget is covered by [pacingPerDay](#pacingperday).

## pacingPerDay

**No more than 9,000 iNaturalist requests in a rolling 24-hour window**

Covers: [requestRate](#requestrate)

`inat.Client` defaults to a 9,000-request cap and refuses requests past its configured cap
before sending them, and `birdsync()` estimates prior 24-hour write requests from the
downloaded observations' `created_at` timestamps and media (ignoring observations created more
than 24 hours ago), counts the download and write requests of the current run, stops cleanly at
an observation boundary before exceeding the 24-hour budget without leaving a partial
observation, and reports the stop in the end-of-run summary.

Limit: parallel processes and `--after`/`--before` runs that omit recently created observations
with older observation dates are accepted residuals
([decisions.dailyRequestBudget](../../../decisions.md#dailyrequestbudget)).

## downloadOrderedByID

**The download sorts by id ascending and never uses page**

Covers: [observationDownloadByIDAbove](#observationdownloadbyidabove)

Every request of the existing-observation download carries `order_by=id` and `order=asc`, and
none carries a `page` parameter.

## downloadStopsOnShortPage

**The download ends at the first short page**

Covers: [observationDownloadByIDAbove](#observationdownloadbyidabove)

When the test server answers a download request with fewer observations than the page size,
birdsync sends no further download request.
