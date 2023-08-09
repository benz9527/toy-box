package test

import (
	"github.com/benz9527/toy-box/lintcode/recursion"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTowerOfHanoi(t *testing.T) {
	results := recursion.TowerOfHanoi(3)
	assert.Equal(t, []string{"from A to C", "from A to B", "from C to B", "from A to C", "from B to A", "from B to C", "from A to C"}, results)
}
