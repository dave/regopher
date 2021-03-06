package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dave/dst/decorator"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func TestExtractParameterObject(t *testing.T) {
	testCases := []struct {
		pos             string
		additionalFiles []string
	}{
		{
			pos: "params_to_struct_basic.go:#42",
		},
		{
			pos: "params_to_struct_used.go:#35",
		},
		{
			pos: "params_to_struct_referenced.go:#35",
		},
		{
			pos: "params_to_struct_method.go:#70",
		},
		{
			pos: "params_to_struct_collision_by_type.go:#73",
		},
		{
			pos: "params_to_struct_collision_same_type.go:#73",
		},
		{
			pos:             "params_to_struct_referenced_other_file.go:#35",
			additionalFiles: []string{"params_to_struct_referenced_additional_file.go"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.pos, func(t *testing.T) {
			pos, err := parseInputPositionString("testdata/before/" + testCase.pos)
			if err != nil {
				t.Fatal(err)
			}
			file := filepath.Clean(pos.file)
			beforeFiles := []string{file}

			for _, of := range testCase.additionalFiles {
				beforeFiles = append(beforeFiles, filepath.Clean("testdata/before/"+of))
			}

			d, files, err := loadNamedFiles(pos, beforeFiles)
			if err != nil {
				t.Fatal(err)
			}
			funcDecl, err := getFuncAt(d, files[file], pos.pos)
			if err != nil {
				t.Fatal(err)
			}
			_, err = regopherParamsToStruct(pos, files, funcDecl)
			if err != nil {
				t.Fatal(err)
			}
			w := &bytes.Buffer{}
			if err := decorator.Fprint(w, files[file]); err != nil {
				t.Fatal(err)
			}
			actual := string(w.Bytes())
			expected, err := ioutil.ReadFile(strings.Replace(pos.file, "testdata/before", "testdata/after", 1))
			if err != nil {
				t.Fatal(err)
			}
			if string(expected) != actual {
				dmp := diffmatchpatch.New()

				diffs := dmp.DiffMain(string(expected), actual, false)

				t.Error(dmp.DiffPrettyText(diffs))
			}
		})
	}
}
