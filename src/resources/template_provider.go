package resources

import (
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/yoquec/documenter/src/document"
)

const (
	UNIX    = iota
	WINDOWS = iota
)

var (
	ErrUnknownOs          = errors.New("Specified os is not either UNIX nor WINDOWS")
	ErrTemplateNotFound   = errors.New("Template could not be found")
	ErrStylesheetNotFound = errors.New("Could not find stylesheet")
)

type TemplateProvider struct {
	os uint
}

func NewTemplateProvider(os uint) *TemplateProvider {
	return &TemplateProvider{
		os: os,
	}
}

func GetCurrentOs() uint {
	var currentOs uint
	if strings.ToLower(runtime.GOOS) == "windows" {
		currentOs = WINDOWS
	} else {
		currentOs = UNIX
	}
	return currentOs
}

// Get templates by their name. Filename and template name must coincide for templates
// to be considered valid.
func (tp *TemplateProvider) GetTemplateByName(templateName string) (*document.TemplateInfo, error) {
	resourceDirs, err := getResourceDirectories(tp.os)
	if err != nil {
		return new(
				document.TemplateInfo,
			), fmt.Errorf(
				"Could not retrieve template directories: %w",
				err,
			)
	}

	for i := 0; i < len(resourceDirs); i++ {
		templatePath := path.Join(resourceDirs[i], "templates", templateName+".html")
		_, err := os.Stat(templatePath)

		// TEST: Write a test for this
		if err == nil {
			return &document.TemplateInfo{Name: templateName, Path: templatePath}, nil
		}
	}

	return new(
			document.TemplateInfo,
		), fmt.Errorf(
			"Could not find template named '%s'. Tried directories %v under the 'templates' subdir: %w",
			templateName,
			resourceDirs,
			ErrTemplateNotFound,
		)
}

func (tp *TemplateProvider) GetStyleSheetPath() (string, error) {
	resourceDirs, err := getResourceDirectories(tp.os)
	if err != nil {
		return "", fmt.Errorf(
			"Could not retrieve template directories: %w",
			err,
		)
	}

	for i := 0; i < len(resourceDirs); i++ {
		cssPath := path.Join(resourceDirs[i], "css", "style.css")

		_, err := os.Stat(cssPath)

		// TEST: Write a test for this
		if err == nil {
			return cssPath, nil
		}
	}
	return "", ErrStylesheetNotFound
}

func getResourceDirectories(osString uint) ([]string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return *new([]string), fmt.Errorf("Could not retrieve the current user's directory")
	}

	switch osString {
	case UNIX:
		return []string{
			"/usr/share/go-documenter",
			path.Join(homedir, ".local/share/go-documenter"),
		}, nil
	case WINDOWS:
		// TODO: implement for windows
		return *new([]string), errors.New("Not implemented")
	default:
		return *new([]string), ErrUnknownOs
	}
}
