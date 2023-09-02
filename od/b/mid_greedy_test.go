package b

import "testing"

func TestGetMaxScoreOfPokers(t *testing.T) {
	type args struct {
		pokerStr string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				pokerStr: "33445677",
			},
			want: 67,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMaxScoreOfPokers(tt.args.pokerStr); got != tt.want {
				t.Errorf("GetMaxScoreOfPokers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCompetitionAWinMaxScore(t *testing.T) {
	type args struct {
		seqA []int
		seqB []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				seqA: []int{4, 8, 10},
				seqB: []int{3, 6, 4},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCompetitionAWinMaxScore(tt.args.seqA, tt.args.seqB); got != tt.want {
				t.Errorf("GetCompetitionAWinMaxScore() = %v, want %v", got, tt.want)
			}
		})
	}
}
