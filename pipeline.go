package main

type Processor interface {
    Apply(df *DocumenterFile)
}

// Holds and runs buffer processors
type Pipeline struct {
    processors []Processor
}

func newPipeline(processors ...Processor) *Pipeline {
    return &Pipeline{
        processors: processors,
    }
}

func (p *Pipeline) Execute(data *DocumenterFile) {
    for _, processor := range p.processors {
        processor.Apply(data)
    }
}
