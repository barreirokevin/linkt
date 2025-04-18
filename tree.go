package main

import (
	"errors"
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

// Returns an array that contains the nodes in the tree after having traversed it
// with a preorder algorithm.
func (t *Tree[T]) Preorder() []*Node[T] {
	snapshot := []*Node[T]{}

	// define recursive preorder traversal algorithm
	var preorderSubtree func(n *Node[T], snapshot []*Node[T])
	preorderSubtree = func(n *Node[T], snapshot []*Node[T]) {
		// preorder requires adding n before exploring subtrees
		snapshot = append(snapshot, n)
		for _, c := range n.Children() {
			preorderSubtree(c, snapshot)
		}
	}

	if t.size > 0 { // if tree is not empty
		preorderSubtree(t.root, snapshot)
	}

	return snapshot
}

// Returns an array that contains the nodes in the tree after having traversed it
// with a postorder algorithm.
func (t *Tree[T]) Postorder() []*Node[T] {
	snapshot := []*Node[T]{}

	// define recursive postorder traversal algorithm
	var postorderSubtree func(n *Node[T], snapshot []*Node[T])
	postorderSubtree = func(n *Node[T], snapshot []*Node[T]) {
		for _, c := range n.Children() {
			postorderSubtree(c, snapshot)
		}
		// postorder requires adding n after exploring subtrees
		snapshot = append(snapshot, n)
	}

	if t.size > 0 { // if tree is not empty
		postorderSubtree(t.root, snapshot)
	}

	return snapshot
}
