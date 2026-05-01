package rules

import "testing"

func TestHTMLArtifacts_LeakedPatterns(t *testing.T) {
	cases := []struct {
		text, want string
	}{
		{"see (http://x.com", "leaked '(http'"},
		{"see )http", "leaked ')http'"},
		{"text [http://x.com text", "leaked '[http'"},
		{"]http something", "leaked ']http'"},
		{"<!-- raw -->", "literal '<!--'"},
		{"end -->", "literal '-->'"},
		{"<-- arrow", "literal '<--'"},
		{"<— arrow", "literal '<—'"},
		{"arrow —>", "literal '—>'"},
		{"<q>quote</q>", "literal '<q>'"},
		{"stray </q> only", "literal '</q>'"},
		{"</q< broken", "literal '</q<'"},
		{"<del>strike</del>", "literal '<del>'"},
		{"unparsed ** bold", "unparsed bold"},
		{"stray /* comment", "stray code-comment marker"},
		{"stray */ end", "stray code-comment marker"},
		{"shortcode {{<", "shortcode delimiter '{{<'"},
		{"shortcode >}}", "shortcode delimiter '>}}'"},
		{"shortcode {{%", "shortcode delimiter '{{%'"},
		{"shortcode %}}", "shortcode delimiter '%}}'"},
	}
	for _, tc := range cases {
		diags := htmlArtifacts{}.Check(htmlFile(tc.text), nil)
		if !containsMsg(diags, tc.want) {
			t.Errorf("text %q: want %q, got %v", tc.text, tc.want, messages(diags))
		}
	}
}

func TestHTMLArtifacts_EmptyText(t *testing.T) {
	diags := htmlArtifacts{}.Check(&HTMLFile{Path: "test.html"}, nil)
	assertNoDiags(t, diags)
}

func TestHTMLArtifacts_CleanText(t *testing.T) {
	diags := htmlArtifacts{}.Check(htmlFile("Just plain prose with no artifacts."), nil)
	assertNoDiags(t, diags)
}

func TestHTMLArtifacts_ID(t *testing.T) {
	if (htmlArtifacts{}).ID() != "rendered-artifacts" {
		t.Fatal("wrong ID")
	}
}
