# birdsync requirements

What birdsync must do and the constraints on how it does it, kept by River's method (the skill
at [`.agents/skills/river/`](../../.agents/skills/river/SKILL.md)). How this project applies it
is in [`../bindings.md`](../bindings.md).

| Path | What it holds |
| --- | --- |
| [`local/product/`](local/product/product.md) | `package product`: what birdsync does, from a user's point of view |
| [`local/tech/`](local/tech/tech.md) | `package tech`: constraints on the implementation |
| [`imported/`](imported/) | Other parties' terms, pinned, with `PROVENANCE.md` beside each |
| [`imports.md`](imports.md) | Every imported source, considered source, and its status |
| [`precedence.md`](precedence.md) | Which source outranks which |
| [`decisions.md`](decisions.md) | How each conflict and judgment was settled, oldest first |
| [`composed/internal/birdsync/`](composed/internal/birdsync/birdsync.md) | The composed requirements and their acceptance criteria |

Never edit `composed/` by hand: change a source and recompose. Search `^## ` in any file for its
items, or `## <name>` for one.

## Module

```
module github.com/Sajmani/birdsync/river
```

The module path names this bundle. Package names are unique within it, so a qualified name
such as [product.syncKey](local/product/product.md#synckey) resolves without an import block.
Nothing is exported yet, so no other module imports from it.
