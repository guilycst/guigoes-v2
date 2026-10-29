package main

import (
	"os"
	"path/filepath"
)

// Retain published brand assets and legacy default icon URLs.
func copyBrandAssets(root, render string) error {
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
	return nil
}
