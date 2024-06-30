package documenter

import (
	"fmt"
	"os"

	"github.com/yoquec/documenter/src/document"
	"github.com/yoquec/documenter/src/plugins"
)

type Engine interface {
	RenderToString(contents []byte) string
}

type ResourceProvider interface {
	GetTemplateByName(name string) (*document.TemplateInfo, error)
	GetStyleSheetPath() (string, error)
}

type Documenter struct {
	engine    Engine
	plugins   []plugins.Plugin
	resources ResourceProvider
}

func New(engine Engine, plugins []plugins.Plugin, provider ResourceProvider) *Documenter {
	return &Documenter{
		engine,
		plugins,
		provider,
	}
}

func (d *Documenter) RenderDoc(title, path string) ([]byte, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return *new([]byte), fmt.Errorf("Could not read file '%s': %w", path, err)
	}
	template, err := d.resources.GetTemplateByName("index")
	if err != nil {
		return *new([]byte), err
	}
	stylesheetPath, err := d.resources.GetStyleSheetPath()
    if err != nil {
        return *new([]byte), err
    }

	html_body := d.engine.RenderToString(contents)
	doc := document.New(title, d.plugins, html_body, *template, stylesheetPath)

    render, err := doc.Render()
    if err != nil {
        return *new([]byte), fmt.Errorf("Could not render document: %w", err)
    }

    return render, nil
}
