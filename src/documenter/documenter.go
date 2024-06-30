package documenter

import (
	"fmt"
	"os"

	"github.com/yoquec/documenter/src/document"
	"github.com/yoquec/documenter/src/plugins"
)

// Interface for a markdown processor
type Processor interface {
	Render(contents []byte) ([]byte, error)
}

type ResourceProvider interface {
	GetTemplateByName(name string) (*document.TemplateInfo, error)
	GetStyleSheetPath() (string, error)
}

type Documenter struct {
	processor Processor
	plugins   []plugins.Plugin
	resources ResourceProvider
}

func New(engine Processor, plugins []plugins.Plugin, provider ResourceProvider) *Documenter {
	return &Documenter{
		engine,
		plugins,
		provider,
	}
}

func (d *Documenter) RenderDoc(title, path string) ([]byte, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Could not read file '%s': %w", path, err)
	}
	template, err := d.resources.GetTemplateByName("index")
	if err != nil {
		return nil, err
	}
	stylesheetPath, err := d.resources.GetStyleSheetPath()
	if err != nil {
		return nil, err
	}
	html_body, err := d.processor.Render(contents)
    if err != nil {
        return nil, fmt.Errorf("Could not render markdown file '%s': %w", path, err)
    }

	doc := document.New(title, d.plugins, string(html_body), *template, stylesheetPath)
	render, err := doc.Render()
	if err != nil {
		return nil, fmt.Errorf("Could not render document: %w", err)
	}

	return render, nil
}
