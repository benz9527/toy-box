package test

import (
	"github.com/benz9527/toy-box/leetcode/dp"
	"testing"
)

func TestFindTargetSumWaysByBackTracing(t *testing.T) {
	type args struct {
		nums   []int
		target int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				nums: []int{
					1, 1, 1, 1, 1,
				},
				target: 3,
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dp.FindTargetSumWaysByBackTracing(tt.args.nums, tt.args.target); got != tt.want {
				t.Errorf("FindTargetSumWaysByBackTracing() = %v, want %v", got, tt.want)
			}
		})
	}
}
