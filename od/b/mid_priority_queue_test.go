package b

import (
	"reflect"
	"testing"
)

func TestShoppingPlans(t *testing.T) {
	type args struct {
		n      int
		k      int
		prices []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "1",
			args: args{
				n: 5,
				k: 6,
				prices: []int{
					1, 1, 2, 3, 3,
				},
			},
			want: []int{
				1, 1, 2, 2, 3, 3,
			},
		},
		{
			name: "2",
			args: args{
				n: 3,
				k: 5,
				prices: []int{
					1, 100, 101,
				},
			},
			want: []int{
				1, 100, 101, 101, 102,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShoppingPlans(tt.args.n, tt.args.k, tt.args.prices); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ShoppingPlans() = %v, want %v", got, tt.want)
			}
		})
	}
}
