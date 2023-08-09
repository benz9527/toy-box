package tree

import (
	"fmt"
)

type BTreeNode struct {
	// current number of keys
	N int
	// key values array
	Keys []int
	// child pointers array (插空数量)
	Children []*BTreeNode
	IsLeaf   bool
}

func (n *BTreeNode) Find(key int) (idx int) {
	idx = 0
	for idx < n.N && n.Keys[idx] < key {
		idx++
	}
	return idx
}

type BTree struct {
	// 整棵树的 degree，限制 key 的数量
	T    int
	Root *BTreeNode
}

func NewBTree(t int) *BTree {
	btree := &BTree{T: t}
	// 设置 key 的最大数量为 2T - 1， child nodes 的最大数量为 2T
	root := &BTreeNode{N: 0, IsLeaf: true, Keys: make([]int, 2*t-1), Children: make([]*BTreeNode, 2*t)}
	btree.Root = root
	return btree
}

func (btree *BTree) Search(n *BTreeNode, key int) *BTreeNode {
	i := 0
	if n == nil {
		return n
	}

	for ; i < n.N; i++ {
		// 遍历完同一层级的所有key
		if key < n.Keys[i] {
			break
		}

		if key == n.Keys[i] {
			return n
		}
	}

	if n.IsLeaf {
		return nil
	} else {
		return btree.Search(n.Children[i], key)
	}
}

func (btree *BTree) Split(nx *BTreeNode, pos int, ny *BTreeNode) {
	// x node is the upper level node
	// y node is the left part
	// z node is the right part
	// 该新节点，准备接收 y node（处于满状态节点）的 T - 1 个 keys
	z := &BTreeNode{
		N:        btree.T - 1,
		IsLeaf:   ny.IsLeaf,
		Keys:     make([]int, 2*btree.T-1),
		Children: make([]*BTreeNode, 2*btree.T),
	}
	// 把后 T - 1 个 key 数据复制到 z node
	for i := 0; i < btree.T-1; i++ {
		z.Keys[i] = ny.Keys[i+btree.T]
	}
	// 如果是叶子节点，它根本没有 children 能够移动
	if !ny.IsLeaf {
		// y node 如果是非叶子节点，把后 T（也就是一半的）个 children 节点复制到 z node
		for i := 0; i < btree.T; i++ {
			z.Children[i] = ny.Children[i+btree.T]
		}
	}
	// 更新 y node 的 key 数量为 T - 1, 也就是 y node children & keys 还保留着原本的信息
	// 但是更新了现有的 keys 的数量信息，这些被保留的冗余会在之后的 insert 中被覆盖掉
	ny.N = btree.T - 1

	// x node 准备拥有新的 children nodes，腾挪位置
	for i := nx.N; i >= pos+1; i-- {
		nx.Children[i+1] = nx.Children[i]
	}
	// y node 本身就在 children 里面了，所以只需要把 z node 加入到 children 就行
	nx.Children[pos+1] = z

	// keys 也要腾挪位置
	for i := nx.N - 1; i >= pos; i-- {
		nx.Keys[i+1] = nx.Keys[i]
	}
	// copy the middle key of y node to this node
	nx.Keys[pos] = ny.Keys[btree.T-1]
	nx.N = nx.N + 1
}

func (btree *BTree) Insert(key int) {
	// 从 root 插入
	r := btree.Root
	// 检查 root 的 key 数量是否达到最大，达到最大时触发新的 root 选取
	if r.N == 2*btree.T-1 {
		// 已经达到最大，需要选出新的 root， 就是中间值
		s := &BTreeNode{IsLeaf: true, Keys: make([]int, 2*btree.T-1), Children: make([]*BTreeNode, 2*btree.T)}
		// 新的 root
		btree.Root = s
		s.IsLeaf = false
		// 临时指向原本的 root 首节点
		s.Children[0] = r
		// 当前 r 是处于节点的满状态， key 的数量达到最大
		btree.Split(s, 0, r)
		btree.InsertValue(s, key)
	} else {
		btree.InsertValue(r, key)
	}
}

func (btree *BTree) InsertValue(n *BTreeNode, key int) {
	// initialize index as rightmost element index
	i := n.N - 1
	if n.IsLeaf {
		// 叶子节点的操作

		// 插入的新值比较小，要不断往前找合适的位置
		// a) finds the location of new key to be inserted
		// b) moves all greater keys to one place ahead
		for ; i >= 0 && key < n.Keys[i]; i-- {
			n.Keys[i+1] = n.Keys[i]
		}
		// 找到合适的位置时就填入值
		n.Keys[i+1] = key
		n.N = n.N + 1
	} else {
		// 非叶子节点的操作
		// 从后往前找位置，但是不腾挪数据先
		// find the child which is going to have the new key
		for ; i >= 0 && key < n.Keys[i]; i-- {
			// 当当前 new key 的值比这里所有的 keys 的值都要小，只能往更小的 child 去插入
		}
		// 找到了合适的 child 位置，（如果 children is full）这个也会是调整后 keys 的新索引
		i++
		temp := n.Children[i]
		// see if the found child is full
		// children 满了就要产生一个 parent node 添加到上一层
		if temp.N == 2*btree.T-1 {
			btree.Split(n, i, temp)
			// after split, the middle key of children[i] goes up and
			// children[i] is split into two. See which of the two
			// is going to have the new key
			if key > n.Keys[i] {
				i++
			}
		}
		btree.InsertValue(n.Children[i], key)
	}
}

func Show(n *BTreeNode) {
	// 先序遍历
	if n == nil {
		fmt.Println("The node can not be nil!")
		return
	}

	for i := 0; i < n.N; i++ {
		fmt.Print(n.Keys[i], " ")
	}

	if !n.IsLeaf {
		for i := 0; i < n.N+1; i++ {
			Show(n.Children[i])
		}
	}
}

func (btree *BTree) Show() {
	Show(btree.Root)
}

func PreorderShow(n *BTreeNode) {
	i := 0
	for i = 0; i < n.N; i++ {
		if !n.IsLeaf {
			PreorderShow(n.Children[i])
		}
		fmt.Print(n.Keys[i], " ")
	}

	if !n.IsLeaf {
		PreorderShow(n.Children[i])
	}
}

func (btree *BTree) PreorderShow() {
	PreorderShow(btree.Root)
	return
}

func (btree *BTree) Contain(k int) bool {
	if btree.Search(btree.Root, k) != nil {
		return true
	} else {
		return false
	}
}

func (btree *BTree) RemoveKey(k int) {
	if btree.Root == nil {
		return
	}

	btree.Remove(btree.Root, k)
	if btree.Root.N == 0 {
		if btree.Root.IsLeaf {
			btree.Root = nil
		} else {
			btree.Root = btree.Root.Children[0]
		}
	}
}

func (btree *BTree) Remove(parentNode *BTreeNode, k int) {
	// 执行删除的之前，先要分清楚不同的场景
	// 1。是否为叶子结点
	// 是：
	// 叶子结点的 keys 是否大于 btree properties 中要求的数量至少是 ceiling([t-1]/2)
	// -- 是：可以直接删除
	// -- 否：需要通过 fill/merge 填充该叶子节点的 keys 数量 > ceiling([t-1]/2) 才不会在删除后违反 btree properties \
	// fill 需要 inorder predecessor & inorder successor 来， pred 和 succ 主要是填充叶子节点的 parent node， 而
	// parent node 的 key 是用来填充叶子节点
	// merge 是当 pred & succ 的数量都刚好满足 == ceiling([t-1]/2)，不能发生借 key 的行为，那就固定选择一边的 parent key &
	// sibling keys 来合并成为一个大的叶子节点，再执行删除
	// 否：
	// 基本就是要进行借 key 填充的操作或者 merge
	// if the key k is in node x (parent node) and x is an internal node
	// 1) if the child y that precedes k in node x has at least keys, then find the predecessor k0 of k in the sub-tree
	// rooted at y. recursively delete k0, and replace k by k0 in x.(we can find k0 and delete it in a single downward pass)
	// 2) if y has fewer than t keys, then, symmetrically, examine the child z that follows k in node x, if z has at least
	// t keys, then find the successors k0 of k in the subtree rooted at z, recursively delete k0, and replace k by k0 in x.
	// 3) otherwise, if both y and z have only t-1 keys, merge k and all of z into y, so that x losses both k and the pointer to z
	// and y now contains 2t-1 keys. then free z and recursively delete k from y

	// 在当前节点里找 k
	idx := parentNode.Find(k)

	// k 在当前节点里找到
	if idx < parentNode.N && parentNode.Keys[idx] == k {
		if parentNode.IsLeaf {
			btree.RemoveFromLeaf(parentNode, idx)
		} else {
			btree.RemoveFromNonLeaf(parentNode, idx)
		}
	} else {
		// k 在当前节点里找不到，就可能要转向 child nodes 里去找
		if parentNode.IsLeaf {
			// 这个直接不用找
			return
		}

		// the key to be removed is present in the sub-tree rooted with this node
		// the flag indicates whether the key is present in the sub-tree rooted
		// with the last child of this node
		flag := false
		if idx == parentNode.N {
			// 相等说明是最右边的子树，不等就是中间或者最前的子树
			flag = true
		}

		// if the child where the key is supposed to exist has less that t keys,
		// we fill that child
		if parentNode.Children[idx].N < btree.T {
			// 先调整不符合要求的子树再删除
			// 这里存在 merge 的可能
			btree.Fill(parentNode, idx)
		}

		// if the last child has been merged, it must have merged with the previous
		// child and so we recurse on the (idx-1)th child.
		// else, we recurse on the (idx)th child which now has at least t keys
		if flag && idx > parentNode.N {
			btree.Remove(parentNode.Children[idx-1], k)
		} else {
			btree.Remove(parentNode.Children[idx], k)
		}
	}
	return
}

// RemoveFromLeaf remove the key present in idx-th position in this node which is a leaf
func (btree *BTree) RemoveFromLeaf(n *BTreeNode, idx int) {
	// move all the keys after the idx-th pos one place backward
	for i := idx + 1; i < n.N; i++ {
		n.Keys[i-1] = n.Keys[i]
	}
	// reduce the count of keys
	n.N--
	return
}

// RemoveFromNonLeaf remove the key present in idx-th position in this node which is a leaf
func (btree *BTree) RemoveFromNonLeaf(n *BTreeNode, idx int) {
	k := n.Keys[idx]

	// if the child that precedes k (children[idx]) has at least t keys,
	// find the predecessor 'pred' of k in the subtree rooted at children[idx].
	// replace k by pred. recursively delete pred in children[idx]
	if n.Children[idx].N >= btree.T {
		pred := btree.GetPred(n, idx)
		n.Keys[idx] = pred
		btree.Remove(n.Children[idx], pred)
	} else if n.Children[idx+1].N >= btree.T {
		// if the child children[idx] has less that t keys, examine children[idx + 1].
		// if children[idx + 1] has at least t keys, find the successor 'succ' of k in
		// the subtree rooted at children[idx + 1]
		// replace k by succ
		// recursively delete succ in children[idx + 1]
		succ := btree.GetSucc(n, idx)
		n.Keys[idx] = succ
		btree.Remove(n.Children[idx+1], succ)
	} else {
		// if both children[idx + 1] and children[idx + 1] has less that t keys, merge k and all of
		// into children[idx]
		// now children[idx] contains 2t - 1 keys
		// free children[idx + 1] and recursively delete k from children[idx]
		btree.Merge(n, idx)
		btree.Remove(n.Children[idx], k)
	}
}

// GetPred get the predecessor of the key, where the key is present in the idx-th position in the node
func (btree *BTree) GetPred(n *BTreeNode, idx int) int {
	// keep moving to the right most node until we reach a leaf
	cur := n.Children[idx]
	for !cur.IsLeaf {
		cur = cur.Children[cur.N]
	}
	// return the last key of the leaf
	return cur.Keys[cur.N-1]
}

// GetSucc get the successor of the key, where the key is present in the idx-th position in the node
func (btree *BTree) GetSucc(n *BTreeNode, idx int) int {
	// keep moving the left most node starting from children[idx + 1] until we reach a leaf
	cur := n.Children[idx+1]
	for !cur.IsLeaf {
		cur = cur.Children[0]
	}
	// return the first key of the leaf
	return cur.Keys[0]
}

// Fill to fill up the child node present in the idx-th position in the children[] array
// if that child has less than T - 1 keys
func (btree *BTree) Fill(parentNode *BTreeNode, idx int) {
	// if the previous children[idx - 1] has more than t - 1 keys, borrow a key
	// from that child
	if idx != 0 && parentNode.Children[idx-1].N >= btree.T {
		btree.BorrowFromPrev(parentNode, idx)
	} else if idx != parentNode.N && parentNode.Children[idx+1].N >= btree.T {
		// if the next children[idx + 1] has more than t - 1 keys, borrow a key
		// from that child
		btree.BorrowFromNext(parentNode, idx)
	} else {
		if idx != parentNode.N {
			btree.Merge(parentNode, idx)
		} else {
			btree.Merge(parentNode, idx)
		}
	}
	return
}

// BorroweFromPrev to borrow a key from the children[idx - 1]-th node and
// place it in children[idx]-th node
func (btree *BTree) BorrowFromPrev(parentNode *BTreeNode, idx int) {
	deletedChild := parentNode.Children[idx]
	leftSibling := parentNode.Children[idx-1]

	// the last key from children[idx - 1] goes up to the parent and keys[idx - 1]
	// from parent is inserted as the first key in children[idx].
	// thus, the losses leftSibling one key and deletedChild gains one key

	// move one step ahead
	// 留出空位，（deleted child）left most 接收 new key from parent & new child from sibling （sibling right most）
	for i := deletedChild.N - 1; i >= 0; i-- {
		deletedChild.Keys[i+1] = deletedChild.Keys[i]
	}

	if !deletedChild.IsLeaf {
		for i := deletedChild.N; i >= 0; i-- {
			deletedChild.Children[i+1] = deletedChild.Children[i]
		}
	}

	// setting deletedChild's first key equal to keys[idx - 1] from the current node
	deletedChild.Keys[0] = parentNode.Keys[idx-1]

	// moving leftSibling's last deletedChild as children[idx]'s first deletedChild
	if !deletedChild.IsLeaf {
		deletedChild.Children[0] = leftSibling.Children[leftSibling.N]
	}

	// move the key at right most in left sibling to parent
	parentNode.Keys[idx-1] = leftSibling.Keys[leftSibling.N-1]
	deletedChild.N += 1
	leftSibling.N -= 1
	return
}

// BorrowFromNext to borrow a key from the children[idx + 1]-th node and
// place it in children[idx]-th node
func (btree *BTree) BorrowFromNext(parentNode *BTreeNode, idx int) {
	deletedChild := parentNode.Children[idx]
	rightSibling := parentNode.Children[idx+1]

	// 留出空位，（deleted child）right most 接收 new key from parent & new child from sibling （sibling left most index）
	deletedChild.Keys[deletedChild.N] = parentNode.Keys[idx]

	// rightSibling's first deletedChild is inserted as the last deletedChild into children[idx]
	if !deletedChild.IsLeaf {
		deletedChild.Children[deletedChild.N+1] = rightSibling.Children[0]
	}

	// the first key from rightSibling is inserted into keys[idx]
	parentNode.Keys[idx] = rightSibling.Keys[0]

	// moving all keys in rightSibling one step behind
	// 往前填补空位
	for i := 1; i < rightSibling.N; i++ {
		rightSibling.Keys[i-1] = rightSibling.Keys[i]
	}

	// moving the deletedChild pointers one step behind
	if !rightSibling.IsLeaf {
		for i := 1; i <= rightSibling.N; i++ {
			rightSibling.Children[i-1] = rightSibling.Children[i]
		}
	}

	// increasing and decreasing the key count of children[idx] and children[idx + 1]
	// respectively
	deletedChild.N += 1
	rightSibling.N -= 1
	return
}

// Merge to merge idx-th child of the node with (idx+1)th child of the node
func (btree *BTree) Merge(parentNode *BTreeNode, idx int) {
	deletedChild := parentNode.Children[idx]
	sibling := parentNode.Children[idx+1]

	// pulling a key from the current node and inserting it into (t-1)th
	// position of children[idx]
	deletedChild.Keys[btree.T-1] = parentNode.Keys[idx]

	// copying the keys from children[idx + 1] to children[idx] at the end
	for i := 0; i < sibling.N; i++ {
		deletedChild.Keys[i+btree.T] = sibling.Keys[i]
	}

	// copying the deletedChild pointers from children[idx + 1] to children[idx]
	if !deletedChild.IsLeaf {
		for i := 0; i <= sibling.N; i++ {
			deletedChild.Children[i+btree.T] = sibling.Children[i]
		}
	}

	// moving all keys after idx in the current node one step before
	// to fill the gap created by moving keys[idx] to children[idx]
	// moving the deletedChild pointers after (idx + 1) in the current node one step before
	for i := idx + 1; i < parentNode.N; i++ {
		parentNode.Keys[i-1] = parentNode.Keys[i]
	}

	// moving the deletedChild pointers after (idx+1) in the current node one step before
	for i := idx + 2; i <= parentNode.N; i++ {
		parentNode.Children[i-1] = parentNode.Children[i]
	}

	// updating the key count of deletedChild and the current node
	deletedChild.N += sibling.N + 1
	parentNode.N--
	return
}
