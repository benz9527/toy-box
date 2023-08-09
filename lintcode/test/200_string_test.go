package test

import (
	"github.com/benz9527/toy-box/lintcode/str"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLongestPalindrome(t *testing.T) {
	res := str.LongestPalindrome("abcdzdcab")
	assert.Equal(t, "cdzdc", res)

}
