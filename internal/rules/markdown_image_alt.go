package rules

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

func init() {
	RegisterMarkdown(&markdownImageAlt{})
}

type markdownImageAlt struct{}

func (markdownImageAlt) ID() string { return "image-alt" }

var mdImage = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)

var genericAlts = map[string]bool{
	"":           true,
	"image":      true,
	"img":        true,
	"picture":    true,
	"pic":        true,
	"photo":      true,
	"screenshot": true,
	"figure":     true,
	"alt":        true,
	"alt text":   true,
}

func (markdownImageAlt) Check(f *MarkdownFile, _ *MarkdownContext) []Diagnostic {
	var diags []Diagnostic
	scanner := bufio.NewScanner(bytes.NewReader(f.Content))
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	line := 0
	for scanner.Scan() {
		line++
		for _, m := range mdImage.FindAllStringSubmatch(scanner.Text(), -1) {
			alt := strings.ToLower(strings.TrimSpace(m[1]))
			if genericAlts[alt] {
				diags = append(diags, Diagnostic{
					Path: f.Path, Line: line, Rule: "image-alt",
					Message: fmt.Sprintf("useless image alt text: %q", m[1]),
				})
			}
		}
	}
	return diags
}
