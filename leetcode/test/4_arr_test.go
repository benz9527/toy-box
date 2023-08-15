package test

import (
	"github.com/benz9527/toy-box/leetcode/arr"
	"testing"
)

func TestFindMedianSortedArrays(t *testing.T) {
	type args struct {
		nums1 []int
		nums2 []int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "1",
			args: args{
				nums1: []int{1},
				nums2: []int{2, 3},
			},
			want: 2.0,
		},
		{
			name: "2",
			args: args{
				nums1: []int{1, 2},
				nums2: []int{2, 3},
			},
			want: 2.0,
		},
		{
			name: "3",
			args: args{
				nums1: []int{1},
				nums2: []int{2},
			},
			want: 1.5,
		},
		{
			name: "4",
			args: args{
				nums1: []int{1},
				nums2: []int{},
			},
			want: 1.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := arr.FindMedianSortedArrays(tt.args.nums1, tt.args.nums2); got != tt.want {
				t.Errorf("findMedianSortedArrays() = %v, want %v", got, tt.want)
			}
		})
	}
}
