# Goals and gradients

Part of River's method. Search `^## ` for the sections; read only the one you need.

## Goals

A goal is an **objective**: why the software exists. It is what a requirement is in service of,
not the requirement restated in looser words.

**A goal need not be measurable, but it may be.** Measurement enters through gradients.

**Goals and requirements differ in how they're met.** A release meets a requirement: ship it
and the requirement is true, and it has to stay true while the release runs. You may need a week
of traffic to prove p99 is under 200 ms, but the release is what makes it so. Releases meet a
goal over time, as people respond to them. No release gets you 10,000 weekly users on its own,
and some goals you only ever get closer to.

So **a goal never blocks a release. A requirement can.**

**Requirements can cover how you ship, not just what you ship.** "Fix critical vulnerabilities
within seven days" is checked every time one is published, and it can block a release that would
miss the deadline.

**A goal names who benefits and what they gain.** The usual beneficiaries are the software's
users, the communities it depends on, and the people who build it; a project with no goal for one
of them is usually missing one. A goal that names a mechanism is a requirement in disguise: it
should survive a rewrite of the implementation. One paragraph is enough.

Write a goal as what someone gets, not what the software lets them do: *readers learn how big
birds are*, not *readers can see a bird's size*. If one release can ship it, it's a requirement.

Requirements cite goals by package-qualified ID in a `Goals:` field on each local source item and
each composed item (included bullets in a composed package link to their source item rather than
repeating `Goals:`). A goal never lists its requirements, so the two cannot drift apart.

A goal nothing cites yet, and a requirement citing no goal, are both **counted, not rejected**.
A brownfield project starts with the second; a goal added because the market moved starts with
the first. Requiring either to arrive with the other attached produces invented ones.

## Gradients

A **gradient** is a metric with a declared **direction of improvement**. The direction is the
only difference from a KPI: a KPI reports where you are; a gradient says which of two candidates
is better, which is what an agent needs when it has produced several.

A gradient may carry a **target**, a **backstop**, both, or neither.

```
  worse  ├──────────────●───────────→──────────┤  better
      Backstop        current               Target
```

- A **target** is absolute: where you are trying to get to. A gradient with one is an SLO.
- A **backstop** is a fixed limit the value must not cross. It does not move when you
  improve: reach 90 ms against a 100 ms backstop and the 10 ms between is slack you may spend
  elsewhere, as an error budget is spent. Performance budgets are backstops.
- **Both** brackets the work: *reach 100 ms, fail above 300 ms, improve in between.*

**The kind of threshold does not decide whether you have written a goal or a requirement** —
attainment does. Either kind can produce either, though backstops are usually requirements,
because every change can be checked against one.

**A ratchet is not a backstop.** A ratchet moves the limit to each new best, which throws the
slack away. It is occasionally the right policy, but moving a backstop is a decision, and a
person's to make.

**Most gradients carry no threshold, and should not.** Where several exist, nearly every change
regresses one of them; thresholding all of them stops delivery, and teams learn to override
gates. Online-experiment practice calls a thresholded metric a **guardrail metric** and warns
against having many, for the same reason.

**Achieving a target may promote it to a backstop.** A target gives direction until you reach
it and then gives none; converting it holds the position while effort moves to goals still
unmet. That is what makes a set of goals prioritizable: an agent reads the achieved ones as
constraints and the rest as direction. Promotion is a person's call. It needs a firing rule or
it flaps, a level below the target or it gates on the first ordinary dip, and a cap, because
ten achieved goals would otherwise become ten permanent gates.

Not every goal has a gradient. A fabricated measure is worse than none, because it looks like
evidence.
