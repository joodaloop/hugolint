package config

import "testing"

func TestSchemaFor_RootSection(t *testing.T) {
	c := &Config{
		Paths:    Paths{MarkdownRoot: "content"},
		Sections: map[string]map[string]FieldSpec{"root": {"title": {Type: "string"}}},
	}
	key, schema := c.SchemaFor("content/about.md")
	if key != "root" || schema == nil || schema["title"].Type != "string" {
		t.Fatalf("want root match, got key=%q schema=%v", key, schema)
	}
}

func TestSchemaFor_LongestPrefix(t *testing.T) {
	c := &Config{
		Paths: Paths{MarkdownRoot: "content"},
		Sections: map[string]map[string]FieldSpec{
			"writing":         {"a": {Type: "string"}},
			"writing/drafts":  {"b": {Type: "string"}},
		},
	}
	key, schema := c.SchemaFor("content/writing/drafts/x.md")
	if key != "writing/drafts" || schema["b"].Type != "string" {
		t.Fatalf("want writing/drafts, got %q %v", key, schema)
	}
	key, _ = c.SchemaFor("content/writing/x.md")
	if key != "writing" {
		t.Fatalf("want writing, got %q", key)
	}
}

func TestSchemaFor_IndexPages(t *testing.T) {
	c := &Config{
		Paths:         Paths{MarkdownRoot: "content"},
		IndexSections: map[string]map[string]FieldSpec{"writing": {"title": {Type: "string"}}},
		Sections:      map[string]map[string]FieldSpec{"writing": {"author": {Type: "string"}}},
	}
	key, schema := c.SchemaFor("content/writing/_index.md")
	if key != "writing" || schema["title"].Type != "string" {
		t.Fatalf("want index schema, got %q %v", key, schema)
	}
	if _, ok := schema["author"]; ok {
		t.Fatal("index pages should not pull from Sections")
	}
}

func TestSchemaFor_NoMatch(t *testing.T) {
	c := &Config{
		Paths:    Paths{MarkdownRoot: "content"},
		Sections: map[string]map[string]FieldSpec{"writing": {"a": {Type: "string"}}},
	}
	key, schema := c.SchemaFor("content/other/x.md")
	if key != "" || schema != nil {
		t.Fatalf("want no match, got %q %v", key, schema)
	}
}

func TestSkipDir(t *testing.T) {
	c := &Config{Paths: Paths{SkipDirs: []string{"drafts", "private"}}}
	if !c.SkipDir("drafts") || !c.SkipDir("private") {
		t.Fatal("want skip")
	}
	if c.SkipDir("public") {
		t.Fatal("want no skip")
	}
}
