package sort

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"log"

	"github.com/Hofled/whimify/internal/unique"
)

// ByDependencyOrder sorts the functions in a file by dependency order.
func FileByDependencyOrder(filePath string) error {
	if filePath == "" {
		return errors.New("file path cannot be empty")
	}

	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, filePath, nil, parser.AllErrors)
	if err != nil {
		return err
	}

	forest := newForest(len(f.Scope.Objects))

	for _, object := range f.Scope.Objects {
		declaration := object.Decl
		funcDeclaration, ok := declaration.(*ast.FuncDecl)
		if !ok {
			log.Printf("declaration is not a function: %+v", object)
			continue
		}

		forest.AddTreeFromFunction(funcDeclaration)
	}

	if len(forest.roots) > 0 {
		return sortRoots(filePath, fileSet, forest.roots)
	}

	return nil
}

func sortRoots(filePath string, fileSet *token.FileSet, unsortedRoots unique.Map[string, *dependencyTree]) error {
	return errors.New("not implemented")
}
