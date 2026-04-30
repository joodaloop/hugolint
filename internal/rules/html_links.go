package rules

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

func init() {
	RegisterHTML(&relativeLinks{})
}

type relativeLinks struct{}

func (relativeLinks) ID() string { return "relative-link-exists" }

func (relativeLinks) Check(f *HTMLFile, ctx *HTMLContext) []Diagnostic {
	var diags []Diagnostic
	for _, href := range f.Links {
		if !isRelative(href) {
			continue
		}
		resolved, ok := resolve(f.URLPath, href)
		if !ok {
			continue
		}
		if resolved != f.URLPath {
			ctx.MarkLinked(resolved)
		}
		if !ctx.Pages[resolved] {
			diags = append(diags, Diagnostic{
				Path:    f.Path,
				Rule:    "relative-link-exists",
				Message: fmt.Sprintf("link %q resolves to %s which does not exist", href, resolved),
			})
		}
	}
	return diags
}

func isRelative(href string) bool {
	href = strings.TrimSpace(href)
	if href == "" {
		return false
	}
	if strings.HasPrefix(href, "#") {
		return false
	}
	if strings.HasPrefix(href, "//") {
		return false
	}
	if strings.HasPrefix(href, "mailto:") || strings.HasPrefix(href, "tel:") || strings.HasPrefix(href, "javascript:") || strings.HasPrefix(href, "data:") {
		return false
	}
	u, err := url.Parse(href)
	if err != nil {
		return false
	}
	return u.Scheme == ""
}

func resolve(pageURL, href string) (string, bool) {
	u, err := url.Parse(href)
	if err != nil {
		return "", false
	}
	target := u.Path
	if target == "" {
		return "", false
	}
	if !strings.HasPrefix(target, "/") {
		target = path.Join(path.Dir(strings.TrimSuffix(pageURL, "/")+"/_"), target)
	}
	target = path.Clean(target)
	if strings.HasSuffix(href, "/") && !strings.HasSuffix(target, "/") {
		target += "/"
	}
	return target, true
}
