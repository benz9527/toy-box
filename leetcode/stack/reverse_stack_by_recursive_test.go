package stack

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReverseStack(t *testing.T) {
	type args struct {
		s *myStack
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "1",
			args: args{
				s: func() *myStack {
					s := &myStack{arr: make([]int, 0, 8)}
					s.Push(1)
					s.Push(2)
					s.Push(3)
					return s
				}(),
			},
			want: []int{3, 2, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReverseStack(tt.args.s)
			assert.Equal(t, tt.want, tt.args.s.GetArray())
		})
	}
}
