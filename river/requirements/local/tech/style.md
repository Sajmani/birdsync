package tech

# Tech requirements: Code conventions

How the code and its prose are written. One item per heading, and each item's first paragraph
after its fields is the complete rule. Search `^## ` for the list, or `## <name>` for one item.
Other tech requirements are in the other files in this directory, all one package, `tech`.

## gofmtClean

**gofmt -l . prints nothing**

Formerly: T-024
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

`gofmt -l .` prints nothing.

## vetClean

**go vet is clean**

Formerly: T-025
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

`go vet ./...` is clean.

## errorWrapping

**Errors are wrapped with the calling function's name**

Formerly: T-026
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

Errors are wrapped with the calling function's name:
`fmt.Errorf("DownloadMLAsset(%s): %w", id, err)`. `inat.Client.roundTrip` is the
exception, wrapping with a phrase describing the step.

## noLogFatalInLibraries

**log.Fatal only in main and tools/**

Formerly: T-027
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

`log.Fatal` is acceptable in `main` and in `tools/`; the `ebird` and `inat`
packages return errors to their caller.

*Enforced by static analysis rather than convention, since a test can show that one
function returns an error but only analysis can show that no function exits. `os.Exit` is
checked too, being the same thing under another name.*

## progressLogging

**Progress via log.Printf, detail via debugf**

Formerly: T-028
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

User-visible progress goes through `log.Printf`; verbose detail goes through
`debugf`, gated on `--debug`.

## americanSpelling

**Prose uses American spelling**

Formerly: T-037
Subject: `docs.spelling`
Value: `en-US`
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

Prose in this repository — documentation, comments, commit messages — uses
American spellings.

*Rationale: consistency, and the maintainer's preference. It is the weakest requirement here
on its own; it is written down for the exception below.*

## quotationsNeverRespelled

**Spelling changes never touch quotations or vendored documents**

Formerly: T-038
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

A spelling change never touches a quotation or a vendored document.

*Rationale: this is the part with consequences. A bulk find-and-replace across `river/` would
silently rewrite text quoted from someone else's terms of service, turning a verbatim citation
into a paraphrase that still looks verbatim, and would break the hashes recorded in each
source's `PROVENANCE.md`. The rule is not "be careful"; the check skips blockquote lines and
never opens `river/requirements/imported/**` or the skill copy in `.agents/skills/` (under
specman these were `spec/` and `spec/sources/**`).*

## commentsExplainWhy

**Comments explain why**

Formerly: T-029
Goals: [goals.maintainerChangesSafely](../../../goals/goals.md#maintainerchangessafely)

Comments explain *why*. The explanations of eBird and iNaturalist quirks are
the most valuable comments in the repository and are preserved through refactors.
