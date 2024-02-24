package pandoc

import (
	"os/exec"
)

const (
	pandocCommand      = "pandoc"
	markdownSourceType = "markdown_phpextra"
)

func MarkdownToHtml(filename string) ([]byte, error) {
	return exec.Command(pandocCommand,
		"-f", markdownSourceType,
		"-t", "html",
		filename,
	).Output()
}
