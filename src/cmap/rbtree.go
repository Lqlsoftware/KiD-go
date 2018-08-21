package cmap

type COLOR bool
const (
	RED		COLOR = false
	BLACK	COLOR = true
)

// Red-Black BinaryTree
type Node struct {
	color 	COLOR
	Key 	uint32
	Left 	*Node
	Right	*Node

	Value 	interface{} // TODO
}

func NewTree() *Node {
	var root = &Node{
		color:	BLACK,
	}

}