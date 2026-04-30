package rules

import (
	"github.com/yuin/goldmark/ast"
)

func init() {
	RegisterMarkdownAST(&markdownHeadingsAST{})
}

// markdownHeadingsAST flags any H1 (ATX or setext) using the parsed AST.
// CommonMark normalizes both `# Title` and `Title\n=====` into the same
// *ast.Heading{Level:1}, so this is one branch instead of two.
type markdownHeadingsAST struct{}

func (markdownHeadingsAST) ID() string { return "headings" }

func (markdownHeadingsAST) Check(f *MarkdownFile, _ *MarkdownContext) []Diagnostic {
	if f.AST == nil {
		return nil
	}
	var diags []Diagnostic
	_ = ast.Walk(f.AST, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		h, ok := n.(*ast.Heading)
		if !ok {
			return ast.WalkContinue, nil
		}
		if h.Level == 1 {
			diags = append(diags, Diagnostic{
				Path: f.Path, Line: f.NodeLine(h), Rule: "headings",
				Message: "h1 headings are not allowed",
			})
		}
		return ast.WalkContinue, nil
	})
	return diags
}
