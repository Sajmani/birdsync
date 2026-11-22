package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/Sajmani/birdsync/inat"
	"github.com/google/uuid"
)

const (
	stateFileVersion = 1
	// syncOverlap subtracts 1 minute from the cached LastSyncedAt watermark
	// when querying updated_since, guarding against minor clock skew or writes
	// that raced with the previous download.
	syncOverlap = time.Minute
)

// defaultStatePath returns the default path for the local sync state cache:
// $XDG_STATE_HOME/birdsync/state.json when XDG_STATE_HOME is set, falling back
// to os.UserCacheDir()/birdsync/state.json.
func defaultStatePath() string {
	if xdg := strings.TrimSpace(os.Getenv("XDG_STATE_HOME")); xdg != "" {
		return filepath.Join(xdg, "birdsync", "state.json")
	}
	dir, err := os.UserCacheDir()
	if err != nil || dir == "" {
		return ""
	}
	return filepath.Join(dir, "birdsync", "state.json")
}

// stateFile holds cached sync state across iNaturalist accounts.
// It never stores credentials or API tokens.
type stateFile struct {
	Version  int                      `json:"version"`
	Accounts map[string]*accountState `json:"accounts"`
}

type accountState struct {
	// LastSyncedAt is the UTC watermark for full-history completeness.
	// Only unwindowed runs advance it.
	LastSyncedAt time.Time                    `json:"last_synced_at,omitempty"`
	Observations map[string]cachedObservation `json:"observations,omitempty"`
	// RecentRequests records exact UTC timestamps of iNaturalist API requests
	// made in the rolling 24-hour window.
	RecentRequests []time.Time `json:"recent_requests,omitempty"`
}

type cachedObservation struct {
	ID                  int        `json:"id,omitempty"`
	UUID                uuid.UUID  `json:"uuid"`
	CreatedAt           string     `json:"created_at,omitempty"`
	Description         string     `json:"description,omitempty"`
	ObservedOn          string     `json:"observed_on,omitempty"`
	Ofvs                []inat.Ofv `json:"ofvs,omitempty"`
	PhotoCount          int        `json:"photo_count,omitempty"`
	SoundCount          int        `json:"sound_count,omitempty"`
	PreferredCommonName string     `json:"preferred_common_name,omitempty"`
	Taxon               inat.Taxon `json:"taxon,omitempty"`
}

func cachedFromResult(r inat.Result) cachedObservation {
	var ofvs []inat.Ofv
	for _, o := range r.Ofvs {
		if o.FieldID == inat.EBirdField || o.FieldID == inat.EBirdScientificNameField {
			ofvs = append(ofvs, inat.Ofv{FieldID: o.FieldID, Value: o.Value})
		}
	}
	return cachedObservation{
		ID:                  r.ID,
		UUID:                r.UUID,
		CreatedAt:           r.CreatedAt,
		Description:         r.Description,
		ObservedOn:          r.ObservedOn,
		Ofvs:                ofvs,
		PhotoCount:          len(r.Photos),
		SoundCount:          len(r.Sounds),
		PreferredCommonName: r.PreferredCommonName,
		Taxon: inat.Taxon{
			Name:                r.Taxon.Name,
			PreferredCommonName: r.Taxon.PreferredCommonName,
		},
	}
}

func (c cachedObservation) toResult() inat.Result {
	var photos []inat.Photo
	if c.PhotoCount > 0 {
		photos = make([]inat.Photo, c.PhotoCount)
	}
	var sounds []inat.Sound
	if c.SoundCount > 0 {
		sounds = make([]inat.Sound, c.SoundCount)
	}
	return inat.Result{
		ID:                  c.ID,
		UUID:                c.UUID,
		CreatedAt:           c.CreatedAt,
		Description:         c.Description,
		ObservedOn:          c.ObservedOn,
		Ofvs:                slices.Clone(c.Ofvs),
		Photos:              photos,
		Sounds:              sounds,
		PreferredCommonName: c.PreferredCommonName,
		Taxon:               c.Taxon,
	}
}

func newEmptyState() *stateFile {
	return &stateFile{
		Version:  stateFileVersion,
		Accounts: map[string]*accountState{},
	}
}

func (s *stateFile) account(userID string) *accountState {
	if s.Accounts == nil {
		s.Accounts = map[string]*accountState{}
	}
	acct, ok := s.Accounts[userID]
	if !ok || acct == nil {
		acct = &accountState{
			Observations: map[string]cachedObservation{},
		}
		s.Accounts[userID] = acct
	}
	if acct.Observations == nil {
		acct.Observations = map[string]cachedObservation{}
	}
	return acct
}

func (a *accountState) pruneRecentRequests(now time.Time) {
	cutoff := now.Add(-24 * time.Hour)
	var kept []time.Time
	for _, t := range a.RecentRequests {
		if !t.Before(cutoff) {
			kept = append(kept, t.UTC())
		}
	}
	a.RecentRequests = kept
}

func (a *accountState) recordRequests(t time.Time, count int) {
	utc := t.UTC()
	for range count {
		a.RecentRequests = append(a.RecentRequests, utc)
	}
}

func (a *accountState) replaceResults(results []inat.Result) {
	a.Observations = make(map[string]cachedObservation, len(results))
	a.mergeResults(results)
}

func (a *accountState) mergeResults(results []inat.Result) {
	if a.Observations == nil {
		a.Observations = map[string]cachedObservation{}
	}
	for _, r := range results {
		a.Observations[r.UUID.String()] = cachedFromResult(r)
	}
}

func (a *accountState) allResults() []inat.Result {
	out := make([]inat.Result, 0, len(a.Observations))
	for _, c := range a.Observations {
		out = append(out, c.toResult())
	}
	slices.SortFunc(out, func(x, y inat.Result) int {
		if c := cmp.Compare(x.ID, y.ID); c != 0 {
			return c
		}
		return cmp.Compare(x.UUID.String(), y.UUID.String())
	})
	return out
}

func (a *accountState) recordCreatedObservation(obs inat.Observation, observedDate string, now time.Time) {
	if a.Observations == nil {
		a.Observations = map[string]cachedObservation{}
	}
	var ofvs []inat.Ofv
	for _, attr := range obs.ObservationFieldValuesAttributes {
		if attr.ObservationFieldID == inat.EBirdField || attr.ObservationFieldID == inat.EBirdScientificNameField {
			if val, ok := attr.Value.(string); ok {
				ofvs = append(ofvs, inat.Ofv{FieldID: attr.ObservationFieldID, Value: val})
			}
		}
	}
	a.Observations[obs.UUID.String()] = cachedObservation{
		UUID:        obs.UUID,
		CreatedAt:   now.UTC().Format(time.RFC3339),
		Description: obs.Description,
		ObservedOn:  observedDate,
		Ofvs:        ofvs,
		Taxon:       inat.Taxon{Name: obs.SpeciesGuess},
	}
}

func (a *accountState) recordUpdatedObservation(u uuid.UUID, desc string, addedPhotos, addedSounds int) {
	if a.Observations == nil {
		a.Observations = map[string]cachedObservation{}
	}
	c := a.Observations[u.String()]
	c.UUID = u
	c.Description = desc
	c.PhotoCount += addedPhotos
	c.SoundCount += addedSounds
	a.Observations[u.String()] = c
}

// loadState reads the state file at path. Missing, unreadable, or corrupt state
// files return a fresh empty state so the run falls back cleanly to a full
// download.
func loadState(path string) *stateFile {
	if path == "" {
		return newEmptyState()
	}
	b, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Warning: couldn't read state file %s (%v); performing full download", path, err)
		}
		return newEmptyState()
	}
	var st stateFile
	if err := json.Unmarshal(b, &st); err != nil || st.Version != stateFileVersion {
		log.Printf("Warning: ignoring invalid or incompatible state file %s; performing full download", path)
		return newEmptyState()
	}
	if st.Accounts == nil {
		st.Accounts = map[string]*accountState{}
	}
	return &st
}

// saveState writes st atomically to path using a temp file in the same
// directory.
func saveState(path string, st *stateFile) error {
	if path == "" || st == nil {
		return nil
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("creating state directory %s: %w", dir, err)
	}
	b, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling state: %w", err)
	}
	tmp, err := os.CreateTemp(dir, "state-*.tmp")
	if err != nil {
		return fmt.Errorf("creating temp state file: %w", err)
	}
	tmpName := tmp.Name()
	renamed := false
	defer func() {
		tmp.Close()
		if !renamed {
			os.Remove(tmpName)
		}
	}()
	if _, err := tmp.Write(b); err != nil {
		return fmt.Errorf("writing temp state file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("closing temp state file: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("renaming state file into place: %w", err)
	}
	renamed = true
	return nil
}
