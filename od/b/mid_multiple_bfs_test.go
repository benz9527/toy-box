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
