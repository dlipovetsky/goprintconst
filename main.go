package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strconv"
)

func main() {
	var help bool
	var filepath string
	var name string
	var raw bool

	flag.BoolVar(&help, "help", false, "Print usage.")
	flag.StringVar(&filepath, "file", "", "Path to Go source file.")
	flag.StringVar(&name, "name", "", "Name of top-level constant.")
	flag.BoolVar(&raw, "raw", true, "Remove quotes from string and character values.")
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

	tok, value, ok := FindTopLevelConstValue(fileAST, name)
	if !ok {
		log.Fatalf("failed to find top-level constant %s", name)
	}

	if tok == token.STRING || tok == token.CHAR {
		if raw {
			value, err = strconv.Unquote(value)
			if err != nil {
				log.Fatalf("failed to unquote value %s: %s", value, err)
			}
		}
	}

	_, err = fmt.Println(value)
	if err != nil {
		log.Printf("failed to print value of constant %q: %s", name, err)
	}
}

func FindTopLevelConstValue(fileAST *ast.File, name string) (token.Token, string, bool) {
	for _, d := range fileAST.Decls {
		if decl, ok := d.(*ast.GenDecl); ok {
			if decl.Tok != token.CONST {
				continue
			}

			constMap := map[string]*ast.BasicLit{}
			for _, s := range decl.Specs {
				if spec, ok := s.(*ast.ValueSpec); ok {
					for _, id := range spec.Names {
						switch val := id.Obj.Decl.(*ast.ValueSpec).Values[0].(type) {
						case *ast.BasicLit:
							if id.Name == name {
								return val.Kind, val.Value, true
							}
							constMap[id.Name] = val
						case *ast.Ident:
							if id.Name == name {
								name = val.Name
								if v, ok := constMap[name]; ok {
									return v.Kind, v.Value, true
								}
							}
						}
					}
				}
			}
		}
	}
	return 0, "", false
}
