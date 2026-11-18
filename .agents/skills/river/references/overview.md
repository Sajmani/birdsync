# How River works

Part of River's method. Search `^## ` for the sections; read only the one you need.

## How River works

**River is a process whose state lives in your repository.** Its requirements, the decisions
that settled their conflicts, and the record of what has been built against which version are
plain files in a `river/` directory, the **River module**, beside the code they govern. Every
change to them is reviewed and versioned the way code is, and nothing about the process lives only
in a server, a tracker or someone's memory.

A **software factory** is the people, agents and automation that change a piece of software:
whatever writes, reviews, tests and ships it. River does not replace any of them. All of them
read the same files, so the factory works from one statement of what the software must do, and is
checked against it.

River 3.0 defines its core phases, which run continuously rather than in sequence:

| Phase | Produces | So that |
| --- | --- | --- |
| **Requirements** | Local and imported sources, with pinned provenance | Every obligation is written down once, where it came from |
| **Composition** | Composed packages, from sources, precedence, and decisions | Obligations from many sources become one consistent set, none dropped silently |
| **Acceptance criteria** | Checkable behaviors, each covering requirements | Every requirement has a test of whether it is met |
| **Build-time guardrails** | Tests and checkers that implement criteria, watched failing | No change silently breaks a requirement |

Implementation works against a composed version, and revision returns what it learned to the
sources ([`lifecycle.md`](lifecycle.md)). Goals, gradients and optimization, release guardrails,
and experimentation are future phases, described in `ROADMAP.md` at the root of River's own
repository.

## Who decides

**An agent drafts; a person decides.**

An agent transcribes sources, drafts requirements, surfaces conflicts, generates criteria, runs
and runs guardrails.

A person decides whether those are the right requirements, which one loses a conflict, whether a
gap may be accepted, and whether a change to the requirements merges.

The division is by accountability, not capability. **A criterion an agent both wrote and
accepted verifies nothing**, because the agent will simply have satisfied itself.

## What is new here, and what is not

River assembles a great deal of existing practice and says so: requirements management,
acceptance criteria, CI gates, code review, dependency pinning, Go's package and module conventions, which it adopts directly rather than
inventing syntax, and design-first development kept in step with the code, from the
Elephant-Goldfish Model.

What is new is mostly how these are joined. Requirements are composed from many sources with
explicit precedence and a durable record of every resolution. Criteria are generated from
requirements and cite them, so coverage is derived rather than asserted. And the whole of it is files in your
repository, legible to a person and loadable as grounding context by an agent.
