# Implementation, release, optimization, and revision

Part of River's method. Search `^## ` for the sections; read only the one you need.

## Implementation

Implementation produces the code and the **designs that describe it**, kept in step at every
commit. A design may lead the code: written and reviewed before the code exists, it says it is a
proposal, and it becomes the code's description in the commit where the code lands. From then on
a change that makes the two disagree updates both, and a guardrail can detect drift between them.
A design is never a record written afterwards. This follows the Elephant-Goldfish Model, which
reviews a design with zero-context agents before code is written and reports where they drift.

**A design lives with the code, not in the River module.** The module holds intent and conformance;
keeping a design in it would put it where people edit requirements rather than where they edit
code, which is backwards for a document whose only job is to stay in step. `conformance.md` cites
where the designs are.

## Release and Experimentation

Release produces a release version from an implementation version, verifies it against the
release guardrails, and ramps traffic according to rollout policy.

**Experimentation** aligns releases with goals using feature flags, A/B tests, impact analyses
and traffic shifts. It hill-climbs gradients toward their targets and rejects candidates that
hit a backstop.

## Optimization

**Optimization reduces the cost, time and risk of future changes** to Requirements,
Implementation, Verification and Release. It is continuous and belongs to no single stage. It
acts on four things: **structure** (refactoring, debt, reorganizing an overlarge source),
**tooling** at any stage, the **tech stack**, and **agent capability** — the skills and
instruction files an agent reads.

**Optimization is distinguished from Implementation by purpose, not by what it touches.** Both
change code. Implementation changes the system to satisfy a requirement; Optimization changes it
to make future satisfaction cheaper. Porting to a faster language is Optimization when it makes
later performance work tractable, and Implementation when it meets a stated performance
requirement. Say which: only one of them answers to an acceptance criterion.

**It acts on the River module too.** Restructuring requirements so that writing code and guardrails
from them costs less — a composed package that lists rather than copies, a requirement whose
rule comes first — is Optimization: nothing required changes.

**It has no target.** Cost, time and risk are gradients with a direction and no finish line.

**It depends on Verification more than any other stage.** Restructuring is safe only because
guardrails hold behavior still while the structure moves. A project with no guardrails cannot
optimize; it can only rewrite and hope.

## Revision

**Revision** is the loop's return arc: what a later stage learns becomes a change to intent.
Without it a River module records what someone believed once.

**Revision edits what was authored, never what was derived.** A source and a goal are written by
people; a composed package is produced from them and is immutable once tagged. An edit there is
discarded by the next composition and invisible to anyone who imported it. So feedback lands on
the source or the goal, and composition runs again.

Where the source is **imported** you cannot edit it. The change is either a local requirement
overriding it under precedence, or a request to whoever owns it — and under a Mandatory or
Governing tier, only the request.

| From | What comes back |
| --- | --- |
| **Implementation** | Newly discovered requirements; requirements that cannot be satisfied as stated |
| **Verification** | Requirements that cannot be verified |
| **Release** | Telemetry and product feedback |
| **Goals** | Intent moved |

**Newly discovered** is one of three things: a requirement, so add it; an existing requirement
that was too vague, so sharpen that one; or an implementation decision that is not a requirement
at all. The last is the most common, and the usual mistake is to write it down anyway.

**Unsatisfiable** has three dispositions: weaken it to what is achievable, withdraw it
(`Status: withdrawn — <why>`, [`requirements.md`](requirements.md), *Status of a source item*),
or **accept the gap** — keep the requirement and record in `decisions.md` that it is not met, by
whom, on what grounds, and when it will be looked at again, in its `Expires:` field. An accepted
gap with no review date is indistinguishable from a forgotten one.

**Release feedback is not a weaker verification.** A guardrail catches only violations of
criteria that exist, so it cannot see anything outside the written set, and that is where most
of what a running system teaches you lies. Bug reports, feature requests, UX studies, traces,
profiles and exceptions are the one channel that reaches the unknown. It is also the only
feedback that is not always about a defect: the answer to an opportunity may be a new **goal**.

**Goal movement has three shapes and only one invalidates anything.** A goal **added** needs new
requirements and touches none. A goal **expanded** means the requirements serving it still serve
it and no longer suffice. A goal **changed or withdrawn** puts every requirement citing it in
doubt. Treating all movement as invalidation makes teams reluctant to move goals at all.

**Realign by re-citing, not by editing.** When goals move, and from time to time regardless, cite
every requirement afresh rather than adjusting its old `Goals:` line, and mark the citations that
fit only by stretching. Weak fits that cluster are a missing goal. A goal cited by unrelated
requirements is a catch-all hiding them. Two goals cited by the same requirements are one goal.
Two goals that pull against each other are a tension to record, not a mistake: each new
requirement says which side it serves.

**Revision is not composition.** Revision edits a source or a goal; composition then produces a
new requirements version from it.
