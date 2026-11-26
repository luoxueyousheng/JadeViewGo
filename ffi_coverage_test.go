package jadeview

import (
	"os"
	"regexp"
	"slices"
	"strings"
	"testing"
)

// This checks symbol coverage, not behavior or exported Go names.
// The documented adapters for C-only lifetimes are asserted explicitly.
func TestBeta4HeaderCoverage(t *testing.T) {
	header, err := os.ReadFile("include/JadeView.h")
	if err != nil {
		t.Fatal(err)
	}
	decl := regexp.MustCompile(`\bJADEVIEW_CALL\s+([a-zA-Z_]\w*)\s*\(`)
	bindings := regexp.MustCompile(`(?:\bnp\("|(?:Dlsym\(handle,\s*|NewProc\()")([a-zA-Z_]\w*)"`)
	seen := map[string]bool{}
	files, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".go") || strings.HasSuffix(f.Name(), "_test.go") {
			continue
		}
		b, err := os.ReadFile(f.Name())
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range bindings.FindAllSubmatch(b, -1) {
			seen[string(m[1])] = true
		}
	}
	var missing []string
	for _, m := range decl.FindAllSubmatch(header, -1) {
		if !seen[string(m[1])] {
			missing = append(missing, string(m[1]))
		}
	}
	slices.Sort(missing)
	// Both are covered by a copied Go TypedEvent with GC ownership; Clone makes an independent copy.
	if !slices.Equal(missing, []string{"event_clone", "event_free"}) {
		t.Fatalf("unexpected unbound ABI: %v", missing)
	}
}
