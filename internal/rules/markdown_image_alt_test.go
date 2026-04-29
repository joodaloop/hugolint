package rules

import "testing"

func TestImageAlt_GenericAlts(t *testing.T) {
	cases := []string{
		"![](/foo.png)\n",
		"![image](/foo.png)\n",
		"![img](/foo.png)\n",
		"![picture](/foo.png)\n",
		"![pic](/foo.png)\n",
		"![photo](/foo.png)\n",
		"![screenshot](/foo.png)\n",
		"![figure](/foo.png)\n",
		"![alt](/foo.png)\n",
		"![alt text](/foo.png)\n",
		"![ Image ](/foo.png)\n", // trimmed, lowercased
	}
	for _, in := range cases {
		diags := markdownImageAlt{}.Check(mdFile(in), nil)
		if len(diags) == 0 {
			t.Errorf("input %q: expected diag", in)
		}
	}
}

func TestImageAlt_DescriptiveAltOK(t *testing.T) {
	diags := markdownImageAlt{}.Check(mdFile("![A black cat sleeping](/cat.png)\n"), nil)
	assertNoDiags(t, diags)
}

func TestImageAlt_MultiplePerLine(t *testing.T) {
	src := "![](/a.png) and ![real text](/b.png) and ![image](/c.png)\n"
	diags := markdownImageAlt{}.Check(mdFile(src), nil)
	if len(diags) != 2 {
		t.Fatalf("want 2 diags, got %d: %v", len(diags), messages(diags))
	}
}

func TestImageAlt_LineNumber(t *testing.T) {
	src := "para\n\n![](/x.png)\n"
	diags := markdownImageAlt{}.Check(mdFile(src), nil)
	if len(diags) != 1 || diags[0].Line != 3 {
		t.Fatalf("want line 3, got %+v", diags)
	}
}

func TestImageAlt_ID(t *testing.T) {
	if (markdownImageAlt{}).ID() != "image-alt" {
		t.Fatal("wrong ID")
	}
}
