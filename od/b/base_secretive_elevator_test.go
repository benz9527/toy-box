package b

import (
	"reflect"
	"testing"
)

func TestSecretiveElevator(t *testing.T) {
	type args struct {
		nums   []int
		target int
		n      int
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
					1, 2, 6,
				},
				target: 5,
				n:      3,
			},
			want: []int{
				6, 2, 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SecretiveElevator(tt.args.nums, tt.args.target, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SecretiveElevator() = %v, want %v", got, tt.want)
			}
		})
	}
}
