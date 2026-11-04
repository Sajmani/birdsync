package product

# Product requirements: Media

Copying photos and sounds from the Macaulay Library to iNaturalist. One item per heading, and
each item's first paragraph after its fields is the complete rule. Search `^## ` for the list,
or `## <name>` for one item. Other product requirements are in the other files in this
directory, all one package, `product`.

## createBeforeMedia

**The observation is created before media is attached**

Formerly: P-042
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

The observation is created first and media attached afterward.

*Rationale: iNaturalist cannot attach media to an observation that does not yet exist.*

## mediaFormats

**Photos at 2400px, sounds as MP3**

Formerly: P-043
Subject: `media.photo.max_dimension_px`
Value: `2400`
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord), [goals.communityGainsResearchGrade](../../../goals/goals.md#communitygainsresearchgrade)

Photos are fetched at 2400px and sounds as MP3.

## photoThenSoundFallback

**Try the photo URL, then the sound URL**

Formerly: P-044
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord), [goals.communityGainsResearchGrade](../../../goals/goals.md#communitygainsresearchgrade)

A Macaulay Library asset ID does not say whether it is a photo or a sound, so
birdsync tries the photo URL and falls back to the sound URL on a 404.

## mediaFilename

**Media is uploaded as ML<asset ID> with a canonical extension**

Formerly: P-045
Subject: `media.filename.extension`
Value: `{image/jpeg: .jpg, image/png: .png, audio/mpeg: .mp3, audio/mpeg3: .mp3, audio/wav: .wav}`, falling back to `.jpg` or `.mp3` by endpoint
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

Media is uploaded under the filename `ML<asset ID>` plus a canonical extension
for the response content type — `.jpg`, `.png`, `.mp3`, `.wav` — so any file in iNaturalist
can be traced back to its Macaulay Library asset.

*Rationale: the extension must not depend on the machine running the sync. Amended by
[decisions.canonicalFileExtensions](../../decisions.md#canonicalfileextensions); previously it was
whichever extension the system's mime database happened to list first.*

## descriptionUpdateKeepsMedia

**The description update keeps attached media**

Formerly: P-046
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

After uploading, the observation description is updated with the asset URLs,
without detaching the media already attached.

## additiveMediaResync

**Media re-sync is additive**

Formerly: P-047
Goals: [goals.birdersKeepOneRecord](../../../goals/goals.md#birderskeeponerecord)

Media re-sync is additive. Assets added to an eBird checklist after a sync
are uploaded on the next run.

## removedAssetsReported

**Assets removed from eBird are reported, not removed**

Formerly: P-048
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

Assets removed from eBird are reported, not removed from iNaturalist.

## mediaCountMismatchReported

**A media count mismatch is reported, not corrected**

Formerly: P-049
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

A mismatch between the asset count in the description and the media actually
attached is reported, not corrected.

## mediaFailuresTolerated

**A failed media transfer is logged, counted, and survived**

Formerly: P-050
Goals: [goals.accountsStaySafe](../../../goals/goals.md#accountsstaysafe)

A failed media download or upload is logged and counted, and does not stop
the run.

*The description lists only assets that actually uploaded, so a failure is retried on the
next run rather than recorded as done. See
[decisions.descriptionListsOnlyUploaded](../../decisions.md#descriptionlistsonlyuploaded).*
