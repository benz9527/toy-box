package tree

type Trie struct {
	nodes   [26]*Trie
	isEnded bool
}

func Constructor() Trie {
	return Trie{nodes: [26]*Trie{}, isEnded: false}
}

func (this *Trie) Insert(word string) {
	cur := this
	for i := 0; i < len(word); i++ {
		idx := word[i] - 'a'
		if cur.nodes[idx] == nil {
			node := Constructor()
			cur.nodes[idx] = &node
		}
		cur = cur.nodes[idx]
	}
	cur.isEnded = true
}

func (this *Trie) searchPrefix(word string) *Trie {
	cur := this
	for i := 0; i < len(word); i++ {
		idx := word[i] - 'a'
		if cur.nodes[idx] == nil {
			return nil
		}
		cur = cur.nodes[idx]
	}
	return cur
}

func (this *Trie) Search(word string) bool {
	n := this.searchPrefix(word)
	return n != nil && n.isEnded
}

func (this *Trie) StartsWith(prefix string) bool {
	n := this.searchPrefix(prefix)
	return n != nil
}
