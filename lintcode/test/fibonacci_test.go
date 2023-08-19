package test

import (
	"github.com/benz9527/toy-box/lintcode/dp"
	"testing"
)

func TestFibonacci(t *testing.T) {
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
				n: 1,
			},
			want: 1,
		},
		{
			name: "2",
			args: args{
				n: 0,
			},
			want: 0,
		},
		{
			name: "3",
			args: args{
				n: 10,
			},
			want: 55,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g1, g2 := dp.Fibonacci(tt.args.n), dp.FibonacciOptimize(tt.args.n)
			if g1 != g2 && g1 != tt.want || g2 != tt.want {
				t.Errorf("Fibonacci() = %v, FibonacciOptimize() = %v, want %v", g1, g2, tt.want)
			}
		})
	}
}
