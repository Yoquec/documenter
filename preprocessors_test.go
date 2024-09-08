package main

import (
	"strings"
	"testing"
)

func TestIsYamlFrontMatterDelimeter(t *testing.T) {
	tests := []struct {
		line        string
		isDelimeter bool
	}{
		{"---", true},
		{"------", true},
		{"------foo", false},
		{"key: value", false},
		{"\n", false},
		{"", false},
	}

	for _, test := range tests {
		actual := isYamlFrontmatterDelimiter([]byte(test.line))
		if actual != test.isDelimeter {
			t.Errorf(
				"Line %q has not been correctly classified as a delimiter. Expected %v, got %v",
				test.line,
				test.isDelimeter,
				actual,
			)
		}
	}
}

func TestFrontmatterProcessorSetsTitleInMetadata(t *testing.T) {
	processor := newFrontMatterPreprocessor()
	tests := []struct {
		contents      string
		expectedTitle string
	}{
		{strings.Join([]string{
			`---`,
			`aliases:`,
			`   - foo`,
			`   - bar`,
			`---`,
		}, "\n"), "foo"},

		// empty alias section, should leave empty title
		{strings.Join([]string{
			`---`,
			`id: bar`,
			`---`,
		}, "\n"), ""},

		// No frontmatter should not modify empty title
		{"Foobarbaz", ""},
	}

	for _, test := range tests {
		file := newDocumenterFileFromString(test.contents)
		processor.Apply(file)

		if file.metadata.title != test.expectedTitle {
			t.Errorf(
				"Failed to set title correctly in metadata. Expected %q and got %q",
				test.expectedTitle,
				file.metadata.title,
			)
		}
	}
}

func TestFrontmatterProcessorTrimsFrontmatterFromContents(t *testing.T) {
	processor := newFrontMatterPreprocessor()

	tests := []struct {
		contents string
		expected string
	}{
		{strings.Join([]string{
			`---`,
			`id: foobar`,
			`aliases:`,
			`   - foo`,
			`   - bar`,
			`tags:`,
			`   - baz`,
			`---`,
			``,
			`foobarbaz`,
		}, "\n"), "\nfoobarbaz"},
	}

	for _, test := range tests {
		file := newDocumenterFileFromString(test.contents)
		processor.Apply(file)
		actual := string(file.contents)

		if actual != test.expected {
			t.Errorf(
				"Processed contents didn't match expected contents.\nExpected:\n```\n%s\n```\n\nGot:\n```\n%s\n```",
				test.expected,
				actual,
			)
		}
	}
}
