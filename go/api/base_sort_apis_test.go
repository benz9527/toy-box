package api

import (
	"github.com/stretchr/testify/assert"
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

type participant struct {
	height, weight int
}

type participantSlice []participant

func (s participantSlice) Less(i, j int) bool {
	if s[i].height == s[j].height {
		return s[i].weight > s[j].weight
	}
	return s[i].height > s[j].height
}
func (s participantSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s participantSlice) Len() int {
	return len(s)
}

func TestParticipants(t *testing.T) {
	participants := []participant{
		{
			height: 181,
			weight: 70,
		},
		{
			height: 182,
			weight: 70,
		},
		{
			height: 183,
			weight: 70,
		},
		{
			height: 184,
			weight: 70,
		},
		{
			height: 185,
			weight: 70,
		},
		{
			height: 186,
			weight: 70,
		},
		{
			height: 180,
			weight: 71,
		},
		{
			height: 180,
			weight: 72,
		},
		{
			height: 180,
			weight: 73,
		},
		{
			height: 180,
			weight: 74,
		},
		{
			height: 180,
			weight: 75,
		},
	}
	expected := []participant{
		{
			height: 186,
			weight: 70,
		},
		{
			height: 185,
			weight: 70,
		},
		{
			height: 184,
			weight: 70,
		},
		{
			height: 183,
			weight: 70,
		},
		{
			height: 182,
			weight: 70,
		},
		{
			height: 181,
			weight: 70,
		},
		{
			height: 180,
			weight: 75,
		},
		{
			height: 180,
			weight: 74,
		},
		{
			height: 180,
			weight: 73,
		},
		{
			height: 180,
			weight: 72,
		},
		{
			height: 180,
			weight: 71,
		},
	}
	isort.Sort(participantSlice(participants))
	assert.Equal(t, expected, participants)
	t.Log(participants[:10])
}
