package main

// A node in a general tree.
type Node[T any] struct {
	element  T
	parent   *Node[T]
	children []*Node[T]
}

// Creates and returns a new general node.
func NewNode[T any]() *Node[T] {
	return &Node[T]{
		parent:   nil,
		children: nil,
	}
}

// Returns the element stored in this node.
func (n Node[T]) GetElement() T {
	return n.element
}

// Returns the parent node of this node.
func (n Node[T]) GetParent() *Node[T] {
	return n.parent
}

// Sets element enas this nodes element.
func (n Node[T]) SetElement(e T) {
	n.element = e
}

// Sets parent p as this nodes parent.
func (n Node[T]) SetParent(p *Node[T]) {
	n.parent = p
}

// Returns a slice of nodes that are the children of this node.
func (n Node[T]) Children() []*Node[T] {
	snapshot := []*Node[T]{}
	for _, c := range n.children {
		snapshot = append(snapshot, c)
	}
	return snapshot
}

// Returns true if this node has children, otherwise false.
func (n Node[T]) IsInternal() bool {
	return len(n.Children()) > 0
}

// Returns true if this node does not have children, otherwise false.
func (n Node[T]) IsExternal() bool {
	return len(n.Children()) == 0
}
