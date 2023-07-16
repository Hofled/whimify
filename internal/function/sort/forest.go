package sort

import (
	"go/ast"
	"log"

	"github.com/Hofled/whimify/internal/unique"
)

type forest struct {
	all   unique.Map[string, *dependencyTree]
	roots unique.Map[string, *dependencyTree]
}

func newForest(size int) *forest {
	return &forest{
		all:   make(unique.Map[string, *dependencyTree], size),
		roots: make(unique.Map[string, *dependencyTree], size),
	}
}

func (f *forest) Exists(name string) (*dependencyTree, bool) {
	return f.all.Exists(name)
}

func (f *forest) AddTreeFromFunction(declaration *ast.FuncDecl) *dependencyTree {
	if tree, exists := f.Exists(declaration.Name.Name); exists {
		return tree
	}

	tree := newDependencyTree(declaration)

	f.all.Add(tree.Name(), tree)
	f.roots.Add(tree.Name(), tree)

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

		if identifier.Obj.Kind == ast.Fun && identifier.Obj.Decl != declaration && identifier.Obj.Pos() < declaration.Pos() {
			funcDeclaration, ok := identifier.Obj.Decl.(*ast.FuncDecl)
			if ok {
				subTree := f.AddTreeFromFunction(funcDeclaration)
				// remove subtree from being a root
				delete(f.roots, subTree.Name())
				tree.children.Add(subTree)
				subTree.parents.Add(tree)
			}
		}
	}

	if tree.IsLeaf() {
		delete(f.roots, tree.Name())
	}

	return tree
}
