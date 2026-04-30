package rules

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/joodaloop/hugolint/internal/config"
)

func init() {
	RegisterMarkdown(&markdownSpelling{})
}

// Hyphenated suffixes stripped from words before aspell sees them
// (e.g. "joodaloop-ish" → "joodaloop"). Bare uses still get flagged.
var spellingSuffixes = []string{"ish", "esque", "like", "y", "ness", "ery", "px"}

type markdownSpelling struct {
	once        sync.Once
	dict        map[string]bool
	suffixStrip *regexp.Regexp
	aspellPath  string
	initErr     error
	enabled     bool
	errOnce     sync.Once
}

func (*markdownSpelling) ID() string { return "spelling" }

func (m *markdownSpelling) Check(f *MarkdownFile, ctx *MarkdownContext) []Diagnostic {
	if ctx == nil || ctx.Config == nil || ctx.Config.Spelling.Dict == "" {
		return nil
	}
	m.once.Do(func() { m.init(ctx.Config) })
	if m.initErr != nil {
		var d []Diagnostic
		m.errOnce.Do(func() {
			d = []Diagnostic{{Path: f.Path, Rule: "spelling", Message: m.initErr.Error()}}
		})
		return d
	}
	if !m.enabled {
		return nil
	}

	body := f.Body
	if m.suffixStrip != nil {
		body = m.suffixStrip.ReplaceAll(body, []byte(""))
	}

	cmd := exec.Command(m.aspellPath, "--mode=markdown", "--lang=en", "list")
	cmd.Stdin = bytes.NewReader(body)
	out, err := cmd.Output()
	if err != nil {
		return []Diagnostic{{Path: f.Path, Rule: "spelling", Message: fmt.Sprintf("aspell failed: %v", err)}}
	}

	unknown := map[string]bool{}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		w := strings.TrimSpace(scanner.Text())
		if w == "" || m.dict[w] {
			continue
		}
		unknown[w] = true
	}
	if len(unknown) == 0 {
		return nil
	}

	// Locate each unknown word's first occurrence in the original content
	// for line-numbered diagnostics.
	wordRe := buildWordRegex(unknown)
	var diags []Diagnostic
	seen := map[string]bool{}
	lineScanner := bufio.NewScanner(bytes.NewReader(f.Body))
	lineScanner.Buffer(make([]byte, 64*1024), 1024*1024)
	line := f.BodyStartLine - 1
	for lineScanner.Scan() {
		line++
		for _, m := range wordRe.FindAllString(lineScanner.Text(), -1) {
			if seen[m] {
				continue
			}
			seen[m] = true
			diags = append(diags, Diagnostic{
				Path: f.Path, Line: line, Rule: "spelling",
				Message: fmt.Sprintf("unknown word: %q", m),
			})
		}
	}
	return diags
}

func (m *markdownSpelling) init(cfg *config.Config) {
	path, err := exec.LookPath("aspell")
	if err != nil {
		fmt.Println("hugolint: aspell not installed - skipping spell-checking. (brew install aspell)")
		m.enabled = false
		return
	}
	m.aspellPath = path

	m.dict = map[string]bool{}
	if b, err := os.ReadFile(cfg.Spelling.Dict); err == nil {
		s := bufio.NewScanner(bytes.NewReader(b))
		for s.Scan() {
			w := strings.TrimSpace(s.Text())
			if w != "" {
				m.dict[w] = true
			}
		}
	} else {
		m.initErr = fmt.Errorf("spelling dict %q: %v", cfg.Spelling.Dict, err)
		return
	}

	quoted := make([]string, len(spellingSuffixes))
	for i, s := range spellingSuffixes {
		quoted[i] = regexp.QuoteMeta(s)
	}
	m.suffixStrip = regexp.MustCompile(`-(` + strings.Join(quoted, "|") + `)\b`)

	m.enabled = true
}

func buildWordRegex(words map[string]bool) *regexp.Regexp {
	parts := make([]string, 0, len(words))
	for w := range words {
		parts = append(parts, regexp.QuoteMeta(w))
	}
	return regexp.MustCompile(`\b(?:` + strings.Join(parts, "|") + `)\b`)
}
