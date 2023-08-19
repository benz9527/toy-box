package test

import (
	"github.com/benz9527/toy-box/leetcode/dp"
	"testing"
)

func TestIntegerBreak(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				n: 4,
			},
			want: 4,
		},
		{
			name: "2",
			args: args{
				n: 10,
			},
			want: 36,
		},
		{
			name: "3",
			args: args{
				n: 28,
			},
			want: 26244,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g1 := dp.IntegerBreak(tt.args.n)
			g2 := dp.IntegerBreakOptimize(tt.args.n)
			if g1 != g2 && g1 != tt.want || g2 != tt.want {
				t.Errorf("IntegerBreak() = %v, IntegerBreakOptimize() = %v, want %v", g1, g2, tt.want)
			}
		})
	}
}
