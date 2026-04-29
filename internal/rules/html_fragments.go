package rules

import (
	"fmt"
	"net/url"
	"strings"
)

func init() {
	RegisterHTML(&fragmentLinks{})
}

type fragmentLinks struct{}

func (fragmentLinks) ID() string { return "fragment-link-exists" }

func (fragmentLinks) Check(f *HTMLFile, ctx *HTMLContext) []Diagnostic {
	var diags []Diagnostic
	for _, href := range f.Links {
		trimmed := strings.TrimSpace(href)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "//") {
			continue
		}
		if strings.HasPrefix(trimmed, "mailto:") || strings.HasPrefix(trimmed, "tel:") || strings.HasPrefix(trimmed, "javascript:") || strings.HasPrefix(trimmed, "data:") {
			continue
		}
		u, err := url.Parse(trimmed)
		if err != nil {
			continue
		}
		if u.Scheme != "" {
			continue
		}
		if u.Fragment == "" {
			continue
		}

		var target string
		if u.Path == "" {
			target = f.URLPath
		} else {
			r, ok := resolve(f.URLPath, u.Path)
			if !ok {
				continue
			}
			target = r
		}

		ids, ok := ctx.PageIDs[target]
		if !ok {
			continue
		}
		count := ids[u.Fragment]
		if count == 0 {
			diags = append(diags, Diagnostic{
				Path:    f.Path,
				Rule:    "fragment-link-exists",
				Message: fmt.Sprintf("link %q has fragment #%s not found on %s", href, u.Fragment, target),
			})
		} else if count > 1 {
			diags = append(diags, Diagnostic{
				Path:    f.Path,
				Rule:    "fragment-link-exists",
				Message: fmt.Sprintf("link %q fragment #%s appears %d times on %s", href, u.Fragment, count, target),
			})
		}
	}
	return diags
}
