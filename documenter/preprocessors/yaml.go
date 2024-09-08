package preprocessors

import (
	"bytes"
	"log/slog"

	"github.com/yoquec/documenter/documenter/content"
	"gopkg.in/yaml.v3"
)

type yamlFrontmatter struct {
	Id      string   `yaml:"id"`
	Aliases []string `yaml:"aliases,omitempty"`
	Tags    []string `yaml:"tags,omitempty"`
}

type YamlFrontmatterProcessor struct{}

func NewYamlFrontmatterProcessor() *YamlFrontmatterProcessor {
	return &YamlFrontmatterProcessor{}
}

func isYamlFrontmatterDelimiter(line []byte) bool {
	if len(line) == 0 {
		return false
	}

	// The first line must only contain '-' code points.
	// Otherwise, it will be considered as no YAML frontmatter is present
	return !bytes.ContainsFunc(line, func(r rune) bool {
		return r != '-'
	})
}

func (yp *YamlFrontmatterProcessor) shouldExecute(md *content.MarkdownContent) bool {
    bufferBytes := md.Content.Buffer.Bytes()
	firstLineEnd := bytes.Index(bufferBytes, []byte("\n"))

	// Case that the file has a single line
	if firstLineEnd == -1 {
		return false
	}

	firstLine := bufferBytes[:firstLineEnd]
	return isYamlFrontmatterDelimiter(firstLine)
}

func (yp *YamlFrontmatterProcessor) getFrontmatter(df *content.MarkdownContent) ([]byte, int) {
	frontmatterLines := make([][]byte, 0)
	firstLine := true

	for line := range df.ReadLines() {
		if firstLine {
			firstLine = false
			continue
		}
		if isYamlFrontmatterDelimiter(line) {
			break
		}
		frontmatterLines = append(frontmatterLines, line)
	}

	if len(frontmatterLines) == 0 {
		return nil, -1
	}

	return bytes.Join(frontmatterLines, []byte("\n")), len(frontmatterLines) + 2
}

func (yp *YamlFrontmatterProcessor) Apply(md *content.MarkdownContent) {
	if !yp.shouldExecute(md) {
		slog.Debug("Didn't find YAML frontmatter in file")
		return
	}

	frontmatter, lines := yp.getFrontmatter(md)
	if frontmatter == nil {
		slog.Debug("Could not get contents of the YAML frontmatter. Is the frontmatter empty?")
		return
	}

	parsedFrontmatter := &yamlFrontmatter{}
	if err := yaml.Unmarshal(frontmatter, parsedFrontmatter); err != nil {
		slog.Warn("Could not parse contents of the YAML frontmatter")
	}

	if len(parsedFrontmatter.Aliases) > 0 {
		md.Metadata.Title = parsedFrontmatter.Aliases[0]
	}

	if len(parsedFrontmatter.Tags) > 0 {
		md.Metadata.Tags = parsedFrontmatter.Tags
	}

	// Cut the frontmatter from the file contents so it does not get converted into HTML
	for i := 0; i < lines; i++ {
        md.Content.Buffer.ReadBytes('\n')
	}
}
