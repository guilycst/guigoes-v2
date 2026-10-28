package main

import (
	"encoding/xml"
	"fmt"
	"html"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

func writeFile(file string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return err
	}
	return os.WriteFile(file, data, 0644)
}

func copyTree(source, target string) error {
	return filepath.WalkDir(source, func(file string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("cannot copy symlink %s", file)
		}
		if entry.IsDir() {
			return nil
		}
		if !entry.Type().IsRegular() {
			return fmt.Errorf("not a regular file: %s", file)
		}
		relative, err := filepath.Rel(source, file)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		return writeFile(filepath.Join(target, relative), data)
	})
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	GUID        string `xml:"guid"`
	Description string `xml:"description"`
	Date        string `xml:"pubDate"`
}
type rssChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Language    string    `xml:"language"`
	Items       []rssItem `xml:"item"`
}
type rss struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

func writeFeed(file, title string, posts []document) error {
	feed := rss{Version: "2.0", Channel: rssChannel{Title: title, Link: baseURL + "/", Description: "Blog de Guilherme de Castro", Language: "pt-BR"}}
	for _, post := range posts {
		feed.Channel.Items = append(feed.Channel.Items, rssItem{Title: post.Meta.Title, Link: baseURL + post.Route, GUID: baseURL + post.Route, Description: post.Meta.Description, Date: post.Published.Format("Mon, 02 Jan 2006 15:04:05 -0700")})
	}
	data, err := xml.Marshal(feed)
	if err != nil {
		return err
	}
	return writeFile(file, append([]byte(xml.Header), data...))
}

func generatedPage(stage, source, title, body string, items ...schemaListItem) error {
	return writePage(stage, document{SchemaItems: items, Source: source, Body: body, Meta: metadata{Title: title, Language: "pt-BR"}, Front: map[string]any{"title": title, "language": "pt-BR"}})
}

func buildSite(root string, now time.Time) error {
	docs, posts, err := loadDocuments(root, now)
	if err != nil {
		return err
	}
	stage, render := filepath.Join(root, ".build/content"), filepath.Join(root, ".build/render")
	for _, dir := range []string{stage, render} {
		if err := os.RemoveAll(dir); err != nil {
			return err
		}
	}
	if err := os.MkdirAll(stage, 0755); err != nil {
		return err
	}
	for _, d := range docs {
		sourceAssets := filepath.Join(root, "content", filepath.Dir(d.Source), "assets")
		if _, err := os.Stat(sourceAssets); err == nil {
			if err := copyTree(sourceAssets, filepath.Join(stage, filepath.Dir(d.Source), "assets")); err != nil {
				return err
			}
		} else if !os.IsNotExist(err) {
			return err
		}
		if d.Route != "" {
			index, pagination := seriesNavigation(d, posts)
			d.Body = index + d.Body + "\n\n" + pagination
		}
		if d.Route != "" && d.Meta.Cover != "" {
			alt := d.Meta.CoverAlt
			if alt == "" {
				alt = d.Meta.Title
			}
			d.Body = "![" + alt + "](" + d.Meta.Cover + ")\n\n" + d.Body
		}
		if d.Source == "index.md" {
			d.SchemaItems = postSchemaItems(posts[:min(4, len(posts))])
			list, err := listing(posts[:min(4, len(posts))], d.Source)
			if err != nil {
				return err
			}
			d.Body += "\n\n" + seriesEntry(posts) + "\n\n" + list + "\n\n[Todos os artigos](/blogs/)\n"
		}
		if err := writePage(stage, d); err != nil {
			return err
		}
	}
	list, err := listing(posts, "blogs/index.md")
	if err != nil {
		return err
	}
	if err := generatedPage(stage, "blogs/index.md", "Blog", seriesEntry(posts)+"\n\n"+list, postSchemaItems(posts)...); err != nil {
		return err
	}
	type topic struct {
		label string
		posts []document
	}
	tags := map[string]*topic{}
	for _, post := range posts {
		for _, tag := range post.Meta.Tags {
			slug := tagSlug(tag)
			if tags[slug] == nil {
				tags[slug] = &topic{label: tag}
			}
			tags[slug].posts = append(tags[slug].posts, post)
		}
	}
	slugs := make([]string, 0, len(tags))
	for slug := range tags {
		slugs = append(slugs, slug)
	}
	sort.Strings(slugs)
	var index strings.Builder
	for _, slug := range slugs {
		tag := tags[slug]
		source := "tags/" + slug + "/index.md"
		fmt.Fprintf(&index, "- [%s](/tags/%s/)\n", tag.label, slug)
		list, err := listing(tag.posts, source)
		if err != nil {
			return err
		}
		if err := generatedPage(stage, source, tag.label, list, postSchemaItems(tag.posts)...); err != nil {
			return err
		}
	}
	for _, page := range []struct{ source, title, body string }{
		{"tags/index.md", "Tags", index.String()},
		{"categories/index.md", "Categorias", "As publicações estão organizadas por [tags](/tags/)."},
		{"404.md", "Página não encontrada", "Este endereço não existe. [Voltar ao blog](/blogs/)."},
	} {
		if err := generatedPage(stage, page.source, page.title, page.body); err != nil {
			return err
		}
	}
	// Existing articles contain reviewed HTML/JS demos. Never build untrusted Markdown.
	binary, err := filepath.Abs(filepath.Join(root, ".tools/margo/margo"))
	if err != nil {
		return err
	}
	command := exec.Command(binary, "site", "site.yaml", "--output-dir", ".build/render", "--allow-unsafe-html")
	command.Dir = root
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		return err
	}
	// Retain legacy image URLs even when Margo does not render a reference to them.
	for _, d := range docs {
		assets := filepath.Join(stage, filepath.Dir(d.Source), "assets")
		if _, err := os.Stat(assets); err == nil {
			if err := copyTree(assets, filepath.Join(render, filepath.Dir(d.Source), "assets")); err != nil {
				return err
			}
		} else if !os.IsNotExist(err) {
			return err
		}
	}
	if err := copyTree(filepath.Join(root, "static"), render); err != nil {
		return err
	}
	for _, feed := range []struct {
		file, title string
		posts       []document
	}{
		{"index.xml", "Guigoes", posts}, {"blogs/index.xml", "Guigoes · Blog", posts},
		{"tags/index.xml", "Guigoes · Tags", nil}, {"categories/index.xml", "Guigoes · Categorias", nil},
	} {
		if err := writeFeed(filepath.Join(render, feed.file), feed.title, feed.posts); err != nil {
			return err
		}
	}
	for _, slug := range slugs {
		if err := writeFeed(filepath.Join(render, "tags", slug, "index.xml"), "Guigoes · "+tags[slug].label, tags[slug].posts); err != nil {
			return err
		}
	}
	if err := writeAliases(render, posts); err != nil {
		return err
	}
	if err := writeFile(filepath.Join(render, "robots.txt"), []byte("User-agent: *\nAllow: /\nSitemap: "+baseURL+"/sitemap.xml\n")); err != nil {
		return err
	}
	if err := applySocialCovers(root, render, posts); err != nil {
		return err
	}
	if err := applyBranding(root, render); err != nil {
		return err
	}
	if err := checkSite(root, render, now); err != nil {
		return err
	}
	if err := publishDirectory(render, filepath.Join(root, "public")); err != nil {
		return err
	}
	fmt.Printf("Built %d posts with Margo into public/\n", len(posts))
	return nil
}

var aliasPattern = regexp.MustCompile(`^/[a-z0-9_/-]+/$`)

func writeAliases(render string, posts []document) error {
	var redirects []string
	for _, post := range posts {
		for _, alias := range post.Meta.Aliases {
			if !aliasPattern.MatchString(alias) || strings.Contains(alias, "..") {
				return fmt.Errorf("unsafe alias: %s", alias)
			}
			target := filepath.Join(render, strings.Trim(alias, "/"), "index.html")
			if _, err := os.Stat(target); err == nil {
				return fmt.Errorf("alias collision: %s", alias)
			} else if !os.IsNotExist(err) {
				return err
			}
			canonical := html.EscapeString(baseURL + post.Route)
			if err := writeFile(target, []byte(`<!doctype html><html lang="pt-BR"><head><meta charset="utf-8"><title>Redirecionando</title><link rel="canonical" href="`+canonical+`"><meta http-equiv="refresh" content="0; url=`+canonical+`"></head><body><a href="`+canonical+`">Abrir artigo</a></body></html>`)); err != nil {
				return err
			}
			redirects = append(redirects, alias+" "+post.Route+" 301", strings.TrimSuffix(alias, "/")+" "+post.Route+" 301")
		}
	}
	return writeFile(filepath.Join(render, "_redirects"), []byte(strings.Join(redirects, "\n")+"\n"))
}

func publishDirectory(render, public string) error {
	backup := public + ".previous"
	if _, err := os.Stat(backup); err == nil {
		return fmt.Errorf("previous publication still exists: %s", backup)
	} else if !os.IsNotExist(err) {
		return err
	}
	hadPublic := false
	if _, err := os.Stat(public); err == nil {
		if err := os.Rename(public, backup); err != nil {
			return err
		}
		hadPublic = true
	} else if !os.IsNotExist(err) {
		return err
	}
	if err := os.Rename(render, public); err != nil {
		if hadPublic {
			if restore := os.Rename(backup, public); restore != nil {
				return fmt.Errorf("publish: %v; restore %s: %w", err, backup, restore)
			}
		}
		return err
	}
	if hadPublic {
		return os.RemoveAll(backup)
	}
	return nil
}
