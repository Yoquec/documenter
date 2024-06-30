package documenter

type Renderer interface {
	RenderDoc(title, path string) (string, error)
}
