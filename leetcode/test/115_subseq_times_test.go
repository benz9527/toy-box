package test

import (
	"github.com/benz9527/toy-box/leetcode/dp"
	"testing"
)

func TestNumDistinct(t *testing.T) {
	type args struct {
		s string
		t string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				s: "abcbc",
				t: "ab",
			},
			want: 2,
		},
		{
			name: "1",
			args: args{
				s: "abbcbc",
				t: "ab",
			},
			want: 3,
		},
		{
			name: "1",
			args: args{
				s: "rabbbit",
				t: "rabbit",
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dp.NumDistinct(tt.args.s, tt.args.t); got != tt.want {
				t.Errorf("NumDistinct() = %v, want %v", got, tt.want)
			}
		})
	}
}
