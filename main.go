package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

var (
	outputFormat, outputPath string
	allowedFormats           = [...]string{"pdf", "html"}
)

const (
	programName = "documenter"
)

// Function that sets the help message
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

func getFilename() string {
	filename := flag.Args()[0]
	if !strings.HasSuffix(filename, ".md") {
		panic("Input filename does not appear to be a markdown file")
	}
	return filename
}

func main() {
	// NOTE: flag.Parse() cannot be inside the init function, otherwise
	// running tests on package `main` will fail
	// See https://stackoverflow.com/questions/60235896/flag-provided-but-not-defined-test-v
	flag.Parse()

	if flag.NArg() == 0 {
		failWithHelpMessage("Input file was not specified")
	}

	filename := getFilename()

	file, err := newDocumenterFileFromPath(filename)
	if err != nil {
		panic(fmt.Sprintf("Could not read filepath %q\n", filename))
	}

	preprocessorPipeline := newPipeline(
		newFrontMatterPreprocessor(),
	)
    preprocessorPipeline.Execute(file)

	// fakeDocumenterFile := newDocumenterFileFromString("hola\nque\ntal")
	// for line := range fakeDocumenterFile.ReadLines() {
	// 	fmt.Println("line: ", string(line))
	// }
}
