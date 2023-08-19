package test

import (
	"github.com/benz9527/toy-box/leetcode/dp"
	"testing"
)

func TestCanPartition(t *testing.T) {
	type args struct {
		nums []int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1",
			args: args{
				nums: []int{
					1, 2, 3, 3, 5,
				},
			},
			want: true,
		},
		{
			name: "2",
			args: args{
				nums: []int{
					1, 5, 11, 5,
				},
			},
			want: true,
		},
		{
			name: "3",
			args: args{
				nums: []int{
					5, 5, 5, 5,
				},
			},
			want: true,
		},
		{
			name: "4",
			args: args{
				nums: []int{
					66, 90, 7, 6, 32, 16, 2, 78, 69, 88, 85, 26, 3, 9, 58, 65, 30, 96, 11,
					31, 99, 49, 63, 83, 79, 97, 20, 64, 81,
					80, 25, 69, 9, 75, 23, 70, 26, 71, 25, 54, 1, 40, 41, 82, 32, 10, 26,
					33, 50, 71, 5, 91, 59, 96, 9, 15, 46, 70,
					26, 32, 49, 35, 80, 21, 34, 95, 51, 66, 17, 71, 28, 88, 46, 21, 31,
					71, 42, 2, 98, 96, 40, 65, 92, 43, 68, 14,
					98, 38, 13, 77, 14, 13, 60, 79, 52, 46, 9, 13, 25, 8,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g1 := dp.CanPartition(tt.args.nums)
			g2 := dp.CanPartitionOptimize(tt.args.nums)
			if g1 != g2 && g1 != tt.want || g2 != tt.want {
				t.Errorf("CanPartition() = %v, CanPartitionOptimize() = %v, want %v", g1, g2, tt.want)
			}
		})
	}
}
