package b

import "testing"

func TestFullCarForTravel(t *testing.T) {
	type args struct {
		nums []int
		n    int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				nums: []int{
					5, 4, 2, 3, 2, 4, 9,
				},
				n: 10,
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FullCarForTravel(tt.args.nums, tt.args.n); got != tt.want {
				t.Errorf("FullCarForTravel() = %v, want %v", got, tt.want)
			}
		})
	}
}
