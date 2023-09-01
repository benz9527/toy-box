package b

import "testing"

func TestBuyMachines(t *testing.T) {
	type args struct {
		totalSum        int
		totalComponents int
		components      []component
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				totalSum:        500,
				totalComponents: 3,
				components: []component{
					{
						kind:        0,
						reliability: 80,
						price:       100,
					},
					{
						kind:        0,
						reliability: 90,
						price:       200,
					},
					{
						kind:        1,
						reliability: 50,
						price:       50,
					},
					{
						kind:        1,
						reliability: 70,
						price:       210,
					},
					{
						kind:        2,
						reliability: 50,
						price:       100,
					},
					{
						kind:        2,
						reliability: 60,
						price:       150,
					},
				},
			},
			want: 60,
		},
		{
			name: "2",
			args: args{
				totalSum:        100,
				totalComponents: 1,
				components: []component{
					{
						kind:        0,
						reliability: 90,
						price:       200,
					},
				},
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuyMachines(tt.args.totalSum, tt.args.totalComponents, tt.args.components); got != tt.want {
				t.Errorf("BuyMachines() = %v, want %v", got, tt.want)
			}
		})
	}
}
