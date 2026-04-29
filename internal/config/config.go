package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Paths         Paths                           `yaml:"paths"`
	Sections      map[string]map[string]FieldSpec `yaml:"sections"`
	IndexSections map[string]map[string]FieldSpec `yaml:"index_pages"`
	Links         Links                           `yaml:"links"`
	Spelling      Spelling                        `yaml:"spelling"`
}

type Links struct {
	SiteHosts []string `yaml:"site_hosts"`
}

type Spelling struct {
	Dict string `yaml:"dict"`
}

type Paths struct {
	MarkdownRoot string   `yaml:"markdown_root"`
	BuildRoot    string   `yaml:"build_root"`
	SkipDirs     []string `yaml:"skip_dirs"`
}

type FieldSpec struct {
	Type     string   `yaml:"type"`
	Required bool     `yaml:"required"`
	Values   []string `yaml:"values"`
	Items    string   `yaml:"items"`
	Min      any      `yaml:"min"`
	Max      any      `yaml:"max"`
}

func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if c.Paths.MarkdownRoot == "" {
		c.Paths.MarkdownRoot = "content"
	}
	if c.Paths.BuildRoot == "" {
		c.Paths.BuildRoot = "public"
	}
	return &c, nil
}

// SchemaFor returns the schema to apply to a markdown file given its path.
// For files named _index.md it consults index_pages; otherwise sections.
// Section match is longest-prefix relative to MarkdownRoot. Files directly
// under MarkdownRoot use the special section key "root".
func (c *Config) SchemaFor(filePath string) (string, map[string]FieldSpec) {
	table := c.Sections
	if filepath.Base(filePath) == "_index.md" {
		table = c.IndexSections
	}
	rel, err := filepath.Rel(c.Paths.MarkdownRoot, filePath)
	if err != nil {
		return "", nil
	}
	rel = filepath.ToSlash(rel)
	if !strings.Contains(rel, "/") {
		if schema, ok := table["root"]; ok {
			return "root", schema
		}
	}
	var bestKey string
	for key := range table {
		if key == "root" {
			continue
		}
		if rel == key || strings.HasPrefix(rel, key+"/") {
			if len(key) > len(bestKey) {
				bestKey = key
			}
		}
	}
	if bestKey == "" {
		return "", nil
	}
	return bestKey, table[bestKey]
}

// SkipDir reports whether a directory name should be skipped during traversal.
func (c *Config) SkipDir(name string) bool {
	for _, d := range c.Paths.SkipDirs {
		if d == name {
			return true
		}
	}
	return false
}
