package main

import (
	"iter"
	"time"

	"github.com/Sajmani/birdsync/ebird"
	"github.com/Sajmani/birdsync/inat"
)

// dateTimeFlag validates a command line flag containing a date or a date & time.
type dateTimeFlag struct {
	t time.Time
}

func (f *dateTimeFlag) String() string {
	return f.t.Format(time.DateTime)
}

func (f *dateTimeFlag) Set(s string) error {
	t, err := time.Parse(time.DateTime, s)
	if err != nil {
		t, err = time.Parse(time.DateOnly, s)
		if err != nil {
			return err
		}
	}
	f.t = t
	return nil
}

func (f *dateTimeFlag) Time() time.Time {
	return f.t
}

// ebirdClient encapsulates the ebird package functions for testing.
type ebirdClient interface {
	Records(string) (iter.Seq[ebird.Record], error)
	DownloadMLAsset(string) (string, bool, error)
}

type ebirdClientImpl struct{}

func (ebirdClientImpl) Records(path string) (iter.Seq[ebird.Record], error) {
	return ebird.Records(path)
}

func (ebirdClientImpl) DownloadMLAsset(id string) (string, bool, error) {
	return ebird.DownloadMLAsset(id)
}

// inatClient encapsulates the inat package functions for testing.
type inatClient interface {
	GetUserID() string
	GetAPIToken() string
	DownloadObservations(string, time.Time, time.Time, ...string) []inat.Result
	CreateObservation(inat.Observation) error
	UpdateObservation(inat.Observation) error
	UploadMedia(string, bool, string, string) error
}

type inatClientImpl struct {
	client *inat.Client
}

func (c inatClientImpl) GetUserID() string {
	return inat.GetUserID()
}

func (c inatClientImpl) GetAPIToken() string {
	return inat.GetAPIToken()
}

func (c inatClientImpl) DownloadObservations(userID string, after, before time.Time, fields ...string) []inat.Result {
	return c.client.DownloadObservations(userID, after, before, fields...)
}

func (c inatClientImpl) CreateObservation(obs inat.Observation) error {
	return c.client.CreateObservation(obs)
}

func (c inatClientImpl) UpdateObservation(obs inat.Observation) error {
	return c.client.UpdateObservation(obs)
}

func (c inatClientImpl) UploadMedia(filename string, isPhoto bool, assetID, obsUUID string) error {
	return c.client.UploadMedia(filename, isPhoto, assetID, obsUUID)
}
