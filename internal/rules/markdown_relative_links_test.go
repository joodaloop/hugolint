package rules

import "testing"

func TestRelativeLinks_RelativePathFlagged(t *testing.T) {
	src := "see [foo](foo/bar.md) please\n"
	diags := markdownRelativeLinks{}.Check(mdFile(src), nil)
	if len(diags) != 1 {
		t.Fatalf("want 1 diag, got %v", messages(diags))
	}
}

func TestRelativeLinks_RootRelativeOK(t *testing.T) {
	diags := markdownRelativeLinks{}.Check(mdFile("[foo](/foo/bar)\n"), nil)
	assertNoDiags(t, diags)
}

func TestRelativeLinks_FragmentOK(t *testing.T) {
	diags := markdownRelativeLinks{}.Check(mdFile("[foo](#anchor)\n"), nil)
	assertNoDiags(t, diags)
}

func TestRelativeLinks_AbsoluteSchemeOK(t *testing.T) {
	for _, in := range []string{
		"[foo](https://example.com)\n",
		"[foo](http://example.com)\n",
		"[foo](mailto:x@y.com)\n",
	} {
		diags := markdownRelativeLinks{}.Check(mdFile(in), nil)
		if len(diags) != 0 {
			t.Errorf("input %q should not flag, got %v", in, messages(diags))
		}
	}
}

func TestRelativeLinks_DotPathFlagged(t *testing.T) {
	diags := markdownRelativeLinks{}.Check(mdFile("[foo](./bar.md)\n"), nil)
	if len(diags) != 1 {
		t.Fatalf("want 1 diag, got %v", messages(diags))
	}
}

func TestRelativeLinks_ID(t *testing.T) {
	if (markdownRelativeLinks{}).ID() != "relative-link" {
		t.Fatal("wrong ID")
	}
}
