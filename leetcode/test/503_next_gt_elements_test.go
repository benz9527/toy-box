package test

import (
	"github.com/benz9527/toy-box/leetcode/stack"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNextGreaterElements(t *testing.T) {
	type args struct {
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
				nums: []int{
					1, 2, 1,
				},
			},
			want: []int{
				2, -1, 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, stack.NextGreaterElements(tt.args.nums), "NextGreaterElements(%v)", tt.args.nums)
		})
	}
}
