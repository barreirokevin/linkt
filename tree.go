package main

import (
	"errors"
	"iter"
	"reflect"
)

// Consists of a root Node and a size that indicates the number of
// Nodes that make up the general tree. Any Node in a general tree
// has any number of children, in contrast with a binary tree wherein a node has
// at most two children.
type Tree[T any] struct {
	root *Node[T]
	size int
}

// Returns a pointer to a new general tree with root set to nil and size set to 0.
func NewTree[T any]() *Tree[T] {
	return &Tree[T]{root: nil, size: 0}
}

// Returns the size of the general tree, i.e. the total number of Nodes.
func (t *Tree[T]) Size() int {
	return t.size
}

// Returns a pointer to the root position of the general tree (or nil if empty)
func (t *Tree[T]) Root() *Node[T] {
	return t.root
}

// Creates a root for an empty general tree, storing e as the root's element,
// and returns the Node at the root. If the general tree has a root, then
// e is not added as the root and an error is returned.
func (t *Tree[T]) AddRoot(e T) (*Node[T], error) {
	if t.root != nil {
		return nil, errors.New("tree is not empty")
	}
	t.root = &Node[T]{
		element:  e,
		parent:   nil,
		children: nil,
	}
	t.size = 1
	return t.root, nil
}

// Creates a child for Node n, storing element e, and returns child.
func (t *Tree[T]) AddChild(n *Node[T], e T) *Node[T] {
	child := &Node[T]{element: e, parent: n, children: nil}
	n.children = append(n.children, child)
	t.size += 1
	return child
}

// TODO: func (t *Tree[T]) Remove(n *Node[T]) (*T, error)

// Returns the height of the subtree rooted at a Node in the general tree.
func (t *Tree[T]) Height(n *Node[T]) int {
	h := 0 // base case if n is external
	for _, c := range n.Children() {
		h = max(h, 1+t.Height(c))
	}
	return h
}

// Returns the number of levels separating Node n from the root Node in
// the general tree.
func (t *Tree[T]) Depth(n *Node[T]) int {
	if reflect.DeepEqual(t.root, n) {
		return 0
	} else {
		return 1 + t.Depth(n.parent)
	}
}

// Adds general nodes of the subtree rooted at general node n to the given snapshot using
// the preorder traversal algorithm.
func (t *Tree[T]) preorderSubtree(n *Node[T], snapshot []*Node[T]) {
	snapshot = append(snapshot, n) // preorder algorithm requires adding n before exploring subtrees
	for _, c := range n.Children() {
		t.preorderSubtree(c, snapshot)
	}
}

// Returns an iterator that performs a preorder traversal of the general tree.
func (t *Tree[T]) Preorder() iter.Seq[*Node[T]] {
	snapshot := []*Node[T]{}
	if t.size > 0 { // if tree is not empty
		t.preorderSubtree(t.root, snapshot)
	}
	return func(yield func(*Node[T]) bool) {
		for _, n := range snapshot {
			if !yield(n) {
				return
			}
		}
	}
}

// Adds general nodes of the subtree rooted at general node n to the given snapshot using
// the prostorder traversal algorithm.
func (t *Tree[T]) postorderSubtree(n *Node[T], snapshot []*Node[T]) {
	for _, c := range n.Children() {
		t.postorderSubtree(c, snapshot)
	}
	snapshot = append(snapshot, n) // postorder algorithm requires adding n after exploring subtrees
}

// Returns an iterator that performs a postorder traversal of the general tree.
func (t *Tree[T]) Postorder() iter.Seq[*Node[T]] {
	snapshot := []*Node[T]{}
	if t.size > 0 { // if tree is not empty
		t.postorderSubtree(t.root, snapshot)
	}
	return func(yield func(*Node[T]) bool) {
		for _, n := range snapshot {
			if !yield(n) {
				return
			}
		}
	}
}
