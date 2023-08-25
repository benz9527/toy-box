package test

import (
	string2 "github.com/benz9527/toy-box/leetcode/string"
	Lo "github.com/samber/lo"
	"testing"
)

func TestLongestPalindrome(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name  string
		args  args
		wants []string
	}{
		{
			name: "",
			args: args{
				s: "babad",
			},
			wants: []string{"bab", "aba"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string2.LongestPalindrome(tt.args.s)
			if !Lo.Contains[string](tt.wants, got) {
				t.Errorf("LongestPalindrome() = %v, want %v", got, tt.wants)
			}
		})
	}
}
