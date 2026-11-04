# ml-terms — provenance

Vendored copy of the Cornell Lab / eBird / Macaulay Library terms that govern what birdsync
does with media. See `river/requirements/imports.md` for the adoption record and
[mlterms.md](mlterms.md) for the transcription.

## What was retrieved

| File | Origin | SHA-256 |
| --- | --- | --- |
| `using-and-requesting-media.html` | <https://support.ebird.org/en/support/solutions/articles/48001064551-using-and-requesting-media> | `40b96ef7d4e1c1868fc93cd83d8b9f590487e0ef4ada384ac67f5a5615b74a62` |
| `ebird-media-upload-faq.html` | <https://support.ebird.org/en/support/solutions/articles/48000960529-ebird-media-upload-faq> | `7cd20a7ea09f8d7d78bb447f32913264dcaa3c24704bbc4d31a15f39d566f816` |
| `cornell-media-licensing-agreement.html` | <https://support.ebird.org/en/support/solutions/articles/48000952192-cornell-media-licensing-agreement> | `8db7e25c0b3c8163a741fe1446507fba18c268181573ecd1b3b2e4cc6b60d558` |

- **Retrieved:** 2026-08-09, by `curl`, as served (raw HTML, not a rendering).
- **Retrieved by:** an AI agent working with the maintainer, unreviewed by counsel.
- **Pages' own "modified on" dates:** 29 May 2025, 1 May 2025, and 6 July 2026 respectively.

## What this is not

**This is not legal advice and has not been reviewed by a lawyer.** It is a snapshot of
publicly published guidance, kept so that a future reader can see exactly what the terms said
when the requirements were written, and diff them when they change.

The transcription in `mlterms.md` is an interpretation. Each entry quotes the source
verbatim so the interpretation can be checked against it, per
River's rule on sources (requirements.md, *Sources*, in the River skill): a transcription is a reviewable artifact, not
ground truth.

## Pages considered and not vendored

- The Cornell Lab Terms of Use (linked from the licensing agreement) — the media-specific
  pages above are more precise about the case birdsync is in. Worth adding if a question
  arises that they don't answer.
- The eBird data privacy and data use page — concerns observation records rather than media.

## Refreshing

Re-fetch the URLs, compare hashes, and read the diff before updating the pin. The licensing
agreement changed as recently as July 2026, so it is the one most likely to move.

## Credential-shaped strings

The three pages, as served, embed Freshdesk presigned logo URLs whose `X-Amz-Credential`
parameter holds an AWS-style access key ID. The URLs are public and expired, but secret
scanners flag them. They are kept as served, because the hashes above describe these bytes.

**Policy for the next re-fetch:** strip the presigned `X-Amz-*` query parameters from vendored
HTML before committing it. Record here what was removed, and both hashes, as served and as
stored, the way `inatterms` records its removed `<script>` elements.
