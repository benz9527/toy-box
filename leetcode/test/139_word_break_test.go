package test

import (
	"github.com/benz9527/toy-box/leetcode/dp"
	"testing"
)

func TestWordBreak(t *testing.T) {
	type args struct {
		s        string
		wordDict []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1",
			args: args{
				s: "leetcode",
				wordDict: []string{
					"leet", "code",
				},
			},
			want: true,
		},
		{
			name: "2",
			args: args{
				s: "leetcode",
				wordDict: []string{
					"t", "code",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dp.WordBreak(tt.args.s, tt.args.wordDict); got != tt.want {
				t.Errorf("WordBreak() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWordBreakOptimize(t *testing.T) {
	type args struct {
		s        string
		wordDict []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1",
			args: args{
				s: "leetcode",
				wordDict: []string{
					"leet", "code",
				},
			},
			want: true,
		},
		{
			name: "2",
			args: args{
				s: "leetcode",
				wordDict: []string{
					"t", "code",
				},
			},
			want: false,
		},
		{
			name: "3",
			args: args{
				s: "aaaaaaa",
				wordDict: []string{
					"aaaa", "aaa",
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dp.WordBreakOptimize(tt.args.s, tt.args.wordDict); got != tt.want {
				t.Errorf("WordBreakOptimize() = %v, want %v", got, tt.want)
			}
		})
	}
}
