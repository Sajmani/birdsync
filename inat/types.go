package inat

import "github.com/google/uuid"

const (
	// iNaturalist observation fields. Look up IDs using:
	// https://www.inaturalist.org/observation_fields?order=asc&order_by=created_at
	CountField               = 1   // https://www.inaturalist.org/observation_fields/1
	LocationField            = 157 // https://www.inaturalist.org/observation_fields/157
	CountyField              = 245
	CommonNameField          = 256
	DistanceField            = 396
	NumObserversField        = 2527
	EBirdField               = 6033
	StateOrProvinceField     = 7739
	EBirdScientificNameField = 20215
)

type CreateObservation struct {
	Fields      any         `json:"fields,omitempty"`
	Observation Observation `json:"observation,omitempty"`
}

type UpdateObservation struct {
	Fields       any         `json:"fields,omitempty"`
	IgnorePhotos bool        `json:"ignore_photos,omitempty"`
	Observation  Observation `json:"observation,omitempty"`
}

type Observation struct {
	CaptiveFlag                      bool                    `json:"captive_flag,omitempty"`
	CoordinateSystem                 string                  `json:"coordinate_system,omitempty"`
	Description                      string                  `json:"description,omitempty"`
	Geoprivacy                       string                  `json:"geoprivacy,omitempty"`
	GeoX                             float64                 `json:"geo_x,omitempty"`
	GeoY                             float64                 `json:"geo_y,omitempty"`
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
	PrefersCommunityTaxon            bool                    `json:"prefers_community_taxon,omitempty"`
	ProjectID                        int                     `json:"project_id,omitempty"`
	SiteID                           int                     `json:"site_id,omitempty"`
	SpeciesGuess                     string                  `json:"species_guess,omitempty"`
	TagList                          string                  `json:"tag_list,omitempty"`
	TaxonID                          float64                 `json:"taxon_id,omitempty"`
	TaxonName                        float64                 `json:"taxon_name,omitempty"`
	TimeZone                         string                  `json:"time_zone,omitempty"`
	User                             User                    `json:"user,omitempty"`
	UUID                             uuid.UUID               `json:"uuid,omitempty"`
}

type User struct {
	Login string `json:"login,omitempty"`
	ID    string `json:"id,omitempty"`
}

type ObservationFieldValue struct {
	ObservationFieldID int `json:"observation_field_id,omitempty"`
	Value              any `json:"value,omitempty"`
}

// Observations is returned by https://api.inaturalist.org/v2/observations
type Observations struct {
	Page         int      `json:"page,omitempty"`
	PerPage      int      `json:"per_page,omitempty"`
	Results      []Result `json:"results,omitempty"`
	TotalResults int      `json:"total_results,omitempty"`
}

type Result struct {
	CreatedAt            string    `json:"created_at,omitempty"`
	Description          string    `json:"description,omitempty"`
	IdentificationsCount int       `json:"identifications_count,omitempty"`
	ObservedOn           string    `json:"observed_on,omitempty"`
	Ofvs                 []Ofv     `json:"ofvs,omitempty"`
	Photos               []Photo   `json:"photos,omitempty"`
	PositionalAccuracy   int       `json:"positional_accuracy,omitempty"`
	PreferredCommonName  string    `json:"preferred_common_name,omitempty"`
	QualityGrade         string    `json:"quality_grade,omitempty"`
	Sounds               []Sound   `json:"sounds,omitempty"`
	Taxon                Taxon     `json:"taxon,omitempty"`
	UUID                 uuid.UUID `json:"uuid,omitempty"`
}

// ObservationFieldValue returns the value of the observation field
// with the given field ID.
// It returns "" if the field is empty or not found.
func (r Result) ObservationFieldValue(fieldID int) string {
	for _, ofv := range r.Ofvs {
		if ofv.FieldID == fieldID {
			return ofv.Value
		}
	}
	return ""
}

type Ofv struct {
	FieldID int    `json:"field_id,omitempty"`
	ID      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Value   string `json:"value,omitempty"`
}

type Photo struct {
	Attribution string `json:"attribution,omitempty"`
	ID          int    `json:"id,omitempty"`
	LicenseCode string `json:"license_code,omitempty"`
	URL         string `json:"url,omitempty"`
}

type Sound struct {
	Attribution string `json:"attribution,omitempty"`
	ID          int    `json:"id,omitempty"`
	LicenseCode string `json:"license_code,omitempty"`
	FileURL     string `json:"file_url,omitempty"`
}

type Taxon struct {
	IconicTaxonName     string `json:"iconic_taxon_name,omitempty"`
	ID                  int    `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
	PreferredCommonName string `json:"preferred_common_name,omitempty"`
}
