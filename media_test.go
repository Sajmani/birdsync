package main

import (
	"slices"
	"testing"

	"github.com/Sajmani/birdsync/ebird"
	"github.com/Sajmani/birdsync/inat"
)

func TestMediaChange(t *testing.T) {
	tests := []struct {
		name          string
		rec           ebird.Record
		r             inat.Result
		wantSet       mlAssetSet
		wantDesc      string
		wantNeedsUpdt bool
	}{
		{
			name: "no media",
			rec:  ebird.Record{},
			r:    inat.Result{},
		},
		{
			name: "same media",
			rec: ebird.Record{
				MLCatalogNumbers: "12345",
			},
			r: inat.Result{
				Description: mlAssetURL("12345"),
				Photos:      []inat.Photo{{OriginalFilename: "ML12345.jpg"}},
			},
		},
		{
			name: "added to ebird",
			rec: ebird.Record{
				MLCatalogNumbers: "12345 67890",
			},
			r: inat.Result{
				Photos: []inat.Photo{{OriginalFilename: "ML12345.jpg"}},
			},
			wantSet:       mlAssetSet{ids: []string{"67890"}},
			wantDesc:      "1 ML Asset IDs added to eBird: 67890; 1 ML Asset IDs in description need to be added: 12345",
			wantNeedsUpdt: true,
		},
		{
			name: "removed from ebird",
			rec: ebird.Record{
				MLCatalogNumbers: "12345",
			},
			r: inat.Result{
				Photos: []inat.Photo{
					{OriginalFilename: "ML12345.jpg"},
					{OriginalFilename: "ML67890.jpg"},
				},
			},
			wantDesc:      "1 ML Asset IDs removed from eBird: 67890; 2 ML Asset IDs in description need to be added: 12345 67890",
			wantNeedsUpdt: true,
		},
		{
			name: "description out of sync",
			rec: ebird.Record{
				MLCatalogNumbers: "12345",
			},
			r: inat.Result{
				Description: mlAssetURL("67890"),
				Photos: []inat.Photo{
					{OriginalFilename: "ML12345.jpg"},
				},
			},
			wantDesc:      "1 ML Asset IDs in description need to be added: 12345; 1 ML Asset IDs in description need to be removed: 67890",
			wantNeedsUpdt: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSet, got, needsUpdt := mediaChange(tt.rec, tt.r)
			if got != tt.wantDesc {
				t.Errorf("mediaChange() got = %q, want %q", got, tt.wantDesc)
			}
			if !slices.Equal(gotSet.ids, tt.wantSet.ids) {
				t.Errorf("mediaChange() gotSet = %v, wantSet %v", gotSet, tt.wantSet)
			}
			if needsUpdt != tt.wantNeedsUpdt {
				t.Errorf("mediaChange() needsUpdt = %v, want %v", needsUpdt, tt.wantNeedsUpdt)
			}
		})
	}
}

func TestRebuildDescription(t *testing.T) {
	tests := []struct {
		name    string
		desc    string
		assets  mlAssetSet
		want    string
	}{
		{
			name: "no assets",
			desc: "foo\nbar",
			want: "foo\nbar",
		},
		{
			name: "add assets",
			desc: "foo\nbar",
			assets: mlAssetSet{ids: []string{"12345", "67890"}},
			want: "foo\nbar\n" + mlAssetURL("12345") + "\n" + mlAssetURL("67890"),
		},
		{
			name: "remove assets",
			desc: "foo\n" + mlAssetURL("12345") + "\nbar",
			want: "foo\nbar",
		},
		{
			name: "add and remove assets",
			desc: "foo\n" + mlAssetURL("12345") + "\nbar",
			assets: mlAssetSet{ids: []string{"67890"}},
			want: "foo\nbar\n" + mlAssetURL("67890"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rebuildDescription(tt.desc, tt.assets); got != tt.want {
				t.Errorf("rebuildDescription() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMLAssetID(t *testing.T) {
	tests := []struct {
		filename string
		want     string
	}{
		{"", ""},
		{"foo.jpg", ""},
		{"ML.jpg", ""},
		{"ML12345.jpg", "12345"},
		{"ML12345", "12345"},
		{"ML12345.foo.bar", "12345.foo"},
	}
	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			if got := mlAssetID(tt.filename); got != tt.want {
				t.Errorf("mlAssetID(%q) = %q, want %q", tt.filename, got, tt.want)
			}
		})
	}
}
