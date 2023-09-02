package b

import (
	"reflect"
	"testing"
)

func TestKMPGetNextTable(t *testing.T) {
	type args struct {
		tmpl string
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "1",
			args: args{
				tmpl: "aabaa",
			},
			want: []int{0, 1, 0, 1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KMPGetNextTable(tt.args.tmpl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KMPGetNextTable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKMPIndexOf(t *testing.T) {
	type args struct {
		src  string
		tmpl string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				src:  "xyz",
				tmpl: "z",
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KMPIndexOf(tt.args.src, tt.args.tmpl); got != tt.want {
				t.Errorf("KMPIndexOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMinLoopSubStr(t *testing.T) {
	type args struct {
		n    int
		nums []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "1",
			args: args{
				n: 9,
				nums: []int{
					1, 2, 1, 1, 2, 1, 1, 2, 1,
				},
			},
			want: []int{1, 2, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMinLoopSubStr(tt.args.n, tt.args.nums); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMinLoopSubStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
