package cmap

// COLOR Definition
type COLOR bool
const (
	RED		COLOR = false
	BLACK	COLOR = true
)

// Red-Black Binary Tree
type RBTree struct {
	root   *treeNode
	mLeft  *treeNode
	mRight *treeNode
	// attributes
	Size   uint32
	Index  uint8
}

// Red-Black Binary Tree treeNode
type treeNode struct {
	color  COLOR
	father *treeNode
	left   *treeNode
	right  *treeNode
	// key & data
	Key    MapKey
	Value  *MapValue
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
func (rbTree *RBTree)Put(key MapKey, value *MapValue) {
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
func (rbTree *RBTree)Get(key MapKey) *MapValue {
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
func (rbTree *RBTree)Delete(key MapKey) *MapValue {
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
func (rbTree *RBTree)insert(curr *treeNode, prev *treeNode, key MapKey, value *MapValue) {
	// replace value
	if curr != nil {
		curr.Value = value
		return
	}
	// new node
	var node = &treeNode{
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
func (rbTree *RBTree)delete(curr *treeNode) *MapValue {
	var value = curr.Value
	// Both left child & right child is not nil
	//     move right's most-left child(successor) to curr
	//     and change curr to mright to make problem easy
	if curr.left != nil && curr.right != nil {
		var mright = mostLeftChild(curr.right)
		curr.Value = mright.Value
		curr.Key = mright.Key
		curr = mright
	}
	// get node's father and left, right child
	var father, left, right = curr.father, curr.left, curr.right
	// replace node
	var replace *treeNode = nil
	if left != nil {
		// curr's left child is nil
		//     so we can just use curr's right child to fill the empty of curr's father
		replace = left
		replace.father = father
		// update rbtree
		// curr cannot be mLeft because curr's left child is not nil
		if curr == rbTree.mRight {
			rbTree.mRight = mostRightChild(left)
		}
	} else if right != nil {
		// curr's right child is nil
		//     so we can just use curr's left child to fill the empty of curr's father
		replace = right
		replace.father = father
		// update rbtree
		// curr cannot be mRight because curr's right child is not nil
		if curr == rbTree.mLeft {
			rbTree.mLeft = mostLeftChild(right)
		}
	} else {
		// curr is a leaf curr
		// update rbtree
		if curr == rbTree.root {
			rbTree.mLeft = nil
			rbTree.mRight = nil
		} else if curr == father.left && curr == rbTree.mLeft {
			rbTree.mLeft = father
		} else if curr == father.right && curr == rbTree.mRight {
			rbTree.mRight = father
		}
	}
	// change father's attribute
	switch curr {
	case rbTree.root:	rbTree.root = replace
	case father.left: 	father.left = replace
	case father.right: 	father.right = replace
	}
	rbTree.Size --
	// reshape
	if curr.color == BLACK {
		if replace != nil {
			rbTree.deleteReshape(replace)
		} else {
			rbTree.deleteReshape(curr)
		}
	}
	return value
}

// reshape after delete node X
// There are 4 probably situation:
// 1. X's brother is RED
//    Step 1: change X's father to RED, brother to BLACK
//    Step 2: rotateLeft with X's father as pivot
//
//         |                              |
//       1 ●                            3 ●
//        / \                            / \
// X-> 2 ●   ○ 3 <-brother    =>      1 ○   ● 5
//          / \                        / \
//       4 ●   ● 5              X-> 2 ●   ● 4
//
//
// 2. X's brother is BLACK and brother's two child is BLACK
//    Step 1: change brother to RED
//    Step 2: set X point to X's father
//
//         |                              |
//       1 ○                        X-> 1 ○
//        / \                            / \
// X-> 2 ●   ● 3 <-brother    =>      2 ●   ○ 3
//          / \                            / \
//       4 ●   ● 5                      4 ●   ● 5
//
//
// 3. X's brother is BLACK and brother's left child is RED, right child is BLACK
//    Step 1: change brother to RED, brother's left child to BLACK
//    Step 2: rotateRight with brother as pivot
//
//         |                             |
//       1 ○                           1 ○
//        / \                           / \
// X-> 2 ●   ● 3 <-brother    => X-> 2 ●   ● 4
//          / \                             \
//       4 ○   ● 5                           ○ 3
//                                            \
//                                             ● 5
//
//
// 4. X's brother is BLACK and brother's right child is RED
//    Step 1: change X's father to BLACK, brother to RED, brother's right child to BLACK
//    Step 2: rotateLeft with X's father as pivot
//
//         |                              |
//       1 ○                            3 ○
//        / \                            / \
// X-> 2 ●   ● 3 <-brother    =>      1 ●   ● 5
//          / \                        / \
//       4 ○   ○ 5                  2 ●   ○ 4
//
//
func (rbTree *RBTree)deleteReshape(X *treeNode) {
	for X != rbTree.root {
		if X == X.father.left && X.father.right != nil {
			var brother = X.father.right
			if brother.color == RED {
				X.father.color = RED
				brother.color = BLACK
				rbTree.rotateLeft(X.father)
			} else if brother.left != nil && brother.left.color == BLACK && brother.right != nil && brother.right.color == BLACK {
				brother.color = RED
				X = X.father
			} else if brother.left != nil && brother.left.color == RED && brother.right != nil && brother.right.color == BLACK {
				brother.color = RED
				brother.left.color = BLACK
				rbTree.rotateRight(brother)
			} else if brother.right != nil && brother.right.color == RED {
				brother.father.color = BLACK
				brother.color = RED
				brother.right.color = BLACK
				rbTree.rotateLeft(brother.father)
				break
			}
		} else if X == X.father.right && X.father.left != nil {
			var brother = X.father.left
			if brother.color == RED {
				brother.color = BLACK
				X.father.color = RED
				rbTree.rotateRight(X.father)
			} else if brother.left != nil && brother.left.color == BLACK && brother.right != nil && brother.right.color == BLACK {
				brother.color = RED
				X = X.father
			} else if brother.left != nil && brother.left.color == BLACK && brother.right != nil && brother.right.color == RED {
				brother.color = RED
				brother.right.color = BLACK
				rbTree.rotateLeft(brother)
			} else if brother.left != nil && brother.left.color == RED {
				brother.color = RED
				brother.left.color = BLACK
				brother.father.color = BLACK
				rbTree.rotateRight(brother.father)
				X = rbTree.root
			}
		} else {
			break
		}
	}
}

// reshape after insert node X
func (rbTree *RBTree)insertReshape(X *treeNode) {
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
func (rbTree *RBTree)rotateLeft(X *treeNode) {
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
func (rbTree *RBTree)rotateRight(X *treeNode) {
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
func mostLeftChild(root *treeNode) *treeNode {
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
func mostRightChild(root *treeNode) *treeNode {
	if root == nil {
		return root
	}
	p := root
	for p.right != nil {
		p = p.right
	}
	return p
}