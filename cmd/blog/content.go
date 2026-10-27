package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/unicode/norm"
)

const baseURL = "https://guigoes.com"
const nav = `<nav class="guigoes-nav" aria-label="Navegação principal"><a class="guigoes-brand" href="/">guigoes<span aria-hidden="true">.</span></a><span class="guigoes-nav-links"><a href="/blogs/">Blog</a><a href="/projects/">Projetos</a><a href="/about/">Sobre</a><a href="/contact/">Contato</a><a href="/index.xml">RSS</a>` + themeToggle + `</span></nav>`

const themeToggle = `<button type="button" class="guigoes-theme-toggle" data-guigoes-color-mode role="switch" aria-label="Tema escuro" aria-checked="false" hidden><svg class="sun" aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"><circle cx="12" cy="12" r="4"/><path d="M12 2v2m0 16v2M2 12h2m16 0h2M5 5l1.5 1.5m11 11L19 19M5 19l1.5-1.5m11-11L19 5"/></svg><svg class="moon" aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"><path d="M20.5 13A9 9 0 0 1 11 3.5 9 9 0 1 0 20.5 13Z"/></svg></button>`

//go:embed theme.js
var themeScript string

const footer = `<footer class="guigoes-footer"><span>Guilherme de Castro</span><span><a class="guigoes-footer-social" href="https://github.com/guilycst"><img class="guigoes-social-icon" src="%[1]s/github.svg" alt="" width="16" height="16">GitHub</a> · <a class="guigoes-footer-social" href="https://www.linkedin.com/in/guicas"><img class="guigoes-social-icon" src="%[1]s/linkedin.svg" alt="" width="16" height="16">LinkedIn</a><br>Publicado com <a href="https://github.com/araihu/margo">Margo</a></span></footer>`

type metadata struct {
	Title       string   `toml:"title"`
	Description string   `toml:"description"`
	Language    string   `toml:"language"`
	Authors     []string `toml:"authors"`
	PublishedAt string   `toml:"publishedAt"`
	ModifiedAt  string   `toml:"modifiedAt"`
	Cover       string   `toml:"cover"`
	CoverAlt    string   `toml:"cover_alt"`
	Tags        []string `toml:"tags"`
	Aliases     []string `toml:"aliases"`
	Draft       bool     `toml:"draft"`
}

type document struct {
	Meta                metadata
	Front               map[string]any
	Body, Source, Route string
	Published           time.Time
}

func parseDocument(data []byte) (document, error) {
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	if !strings.HasPrefix(text, "+++\n") {
		return document{}, fmt.Errorf("missing TOML frontmatter")
	}
	front, body, ok := strings.Cut(text[4:], "\n+++\n")
	if !ok {
		return document{}, fmt.Errorf("unclosed TOML frontmatter")
	}
	d := document{Body: strings.TrimLeft(body, "\n")}
	if err := toml.Unmarshal([]byte(front), &d.Meta); err != nil {
		return d, err
	}
	if err := toml.Unmarshal([]byte(front), &d.Front); err != nil {
		return d, err
	}
	if strings.TrimSpace(d.Meta.Title) == "" {
		return d, fmt.Errorf("missing title")
	}
	var err error
	if d.Meta.PublishedAt != "" {
		d.Published, err = time.Parse(time.RFC3339, d.Meta.PublishedAt)
		if err != nil {
			return d, fmt.Errorf("publishedAt must be RFC 3339 with timezone: %w", err)
		}
	}
	if d.Meta.ModifiedAt != "" {
		modified, err := time.Parse(time.RFC3339, d.Meta.ModifiedAt)
		if err != nil {
			return d, fmt.Errorf("invalid modifiedAt: %w", err)
		}
		if modified.Before(d.Published) {
			return d, fmt.Errorf("modifiedAt precedes publishedAt")
		}
	}
	return d, nil
}

func (d document) visible(now time.Time) bool {
	return !d.Meta.Draft && (d.Published.IsZero() || !d.Published.After(now))
}

func loadDocuments(root string, now time.Time) ([]document, []document, error) {
	var docs, posts []document
	err := filepath.WalkDir(filepath.Join(root, "content"), func(file string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("content symlink unsupported: %s", file)
		}
		if entry.IsDir() || filepath.Ext(file) != ".md" {
			return nil
		}
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		d, err := parseDocument(data)
		if err != nil {
			return fmt.Errorf("%s: %w", file, err)
		}
		if !d.visible(now) {
			return nil
		}
		relative, err := filepath.Rel(filepath.Join(root, "content"), file)
		if err != nil {
			return err
		}
		d.Source = filepath.ToSlash(relative)
		if strings.HasPrefix(d.Source, "blogs/") {
			if path.Base(d.Source) != "index.md" || path.Dir(d.Source) == "blogs" {
				return fmt.Errorf("post must be blogs/<slug>/index.md: %s", file)
			}
			d.Route = "/" + path.Dir(d.Source) + "/"
			if d.Published.IsZero() || len(d.Meta.Authors) == 0 {
				return fmt.Errorf("%s: posts require publishedAt and authors", file)
			}
			for _, author := range d.Meta.Authors {
				if strings.TrimSpace(author) == "" {
					return fmt.Errorf("%s: empty author", file)
				}
			}
			if d.Meta.Cover != "" {
				cover := strings.TrimPrefix(d.Meta.Cover, "./")
				if !strings.HasPrefix(cover, "assets/") || !fs.ValidPath(cover) {
					return fmt.Errorf("%s: cover must be a local assets/ path", file)
				}
				info, err := os.Stat(filepath.Join(filepath.Dir(file), filepath.FromSlash(cover)))
				if err != nil {
					return fmt.Errorf("%s: cover: %w", file, err)
				}
				if !info.Mode().IsRegular() {
					return fmt.Errorf("%s: cover is not a regular file", file)
				}
			}
			for _, tag := range d.Meta.Tags {
				if tagSlug(tag) == "" {
					return fmt.Errorf("%s: tag has no URL slug: %q", file, tag)
				}
			}
			posts = append(posts, d)
		}
		docs = append(docs, d)
		return nil
	})
	sort.Slice(posts, func(i, j int) bool {
		if posts[i].Published.Equal(posts[j].Published) {
			return posts[i].Route > posts[j].Route
		}
		return posts[i].Published.After(posts[j].Published)
	})
	return docs, posts, err
}

var slugSeparators = regexp.MustCompile(`[^a-z0-9]+`)

func tagSlug(tag string) string {
	var ascii strings.Builder
	for _, r := range norm.NFKD.String(tag) {
		if r < 128 {
			ascii.WriteRune(r)
		}
	}
	return strings.Trim(slugSeparators.ReplaceAllString(strings.ToLower(ascii.String()), "-"), "-")
}

func publicationDetails(m metadata) string {
	authors := make([]string, len(m.Authors))
	for i, author := range m.Authors {
		authors[i] = `<span rel="author">` + html.EscapeString(author) + `</span>`
		if profile := authorURL(author); profile != "" {
			authors[i] = `<a rel="author" href="` + html.EscapeString(profile) + `">` + html.EscapeString(author) + `</a>`
		}
	}
	var dates, tags []string
	for _, item := range []struct{ value, label string }{{m.PublishedAt, "Publicado"}, {m.ModifiedAt, "Atualizado"}} {
		if item.value != "" {
			date, _ := time.Parse(time.RFC3339, item.value)
			dates = append(dates, item.label+` <time datetime="`+html.EscapeString(item.value)+`">`+date.Format("02/01/2006")+`</time>`)
		}
	}
	for _, tag := range m.Tags {
		tags = append(tags, `<li><a href="/tags/`+tagSlug(tag)+`/">`+html.EscapeString(tag)+`</a></li>`)
	}
	// Margo's marker preserves this header while still generating head metadata.
	return `<header class="margo-document__metadata" data-margo-publication-metadata="true" aria-label="Sobre este artigo"><address>` + strings.Join(authors, ", ") + `</address><span>` + strings.Join(dates, " · ") + `</span><ul aria-label="Tags">` + strings.Join(tags, "") + `</ul></header>`
}

func listing(posts []document, source string) (string, error) {
	var cards strings.Builder
	cards.WriteString(`<ol class="guigoes-posts" aria-label="Artigos">`)
	for _, post := range posts {
		m := post.Meta
		cards.WriteString(`<li class="guigoes-post">`)
		if m.Cover != "" {
			asset := path.Join(strings.TrimPrefix(post.Route, "/"), m.Cover)
			src, err := filepath.Rel(filepath.FromSlash(path.Dir(source)), filepath.FromSlash(asset))
			if err != nil {
				return "", err
			}
			fmt.Fprintf(&cards, `<a class="guigoes-post-image" href="%s" tabindex="-1" aria-hidden="true"><img src="%s" alt="%s" loading="lazy" width="640" height="360"></a>`, html.EscapeString(post.Route), html.EscapeString(filepath.ToSlash(src)), html.EscapeString(m.CoverAlt))
		}
		fmt.Fprintf(&cards, `<div class="guigoes-post-text"><time datetime="%s">%s</time><h2><a href="%s">%s</a></h2><p>%s</p></div></li>`, html.EscapeString(m.PublishedAt), post.Published.Format("02/01/2006"), html.EscapeString(post.Route), html.EscapeString(m.Title), html.EscapeString(m.Description))
	}
	cards.WriteString("</ol>\n")
	return cards.String(), nil
}

func writePage(stage string, d document) error {
	body := d.Body
	heading := "# " + d.Meta.Title + "\n\n"
	if strings.HasPrefix(body, "# ") {
		var first string
		first, body, _ = strings.Cut(body, "\n")
		heading = first + "\n\n"
	}
	if d.Meta.PublishedAt != "" {
		body = publicationDetails(d.Meta) + "\n\n" + body
	}
	if d.Route != "" {
		schema, err := schemaScript(d)
		if err != nil {
			return err
		}
		body += "\n\n" + schema + "\n"
	}
	class := "guigoes-content"
	if d.Source == "index.md" {
		class = "guigoes-home"
	}
	front, err := json.Marshal(d.Front)
	if err != nil {
		return err
	}
	iconDir, err := filepath.Rel(filepath.FromSlash(path.Dir(d.Source)), filepath.FromSlash("contact/assets/icons"))
	if err != nil {
		return err
	}
	pageFooter := fmt.Sprintf(footer, html.EscapeString(filepath.ToSlash(iconDir)))
	return writeFile(filepath.Join(stage, filepath.FromSlash(d.Source)), []byte("---\n"+string(front)+"\n---\n\n<script>"+themeScript+"</script>\n\n"+nav+"\n\n<div class=\""+class+"\">\n\n"+heading+body+"\n\n</div>\n\n"+pageFooter+"\n"))
}
