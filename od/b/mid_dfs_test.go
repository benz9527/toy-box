package b

import (
	"reflect"
	"testing"
)

func TestCommentsTransfer(t *testing.T) {
	type args struct {
		comments string
	}
	tests := []struct {
		name         string
		args         args
		wantLevels   int
		wantComments []string
	}{
		{
			name: "1",
			args: args{
				comments: "hello,2,ok,0,bye,0,test,0,one,1,two,1,a,0",
			},
			wantLevels: 3,
			wantComments: []string{
				"hello test one",
				"ok bye two",
				"a",
			},
		},
		{
			name: "2",
			args: args{
				comments: "A,5,A,0,a,0,A,0,a,0,A,0",
			},
			wantLevels: 2,
			wantComments: []string{
				"A",
				"A a A a A",
			},
		},
		{
			name: "3",
			args: args{
				comments: "A,3,B,2,C,0,D,1,E,0,F,1,G,0,H,1,I,1,J,0,K,1,L,0,M,2,N,0,O,1,P,0",
			},
			wantLevels: 4,
			wantComments: []string{
				"A K M",
				"B F H L N O",
				"C D G I P",
				"E J",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := CommentsTransfer(tt.args.comments)
			if got != tt.wantLevels {
				t.Errorf("CommentsTransfer() got = %v, wantLevels %v", got, tt.wantLevels)
			}
			if !reflect.DeepEqual(got1, tt.wantComments) {
				t.Errorf("CommentsTransfer() got1 = %v, wantLevels %v", got1, tt.wantComments)
			}
		})
	}
}
