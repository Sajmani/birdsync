# inat-terms — transcribed requirements

Tier: **mandatory**. Scope: everything birdsync writes to a user's iNaturalist account, and
the manner in which it writes it.

IDs are assigned by this transcription. Each entry quotes the source. See
[PROVENANCE.md](PROVENANCE.md), including what is missing.

---

### `inat-terms/R1` — A contributor keeps their rights, and licenses by default as CC BY-NC

> "If you own the Content prior to contributing the Content to the Site, and that Content is
> subject to Intellectual Property rights, you retain ownership of those rights. Unless you
> specify otherwise when you post Content, you agree to license Content you contribute to the
> Platform under the Creative Commons Attribution Noncommercial license (CC BY-NC)."
> — Terms of Use, Responsibility of Contributors

**For birdsync:** uploading a user's own Macaulay Library media applies that user's own
default licence to their own work. Note birdsync never sets a licence explicitly, so every
upload takes the account default — worth knowing, since Macaulay and iNaturalist licences are
chosen separately and need not agree.

### `inat-terms/R2` — The account holder is responsible for everything done under the account

> "You are fully responsible for all activities that occur under the account and any other
> actions taken in connection with the account."
> — Terms of Use, Your iNaturalist Account

**For birdsync:** the user carries the consequences of what the tool posts on their behalf.
That is the reason `--dryrun` exists and the reason its guarantee is worth testing directly
(P-051).

### `inat-terms/R3` — Post at a rate and volume that doesn't hinder other users

> "You will post only Content that is relevant to iNaturalist and at a rate and volume that
> does not hinder other Users' ability to use iNaturalist"
> — Terms of Use, Responsibility of Contributors

**For birdsync:** a second reason for T-035's pacing, independent of `inat-api/R1`. This one
is in the terms rather than in guidance, so it binds harder.

### `inat-terms/R4` — Machine generated content is prohibited *(hard rule)*

> "(!) Machine generated observations, identifications, comments, and messages. We do not
> allow machines to generate and post content on iNat with no human oversight curating each
> piece of content, and any account suspected of doing so is subject to suspension and the
> removal of the content."
> — Community Guidelines, marked (!): "grounds for immediate suspension without warning"

**For birdsync: permitted.** The rule's scope is set by
<https://www.inaturalist.org/pages/machine_generated_content>, vendored here, which lists
"writing a script to create observations from a manually curated local folder of your images
and metadata on your desktop" and "creating a third-party app that enables a real person/people
to create content" among its examples of *acceptable* behavior, and permits using a machine
"to facilitate posting this content" where a human decided each observation. See
[CR-012](../../decisions.md#cr-012--is-birdsync-machine-generated-content), closed.

### `inat-terms/R5` — Don't behave like a machine

> "Any account that adds content we believe decreases the accuracy of iNaturalist data may be
> suspended, particularly if that account behaves like a machine, e.g. adds a lot of content
> very quickly and does not respond to comments and messages."
> — Community Guidelines

**For birdsync:** two obligations. Pacing (T-035) addresses "very quickly". The second half
falls on the user: an account that syncs hundreds of observations and never answers an
identifier's comment looks exactly like the thing this describes. birdsync cannot discharge
that for them, but it can say so — currently it does not.

### `inat-terms/R6` — Duplicates are undesirable but not an offence

> "Duplicate observations. They're not ideal, but they're usually due to oversight or bugs.
> Politely ask people to remove them but if they don't, it's not a big deal unless it becomes
> a habit."
> — Community Guidelines

**For birdsync:** calibrates CR-003 and CR-011. Duplicates are a data-quality problem and an
embarrassment, not a suspension risk — unless they become habitual, which an unattended tool
generating them every run would certainly qualify as.

---

## Not adopted

Hate speech, harassment, sockpuppets, sexually explicit content, and the prohibition on
commercial AI training are all binding on the user but unreachable by birdsync, which posts
only observation records derived from their eBird export.
