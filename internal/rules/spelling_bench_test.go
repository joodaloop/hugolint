package rules

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/joodaloop/joodalint/internal/config"
)

func benchSpellerCfg(b *testing.B) *config.Config {
	b.Helper()
	if _, err := exec.LookPath("aspell"); err != nil {
		b.Skip("aspell not installed")
	}
	dict := filepath.Join(b.TempDir(), "dict.txt")
	if err := os.WriteFile(dict, []byte(""), 0o644); err != nil {
		b.Fatal(err)
	}
	return &config.Config{Spelling: config.Spelling{Dict: dict}}
}

// BenchmarkAspellStartup measures the cost of one full aspell invocation —
// fork+exec, dictionary load, scan of a tiny input. This is the per-file
// tax PERF.md attributes to the spelling rule.
func BenchmarkAspellStartup(b *testing.B) {
	cfg := benchSpellerCfg(b)
	resetSpeller()
	sharedSpeller.ensureInit(cfg)
	if !sharedSpeller.enabled {
		b.Skip("speller failed to enable")
	}
	body := []byte("the quick brown fox\n")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := sharedSpeller.unknown(body); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkAspellThroughput runs aspell on a substantial body so per-byte
// scanning outweighs startup, exposing aspell's actual scan rate.
func BenchmarkAspellThroughput(b *testing.B) {
	cfg := benchSpellerCfg(b)
	resetSpeller()
	sharedSpeller.ensureInit(cfg)
	if !sharedSpeller.enabled {
		b.Skip("speller failed to enable")
	}
	body := bytes.Repeat([]byte("the quick brown fox jumps over the lazy dog\n"), 1200) // ~50 KB
	b.SetBytes(int64(len(body)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := sharedSpeller.unknown(body); err != nil {
			b.Fatal(err)
		}
	}
}
