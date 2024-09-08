package main

import (
	"bytes"
	"log/slog"

	"gopkg.in/yaml.v3"
)

type obsidianNvimYamlFrontmatter struct {
	Id      string   `yaml:"id,omitempty"`
	Aliases []string `yaml:"aliases,omitempty"`
	Tags    []string `yaml:"tags,omitempty"`
}

type ObsidianNvimYamlFrontmatterPreprocessor struct{}

func newFrontMatterPreprocessor() *ObsidianNvimYamlFrontmatterPreprocessor {
	return &ObsidianNvimYamlFrontmatterPreprocessor{}
}

func isYamlFrontmatterDelimiter(line []byte) bool {
	if len(line) == 0 {
		return false
	}

	// The first line must only contain '-' codepoints.
	// Otherwise, it will be considered as no YAML frontmatter
	// is present
	return !bytes.ContainsFunc(line, func(r rune) bool {
		return r != '-'
	})
}

func (fpp *ObsidianNvimYamlFrontmatterPreprocessor) shouldExecute(df *DocumenterFile) bool {
	firstLineEnd := bytes.Index(df.contents, []byte("\n"))

    // Case that the file has a single line
    if firstLineEnd == -1 {
        return false
    }

	firstLine := []byte(df.contents)[:firstLineEnd]
	return isYamlFrontmatterDelimiter(firstLine) 
}

func (fpp *ObsidianNvimYamlFrontmatterPreprocessor) getFrontmatter(
	df *DocumenterFile,
) ([]byte, int) {
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

func (fpp *ObsidianNvimYamlFrontmatterPreprocessor) Apply(df *DocumenterFile) {
	if !fpp.shouldExecute(df) {
		slog.Debug("Didn't find YAML frontmatter in file")
		return
	}

	frontmatter, lines := fpp.getFrontmatter(df)
	if frontmatter == nil {
		slog.Debug("Could not get contents of the YAML frontmatter. Is the frontmatter empty?")
		return
	}

	parsedFrontmatter := &obsidianNvimYamlFrontmatter{}
	if err := yaml.Unmarshal(frontmatter, parsedFrontmatter); err != nil {
		slog.Warn("Could not parse contents of the YAML frontmatter")
	}

	if len(parsedFrontmatter.Aliases) > 0 {
		df.metadata.title = parsedFrontmatter.Aliases[0]
	}

	if len(parsedFrontmatter.Tags) > 0 {
		df.metadata.tags = parsedFrontmatter.Tags
	}

	// cut the frontmatter from the file contents so it does not get converted into HTML
	for i := 0; i < lines; i++ {
		_, df.contents, _ = bytes.Cut(df.contents, []byte("\n"))
	}
}
