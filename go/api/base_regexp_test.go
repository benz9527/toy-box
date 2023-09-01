package api

import (
	"regexp"
	"testing"
)

func TestRegexpReturnMatchIndex(t *testing.T) {
	c, _ := regexp.Compile("b[cd]")
	idxList := c.FindStringSubmatchIndex("abcd")
	t.Log(idxList)
	idxList = c.FindStringIndex("abcd")
	t.Log(idxList)

	c, _ = regexp.Compile("[ae]b[cd]")
	t.Log(c.FindStringSubmatchIndex("aebcd"))
	t.Log(c.FindAllStringIndex("aebcdabc", -1))
	t.Log(c.FindAllStringSubmatch("aebcdabc", -1))
}
