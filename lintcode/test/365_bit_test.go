package test

import (
	"github.com/benz9527/toy-box/lintcode/bit"
	"testing"
)

func TestCountOnes2(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				n: 7,
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bit.CountOnes2(tt.args.n); got != tt.want {
				t.Errorf("CountOnes2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCountOnes3(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				n: 7,
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bit.CountOnes3(tt.args.n); got != tt.want {
				t.Errorf("CountOnes2() = %v, want %v", got, tt.want)
			}
		})
	}
}
