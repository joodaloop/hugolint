package rules

import "testing"

func TestFragmentLinks_MissingFragmentSamePage(t *testing.T) {
	f := &HTMLFile{
		Path: "/site/public/x.html", URLPath: "/x.html",
		Links: []string{"#missing"},
	}
	ctx := htmlCtx(
		map[string]bool{"/x.html": true},
		map[string]map[string]int{"/x.html": {"present": 1}},
	)
	diags := fragmentLinks{}.Check(f, ctx)
	if !containsMsg(diags, "fragment #missing not found") {
		t.Fatalf("want missing-fragment, got %v", messages(diags))
	}
}

func TestFragmentLinks_PresentFragmentSamePage(t *testing.T) {
	f := &HTMLFile{
		Path: "/site/public/x.html", URLPath: "/x.html",
		Links: []string{"#here"},
	}
	ctx := htmlCtx(
		map[string]bool{"/x.html": true},
		map[string]map[string]int{"/x.html": {"here": 1}},
	)
	assertNoDiags(t, fragmentLinks{}.Check(f, ctx))
}

func TestFragmentLinks_DuplicateFragment(t *testing.T) {
	f := &HTMLFile{
		Path: "/site/public/x.html", URLPath: "/x.html",
		Links: []string{"#dup"},
	}
	ctx := htmlCtx(
		map[string]bool{"/x.html": true},
		map[string]map[string]int{"/x.html": {"dup": 2}},
	)
	diags := fragmentLinks{}.Check(f, ctx)
	if !containsMsg(diags, "appears 2 times") {
		t.Fatalf("want duplicate diag, got %v", messages(diags))
	}
}

func TestFragmentLinks_CrossPage(t *testing.T) {
	f := &HTMLFile{
		Path: "/site/public/a.html", URLPath: "/a.html",
		Links: []string{"/b.html#sec"},
	}
	ctx := htmlCtx(
		map[string]bool{"/a.html": true, "/b.html": true},
		map[string]map[string]int{"/b.html": {"sec": 1}},
	)
	assertNoDiags(t, fragmentLinks{}.Check(f, ctx))
}

func TestFragmentLinks_NoFragmentSkipped(t *testing.T) {
	f := &HTMLFile{
		Path: "/site/public/x.html", URLPath: "/x.html",
		Links: []string{"/foo"},
	}
	ctx := htmlCtx(map[string]bool{}, nil)
	assertNoDiags(t, fragmentLinks{}.Check(f, ctx))
}

func TestFragmentLinks_AbsoluteSkipped(t *testing.T) {
	f := &HTMLFile{
		Path: "/site/public/x.html", URLPath: "/x.html",
		Links: []string{"https://x.com/y#sec", "mailto:a@b#anchor"},
	}
	ctx := htmlCtx(map[string]bool{}, nil)
	assertNoDiags(t, fragmentLinks{}.Check(f, ctx))
}
