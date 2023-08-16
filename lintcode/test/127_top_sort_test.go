package test

import (
	G "github.com/benz9527/toy-box/lintcode/graph"
	"reflect"
	"testing"
)

func TestTopoSort(t *testing.T) {
	type args struct {
		graph []*G.DirectedGraphNode
	}
	tests := []struct {
		name  string
		args  args
		wants [][]*G.DirectedGraphNode
	}{
		{
			name: "1",
			args: args{
				graph: func() []*G.DirectedGraphNode {
					_0 := G.NewDirectedGraphNode(0)
					_1 := G.NewDirectedGraphNode(1)
					_2 := G.NewDirectedGraphNode(2)
					_3 := G.NewDirectedGraphNode(3)
					_4 := G.NewDirectedGraphNode(4)
					_5 := G.NewDirectedGraphNode(5)
					_0.AddNeighbor(_1, _2, _3)
					_1.AddNeighbor(_4)
					_2.AddNeighbor(_4, _5)
					_3.AddNeighbor(_4, _5)
					return []*G.DirectedGraphNode{_0, _1, _2, _3, _4, _5}
				}(),
			},
			wants: func() [][]*G.DirectedGraphNode {
				_0 := G.NewDirectedGraphNode(0)
				_1 := G.NewDirectedGraphNode(1)
				_2 := G.NewDirectedGraphNode(2)
				_3 := G.NewDirectedGraphNode(3)
				_4 := G.NewDirectedGraphNode(4)
				_5 := G.NewDirectedGraphNode(5)
				return [][]*G.DirectedGraphNode{
					0:  {_0, _1, _2, _3, _4, _5},
					1:  {_0, _1, _3, _2, _4, _5},
					2:  {_0, _1, _3, _2, _5, _4},
					3:  {_0, _1, _2, _3, _5, _4},
					4:  {_0, _2, _1, _3, _4, _5},
					5:  {_0, _2, _1, _3, _5, _4},
					6:  {_0, _2, _3, _1, _5, _4},
					7:  {_0, _2, _3, _1, _4, _5},
					8:  {_0, _3, _1, _2, _4, _5},
					9:  {_0, _3, _1, _2, _5, _4},
					10: {_0, _3, _2, _1, _5, _4},
					11: {_0, _3, _2, _1, _4, _5},
				}
			}(),
		},
		{
			name: "2",
			args: args{
				graph: func() []*G.DirectedGraphNode {
					_1 := G.NewDirectedGraphNode(1)
					_2 := G.NewDirectedGraphNode(2)
					_3 := G.NewDirectedGraphNode(3)
					_4 := G.NewDirectedGraphNode(4)
					_5 := G.NewDirectedGraphNode(5)
					_1.AddNeighbor(_2, _4)
					_2.AddNeighbor(_1, _4)
					_4.AddNeighbor(_1, _2)
					_5.AddNeighbor(_3)
					_3.AddNeighbor(_5)
					return []*G.DirectedGraphNode{_1, _2, _3, _4, _5}
				}(),
			},
			wants: [][]*G.DirectedGraphNode{},
		},
		{
			name: "3",
			args: args{
				graph: func() []*G.DirectedGraphNode {
					_1 := G.NewDirectedGraphNode(1)
					_2 := G.NewDirectedGraphNode(2)
					_3 := G.NewDirectedGraphNode(3)
					_2.AddNeighbor(_3)
					_3.AddNeighbor(_2)
					return []*G.DirectedGraphNode{_1, _2, _3}
				}(),
			},
			wants: [][]*G.DirectedGraphNode{
				0: {G.NewDirectedGraphNode(1)},
			},
		},
		{
			name: "4",
			args: args{
				graph: func() []*G.DirectedGraphNode {
					_1 := G.NewDirectedGraphNode(1)
					_2 := G.NewDirectedGraphNode(2)
					_3 := G.NewDirectedGraphNode(3)
					_2.AddNeighbor(_3)
					return []*G.DirectedGraphNode{_1, _2, _3}
				}(),
			},
			wants: [][]*G.DirectedGraphNode{
				0: {
					G.NewDirectedGraphNode(1),
					G.NewDirectedGraphNode(2),
					G.NewDirectedGraphNode(3),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := G.TopoSort(tt.args.graph)
			if !topoSortResultCompare(t, res, tt.wants...) {
				t.Errorf("TopoSort() = %v, wants %v", res, tt.wants)
			}
		})
	}
}

func topoSortResultCompare(t *testing.T, result []*G.DirectedGraphNode, wants ...[]*G.DirectedGraphNode) bool {
	convert := func(srcs ...[]*G.DirectedGraphNode) [][]int {
		results := make([][]int, 0, 16)
		for _, src := range srcs {
			res := make([]int, 0, len(result))
			for _, s := range src {
				res = append(res, s.Label)
			}
			results = append(results, res)
		}
		return results
	}

	if len(result) <= 0 && len(wants) <= 0 {
		t.Log([]int{})
		return true
	} else if len(result) > 0 && len(wants) <= 0 ||
		len(result) <= 0 && len(wants) > 0 {
		t.Log(convert(result)[0])
		return false
	}

	src := convert(result)[0]
	t.Log(src)
	results := convert(wants...)
	for _, res := range results {
		if reflect.DeepEqual(src, res) {
			return true
		}
	}
	return false
}

func TestTopoSort2(t *testing.T) {
	type args struct {
		graph []*G.DirectedGraphNode
	}
	tests := []struct {
		name  string
		args  args
		wants [][]*G.DirectedGraphNode
	}{
		{
			name: "1",
			args: args{
				graph: func() []*G.DirectedGraphNode {
					_0 := G.NewDirectedGraphNode(0)
					_1 := G.NewDirectedGraphNode(1)
					_2 := G.NewDirectedGraphNode(2)
					_3 := G.NewDirectedGraphNode(3)
					_4 := G.NewDirectedGraphNode(4)
					_5 := G.NewDirectedGraphNode(5)
					_0.AddNeighbor(_1, _2, _3)
					_1.AddNeighbor(_4)
					_2.AddNeighbor(_4, _5)
					_3.AddNeighbor(_4, _5)
					return []*G.DirectedGraphNode{_0, _1, _2, _3, _4, _5}
				}(),
			},
			wants: func() [][]*G.DirectedGraphNode {
				_0 := G.NewDirectedGraphNode(0)
				_1 := G.NewDirectedGraphNode(1)
				_2 := G.NewDirectedGraphNode(2)
				_3 := G.NewDirectedGraphNode(3)
				_4 := G.NewDirectedGraphNode(4)
				_5 := G.NewDirectedGraphNode(5)
				return [][]*G.DirectedGraphNode{
					0:  {_0, _1, _2, _3, _4, _5},
					1:  {_0, _1, _3, _2, _4, _5},
					2:  {_0, _1, _3, _2, _5, _4},
					3:  {_0, _1, _2, _3, _5, _4},
					4:  {_0, _2, _1, _3, _4, _5},
					5:  {_0, _2, _1, _3, _5, _4},
					6:  {_0, _2, _3, _1, _5, _4},
					7:  {_0, _2, _3, _1, _4, _5},
					8:  {_0, _3, _1, _2, _4, _5},
					9:  {_0, _3, _1, _2, _5, _4},
					10: {_0, _3, _2, _1, _5, _4},
					11: {_0, _3, _2, _1, _4, _5},
				}
			}(),
		},
		{
			name: "2",
			args: args{
				graph: func() []*G.DirectedGraphNode {
					_1 := G.NewDirectedGraphNode(1)
					_2 := G.NewDirectedGraphNode(2)
					_3 := G.NewDirectedGraphNode(3)
					_4 := G.NewDirectedGraphNode(4)
					_5 := G.NewDirectedGraphNode(5)
					_1.AddNeighbor(_2, _4)
					_2.AddNeighbor(_1, _4)
					_4.AddNeighbor(_1, _2)
					_5.AddNeighbor(_3)
					_3.AddNeighbor(_5)
					return []*G.DirectedGraphNode{_1, _2, _3, _4, _5}
				}(),
			},
			wants: [][]*G.DirectedGraphNode{},
		},
		{
			name: "3",
			args: args{
				graph: func() []*G.DirectedGraphNode {
					_1 := G.NewDirectedGraphNode(1)
					_2 := G.NewDirectedGraphNode(2)
					_3 := G.NewDirectedGraphNode(3)
					_2.AddNeighbor(_3)
					_3.AddNeighbor(_2)
					return []*G.DirectedGraphNode{_1, _2, _3}
				}(),
			},
			wants: [][]*G.DirectedGraphNode{
				0: {G.NewDirectedGraphNode(1)},
			},
		},
		{
			name: "4",
			args: args{
				graph: func() []*G.DirectedGraphNode {
					_1 := G.NewDirectedGraphNode(1)
					_2 := G.NewDirectedGraphNode(2)
					_3 := G.NewDirectedGraphNode(3)
					_2.AddNeighbor(_3)
					return []*G.DirectedGraphNode{_1, _2, _3}
				}(),
			},
			wants: [][]*G.DirectedGraphNode{
				0: {
					G.NewDirectedGraphNode(1),
					G.NewDirectedGraphNode(2),
					G.NewDirectedGraphNode(3),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := G.TopoSort2(tt.args.graph)
			if !topoSortResultCompare(t, res, tt.wants...) {
				t.Errorf("TopoSort() = %v, wants %v", res, tt.wants)
			}
		})
	}
}
