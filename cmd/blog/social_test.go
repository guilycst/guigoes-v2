package main

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSocialCoverReplacesGlobalImage(t *testing.T) {
	root := t.TempDir()
	render := filepath.Join(root, "public")
	post := testDocument(t, "cover = './assets/a capa.png'\ncover_alt = 'Código & dados <juntos>'")
	var cover bytes.Buffer
	if err := png.Encode(&cover, image.NewRGBA(image.Rect(0, 0, 120, 63))); err != nil {
		t.Fatal(err)
	}
	coverFile := filepath.Join(root, "content/blogs/test/assets/a capa.png")
	if err := writeFile(coverFile, cover.Bytes()); err != nil {
		t.Fatal(err)
	}
	body := `<body><script type="application/ld+json">{"image":"keep me"}</script><p>Keep &amp; demos</p></body></html>`
	input := `<html><head><title>Test</title><meta property='og:image' content='old.png'><meta content="duplicate.png" property="og:image"/><meta property="og:image:width" content="1280"><meta property="og:image:secure_url" content="old.png"><meta name="twitter:image" content="old.png"><meta name="twitter:card" content="summary_large_image"></head>` + body
	file := filepath.Join(render, "blogs/test/index.html")
	for _, path := range []string{file, filepath.Join(render, "index.html")} {
		if err := writeFile(path, []byte(input)); err != nil {
			t.Fatal(err)
		}
	}
	withoutCover := post
	withoutCover.Meta.Cover = ""
	if err := applySocialCovers(root, render, []document{withoutCover}); err != nil {
		t.Fatal(err)
	}
	unchanged, _ := os.ReadFile(file)
	if string(unchanged) != input {
		t.Fatal("post without cover lost the site fallback")
	}
	if err := applySocialCovers(root, render, []document{post}); err != nil {
		t.Fatal(err)
	}
	output, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(string(output), body) || bytes.Contains(output, []byte("old.png")) || bytes.Contains(output, []byte("duplicate.png")) {
		t.Fatal("article changed or old image remains")
	}
	page, err := readHTML(output)
	if err != nil {
		t.Fatal(err)
	}
	for key, want := range map[string]string{
		"og:image":       "https://guigoes.com/blogs/test/assets/a%20capa.png",
		"twitter:image":  "https://guigoes.com/blogs/test/assets/a%20capa.png",
		"og:image:width": "120", "og:image:height": "63", "og:image:type": "image/png",
		"og:image:alt": "Código & dados <juntos>", "twitter:image:alt": "Código & dados <juntos>",
		"twitter:card": "summary_large_image",
	} {
		got := page.HeadMetadata[key]
		if len(got) != 1 || got[0] != want {
			t.Errorf("%s: got %v, want %q once", key, got, want)
		}
	}
	home, _ := os.ReadFile(filepath.Join(render, "index.html"))
	if string(home) != input {
		t.Fatal("changed non-post social image")
	}
	post.Meta.CoverAlt = ""
	values, err := coverMetadata(root, post)
	if err != nil || values["og:image:alt"] != post.Meta.Title {
		t.Fatalf("title fallback: %v, %v", values, err)
	}
	if err := os.WriteFile(coverFile, []byte("not an image"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := applySocialCovers(root, render, []document{post}); err == nil {
		t.Fatal("accepted invalid cover")
	}
	if _, err := replaceSocialImages([]byte("<body>no head</body>"), values); err == nil {
		t.Fatal("accepted missing head")
	}
}
