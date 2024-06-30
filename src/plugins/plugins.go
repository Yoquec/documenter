package plugins

import (
	"html/template"
	"strings"
)

type Plugin struct {
	// HTML import code for the plugin
	PluginHead template.HTML
	// Code to execute when the window is loaded for this plugin
	WindowLoadCode template.JS
}

var twemoji Plugin = Plugin{
	template.HTML(
		`<script src="https://unpkg.com/twemoji@latest/dist/twemoji.min.js" crossorigin="anonymous"></script>`,
	),
	template.JS("twemoji.parse(document.body);"),
}

var katex Plugin = Plugin{
	template.HTML(strings.Join([]string{
		`<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/katex.min.css" integrity="sha384-n8MVd4RsNIU0tAv4ct0nTaAbDJwPJzDEaqSD1odI+WdtXRGWt2kTvGFasHpSy3SV" crossorigin="anonymous">`,
		`<script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/katex.min.js" integrity="sha384-XjKyOOlGwcjNTAIQHIpgOno0Hl1YQqzUOEleOLALmuqehneUG+vnGctmUb0ZY0l8" crossorigin="anonymous"></script>`,
		`<script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/contrib/auto-render.min.js" integrity="sha384-+VBxd3r6XgURycqtZ117nYw44OOcIax56Z4dCRWbxyPt0Koah1uHoK0o4+/RRE05" crossorigin="anonymous"></script>`,
	}, "\n")),
	template.JS("renderMathInElement(document.body);"),
}

var highlightjs Plugin = Plugin{
	template.HTML(
		strings.Join([]string{
			`<link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/styles/default.min.css">`,
			`<script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/highlight.min.js"></script>`,
		}, "\n")),
	template.JS("hljs.highlightAll();"),
}

var mermaid Plugin = Plugin{
	template.HTML(
		`<script src="https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js"></script>`,
	),
	template.JS("mermaid.initialize();"),
}

var DefaultPlugins = []Plugin{
	twemoji,
	katex,
	highlightjs,
	mermaid,
}
