package tech

# Tech requirements: External service etiquette

How birdsync behaves toward the iNaturalist API. One item per heading, and each item's first
paragraph after its fields is the complete rule. Search `^## ` for the list, or `## <name>` for
one item. Other tech requirements are in the other files in this directory, all one package,
`tech`.

## userAgent

**Every iNaturalist request carries a birdsync User-Agent**

Formerly: T-016
Subject: `http.user_agent`
Value: `birdsync/0.1`
Goals: [goals.communityNotBurdened](../../../goals/goals.md#communitynotburdened)

Every request to the iNaturalist API carries a `User-Agent` identifying
birdsync and its version.

## maxPageSize

**Downloads request the maximum page size**

Formerly: T-017
Subject: `inat.page_size`
Value: `200`
Goals: [goals.communityNotBurdened](../../../goals/goals.md#communitynotburdened)

Observation downloads request the maximum supported page size.

*Rationale: required by `inat/api-recommended-practices` (see
[imports.md](../../imports.md)).*

## errorsIncludeBody

**API errors include the response body**

Formerly: T-034
Goals: [goals.communityNotBurdened](../../../goals/goals.md#communitynotburdened)

An error from a non-200 API response includes the response body, truncated, not
just the status line.

*Rationale: iNaturalist explains its refusals in the body. Discarding it left every failure
looking identical — `bad HTTP status: 422 Unprocessable Entity` whether the file was too
large, the format unsupported, or the asset withdrawn. Classification under [product.permanentUploadFailureRecorded](../product/failures.md#permanentuploadfailurerecorded) depends on
the status code, but the human reading the log needs the reason.*

*An HTML body is dropped rather than included: it comes from a proxy rather than the API, and
a 413 arrives as a seven-line nginx page. A message that is kept is collapsed onto one line.*

## requestPacing

**iNaturalist requests are paced**

Formerly: T-035
Subject: `inat.request_rate`
Value: `1/second, 10000/day`
Goals: [goals.communityNotBurdened](../../../goals/goals.md#communitynotburdened)

Requests to the iNaturalist API are paced to the rate its recommended practices
ask for: about one per second, and about 10,000 per day.

*Rationale: [inatapi.requestRate](../../imported/inaturalist.org/inatapi/inatapi.md#requestrate), a governing source. Exceeding it returns HTTP 429, and the page
warns that persistent offenders may be IP-blocked. A first sync of a media-heavy account is
thousands of writes: reads are not the exposure.*

*Enforced in `Client.roundTrip`, which every request in the package passes through — the only
reason a single choke point suffices. The daily cap is not enforced: birdsync keeps no state
between runs, so it cannot know how many requests today has already seen.*

*This limit was already being met before the limiter existed, by an argument the maintainer
had reasoned through and not written down: observations are fetched in large pages, and
between writes birdsync downloads a media file, so the operations are slow enough to stay
under one per second on their own. Measurement agrees — a real run averaged about five seconds
per upload and three per download page, roughly 0.2 requests per second.*

*The argument is nonetheless worth replacing with a mechanism, for two reasons. It is
**conditional**: with `--verifiable=false` there is no media to fetch between creates, so the
spacing that made it true disappears. And it is **load-bearing on slowness** — the first
change that adds concurrency, or that makes uploads faster, silently removes the guarantee
with nothing to notice. An explicit limiter inverts that: compliance stops depending on the
code staying slow, which is what makes it safe to parallelize later.*

## idAbovePaging

**The download pages by id_above**

Formerly: T-036
Goals: [goals.communityNotBurdened](../../../goals/goals.md#communitynotburdened)

The observation download pages with `id_above` and an ascending id sort, never
with page numbers, and stops when a page comes back short.

*Rationale: [inatapi.idAbovePastTenThousand](../../imported/inaturalist.org/inatapi/inatapi.md#idabovepasttenthousand) and [decisions.pageByIDAbove](../../decisions.md#pagebyidabove).
`total_results` shrinks as the cursor advances, so it cannot decide when to stop. The client
requires the cursor to advance on every page: without that check, a server ignoring `id_above`
makes the loop re-fetch one page forever.*
