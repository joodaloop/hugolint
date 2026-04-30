package rules

import (
	"fmt"
	"strings"
)

func init() {
	RegisterHTML(&imageSrcs{})
}

type imageSrcs struct{}

func (imageSrcs) ID() string { return "image-src-exists" }

func (imageSrcs) Check(f *HTMLFile, ctx *HTMLContext) []Diagnostic {
	var diags []Diagnostic
	for _, src := range f.Images {
		if strings.HasPrefix(strings.TrimSpace(src), "data:") {
			continue
		}
		if !isRelative(src) {
			continue
		}
		resolved, ok := resolve(f.URLPath, src)
		if !ok {
			continue
		}
		ctx.MarkLinked(resolved)
		if !ctx.Pages[resolved] {
			diags = append(diags, Diagnostic{
				Path:    f.Path,
				Rule:    "image-src-exists",
				Message: fmt.Sprintf("img src %q resolves to %s which does not exist", src, resolved),
			})
		}
	}
	return diags
}
