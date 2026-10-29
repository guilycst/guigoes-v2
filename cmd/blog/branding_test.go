package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/araihu/margo/site"
)

func TestBrandAssetDimensionsAndHead(t *testing.T) {
	root := "../.."
	config, err := site.LoadConfig(filepath.Join(root, "site.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	for _, icon := range config.Site.Icons {
		if !strings.HasSuffix(icon.Path, ".png") {
			continue
		}
		f, err := os.Open(filepath.Join(root, icon.Path))
		if err != nil {
			t.Fatal(err)
		}
		cfg, _, err := image.DecodeConfig(f)
		f.Close()
		if err != nil {
			t.Fatal(err)
		}
		if icon.Sizes != fmtSize(cfg.Width, cfg.Height) {
			t.Fatalf("wrong declared dimensions for %s", icon.Path)
		}
	}
	data, err := os.ReadFile(filepath.Join(root, "site.webmanifest"))
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
