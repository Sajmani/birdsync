# ml-terms — transcribed requirements

Tier: **mandatory**. Scope: downloading assets from the Macaulay Library CDN and re-uploading
them to iNaturalist.

These pages carry no numbering of their own, so the IDs below are assigned by this
transcription — unlike, say, a WCAG success criterion, there is no upstream number to
preserve. Each entry quotes the source so the reading can be checked. See
[PROVENANCE.md](PROVENANCE.md) for what this is and is not.

---

### `ml-terms/R1` — A contributor keeps copyright in their own media

> "You retain all intellectual property rights and copyright to any Content you submit."
> — Cornell Media Licensing Agreement, §1 Rights and Responsibility

**For birdsync:** a user re-publishing their own photo or recording to iNaturalist is
exercising their own copyright. Nothing in the agreement restricts it — the licence granted
to Cornell is explicitly *non-exclusive*.

### `ml-terms/R2` — A user may download their own media

> "note that media are not available for personal download unless they belong to you"
>
> "Have you archived your media in the Macaulay Library through eBird? If so, your original
> uploaded files are available for you to download at any time... the archive serves as your
> personal cloud storage and backup of your media files."
> — Using and requesting media

**For birdsync:** downloading a user's own assets is the sanctioned case, and the sentence
frames the archive as the user's own backup. This is what birdsync is built to do.

### `ml-terms/R3` — Another contributor's media may not be downloaded without permission

> "Assets in the collection are owned by the original author (e.g., photographer/recordist),
> and the copyright lies with each author unless otherwise indicated. It is not appropriate to
> download these assets for general third party use, without permission from the author."
> — eBird Media Upload FAQ, "Can I use photos from other people's eBird checklists?"

**For birdsync:** this is the constraint that matters. birdsync downloads from the Macaulay
CDN by asset ID, unauthenticated, and the CDN serves any asset to anyone — verified
2026-08-09 against two assets belonging to other contributors, both returning HTTP 200. So
nothing in the system prevents birdsync from copying another person's media into a user's
iNaturalist account under that user's name. See
[CR-010](../../decisions.md#cr-010--birdsync-cannot-tell-whose-media-it-is-downloading).

### `ml-terms/R4` — Media is owned by the account that uploaded it, not the checklist

> "Photos are associated with the eBird account they are uploaded under, and the rights to the
> photo belong to that account."
>
> "when someone shares a checklist with you, or you share with someone, you can see their
> media and they can see yours. However, they cannot edit or change your media, and you cannot
> change theirs."
> — eBird Media Upload FAQ

**For birdsync:** a shared checklist shows both contributors' media. Whether the recipient's
`MyEBirdData.csv` lists the *other* contributor's ML catalog numbers is **not stated on any
page vendored here**, and it is the fact CR-010 turns on.

### `ml-terms/R5` — Attribution is expected when Macaulay assets are used

> "For all uses, we ask that you follow our guidelines for media attribution when using
> Macaulay Library assets."
> — Using and requesting media

**For birdsync:** each uploaded asset's iNaturalist filename is `ML<asset id>` and the
description carries the Macaulay Library URL, so the origin is recorded. Whether that
satisfies the attribution guidelines has not been checked — the guidelines page is not
vendored here.

---

## Not adopted

Requesting media for research, education, or commercial use, and the licensing process around
it, are out of scope: birdsync never requests assets it doesn't already have IDs for, and
makes no commercial use.
