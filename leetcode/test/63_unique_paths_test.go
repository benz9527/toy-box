package test

import (
	"github.com/benz9527/toy-box/leetcode/dp"
	"testing"
)

func TestUniquePathsWithObstacles(t *testing.T) {
	type args struct {
		obstacleGrid [][]int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				obstacleGrid: [][]int{
					{0, 0}, {1, 1}, {0, 0},
				},
			},
			want: 0,
		},
		{
			name: "2",
			args: args{
				obstacleGrid: [][]int{
					{0, 1, 0, 0, 0}, {1, 0, 0, 0, 0}, {0, 0, 0, 0, 0}, {0, 0, 0, 0, 0},
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dp.UniquePathsWithObstacles(tt.args.obstacleGrid); got != tt.want {
				t.Errorf("UniquePathsWithObstacles() = %v, want %v", got, tt.want)
			}
		})
	}
}
