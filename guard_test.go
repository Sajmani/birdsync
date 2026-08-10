package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"html"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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

// TestNoLogFatalInLibraryPackages enforces T-027: log.Fatal belongs in main and
// in tools/, not in a package that has a caller. A library that calls it takes
// the decision to abort away from the program using it — and in this repository
// the program using it is sometimes a test, which the call takes down with it.
//
// A convention like this cannot be checked by example. A test can show that one
// function returns an error; only analysis can show that no function exits.
//
// Verifies: T-027.
func TestNoLogFatalInLibraryPackages(t *testing.T) {
	root, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}

	fset := token.NewFileSet()
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
		name := d.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			return nil
		}
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return err
		}
		if f.Name.Name == "main" {
			return nil // main and tools/ may abort
		}
		scanned++
		rel, _ := filepath.Rel(root, path)
		ast.Inspect(f, func(n ast.Node) bool {
			sel, ok := n.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			pkg, ok := sel.X.(*ast.Ident)
			if !ok {
				return true
			}
			// os.Exit is the same sin under another name.
			bad := (pkg.Name == "log" && strings.HasPrefix(sel.Sel.Name, "Fatal")) ||
				(pkg.Name == "os" && sel.Sel.Name == "Exit")
			if bad {
				t.Errorf("%s:%d calls %s.%s in package %s: library packages return errors "+
					"to their caller (T-027)",
					rel, fset.Position(sel.Pos()).Line, pkg.Name, sel.Sel.Name, f.Name.Name)
			}
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatalf("Walking %s: %v", root, err)
	}
	if scanned == 0 {
		t.Fatal("Scanned no non-main .go files; the check is not doing anything")
	}
}

// TestTranscribedQuotesAppearInSources checks every quotation in a vendored
// source's requirements.md against the documents in the same directory.
//
// A transcription is where an outside rule becomes a requirement, and these
// were produced by an agent reading terms of service. The worst thing such a
// file can do is quote something its source does not say — a paraphrase that
// drifted, a sentence assembled from two places, a passage that was edited
// after being quoted. None of that is visible on the page; all of it is
// visible here.
//
// Matching ignores whitespace, because stripping tags introduces spaces the
// document does not contain ("with you , or"), and splits quotations on an
// ellipsis, because an elided quote is two fragments rather than one string.
//
// Verifies: the accuracy of every `<source>/R#` transcription.
func TestTranscribedQuotesAppearInSources(t *testing.T) {
	tagRE := regexp.MustCompile(`<[^>]+>`)
	ellipsisRE := regexp.MustCompile(`\.\.\.|…`)
	spaceRE := regexp.MustCompile(`\s+`)
	normalize := func(s string) string {
		s = html.UnescapeString(s)
		s = strings.NewReplacer("\u2019", "'", "\u2018", "'", "\u201c", `"`,
			"\u201d", `"`, "\u2014", "-", "\u2013", "-", "\u00a0", " ").Replace(s)
		return spaceRE.ReplaceAllString(s, "")
	}

	dirs, err := filepath.Glob(filepath.Join("spec", "sources", "*"))
	if err != nil {
		t.Fatal(err)
	}
	var checked int
	for _, dir := range dirs {
		req := filepath.Join(dir, "requirements.md")
		if _, err := os.Stat(req); err != nil {
			continue
		}
		docs, _ := filepath.Glob(filepath.Join(dir, "*.html"))
		var corpus strings.Builder
		for _, d := range docs {
			b, err := os.ReadFile(d)
			if err != nil {
				t.Fatal(err)
			}
			corpus.WriteString(normalize(tagRE.ReplaceAllString(string(b), " ")))
		}
		haystack := corpus.String()
		if haystack == "" {
			t.Errorf("%s: no vendored documents to check quotations against", dir)
			continue
		}

		b, err := os.ReadFile(req)
		if err != nil {
			t.Fatal(err)
		}
		var quote []string
		flush := func() {
			defer func() { quote = nil }()
			joined := strings.Join(quote, " ")
			// An attribution line ("— Terms of Use, §1") is not a quotation.
			if i := strings.Index(joined, "\u2014"); i >= 0 {
				joined = joined[:i]
			}
			joined = strings.ReplaceAll(joined, "**", "")
			for _, frag := range ellipsisRE.Split(joined, -1) {
				// Quotation marks and trailing punctuation belong to the
				// citation, not to the quoted text.
				frag = strings.Trim(normalize(frag), `"'.,;:`)
				if len(frag) < 40 { // too short to identify a passage
					continue
				}
				checked++
				if !strings.Contains(haystack, frag) {
					t.Errorf("%s quotes a passage absent from its vendored source:\n  %.90s...", req, frag)
				}
			}
		}
		for _, line := range strings.Split(string(b), "\n") {
			// A bare ">" separates two blockquote paragraphs, which are two
			// quotations rather than one passage. Treating it as a
			// continuation makes the checker hunt for their concatenation,
			// which no source contains.
			if strings.TrimSpace(line) == ">" {
				flush()
				continue
			}
			if after, ok := strings.CutPrefix(line, "> "); ok {
				quote = append(quote, strings.TrimSpace(after))
			} else if len(quote) > 0 {
				flush()
			}
		}
		flush()
	}
	if checked == 0 {
		t.Fatal("Checked no quotations; the check is not doing anything")
	}
	t.Logf("verified %d quoted passages against their vendored sources", checked)
}

// TestTalkLinksResolve checks every github.com link in talks/ against the
// working tree: the file must exist, a line range must be within it, and a
// heading anchor must match a real heading.
//
// A talk narrates the repository, so its links rot as the repository changes —
// and they rot silently, since nobody clicks them until they are on a screen in
// front of an audience. Line numbers are the worst offender: an edit anywhere
// above the target moves it, and the link still resolves, just to the wrong
// place. Anchors survive that, which is why the talk prefers them.
func TestTalkLinksResolve(t *testing.T) {
	const prefix = "https://github.com/Sajmani/birdsync/blob/main/"
	linkRE := regexp.MustCompile(regexp.QuoteMeta(prefix) + `([^)\s]+)`)
	lineRE := regexp.MustCompile(`^L(\d+)(?:-L(\d+))?$`)
	// GitHub's heading slug: lowercase, drop punctuation, each space a hyphen.
	slugPunct := regexp.MustCompile(`[^\w\s-]`)
	slugify := func(h string) string {
		h = strings.ToLower(strings.TrimSpace(h))
		h = strings.NewReplacer("`", "", "*", "", "_", "").Replace(h)
		return strings.ReplaceAll(slugPunct.ReplaceAllString(h, ""), " ", "-")
	}

	// Relative links and images, e.g. an inline slide: ](sdd-example.png)
	relRE := regexp.MustCompile(`\]\(([^)#:\s]+)\)`)

	talks, err := filepath.Glob(filepath.Join("talks", "*.md"))
	if err != nil {
		t.Fatal(err)
	}
	if len(talks) == 0 {
		return
	}
	var checked int
	for _, talk := range talks {
		b, err := os.ReadFile(talk)
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range relRE.FindAllStringSubmatch(string(b), -1) {
			rel := filepath.Join(filepath.Dir(talk), m[1])
			checked++
			if _, err := os.Stat(rel); err != nil {
				t.Errorf("%s links to %s, which does not exist", talk, rel)
			}
		}
		for _, m := range linkRE.FindAllStringSubmatch(string(b), -1) {
			target := strings.TrimRight(m[1], ".,)>")
			path, frag, _ := strings.Cut(target, "#")
			checked++

			content, err := os.ReadFile(path)
			if err != nil {
				t.Errorf("%s links to %s, which does not exist", talk, path)
				continue
			}
			lines := strings.Split(string(content), "\n")

			switch {
			case frag == "":
			case lineRE.MatchString(frag):
				g := lineRE.FindStringSubmatch(frag)
				last := g[1]
				if g[2] != "" {
					last = g[2]
				}
				n, _ := strconv.Atoi(last)
				if n > len(lines) {
					t.Errorf("%s links to %s#%s but the file has %d lines", talk, path, frag, len(lines))
				}
			default:
				var found bool
				for _, line := range lines {
					if strings.HasPrefix(line, "#") && slugify(strings.TrimLeft(line, "# ")) == frag {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("%s links to %s#%s, but no heading there has that anchor", talk, path, frag)
				}
			}
		}
	}
	if checked == 0 {
		t.Fatal("Found no links to check in talks/; the check is not doing anything")
	}
	t.Logf("verified %d links from %d talk(s)", checked, len(talks))
}

// britishSpellings maps a British spelling to its American equivalent.
//
// The regular -ise/-isation families are generated rather than listed, because
// listing them by hand is how a word list ends up with "recognise" and
// "recognised" but not "recognises" — which is exactly what mutation-testing
// this check found. Only the irregulars are enumerated.
//
// Matching is on whole words. An earlier throwaway version matched prefixes and
// flagged "realistic" for "realise". Words genuinely ambiguous between the
// dialects, such as "analyses" — a noun plural in both, a verb only in one —
// are deliberately absent.
var britishSpellings = map[string]string{
	"acknowledgement": "acknowledgment", "amongst": "among", "analyse": "analyze",
	"analysed": "analyzed", "analyses": "analyzes", "analysing": "analyzing",
	"behaviour": "behavior", "behavioural": "behavioral", "behaviours": "behaviors",
	"cancelled": "canceled", "cancelling": "canceling", "catalogue": "catalog",
	"catalogues": "catalogs", "centre": "center", "centred": "centered",
	"centres": "centers", "colour": "color", "coloured": "colored",
	"colouring": "coloring", "colours": "colors", "defence": "defense",
	"fulfil": "fulfill", "grey": "gray", "greyed": "grayed",
	"honour": "honor", "honoured": "honored", "honouring": "honoring",
	"honours": "honors", "judgement": "judgment", "labelled": "labeled",
	"labelling": "labeling", "labour": "labor", "licence": "license",
	"licences": "licenses", "manoeuvre": "maneuver", "modelled": "modeled",
	"modelling": "modeling", "offence": "offense", "practise": "practice",
	"programme": "program", "relabelling": "relabeling", "sceptic": "skeptic",
	"sceptical": "skeptical", "signalled": "signaled", "signalling": "signaling",
	"travelled": "traveled", "travelling": "traveling", "whilst": "while",
}

// iseStems take the regular -ise/-ises/-ised/-ising and -isation endings.
var iseStems = []string{
	"apolog", "author", "categor", "character", "critic", "digit", "emphas",
	"familiar", "final", "formal", "general", "initial", "jeopard", "margin",
	"material", "maxim", "mobil", "modern", "normal", "optim", "organ",
	"parallel", "penal", "priorit", "real", "recogn", "scrutin", "serial",
	"special", "stabil", "standard", "summar", "unrecogn", "util", "visual",
}

func init() {
	for _, stem := range iseStems {
		for brit, amer := range map[string]string{
			"ise": "ize", "ises": "izes", "ised": "ized", "ising": "izing",
			"isation": "ization", "isations": "izations",
		} {
			britishSpellings[stem+brit] = stem + amer
		}
	}
}

// TestAmericanSpellings enforces T-037 and T-038 together, and the second is
// the reason the first is worth a test at all.
//
// Spelling on its own is taste. But a bulk find-and-replace across spec/ would
// rewrite text quoted from somebody else's terms of service — turning a
// verbatim citation into a paraphrase that still looks verbatim, and breaking
// the hashes recorded in each source's PROVENANCE.md. So this check skips two
// things by construction rather than by care: blockquote lines, which are
// quotations, and everything under spec/sources/, which is either a vendored
// document or a transcription of one.
//
// Verifies: T-037, T-038.
func TestAmericanSpellings(t *testing.T) {
	word := regexp.MustCompile(`[A-Za-z]+`)
	var scanned, flagged int

	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch {
			case d.Name() == ".git":
				return filepath.SkipDir
			// Vendored sources and their transcriptions quote publishers
			// verbatim. Their spelling is not ours to correct (T-038).
			case path == filepath.Join("spec", "sources"):
				return filepath.SkipDir
			}
			return nil
		}
		ext := filepath.Ext(d.Name())
		if ext != ".md" && ext != ".go" {
			return nil
		}
		if d.Name() == "guard_test.go" {
			// This file lists the British spellings in order to forbid them,
			// exactly as the live-hostname check names the hosts it bans.
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		scanned++
		for i, line := range strings.Split(string(b), "\n") {
			if strings.HasPrefix(strings.TrimSpace(line), ">") {
				continue // a quotation: reproduce it as published (T-038)
			}
			for _, w := range word.FindAllString(line, -1) {
				if american, bad := britishSpellings[strings.ToLower(w)]; bad {
					flagged++
					t.Errorf("%s:%d: %q should be %q (T-037)", path, i+1, w, american)
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Walking the tree: %v", err)
	}
	if scanned == 0 {
		t.Fatal("Scanned no files; the check is not doing anything")
	}
	if flagged == 0 {
		t.Logf("scanned %d files", scanned)
	}
}
