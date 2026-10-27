package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

const brandHead = `<link rel="icon" type="image/svg+xml" sizes="any" href="/assets/brand/icon-v2.svg">
<link rel="icon" type="image/png" sizes="16x16" href="/assets/brand/icon-16-v2.png">
<link rel="icon" type="image/png" sizes="32x32" href="/assets/brand/icon-32-v2.png">
<link rel="icon" type="image/png" sizes="48x48" href="/assets/brand/icon-48-v2.png">
<link rel="icon" type="image/png" sizes="96x96" href="/assets/brand/icon-96-v2.png">
<link rel="shortcut icon" type="image/x-icon" href="/assets/brand/favicon-v2.ico">
<link rel="apple-touch-icon" sizes="152x152" href="/assets/brand/apple-touch-icon-152-v2.png">
<link rel="apple-touch-icon" sizes="167x167" href="/assets/brand/apple-touch-icon-167-v2.png">
<link rel="apple-touch-icon" sizes="180x180" href="/assets/brand/apple-touch-icon-180-v2.png">
<link rel="manifest" href="/site.webmanifest" type="application/manifest+json">
<meta name="application-name" content="Guigoes">
<meta name="apple-mobile-web-app-title" content="Guigoes">
<meta name="theme-color" content="#007d9c">`

var generatedIcon = regexp.MustCompile(`<link\b[^>]*\brel="icon"[^>]*>`)

func addBrandHead(data []byte) ([]byte, error) {
	head, body, ok := bytes.Cut(data, []byte("</head>"))
	if !ok {
		return nil, fmt.Errorf("brand metadata: missing head")
	}
	if bytes.Contains(head, []byte(`rel="manifest"`)) {
		return nil, fmt.Errorf("brand metadata already present")
	}
	// Change only the head; Margo's document digest still identifies the article bytes.
	head = generatedIcon.ReplaceAll(head, nil)
	return append(append(append(head, []byte(brandHead)...), []byte("</head>")...), body...), nil
}

func applyBranding(root, render string) error {
	if err := copyTree(filepath.Join(root, "assets/brand"), filepath.Join(render, "assets/brand")); err != nil {
		return err
	}
	for _, file := range []struct{ source, target string }{
		{"favicon-v2.ico", "favicon.ico"},
		{"apple-touch-icon-180-v2.png", "apple-touch-icon.png"},
	} {
		data, err := os.ReadFile(filepath.Join(root, "assets/brand", file.source))
		if err != nil {
			return err
		}
		if err := writeFile(filepath.Join(render, file.target), data); err != nil {
			return err
		}
	}
	err := filepath.WalkDir(render, func(file string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(file) != ".html" {
			return nil
		}
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		data, err = addBrandHead(data)
		if err != nil {
			return fmt.Errorf("%s: %w", file, err)
		}
		return writeFile(file, data)
	})
	if err != nil {
		return err
	}
	return refreshBrandManifest(render)
}

// Margo's artifact manifest includes exact output hashes, so head edits must update it.
func refreshBrandManifest(render string) error {
	file := filepath.Join(render, "margo-manifest.json")
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	var manifest map[string]json.RawMessage
	if err := json.Unmarshal(data, &manifest); err != nil {
		return err
	}
	type entry struct {
		Path   string `json:"path"`
		Digest string `json:"digest"`
	}
	var entries []entry
	err = filepath.WalkDir(render, func(file string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Base(file) == "margo-manifest.json" {
			return nil
		}
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(render, file)
		if err != nil {
			return err
		}
		digest := sha256.Sum256(data)
		entries = append(entries, entry{filepath.ToSlash(rel), hex.EncodeToString(digest[:])})
		return nil
	})
	if err != nil {
		return err
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Path < entries[j].Path })
	// Margo hashes digest byte arrays internally, although its public JSON uses hex.
	type digestEntry struct {
		Path   string   `json:"path"`
		Digest [32]byte `json:"digest"`
	}
	canonical := struct {
		Entries []digestEntry `json:"entries"`
	}{Entries: make([]digestEntry, len(entries))}
	for i, item := range entries {
		raw, err := hex.DecodeString(item.Digest)
		if err != nil {
			return err
		}
		canonical.Entries[i].Path = item.Path
		copy(canonical.Entries[i].Digest[:], raw)
	}
	var encoded bytes.Buffer
	encoder := json.NewEncoder(&encoded)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(canonical); err != nil {
		return err
	}
	digest := sha256.Sum256(append([]byte("margo/manifest/v1\n"), bytes.TrimSuffix(encoded.Bytes(), []byte("\n"))...))
	manifest["digest"], err = json.Marshal(hex.EncodeToString(digest[:]))
	if err != nil {
		return err
	}
	manifest["entries"], err = json.Marshal(entries)
	if err != nil {
		return err
	}
	output, err := json.Marshal(manifest)
	if err != nil {
		return err
	}
	return writeFile(file, append(output, '\n'))
}
