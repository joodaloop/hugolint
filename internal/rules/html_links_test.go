package rules

import "testing"

func htmlCtx(pages map[string]bool, ids map[string]map[string]int) *HTMLContext {
	return &HTMLContext{Root: "/site/public", Pages: pages, PageIDs: ids}
}

func TestHTMLLinks_RelativeMissing(t *testing.T) {
	f := &HTMLFile{Path: "/site/public/foo/index.html", URLPath: "/foo/", Links: []string{"/missing/"}}
	ctx := htmlCtx(map[string]bool{"/foo/": true}, nil)
	diags := relativeLinks{}.Check(f, ctx)
	if !containsMsg(diags, "/missing/") {
		t.Fatalf("want missing-link diag, got %v", messages(diags))
	}
}

func TestHTMLLinks_RelativeExists(t *testing.T) {
	f := &HTMLFile{Path: "/site/public/foo/index.html", URLPath: "/foo/", Links: []string{"/bar/"}}
	ctx := htmlCtx(map[string]bool{"/foo/": true, "/bar/": true}, nil)
	diags := relativeLinks{}.Check(f, ctx)
	assertNoDiags(t, diags)
}

func TestHTMLLinks_AbsoluteSkipped(t *testing.T) {
	f := &HTMLFile{Path: "/site/public/x.html", URLPath: "/x.html", Links: []string{"https://x.com", "mailto:a@b", "//cdn.x.com/y"}}
	diags := relativeLinks{}.Check(f, &HTMLContext{Pages: map[string]bool{}})
	assertNoDiags(t, diags)
}

func TestIsRelative(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"", false},
		{"#frag", false},
		{"//cdn.x.com", false},
		{"mailto:a@b.com", false},
		{"javascript:void(0)", false},
		{"data:text/plain,hi", false},
		{"https://x.com", false},
		{"/foo/bar", true},
		{"foo/bar", true},
	}
	for _, tc := range cases {
		if got := isRelative(tc.in); got != tc.want {
			t.Errorf("isRelative(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestResolve(t *testing.T) {
	cases := []struct {
		page, href, want string
		ok               bool
	}{
		{"/foo/", "/bar", "/bar", true},
		{"/foo/", "bar/", "/foo/bar/", true},
		{"/foo/index.html", "../bar", "/foo/bar", true},
		{"/foo/", "", "", false},
	}
	for _, tc := range cases {
		got, ok := resolve(tc.page, tc.href)
		if ok != tc.ok || (ok && got != tc.want) {
			t.Errorf("resolve(%q,%q) = (%q,%v), want (%q,%v)", tc.page, tc.href, got, ok, tc.want, tc.ok)
		}
	}
}
