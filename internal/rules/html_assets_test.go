package rules

import "testing"

func TestHTMLAssets_MissingAsset(t *testing.T) {
	f := &HTMLFile{
		Path: "/site/public/x.html", URLPath: "/x.html",
		Assets: []Asset{{Tag: "script", Attr: "src", URL: "/missing.js"}},
	}
	ctx := htmlCtx(map[string]bool{"/x.html": true}, nil)
	diags := assetSrcs{}.Check(f, ctx)
	if !containsMsg(diags, "/missing.js") {
		t.Fatalf("want missing diag, got %v", messages(diags))
	}
}

func TestHTMLAssets_AssetExists(t *testing.T) {
	f := &HTMLFile{
		Path: "/site/public/x.html", URLPath: "/x.html",
		Assets: []Asset{{Tag: "link", Attr: "href", URL: "/style.css"}},
	}
	ctx := htmlCtx(map[string]bool{"/style.css": true}, nil)
	diags := assetSrcs{}.Check(f, ctx)
	assertNoDiags(t, diags)
}

func TestHTMLAssets_DataURISkipped(t *testing.T) {
	f := &HTMLFile{
		Path: "/site/public/x.html", URLPath: "/x.html",
		Assets: []Asset{{Tag: "link", Attr: "href", URL: "data:text/css,body{}"}},
	}
	ctx := htmlCtx(map[string]bool{}, nil)
	assertNoDiags(t, assetSrcs{}.Check(f, ctx))
}

func TestHTMLAssets_AbsoluteSkipped(t *testing.T) {
	f := &HTMLFile{
		Path: "/site/public/x.html", URLPath: "/x.html",
		Assets: []Asset{{Tag: "script", Attr: "src", URL: "https://cdn.example.com/x.js"}},
	}
	ctx := htmlCtx(map[string]bool{}, nil)
	assertNoDiags(t, assetSrcs{}.Check(f, ctx))
}

func TestHTMLImages_MissingImage(t *testing.T) {
	f := &HTMLFile{Path: "/site/public/x.html", URLPath: "/x.html", Images: []string{"/missing.png"}}
	ctx := htmlCtx(map[string]bool{"/x.html": true}, nil)
	diags := imageSrcs{}.Check(f, ctx)
	if !containsMsg(diags, "/missing.png") {
		t.Fatalf("want missing img diag, got %v", messages(diags))
	}
}

func TestHTMLImages_DataURISkipped(t *testing.T) {
	f := &HTMLFile{Path: "/site/public/x.html", URLPath: "/x.html", Images: []string{"data:image/png;base64,xx"}}
	ctx := htmlCtx(map[string]bool{}, nil)
	assertNoDiags(t, imageSrcs{}.Check(f, ctx))
}
