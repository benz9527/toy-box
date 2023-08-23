package test

import (
	"github.com/benz9527/toy-box/leetcode/dp"
	"testing"
)

func TestLengthOfLIS(t *testing.T) {
	type args struct {
		nums []int
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
					1, 3, 5, 4, 7,
				},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dp.LengthOfLIS(tt.args.nums); got != tt.want {
				t.Errorf("LengthOfLIS() = %v, want %v", got, tt.want)
			}
		})
	}
}
