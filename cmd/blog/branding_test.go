package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestBrandHeadPreservesArticleAndSchema(t *testing.T) {
	body := `<body><article><script type="application/ld+json">{"headline":"Go & HTML"}</script><p>Article &amp; demos</p></article></body></html>`
	input := []byte(`<html><head><link rel="icon" href="/old.svg"><meta property="og:image" content="https://guigoes.com/assets/brand/social-v3.png"></head>` + body)
	output, err := addBrandHead(input)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(string(output), body) {
		t.Fatal("article bytes changed")
	}
	if bytes.Contains(output, []byte("old.svg")) {
		t.Fatal("old icon remains")
	}
	if strings.Count(string(output), `rel="manifest"`) != 1 || strings.Count(string(output), `property="og:image"`) != 1 {
		t.Fatal("missing or duplicate metadata")
	}
	if _, err := addBrandHead(output); err == nil {
		t.Fatal("repeated injection must fail")
	}
}

func TestBrandAssetDimensionsAndHead(t *testing.T) {
	root := "../.."
	tokenizer := html.NewTokenizer(strings.NewReader(brandHead))
	for {
		kind := tokenizer.Next()
		if kind == html.ErrorToken {
			break
		}
		if kind != html.StartTagToken && kind != html.SelfClosingTagToken {
			continue
		}
		token := tokenizer.Token()
		attrs := map[string]string{}
		for _, a := range token.Attr {
			attrs[a.Key] = a.Val
		}
		if token.Data != "link" || attrs["type"] != "image/png" && attrs["rel"] != "apple-touch-icon" {
			continue
		}
		f, err := os.Open(filepath.Join(root, strings.TrimPrefix(attrs["href"], "/")))
		if err != nil {
			t.Fatal(err)
		}
		cfg, _, err := image.DecodeConfig(f)
		f.Close()
		if err != nil {
			t.Fatal(err)
		}
		size := fmtSize(cfg.Width, cfg.Height)
		if attrs["sizes"] != size {
			t.Fatalf("%s declares %s, has %s", attrs["href"], attrs["sizes"], size)
		}
	}
	data, err := os.ReadFile(filepath.Join(root, "static/site.webmanifest"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		Icons []struct{ Src, Sizes, Type, Purpose string }
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	if len(manifest.Icons) != 4 {
		t.Fatal("missing Android icons")
	}
	for _, icon := range manifest.Icons {
		file, err := os.Open(filepath.Join(root, strings.TrimPrefix(icon.Src, "/")))
		if err != nil {
			t.Fatal(err)
		}
		img, _, err := image.Decode(file)
		file.Close()
		if err != nil {
			t.Fatal(err)
		}
		if fmtSize(img.Bounds().Dx(), img.Bounds().Dy()) != icon.Sizes {
			t.Fatal("manifest dimensions differ")
		}
		if icon.Purpose == "maskable" {
			// All meaningful pixels stay inside the central 80%-diameter safe circle.
			n := img.Bounds().Dx()
			radius := float64(n) * .4
			for y := 0; y < n; y++ {
				for x := 0; x < n; x++ {
					r, g, b, a := img.At(x, y).RGBA()
					if a != 65535 {
						t.Fatal("maskable must be opaque")
					}
					if b > g && g > r {
						dx, dy := float64(x)-float64(n)/2, float64(y)-float64(n)/2
						if dx*dx+dy*dy > radius*radius {
							t.Fatal("maskable symbol outside safe circle")
						}
					}
				}
			}
		}
	}
	for _, name := range []string{"logo-v2.svg", "logo-light-v2.svg", "icon-v2.svg", "social-v3.svg", "social-square-v3.svg"} {
		data, err := os.ReadFile(filepath.Join(root, "assets/brand", name))
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Contains(data, []byte("<image")) || bytes.Contains(data, []byte("<text")) {
			t.Fatalf("%s must contain outlined vectors", name)
		}
	}
}

func fmtSize(w, h int) string { return fmt.Sprintf("%dx%d", w, h) }

func TestBrandManifestHashesFinalBytes(t *testing.T) {
	root := t.TempDir()
	if err := writeFile(filepath.Join(root, "index.html"), []byte("final page")); err != nil {
		t.Fatal(err)
	}
	if err := writeFile(filepath.Join(root, "margo-manifest.json"), []byte(`{"schemaVersion":"margo-site-manifest/v1","digest":"old","entries":[]}`)); err != nil {
		t.Fatal(err)
	}
	if err := refreshBrandManifest(root); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, "margo-manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var m struct {
		Digest  string
		Entries []struct{ Path, Digest string }
	}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256([]byte("final page"))
	if len(m.Entries) != 1 || m.Entries[0].Digest != hex.EncodeToString(digest[:]) {
		t.Fatal("artifact hash differs from output bytes")
	}
	if err := refreshBrandManifest(root); err != nil {
		t.Fatal(err)
	}
	again, err := os.ReadFile(filepath.Join(root, "margo-manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, again) {
		t.Fatal("manifest is not deterministic")
	}
}
