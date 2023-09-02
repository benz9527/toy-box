package b

import (
	"reflect"
	"testing"
)

func TestKSum(t *testing.T) {
	type args struct {
		nums   []int64
		k      int
		target int64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				nums: []int64{
					2, 7, 11, 15,
				},
				k:      2,
				target: 9,
			},
			want: 1,
		},
		{
			name: "2",
			args: args{
				nums: []int64{
					-1, 0, 1, 2, -1, -4,
				},
				k:      3,
				target: 0,
			},
			want: 2,
		},
		{
			name: "3",
			args: args{
				nums: []int64{
					-1, 0, 1, 2, -1, -4,
				},
				k:      4,
				target: 1,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KSum(tt.args.nums, tt.args.k, tt.args.target); got != tt.want {
				t.Errorf("KSum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJumpHousesII(t *testing.T) {
	type args struct {
		steps []int
		count int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "1",
			args: args{
				steps: []int{1, 4, 5, 2, 0, 2},
				count: 9,
			},
			want: []int{4, 5, 0},
		},
		{
			name: "2",
			args: args{
				steps: []int{1, 5, 2, 0, 2, 4},
				count: 9,
			},
			want: []int{5, 2, 2},
		},
		{
			name: "2",
			args: args{
				steps: []int{-1, 2, 4, 9},
				count: 12,
			},
			want: []int{-1, 4, 9},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JumpHousesII(tt.args.steps, tt.args.count); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JumpHousesII() = %v, want %v", got, tt.want)
			}
		})
	}
}
