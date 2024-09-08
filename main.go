package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/yoquec/documenter/documenter"
	"github.com/yoquec/documenter/documenter/content"
	"github.com/yoquec/documenter/documenter/preprocessors"
)

var (
	outputFormat, outputPath string
	allowedFormats           = []string{"pdf", "html"}
)

const (
	programName = "documenter"
)

// Sets the help message
func setUsageMessage() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v [OPTIONS] FILE\n", programName)
		fmt.Fprintf(os.Stderr, "Example: %v -t pdf -o test.pdf test.md\n\n", programName)
		fmt.Fprintln(os.Stderr, "Flags:")
		flag.PrintDefaults()
	}
}

func failWithHelpMessage(errorMessage string) {
	slog.Error(errorMessage + "\n\n")
	flag.Usage()
	os.Exit(1)
}

func getFilename() string {
	filename := flag.Args()[0]
	if !strings.HasSuffix(filename, ".md") {
		panic("Input filename does not appear to be a markdown file")
	}
	return filename
}

func parseArguments() {
	flag.Parse()

	if flag.NArg() == 0 {
		failWithHelpMessage("Input file was not specified")
	}

	validateArguments()
}

func init() {
	flag.StringVar(
		&outputFormat,
		"t",
		"html",
		fmt.Sprintf("Output format. Possible ouput formats: %v", allowedFormats),
	)
	flag.StringVar(
		&outputPath,
		"o",
		"",
		"Name of the output file. It defaults to printing to stdout.",
	)
	setUsageMessage()
}

func main() {
	parseArguments()
	filename := getFilename()

	markdownFile, err := content.FromPath(filename)
	if err != nil {
		panic(err)
	}

	pipeline := documenter.NewPipeline(
		preprocessors.NewYamlFrontmatterProcessor(),
	)

	markdownFile.Apply(pipeline)

	htmlFile, err := markdownFile.ToHtml()
	if err != nil {
		panic(err)
	}

	var output *bytes.Buffer
	if outputFormat == "html" {
		output, err = htmlFile.Render()
	} else {
		panic("PDF conversion is not yet implemented")
		// output, err = htmlFile.RenderPdf()
	}
	if err != nil {
		panic(err)
	}

	var file io.Writer
	if outputPath != "" {
		file, err = os.Create(outputPath)
		if err != nil {
			panic(err)
		}
	} else {
		file = os.Stdout
	}

	output.WriteTo(file)
}
