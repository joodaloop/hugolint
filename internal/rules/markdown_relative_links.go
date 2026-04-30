package rules

import (
	"fmt"
	"net/url"

	"github.com/yuin/goldmark/ast"
)

func init() {
	RegisterMarkdownAST(&markdownRelativeLinks{})
}

type markdownRelativeLinks struct{}

func (markdownRelativeLinks) ID() string { return "relative-link" }

func (markdownRelativeLinks) Check(f *MarkdownFile, _ *MarkdownContext) []Diagnostic {
	var diags []Diagnostic
	ast.Walk(f.AST, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		raw, ok := linkDestination(n)
		if !ok || raw == "" {
			return ast.WalkContinue, nil
		}
		if raw[0] == '/' || raw[0] == '#' {
			return ast.WalkContinue, nil
		}
		u, err := url.Parse(raw)
		if err != nil || u.Scheme != "" {
			return ast.WalkContinue, nil
		}
		diags = append(diags, Diagnostic{
			Path: f.Path, Line: f.NodeLine(n), Rule: "relative-link",
			Message: fmt.Sprintf("relative link: %s (use root-relative path starting with /)", raw),
		})
		return ast.WalkContinue, nil
	})
	return diags
}

// linkDestination returns the destination URL for a Link or Image node.
func linkDestination(n ast.Node) (string, bool) {
	switch v := n.(type) {
	case *ast.Link:
		return string(v.Destination), true
	case *ast.Image:
		return string(v.Destination), true
	}
	return "", false
}
