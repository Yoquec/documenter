package documenter

type processor[T any] interface {
    Apply(df *T)
}

// Holds and runs buffer processors
type Pipeline[T any] struct {
    processors []processor[T]
}

func NewPipeline[T any](processors ...processor[T]) *Pipeline[T] {
    return &Pipeline[T]{
        processors: processors,
    }
}

func (p *Pipeline[T]) Execute(data *T) {
    for _, processor := range p.processors {
        processor.Apply(data)
    }
}
