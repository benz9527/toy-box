package dichotomy

import "testing"

func TestDichotomyIndex(t *testing.T) {
	type args struct {
		nums   []int
		target int
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
					1, 2, 3, 4, 5, 5, 5, 5, 5, 6, 7, 8, 8, 9, 10,
				},
				target: 5,
			},
			want: 8,
		},
		{
			name: "2",
			args: args{
				nums: []int{
					1, 2, 3, 4, 5, 5, 5, 5, 5, 6, 7, 8, 8, 9, 11,
				},
				target: 10,
			},
			want: -15, // -14-1
		},
		{
			name: "3",
			args: args{
				nums: []int{
					1, 2, 3, 4, 5, 5, 5, 5, 5, 6, 7, 8, 8, 9, 11,
				},
				target: 12,
			},
			want: -16, // -15-1
		},
		{
			name: "4",
			args: args{
				nums: []int{
					1, 2, 2, 4, 5, 5, 5, 5, 5, 6, 7, 8, 8, 9, 11,
				},
				target: 3,
			},
			want: -4, // -3-1
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BinarySearchLastIndex(tt.args.nums, tt.args.target); got != tt.want {
				t.Errorf("BinarySearchLastIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBinarySearchFirstIndex(t *testing.T) {
	type args struct {
		nums   []int
		target int
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
					1, 2, 3, 4, 5, 5, 5, 5, 5, 6, 7, 8, 8, 9, 10,
				},
				target: 5,
			},
			want: 4,
		},
		{
			name: "2",
			args: args{
				nums: []int{
					1, 2, 3, 4, 5, 5, 5, 5, 5, 6, 7, 8, 8, 9, 11,
				},
				target: 10,
			},
			want: -15, // -14-1
		},
		{
			name: "3",
			args: args{
				nums: []int{
					1, 2, 3, 4, 5, 5, 5, 5, 5, 6, 7, 8, 8, 9, 11,
				},
				target: 12,
			},
			want: -16, // -15-1
		},
		{
			name: "4",
			args: args{
				nums: []int{
					1, 2, 2, 4, 5, 5, 5, 5, 5, 6, 7, 8, 8, 9, 11,
				},
				target: 3,
			},
			want: -4, // -3-1
		},
		{
			name: "5",
			args: args{
				nums: []int{
					2, 2, 2, 4, 5, 5, 5, 5, 5, 6, 7, 8, 8, 9, 11,
				},
				target: 2,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BinarySearchFirstIndex(tt.args.nums, tt.args.target); got != tt.want {
				t.Errorf("BinarySearchFirstIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}
