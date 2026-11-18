# Goals

Part of River's method.

## Goals are dormant

**From River 3.0, goals are optional documentation.** A module may keep a `goals/` directory
(`package goals`, one goal per `## <id>`), and a requirement may cite the goals it serves in a
`Goals:` field. River checks only that those references resolve: a `Goals:` or `Serves:` naming a
goal that does not exist fails as `unknown-goal`, and a link into `goals/` must lead to a heading
that exists. Nothing else about goals is checked, counted, summed, or stamped, and a goal may change
without anything going stale.

Keep goals if they help people read the requirements. Do not draft goals for the project unless a
person asks; whether a goal is right is never your call.

A `gradients/` directory is kept as documentation too, and noted (`gradients-set-aside`). Goals as a
checked graph, goal versions, gradients, and optimization are future phases:
`ROADMAP.md`, at the root of River's own repository, says what they were and what would bring
them back.
