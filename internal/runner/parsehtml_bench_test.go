package runner

import (
	"fmt"
	"strings"
	"testing"
)

// synthHTML builds a Hugo-output-shaped HTML page of roughly `paragraphs`
// body paragraphs, with realistic <head>, navigation, and footer sections.
func synthHTML(paragraphs int) []byte {
	var b strings.Builder
	b.WriteString(`<!doctype html><html lang="en"><head>`)
	b.WriteString(`<meta charset="utf-8">`)
	b.WriteString(`<meta name="viewport" content="width=device-width,initial-scale=1">`)
	b.WriteString(`<meta name="description" content="A representative description for the benchmark page.">`)
	b.WriteString(`<meta property="og:title" content="Benchmark Page">`)
	b.WriteString(`<meta property="og:type" content="article">`)
	b.WriteString(`<title>Benchmark Page Title</title>`)
	b.WriteString(`<link rel="stylesheet" href="/assets/site.css">`)
	b.WriteString(`<link rel="canonical" href="/benchmark/">`)
	b.WriteString(`<link rel="alternate" type="application/rss+xml" href="/index.xml">`)
	b.WriteString(`<script src="/assets/site.js" defer></script>`)
	b.WriteString(`</head><body>`)
	b.WriteString(`<nav><ul>`)
	for i := 0; i < 8; i++ {
		fmt.Fprintf(&b, `<li><a href="/section/%d/">Section %d</a></li>`, i, i)
	}
	b.WriteString(`</ul></nav>`)
	b.WriteString(`<main><article id="main"><h1>Benchmark Page Title</h1>`)
	for i := 0; i < paragraphs; i++ {
		fmt.Fprintf(&b, `<h2 id="section-%d">Section %d</h2>`, i, i)
		fmt.Fprintf(&b,
			`<p>Paragraph %d with <a href="/page/%d/">an internal link</a> and `+
				`<a href="https://example.com/%d">an external one</a>, plus `+
				`<img src="/img/%d.png" alt="figure %d"> inline.</p>`,
			i, i, i, i, i)
		b.WriteString(`<pre><code>code := "block " + "ignored"</code></pre>`)
	}
	b.WriteString(`</article></main>`)
	b.WriteString(`<footer><p>Footer text with <a href="/about/">about</a>.</p></footer>`)
	b.WriteString(`</body></html>`)
	return []byte(b.String())
}

func BenchmarkParseHTML_Small(b *testing.B) {
	doc := synthHTML(20)
	b.SetBytes(int64(len(doc)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _, _, _, _, _, _, _ = parseHTML(doc)
	}
}

func BenchmarkParseHTML_Medium(b *testing.B) {
	doc := synthHTML(100)
	b.SetBytes(int64(len(doc)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _, _, _, _, _, _, _ = parseHTML(doc)
	}
}

func BenchmarkParseHTML_Large(b *testing.B) {
	doc := synthHTML(500)
	b.SetBytes(int64(len(doc)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _, _, _, _, _, _, _ = parseHTML(doc)
	}
}
