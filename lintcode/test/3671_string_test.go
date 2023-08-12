package test

import (
	"github.com/benz9527/toy-box/lintcode/str"
	"reflect"
	"testing"
)

func TestReplaceSynonyms(t *testing.T) {
	type args struct {
		synonyms [][]string
		text     string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "1",
			args: args{
				synonyms: [][]string{
					{"happy", "joy"}, {"sad", "sorrow"}, {"joy", "cheerful"},
				},
				text: "I am happy today but was sad yesterday",
			},
			want: []string{
				"I am cheerful today but was sad yesterday",
				"I am cheerful today but was sorrow yesterday",
				"I am happy today but was sad yesterday",
				"I am happy today but was sorrow yesterday",
				"I am joy today but was sad yesterday",
				"I am joy today but was sorrow yesterday",
			},
		},
		{
			name: "2",
			args: args{
				synonyms: [][]string{
					{"happy", "joy"}, {"glad", "cheerful"},
				},
				text: "I am happy today but was sad yesterday",
			},
			want: []string{
				"I am happy today but was sad yesterday",
				"I am joy today but was sad yesterday",
			},
		},
		{
			name: "3",
			args: args{
				synonyms: [][]string{
					{"check", "see"}, {"see", "look"},
				},
				text: "Let me see see",
			},
			want: []string{
				"Let me check check",
				"Let me check look",
				"Let me check see",
				"Let me look check",
				"Let me look look",
				"Let me look see",
				"Let me see check",
				"Let me see look",
				"Let me see see",
			},
		},
		{
			name: "4",
			args: args{
				synonyms: [][]string{
					{"Pleasure", "Delight"}, {"Cheer", "Joy"}, {"Happiness", "Pleasure"}, {"Delight", "Joy"},
				},
				text: "Happiness is not about being immortal nor having food.",
			},
			want: []string{
				"Cheer is not about being immortal nor having food.",
				"Delight is not about being immortal nor having food.",
				"Happiness is not about being immortal nor having food.",
				"Joy is not about being immortal nor having food.",
				"Pleasure is not about being immortal nor having food.",
			},
		},
		{
			name: "5",
			args: args{
				synonyms: [][]string{
					{"Pleasure", "Delight"}, {"Cheer", "Joy"}, {"Happiness", "Pleasure"}, {"Delight", "Joy"},
					{"immortal", "timeless"}, {"timeless", "eternal"}, {"eternal", "lasting"},
				},
				text: "Happiness is not about being immortal nor having food.",
			},
			want: []string{

				"Cheer is not about being eternal nor having food.",
				"Cheer is not about being immortal nor having food.",
				"Cheer is not about being lasting nor having food.",
				"Cheer is not about being timeless nor having food.",
				"Delight is not about being eternal nor having food.",
				"Delight is not about being immortal nor having food.",
				"Delight is not about being lasting nor having food.",
				"Delight is not about being timeless nor having food.",
				"Happiness is not about being eternal nor having food.",
				"Happiness is not about being immortal nor having food.",
				"Happiness is not about being lasting nor having food.",
				"Happiness is not about being timeless nor having food.",
				"Joy is not about being eternal nor having food.",
				"Joy is not about being immortal nor having food.",
				"Joy is not about being lasting nor having food.",
				"Joy is not about being timeless nor having food.",
				"Pleasure is not about being eternal nor having food.",
				"Pleasure is not about being immortal nor having food.",
				"Pleasure is not about being lasting nor having food.",
				"Pleasure is not about being timeless nor having food.",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := str.ReplaceSynonyms(tt.args.synonyms, tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReplaceSynonyms() = %v, want %v", got, tt.want)
			}
		})
	}
}
