package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/yoquec/documenter/src/documenter"
	"github.com/yoquec/documenter/src/pdf"
	"github.com/yoquec/documenter/src/plugins"
	"github.com/yoquec/documenter/src/processor"
	"github.com/yoquec/documenter/src/resources"
	"gitlab.com/golang-commonmark/markdown"
)

var (
	mode, outputPath string
	allowedFormats   = [...]string{"pdf", "html"}
)

const programName = "go-documenter"

// Function that sets the message when -h is called
func setUsage() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v [OPTIONS] FILE\n", programName)
		fmt.Fprintf(os.Stderr, "Example: %v -t pdf -o test.pdf test.md\n\n", programName)
		fmt.Fprintln(os.Stderr, "Flags:")
		flag.PrintDefaults()
	}
}

func init() {
	flag.StringVar(
		&mode,
		"t",
		"pdf",
		fmt.Sprintf("Output format. Possible ouput formats: %v", allowedFormats),
	)
	flag.StringVar(
		&outputPath,
		"o",
		"",
		"Name of the output file. It defaults to printing to stdout.",
	)
	setUsage()
	flag.Parse()

	// check for the input file
	if flag.NArg() == 0 {
		log.Print("Input file was not speciefied\n\n")
		flag.Usage()
		os.Exit(2)
	}
}

func main() {
	filename := flag.Args()[0]
	name := strings.TrimSuffix(filename, ".md")

	engine := markdown.New(
		markdown.Tables(true),  // Render GFM tables
		markdown.Linkify(true), // Generate links for urls automatically
		markdown.HTML(true),    // Ignore html inside the markdown
        markdown.Typographer(true),
	)
	processor := processor.FromEngine(engine)
	provider := resources.NewTemplateProvider(resources.GetCurrentOs())
	renderer := documenter.New(processor, plugins.DefaultPlugins, provider)

	output, err := renderer.RenderDoc(name, filename)
	if err != nil {
		log.Fatal(err)
	}

	if mode == "pdf" {
		output, err = pdf.GenerateFromHtml(output)
		if err != nil {
			log.Fatal(err)
		}
	}

	if outputPath == "" {
		fmt.Println(string(output))
	} else {
		err = os.WriteFile(outputPath, output, os.FileMode(0644))
		if err != nil {
			log.Fatal(err)
		}
	}
}
