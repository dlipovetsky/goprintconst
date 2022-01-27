package main

import (
	"go/parser"
	"go/token"
	"reflect"
	"testing"
)

func TestTopLevelConsts(t *testing.T) {
	tcs := []struct {
		name string
		src  string
		want []ConstDecl
	}{
		{
			name: "list top-level constants",
			src: `package example
const (
	Foo = "foo"
	One = 1
)
`,
			want: []ConstDecl{
				{
					name:  "Foo",
					value: `"foo"`,
				},
				{
					name:  "One",
					value: "1",
				},
			},
		},
		{
			name: "ignore other constants",
			src: `package example

func Foo() {
	const Bar = "bar"
}
`,
			want: []ConstDecl{},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			fset := token.NewFileSet()
			fileAST, err := parser.ParseFile(fset, "src.go", tc.src, parser.AllErrors)
			if err != nil {
				t.Fatalf("failed to parse %q: %s", "src.go", err)
			}

			cds, err := TopLevelConsts(fileAST)
			if err != nil {
				t.Errorf("failed to parse source: %s", err)
			}
			if !reflect.DeepEqual(tc.want, cds) {
				t.Errorf("got %s, want %s", cds, tc.want)
			}
		})
	}
}
