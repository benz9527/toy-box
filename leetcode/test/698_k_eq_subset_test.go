package test

import (
	"github.com/benz9527/toy-box/leetcode/backstrace"
	"testing"
)

func TestCanPartitionKSubsets(t *testing.T) {
	type args struct {
		nums []int
		k    int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1",
			args: args{
				nums: []int{4, 3, 2, 3, 5, 2, 1},
				k:    4,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := backstrace.CanPartitionKSubsets(tt.args.nums, tt.args.k); got != tt.want {
				t.Errorf("CanPartitionKSubsets() = %v, want %v", got, tt.want)
			}
		})
	}
}
