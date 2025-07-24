// Package inat provides types for messages defined in the iNaturalist API.
package inat

import "github.com/google/uuid"

const (
	// iNaturalist observation fields. Look up IDs using:
	// https://www.inaturalist.org/observation_fields?order=asc&order_by=created_at
	CountField           = 1   // https://www.inaturalist.org/observation_fields/1
	LocationField        = 157 // https://www.inaturalist.org/observation_fields/157
	CountyField          = 245
	CommonNameField      = 256
	DistanceField        = 396
	ProtocolField        = 1285
	NumObserversField    = 2527
	EBirdField           = 6033
	StateOrProvinceField = 7739
)

type CreateObservation struct {
	Fields      any
	Observation Observation
}

type Observation struct {
	UUID                             uuid.UUID               `json:"uuid"`
	CaptiveFlag                      bool                    `json:"captive_flag"`
	CoordinateSystem                 string                  `json:"coordinate_system"`
	Description                      string                  `json:"description"`
	GeoX                             float64                 `json:"geo_x"`
	GeoY                             float64                 `json:"geo_y"`
	Geoprivacy                       string                  `json:"geoprivacy"`
	Latitude                         float64                 `json:"latitude"`
	License                          string                  `json:"license"`
	LocationIsExact                  bool                    `json:"location_is_exact"`
	Longitude                        float64                 `json:"longitude"`
	MakeLicenseDefault               bool                    `json:"make_license_default"`
	MakeLicensesSame                 bool                    `json:"make_licenses_same"`
	MapScale                         int                     `json:"map_scale"`
	ObservationFieldValuesAttributes []ObservationFieldValue `json:"observation_field_values_attributes"`
	ObservedOnString                 string                  `json:"observed_on_string"`
	OwnersIdentificationFromVision   bool                    `json:"owners_identification_from_vision"`
	PlaceGuess                       string                  `json:"place_guess"`
	PositionalAccuracy               float64                 `json:"positional_accuracy"`
	PositioningDevice                string                  `json:"positioning_device"`
	PositioningMethod                string                  `json:"positioning_method"`
	ProjectID                        int                     `json:"project_id"`
	PrefersCommunityTaxon            bool                    `json:"prefers_community_taxon"`
	SiteID                           int                     `json:"site_id"`
	SpeciesGuess                     string                  `json:"species_guess"`
	TagList                          string                  `json:"tag_list"`
	TaxonID                          float64                 `json:"taxon_id"`
	TaxonName                        float64                 `json:"taxon_name"`
	TimeZone                         string                  `json:"time_zone"`
}

type ObservationFieldValue struct {
	ObservationFieldID int `json:"observation_field_id"`
	Value              any `json:"value"`
}

// Returned by https://api.inaturalist.org/v2/observations
type Observations struct {
	TotalResults int      `json:"total_results"`
	Page         int      `json:"page"`
	PerPage      int      `json:"per_page"`
	Results      []Result `json:"results"`
}

type Result struct {
	UUID uuid.UUID `json:"uuid"`
	Ofvs []Ofv     `json:"ofvs"`
}

type Ofv struct {
	ID      int    `json:"id"`
	FieldID int    `json:"field_id"`
	Name    string `json:"name"`
	Value   string `json:"value"`
}
