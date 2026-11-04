package inatterms

# inat-terms — transcribed requirements

Tier: **mandatory**. Scope: everything birdsync writes to a user's iNaturalist account, and
the manner in which it writes it.

IDs are assigned by this transcription. Each entry quotes the source. See
[PROVENANCE.md](PROVENANCE.md), including what is missing.

Identifiers were renamed from the transcription's `R#` numbers to descriptive names when
birdsync adopted River; `Formerly:` gives the old one. The quotations are unchanged.

## contributorKeepsRights

**A contributor keeps their rights, and licenses by default as CC BY-NC**

Formerly: inat-terms/R1

> "If you own the Content prior to contributing the Content to the Site, and that Content is
> subject to Intellectual Property rights, you retain ownership of those rights. Unless you
> specify otherwise when you post Content, you agree to license Content you contribute to the
> Platform under the Creative Commons Attribution Noncommercial license (CC BY-NC)."
> — Terms of Use, Responsibility of Contributors

**For birdsync:** uploading a user's own Macaulay Library media applies that user's own
default license to their own work. Note birdsync never sets a license explicitly, so every
upload takes the account default — worth knowing, since Macaulay and iNaturalist licenses are
chosen separately and need not agree.

## accountHolderResponsible

**The account holder is responsible for everything done under the account**

Formerly: inat-terms/R2

> "You are fully responsible for all activities that occur under the account and any other
> actions taken in connection with the account."
> — Terms of Use, Your iNaturalist Account

**For birdsync:** the user carries the consequences of what the tool posts on their behalf.
That is the reason `--dryrun` exists and the reason its guarantee is worth testing directly
(`dryRunNoWrites`).

## rateAndVolume

**Post at a rate and volume that doesn't hinder other users**

Formerly: inat-terms/R3

> "You will post only Content that is relevant to iNaturalist and at a rate and volume that
> does not hinder other Users' ability to use iNaturalist"
> — Terms of Use, Responsibility of Contributors

**For birdsync:** a second reason for `requestPacing`'s pacing, independent of `inat-api/R1`. This one
is in the terms rather than in guidance, so it binds harder.

## noMachineGeneratedContent

**Machine generated content is prohibited**

Formerly: inat-terms/R4
Hard rule: yes, marked (!) by the source

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
the decision `notMachineGeneratedContent`, closed.

## dontBehaveLikeMachine

**Don't behave like a machine**

Formerly: inat-terms/R5

> "Any account that adds content we believe decreases the accuracy of iNaturalist data may be
> suspended, particularly if that account behaves like a machine, e.g. adds a lot of content
> very quickly and does not respond to comments and messages."
> — Community Guidelines

**For birdsync:** two obligations. Pacing (`requestPacing`) addresses "very quickly". The second half
falls on the user: an account that syncs hundreds of observations and never answers an
identifier's comment looks exactly like the thing this describes. birdsync cannot discharge
that for them, but it can say so — currently it does not.

## duplicatesUndesirable

**Duplicates are undesirable but not an offense**

Formerly: inat-terms/R6

> "Duplicate observations. They're not ideal, but they're usually due to oversight or bugs.
> Politely ask people to remove them but if they don't, it's not a big deal unless it becomes
> a habit."
> — Community Guidelines

**For birdsync:** calibrates the decision `noTaxonFilterOnDownload` and the decision `pageByIDAbove`. Duplicates are a data-quality problem and an
embarrassment, not a suspension risk — unless they become habitual, which an unattended tool
generating them every run would certainly qualify as.

---

## Not transcribed

Hate speech, harassment, sockpuppets, sexually explicit content, and the prohibition on
commercial AI training are all binding on the user but unreachable by birdsync, which posts
only observation records derived from their eBird export.
