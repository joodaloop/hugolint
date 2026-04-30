package rules

import (
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
)

func init() {
	RegisterMarkdownAST(&markdownImageAlt{})
}

type markdownImageAlt struct{}

func (markdownImageAlt) ID() string { return "image-alt" }

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
	ast.Walk(f.AST, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		img, ok := n.(*ast.Image)
		if !ok {
			return ast.WalkContinue, nil
		}
		raw := imageAltText(img, f.Body)
		alt := strings.ToLower(strings.TrimSpace(raw))
		if genericAlts[alt] {
			diags = append(diags, Diagnostic{
				Path: f.Path, Line: f.NodeLine(img), Rule: "image-alt",
				Message: fmt.Sprintf("useless image alt text: %q", raw),
			})
		}
		return ast.WalkSkipChildren, nil
	})
	return diags
}

func imageAltText(img *ast.Image, body []byte) string {
	var b strings.Builder
	ast.Walk(img, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if t, ok := n.(*ast.Text); ok {
			b.Write(t.Segment.Value(body))
		}
		return ast.WalkContinue, nil
	})
	return b.String()
}
