package cmap

// COLOR Definition
type COLOR bool
const (
	RED		COLOR = false
	BLACK	COLOR = true
)

// Red-Black Binary Tree
type RBTree struct {
	Root	*Node
	MLeft	*Node
	MRight	*Node
	Size 	uint32
	Index	uint8
}

// Red-Black Binary Tree Node
type Node struct {
	color 	COLOR
	Father	*Node
	Left 	*Node
	Right	*Node
	Key 	MapKey
	Value 	*MapData
}

// Generate a new Tree
func NewTree(index uint8) *RBTree {
	var rbTree = &RBTree{
		Root: 	nil,
		MLeft:	nil,
		MRight:	nil,
		Size:	0,
		Index:	index,
	}
	return rbTree
}

// Put method on RB-Tree with a key and value
func (rbTree *RBTree)Put(key MapKey, value *MapData) {
	// curr view tree node
	curr := rbTree.Root
	// upstairs tree node
	prev := curr
	// find insert position
	for curr != nil {
		// update prev
		prev = curr
		if key == curr.Key {
			// find position && replace
			rbTree.insert(curr, prev, key, value)
			return
		} else if key < curr.Key {
			// goto left
			curr = curr.Left
		} else {
			// goto right
			curr = curr.Right
		}
	}
	// insert
	rbTree.insert(curr, prev, key, value)
}

// Get method on RB-Tree with a key
func (rbTree *RBTree)Get(key MapKey) *MapData {
	// curr view tree node
	curr := rbTree.Root
	// find position
	for curr != nil {
		// update prev
		if key == curr.Key {
			// find
			return curr.Value
		} else if key < curr.Key {
			// goto left
			curr = curr.Left
		} else {
			// goto right
			curr = curr.Right
		}
	}
	return nil
}

// Delete method on RB-Tree with a key
func (rbTree *RBTree)Delete(key MapKey) *MapData {
	// curr view tree node
	curr := rbTree.Root
	// find insert position
	for curr != nil {
		if key == curr.Key {
			return rbTree.delete(curr)
		} else if key < curr.Key {
			// goto left
			curr = curr.Left
		} else {
			// goto right
			curr = curr.Right
		}
	}
	return nil
}

// insert a node into tree
// after insert must there must be a reshape operate
func (rbTree *RBTree)insert(curr *Node, prev *Node, key MapKey, value *MapData) {
	// replace value
	if curr != nil {
		curr.Value = value
		return
	}
	// new node
	var node = &Node{
		Left:	nil,
		Right: 	nil,
		Key:	key,
		Value: 	value,
	}
	if curr == prev {
		// insert root node
		node.Father = nil
		rbTree.Root = node
		rbTree.MLeft = node
		rbTree.MRight = node
	} else if key < prev.Key {
		// insert prev's left
		prev.Left = node
		node.Father = prev
		if prev == rbTree.MLeft {
			rbTree.MLeft = node
		}
	} else {
		// insert prev's right
		prev.Right = node
		node.Father = prev
		if prev == rbTree.MRight {
			rbTree.MRight = node
		}
	}
	rbTree.Size++
	rbTree.insertReshape(node)
}

// delete a node in tree with a point to node
// after delete must there must be a reshape operate
func (rbTree *RBTree)delete(curr *Node) *MapData {
	// get node's father and left, right child
	father := curr.Father
	left := curr.Left
	right := curr.Right
	// curr is a leaf node
	// just delete it
	if left == nil && right == nil {
		if curr == father.Left {
			if curr == rbTree.MLeft {
				rbTree.MLeft = father
			}
			father.Left = nil
		} else {
			if curr == rbTree.MRight {
				rbTree.MRight = father
			}
			father.Right = nil
		}
	} else {
		// curr's left or right child is nil
		// so we can just use curr's another child to fill the empty of curr's father
		if left == nil {
			// update rbtree
			// curr cannot be MRight because curr's right child is not nil
			if curr == rbTree.MLeft {
				rbTree.MLeft = mostLeftChild(right)
			}
			// move
			if curr == rbTree.Root {
				// root node
				rbTree.Root = right
				right.Father = nil
			} else if curr == father.Left {
				// curr is father's left child
				father.Left = right
				right.Father = father
			} else {
				// curr is father's right child
				father.Right = right
				right.Father = father
			}
		} else if right == nil {
			// update rbtree
			// curr cannot be MLeft because curr's left child is not nil
			if curr == rbTree.MRight {
				rbTree.MRight = mostRightChild(left)
			}
			if curr == rbTree.Root {
				// root node
				rbTree.Root = left
				right.Father = nil
			} else if curr == father.Left {
				// curr is father's left child
				father.Left = left
				right.Father = father
			} else {
				// curr is father's right child
				father.Right = left
				right.Father = father
			}
		} else {
			// left and right child all not nil
			// move left to right's most-left child's left
			var mright = mostLeftChild(right)
			mright.Left = left
			left.Father = mright
			// connect right to father
			right.Father = father
			if curr == father.Left {
				// curr is father's left child
				father.Left = right
			} else {
				// curr is father's right child
				father.Right = right
			}
			// update rbtree
			rbTree.MLeft = mostLeftChild(rbTree.Root)
			rbTree.MRight = mostRightChild(rbTree.Root)
		}
	}
	rbTree.Size --
	// TODO reshape

	return curr.Value
}

// reshape after insert node X
func (rbTree *RBTree)insertReshape(X *Node) {
	X.color = RED
	// X not root && X is red and his father is red too
	for X != rbTree.Root && X.Father.color == RED {
		if X.Father == X.Father.Father.Left {
			//      FF          FF
			//     /  \        /  \
			//    F    Y  or  F    Y
			//   /             \
			//  X               X
			// Y is FF's right child
			var Y = X.Father.Father.Right
			// Y is not nil and red
			if Y != nil && Y.color == RED {
				X.Father.color = BLACK
				Y.color = BLACK
				X.Father.Father.color = RED
				X = X.Father.Father
			} else {
				//    FF
				//   /  \
				//  F    Y
				//   \
				//    X
				// X is F's right child
				if X == X.Father.Right {
					// rotate left with pivot X's father
					X = X.Father
					rbTree.rotateLeft(X)
				}
				//      FF
				//     /  \
				//    F    Y
				//   /
				//  X
				// X is F's right child
				// change color
				X.Father.color = BLACK
				X.Father.Father.color = RED
				// rotate right with pivot X's grandfather
				rbTree.rotateRight(X.Father.Father)
			}
		} else {
			//      FF          FF
			//     /  \        /  \
			//    Y    F  or  Y    F
			//        /             \
			//       X               X
			// Y is FF's left child
			var Y = X.Father.Father.Left
			// Y is not nil and red
			if Y != nil && Y.color == RED {
				X.Father.color = BLACK
				Y.color = BLACK
				X.Father.Father.color = RED
				X = X.Father.Father
			} else {
				//       FF
				//      /  \
				//     Y    F
				//         /
				//        X
				// X is F's left child
				if X == X.Father.Left {
					// rotate right with pivot X's father
					X = X.Father
					rbTree.rotateRight(X)
				}
				//      FF
				//     /  \
				//    Y    F
				//          \
				//           X
				// X is F's right child
				// change color
				X.Father.color = BLACK
				X.Father.Father.color = RED
				// rotate left with pivot X's grandfather
				rbTree.rotateLeft(X.Father.Father)
			}
		}
	}
	// root is black
	rbTree.Root.color = BLACK
}

// Left Rotate with pivot X
//
// 	    R                  R
//      |                  |
//      X                  Y
//       \       =>       / \
//        Y              X   Z
//       / \              \
//      W   Z              W
//
// make tree X to a balance binary tree
func (rbTree *RBTree)rotateLeft(X *Node) {
	var Y = X.Right
	// move W to X's right child
	X.Right = Y.Left
	if Y.Left != nil {
		Y.Left.Father = X
	}
	// move Y to X's father's child
	Y.Father = X.Father
	if X == rbTree.Root {
		// X is root
		rbTree.Root = Y
	} else if X == X.Father.Left {
		// X is father's left child
		X.Father.Left = Y
	} else {
		// X is father's right child
		X.Father.Right = Y
	}
	// move X to Y's left child
	Y.Left = X
	X.Father = Y
}

// Right Rotate with pivot X
//
// 	     R                R
//       |                |
//       X                Y
//      /       =>       / \
//     Y                Z   X
//    / \                  /
//   Z   W                W
//
// make tree X to a balance binary tree
func (rbTree *RBTree)rotateRight(X *Node) {
	var Y = X.Left
	// move W to X's left child
	X.Left = Y.Right
	if Y.Right != nil {
		Y.Right.Father = X
	}
	// move Y to X's father's child
	Y.Father = X.Father
	if X == rbTree.Root {
		// X is root
	 	rbTree.Root = Y
	} else if X == X.Father.Right {
		// X is father's left child
		X.Father.Right = Y
	} else {
		// X is father's right child
		X.Father.Left = Y
	}
	// move X to Y's right child
	Y.Right = X
	X.Father = Y
}

// get most left child
func mostLeftChild(root *Node) *Node {
	if root == nil {
		return root
	}
	p := root
	for p.Left != nil {
		p = p.Left
	}
	return p
}

// get most right child
func mostRightChild(root *Node) *Node {
	if root == nil {
		return root
	}
	p := root
	for p.Right != nil {
		p = p.Right
	}
	return p
}