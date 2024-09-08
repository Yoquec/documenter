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

var Twemoji Plugin = Plugin{
	template.HTML(
		`<script src="https://unpkg.com/twemoji@latest/dist/twemoji.min.js" crossorigin="anonymous"></script>`,
	),
	template.JS("twemoji.parse(document.body);"),
}

var MathJax Plugin = Plugin{
	template.HTML(
		strings.Join([]string{
			`<script>MathJax = { tex: { inlineMath: [['$', '$'], ['\\(', '\\)']] } }; </script>`,
			`<script id="MathJax-script" async src="https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-chtml.js"> </script>`,
		}, "\n"),
	),
	"",
}

var HighlightJs Plugin = Plugin{
	template.HTML(
		strings.Join([]string{
			`<link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/styles/default.min.css">`,
			`<script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/highlight.min.js"></script>`,
		}, "\n")),
	template.JS("hljs.highlightAll();"),
}

var Mermaid Plugin = Plugin{
	template.HTML(
		`<script src="https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js"></script>`,
	),
	template.JS("mermaid.initialize();"),
}

var All []Plugin = []Plugin {
    Twemoji,
    MathJax,
    HighlightJs,
    Mermaid,
}
