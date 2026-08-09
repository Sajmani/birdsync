package main

import (
	"fmt"
	"strings"

	"github.com/Sajmani/birdsync/ebird"
	"github.com/Sajmani/birdsync/inat"
)

// mediaChange returns a string describing any differences in the ML Asset IDs
// in the eBird record vs. the iNaturalist observation. It also reports
// whether the number of assets in the iNaturalist description matches
// the number of photos and sounds in the observation itself.
//
// TODO: Correct these differences by resyncing the media.
func mediaChange(rec ebird.Record, r inat.Result) (mlAssetSet, string) {
	eSet := eBirdMLAssets(rec.MLCatalogNumbers)
	iSet, fSet := iNatMLAssets(r)
	// An asset already recorded either way is not offered for upload again: a
	// permanent failure is remembered precisely so it isn't retried (P-063).
	known := iSet
	for _, id := range fSet.ids {
		known.Add(id)
	}
	var diffs []string
	var addedMediaIDs mlAssetSet
	if diff := mlAssetDiff(eSet, known); diff.Len() > 0 {
		addedMediaIDs = diff
		diffs = append(diffs, fmt.Sprintf("%d ML Asset IDs added to eBird: %s", diff.Len(), diff))
	}
	if diff := mlAssetDiff(known, eSet); diff.Len() > 0 {
		diffs = append(diffs, fmt.Sprintf("%d ML Asset IDs removed from eBird: %s", diff.Len(), diff))
	}
	if fSet.Len() > 0 {
		// Reported on every run so the user knows something needs attention,
		// without birdsync retrying it (P-064). Deleting the line from the
		// description asks for a retry.
		diffs = append(diffs, fmt.Sprintf("%d ML Asset IDs previously failed to upload and will not be retried: %s",
			fSet.Len(), fSet))
	}
	photoCount := len(r.Photos)
	soundCount := len(r.Sounds)
	mediaCount := photoCount + soundCount
	descCount := iSet.Len()
	if descCount != mediaCount {
		diffs = append(diffs, fmt.Sprintf("iNat description lists %d ML Asset IDs, but observation has %d media files (%d photos + %d sounds)",
			descCount, mediaCount, photoCount, soundCount))
	}
	if len(diffs) == 0 {
		return mlAssetSet{}, ""
	}
	return addedMediaIDs, strings.Join(diffs, "; ")
}

type mlAssetSet struct {
	ids []string // ordered
}

func (set mlAssetSet) Len() int {
	return len(set.ids)
}

func (set mlAssetSet) String() string {
	return strings.Join(set.ids, " ")
}

func (set mlAssetSet) Has(id string) bool {
	for _, x := range set.ids {
		if x == id {
			return true
		}
	}
	return false
}

func (set *mlAssetSet) Add(id string) {
	if !set.Has(id) {
		set.ids = append(set.ids, id)
	}
}

func eBirdMLAssets(mlAssets string) mlAssetSet {
	var set mlAssetSet
	if mlAssets == "" {
		return set
	}
	for _, id := range strings.Split(mlAssets, " ") {
		set.Add(strings.TrimSpace(id))
	}
	return set
}

// failedMarker distinguishes an asset birdsync could not upload from one it
// did. Both suppress a retry, but only the uploaded set is compared against the
// media actually attached to the observation, so a recorded failure doesn't
// report a count mismatch on every run (P-063).
//
// Only the opening fragment is matched, so the guidance that follows it can be
// reworded without orphaning descriptions already written.
const failedMarker = "(upload failed"

// failedNote is the full parenthetical written into the description. It
// explains itself because the person who finds it is reading an observation,
// not this source file, and the retry mechanism is otherwise undiscoverable.
const failedNote = "(upload failed permanently; delete this line from the description to retry)"

// iNatMLAssets parses the Macaulay Library assets recorded in an observation's
// description, separating those birdsync uploaded from those the service
// permanently refused.
func iNatMLAssets(r inat.Result) (uploaded, failed mlAssetSet) {
	for _, line := range strings.Split(r.Description, "\n") {
		i := strings.Index(line, "macaulaylibrary.org/asset/")
		if i < 0 {
			continue
		}
		id := strings.TrimSpace(line[i+len("macaulaylibrary.org/asset/"):])
		if strings.Contains(line, failedMarker) {
			failed.Add(id)
		} else {
			uploaded.Add(id)
		}
	}
	return uploaded, failed
}

// assetLine renders one description line for an asset.
func assetLine(id string, ok bool) string {
	if ok {
		return "Macaulay Library Asset: " + mlAssetURL(id) + "\n"
	}
	return "Macaulay Library Asset " + failedNote + ": " + mlAssetURL(id) + "\n"
}

func mlAssetURL(id string) string {
	return "https://macaulaylibrary.org/asset/" + id
}

func mlAssetDiff(a, b mlAssetSet) mlAssetSet {
	var diff mlAssetSet
	for _, id := range a.ids {
		if !b.Has(id) {
			diff.Add(id)
		}
	}
	return diff
}
