package rules

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCollectJSFiles_OnlyJS(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "foo.js", 100)
	writeFile(t, dir, "bar.css", 200)
	writeFile(t, dir, "baz.html", 300)

	files := []BuiltFile{
		{Path: filepath.Join(dir, "foo.js"), URLPath: "/foo.js"},
		{Path: filepath.Join(dir, "bar.css"), URLPath: "/bar.css"},
		{Path: filepath.Join(dir, "baz.html"), URLPath: "/baz.html"},
	}
	ctx := &HTMLContext{LinkedPages: map[string]bool{"/foo.js": true}}

	got := collectJSFiles(files, ctx)
	if len(got) != 1 {
		t.Fatalf("want 1 JS file, got %d", len(got))
	}
	if got[0].urlPath != "/foo.js" {
		t.Errorf("urlPath = %q, want /foo.js", got[0].urlPath)
	}
}

func TestCollectJSFiles_OrphanExcluded(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "linked.js", 100)
	writeFile(t, dir, "orphan.js", 200)

	files := []BuiltFile{
		{Path: filepath.Join(dir, "linked.js"), URLPath: "/linked.js"},
		{Path: filepath.Join(dir, "orphan.js"), URLPath: "/orphan.js"},
	}
	ctx := &HTMLContext{LinkedPages: map[string]bool{"/linked.js": true}}

	got := collectJSFiles(files, ctx)
	if len(got) != 1 {
		t.Fatalf("want 1 non-orphaned JS, got %d", len(got))
	}
	if got[0].urlPath != "/linked.js" {
		t.Errorf("urlPath = %q, want /linked.js", got[0].urlPath)
	}
}

func TestCollectJSFiles_WellKnownIncluded(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "sw.js", 50)

	files := []BuiltFile{
		{Path: filepath.Join(dir, "sw.js"), URLPath: "/sw.js"},
	}
	ctx := &HTMLContext{LinkedPages: map[string]bool{}}

	got := collectJSFiles(files, ctx)
	if len(got) != 1 {
		t.Fatalf("want 1 (well-known) JS, got %d", len(got))
	}
	if got[0].urlPath != "/sw.js" {
		t.Errorf("urlPath = %q, want /sw.js", got[0].urlPath)
	}
}

func TestCollectJSFiles_SortedBySize(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "small.js", 100)
	writeFile(t, dir, "big.js", 500)
	writeFile(t, dir, "mid.js", 200)

	files := []BuiltFile{
		{Path: filepath.Join(dir, "small.js"), URLPath: "/small.js"},
		{Path: filepath.Join(dir, "big.js"), URLPath: "/big.js"},
		{Path: filepath.Join(dir, "mid.js"), URLPath: "/mid.js"},
	}
	ctx := &HTMLContext{LinkedPages: map[string]bool{
		"/small.js": true, "/big.js": true, "/mid.js": true,
	}}

	got := collectJSFiles(files, ctx)
	if len(got) != 3 {
		t.Fatalf("want 3 JS files, got %d", len(got))
	}
	if got[0].urlPath != "/big.js" {
		t.Errorf("first = %q, want /big.js", got[0].urlPath)
	}
	if got[1].urlPath != "/mid.js" {
		t.Errorf("second = %q, want /mid.js", got[1].urlPath)
	}
	if got[2].urlPath != "/small.js" {
		t.Errorf("third = %q, want /small.js", got[2].urlPath)
	}
}

func TestCollectJSFiles_EmptyWhenNoJS(t *testing.T) {
	files := []BuiltFile{
		{Path: "/site/public/index.html", URLPath: "/"},
	}
	ctx := &HTMLContext{LinkedPages: map[string]bool{}}
	got := collectJSFiles(files, ctx)
	if len(got) != 0 {
		t.Fatalf("want 0 JS files, got %d", len(got))
	}
}

func TestCollectJSFiles_LinkedViaAlias(t *testing.T) {
	dir := t.TempDir()
	// JS files in subdirs may end with "/" in URLPath, triggering aliases
	writeFile(t, dir, "robots.txt", 100)

	files := []BuiltFile{
		{Path: filepath.Join(dir, "robots.txt"), URLPath: "/robots.txt"},
	}
	ctx := &HTMLContext{LinkedPages: map[string]bool{}}
	got := collectJSFiles(files, ctx)
	if len(got) != 0 {
		t.Fatalf("want 0 (not JS), got %d", len(got))
	}
}

func writeFile(t *testing.T, dir, name string, size int) {
	t.Helper()
	content := make([]byte, size)
	if err := os.WriteFile(filepath.Join(dir, name), content, 0644); err != nil {
		t.Fatal(err)
	}
}
