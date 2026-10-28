package main

import (
	"fmt"
	"html"
	"sort"
	"strings"
)

// posts contains only articles included in this build, never unpublished parts.
func seriesPosts(name string, posts []document) []document {
	var parts []document
	for _, post := range posts {
		if name != "" && post.Meta.Series == name {
			parts = append(parts, post)
		}
	}
	sort.Slice(parts, func(i, j int) bool { return parts[i].Meta.SeriesPart < parts[j].Meta.SeriesPart })
	return parts
}

func seriesNavigation(current document, posts []document) (string, string) {
	parts := seriesPosts(current.Meta.Series, posts)
	if len(parts) == 0 {
		return "", ""
	}
	var index, pagination strings.Builder
	fmt.Fprintf(&index, `<nav class="guigoes-series" aria-label="Índice da série"><p>Série %s · Parte %d</p><ol>`, html.EscapeString(current.Meta.Series), current.Meta.SeriesPart)
	pagination.WriteString(`<nav class="guigoes-series-pagination" aria-label="Continuar a série">`)
	for i, part := range parts {
		label := strings.TrimPrefix(part.Meta.Title, fmt.Sprintf("%s, parte %d: ", part.Meta.Series, part.Meta.SeriesPart))
		active := ""
		if part.Route == current.Route {
			active = ` aria-current="page"`
			for _, adjacent := range []struct {
				offset     int
				label, rel string
			}{{-1, "Anterior", "prev"}, {1, "Próximo", "next"}} {
				j := i + adjacent.offset
				if j >= 0 && j < len(parts) {
					fmt.Fprintf(&pagination, `<a rel="%s" href="%s">%s: %s</a>`, adjacent.rel, html.EscapeString(parts[j].Route), adjacent.label, html.EscapeString(parts[j].Meta.Title))
				}
			}
		}
		fmt.Fprintf(&index, `<li value="%d"><a href="%s"%s>%s</a></li>`, part.Meta.SeriesPart, html.EscapeString(part.Route), active, html.EscapeString(label))
	}
	index.WriteString("</ol></nav>\n\n")
	if len(parts) == 1 {
		return index.String(), ""
	}
	pagination.WriteString("</nav>\n")
	return index.String(), pagination.String()
}

func seriesEntry(posts []document) string {
	if len(posts) == 0 {
		return ""
	}
	for _, part := range seriesPosts(posts[0].Meta.Series, posts) {
		if part.Meta.SeriesPart == 1 {
			return fmt.Sprintf(`<aside class="guigoes-series-entry" aria-label="Começar a série"><span>Série %s</span><a href="%s">Comece pela parte 1: %s</a></aside>`, html.EscapeString(part.Meta.Series), html.EscapeString(part.Route), html.EscapeString(strings.TrimPrefix(part.Meta.Title, part.Meta.Series+", parte 1: ")))
		}
	}
	return ""
}
