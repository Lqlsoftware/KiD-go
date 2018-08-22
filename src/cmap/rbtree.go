package cmap

// COLOR Definition
type COLOR bool
const (
	RED		COLOR = false
	BLACK	COLOR = true
)

// Red-Black Binary Tree
type RBTree struct {
	root   *RBtree
	mLeft  *RBtree
	mRight *RBtree
	// attributes
	Size   uint32
	Index  uint8
}

// Red-Black Binary Tree RBtree
type RBtree struct {
	color  COLOR
	father *RBtree
	left   *RBtree
	right  *RBtree
	// key & data
	Key    MapKey
	Value  *MapData
}

// Generate a new Tree
func NewTree(index uint8) *RBTree {
	var rbTree = &RBTree{
		root:   nil,
		mLeft:  nil,
		mRight: nil,
		Size:   0,
		Index:  index,
	}
	return rbTree
}

// Put method on RB-Tree with a key and value
func (rbTree *RBTree)Put(key MapKey, value *MapData) {
	// curr view tree node
	curr := rbTree.root
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
			curr = curr.left
		} else {
			// goto right
			curr = curr.right
		}
	}
	// insert
	rbTree.insert(curr, prev, key, value)
}

// Get method on RB-Tree with a key
func (rbTree *RBTree)Get(key MapKey) *MapData {
	// curr view tree node
	curr := rbTree.root
	// find position
	for curr != nil {
		// update prev
		if key == curr.Key {
			// find
			return curr.Value
		} else if key < curr.Key {
			// goto left
			curr = curr.left
		} else {
			// goto right
			curr = curr.right
		}
	}
	return nil
}

// Delete method on RB-Tree with a key
func (rbTree *RBTree)Delete(key MapKey) *MapData {
	// curr view tree node
	curr := rbTree.root
	// find insert position
	for curr != nil {
		if key == curr.Key {
			return rbTree.delete(curr)
		} else if key < curr.Key {
			// goto left
			curr = curr.left
		} else {
			// goto right
			curr = curr.right
		}
	}
	return nil
}

// insert a node into tree
// after insert must there must be a reshape operate
func (rbTree *RBTree)insert(curr *RBtree, prev *RBtree, key MapKey, value *MapData) {
	// replace value
	if curr != nil {
		curr.Value = value
		return
	}
	// new node
	var node = &RBtree{
		left:  nil,
		right: nil,
		Key:   key,
		Value: value,
	}
	if curr == prev {
		// insert root node
		node.father = nil
		rbTree.root = node
		rbTree.mLeft = node
		rbTree.mRight = node
	} else if key < prev.Key {
		// insert prev's left
		prev.left = node
		node.father = prev
		if prev == rbTree.mLeft {
			rbTree.mLeft = node
		}
	} else {
		// insert prev's right
		prev.right = node
		node.father = prev
		if prev == rbTree.mRight {
			rbTree.mRight = node
		}
	}
	rbTree.Size++
	rbTree.insertReshape(node)
}

// delete a node in tree with a point to node
// after delete must there must be a reshape operate
func (rbTree *RBTree)delete(curr *RBtree) *MapData {
	// get node's father and left, right child
	father := curr.father
	left := curr.left
	right := curr.right
	// curr is a leaf node
	// just delete it
	if left == nil && right == nil {
		if curr == father.left {
			if curr == rbTree.mLeft {
				rbTree.mLeft = father
			}
			father.left = nil
		} else {
			if curr == rbTree.mRight {
				rbTree.mRight = father
			}
			father.right = nil
		}
	} else {
		// curr's left or right child is nil
		// so we can just use curr's another child to fill the empty of curr's father
		if left == nil {
			// update rbtree
			// curr cannot be mRight because curr's right child is not nil
			if curr == rbTree.mLeft {
				rbTree.mLeft = mostLeftChild(right)
			}
			// move
			right.father = father
			if curr == rbTree.root {
				// root node
				rbTree.root = right
			} else if curr == father.left {
				// curr is father's left child
				father.left = right
			} else {
				// curr is father's right child
				father.right = right
			}
		} else if right == nil {
			// update rbtree
			// curr cannot be mLeft because curr's left child is not nil
			if curr == rbTree.mRight {
				rbTree.mRight = mostRightChild(left)
			}
			left.father = father
			if curr == rbTree.root {
				// root node
				rbTree.root = left
			} else if curr == father.left {
				// curr is father's left child
				father.left = left
			} else {
				// curr is father's right child
				father.right = left
			}
		} else {
			// left and right child all not nil
			// move left to right's most-left child's left
			var mright = mostLeftChild(right)
			mright.left = left
			left.father = mright
			// connect right to father
			right.father = father
			if curr == rbTree.root {
				// root node
				rbTree.root = right
			} else if curr == father.left {
				// curr is father's left child
				father.left = right
			} else {
				// curr is father's right child
				father.right = right
			}
			// update rbtree
			rbTree.mLeft = mostLeftChild(rbTree.root)
			rbTree.mRight = mostRightChild(rbTree.root)
		}
	}
	rbTree.Size --
	// TODO reshape

	return curr.Value
}

// reshape after insert node X
func (rbTree *RBTree)insertReshape(X *RBtree) {
	X.color = RED
	// X not root && X is red and his father is red too
	for X != rbTree.root && X.father.color == RED {
		if X.father == X.father.father.left {
			//      FF          FF
			//     /  \        /  \
			//    F    Y  or  F    Y
			//   /             \
			//  X               X
			// Y is FF's right child
			var Y = X.father.father.right
			// Y is not nil and red
			if Y != nil && Y.color == RED {
				X.father.color = BLACK
				Y.color = BLACK
				X.father.father.color = RED
				X = X.father.father
			} else {
				//    FF
				//   /  \
				//  F    Y
				//   \
				//    X
				// X is F's right child
				if X == X.father.right {
					// rotate left with pivot X's father
					X = X.father
					rbTree.rotateLeft(X)
				}
				//      FF
				//     /  \
				//    F    Y
				//   /
				//  X
				// X is F's right child
				// change color
				X.father.color = BLACK
				X.father.father.color = RED
				// rotate right with pivot X's grandfather
				rbTree.rotateRight(X.father.father)
			}
		} else {
			//      FF          FF
			//     /  \        /  \
			//    Y    F  or  Y    F
			//        /             \
			//       X               X
			// Y is FF's left child
			var Y = X.father.father.left
			// Y is not nil and red
			if Y != nil && Y.color == RED {
				X.father.color = BLACK
				Y.color = BLACK
				X.father.father.color = RED
				X = X.father.father
			} else {
				//       FF
				//      /  \
				//     Y    F
				//         /
				//        X
				// X is F's left child
				if X == X.father.left {
					// rotate right with pivot X's father
					X = X.father
					rbTree.rotateRight(X)
				}
				//      FF
				//     /  \
				//    Y    F
				//          \
				//           X
				// X is F's right child
				// change color
				X.father.color = BLACK
				X.father.father.color = RED
				// rotate left with pivot X's grandfather
				rbTree.rotateLeft(X.father.father)
			}
		}
	}
	// root is black
	rbTree.root.color = BLACK
}

// left Rotate with pivot X
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
func (rbTree *RBTree)rotateLeft(X *RBtree) {
	var Y = X.right
	// move W to X's right child
	X.right = Y.left
	if Y.left != nil {
		Y.left.father = X
	}
	// move Y to X's father's child
	Y.father = X.father
	if X == rbTree.root {
		// X is root
		rbTree.root = Y
	} else if X == X.father.left {
		// X is father's left child
		X.father.left = Y
	} else {
		// X is father's right child
		X.father.right = Y
	}
	// move X to Y's left child
	Y.left = X
	X.father = Y
}

// right Rotate with pivot X
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
func (rbTree *RBTree)rotateRight(X *RBtree) {
	var Y = X.left
	// move W to X's left child
	X.left = Y.right
	if Y.right != nil {
		Y.right.father = X
	}
	// move Y to X's father's child
	Y.father = X.father
	if X == rbTree.root {
		// X is root
	 	rbTree.root = Y
	} else if X == X.father.right {
		// X is father's left child
		X.father.right = Y
	} else {
		// X is father's right child
		X.father.left = Y
	}
	// move X to Y's right child
	Y.right = X
	X.father = Y
}

// get most left child
func mostLeftChild(root *RBtree) *RBtree {
	if root == nil {
		return root
	}
	p := root
	for p.left != nil {
		p = p.left
	}
	return p
}

// get most right child
func mostRightChild(root *RBtree) *RBtree {
	if root == nil {
		return root
	}
	p := root
	for p.right != nil {
		p = p.right
	}
	return p
}