package test

import (
	"github.com/benz9527/toy-box/leetcode/dp"
	"testing"
)

func TestMinCostClimbingStairs(t *testing.T) {
	type args struct {
		cost []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				cost: []int{
					10, 15, 20,
				},
			},
			want: 15,
		},
		{
			name: "2",
			args: args{
				cost: []int{
					1, 100, 1, 1, 1, 100, 1, 1, 100, 1,
				},
			},
			want: 6,
		},
		{
			name: "3",
			args: args{
				cost: []int{
					0, 0, 0, 1,
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g1 := dp.MinCostClimbingStairs(tt.args.cost)
			g2 := dp.MinCostClimbingStairsOptimize(tt.args.cost)
			if g1 != g2 && g1 != tt.want || g2 != tt.want {
				t.Errorf("MinCostClimbingStairs() = %v, MinCostClimbingStairsOptimize() = %v, want %v", g1, g2, tt.want)
			}
		})
	}
}
