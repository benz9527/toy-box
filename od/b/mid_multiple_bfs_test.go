package b

import "testing"

func TestInfectionDays(t *testing.T) {
	type args struct {
		area []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				area: []int{1, 0, 1, 0, 0, 0, 1, 0, 1},
			},
			want: 2,
		},
		{
			name: "2",
			args: args{
				area: []int{
					0, 0, 0, 0,
				},
			},
			want: -1,
		},
		{
			name: "3",
			args: args{
				area: []int{
					1, 1, 1, 1, 1, 1, 1, 1, 1,
				},
			},
			want: -1,
		},
		{
			name: "4",
			args: args{
				area: []int{
					0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0,
				},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InfectionDays(tt.args.area); got != tt.want {
				t.Errorf("InfectionDays() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInfectionDays2(t *testing.T) {
	type args struct {
		area []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				area: []int{1, 0, 1, 0, 0, 0, 1, 0, 1},
			},
			want: 2,
		},
		{
			name: "2",
			args: args{
				area: []int{
					0, 0, 0, 0,
				},
			},
			want: -1,
		},
		{
			name: "3",
			args: args{
				area: []int{
					1, 1, 1, 1, 1, 1, 1, 1, 1,
				},
			},
			want: -1,
		},
		{
			name: "4",
			args: args{
				area: []int{
					0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0,
				},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InfectionDays2(tt.args.area); got != tt.want {
				t.Errorf("InfectionDays2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindMaxValueStockHeapByBFS(t *testing.T) {
	type args struct {
		matrix [][]int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				matrix: [][]int{
					{2, 2, 2, 2, 0},
					{0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0},
					{0, 1, 1, 1, 1},
				},
			},
			want: 8,
		},
		{
			name: "2",
			args: args{
				matrix: [][]int{
					{2, 2, 2, 2, 0},
					{0, 0, 0, 2, 0},
					{0, 0, 0, 1, 0},
					{0, 1, 1, 1, 1},
				},
			},
			want: 15,
		},
		{
			name: "3",
			args: args{
				matrix: [][]int{
					{2, 0, 0, 0, 0},
					{0, 0, 0, 2, 0},
					{0, 0, 0, 0, 0},
					{0, 0, 1, 1, 1},
				},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindMaxValueStockHeapByBFS(tt.args.matrix); got != tt.want {
				t.Errorf("FindMaxValueStockHeapByBFS() = %v, want %v", got, tt.want)
			}
		})
	}
}
