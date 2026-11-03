# Composing requirements

The rules are in [`requirements.md`](requirements.md), sections *Precedence*, *Composition*, and
*Conflicts*, and in [`verification.md`](verification.md), *Acceptance criteria*. This is the
procedure; read those sections for the format.

1. List every source: `river/requirements/local/*/` and `river/requirements/imported/*/`.
2. Read `river/requirements/precedence.md` for the tiers. **Mandatory** always wins.
   **Governing** beats local requirements unless its owner grants an exception. **Local** is the
   project's own. **Advisory** loses to local and fills in where local says nothing.
3. For each composed package (named for its audience or scope, never a version):
   - List requirements that apply unchanged, one line each, as a link to the source item plus its
     title, grouped by area.
   - Write out in full only what composition made: a rule narrowed by a higher tier, a rule over a
     domain, an intermediary's contract. Each such item has `Composed from:` linking every source.
   - List requirements left out under `## Excluded`, each naming the decision that excluded it.
   - Every source requirement must be included, excluded, or composed from, in some package.
4. Write one acceptance criterion per independently checkable behavior in the package's
   `acceptance.md`, each with `Covers:` naming what it covers.
5. Where precedence does not decide a conflict, stop. Ask the owner; record the ruling in
   `decisions.md` before composing on it.
6. Run `river check`, then ask the owner to approve (`river hash <package-dir>` gives the line).
