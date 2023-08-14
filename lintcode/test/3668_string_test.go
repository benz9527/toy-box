package test

import (
	"github.com/benz9527/toy-box/lintcode/str"
	"reflect"
	"testing"
)

func TestBracketExpansion(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "1",
			args: args{
				s: "a{b,c}",
			},
			want: []string{
				"ab", "ac",
			},
		},
		{

			name: "2",
			args: args{
				s: "abcde",
			},
			want: []string{
				"abcde",
			},
		},
		{

			name: "3",
			args: args{
				s: "ab{c,d}e",
			},
			want: []string{
				"abce", "abde",
			},
		},
		{

			name: "4",
			args: args{
				s: "{a,b}{c,d}e",
			},
			want: []string{
				"ace", "ade", "bce", "bde",
			},
		},
		{

			name: "5",
			args: args{
				s: "{a,b}{c,d}{e,f}",
			},
			want: []string{
				"ace", "acf", "ade", "adf", "bce", "bcf", "bde", "bdf",
			},
		},
		{

			name: "6",
			args: args{
				s: "{b,a}{c,d}{e,f}",
			},
			want: []string{
				"ace", "acf", "ade", "adf", "bce", "bcf", "bde", "bdf",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := str.BracketExpansion(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BracketExpansion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBracketExpansion2(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "1",
			args: args{
				s: "a{b,c}",
			},
			want: []string{
				"ab", "ac",
			},
		},
		{

			name: "2",
			args: args{
				s: "abcde",
			},
			want: []string{
				"abcde",
			},
		},
		{

			name: "3",
			args: args{
				s: "ab{c,d}e",
			},
			want: []string{
				"abce", "abde",
			},
		},
		{

			name: "4",
			args: args{
				s: "{a,b}{c,d}e",
			},
			want: []string{
				"ace", "ade", "bce", "bde",
			},
		},
		{

			name: "5",
			args: args{
				s: "{a,b}{c,d}{e,f}",
			},
			want: []string{
				"ace", "acf", "ade", "adf", "bce", "bcf", "bde", "bdf",
			},
		},
		{

			name: "6",
			args: args{
				s: "{b,a}{c,d}{e,f}",
			},
			want: []string{
				"ace", "acf", "ade", "adf", "bce", "bcf", "bde", "bdf",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := str.BracketExpansion2(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BracketExpansion() = %v, want %v", got, tt.want)
			}
		})
	}
}
