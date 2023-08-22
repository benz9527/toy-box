package test

import (
	"github.com/benz9527/toy-box/leetcode/dp"
	"testing"
)

func TestMaxProfitIV(t *testing.T) {
	type args struct {
		k      int
		prices []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				k: 4,
				prices: []int{
					1, 2, 4, 2, 5, 7, 2, 4, 9, 0,
				},
			},
			want: 15,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dp.MaxProfitIV(tt.args.k, tt.args.prices); got != tt.want {
				t.Errorf("MaxProfitIV() = %v, want %v", got, tt.want)
			}
		})
	}
}
