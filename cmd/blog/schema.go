package main

import (
	"encoding/json"
	"net/url"
	"path"
)

type person struct {
	Type string `json:"@type"`
	Name string `json:"name"`
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
		data.Authors = append(data.Authors, person{Type: "Person", Name: author})
	}
	if d.Meta.Cover != "" {
		data.Image = baseURL + (&url.URL{Path: path.Join(d.Route, d.Meta.Cover)}).String()
	}
	return data
}

func schemaScript(d document) (string, error) {
	// encoding/json escapes '<', including '</script>' embedded in metadata.
	data, err := json.Marshal(structuredData(d))
	if err != nil {
		return "", err
	}
	return `<script type="application/ld+json">` + string(data) + `</script>`, nil
}
