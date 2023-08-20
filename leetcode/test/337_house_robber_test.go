package test

import (
	"github.com/benz9527/toy-box/leetcode/dp"
	"testing"
)

func TestRobIII(t *testing.T) {
	type args struct {
		root *dp.TreeNode
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				root: &dp.TreeNode{
					Val: 3,
					Left: &dp.TreeNode{
						Val: 2,
						Right: &dp.TreeNode{
							Val: 3,
						},
					},
					Right: &dp.TreeNode{
						Val: 3,
						Right: &dp.TreeNode{
							Val: 1,
						},
					},
				},
			},
			want: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dp.RobIII(tt.args.root); got != tt.want {
				t.Errorf("RobIII() = %v, want %v", got, tt.want)
			}
		})
	}
}
