package test

import (
	"github.com/benz9527/toy-box/leetcode/dp"
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
				coins: []int{
					2, 5, 10, 1,
				},
				amount: 27,
			},
			want: 4,
		},
		{
			name: "2",
			args: args{
				coins: []int{
					186, 419, 83, 408,
				},
				amount: 6249,
			},
			want: 20,
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
