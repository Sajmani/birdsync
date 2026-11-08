# River module

This repository's goals, requirements, and the records River keeps about them. Each entry here
says what it holds in its own README or first lines.

- [`bindings.md`](bindings.md)
- [`conformance.md`](conformance.md)
- [`goals/`](goals/)
- [`proposals.md`](proposals.md)
- [`requirements/`](requirements/)

## Module

```
module github.com/Sajmani/birdsync/river
river 2.0
```

The module path names this module. Package names are unique within it, so a qualified name
such as [product.syncKey](requirements/local/product/identity.md#synckey) resolves without an import block.
Nothing is exported yet, so no other module imports from it.
