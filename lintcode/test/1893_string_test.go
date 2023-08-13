package test

import (
	"github.com/benz9527/toy-box/lintcode/str"
	"testing"
)

func TestIsValid(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				s: "aabbccc",
			},
			want: "YES",
		},
		{
			name: "2",
			args: args{
				s: "aabbcd",
			},
			want: "NO",
		},
		{
			name: "3",
			args: args{
				s: "zzzzzza",
			},
			want: "YES",
		},
		{
			name: "4",
			args: args{
				s: "zzzzzzzzzzzzzzzzzzzzz",
			},
			want: "YES",
		},
		{
			name: "5",
			args: args{
				s: "zzzzzzzzzzzzzzzzzzzzzahfdl",
			},
			want: "NO",
		},
		{
			name: "6",
			args: args{
				s: "abbcc",
			},
			want: "YES",
		},
		{
			name: "7",
			args: args{
				s: "abbaaa",
			},
			want: "NO",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := str.IsValid(tt.args.s); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
