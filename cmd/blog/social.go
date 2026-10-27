package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func coverMetadata(root string, post document) (map[string]string, error) {
	f, err := os.Open(filepath.Join(root, "content", filepath.Dir(post.Source), post.Meta.Cover))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	cfg, format, err := image.DecodeConfig(f)
	if err != nil {
		return nil, fmt.Errorf("%s: decode social cover: %w", post.Source, err)
	}
	alt := post.Meta.CoverAlt
	if alt == "" {
		alt = post.Meta.Title
	}
	url := structuredData(post).Image
	return map[string]string{
		"og:image": url, "og:image:alt": alt, "og:image:type": "image/" + format,
		"og:image:width": strconv.Itoa(cfg.Width), "og:image:height": strconv.Itoa(cfg.Height),
		"twitter:image": url, "twitter:image:alt": alt,
	}, nil
}

func replaceSocialImages(data []byte, values map[string]string) ([]byte, error) {
	head, body, ok := bytes.Cut(data, []byte("</head>"))
	if !ok {
		return nil, fmt.Errorf("social metadata: missing head")
	}
	var output bytes.Buffer
	tokens := html.NewTokenizer(bytes.NewReader(head))
	for {
		kind := tokens.Next()
		if kind == html.ErrorToken {
			if tokens.Err() != io.EOF {
				return nil, tokens.Err()
			}
			break
		}
		raw := append([]byte(nil), tokens.Raw()...)
		remove := false
		if kind == html.StartTagToken || kind == html.SelfClosingTagToken {
			token := tokens.Token()
			if token.Data == "meta" {
				for _, attr := range token.Attr {
					if attr.Key == "property" || attr.Key == "name" {
						key := strings.ToLower(attr.Val)
						remove = remove || key == "og:image" || strings.HasPrefix(key, "og:image:") || key == "twitter:image" || strings.HasPrefix(key, "twitter:image:")
					}
				}
			}
		}
		if !remove {
			output.Write(raw)
		}
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	for _, key := range keys {
		attr := "property"
		if strings.HasPrefix(key, "twitter:") {
			attr = "name"
		}
		fmt.Fprintf(&output, `<meta %s="%s" content="%s">`, attr, key, html.EscapeString(values[key]))
	}
	output.WriteString("</head>")
	output.Write(body)
	return output.Bytes(), nil
}

// Margo v0.0.25 supplies one site-wide social image; replace only covered posts.
func applySocialCovers(root, render string, posts []document) error {
	for _, post := range posts {
		if post.Meta.Cover == "" {
			continue
		}
		values, err := coverMetadata(root, post)
		if err != nil {
			return err
		}
		file := filepath.Join(render, filepath.Dir(post.Source), "index.html")
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		data, err = replaceSocialImages(data, values)
		if err != nil {
			return fmt.Errorf("%s: %w", post.Source, err)
		}
		if err := writeFile(file, data); err != nil {
			return err
		}
	}
	return nil
}
