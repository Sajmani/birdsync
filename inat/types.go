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
	Fields      any         `json:"fields,omitempty"`
	Observation Observation `json:"observation,omitempty"`
}

type Observation struct {
	UUID                             uuid.UUID               `json:"uuid,omitempty"`
	CaptiveFlag                      bool                    `json:"captive_flag,omitempty"`
	CoordinateSystem                 string                  `json:"coordinate_system,omitempty"`
	Description                      string                  `json:"description,omitempty"`
	GeoX                             float64                 `json:"geo_x,omitempty"`
	GeoY                             float64                 `json:"geo_y,omitempty"`
	Geoprivacy                       string                  `json:"geoprivacy,omitempty"`
	Latitude                         float64                 `json:"latitude,omitempty"`
	License                          string                  `json:"license,omitempty"`
	LocationIsExact                  bool                    `json:"location_is_exact,omitempty"`
	Longitude                        float64                 `json:"longitude,omitempty"`
	MakeLicenseDefault               bool                    `json:"make_license_default,omitempty"`
	MakeLicensesSame                 bool                    `json:"make_licenses_same,omitempty"`
	MapScale                         int                     `json:"map_scale,omitempty"`
	ObservationFieldValuesAttributes []ObservationFieldValue `json:"observation_field_values_attributes,omitempty"`
	ObservedOnString                 string                  `json:"observed_on_string,omitempty"`
	OwnersIdentificationFromVision   bool                    `json:"owners_identification_from_vision,omitempty"`
	PlaceGuess                       string                  `json:"place_guess,omitempty"`
	PositionalAccuracy               float64                 `json:"positional_accuracy,omitempty"`
	PositioningDevice                string                  `json:"positioning_device,omitempty"`
	PositioningMethod                string                  `json:"positioning_method,omitempty"`
	ProjectID                        int                     `json:"project_id,omitempty"`
	PrefersCommunityTaxon            bool                    `json:"prefers_community_taxon,omitempty"`
	SiteID                           int                     `json:"site_id,omitempty"`
	SpeciesGuess                     string                  `json:"species_guess,omitempty"`
	TagList                          string                  `json:"tag_list,omitempty"`
	TaxonID                          float64                 `json:"taxon_id,omitempty"`
	TaxonName                        float64                 `json:"taxon_name,omitempty"`
	TimeZone                         string                  `json:"time_zone,omitempty"`
	User                             User                    `json:"user,omitempty"`
}

type User struct {
	Login string `json:"login,omitempty"`
	ID    string `json:"id,omitempty"`
}

type ObservationFieldValue struct {
	ObservationFieldID int `json:"observation_field_id,omitempty"`
	Value              any `json:"value,omitempty"`
}

// Returned by https://api.inaturalist.org/v2/observations
type Observations struct {
	TotalResults int      `json:"total_results,omitempty"`
	Page         int      `json:"page,omitempty"`
	PerPage      int      `json:"per_page,omitempty"`
	Results      []Result `json:"results,omitempty"`
}

type Result struct {
	UUID        uuid.UUID `json:"uuid,omitempty"`
	Description string    `json:"description,omitempty"`
	Ofvs        []Ofv     `json:"ofvs,omitempty"`
	Photos      []Photo   `json:"photos,omitempty"`
	Taxon       Taxon     `json:"taxon,omitempty"`
}

type Ofv struct {
	ID      int    `json:"id,omitempty"`
	FieldID int    `json:"field_id,omitempty"`
	Name    string `json:"name,omitempty"`
	Value   string `json:"value,omitempty"`
}

type Photo struct {
	ID          int    `json:"id,omitempty"`
	Attribution string `json:"attribution,omitempty"`
	LicenseCode string `json:"license_code,omitempty"`
	URL         string `json:"url,omitempty"`
}

type Taxon struct {
	ID                  int    `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
	IconicTaxonName     string `json:"iconic_taxon_name,omitempty"`
	PreferredCommonName string `json:"preferred_common_name,omitempty"`
}
