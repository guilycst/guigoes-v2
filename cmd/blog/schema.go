package main

import (
	"net/url"
	"path"
	"strings"
)

func authorURL(name string) string {
	if name == "Guilherme de Castro" {
		return baseURL + "/about/"
	}
	return ""
}

// pageStructuredData declares editorial meaning; Margo owns JSON-LD emission.
func pageStructuredData(d document) any {
	if d.Route != "" {
		options := map[string]any{"type": "BlogPosting"}
		if image := coverURL(d); image != "" {
			options["image"] = image
		}
		profiles := map[string]string{}
		for _, author := range d.Meta.Authors {
			if profile := authorURL(author); profile != "" {
				profiles[author] = profile
			}
		}
		if len(profiles) > 0 {
			options["author_urls"] = profiles
		}
		return options
	}
	if d.Source == "404.md" {
		return nil
	}
	route := "/" + strings.TrimSuffix(d.Source, "index.md")
	canonical := baseURL + (&url.URL{Path: route}).String()
	kind := "WebPage"
	switch {
	case d.Source == "about/index.md":
		kind = "AboutPage"
	case d.Source == "contact/index.md":
		kind = "ContactPage"
	case d.Source == "projects/index.md", d.Source == "blogs/index.md", strings.HasPrefix(d.Source, "tags/"):
		kind = "CollectionPage"
	}
	page := map[string]any{
		"@context": "https://schema.org", "@type": kind, "@id": canonical + "#webpage", "url": canonical,
		"name": d.Meta.Title, "inLanguage": "pt-BR", "isPartOf": map[string]string{"@id": baseURL + "/#website"},
	}
	if d.Meta.Description != "" {
		page["description"] = d.Meta.Description
	}
	if d.Source == "about/index.md" {
		page["mainEntity"] = map[string]string{"@type": "Person", "@id": baseURL + "/about/#person", "name": "Guilherme de Castro", "url": baseURL + "/about/"}
	}
	if d.SchemaItems != nil {
		page["mainEntity"] = map[string]any{"@type": "ItemList", "numberOfItems": len(d.SchemaItems), "itemListElement": d.SchemaItems}
	}
	if d.Source == "index.md" {
		delete(page, "@context")
		return map[string]any{"@context": "https://schema.org", "@graph": []any{
			map[string]any{"@type": "WebSite", "@id": baseURL + "/#website", "url": baseURL + "/", "name": "Guigoes", "inLanguage": "pt-BR"}, page,
		}}
	}
	return page
}

type schemaListItem struct {
	Type     string `json:"@type"`
	Position int    `json:"position"`
	Name     string `json:"name"`
	URL      string `json:"url"`
}

func postSchemaItems(posts []document) []schemaListItem {
	items := make([]schemaListItem, 0, len(posts))
	for _, post := range posts {
		items = append(items, schemaListItem{Type: "ListItem", Position: len(items) + 1, Name: post.Meta.Title, URL: baseURL + post.Route})
	}
	return items
}

func coverURL(d document) string {
	if d.Meta.Cover == "" {
		return ""
	}
	return baseURL + (&url.URL{Path: path.Join(d.Route, d.Meta.Cover)}).String()
}
