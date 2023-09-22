package cgo

import "testing"

func Test_printByC(t *testing.T) {
	type args struct {
		num int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{
				num: 7,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printByC(tt.args.num)
		})
	}
}
