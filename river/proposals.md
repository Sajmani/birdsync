package proposals

# Proposals for the owner

**Not normative.** Everything here was drafted by an agent while birdsync was adopting River, and **none of it
has been decided yet**. Items the owner settled on 2026-11-04 were removed from this file once
recorded in `goals/`, `requirements/decisions.md`, `bindings.md`, or `precedence.md`. River's rule is that an agent drafts
and a person decides. The owner (Sajmani) accepts, changes, or rejects each item. Accepted
goals then go into `goals/`, accepted decisions into `requirements/decisions.md`, and each
item is deleted from here once settled. Search `^## ` for the list.

River keeps undecided drafts here, outside `goals/` and `decisions.md`. See
[`bindings.md`](bindings.md).

## githubActionsTerms

**Read GitHub's terms for Actions, and the licenses of `actions/checkout` and `actions/setup-go`**

Found by discovery and marked `partial` in `imports.md`; nothing was retrieved. Someone reads
the owner's pages and records their obligations, or the owner records why they impose nothing.
The owner adopted `uuid` on 2026-11-04 ([decisions.uuidAdopted](requirements/decisions.md#uuidadopted)).

## compositionQuestions

**Questions the recomposition of 2026-11-05 raised that precedence does not settle**

- **"About" one per second.** The criterion `requestsPaced` checks a minimum interval between
  requests. Is a minimum gap the threshold you want, or an average over a window?
- **OAuth.** [inatapi.noGroupAccounts](requirements/imported/inaturalist.org/inatapi/inatapi.md#nogroupaccounts) says to "use one of the various OAuth authentication
  flows". Does the pasted api_token JWT, which [inatapi.jwtAuthorization](requirements/imported/inaturalist.org/inatapi/inatapi.md#jwtauthorization) calls the preferred
  way, meet it, or should iNaturalist be asked? Composed as [birdsync.ownAccountOnly](requirements/composed/internal/birdsync/auth.md#ownaccountonly).
- **License on upload.** [inatterms.contributorKeepsRights](requirements/imported/inaturalist.org/inatterms/inatterms.md#contributorkeepsrights) makes uploads CC BY-NC unless the
  user says otherwise. Should birdsync leave the account default, set a license, or tell the user
  in the documentation?
- **User-Agent on the Macaulay Library CDN.** `inatapi` covers only iNaturalist. Should birdsync
  identify itself to the CDN too? Nothing requires it.

## approveComposedPackage

**Approve `composed/internal/birdsync`, once the items above are settled**

`river hash river/requirements/composed/internal/birdsync` prints the line to record. Only
the owner adds it to `decisions.md`. Recompose first if any decision above changes a source.

## discoveryComplete

**Add `Discovery: complete` to `imports.md` once `githubActionsTerms` is settled**

`river check` then fails, rather than notes, any unrecorded dependency or host.

## carriedOverOpenQuestions

**Open questions carried over from specman, so none drops out with `spec/`**

- Sound files over 50 MB are rejected by iNaturalist, and this surfaces as a generic upload failure. Should it be detected and reported specifically? (`spec/product.md`, open question 2)
- `inat.Client.UploadMedia` has no test. Its multipart request is exercised only through a mock (`spec/tech.md`, open question 1)
- The memory ceiling of [tech.memoryScalesWithExport](requirements/local/tech/resources.md#memoryscaleswithexport) has grown since [product.noTaxonFilter](requirements/local/product/download.md#notaxonfilter) and is unmeasured: worth a benchmark (`spec/tech.md`, open question 2)
- How much of `main` is worth testing: [product.onePositionalArgument](requirements/local/product/cli.md#onepositionalargument) to [product.debugExplainsSkips](requirements/local/product/cli.md#debugexplainsskips) live there untested (`spec/acceptance.md`, Gate 2 question 1)
- Whether to land the seven recommended additions in `spec/acceptance.md` now or as they become relevant: a token-in-logs check, error injection, testing `main` via `os/exec`, `UploadMedia` against `httptest`, a dependency count, and a large-export benchmark (Gate 2 question 2)
- Credentials-never-logged ([product.credentialsNeverLogged](requirements/local/product/auth.md#credentialsneverlogged), [tech.credentialsNeverExposed](requirements/local/tech/safety.md#credentialsneverexposed)) is a security-relevant invariant with no check at all. A test asserting the token never appears in captured log output is cheap (`spec/acceptance.md`, gaps worth naming 1)
- The create and update failure paths are unexercised: no test sets the mock's `createObsErr` or `updateObsErr` ([product.writeFailureAborts](requirements/local/product/failures.md#writefailureaborts); gaps worth naming 2)
- Order-dependent tests: nothing checks that each test calls `resetFlags()` ([tech.resetFlagsInTests](requirements/local/tech/testing.md#resetflagsintests); gaps worth naming 3)
- A successful sound upload has never been observed against the live service (`spec/decisions.md`, open work)
- A download failure is always transient, so a withdrawn asset is retried every run (`spec/decisions.md`, open work)
- Restoring an `Aves` plus `unknown` download filter is a pure bandwidth question, left open on purpose (`spec/decisions.md`)
- Issue #1 is fixed but left open until the reporter confirms on Windows (`spec/decisions.md`)
- [tech.honestDryRunCounters](requirements/local/tech/safety.md#honestdryruncounters) still carries specman's note "not yet satisfied", although `spec/tech.md` says it was implemented. The note was migrated verbatim; the owner may want it removed from the source
- Whether the Macaulay Library attribution guidelines should be vendored, so that `ML<id>` filenames and asset URLs can be checked against [mlterms.attributionExpected](requirements/imported/ebird.org/mlterms/mlterms.md#attributionexpected); `imports.md` keeps `mlterms` partial for this (owner, 2026-11-04: carry forward, don't block)
- Whether the 5 GB/hour and 24 GB/day caps in [inatapi.singleIPAndMediaCaps](requirements/imported/inaturalist.org/inatapi/inatapi.md#singleipandmediacaps), which name downloads, also limit uploads; a first sync can upload gigabytes (owner, 2026-11-04: carry forward, don't block)

The originals are in `spec/` at [`dc25e3b`](https://github.com/Sajmani/birdsync/tree/dc25e3b6756bd59c2de931e6d080d47318617da4/spec).
