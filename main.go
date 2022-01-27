package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
)

func main() {
	var help bool
	var path string
	var names []string

	flag.BoolVar(&help, "help", false, "Print usage.")
	flag.Func("names", "Names of top-level constants to include.", func(s string) error {
		if s == "" {
			names = []string{}
			return nil
		}
		names = strings.Split(s, ",")
		return nil
	})
	flag.StringVar(&path, "path", "", "Path of Go source file.")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if path == "" {
		return
	}

	fset := token.NewFileSet()
	fileAST, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		log.Fatalf("failed to parse %q: %s", path, err)
	}

	printForMatchingNames := DoForMatchingName(names, func(name, value string) error {
		_, err := fmt.Printf("%s=%s\n", name, value)
		return err
	})

	err = ForAllTopLevelConst(fileAST, printForMatchingNames)
	if err != nil {
		log.Fatalf("failed to process const declarations: %s", err)
	}
}

func ForAllTopLevelConst(fileAST *ast.File, f func(name, value string) error) error {
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
						name := spec.Names[i]
						value := spec.Values[i].(*ast.BasicLit).Value
						if err := f(name.Name, value); err != nil {
							log.Printf("error processing const %q: %s", name, err)
						}
					}
				}
			}
		}
	}
	return nil
}

func DoForMatchingName(names []string, f func(name, value string) error) func(name, value string) error {
	if len(names) == 0 {
		return f
	}

	namesMap := make(map[string]struct{}, len(names))
	for _, n := range names {
		namesMap[n] = struct{}{}
	}

	return func(name, value string) error {
		if _, ok := namesMap[name]; !ok {
			return nil
		}
		return f(name, value)
	}
}
