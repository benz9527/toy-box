package tree

// Binary Sort Tree + Heap = Treap（树堆）
// Treap 不等同于 Binary Heap
// 1. Binary Heap 必须是一个完全二叉树， Treap 则不一定是完全二叉树
// 2. Treap 不严格满足 AVL Tree 的要求，因为 Treap 左右子树的高度差绝对值可能会超过 1
//  只能近似满足 AVL Tree 的性质。下面的实现也能看出，进行调整的旋转只进行了一次，而没有
//  AVL Tree 的连续调整。Treap 的旋转调整是根据优先级完成的，就是为了实现随机平衡结构，
//  也能实现快速查找。

type TreapNode struct {
	Left, Right *TreapNode
	// @Field Priority 随机优先级
	Priority, Key int
}

func (t *TreapNode) RotateLeft() *TreapNode {
	// FIXME(Ben) nil ptr
	r := t.Right
	t.Right = r.Left
	r.Left = t
	return r
}

func (t *TreapNode) RotateRight() *TreapNode {
	// FIXME(Ben) nil ptr
	l := t.Left
	t.Left = l.Right
	l.Right = t
	return l
}

type Treap struct {
	Root *TreapNode
}

func (t *Treap) Insert(key, priority int) {
	t.insert(t.Root, key, priority)
}

func (t *Treap) insert(root *TreapNode, key, priority int) {
	if root == nil {
		root = &TreapNode{
			Left:     nil,
			Right:    nil,
			Key:      key,
			Priority: priority,
		}
		return
	}
	// classical dichotomy & recursive
	if key < root.Key {
		// Insert into left node
		t.insert(root.Left, key, priority)
		// rotate adjust
		if root.Left.Priority < root.Priority {
			root = root.RotateRight()
		}
	} else {
		// Insert into right node
		t.insert(root.Right, key, priority)
		// rotate adjust
		if root.Right.Priority < root.Priority {
			root = root.RotateLeft()
		}
	}
}

func (t *Treap) Delete(key int) {
	t.delete(t.Root, key)
}

func (t *Treap) delete(root *TreapNode, key int) {
	if root == nil {
		return
	}

	if key < root.Key {
		t.delete(root.Left, key)
	} else if key > root.Key {
		t.delete(root.Right, key)
	} else {
		if root.Left == nil {
			root = root.Right
		} else if root.Right == nil {
			root = root.Left
		} else {
			if root.Left.Priority < root.Right.Priority {
				root = root.RotateRight()
				t.delete(root.Right, key)
			} else {
				root = root.RotateLeft()
				t.delete(root.Left, key)
			}
		}
	}
}
