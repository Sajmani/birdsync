package birdsync

# Scope and CLI invocation

Purpose, non-goals, command-line invocation, and `tools/` boundaries for [`birdsync`](README.md).

## Included requirements

### Purpose and non-goals

- [product.copiesObservations](../../../local/product/scope.md#copiesobservations) — Copies eBird observations, with their media, into iNaturalist
- [product.oneWaySync](../../../local/product/scope.md#onewaysync) — The sync is one-way, eBird to iNaturalist
- [product.csvExportInput](../../../local/product/scope.md#csvexportinput) — The input is the user's eBird CSV export
- [product.noReverseSync](../../../local/product/scope.md#noreversesync) — No sync from iNaturalist back to eBird
- [product.neverModifyOthersObservations](../../../local/product/scope.md#nevermodifyothersobservations) — Never modifies an observation birdsync did not create
- [product.noTaxonomyReconciliation](../../../local/product/scope.md#notaxonomyreconciliation) — No taxonomy reconciliation between the services
- [product.toolsNotProduct](../../../local/product/scope.md#toolsnotproduct) — Programs in tools/ are not part of the product
- [tech.toolsReadOnly](../../../local/tech/safety.md#toolsreadonly) — No program in tools/ mutates
- [tech.neverRunTools](../../../local/tech/safety.md#neverruntools) — Nothing in the workflow runs a program in tools/

### Command-line invocation

- [product.onePositionalArgument](../../../local/product/cli.md#onepositionalargument) — Exactly one positional argument
- [product.flagsBeforeArgument](../../../local/product/cli.md#flagsbeforeargument) — Flags precede the positional argument
- [product.boolFlagEqualsForm](../../../local/product/cli.md#boolflagequalsform) — Boolean flags are disabled with --flag=false
- [product.unopenableCsvExitsEarly](../../../local/product/cli.md#unopenablecsvexitsearly) — An unopenable CSV exits before contacting iNaturalist
- [product.afterLaterThanBeforeExits](../../../local/product/cli.md#afterlaterthanbeforeexits) — --after later than --before is an error
- [product.debugExplainsSkips](../../../local/product/cli.md#debugexplainsskips) — --debug logs why each record was skipped

## Composed requirements

## postsOnlyUsersOwnRecords

**Every observation birdsync posts is a record from the user's own export, and it posts nothing else of its own**

Composed from: [inatterms.noMachineGeneratedContent](../../../imported/inaturalist.org/inatterms/inatterms.md#nomachinegeneratedcontent),
[decisions.notMachineGeneratedContent](../../../decisions.md#notmachinegeneratedcontent)
Goals: [goals.accountsStaySafe](../../../../goals/goals.md#accountsstaysafe), [goals.communityNotBurdened](../../../../goals/goals.md#communitynotburdened)

Every observation birdsync creates
([product.copiesObservations](../../../local/product/scope.md#copiesobservations),
[product.syncKey](../../../local/product/identity.md#synckey)) corresponds to one record in the
eBird export the user supplied ([product.csvExportInput](../../../local/product/scope.md#csvexportinput)),
and birdsync posts no identification, comment, or message of its own.

[decisions.notMachineGeneratedContent](../../../decisions.md#notmachinegeneratedcontent) holds
that birdsync is not machine-generated content because a human decided each observation and the
machine only facilitates posting it, from the user's own curated data. That holds only while
each observation is one the user recorded, and while birdsync generates no other kind of
content; the Mandatory rule, whose penalty is suspension without warning, is met by keeping
both true.

## Acceptance criteria

## toolsAreReadOnly

**tools/ is read-only**

Formerly: AC-028
Covers: [tech.toolsReadOnly](../../../local/tech/safety.md#toolsreadonly), [product.toolsNotProduct](../../../local/product/scope.md#toolsnotproduct)

No program under `tools/` calls a mutating `inat.Client` method.

## createsMatchExportRecords

**Every created observation corresponds to one export record**

Covers: [postsOnlyUsersOwnRecords](#postsonlyusersownrecords)

Over a run, every observation create the test server receives carries a sync key
([product.syncKey](../../../local/product/identity.md#synckey)) that matches one record of the
input export, and no two creates match the same record.

## noIdentificationsCommentsOrMessages

**birdsync posts no identification, comment, or message**

Covers: [postsOnlyUsersOwnRecords](#postsonlyusersownrecords)

Over a run that creates, uploads, and updates, every write request the test server receives
targets an observation or its media; none targets identifications, comments, or messages.
