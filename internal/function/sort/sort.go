package sort

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

const funcPrefixSizeBytes = 5

type dependencyTree struct {
	root         *ast.FuncDecl
	dependencies []*dependencyTree
	writtenUnder *dependencyTree
}

func (t *dependencyTree) Name() string {
	return t.root.Name.Name
}

func (t *dependencyTree) StartPos() token.Pos {
	return token.Pos(int(t.root.Name.Pos()) - (len(t.root.Name.Name) + funcPrefixSizeBytes))
}

func (t *dependencyTree) EndPos() token.Pos {
	return t.root.End()
}

func (t *dependencyTree) IsEmpty() bool {
	return t.dependencies == nil || len(t.dependencies) == 0
}

type trees map[string]*dependencyTree

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

	allTrees := make(trees, len(f.Scope.Objects))
	roots := make(trees, len(f.Scope.Objects))

	for _, object := range f.Scope.Objects {
		declaration := object.Decl
		funcDeclaration, ok := declaration.(*ast.FuncDecl)
		if !ok {
			log.Printf("declaration is not a function: %+v", object)
			continue
		}

		tree := functionDependencyTree(funcDeclaration, roots, allTrees)
		if !tree.IsEmpty() {
			roots[tree.Name()] = tree
		}
	}

	if len(roots) > 0 {
		return reorderFunctions(filePath, fileSet, roots)
	}

	return nil
}

func functionDependencyTree(declaration *ast.FuncDecl, roots, allTrees trees) *dependencyTree {
	tree := &dependencyTree{
		root:         declaration,
		dependencies: make([]*dependencyTree, 0),
	}

	allTrees[declaration.Name.Name] = tree

	for _, statement := range declaration.Body.List {
		exprStatement, ok := statement.(*ast.ExprStmt)
		if !ok {
			log.Printf("statement is not an expression: %+v", statement)
			continue
		}

		callExpr, ok := exprStatement.X.(*ast.CallExpr)
		if !ok {
			log.Printf("expression is not a call expression: %+v", exprStatement)
			continue
		}

		identifier, ok := callExpr.Fun.(*ast.Ident)
		if !ok {
			log.Printf("function is not an identifier: %+v", callExpr.Fun)
			continue
		}

		if identifier.Obj.Kind == ast.Fun && identifier.Obj.Pos() < declaration.Pos() {
			funcDeclaration, ok := identifier.Obj.Decl.(*ast.FuncDecl)
			if ok {
				subTree := dependencySubTree(funcDeclaration, roots, allTrees)
				tree.dependencies = append(tree.dependencies, subTree)
			}
		}
	}

	return tree
}

func dependencySubTree(declaration *ast.FuncDecl, roots, allTrees trees) (tree *dependencyTree) {
	tree, exists := allTrees[declaration.Name.Name]
	if !exists {
		tree = functionDependencyTree(declaration, roots, allTrees)
	}

	// delete the fetched \ created tree since it is no longer a root
	delete(roots, tree.Name())

	return tree
}

func reorderFunctions(filePath string, fileSet *token.FileSet, roots trees) error {
	return errors.New("not implemented")
}
