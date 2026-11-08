package birdsync

# birdsync: composed package

**Composed, not authored.** Change it by changing a source and recomposing, never by editing
files in this directory by hand. Composed from every source in `river/requirements/local/` and
`river/requirements/imported/`, under [`precedence.md`](../../../precedence.md) and
[`decisions.md`](../../../decisions.md). **Not yet approved**: `decisions.md` carries no
`Approves:` line for this package.

Each area file holds that area's included source requirements, any rules composition derived for
that area, and the acceptance criteria covering them. Which guardrail implements each criterion,
and whether it has been watched failing, is recorded in [`conformance.md`](../../../../conformance.md).

| File | Area |
| --- | --- |
| [`scope.md`](scope.md) | Purpose, non-goals, CLI invocation, `tools/`, and posting only export records |
| [`auth.md`](auth.md) | Credentials, JWT authorization, account ownership, and authenticating only where needed |
| [`download.md`](download.md) | Existing-observation download, `id_above` paging, User-Agent, and API request pacing and daily budget |
| [`filtering.md`](filtering.md) | Sync key, idempotence, CSV and date formats, skip order, and `--fuzzy` matching |
| [`observations.md`](observations.md) | Created observation fields, public visibility, community obligations, and contributor rights |
| [`media.md`](media.md) | Macaulay Library asset downloads, uploads, attribution, and temporary-file cleanup |
| [`failures.md`](failures.md) | Unusable input, malformed CSV rows, write aborts, and permanent versus transient upload failures |
| [`dryrun.md`](dryrun.md) | `--dryrun` call-site gates, logging, and end-of-run summary reporting |
| [`engineering.md`](engineering.md) | Go module and dependencies, testability seams, code and spelling conventions, and CI |

## Composed from

- inatterms (Mandatory)
- mlterms (Mandatory)
- inatapi (Governing)
- tech (Local: tech)
- product (Local: product)
- decisions
- precedence

## Precedence applied

Applied in the rank order of [`precedence.md`](../../../precedence.md): Mandatory (`inatterms`,
`mlterms`), Governing (`inatapi`), Local: tech (`tech`), Local: product (`product`), and Advisory
(`go-practice`, absorbed into `tech`, plus `uuid` and `github-actions` in
[`imports.md`](../../../imports.md)).

## Excluded items

- [product.avesOnlyDownload](../../../local/product/download.md#avesonlydownload) *(excluded)* — The download is restricted to Aves. Excluded by [decisions.noTaxonFilterOnDownload](../../../decisions.md#notaxonfilterondownload)
- [product.malformedRecordAborts](../../../local/product/failures.md#malformedrecordaborts) *(excluded)* — A malformed record aborts the run. Excluded by [decisions.badRecordsSkippedNotFatal](../../../decisions.md#badrecordsskippednotfatal)
- [tech.neverRunMutatingTools](../../../local/tech/safety.md#neverrunmutatingtools) *(excluded)* — Nothing runs the mutating tools. Excluded by [decisions.mutatingToolsDeleted](../../../decisions.md#mutatingtoolsdeleted)
- [tech.maxPageSize](../../../local/tech/etiquette.md#maxpagesize) *(excluded)* — Downloads request the maximum page size. Excluded by [decisions.restatementsComposed](../../../decisions.md#restatementscomposed)
- [tech.requestPacing](../../../local/tech/etiquette.md#requestpacing) *(excluded)* — iNaturalist requests are paced. Excluded by [decisions.restatementsComposed](../../../decisions.md#restatementscomposed)
- [tech.idAbovePaging](../../../local/tech/etiquette.md#idabovepaging) *(excluded)* — The download pages by id_above. Excluded by [decisions.restatementsComposed](../../../decisions.md#restatementscomposed)
