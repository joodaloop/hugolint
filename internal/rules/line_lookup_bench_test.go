package rules

import (
	"testing"

	"github.com/yuin/goldmark/ast"
)

func benchHeadingNode(f *MarkdownFile) ast.Node {
	var found ast.Node
	ast.Walk(f.AST, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering || found != nil {
			return ast.WalkContinue, nil
		}
		if n.Kind() == ast.KindHeading {
			found = n
			return ast.WalkStop, nil
		}
		return ast.WalkContinue, nil
	})
	return found
}

func BenchmarkLineAt(b *testing.B) {
	mf := benchMarkdownFile(2000)
	offset := len(mf.Body) * 3 / 4
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mf.LineAt(offset)
	}
}

func BenchmarkNodeLine(b *testing.B) {
	mf := benchMarkdownFile(2000)
	n := benchHeadingNode(mf)
	if n == nil {
		b.Fatal("no heading node found")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mf.NodeLine(n)
	}
}
