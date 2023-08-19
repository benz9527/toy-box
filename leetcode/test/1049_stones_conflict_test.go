package test

import (
	"github.com/benz9527/toy-box/leetcode/dp"
	"testing"
)

func TestLastStoneWeightII(t *testing.T) {
	type args struct {
		stones []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				stones: []int{
					31, 26, 33, 21, 40, 11,
				},
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dp.LastStoneWeightII(tt.args.stones); got != tt.want {
				t.Errorf("LastStoneWeightII() = %v, want %v", got, tt.want)
			}
		})
	}
}
