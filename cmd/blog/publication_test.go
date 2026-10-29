package main

import (
	"context"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/araihu/gazeta"
)

func copyFixture(t *testing.T, source, target string) {
	t.Helper()
	err := filepath.WalkDir(source, func(file string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(source, file)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		destination := filepath.Join(target, relative)
		if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
			return err
		}
		return os.WriteFile(destination, data, 0644)
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGuigoesConsumesGazeta(t *testing.T) {
	repo, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	for _, dir := range []string{"content", "assets", "static", "tests"} {
		copyFixture(t, filepath.Join(repo, dir), filepath.Join(root, dir))
	}
	for _, file := range []string{"site.yaml", "blog.yaml", "site.webmanifest"} {
		data, err := os.ReadFile(filepath.Join(repo, file))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, file), data, 0644); err != nil {
			t.Fatal(err)
		}
	}
	now := time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC)
	if err := gazeta.Build(context.Background(), root, now); err != nil {
		t.Fatal(err)
	}
	if err := gazeta.Check(root, now); err != nil {
		t.Fatal(err)
	}
	read := func(name string) string {
		t.Helper()
		data, err := os.ReadFile(filepath.Join(root, "public", name))
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}
	post := read("blogs/lazy-dvc-github-ssh-dvc-storage/index.html")
	for _, want := range []string{
		`href="https://guigoes.com/blogs/lazy-dvc-github-ssh-dvc-storage/"`,
		`href="https://guigoes.com/about/"`, `data-margo-publication-date="published"`,
		`data-margo-publication-tags`, `class="blog-content"`, `BlogPosting`, `Guilherme de Castro`,
	} {
		if !strings.Contains(post, want) {
			t.Fatalf("missing consumer contract %q", want)
		}
	}
	about := read("about/index.html")
	if !strings.Contains(about, `"@type":"ProfilePage"`) || !strings.Contains(about, `"name":"Guilherme de Castro"`) {
		t.Fatal("configured author profile schema missing")
	}
	home := read("index.html")
	if !strings.Contains(home, `class="blog-home"`) || !strings.Contains(home, "Guigoes") {
		t.Fatal("consumer home identity lost")
	}
	if !strings.Contains(read("assets/blog.css"), "--guigoes-surface") || !strings.Contains(read("assets/blog.css"), ".blog-home") {
		t.Fatal("consumer styling lost")
	}
	for _, file := range []string{"favicon.ico", "apple-touch-icon.png", "site.webmanifest", "robots.txt", "_redirects"} {
		read(file)
	}
	var manifest struct {
		Entries []struct{ Path, Digest string }
	}
	if err := json.Unmarshal([]byte(read("margo-manifest.json")), &manifest); err != nil {
		t.Fatal(err)
	}
	if len(manifest.Entries) == 0 {
		t.Fatal("final publication manifest missing")
	}
	// An invalid compatibility requirement must preserve the last successful build.
	before := read("margo-manifest.json")
	if err := os.WriteFile(filepath.Join(root, "tests/legacy-urls.txt"), []byte("/missing-legacy-route/"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := gazeta.Build(context.Background(), root, now); err == nil {
		t.Fatal("broken legacy route accepted")
	}
	if read("margo-manifest.json") != before {
		t.Fatal("consumer output changed after failed validation")
	}
}
