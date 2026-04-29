package rules

import (
	"strings"
	"testing"
)

func TestProseHygiene_RepeatedWord(t *testing.T) {
	diags := markdownProseHygiene{}.Check(mdFile("the the cat\n"), nil)
	if len(diags) != 1 || !strings.Contains(diags[0].Message, `"the the"`) {
		t.Fatalf("want one repeated-word diag, got %v", messages(diags))
	}
	if diags[0].Line != 1 {
		t.Fatalf("want line 1, got %d", diags[0].Line)
	}
}

func TestProseHygiene_RepeatedWordCaseInsensitive(t *testing.T) {
	diags := markdownProseHygiene{}.Check(mdFile("The the cat\n"), nil)
	if len(diags) != 1 {
		t.Fatalf("want one diag, got %v", messages(diags))
	}
}

func TestProseHygiene_LiteralPatterns(t *testing.T) {
	cases := []struct {
		in, contains string
	}{
		{"hello —— world\n", "double em dash"},
		{"foo --- bar\n", "literal triple hyphen"},
		{"it's '' a thing\n", "double apostrophe"},
		{"`` two\n", "double backtick"},
		{"text (www.example.com)\n", "missing space/scheme before www"},
		{"text )www.example.com\n", "missing space after closing paren"},
		{"hi (there )\n", "space before closing paren"},
		{"empty []()\n", "empty link"},
		{"empty ![]()\n", "empty image"},
		{"link [foo](//x.com)\n", "protocol-relative link"},
		{"link [foo](/wiki/bar)\n", "wrong wiki path"},
		{"link [foo](wiki/bar)\n", "wrong wiki path"},
		{"img [foo](image/x.png)\n", "wrong image path"},
		{"img [foo](images/x.png)\n", "wrong image path"},
		{"img [foo](/image/x.png)\n", "wrong image path"},
		{"img [foo](/images/x.png)\n", "wrong image path"},
		{` " ](url)` + "\n", "quote glued to link"},
	}
	for _, tc := range cases {
		diags := markdownProseHygiene{}.Check(mdFile(tc.in), nil)
		if !containsMsg(diags, tc.contains) {
			t.Errorf("input %q: missing %q in %v", tc.in, tc.contains, messages(diags))
		}
	}
}

func TestProseHygiene_HRLineDoesNotTrigger(t *testing.T) {
	diags := markdownProseHygiene{}.Check(mdFile("text\n\n---\n\nmore\n"), nil)
	if containsMsg(diags, "literal triple hyphen") {
		t.Fatalf("HR line should not trigger triple-hyphen warning: %v", messages(diags))
	}
}

func TestProseHygiene_FrontmatterSkipped(t *testing.T) {
	src := "---\ntitle: \"hello --- world\"\nthe the: bad\n---\n\nbody\n"
	diags := markdownProseHygiene{}.Check(mdFile(src), nil)
	assertNoDiags(t, diags)
}

func TestProseHygiene_FencedCodeSkipped(t *testing.T) {
	src := "intro\n\n```go\nthe the\n---\n```\n\nafter\n"
	diags := markdownProseHygiene{}.Check(mdFile(src), nil)
	assertNoDiags(t, diags)
}

func TestProseHygiene_StyleScriptBlocksSkipped(t *testing.T) {
	src := "<style>\nthe the\n</style>\nbody\n<script>\nfoo foo\n</script>\n"
	diags := markdownProseHygiene{}.Check(mdFile(src), nil)
	assertNoDiags(t, diags)
}

func TestProseHygiene_LinksAndCodeNotTokenizedAsRepeats(t *testing.T) {
	src := "before [foo](/foo) middle `the the` after <em>x</em>\n"
	diags := markdownProseHygiene{}.Check(mdFile(src), nil)
	if containsMsg(diags, "repeated word") {
		t.Fatalf("inline code/links should not produce repeats: %v", messages(diags))
	}
}

func TestProseHygiene_SpacedColon(t *testing.T) {
	diags := markdownProseHygiene{}.Check(mdFile("note : here\n"), nil)
	if !containsMsg(diags, "spaced colon") {
		t.Fatalf("expected spaced colon: %v", messages(diags))
	}
}

func TestProseHygiene_PlusMinus(t *testing.T) {
	for _, in := range []string{"value +-3\n", "value -+3\n"} {
		diags := markdownProseHygiene{}.Check(mdFile(in), nil)
		if !containsMsg(diags, "malformed plus-minus") {
			t.Errorf("input %q: expected plus-minus diagnostic, got %v", in, messages(diags))
		}
	}
}

func TestProseHygiene_StrayAfterQuoteParen(t *testing.T) {
	cases := []string{
		`hello (foo "), more` + "\n",
		`hello (foo "); more` + "\n",
		`hello (foo "),` + "\n",
		`hello (foo ")` + "\n",
	}
	for _, in := range cases {
		diags := markdownProseHygiene{}.Check(mdFile(in), nil)
		if !containsMsg(diags, "closing paren attached to stray quote/punct") {
			t.Errorf("input %q: missing stray-after-quote diag, got %v", in, messages(diags))
		}
	}
}

func TestProseHygiene_LineNumbersAccurate(t *testing.T) {
	src := "line one\nthe the line two\nline three\nfoo foo line four\n"
	diags := markdownProseHygiene{}.Check(mdFile(src), nil)
	got := linesOf(diags)
	want := []int{2, 4}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("want lines %v, got %v (msgs %v)", want, got, messages(diags))
	}
}

func TestProseHygiene_ID(t *testing.T) {
	if (markdownProseHygiene{}).ID() != "prose-hygiene" {
		t.Fatal("wrong ID")
	}
}
