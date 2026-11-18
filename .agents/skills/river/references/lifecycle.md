# Implementation and revision

Part of River's method. Search `^## ` for the sections; read only the one you need.

## Implementation

Implementation produces the code and the **designs that describe it**, kept in step at every
commit. A design may lead the code: written and reviewed before the code exists, it says it is a
proposal, and it becomes the code's description in the commit where the code lands. From then on
a change that makes the two disagree updates both, and a guardrail can detect drift between them.
This follows the Elephant-Goldfish Model, which reviews a design with zero-context agents before
code is written and reports where they drift.

**A design lives with the code, not in the River module.** The module holds intent and
conformance; `conformance/` cites where the designs are.

## Future phases

River 3.0 covers requirements, composition, acceptance criteria, and build-time guardrails.
Release and runtime guardrails, experimentation, gradients and optimization, and goals as a
checked graph are future phases; `ROADMAP.md`, at the root of River's own repository, describes
each and what would bring it back. Until then, record a project's own practice for them in
`bindings/`.

## Revision

**Revision** is the return arc: what a later stage learns becomes a change to intent. Without it a
River module records what someone believed once.

**Revision edits what was authored, never what was derived.** A source is written by people; a
composed package is produced from it. An edit there is discarded by the next composition and
invisible to anyone who imported it. So feedback lands on the source, and composition runs again.

Where the source is **imported** you cannot edit it. The change is either a local requirement
overriding it under precedence, or a request to whoever owns it, and under a Mandatory or Governing
tier, only the request.

| From | What comes back |
| --- | --- |
| **Implementation** | Newly discovered requirements; requirements that cannot be satisfied as stated |
| **Verification** | Requirements that cannot be verified |
| **Release** | Bug reports, feature requests, telemetry, and product feedback |

**Newly discovered** is one of three things: a requirement, so add it; an existing requirement
that was too vague, so sharpen that one; or an implementation decision that is not a requirement
at all. The last is the most common, and the usual mistake is to write it down anyway.

**Unsatisfiable** has three dispositions: weaken it to what is achievable, withdraw it
(`Status: withdrawn — <why>`, [`requirements.md`](requirements.md), *Status of a source item*),
or **accept the gap**: keep the requirement and record in `decisions/` that it is not met, by
whom, on what grounds, and when it will be looked at again, in its `Expires:` field.

**Release feedback is not a weaker verification.** A guardrail catches only violations of criteria
that exist, so it cannot see anything outside the written set, and that is where most of what a
running system teaches you lies.

**Revision is not composition.** Revision edits a source; composition then produces new composed
packages from it, ending in `river compose`.
