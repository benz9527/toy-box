package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"
)

func TestAstSourceFile(t *testing.T) {
	filename := filepath.Join(os.TempDir(), "main.go")

	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	ast.Print(fset, node)
}
