package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/yoquec/documenter/src/documenter"
	"github.com/yoquec/documenter/src/pandoc"
	"github.com/yoquec/documenter/src/pdf"
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

	content, err := pandoc.MarkdownToHtml(filename)
	if err != nil {
		log.Println("Could not covert the markdown file to html.")
		log.Fatal(err)
	}

	name := strings.TrimSuffix(filename, ".md")

	output, err := documenter.Generate(name, content)
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
