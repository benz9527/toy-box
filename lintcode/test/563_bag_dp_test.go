package test

import (
	"github.com/benz9527/toy-box/lintcode/dp"
	"testing"
)

func TestBackPackV(t *testing.T) {
	type args struct {
		nums   []int
		target int
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
					1, 2, 3, 3, 7,
				},
				target: 7,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dp.BackPackV(tt.args.nums, tt.args.target); got != tt.want {
				t.Errorf("BackPackV() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBackPackV2(t *testing.T) {
	type args struct {
		nums   []int
		target int
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
					1, 2, 3, 3, 7,
				},
				target: 7,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dp.BackPackV2(tt.args.nums, tt.args.target); got != tt.want {
				t.Errorf("BackPackV2() = %v, want %v", got, tt.want)
			}
		})
	}
}
