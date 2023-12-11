package list

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func strCompare(a, b string) int {
	if a == b {
		return 0
	}
	if a < b {
		return -1
	}
	return 1
}

func TestClassicSkipList(t *testing.T) {
	words := []string{
		"foo", "bar", "zap", "pomo", "pera", "arancio", "limone",
	}
	sl := NewClassicSkipList[string](strCompare)
	for _, word := range words {
		sl.Insert(word)
	}
	expected := []string{
		"arancio", "bar", "foo", "limone", "pera", "pomo", "zap",
	}
	actual := make([]string, 0, len(words))
	sl.ForEach(func(idx int64, v string) {
		actual = append(actual, v)
		t.Logf("idx: %d, v: %s", idx, v)
	})
	assert.Equal(t, expected, actual)

	e := sl.Find("ben")
	assert.Nil(t, e)

	e = sl.Find("limone")
	assert.Equal(t, "limone", e.GetObject())

	sl.PopHead()
	expected = []string{
		"bar", "foo", "limone", "pera", "pomo", "zap",
	}
	actual = make([]string, 0, len(words))
	sl.ForEach(func(idx int64, v string) {
		actual = append(actual, v)
		t.Logf("idx: %d, v: %s", idx, v)
	})
	assert.Equal(t, expected, actual)

	sl.Remove("pera")
	expected = []string{
		"bar", "foo", "limone", "pomo", "zap",
	}
	actual = make([]string, 0, len(words))
	sl.ForEach(func(idx int64, v string) {
		actual = append(actual, v)
		t.Logf("idx: %d, v: %s", idx, v)
	})
	assert.Equal(t, expected, actual)
}

func TestMaxLevel(t *testing.T) {
	maxLevels := MaxLevels(math.MaxInt32, 0.25)
	assert.GreaterOrEqual(t, 32, maxLevels)
}
