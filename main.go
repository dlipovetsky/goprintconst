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

func main() {
	var help bool
	var filepath string
	var name string

	flag.BoolVar(&help, "help", false, "Print usage.")
	flag.StringVar(&filepath, "file", "", "Path to Go source file.")
	flag.StringVar(&name, "name", "", "Name of top-level constant.")
	flag.Parse()

	if help || flag.NFlag() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if filepath == "" {
		log.Fatal("path must not be empty")
	}

	if name == "" {
		log.Fatal("name must not be empty")
	}

	fset := token.NewFileSet()
	fileAST, err := parser.ParseFile(fset, filepath, nil, parser.AllErrors)
	if err != nil {
		log.Fatalf("failed to parse %q: %s", filepath, err)
	}

	value, ok := FindTopLevelConstValue(fileAST, name)
	if !ok {
		log.Fatalf("failed to find top-level constant %s", name)
	}

	_, err = fmt.Println(value)
	if err != nil {
		log.Printf("failed to print constant %q: %s", name, err)
	}
}

func FindTopLevelConstValue(fileAST *ast.File, name string) (string, bool) {
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
						if name == spec.Names[i].String() {
							return spec.Values[i].(*ast.BasicLit).Value, true
						}
					}
				}
			}
		}
	}
	return "", false
}
