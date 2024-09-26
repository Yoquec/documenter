package content

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
)

const TEMPLATE_NAME = "index"

//go:embed all:assets
var embeddedFS embed.FS

type embeddedAssets struct {
	fs.FS
	stylesheetInfo *stylesheetInfo
	templateInfo   *templateInfo
}

type stylesheetInfo struct {
	path string
}

type templateInfo struct {
	path string
	name string
}

// Allows access to the static files embedded on the binary
func getEmbeddedAssets() *embeddedAssets {
	fs, _ := fs.Sub(embeddedFS, "assets")
	stylesheetInfo := &stylesheetInfo{
		path.Join("css", "style.css"),
	}
	templateInfo := &templateInfo{
		path: path.Join("templates", "index.html"),
		name: TEMPLATE_NAME,
	}

	return &embeddedAssets{
		fs,
		stylesheetInfo,
		templateInfo,
	}
}

// Attempts to open a file within the embedded asset file-system
func (ea *embeddedAssets) Open(path string) ([]byte, error) {
	file, _ := ea.FS.Open(path)
	contents := make([]byte, 0)
	if _, err := file.Read(contents); err != nil {
		return nil, fmt.Errorf("Could not read embedded file %q: %w", path, err)
	}
	return contents, nil
}
