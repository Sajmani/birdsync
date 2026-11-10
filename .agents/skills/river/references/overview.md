# How River works

Part of River's method. Search `^## ` for the sections; read only the one you need.

## How River works

**River is a process whose state lives in your repository.** Its goals, requirements, the
decisions that settled their conflicts, and the record of what has been built against which
version are plain files in a `river/` directory, the **River module**, beside the code they govern. Every change to
them is reviewed and versioned the way code is, and nothing about the process lives only in a
server, a tracker or someone's memory.

A **software factory** is the people, agents and automation that change a piece of software:
whatever writes, reviews, tests and ships it. River does not replace any of them. All of them
read the same files, so the factory works from one statement of what the software is for and
what it must do, and is checked against it.

River works in four stages, which run continuously rather than in sequence:

| Stage | Produces | So that |
| --- | --- | --- |
| **Requirements** | What must be true, composed from sources, with acceptance criteria | Obligations from many sources become one consistent set, none dropped silently |
| **Implementation** | Code and the designs that describe it | The software meets a published version of that set |
| **Verification** | Guardrails that check criteria | No change silently breaks a requirement |
| **Release** | A running system, ramped under rollout policy | What reaches users stays aligned with the goals |

Four processes act on the arcs between them:

| Process | Where it acts |
| --- | --- |
| **Composition** | Sources → Requirements. Assemble, resolve conflicts, generate criteria |
| **Optimization** | Every stage. Reduce the cost of future change |
| **Experimentation** | Release → Goals. Align releases with intent |
| **Revision** | Back to Goals and Sources. Return what the other stages learned |

The stages are not phases: a project is normally doing all four at once, on different changes.

## Who decides

**An agent drafts; a person decides.**

An agent transcribes sources, drafts requirements, surfaces conflicts, generates criteria, runs
guardrails, clusters unmapped and weakly mapped requirements and proposes candidate goals.

A person decides whether those are the right requirements, which one loses a conflict, whether a
gap may be accepted, and whether a target should become a backstop.

The division is by accountability, not capability. **A criterion an agent both wrote and
approved verifies nothing**, because the agent will simply have satisfied itself.

## What is new here, and what is not

River assembles a great deal of existing practice and says so: requirements management,
acceptance criteria, CI gates, progressive rollout, SLOs and error budgets, guardrail metrics,
dependency pinning, Go's package and module conventions, which it adopts directly rather than
inventing syntax, and design-first development kept in step with the code, from the
Elephant-Goldfish Model.

What is new is mostly how these are joined. Requirements are composed from many sources with
explicit precedence and a durable record of every resolution. Criteria are generated from
requirements and cite them, so coverage is derived rather than asserted. Gradients make "better"
orderable so an agent can choose between candidates. And the whole of it is files in your
repository, legible to a person and loadable as grounding context by an agent.
