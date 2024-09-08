package main

import (
	"slices"
	"testing"
)

func TestReadLinesIter(t *testing.T) {
	tests := []struct {
		data     string
		expected [][]byte
	}{
		{"foo\nbar\nbaz", [][]byte{[]byte("foo"), []byte("bar"), []byte("baz")}},
	}

	for _, test := range tests {
		file := newDocumenterFileFromString(test.data)
		result := slices.Collect(file.ReadLines())

		if len(result) != len(test.expected) {
			t.Errorf(
				"Result and expected lengths differ. Expected %v and got %v",
				test.expected,
				result,
			)
		}

		for i := 0; i < len(result); i++ {
			if !slices.Equal(result[i], test.expected[i]) {
				t.Errorf(
					"Expected and result contents differ. Expected %v and got %v",
					test.expected,
					result,
				)
			}
		}
	}
}
