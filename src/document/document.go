package document

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/yoquec/documenter/src/plugins"
)

// represents the output HTML document of documenter
type Document struct {
	Title          string
	Plugins        []plugins.Plugin
	Body           template.HTML
	Template       TemplateInfo
	StyleSheetPath string
}

type TemplateInfo struct {
	Name string
	Path string
}

func New(
	title string,
	plugins []plugins.Plugin,
	body string,
	templateInfo TemplateInfo,
	stylesheetPath string,
) *Document {
	return &Document{
		title,
		plugins,
		template.HTML(body),
		templateInfo,
		stylesheetPath,
	}
}

func (d *Document) Render() ([]byte, error) {
	t, err := template.ParseFiles(d.Template.Path)
	if err != nil {
		return nil, fmt.Errorf("Could not load template '%s': %w", d.Template.Path, err)
	}

	buffer := bytes.Buffer{}

	err = t.ExecuteTemplate(&buffer, d.Template.Name, d)
	if err != nil {
		return nil, fmt.Errorf(
			"Could not hydrate template '%s' with document %v: %w",
			d.Template.Path,
			d, err,
		)
	}

	return buffer.Bytes(), nil
}
