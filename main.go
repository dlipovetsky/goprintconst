package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
)

type ConstDecl struct {
	name  string
	value string
}

type multivalueFlags struct {
	values map[string]struct{}
}

func (i *multivalueFlags) String() string {
	// fmt prints in key-sorted order, see https://tip.golang.org/doc/go1.12#fmt
	return fmt.Sprint(i.values)
}

func (i *multivalueFlags) Set(value string) error {
	i.values[value] = struct{}{}
	return nil
}

func (i *multivalueFlags) Has(value string) bool {
	_, ok := i.values[value]
	return ok
}

func (i *multivalueFlags) Any() bool {
	return len(i.values) == 0
}

func main() {
	var help bool
	var path string
	names := multivalueFlags{values: make(map[string]struct{})}

	flag.BoolVar(&help, "help", false, "Print usage.")
	flag.Var(&names, "name", "Name of top-level constant to include. Can be used more than once.")
	flag.StringVar(&path, "path", "", "Path of Go source file.")
	flag.Parse()

	if help || flag.NFlag() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	fset := token.NewFileSet()
	fileAST, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		log.Fatalf("failed to parse %q: %s", path, err)
	}

	cds, err := TopLevelConsts(fileAST)
	if err != nil {
		log.Fatalf("failed to get constants: %s", err)
	}

	for _, cd := range cds {
		if !(names.Any() || names.Has(cd.name)) {
			continue
		}

		_, err := fmt.Printf("%s=%s\n", cd.name, cd.value)
		if err != nil {
			log.Printf("failed to print constant %q: %s", cd.name, err)
		}
	}
}

func TopLevelConsts(fileAST *ast.File) ([]ConstDecl, error) {
	cds := []ConstDecl{}

	for _, decl := range fileAST.Decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			if decl.Tok != token.CONST {
				continue
			}

			for _, spec := range decl.Specs {
				switch spec := spec.(type) {
				case *ast.ValueSpec:
					// Assume that spec.Names and spec.Values are parallel arrays.
					for i := range spec.Names {
						cds = append(cds, ConstDecl{
							name:  spec.Names[i].String(),
							value: spec.Values[i].(*ast.BasicLit).Value,
						})
					}
				}
			}
		}
	}
	return cds, nil
}
