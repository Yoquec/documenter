package processor

import (
	"bytes"
	"fmt"
	"gitlab.com/golang-commonmark/markdown"
)

// Markdown processor wrapper for the markdown libray
type Processor struct {
    engine *markdown.Markdown
}

func FromEngine(engine *markdown.Markdown) *Processor {
    return &Processor{
        engine: engine,
    }
}

func (p *Processor) Render(contents []byte) ([]byte, error)  {
    buffer := bytes.Buffer{}
    err := p.engine.Render(&buffer, contents)
    if err != nil {
        return nil, fmt.Errorf("Could not render markdown: %w", err)
    }
    return buffer.Bytes(), nil
}
