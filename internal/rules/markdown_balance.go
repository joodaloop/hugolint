package rules

import (
	"bufio"
	"bytes"
	"fmt"
)

func init() {
	RegisterMarkdown(&markdownBalance{})
}

type markdownBalance struct{}

func (markdownBalance) ID() string { return "balance" }

type openDelim struct {
	ch   byte
	line int
}

func (markdownBalance) Check(f *MarkdownFile, _ *MarkdownContext) []Diagnostic {
	var diags []Diagnostic
	content := stripFrontmatter(f.Content)
	scanner := bufio.NewScanner(bytes.NewReader(content))
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	var stack []openDelim
	quoteOpen := false
	quoteLine := 0

	line := 0
	for scanner.Scan() {
		line++
		text := scanner.Text()
		for i := 0; i < len(text); i++ {
			c := text[i]
			if c == '\\' {
				i++
				continue
			}
			switch c {
			case '(', '[', '{':
				stack = append(stack, openDelim{ch: c, line: line})
			case ')', ']', '}':
				want := matchOpener(c)
				if len(stack) == 0 {
					diags = append(diags, Diagnostic{
						Path: f.Path, Line: line, Rule: "balance",
						Message: fmt.Sprintf("unmatched closing %q", c),
					})
				} else {
					top := stack[len(stack)-1]
					stack = stack[:len(stack)-1]
					if top.ch != want {
						diags = append(diags, Diagnostic{
							Path: f.Path, Line: line, Rule: "balance",
							Message: fmt.Sprintf("mismatched: opened %q on line %d, closed with %q", top.ch, top.line, c),
						})
					}
				}
			case '"':
				if quoteOpen {
					quoteOpen = false
				} else {
					quoteOpen = true
					quoteLine = line
				}
			}
		}
	}

	for _, o := range stack {
		diags = append(diags, Diagnostic{
			Path: f.Path, Line: o.line, Rule: "balance",
			Message: fmt.Sprintf("unclosed %q", o.ch),
		})
	}
	if quoteOpen {
		diags = append(diags, Diagnostic{
			Path: f.Path, Line: quoteLine, Rule: "balance",
			Message: `unbalanced '"' (odd count)`,
		})
	}
	return diags
}

func matchOpener(closer byte) byte {
	switch closer {
	case ')':
		return '('
	case ']':
		return '['
	case '}':
		return '{'
	}
	return 0
}
