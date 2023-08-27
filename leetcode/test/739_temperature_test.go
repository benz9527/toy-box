package test

import (
	"github.com/benz9527/toy-box/leetcode/stack"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDailyTemperatures(t *testing.T) {
	type args struct {
		temperatures []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "1",
			args: args{
				temperatures: []int{
					73, 74, 75, 71, 69, 72, 76, 73,
				},
			},
			want: []int{
				1, 1, 4, 2, 1, 1, 0, 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, stack.DailyTemperatures(tt.args.temperatures), "DailyTemperatures(%v)", tt.args.temperatures)
		})
	}
}
