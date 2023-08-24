package test

import (
	"github.com/benz9527/toy-box/leetcode/dp"
	"testing"
)

func TestLongestCommonSubsequence(t *testing.T) {
	type args struct {
		text1 string
		text2 string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				text1: "abcde",
				text2: "ace",
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dp.LongestCommonSubsequence(tt.args.text1, tt.args.text2); got != tt.want {
				t.Errorf("LongestCommonSubsequence() = %v, want %v", got, tt.want)
			}
		})
	}
}
