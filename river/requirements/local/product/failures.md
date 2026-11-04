package product

# Product requirements: Failure behavior

How birdsync behaves when records, writes, or uploads fail. One item per heading, and each
item's first paragraph after its fields is the complete rule. Search `^## ` for the list, or
`## <name>` for one item. Other product requirements are in the other files in this directory,
all one package, `product`.

## malformedRecordAborts

**A malformed record aborts the run**

Formerly: P-057
Status: withdrawn by [decisions.badRecordsSkippedNotFatal](../../decisions.md#badrecordsskippednotfatal). Superseded by [badRecordsSkipped](#badrecordsskipped).
Goals:

A malformed date, time, latitude, or longitude in any record aborts the
whole run, identifying the CSV line.

## writeFailureAborts

**A failed create or update aborts the run**

Formerly: P-058
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

A failed observation create or update aborts the run; a failed media transfer
does not.

*Rationale: a create or update failure means a broken token or a broken service, where
continuing would produce one identical error per remaining record. Reaffirmed under
[decisions.badRecordsSkippedNotFatal](../../decisions.md#badrecordsskippednotfatal).*

## noRollback

**There is no rollback**

Formerly: P-059
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

There is no rollback. Observations created before an abort remain in the
user's account, and re-running is safe because of [idempotentRerun](identity.md#idempotentrerun).

## unusableInputExplained

**Unusable input is named before parsing**

Formerly: P-066
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

Before parsing, birdsync checks the input and names what is wrong, distinguishing
three cases: the path is a directory, the path is a zip archive, and the file has no
`Submission ID` column and so is not an eBird export.

*Rationale: [issue #1](https://github.com/Sajmani/birdsync/issues/1), where a user passed the
extracted download folder and got `Error reading CSV records from ...: Incorrect function.` —
an operating-system message that explains nothing.*

*The check is on the path and the header, never on the text of an operating-system error: that
symptom is what Windows returns for reading a directory, where Unix says "is a directory", so a
fix written against one would not work on the other ([tech.platformIndependent](../tech/deps.md#platformindependent)).*

*The header case is worse than an unhelpful message. Field lookup is `field[key]`, which yields
0 for an absent column, so a CSV without a `Submission ID` column did not fail — every record
silently took its submission ID from whatever the first column happened to be.*

## permanentUploadFailureRecorded

**A permanent upload failure is recorded and not retried**

Formerly: P-063
Subject: `media.upload.permanent_failure`
Value: `HTTP 4xx except 401, 408, 429`
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

An upload that fails permanently — the service rejected the file itself, rather
than the attempt — is recorded in the description as failed and is not attempted again. A
failure that may be transient is not recorded, and is retried on the next run.

*Rationale: [decisions.permanentFailuresRecorded](../../decisions.md#permanentfailuresrecorded). Without
the distinction, birdsync either gives up on a network blip or re-downloads a 50 MB sound
file on every run forever. Removing the line by hand is the way to ask for a retry.*

## permanentFailuresReported

**Permanently failed assets are reported each run**

Formerly: P-064
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

A permanently failed asset is reported on each run, so the user learns about it
without birdsync retrying it.

## badRecordsSkipped

**Unparseable records are skipped, logged, and counted**

Formerly: P-062
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

A record whose date, time, latitude, or longitude cannot be parsed is
skipped, logged with its CSV line, and counted. The run continues, and the count appears
in the summary when non-zero.

*Rationale: [decisions.badRecordsSkippedNotFatal](../../decisions.md#badrecordsskippednotfatal). This matches how
media failures are already handled ([mediaFailuresTolerated](media.md#mediafailurestolerated)); one bad row in a 12,000-row export should not
end a sync halfway through.*
