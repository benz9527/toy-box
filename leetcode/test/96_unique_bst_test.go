package test

import (
	"github.com/benz9527/toy-box/leetcode/dp"
	"testing"
)

func TestNumTrees(t *testing.T) {
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
				n: 3,
			},
			want: 5,
		},
		{
			name: "1",
			args: args{
				n: 5,
			},
			want: 42,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dp.NumTrees(tt.args.n); got != tt.want {
				t.Errorf("NumTrees() = %v, want %v", got, tt.want)
			}
		})
	}
}
