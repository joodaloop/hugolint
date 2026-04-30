package rules

import (
	"os"
	"path"
	"regexp"
	"strings"
)

type BuiltFile struct {
	Path    string
	URLPath string
}

func ReportOrphans(files []BuiltFile, ctx *HTMLContext) []Diagnostic {
	var diags []Diagnostic
	for _, f := range files {
		if isEntryPoint(f.URLPath) || isWellKnown(f.URLPath) {
			continue
		}
		linked := false
		for _, alias := range pageAliases(f.URLPath) {
			if ctx.LinkedPages[alias] {
				linked = true
				break
			}
		}
		if linked {
			continue
		}
		diags = append(diags, Diagnostic{
			Path:    f.Path,
			Rule:    "orphan-file",
			Message: "file is not linked to from any other page",
		})
	}
	return diags
}

func isEntryPoint(urlPath string) bool {
	return urlPath == "/" || urlPath == "" || urlPath == "/index.html"
}

func isWellKnown(urlPath string) bool {
	base := path.Base(urlPath)
	if strings.HasPrefix(base, ".") {
		return true
	}
	if strings.HasPrefix(base, "favicon.") {
		return true
	}
	switch base {
	case "404.html", "robots.txt", "sitemap.xml", "sw.js", "manifest.json":
		return true
	}
	return false
}

func pageAliases(urlPath string) []string {
	aliases := []string{urlPath}
	if strings.HasSuffix(urlPath, "/") {
		aliases = append(aliases, strings.TrimSuffix(urlPath, "/"))
		aliases = append(aliases, urlPath+"index.html")
	}
	return aliases
}

var cssURLRegex = regexp.MustCompile(`url\(\s*['"]?([^'")\s]+)['"]?\s*\)`)

func ScanCSSLinks(files []BuiltFile, ctx *HTMLContext) error {
	for _, f := range files {
		if !strings.HasSuffix(f.Path, ".css") {
			continue
		}
		b, err := os.ReadFile(f.Path)
		if err != nil {
			return err
		}
		for _, m := range cssURLRegex.FindAllSubmatch(b, -1) {
			ref := string(m[1])
			if strings.HasPrefix(ref, "data:") {
				continue
			}
			if !isRelative(ref) {
				continue
			}
			resolved, ok := resolve(f.URLPath, ref)
			if !ok {
				continue
			}
			ctx.MarkLinked(resolved)
		}
	}
	return nil
}
