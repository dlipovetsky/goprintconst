package main

import (
	"go/parser"
	"go/token"
	"testing"
)

func TestFindTopLevelConstValue(t *testing.T) {
	tcs := []struct {
		description string
		src         string
		name        string
		wantOk      bool
		wantValue   string
	}{
		{
			description: "get value for top-level constant",
			src: `package example
const (
	Foo = "foo"
	One = 1
)
`,
			name:      "One",
			wantOk:    true,
			wantValue: "1",
		},
		{
			description: "no value for other constants",
			src: `package example

func Foo() {
	const Bar = "bar"
}
`,
			name:      "Bar",
			wantOk:    false,
			wantValue: "",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.description, func(t *testing.T) {
			fset := token.NewFileSet()
			fileAST, err := parser.ParseFile(fset, "src.go", tc.src, parser.AllErrors)
			if err != nil {
				t.Fatalf("failed to parse %q: %s", "src.go", err)
			}

			_, got, ok := FindTopLevelConstValue(fileAST, tc.name)
			if ok != tc.wantOk {
				t.Errorf("ok %t, wantOk %t", ok, tc.wantOk)
			}
			if got != tc.wantValue {
				t.Errorf("got %s, want %s", got, tc.wantValue)
			}
		})
	}
}
