package main

import (
	"slices"
	"testing"

	"github.com/Sajmani/birdsync/ebird"
	"github.com/Sajmani/birdsync/inat"
)

func TestMediaChange(t *testing.T) {
	tests := []struct {
		name    string
		rec     ebird.Record
		r       inat.Result
		wantSet mlAssetSet
		want    string
	}{
		{
			name: "no media",
			rec:  ebird.Record{},
			r:    inat.Result{},
			want: "",
		},
		{
			name: "same media",
			rec: ebird.Record{
				MLCatalogNumbers: "12345",
			},
			r: inat.Result{
				Description: "macaulaylibrary.org/asset/12345",
				Photos:      []inat.Photo{{}},
			},
			want: "",
		},
		{
			name: "added to ebird",
			rec: ebird.Record{
				MLCatalogNumbers: "12345 67890",
			},
			r: inat.Result{
				Description: "macaulaylibrary.org/asset/12345",
				Photos:      []inat.Photo{{}},
			},
			wantSet: mlAssetSet{ids: []string{"67890"}},
			want:    "1 ML Asset IDs added to eBird: 67890",
		},
		{
			name: "removed from ebird",
			rec: ebird.Record{
				MLCatalogNumbers: "12345",
			},
			r: inat.Result{
				Description: "macaulaylibrary.org/asset/12345\nmacaulaylibrary.org/asset/67890",
				Photos:      []inat.Photo{{}, {}},
			},
			want: "1 ML Asset IDs removed from eBird: 67890",
		},
		{
			name: "media count mismatch",
			rec: ebird.Record{
				MLCatalogNumbers: "12345",
			},
			r: inat.Result{
				Description: "macaulaylibrary.org/asset/12345",
				Photos:      []inat.Photo{{}, {}},
			},
			want: "iNat description lists 1 ML Asset IDs, but observation has 2 media files (2 photos + 0 sounds)",
		},
		{
			name: "added, removed, and count mismatch",
			rec: ebird.Record{
				MLCatalogNumbers: "12345 99999",
			},
			r: inat.Result{
				Description: "macaulaylibrary.org/asset/12345\nmacaulaylibrary.org/asset/67890",
				Photos:      []inat.Photo{{}},
				Sounds:      []inat.Sound{{}},
			},
			wantSet: mlAssetSet{ids: []string{"99999"}},
			want:    "1 ML Asset IDs added to eBird: 99999; 1 ML Asset IDs removed from eBird: 67890",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSet, got := mediaChange(tt.rec, tt.r)
			if got != tt.want {
				t.Errorf("mediaChange() got = %q, want %q", got, tt.want)
			}
			if !slices.Equal(gotSet.ids, tt.wantSet.ids) {
				t.Errorf("mediaChange() gotSet = %v, wantSet %v", gotSet, tt.wantSet)
			}
		})
	}
}
