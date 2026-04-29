package rules

import (
	"fmt"
	"strings"
)

func init() {
	RegisterHTML(&assetSrcs{})
}

type assetSrcs struct{}

func (assetSrcs) ID() string { return "asset-src-exists" }

func (assetSrcs) Check(f *HTMLFile, ctx *HTMLContext) []Diagnostic {
	var diags []Diagnostic
	for _, a := range f.Assets {
		if strings.HasPrefix(strings.TrimSpace(a.URL), "data:") {
			continue
		}
		if !isRelative(a.URL) {
			continue
		}
		resolved, ok := resolve(f.URLPath, a.URL)
		if !ok {
			continue
		}
		if !ctx.Pages[resolved] {
			diags = append(diags, Diagnostic{
				Path:    f.Path,
				Rule:    "asset-src-exists",
				Message: fmt.Sprintf("<%s %s=%q> resolves to %s which does not exist", a.Tag, a.Attr, a.URL, resolved),
			})
		}
	}
	return diags
}
