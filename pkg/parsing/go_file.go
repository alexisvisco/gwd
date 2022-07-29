package parsing

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"

	"github.com/alexisvisco/gwd/pkg/diff/packages"
)

// GetImports retrieve imports of a go file
func GetImports(filepath string) []packages.ImportPath {
	set := token.NewFileSet()
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil
	}
	astFile, err := parser.ParseFile(set, filepath, bytes, 0)
	if err != nil {
		return nil
	}

	var paths []packages.ImportPath
	importList := imports(set, astFile)
	for _, list := range importList {
		for _, i := range list {
			paths = append(paths, packages.ImportPath(i.Path.Value[1:len(i.Path.Value)-1]))
		}
	}

	return paths
}

func imports(tokenFileSet *token.FileSet, f *ast.File) [][]*ast.ImportSpec {
	var groups [][]*ast.ImportSpec

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT {
			break
		}

		var group []*ast.ImportSpec

		var lastLine int
		for _, spec := range genDecl.Specs {
			importSpec := spec.(*ast.ImportSpec)
			pos := importSpec.Path.ValuePos
			line := tokenFileSet.Position(pos).Line
			if lastLine > 0 && pos > 0 && line-lastLine > 1 {
				groups = append(groups, group)
				group = []*ast.ImportSpec{}
			}
			group = append(group, importSpec)
			lastLine = line
		}
		groups = append(groups, group)
	}

	return groups
}
