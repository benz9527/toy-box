package test

import (
	"github.com/benz9527/toy-box/lintcode/tree"
	"testing"
)

func TestTakeCoins(t *testing.T) {
	type args struct {
		list []int
		k    int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				list: []int{5, 4, 3, 2, 1},
				k:    2,
			},
			want: 9,
		},
		{
			name: "2",
			args: args{
				list: []int{5, 4, 3, 2, 1, 6},
				k:    3,
			},
			want: 15,
		},
		{
			name: "3",
			args: args{
				list: []int{5, 4, 100, 3, 108, 2, 1},
				k:    3,
			},
			want: 111,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tree.TakeCoins(tt.args.list, tt.args.k); got != tt.want {
				t.Errorf("TakeCoins() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTakeCoins2(t *testing.T) {
	type args struct {
		list []int
		k    int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				list: []int{5, 4, 3, 2, 1},
				k:    2,
			},
			want: 9,
		},
		{
			name: "2",
			args: args{
				list: []int{5, 4, 3, 2, 1, 6},
				k:    3,
			},
			want: 15,
		},
		{
			name: "3",
			args: args{
				list: []int{5, 4, 100, 3, 108, 2, 1},
				k:    3,
			},
			want: 111,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tree.TakeCoins2(tt.args.list, tt.args.k); got != tt.want {
				t.Errorf("TakeCoins2() = %v, want %v", got, tt.want)
			}
		})
	}
}
