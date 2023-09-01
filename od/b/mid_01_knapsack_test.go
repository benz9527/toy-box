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
