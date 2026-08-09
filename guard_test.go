package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNoLiveServiceHostsInTests is a static check over this repository's own
// source. T-010 forbids tests from contacting a live service, and the failure
// mode it guards against is quiet: a test that reaches api.inaturalist.org
// still passes on the author's machine, then starts writing to whichever
// account the CI environment's credentials point at.
//
// The check is a substring scan, so it catches the realistic mistake — pasting
// a real URL into a test — and not a determined one, such as assembling the
// host from parts or reading it from the environment. That limit is recorded
// in acceptance.md rather than papered over; the alternative, parsing every
// expression that could produce a hostname, would cost more than it catches.
//
// Verifies: T-010.
func TestNoLiveServiceHostsInTests(t *testing.T) {
	// The real hosts, spelled in pieces so that this file does not trip its
	// own check. Keep in sync with inat.BaseURL and ebird.macaulayBaseURL.
	forbidden := []string{
		"api." + "inaturalist.org",
		"cdn.download.ams." + "birds.cornell.edu",
	}

	root, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}

	var scanned int
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(d.Name(), "_test.go") {
			return nil
		}
		if path == filepath.Join(root, "guard_test.go") {
			return nil // this file names the hosts in order to forbid them
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		scanned++
		for _, host := range forbidden {
			if strings.Contains(string(b), host) {
				rel, _ := filepath.Rel(root, path)
				t.Errorf("%s names the live service host %q; tests must use an "+
					"httptest server via the base-URL seams instead (T-010)", rel, host)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Walking %s: %v", root, err)
	}

	// A scan that silently matched nothing would pass forever. Fail loudly if
	// the walk stops finding test files at all.
	if scanned == 0 {
		t.Fatal("Scanned no _test.go files; the check is not doing anything")
	}
}
