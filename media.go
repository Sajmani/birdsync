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
func mediaChange(rec ebird.Record, r inat.Result) (addedMediaIDs mlAssetSet, desc string, needsUpdate bool) {
	eSet := eBirdMLAssets(rec.MLCatalogNumbers)
	iSet := iNatMLAssets(r)

	var diffs []string
	if diff := mlAssetDiff(eSet, iSet); diff.Len() > 0 {
		addedMediaIDs = diff
		diffs = append(diffs, fmt.Sprintf("%d ML Asset IDs added to eBird: %s", diff.Len(), diff))
	}
	if diff := mlAssetDiff(iSet, eSet); diff.Len() > 0 {
		diffs = append(diffs, fmt.Sprintf("%d ML Asset IDs removed from eBird: %s", diff.Len(), diff))
	}

	iSetDesc := iNatMLAssetsFromDescription(r.Description)
	if diff := mlAssetDiff(iSet, iSetDesc); diff.Len() > 0 {
		diffs = append(diffs, fmt.Sprintf("%d ML Asset IDs in description need to be added: %s", diff.Len(), diff))
	}
	if diff := mlAssetDiff(iSetDesc, iSet); diff.Len() > 0 {
		diffs = append(diffs, fmt.Sprintf("%d ML Asset IDs in description need to be removed: %s", diff.Len(), diff))
	}

	if len(diffs) == 0 {
		return mlAssetSet{}, "", false
	}
	return addedMediaIDs, strings.Join(diffs, "; "), true
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

func iNatMLAssetsFromDescription(desc string) mlAssetSet {
	var set mlAssetSet
	for _, line := range strings.Split(desc, "\n") {
		if i := strings.Index(line, "macaulaylibrary.org/asset/"); i >= 0 {
			id := line[i+len("macaulaylibrary.org/asset/"):]
			set.Add(strings.TrimSpace(id))
		}
	}
	return set
}

func rebuildDescription(desc string, assets mlAssetSet) string {
	var newDesc []string
	for _, line := range strings.Split(desc, "\n") {
		if !strings.Contains(line, "macaulaylibrary.org/asset/") {
			newDesc = append(newDesc, line)
		}
	}
	for _, id := range assets.ids {
		newDesc = append(newDesc, mlAssetURL(id))
	}
	return strings.Join(newDesc, "\n")
}

func iNatMLAssets(r inat.Result) mlAssetSet {
	var set mlAssetSet
	for _, photo := range r.Photos {
		if id := mlAssetID(photo.OriginalFilename); id != "" {
			set.Add(id)
		}
	}
	for _, sound := range r.Sounds {
		if id := mlAssetID(sound.OriginalFilename); id != "" {
			set.Add(id)
		}
	}
	return set
}

// mlAssetID returns the Macaulay Library asset ID from a filename, or "" if
// the filename does not appear to be a Macaulay Library asset.
//
// Example filenames:
//   - ML12345.jpg
//   - ML67890.wav
func mlAssetID(filename string) string {
	if !strings.HasPrefix(filename, "ML") {
		return ""
	}
	filename = strings.TrimPrefix(filename, "ML")
	if i := strings.LastIndex(filename, "."); i >= 0 {
		return filename[:i]
	}
	return filename
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
