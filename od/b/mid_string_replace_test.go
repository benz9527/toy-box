package b

import "testing"

func TestGetReplaceString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				str: "()abd",
			},
			want: "abd",
		},
		{
			name: "2",
			args: args{
				str: "(abd)demand(fb)()for",
			},
			want: "aemanaaor",
		},
		{
			name: "3",
			args: args{
				str: "()happy(xyz)new(wxy)year(t)",
			},
			want: "happwnewwear",
		},
		{
			name: "4",
			args: args{
				str: "()abcdefgAC(a)(Ab)(C)",
			},
			want: "AAcdefgAC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetReplaceString(tt.args.str); got != tt.want {
				t.Errorf("GetReplaceString() = %v, want %v", got, tt.want)
			}
		})
	}
}
