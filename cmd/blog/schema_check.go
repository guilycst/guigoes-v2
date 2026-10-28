package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"slices"
	"strings"
)

func isAlias(route string, posts []document) bool {
	for _, post := range posts {
		if slices.Contains(post.Meta.Aliases, route) {
			return true
		}
	}
	return false
}

func checkNonPostSchema(page htmlPage, route string, posts []document) error {
	if len(page.Schemas) != 1 || page.HeadSchemas != 1 || page.CanonicalCount != 1 {
		return fmt.Errorf("expected one JSON-LD block and one canonical in head")
	}
	var data map[string]any
	if err := json.Unmarshal(page.Schemas[0], &data); err != nil {
		return fmt.Errorf("invalid JSON-LD: %w", err)
	}
	if data["@context"] != "https://schema.org" {
		return fmt.Errorf("invalid JSON-LD context")
	}
	if route == "/" {
		graph, ok := data["@graph"].([]any)
		if !ok || len(graph) != 2 {
			return fmt.Errorf("home requires WebSite and WebPage")
		}
		website, ok := graph[0].(map[string]any)
		if !ok || website["@type"] != "WebSite" || website["@id"] != baseURL+"/#website" || website["name"] != "Guigoes" || website["url"] != baseURL+"/" {
			return fmt.Errorf("incorrect WebSite identity")
		}
		data, ok = graph[1].(map[string]any)
		if !ok {
			return fmt.Errorf("missing home WebPage")
		}
	}
	kind := "WebPage"
	switch {
	case route == "/about/":
		kind = "AboutPage"
	case route == "/contact/":
		kind = "ContactPage"
	case route == "/blogs/", route == "/projects/", strings.HasPrefix(route, "/tags/"):
		kind = "CollectionPage"
	}
	if data["@type"] != kind || data["url"] != baseURL+route || data["@id"] != baseURL+route+"#webpage" || data["inLanguage"] != "pt-BR" || page.Canonical != baseURL+route {
		return fmt.Errorf("incorrect page type, identity or locale")
	}
	if titles := page.HeadMetadata["og:title"]; len(titles) != 1 || data["name"] != strings.TrimSuffix(titles[0], " · Guigoes") {
		return fmt.Errorf("JSON-LD name differs from page title")
	}
	if description, exists := data["description"]; exists {
		values := page.HeadMetadata["description"]
		if len(values) != 1 || description != values[0] {
			return fmt.Errorf("JSON-LD description differs from page description")
		}
	}
	parent, ok := data["isPartOf"].(map[string]any)
	if !ok || parent["@id"] != baseURL+"/#website" {
		return fmt.Errorf("missing WebSite reference")
	}
	if route == "/about/" {
		person, ok := data["mainEntity"].(map[string]any)
		if !ok || person["@type"] != "Person" || person["@id"] != baseURL+"/about/#person" || person["name"] != "Guilherme de Castro" || person["url"] != baseURL+"/about/" {
			return fmt.Errorf("incorrect author identity")
		}
	}
	var listed []document
	switch {
	case route == "/":
		listed = posts[:min(4, len(posts))]
	case route == "/blogs/":
		listed = posts
	case strings.HasPrefix(route, "/tags/") && route != "/tags/":
		listed = make([]document, 0)
		slug := strings.TrimSuffix(strings.TrimPrefix(route, "/tags/"), "/")
		for _, post := range posts {
			for _, tag := range post.Meta.Tags {
				if tagSlug(tag) == slug {
					listed = append(listed, post)
					break
				}
			}
		}
	default:
		return nil
	}
	list, ok := data["mainEntity"].(map[string]any)
	if !ok || list["@type"] != "ItemList" || list["numberOfItems"] != float64(len(listed)) {
		return fmt.Errorf("incorrect list count")
	}
	raw, err := json.Marshal(list["itemListElement"])
	if err != nil {
		return err
	}
	var items []schemaListItem
	if err := json.Unmarshal(raw, &items); err != nil {
		return err
	}
	if !reflect.DeepEqual(items, postSchemaItems(listed)) {
		return fmt.Errorf("structured list differs from visible post selection")
	}
	for _, item := range items {
		if !slices.Contains(page.Refs, strings.TrimPrefix(item.URL, baseURL)) && !slices.Contains(page.Refs, item.URL) {
			return fmt.Errorf("structured list item missing from page: %s", item.URL)
		}
	}
	return nil
}

// Expected projection for checking the renderer against source metadata.
type person struct {
	Type string `json:"@type"`
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

type blogPosting struct {
	Context          string   `json:"@context"`
	Type             string   `json:"@type"`
	ID               string   `json:"@id"`
	URL              string   `json:"url"`
	MainEntityOfPage string   `json:"mainEntityOfPage"`
	Headline         string   `json:"headline"`
	Description      string   `json:"description,omitempty"`
	Language         string   `json:"inLanguage,omitempty"`
	Authors          []person `json:"author"`
	Published        string   `json:"datePublished"`
	Modified         string   `json:"dateModified,omitempty"`
	Image            string   `json:"image,omitempty"`
}

func structuredData(d document) blogPosting {
	canonical := baseURL + (&url.URL{Path: d.Route}).String()
	data := blogPosting{Context: "https://schema.org", Type: "BlogPosting", ID: canonical + "#blogposting", URL: canonical, MainEntityOfPage: canonical, Headline: d.Meta.Title, Description: d.Meta.Description, Language: d.Meta.Language, Published: d.Meta.PublishedAt, Modified: d.Meta.ModifiedAt}
	for _, author := range d.Meta.Authors {
		p := person{Type: "Person", Name: author, URL: authorURL(author)}
		data.Authors = append(data.Authors, p)
	}
	if d.Meta.Cover != "" {
		data.Image = coverURL(d)
	}
	return data
}
