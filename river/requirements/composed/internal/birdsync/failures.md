package birdsync

# Failure behavior

Unusable input, malformed CSV rows, write aborts, and permanent versus transient upload failures
for [`birdsync`](README.md).

## Included requirements

- [product.writeFailureAborts](../../../local/product/failures.md#writefailureaborts) — A failed create or update aborts the run
- [product.noRollback](../../../local/product/failures.md#norollback) — There is no rollback
- [product.unusableInputExplained](../../../local/product/failures.md#unusableinputexplained) — Unusable input is named before parsing
- [product.permanentUploadFailureRecorded](../../../local/product/failures.md#permanentuploadfailurerecorded) — A permanent upload failure is recorded and not retried
- [product.permanentFailuresReported](../../../local/product/failures.md#permanentfailuresreported) — Permanently failed assets are reported each run
- [product.badRecordsSkipped](../../../local/product/failures.md#badrecordsskipped) — Unparseable records are skipped, logged, and counted

## Acceptance criteria

## invalidRecordsSkipped

**Unparseable records are skipped and counted**

Formerly: AC-025
Covers: [product.badRecordsSkipped](../../../local/product/failures.md#badrecordsskipped)

Records with an unparseable date, time, latitude, or longitude are skipped, logged with their line, and counted, and the run continues.

## permanentVersusTransientUpload

**Permanent and transient upload failures are told apart**

Formerly: AC-031
Covers: [product.permanentUploadFailureRecorded](../../../local/product/failures.md#permanentuploadfailurerecorded), [product.permanentFailuresReported](../../../local/product/failures.md#permanentfailuresreported)

A permanently refused upload is recorded and not retried, a transient one is retried, and the status-code classification holds in both directions.

## failedAssetsReported

**Permanently failed assets are reported each run**

Formerly: AC-033
Covers: [product.permanentFailuresReported](../../../local/product/failures.md#permanentfailuresreported)

An asset recorded as permanently failed is reported and not counted as attached or mismatched.

## badInputRejected

**Unusable input is named**

Formerly: AC-036
Covers: [product.unusableInputExplained](../../../local/product/failures.md#unusableinputexplained)

A directory, a zip archive, and a CSV without a `Submission ID` column are each rejected with a message naming the problem.
