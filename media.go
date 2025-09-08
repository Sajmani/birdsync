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
	iSet := iNatMLAssets(r)
	var diffs []string
	var addedMediaIDs mlAssetSet
	if diff := mlAssetDiff(eSet, iSet); diff.Len() > 0 {
		addedMediaIDs = diff
		diffs = append(diffs, fmt.Sprintf("%d ML Asset IDs added to eBird: %s", diff.Len(), diff))
	}
	if diff := mlAssetDiff(iSet, eSet); diff.Len() > 0 {
		diffs = append(diffs, fmt.Sprintf("%d ML Asset IDs removed from eBird: %s", diff.Len(), diff))
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

func iNatMLAssets(r inat.Result) mlAssetSet {
	var set mlAssetSet
	for _, line := range strings.Split(r.Description, "\n") {
		if i := strings.Index(line, "macaulaylibrary.org/asset/"); i >= 0 {
			id := line[i+len("macaulaylibrary.org/asset/"):]
			set.Add(strings.TrimSpace(id))
		}
	}
	return set
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
