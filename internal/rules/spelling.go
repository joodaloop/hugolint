package rules

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/joodaloop/joodalint/internal/config"
)

// Hyphenated suffixes stripped from words before aspell sees them
// (e.g. "joodaloop-ish" → "joodaloop"). Bare uses still get flagged.
var spellingSuffixes = []string{"ish", "esque", "like", "y", "ness", "ery", "px"}

type speller struct {
	once         sync.Once
	hyphenSuffix *regexp.Regexp
	ordinal      *regexp.Regexp
	unitPrefix   *regexp.Regexp
	username     *regexp.Regexp
	aspellPath   string
	personalPath string
	initErr      error
	enabled      bool
}

var sharedSpeller = &speller{}

func (s *speller) ensureInit(cfg *config.Config) {
	s.once.Do(func() { s.init(cfg) })
}

func (s *speller) init(cfg *config.Config) {
	path, err := exec.LookPath("aspell")
	if err != nil {
		fmt.Println("joodalint: aspell not installed - skipping spell-checking. (brew install aspell)")
		s.enabled = false
		return
	}
	s.aspellPath = path

	b, err := os.ReadFile(cfg.Spelling.Dict)
	if err != nil {
		s.initErr = fmt.Errorf("spelling dict %q: %v", cfg.Spelling.Dict, err)
		return
	}
	var words []string
	sc := bufio.NewScanner(bytes.NewReader(b))
	for sc.Scan() {
		w := strings.TrimSpace(sc.Text())
		if w == "" {
			continue
		}
		// aspell rejects personal-dict words that don't start with a letter.
		r := w[0]
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')) {
			continue
		}
		words = append(words, w)
	}

	// Aspell's --personal expects its own format with a header line.
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "personal_ws-1.1 en %d utf-8\n", len(words))
	for _, w := range words {
		fmt.Fprintln(&buf, w)
	}
	s.personalPath = filepath.Join(os.TempDir(), "joodalint-aspell.pws")
	if err := os.WriteFile(s.personalPath, buf.Bytes(), 0o644); err != nil {
		s.initErr = fmt.Errorf("writing personal dict: %v", err)
		return
	}

	quoted := make([]string, len(spellingSuffixes))
	for i, suf := range spellingSuffixes {
		quoted[i] = regexp.QuoteMeta(suf)
	}
	s.hyphenSuffix = regexp.MustCompile(`-(` + strings.Join(quoted, "|") + `)\b`)
	s.ordinal = regexp.MustCompile(`\b\d+(st|nd|rd|th)\b`)
	// Unit-attached numbers: 50kg → kg (strip digits only, so a typo'd
	// unit like "50kgg" still gets caught). RE2 lacks lookahead, so we
	// capture the trailing letter and put it back via $1.
	s.unitPrefix = regexp.MustCompile(`\b\d+([a-zA-Z])`)
	s.username = regexp.MustCompile(`@[A-Za-z0-9_-]+`)

	s.enabled = true
}

// unknown runs aspell on body and returns the set of unknown words. The
// caller is responsible for checking enabled / initErr first via ready().
func (s *speller) unknown(body []byte) (map[string]bool, error) {
	body = s.hyphenSuffix.ReplaceAll(body, []byte(""))
	body = s.ordinal.ReplaceAll(body, []byte(""))
	body = s.unitPrefix.ReplaceAll(body, []byte("$1"))
	body = s.username.ReplaceAll(body, []byte(""))

	cmd := exec.Command(s.aspellPath, "--lang=en", "--encoding=utf-8", "--mode=markdown", "--personal="+s.personalPath, "list")
	cmd.Stdin = bytes.NewReader(body)
	out, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok && len(ee.Stderr) > 0 {
			return nil, fmt.Errorf("aspell failed: %v: %s", err, bytes.TrimSpace(ee.Stderr))
		}
		return nil, fmt.Errorf("aspell failed: %v", err)
	}

	res := map[string]bool{}
	sc := bufio.NewScanner(bytes.NewReader(out))
	for sc.Scan() {
		res[strings.TrimSpace(sc.Text())] = true
	}
	return res, nil
}

func buildWordRegex(words map[string]bool) *regexp.Regexp {
	parts := make([]string, 0, len(words))
	for w := range words {
		parts = append(parts, regexp.QuoteMeta(w))
	}
	return regexp.MustCompile(`\b(?:` + strings.Join(parts, "|") + `)\b`)
}
