package tree

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTrie(t *testing.T) {
	root := Constructor()
	root.Insert("hello")
	ans := root.Search("hell")
	assert.False(t, ans)

	root.Insert("apple")
	ans = root.Search("appl")
	assert.False(t, ans)
}
