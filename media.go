package main

import (
	"fmt"
	"maps"
	"slices"
	"sort"
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
func mediaChange(rec ebird.Record, r inat.Result) string {
	eSet := eBirdMLAssets(rec.MLCatalogNumbers)
	iSet := iNatMLAssets(r)
	var diffs []string
	if diff := mlAssetDiff(eSet, iSet); len(diff) > 0 {
		diffs = append(diffs, fmt.Sprintf("%d ML Asset IDs added to eBird: %s", len(diff), diff))
	}
	if diff := mlAssetDiff(iSet, eSet); len(diff) > 0 {
		diffs = append(diffs, fmt.Sprintf("%d ML Asset IDs removed from eBird: %s", len(diff), diff))
	}
	photoCount := len(r.Photos)
	soundCount := len(r.Sounds)
	mediaCount := photoCount + soundCount
	descCount := len(iSet)
	if descCount != mediaCount {
		diffs = append(diffs, fmt.Sprintf("iNat description lists %d ML Asset IDs, but observation has %d media files (%d photos + %d sounds)",
			descCount, mediaCount, photoCount, soundCount))
	}
	if len(diffs) == 0 {
		return ""
	}
	return strings.Join(diffs, "; ")
}

type mlAssetSet map[string]bool

func (set mlAssetSet) String() string {
	keys := slices.Collect(maps.Keys(set))
	sort.Strings(keys)
	return strings.Join(keys, " ")
}

func eBirdMLAssets(mlAssets string) mlAssetSet {
	set := mlAssetSet{}
	for _, id := range strings.Split(mlAssets, " ") {
		set[strings.TrimSpace(id)] = true
	}
	return set
}

func iNatMLAssets(r inat.Result) mlAssetSet {
	set := mlAssetSet{}
	for _, line := range strings.Split(r.Description, "\n") {
		if i := strings.Index(line, "macaulaylibrary.org/asset/"); i >= 0 {
			id := line[i+len("macaulaylibrary.org/asset/"):]
			set[strings.TrimSpace(id)] = true
		}
	}
	return set
}

func mlAssetDiff(a, b mlAssetSet) mlAssetSet {
	diff := mlAssetSet{}
	for id := range a {
		if !b[id] {
			diff[id] = true
		}
	}
	return diff
}
