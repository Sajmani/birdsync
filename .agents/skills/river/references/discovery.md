# Discovering sources for a new dependency

The rule is in [`requirements.md`](requirements.md), *Sources*: look for sources before integrating
anything new.

1. **Name the dependency** and how the change will use it: what it calls or loads, how often, what
   data it sends, and whether it redistributes anything.
2. **Find what governs it**, from its owner's own pages or files, not from memory. Before you have
   read them, list what you will look for, never what you expect to find. **For a hosted service,
   search at least once**: the web for a public service (for example "iNaturalist API terms",
   "iNaturalist API recommended practices"), the owners' internal documentation, wiki, or legal
   pages for a private one. Record each query, its date, and what it found:
   - terms of service, acceptable-use or usage policy;
   - license (for code, data, and model weights);
   - rate or usage limits, required identification (API keys, user agent), and attribution;
   - data-use and privacy terms: what happens to data you send, retention, training use;
   - for a library: its license and any notice file.
3. **If you cannot retrieve a source** (blocked, behind login, offline), say so and record it as a
   gap. **Never state terms from memory**; a plausible limit is indistinguishable from a real one.
4. **Summarize for the person who decides**: each obligation, in one line, quoted where exact
   wording matters, with its URL and retrieval date, and anything that forbids or limits the
   intended use.
5. **Wait for the decision** before writing integration code. Then record it:
   - an entry in `river/requirements/imports.md` with origin, version or retrieval date, import
     style (`queried`, `quoted`, `consulted`, …), and status (`adopted`, `partial`, or `rejected`,
     with the reason);
   - the obligations you will meet as requirements in a local source (for example a rate limit or
     attribution rule), citing the imported source;
   - any obligation knowingly not met as an accepted gap in `decisions.md`, decided by the person.
6. Run `river check`.

A dependency already in `imports.md` needs no second look unless its version or your use of it
changes.
