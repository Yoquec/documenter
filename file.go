package main

import (
	"bytes"
	"iter"
	"os"
	"strings"
)

type DocumenterFile struct {
	metadata *DocumenterMetadata
	contents []byte
}

type DocumenterMetadata struct {
	title string
	tags  []string
}

func newDocumenterFileFromPath(path string) (*DocumenterFile, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	metadata := newEmptyDocumenterMetadata()
	metadata.title = strings.TrimSuffix(path, ".md")
	return &DocumenterFile{metadata, file}, nil
}

func newDocumenterFileFromString(contents string) *DocumenterFile {
	return &DocumenterFile{newEmptyDocumenterMetadata(), []byte(contents)}
}

func newEmptyDocumenterMetadata() *DocumenterMetadata {
	return &DocumenterMetadata{"", nil}
}

func (df DocumenterFile) ReadLines() iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		data := bytes.Clone(df.contents)
		for len(data) > 0 {
			line, rest, _ := bytes.Cut(data, []byte("\n"))
			if !yield(line) {
				return
			}
			data = rest
		}
	}
}
