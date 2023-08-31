package api

import (
	isort "sort"
	"testing"
)

type person struct {
	age    int
	skills int
}

type personSlice []person

func (p personSlice) Len() int {
	return len(p)
}
func (p personSlice) Less(i, j int) bool {
	if p[i].age == p[j].age {
		return p[i].skills < p[j].skills
	}
	return p[i].age < p[j].age
}
func (p personSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func TestPersonSort(t *testing.T) {
	persons := []person{
		{
			age:    11,
			skills: 2,
		},
		{
			age:    11,
			skills: 13,
		},
		{
			age:    20,
			skills: 3,
		},
		{
			age:    5,
			skills: 100,
		},
	}
	isort.Sort(personSlice(persons))
	t.Log(persons)
}
