package rules

import (
	"bufio"
	"bytes"
	"fmt"
	"net/url"
	"regexp"
)

func init() {
	RegisterMarkdown(&markdownRelativeLinks{})
}

type markdownRelativeLinks struct{}

func (markdownRelativeLinks) ID() string { return "relative-link" }

var mdLinkOrImage = regexp.MustCompile(`!?\]\(([^)]+)\)`)

func (markdownRelativeLinks) Check(f *MarkdownFile, _ *MarkdownContext) []Diagnostic {
	var diags []Diagnostic
	scanner := bufio.NewScanner(bytes.NewReader(f.Content))
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	line := 0
	for scanner.Scan() {
		line++
		for _, m := range mdLinkOrImage.FindAllStringSubmatch(scanner.Text(), -1) {
			raw := stripTitle(m[1])
			if raw == "" {
				continue
			}
			if raw[0] == '/' || raw[0] == '#' {
				continue
			}
			u, err := url.Parse(raw)
			if err != nil {
				continue
			}
			if u.Scheme != "" {
				continue
			}
			diags = append(diags, Diagnostic{
				Path: f.Path, Line: line, Rule: "relative-link",
				Message: fmt.Sprintf("relative link: %s (use root-relative path starting with /)", raw),
			})
		}
	}
	return diags
}
