package b

import "testing"

func TestElectionByMoney(t *testing.T) {
	type args struct {
		nStores int
		votes   [][2]int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				nStores: 5,
				votes: [][2]int{
					{2, 10},
					{3, 20},
					{4, 30},
					{5, 40},
					{5, 90},
				},
			},
			want: 50,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ElectionByMoney(tt.args.nStores, tt.args.votes); got != tt.want {
				t.Errorf("ElectionByMoney() = %v, want %v", got, tt.want)
			}
		})
	}
}
