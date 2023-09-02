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

func TestMaxCopyFileSize(t *testing.T) {
	type args struct {
		files []int
		n     int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				files: []int{
					737270,
					737272,
					737288,
				},
				n: 3,
			},
			want: 1474542,
		},
		{
			name: "2",
			args: args{
				files: []int{
					400000,
					200000,
					200000,
					200000,
					400000,
					400000,
				},
				n: 6,
			},
			want: 1400000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxCopyFileSize(tt.args.files, tt.args.n); got != tt.want {
				t.Errorf("MaxCopyFileSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJumpGrids(t *testing.T) {
	type args struct {
		grids []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				grids: []int{
					1, 2, 3, 1,
				},
			},
			want: 4,
		},
		{
			name: "2",
			args: args{
				grids: []int{
					2, 7, 9, 3, 1,
				},
			},
			want: 12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JumpGrids(tt.args.grids); got != tt.want {
				t.Errorf("JumpGrids() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJumpGridsII(t *testing.T) {
	type args struct {
		grids []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				grids: []int{
					2, 3, 2,
				},
			},
			want: 3,
		},
		{
			name: "2",
			args: args{
				grids: []int{
					2,
				},
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JumpGridsII(tt.args.grids); got != tt.want {
				t.Errorf("JumpGridsII() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMelon(t *testing.T) {
	type args struct {
		stones []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				stones: []int{
					1, 1, 2, 2,
				},
			},
			want: 2,
		},
		{
			name: "2",
			args: args{
				stones: []int{
					1, 3, 2,
				},
			},
			want: 1,
		},
		{
			name: "3",
			args: args{
				stones: []int{
					1, 1, 1, 1, 1, 9, 8, 3, 7, 10,
				},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Melon(tt.args.stones); got != tt.want {
				t.Errorf("Melon() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMelon2(t *testing.T) {
	type args struct {
		stones []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				stones: []int{
					1, 1, 2, 2,
				},
			},
			want: 2,
		},
		{
			name: "2",
			args: args{
				stones: []int{
					1, 3, 2,
				},
			},
			want: 1,
		},
		{
			name: "3",
			args: args{
				stones: []int{
					1, 1, 1, 1, 1, 9, 8, 3, 7, 10,
				},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Melon2(tt.args.stones); got != tt.want {
				t.Errorf("Melon2() = %v, want %v", got, tt.want)
			}
		})
	}
}
