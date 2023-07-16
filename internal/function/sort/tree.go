package sort

import (
	"go/ast"
	"go/token"
	"sort"
)

type dependencyTree struct {
	root     *ast.FuncDecl
	children *treeSet
	parents  *treeSet
}

func newDependencyTree(root *ast.FuncDecl) *dependencyTree {
	return &dependencyTree{
		root:     root,
		children: newTreeSet(lexicographic),
		parents:  newTreeSet(lexicographic),
	}
}

func (t *dependencyTree) Name() string {
	return t.root.Name.Name
}

const funcPrefixSizeBytes = 5

func (t *dependencyTree) StartPos() token.Pos {
	return token.Pos(int(t.root.Name.Pos()) - (len(t.root.Name.Name) + funcPrefixSizeBytes))
}

func (t *dependencyTree) EndPos() token.Pos {
	return t.root.End()
}

func (t *dependencyTree) IsLeaf() bool {
	return t.children == nil || t.children.Len() == 0
}

type by func(t1, t2 *dependencyTree) bool

// treeSet maintains a sorted list of unique trees.
type treeSet struct {
	trees  []*dependencyTree
	keys   map[string]bool
	byFunc by
}

func newTreeSet(byFunc by) *treeSet {
	return &treeSet{
		trees:  make([]*dependencyTree, 0),
		keys:   make(map[string]bool),
		byFunc: byFunc,
	}
}

func (s *treeSet) UpdateBy(byFunc by) {
	s.byFunc = byFunc
	sort.Sort(s)
}

func (s *treeSet) Len() int {
	return len(s.trees)
}

func (s *treeSet) Swap(i, j int) {
	s.trees[i], s.trees[j] = s.trees[j], s.trees[i]
}

func (s *treeSet) Less(i, j int) bool {
	return s.byFunc(s.trees[i], s.trees[j])
}

func (s *treeSet) Exists(tree *dependencyTree) bool {
	return s.keys[tree.Name()]
}

func (s *treeSet) Find(tree *dependencyTree) int {
	return sort.Search(s.Len(), func(i int) bool {
		return s.byFunc(tree, s.trees[i])
	})
}

func (s *treeSet) Add(tree *dependencyTree) {
	if s.Exists(tree) {
		return
	}

	i := s.Find(tree)
	s.keys[tree.Name()] = true
	s.trees = append(s.trees, nil)
	copy(s.trees[i+1:], s.trees[i:])
	s.trees[i] = tree
}

func (s *treeSet) Remove(tree *dependencyTree) {
	if !s.Exists(tree) {
		return
	}

	i := s.Find(tree)
	delete(s.keys, tree.Name())
	s.trees = append(s.trees[:i], s.trees[i+1:]...)
}
