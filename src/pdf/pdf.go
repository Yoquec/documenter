package pdf

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

var chromeCommands = []string{"chromium", "google-chrome", "chrome"}

// GEnerates a PDF file from a given HTML content using the chromium browser pdf engine.
func GenerateFromHtml(contents []byte) ([]byte, error) {
	cmd, err := getChromeCmd()
	if err != nil {
		return nil, err
	}

	tempFile, err := os.CreateTemp("/tmp", "*.html")
	if err != nil {
        log.Println("Could not create an html temporary file.")
		return nil, err
	}

	tempPdfFile, err := os.CreateTemp("/tmp", "*.pdf")
	if err != nil {
        log.Println("Could not create a pdf temporary file.")
		return nil, err
	}

	err = os.WriteFile(tempFile.Name(), contents, os.FileMode(0644))
	if err != nil {
        log.Println("Could not write the html content to the temporary file.")
		return nil, err
	}

	err = exec.Command(
		cmd,
		"--headless",
		"--disable-gpu",
		"--no-pdf-header-footer",
		fmt.Sprintf("--print-to-pdf=%v", tempPdfFile.Name()),
		tempFile.Name(),
	).Run()
	if err != nil {
        log.Println("Could not execute the command to generate the pdf.")
		return nil, err
	}

	output, err := os.ReadFile(tempPdfFile.Name())
	if err != nil {
        log.Println("Could not read the output of the pdf generation.")
		return nil, err
	}

	return output, nil
}

// Returns chome/chromium depending on which one is installed in the system.
func getChromeCmd() (string, error) {
	for _, cmd := range chromeCommands {
		if isCommandAvailable(cmd) {
			return cmd, nil
		}
	}

	return "", fmt.Errorf("Chrome/chromium is not installed or executable.")
}

func isCommandAvailable(cmd string) bool {
	if err := exec.Command(cmd, "--version").Run(); err != nil {
		return false
	}
	return true
}
