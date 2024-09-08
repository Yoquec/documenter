package main

import (
	"fmt"
	"slices"
)

func validateArguments() {
	validateOutputFormat()
}

func validateOutputFormat() {
	if !slices.Contains(allowedFormats, outputFormat) {
		panic(fmt.Errorf("%q is not a valid output format", outputFormat))
	}
}
