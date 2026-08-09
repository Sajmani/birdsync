package main

import (
	"go/ast"
	"go/parser"
	"go/token"
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

// TestToolsAreReadOnly enforces that nothing in tools/ can modify a user's
// account. The directory used to hold six programs that created, updated, and
// deleted observations, guarded only by a `debug` constant that one of them
// shipped with turned off, and by a rule in AGENTS.md telling people not to run
// them. A rule nobody can check is a rule that eventually gets broken, so the
// mutating tools were deleted and this check keeps them from coming back.
//
// It parses rather than greps, so a mention in a comment or a string doesn't
// trip it, and a renamed import doesn't evade it: the test looks for calls to
// the mutating methods by selector name, which survives aliasing of the inat
// package.
//
// Verifies: T-032.
func TestToolsAreReadOnly(t *testing.T) {
	mutating := map[string]bool{
		"CreateObservation": true,
		"UpdateObservation": true,
		"DeleteObservation": true,
		"UploadMedia":       true,
	}

	root, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	toolsDir := filepath.Join(root, "tools")
	if _, err := os.Stat(toolsDir); os.IsNotExist(err) {
		return // no tools to check
	}

	fset := token.NewFileSet()
	var scanned int
	err = filepath.WalkDir(toolsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".go") {
			return nil
		}
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return err
		}
		scanned++
		rel, _ := filepath.Rel(root, path)
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			if mutating[sel.Sel.Name] {
				t.Errorf("%s:%d calls %s: programs in tools/ must be read-only (T-032)",
					rel, fset.Position(sel.Pos()).Line, sel.Sel.Name)
			}
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatalf("Walking %s: %v", toolsDir, err)
	}
	if scanned == 0 {
		t.Fatal("Scanned no .go files under tools/; the check is not doing anything")
	}
}
