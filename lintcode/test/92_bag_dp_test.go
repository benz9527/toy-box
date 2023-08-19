package test

import (
	"github.com/benz9527/toy-box/lintcode/dp"
	"testing"
)

func TestBackPack(t *testing.T) {
	type args struct {
		m int
		a []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				m: 10,
				a: []int{
					3, 4, 8, 5,
				},
			},
			want: 9,
		},
		{
			name: "2",
			args: args{
				m: 12,
				a: []int{
					2, 3, 5, 7,
				},
			},
			want: 12,
		},
		{
			name: "3",
			args: args{
				m: 33,
				a: []int{
					3, 4, 11, 50, 5, 1,
				},
			},
			want: 24,
		},
		{
			name: "4",
			args: args{
				m: 5000,
				a: []int{
					828, 125, 740, 724, 983, 321, 773, 678, 841, 842, 875, 377, 674, 144,
					340, 467, 625, 916, 463, 922, 255, 662, 692, 123, 778, 766, 254, 559,
					480, 483, 904, 60, 305, 966, 872, 935, 626, 691, 832, 998, 508, 657,
					215, 162, 858, 179, 869, 674, 452, 158, 520, 138, 847, 452, 764, 995,
					600, 568, 92, 496, 533, 404, 186, 345, 304, 420, 181, 73, 547, 281,
					374, 376, 454, 438, 553, 929, 140, 298, 451, 674, 91, 531, 685, 862, 446,
					262, 477, 573, 627, 624, 814, 103, 294, 388,
				},
			},
			want: 5000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dp.BackPack(tt.args.m, tt.args.a); got != tt.want {
				t.Errorf("BackPack() = %v, want %v", got, tt.want)
			}
		})
	}
}
