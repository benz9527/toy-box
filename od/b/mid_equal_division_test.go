package b

import "testing"

func TestBricksEqualDivision(t *testing.T) {
	type args struct {
		bricks []int
	}
	type want struct {
		weight   int
		canEqDiv bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "1",
			args: args{
				bricks: []int{6, 3, 5},
			},
			want: want{
				weight:   11,
				canEqDiv: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weight, canEqDiv := BricksEqualDivision(tt.args.bricks)
			if weight != tt.want.weight && canEqDiv != tt.want.canEqDiv {
				t.Errorf("BricksEqualDivision() got weight = %v & can equal division = %v, want %#v", weight, canEqDiv, tt.want)
			}
		})
	}
}
