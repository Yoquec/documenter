package documenter

import (
	"html/template"
	"log"
	"os"
	"path"
	"strings"
)

const shareFolder = "/usr/share/go-documenter/"

var plugins = []DocumenterPlugin{
	// twemoji plugin
	{
		PluginHead: template.HTML(
			`<script src="https://unpkg.com/twemoji@latest/dist/twemoji.min.js" crossorigin="anonymous"></script>`,
		),
		WindowLoadCode: template.JS("twemoji.parse(document.body);"),
	},
	// KaTeX plugin
	{
		PluginHead: template.HTML(
			strings.Join([]string{
				`<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/katex.min.css" integrity="sha384-n8MVd4RsNIU0tAv4ct0nTaAbDJwPJzDEaqSD1odI+WdtXRGWt2kTvGFasHpSy3SV" crossorigin="anonymous">`,
				`<script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/katex.min.js" integrity="sha384-XjKyOOlGwcjNTAIQHIpgOno0Hl1YQqzUOEleOLALmuqehneUG+vnGctmUb0ZY0l8" crossorigin="anonymous"></script>`,
				`<script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/contrib/auto-render.min.js" integrity="sha384-+VBxd3r6XgURycqtZ117nYw44OOcIax56Z4dCRWbxyPt0Koah1uHoK0o4+/RRE05" crossorigin="anonymous"></script>`,
			}, "\n"),
		),
		WindowLoadCode: template.JS("renderMathInElement(document.body);"),
	},
	// highlight.js plugin
	{
		PluginHead: template.HTML(
			strings.Join([]string{
				`<link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/styles/default.min.css">`,
				`<script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/highlight.min.js"></script>`,
			}, "\n"),
		),
		WindowLoadCode: template.JS("hljs.highlightAll();"),
	},
	// mermaid.js plugin
	{
		PluginHead: template.HTML(
			`<script src="https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js"></script>`,
		),
		WindowLoadCode: template.JS("mermaid.initialize();"),
	},
}

type DocumenterDoc struct {
	Head    DocumenterHead
	Content template.HTML
}

type DocumenterHead struct {
	Title      string
	StyleSheet template.URL
	Plugins    []DocumenterPlugin
}

type DocumenterPlugin struct {
	// HTML code for the plugin
	PluginHead template.HTML
	// Code to execute when the window is loaded for this plugin
	WindowLoadCode template.JS
}

// Generates the html file with the content of the markdown file
func Generate(title string, content []byte) ([]byte, error) {
	t, err := template.ParseFiles(
		path.Join(shareFolder, "templates", "index.html"),
	)
	if err != nil {
		log.Println("Could not parse the template file.")
		return nil, err
	}

	tempFile, err := os.CreateTemp("/tmp", "*.html")
	if err != nil {
		log.Println("Could not create an html temporary file.")
		return nil, err
	}

	err = t.ExecuteTemplate(tempFile, "index", DocumenterDoc{
		Head: DocumenterHead{
			Title: title,
			StyleSheet: template.URL(
				"file://" + path.Join(shareFolder, "css", "style.css"),
			),
			Plugins: plugins,
		},
		Content: template.HTML(content),
	})
	if err != nil {
		log.Println("Could not execute the template.")
		return nil, err
	}

	output, err := os.ReadFile(tempFile.Name())
	if err != nil {
		log.Println("Could not read output of the template execution.")
		return nil, err
	}

	return output, nil
}
