package pandoc

import (
	"os/exec"
)

const (
	pandocCommand      = "pandoc"
	markdownSourceType = "commonmark_x"
)

func MarkdownToHtml(filename string) ([]byte, error) {
	return exec.Command(pandocCommand,
		"-f", markdownSourceType,
		"-t", "html",
		filename,
	).Output()
}
