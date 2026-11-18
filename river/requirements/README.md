# birdsync requirements

What birdsync must do and the constraints on how it does it, kept by River's method (the skill
at [`.agents/skills/river/`](../../.agents/skills/river/SKILL.md)). How this project applies it
is in [`../bindings/`](../bindings/bindings.md).

| Path | What it holds |
| --- | --- |
| [`local/product/`](local/product/) | `package product`: what birdsync does, from a user's point of view |
| [`local/tech/`](local/tech/) | `package tech`: constraints on the implementation |
| [`imported/`](imported/) | Other parties' terms, pinned, with `PROVENANCE.md` beside each |
| [`imports/`](imports/imports.md) | Every imported source, considered source, and its status |
| [`precedence/`](precedence/precedence.md) | Which source outranks which |
| [`decisions/`](decisions/decisions.md) | How each conflict and judgment was settled, oldest first |
| [`sums.md`](sums.md) | Per-item SHA-256 of every composition input (`river sum -w`) |
| [`composed/internal/birdsync/`](composed/internal/birdsync/README.md) | The composed requirements and their acceptance criteria, split into area files, with `composed.sum` (`river compose`) |

Never edit `composed/` by hand: change a source and recompose. Search `^## ` in any file for its
items, or `## <name>` for one.
