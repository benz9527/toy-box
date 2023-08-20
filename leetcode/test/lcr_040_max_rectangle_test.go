package test

import (
	"github.com/benz9527/toy-box/leetcode/arr"
	"testing"
)

func TestMaximalRectangle(t *testing.T) {
	type args struct {
		matrix []string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				matrix: []string{
					"10100", "10111", "11111", "10010",
				},
			},
			want: 6,
		},
		{
			name: "2",
			args: args{
				matrix: []string{},
			},
			want: 0,
		},
		{
			name: "3",
			args: args{
				matrix: []string{
					"1",
				},
			},
			want: 1,
		},
		{
			name: "4",
			args: args{
				matrix: []string{
					"00",
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := arr.MaximalRectangle(tt.args.matrix); got != tt.want {
				t.Errorf("MaximalRectangle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxRecArea(t *testing.T) {
	type args struct {
		heights []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				heights: []int{
					3, 4, 5, 6, 4, 0,
				},
			},
			want: 16,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := arr.MaxRecArea(tt.args.heights); got != tt.want {
				t.Errorf("MaxRecArea() = %v, want %v", got, tt.want)
			}
		})
	}
}
