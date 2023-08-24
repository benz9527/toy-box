package test

import (
	"github.com/benz9527/toy-box/leetcode/dp"
	"testing"
)

func TestMaxLengthOfRepeatedSubarray(t *testing.T) {
	type args struct {
		nums1 []int
		nums2 []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				nums1: []int{
					1, 2, 3, 2, 1,
				},
				nums2: []int{
					3, 2, 1, 4, 7,
				},
			},
			want: 3,
		},
		{
			name: "1",
			args: args{
				nums1: []int{
					1, 0, 0, 0, 0,
				},
				nums2: []int{
					0, 0, 0, 0, 1,
				},
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dp.MaxLengthOfRepeatedSubarray(tt.args.nums1, tt.args.nums2); got != tt.want {
				t.Errorf("MaxLengthOfRepeatedSubarray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxLengthOfRepeatedSubarrayOptimize_Compare(t *testing.T) {
	type args struct {
		nums1 []int
		nums2 []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				nums1: []int{
					1, 0, 0, 0, 0,
				},
				nums2: []int{
					0, 0, 0, 0, 1,
				},
			},
			want: 4,
		},
		{
			name: "1",
			args: args{
				nums1: []int{
					0, 0, 0, 0, 1,
				},
				nums2: []int{
					1, 0, 0, 0, 0,
				},
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		g1 := dp.MaxLengthOfRepeatedSubarrayOptimize(tt.args.nums1, tt.args.nums2)
		g2 := dp.MaxLengthOfRepeatedSubarrayOptimize2(tt.args.nums1, tt.args.nums2)
		t.Run(tt.name, func(t *testing.T) {
			if g1 != g2 && g1 != tt.want || g2 != tt.want {
				t.Errorf("MaxLengthOfRepeatedSubarrayOptimize() = %v, MaxLengthOfRepeatedSubarrayOptimize2() = %v, want %v", g1, g2, tt.want)
			}
		})
	}
}
