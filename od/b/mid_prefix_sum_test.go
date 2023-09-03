package b

import (
	"reflect"
	"testing"
)

func TestJudgeContinuousSumIsMultipleOfK(t *testing.T) {
	type args struct {
		nums []int
		m    int
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
					2, 12, 6, 3, 5, 5,
				},
				m: 7,
			},
			want: 1,
		},
		{
			name: "2",
			args: args{
				nums: []int{
					1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
				},
				m: 11,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JudgeContinuousSumIsMultipleOfK(tt.args.nums, tt.args.m); got != tt.want {
				t.Errorf("JudgeContinuousSumIsMultipleOfK() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSameWeightSubContinuousStr(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want [2]int
	}{
		{
			name: "1",
			args: args{
				str: "acdbbbca",
			},
			want: [2]int{2, 5},
		},
		{
			name: "2",
			args: args{
				str: "abcabc",
			},
			want: [2]int{0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SameWeightSubContinuousStr(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SameWeightSubContinuousStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
