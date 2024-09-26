package content

import (
	"bytes"
	"fmt"
	"html/template"
	"iter"
	"os"
	"strings"

	"github.com/yoquec/documenter/documenter/plugins"
	"gitlab.com/golang-commonmark/markdown"
)

type Content struct {
	Metadata *Metadata
	Buffer   *bytes.Buffer
}

type MarkdownContent struct {
	*Content
}

type HtmlContent struct {
	*Content
}

type Metadata struct {
	Title string
	Tags  []string
}

// Data needed to render the output
// HTML template
type renderData struct {
	Title   string
	Style   template.CSS
	Plugins []plugins.Plugin
	Body    template.HTML
}

// Interface for a construct that executes
// a transformation over markdown contents
type markdownTransform interface {
	Execute(*MarkdownContent)
}

// Opens markdown content from the specified path
func FromPath(path string) (*MarkdownContent, error) {
	file, err := os.ReadFile(path)
	if err != nil {
        return nil, fmt.Errorf("Could not read file from path %q: %w", path, err)
	}
	metadata := newEmptyDocumenterMetadata()
	metadata.Title = strings.TrimSuffix(path, ".md")

	return &MarkdownContent{&Content{metadata, bytes.NewBuffer(file)}}, nil
}

func FromString(contents string) *MarkdownContent {
	return &MarkdownContent{
		&Content{
			newEmptyDocumenterMetadata(),
			bytes.NewBufferString(contents),
		},
	}
}

func newEmptyDocumenterMetadata() *Metadata {
	return &Metadata{"", nil}
}

// Returns an iterator with the lines of the file content
func (c *Content) ReadLines() iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		data := bytes.Clone(c.Buffer.Bytes())
		for len(data) > 0 {
			line, rest, _ := bytes.Cut(data, []byte("\n"))
			if !yield(line) {
				return
			}
			data = rest
			// Case for where there is a single line in the contents
			if data == nil {
				return
			}
		}
	}
}

// Applies a transformation to the markdown content
func (md *MarkdownContent) Apply(transform markdownTransform) {
	transform.Execute(md)
}

// Converts contents into HTML
func (md *MarkdownContent) ToHtml() (*HtmlContent, error) {
	engine := markdown.New(
		markdown.Tables(true),  // Render GFM tables
		markdown.Linkify(true), // Generate links for URLs automatically
		markdown.HTML(true),    // Ignore HTML inside the markdown
		markdown.Typographer(true),
	)

	buffer := bytes.NewBuffer(*new([]byte))

	if err := engine.Render(buffer, md.Content.Buffer.Bytes()); err != nil {
		return nil, fmt.Errorf("Could not render HTML from Markdown source: %w", err)
	}

	return &HtmlContent{
		&Content{
			Metadata: md.Metadata,
			Buffer:   buffer,
		},
	}, nil
}

// Renders HTML content into a document including all
// available plugins
func (html *HtmlContent) Render() (*bytes.Buffer, error) {
	return html.RenderWithPlugins(plugins.All)
}

func (h *HtmlContent) RenderWithPlugins(pluginList []plugins.Plugin) (*bytes.Buffer, error) {
	embededAssets := getEmbeddedAssets()

	stylesheetContents, err := embededAssets.Open(embededAssets.stylesheetInfo.path)
	if err != nil {
		return nil, fmt.Errorf(
			"Could not open stylesheet %q: %w",
			embededAssets.stylesheetInfo.path,
			err,
		)
	}

	t, err := template.ParseFS(embededAssets.FS, embededAssets.templateInfo.path)
	if err != nil {
		return nil, fmt.Errorf(
			"Could not open template %q: %w",
			embededAssets.templateInfo.path,
			err,
		)
	}

	data := renderData{
		Title:   h.Metadata.Title,
		Style:   template.CSS(stylesheetContents),
		Plugins: pluginList,
		Body:    template.HTML(h.Content.Buffer.Bytes()),
	}

	buffer := bytes.NewBuffer(*new([]byte))
	t.ExecuteTemplate(buffer, embededAssets.templateInfo.name, data)

	return buffer, nil
}

func (df *HtmlContent) RenderPdf() (*bytes.Buffer, error) {
	return df.RenderPdfWithPlugins(plugins.All)
}

func (df *HtmlContent) RenderPdfWithPlugins(pluginList []plugins.Plugin) (*bytes.Buffer, error) {
	renderedHtml, err := df.RenderWithPlugins(pluginList)
	if err != nil {
		return nil, err
	}

	// TODO:: Implement
	return renderedHtml, nil
}
