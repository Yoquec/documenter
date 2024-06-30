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
	Body           string
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
	template TemplateInfo,
	stylesheetPath string,
) *Document {
	return &Document{
		title,
		plugins,
		body,
		template,
		stylesheetPath,
	}
}

func (d *Document) Render() ([]byte, error) {
	t, err := template.ParseFiles(d.Template.Path)
	if err != nil {
		return *new([]byte), fmt.Errorf("Could not load template '%s': %w", d.Template.Path, err)
	}

	buffer := bytes.Buffer{}

	err = t.ExecuteTemplate(&buffer, d.Template.Name, d)
	if err != nil {
		return *new([]byte), fmt.Errorf(
			"Could not hydrate template '%s' with document %v: %w",
			d.Template.Path,
			d, err,
		)
	}

	return buffer.Bytes(), nil
}
