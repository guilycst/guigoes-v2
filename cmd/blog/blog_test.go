package main

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func testDocument(t *testing.T, extra string) document {
	t.Helper()
	d, err := parseDocument([]byte("+++\ntitle = 'Go & HTML'\nlanguage = 'pt-BR'\npublishedAt = '2024-04-30T01:37:09-03:00'\nauthors = ['Ana', 'Rui']\n" + extra + "\n+++\n# Go & HTML\n\nPost body.\n"))
	if err != nil {
		t.Fatal(err)
	}
	d.Source = "blogs/test/index.md"
	d.Route = "/blogs/test/"
	return d
}

func TestStructuredData(t *testing.T) {
	d := testDocument(t, "cover = './assets/a capa.webp'\nmodifiedAt = '2024-05-01T12:00:00Z'")
	d.Meta.Title = "Go </script><script>alert(1)</script> & HTML"
	script, err := schemaScript(d)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(script, "</script>") != 1 || strings.Contains(script, "<script>alert") {
		t.Fatal("metadata can break out of script")
	}
	page, err := readHTML([]byte(`<html><head><link rel="canonical" href="https://guigoes.com/blogs/test/"></head><body>` + script + `</body></html>`))
	if err != nil {
		t.Fatal(err)
	}
	if err := checkSchema(page, &d); err != nil {
		t.Fatal(err)
	}
	var got blogPosting
	if err := json.Unmarshal(page.Schemas[0], &got); err != nil {
		t.Fatal(err)
	}
	if got.Type != "BlogPosting" || got.Headline != d.Meta.Title || got.Image != "https://guigoes.com/blogs/test/assets/a%20capa.webp" || got.Published != "2024-04-30T01:37:09-03:00" || got.Modified != "2024-05-01T12:00:00Z" {
		t.Fatalf("wrong schema: %+v", got)
	}
	if len(got.Authors) != 2 || got.Authors[0].Name != "Ana" || got.Authors[1].Type != "Person" {
		t.Fatalf("wrong authors: %+v", got.Authors)
	}
	d.Meta.Cover = ""
	d.Meta.ModifiedAt = ""
	script, err = schemaScript(d)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(script, `"image"`) || strings.Contains(script, `"dateModified"`) {
		t.Fatal("invented optional metadata")
	}
}

func TestAuthorProfile(t *testing.T) {
	d := testDocument(t, "")
	d.Meta.Authors = []string{"Guilherme de Castro", "Guest"}
	data := structuredData(d)
	if data.Authors[0].URL != "https://guigoes.com/about/" || data.Authors[1].URL != "" {
		t.Fatalf("incorrect author profiles: %+v", data.Authors)
	}
	header := publicationDetails(d.Meta)
	if !strings.Contains(header, `<a rel="author" href="https://guigoes.com/about/">Guilherme de Castro</a>`) || !strings.Contains(header, `<span rel="author">Guest</span>`) {
		t.Fatalf("incorrect visible author links: %s", header)
	}
}

func TestSchemaValidationRejectsDrift(t *testing.T) {
	d := testDocument(t, "")
	script, err := schemaScript(d)
	if err != nil {
		t.Fatal(err)
	}
	valid := `<html><head><link rel="canonical" href="https://guigoes.com/blogs/test/"></head><body>` + script + `</body></html>`
	for name, markup := range map[string]string{
		"missing":         strings.Replace(valid, script, "", 1),
		"duplicate":       strings.Replace(valid, script, script+script, 1),
		"wrong date":      strings.Replace(valid, "2024-04-30", "2023-04-30", 1),
		"wrong canonical": strings.Replace(valid, `href="https://guigoes.com/blogs/test/"`, `href="https://guigoes.com/"`, 1),
		"invalid JSON":    strings.Replace(valid, `"@context"`, `broken`, 1),
	} {
		t.Run(name, func(t *testing.T) {
			page, err := readHTML([]byte(markup))
			if err != nil {
				t.Fatal(err)
			}
			if checkSchema(page, &d) == nil {
				t.Fatal("accepted invalid schema")
			}
		})
	}
	page, err := readHTML([]byte(valid))
	if err != nil {
		t.Fatal(err)
	}
	if checkSchema(page, nil) == nil {
		t.Fatal("accepted post schema on archive")
	}
}

func TestPublicationDates(t *testing.T) {
	now := time.Date(2026, 10, 27, 0, 0, 0, 0, time.UTC)
	for _, tc := range []struct {
		date        string
		draft, want bool
	}{
		{"2026-10-27T00:00:00-03:00", false, false},
		{"2026-10-26T21:00:00-03:00", false, true},
		{"2024-04-30T01:37:09-03:00", true, false},
	} {
		d := testDocument(t, "")
		d.Meta.Draft = tc.draft
		var err error
		d.Published, err = time.Parse(time.RFC3339, tc.date)
		if err != nil {
			t.Fatal(err)
		}
		if d.visible(now) != tc.want {
			t.Errorf("date %s draft %v visibility wrong", tc.date, tc.draft)
		}
	}
	for _, front := range []string{
		"publishedAt = '2026-10-27T00:00:00'",
		"publishedAt = 'invalid'",
		"publishedAt = '2026-10-27T00:00:00Z'\nmodifiedAt = '2026-10-26T23:59:59Z'",
	} {
		if _, err := parseDocument([]byte("+++\ntitle = 'test'\n" + front + "\n+++\nbody")); err == nil {
			t.Errorf("accepted invalid dates: %s", front)
		}
	}
}

func TestEditorialHTML(t *testing.T) {
	d := testDocument(t, "cover = './assets/cover.webp'\ncover_alt = 'A \"cover\"'\ntags = ['SEO', 'Automação']")
	d.Meta.Title = `<script>alert("title")</script>`
	d.Meta.Description = "Go & HTML"
	for source, want := range map[string]string{"index.md": "blogs/test/assets/cover.webp", "blogs/index.md": "test/assets/cover.webp", "tags/seo/index.md": "../../blogs/test/assets/cover.webp"} {
		output, err := listing([]document{d}, source)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(output, "<script>") || !strings.Contains(output, "&lt;script&gt;") || !strings.Contains(output, "Go &amp; HTML") || !strings.Contains(output, `src="`+want+`"`) {
			t.Fatalf("unsafe or incorrect listing: %s", output)
		}
	}
	for tag, want := range map[string]string{"SEO": "seo", "social media": "social-media", "Automação": "automacao", "///": ""} {
		if got := tagSlug(tag); got != want {
			t.Errorf("slug %q: %q", tag, got)
		}
	}
	output := publicationDetails(d.Meta)
	if !strings.Contains(output, `>30/04/2024</time>`) || !strings.Contains(output, `href="/tags/automacao/"`) {
		t.Fatalf("bad header: %s", output)
	}
}

func TestPageHasOneSchemaAfterHeading(t *testing.T) {
	d := testDocument(t, "")
	stage := t.TempDir()
	if err := writePage(stage, d); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(stage, d.Source))
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	if strings.Index(text, "# Go & HTML") > strings.Index(text, "data-margo-publication-metadata") {
		t.Fatal("metadata before title")
	}
	if strings.Count(text, `application/ld+json`) != 1 {
		t.Fatal("expected one JSON-LD block")
	}
	d.Route = ""
	d.Source = "index.md"
	if err := writePage(stage, d); err != nil {
		t.Fatal(err)
	}
	data, err = os.ReadFile(filepath.Join(stage, d.Source))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "application/ld+json") {
		t.Fatal("schema leaked to home")
	}
}

func TestFeedEscapingAndTimezone(t *testing.T) {
	d := testDocument(t, "")
	d.Meta.Title = "Go & <HTML>"
	file := filepath.Join(t.TempDir(), "index.xml")
	if err := writeFeed(file, "Guigoes", []document{d}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	var feed rss
	if err := xml.Unmarshal(data, &feed); err != nil {
		t.Fatal(err)
	}
	if len(feed.Channel.Items) != 1 || feed.Channel.Items[0].Title != d.Meta.Title || !strings.HasSuffix(feed.Channel.Items[0].Date, "-0300") {
		t.Fatalf("wrong feed: %+v", feed)
	}
}

func TestLoadDocumentsRejectsMissingMetadataAndCover(t *testing.T) {
	for name, extra := range map[string]string{"missing author": "authors = []", "missing date": "publishedAt = ''", "missing cover": "cover = './assets/missing.webp'", "escaping cover": "cover = '../other.webp'"} {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			front := "title = 'Post'\nauthors = ['Ana']\npublishedAt = '2024-01-01T00:00:00Z'\n"
			if name == "missing author" {
				front = "title = 'Post'\npublishedAt = '2024-01-01T00:00:00Z'\n"
			}
			if name == "missing date" {
				front = "title = 'Post'\nauthors = ['Ana']\n"
			}
			if err := writeFile(filepath.Join(root, "content/blogs/test/index.md"), []byte("+++\n"+front+extra+"\n+++\nbody")); err != nil {
				t.Fatal(err)
			}
			if _, _, err := loadDocuments(root, time.Now()); err == nil {
				t.Fatal("accepted invalid post")
			}
		})
	}
}

func TestAliasesAndPublicationRollback(t *testing.T) {
	d := testDocument(t, "")
	render := t.TempDir()
	d.Meta.Aliases = []string{"/posts/test/"}
	if err := writeAliases(render, []document{d}); err != nil {
		t.Fatal(err)
	}
	if err := writeAliases(render, []document{d}); err == nil {
		t.Fatal("alias collision ignored")
	}
	d.Meta.Aliases = []string{"/../../escape/"}
	if err := writeAliases(render, []document{d}); err == nil {
		t.Fatal("unsafe alias accepted")
	}
	root := t.TempDir()
	public := filepath.Join(root, "public")
	if err := writeFile(filepath.Join(public, "index.html"), []byte("old")); err != nil {
		t.Fatal(err)
	}
	if err := publishDirectory(filepath.Join(root, "missing"), public); err == nil {
		t.Fatal("missing render accepted")
	}
	data, err := os.ReadFile(filepath.Join(public, "index.html"))
	if err != nil || string(data) != "old" {
		t.Fatal("old output lost on failed publish")
	}
}

func TestBuildIntegration(t *testing.T) {
	repo, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(repo, ".tools/margo/margo")
	if _, err := os.Stat(binary); os.IsNotExist(err) {
		t.Skip("run make install to enable the real Margo integration test")
	} else if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	for _, dir := range []string{"content", "assets", "static", "tests"} {
		if err := copyTree(filepath.Join(repo, dir), filepath.Join(root, dir)); err != nil {
			t.Fatal(err)
		}
	}
	config, err := os.ReadFile(filepath.Join(repo, "site.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if err := writeFile(filepath.Join(root, "site.yaml"), config); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, ".tools/margo"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(binary, filepath.Join(root, ".tools/margo/margo")); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"draft", "future"} {
		front := "draft = true\npublishedAt = '2020-01-01T00:00:00Z'"
		if name == "future" {
			front = "publishedAt = '2099-01-01T00:00:00Z'"
		}
		if err := writeFile(filepath.Join(root, "content/blogs", name, "index.md"), []byte("+++\ntitle = 'PRIVATE_POST'\nauthors = ['Ana']\n"+front+"\n+++\nPRIVATE_BODY")); err != nil {
			t.Fatal(err)
		}
		if err := writeFile(filepath.Join(root, "content/blogs", name, "assets/private.txt"), []byte("private asset")); err != nil {
			t.Fatal(err)
		}
	}
	now := time.Date(2026, 10, 27, 12, 0, 0, 0, time.UTC)
	if err := buildSite(root, now); err != nil {
		t.Fatal(err)
	}
	public := filepath.Join(root, "public")
	for _, name := range []string{"draft", "future"} {
		if _, err := os.Stat(filepath.Join(public, "blogs", name)); !os.IsNotExist(err) {
			t.Fatalf("unpublished post leaked: %s", name)
		}
	}
	for _, file := range []string{"index.html", "index.xml", "sitemap.xml"} {
		data, err := os.ReadFile(filepath.Join(public, file))
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(data), "PRIVATE_") {
			t.Fatalf("private metadata leaked in %s", file)
		}
	}
	original, err := os.ReadFile(filepath.Join(public, "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	// A valid render with a missing legacy route must not replace the last build.
	if err := writeFile(filepath.Join(root, "tests/legacy-urls.txt"), []byte("/missing-route/\n")); err != nil {
		t.Fatal(err)
	}
	if err := buildSite(root, now); err == nil {
		t.Fatal("invalid site was published")
	}
	preserved, err := os.ReadFile(filepath.Join(public, "index.html"))
	if err != nil || string(preserved) != string(original) {
		t.Fatal("failed validation replaced the previous publication")
	}
}
