package test

import (
	"github.com/benz9527/toy-box/lintcode/dp"
	"testing"
)

func TestValley(t *testing.T) {
	type args struct {
		nums []int
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
					5, 4, 3, 2, 1, 2, 3, 4, 5,
				},
			},
			want: 8,
		},
		{
			name: "2",
			args: args{
				nums: []int{
					1, 2, 3, 4, 5,
				},
			},
			want: 0,
		},
		{
			name: "3",
			args: args{
				nums: []int{
					9, 8, 10, 7, 6, 5, 1, 2, 3, 4, 1, 9,
				},
			},
			want: 4,
		},
		{
			name: "4",
			args: args{
				nums: []int{5},
			},
			want: 0,
		},
		{
			name: "5",
			args: args{
				nums: []int{
					0, 0, 4, 4, 6, 0, 2, 5, 0, 0, 4, 10, 7, 1, 6, 5, 6, 0, 2, 8, 2, 7, 10, 0, 4, 7, 6, 9, 1, 9, 4, 5, 0, 6, 5, 1, 5, 3, 3, 6, 5, 10, 1, 1, 10, 0, 9, 6, 6, 4,
				},
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dp.Valley(tt.args.nums); got != tt.want {
				t.Errorf("Valley() = %v, want %v", got, tt.want)
			}
		})
	}
}
