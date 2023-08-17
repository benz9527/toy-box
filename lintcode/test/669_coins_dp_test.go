package test

import (
	"github.com/benz9527/toy-box/lintcode/dp"
	"testing"
)

func TestCoinChange(t *testing.T) {
	type args struct {
		coins  []int
		amount int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				coins:  []int{1, 2, 5},
				amount: 11,
			},
			want: 3,
		},
		{
			name: "2",
			args: args{
				coins:  []int{1, 2, 3, 100, 5, 88},
				amount: 33,
			},
			want: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dp.CoinChange(tt.args.coins, tt.args.amount); got != tt.want {
				t.Errorf("CoinChange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoinChange2(t *testing.T) {
	type args struct {
		coins  []int
		amount int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				coins:  []int{1, 2, 5},
				amount: 11,
			},
			want: 3,
		},
		{
			name: "2",
			args: args{
				coins:  []int{1, 2, 3, 100, 5, 88},
				amount: 33,
			},
			want: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dp.CoinChange2(tt.args.coins, tt.args.amount); got != tt.want {
				t.Errorf("CoinChange2() = %v, want %v", got, tt.want)
			}
		})
	}
}
