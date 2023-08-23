package test

import (
	"github.com/benz9527/toy-box/leetcode/dp"
	"testing"
)

func TestMaxProfitV(t *testing.T) {
	type args struct {
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
				prices: []int{
					1, 2, 3, 0, 2,
				},
			},
			want: 3,
		},
		{
			name: "2",
			args: args{
				prices: []int{
					1, 2, 3, 0, 2, 9, 6, 100, 5, 2, 0, 999,
				},
			},
			want: 1100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g1 := dp.MaxProfitV(tt.args.prices)
			g2 := dp.MaxProfitVOptimize(tt.args.prices)
			if g1 != g2 && g1 != tt.want || g2 != tt.want {
				t.Errorf("MaxProfitV() = %v, MaxProfitVOptimize() = %v, want %v", g1, g2, tt.want)
			}
		})
	}
}
