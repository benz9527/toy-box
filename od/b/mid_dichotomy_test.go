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

func TestAngryStudentsAreTeachable(t *testing.T) {
	type args struct {
		allStudents    []int
		badStudentIdxs []int
		tolerance      int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				allStudents:    []int{1810, 1809, 1801, 1802},
				badStudentIdxs: []int{0, 1},
				tolerance:      3,
			},
			want: disteachable,
		},
		{
			name: "2",
			args: args{
				allStudents:    []int{1801, 1811, 1811, 1802, 1804, 1803},
				badStudentIdxs: []int{1, 2, 4},
				tolerance:      3,
			},
			want: disteachable,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AngryStudentsAreTeachable(tt.args.allStudents, tt.args.badStudentIdxs, tt.args.tolerance); got != tt.want {
				t.Errorf("AngryStudentsAreTeachable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProgramPractice(t *testing.T) {
	type args struct {
		maxDays       int
		practiceTimes []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				maxDays: 5,
				practiceTimes: []int{
					1, 2, 2, 3, 5, 4, 6, 7, 8,
				},
			},
			want: 4,
		},
		{
			name: "2",
			args: args{
				maxDays: 4,
				practiceTimes: []int{
					999, 999, 999,
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProgramPractice(tt.args.maxDays, tt.args.practiceTimes); got != tt.want {
				t.Errorf("ProgramPractice() = %v, want %v", got, tt.want)
			}
		})
	}
}
