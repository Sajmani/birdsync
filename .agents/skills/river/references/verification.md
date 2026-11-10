# Acceptance criteria, verification, and guardrails

Part of River's method. Search `^## ` for the sections; read only the one you need.

## Acceptance criteria

An **acceptance criterion** is a prose statement of what must be true. **Composition produces
criteria**, so they live in the composed package beside the requirements they decide — they
belong to that version, and an obligation is only *demonstrably* preserved if whoever imports it
also receives how to decide it.

**A criterion cites what it covers**, in a `Covers:` field. A requirement never lists its
criteria and no traceability table is kept; a stored reverse index drifts from what it indexes,
and this one is derivable.

A criterion citing no requirement is an error, not an orphan. A guardrail may predate the
River module; a criterion cannot, because composition produces it from a requirement. Composed
packages hold only **in-force** criteria: when every requirement a criterion covers is withdrawn or
excluded, recomposition drops the criterion from `composed/` (`withdrawn-criterion` if marked
`Status: withdrawn` inside `composed/`, `covers-withdrawn` or `covers-excluded` if left covering a
withdrawn or excluded item), and any guardrail or `conformance.md` entry still citing it fails
(`guardrail-withdrawn`) so obsolete tests and dead code do not linger. Use `river graph
--if-withdrawn <id>` before withdrawing a goal or requirement to see every criterion and guardrail
that becomes unreachable.

**A criterion must have been watched failing before it is trusted.** A check that cannot fail is
indistinguishable from one that passes, and it is worse than nothing because it reports safety.
Likewise, a criterion's expected value comes from the requirement, never from a run of the
current code: a check built from existing output can only confirm that behavior has not
changed, including when it was wrong.

**Record each observation in `river/conformance.md`, in the criterion's section, in a fixed
field**, one line each:

```
Watched failing: 2026-10-16, by adding color:#777777 to an inline style; every species page failed
```

A criterion with no such line has not been watched failing. Prose may tell the story, but a
person or a tool deciding whether the criterion is trusted reads this field. **The observation
never goes in the criterion itself:** the criterion lives in a composed package, which is approved
by hash, so recording an observation there would force a recomposition and a new approval each
time a guardrail is watched failing ([*Recording conformance*](#recording-conformance)).

**A criterion is met or not as a whole.** Where parts of a claim can be met at different times (two
behaviors, two pages, an original check and a later extension), make each part its own criterion,
so each is watched failing, met and protected from regression on its own.

### When a requirement cannot be verified

Verification reports a requirement it cannot turn into a deciding criterion. It is one of three
things, and telling them apart is judgment:

1. **Underspecified.** Vague, not undecidable. *Responsive* becomes *p99 under 200 ms*.
2. **Exact but undecidable.** *Never leaks PII* is precise and cannot be established — instances
   are testable, absence is not. Attach the coverage you can get, state in the criterion what it
   does **not** cover, and count the requirement as uncovered.
3. **Requires judgment.** *Complies with the spirit of the regulation.* The guardrail is a human
   review with a named owner and a cadence. It is still a guardrail, just not automated.

**Never remove a requirement for being unverifiable**, and **never substitute an easier
criterion.** *Must be accessible*, checked by a contrast checker, reports green on a product
nobody can use. A requirement counts as covered only when its criterion **decides** it.

## Verification and guardrails

A **guardrail** is the thing that runs: a test or checker on an implementation, a monitoring or
rollout rule on a release. Criteria say what; guardrails enforce it.

A guardrail either **gates** — the change does not proceed — or **informs**, meaning the result
is recorded and weighed. Which one may differ by stage: an expensive check can inform at code
generation and gate at release. Gating is a property of the guardrail, not a separate kind of
thing.

**Guardrails are not River module artifacts.** They live wherever you already keep your tests and your
deployment configuration. The module **cites** them; it does not contain them, because requiring
otherwise would mean moving your tests as a condition of adoption. A scanner you already run is
a guardrail in this sense, and citing it is the whole of the integration.

The citation does not live in the composed package. A criterion is immutable with its version;
the guardrail implementing it is operational, and swapping one tool for another must not force a
new requirements version. Cite guardrails from a mutable, descriptive record.

**Check exhaustively what varies, and sample the instances.** A thousand pages rendered from one
template need every template checked, and every data-dependent input to the property checked in
full — not a thousand renders. Where an intermediary carries the requirement, three checks cover
it and none is per instance: the intermediary's **guarantee**, once; its **assumption**, for each
kind of user; and the **scope**, so the part it does not cover has its own criterion.

**Then check the result end to end anyway**, gating on a sample and informing on the full set.
The contract checks can all pass over a real violation when the contract is wrong — a stylesheet
check that reads `color` and never sees SVG text painted with `fill` — or when layers interfere,
as when a later stylesheet overrides a color or a middleware logs before the scrubber runs. The
end-to-end check tests the result regardless of how the work was divided, so the more layers
there are, the more it matters.

**A guardrail may have no criterion, and is never removed for lacking one.** A brownfield
project begins with a suite full of orphans. An orphan is one of three things and you cannot
tell which without looking: an undocumented requirement, a test of an implementation detail, or
a test of behavior nobody wants any more. The default is **keep and count**. Deleting
unexplained tests makes the codebase less safe in exchange for a tidier table.

## Recording conformance

**`river/conformance.md` is the mutable record of which guardrails implement which criteria, and
what has been seen.** It is descriptive: where it and the repository disagree, it is wrong. It
holds one section per criterion, headed by the criterion's qualified name, with fixed fields:

```markdown
## app.ServedContrast

Criterion: [app.ServedContrast](requirements/composed/internal/app/app.md#servedcontrast)
Guardrails: `page_test.go` TestServedPageContrast
Gates: CI on every change; release
Watched failing: 2026-10-22, by setting the panel text to #777777; the test failed on the panel
Watched failing: 2026-10-24, by adding an SVG label painted with fill; the test failed
```

- **`Criterion:`** links the criterion in its composed package.
- **`Guardrails:`** names each test, checker, monitor, or human review that implements it, where
  it lives, or `none`.
- **`Gates:`** says where each guardrail gates and where it only informs, or `informs only`.
- **`Watched failing:`** one line per observation: the date, how the behavior was broken, and
  what failed. None means not yet watched failing.

When a criterion's `Covers:` is repointed, as when a composed item replaces a local one, its
section here keeps its observations ([`composition.md`](composition.md), *Replacing a local
item*).
