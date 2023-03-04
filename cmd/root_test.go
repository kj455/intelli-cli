package cmd

import (
	"reflect"
	"testing"
)

func TestParseCompletion(t *testing.T) {
	type args struct {
		res string
	}
	tests := []struct {
		name string
		args args
		want []Suggestion
	}{
		{
			"valid - multiple suggestions",
			args{
				`command: ls -a
note: note1

command: ls
note: note2
`,
			},
			[]Suggestion{
				{
					Command: "ls -a",
					Note:    "note1",
				},
				{
					Command: "ls",
					Note:    "note2",
				},
			},
		},
		{
			"valid - reverse order",
			args{
				`note: note1
command: ls -a

note: note2
command: ls
`,
			},
			[]Suggestion{
				{
					Command: "ls -a",
					Note:    "note1",
				},
				{
					Command: "ls",
					Note:    "note2",
				},
			},
		},
		{
			"invalid - lacking note",
			args{
				`command: ls -a`,
			},
			[]Suggestion{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseCompletion(tt.args.res); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCompletion() = %v, want %v", got, tt.want)
			}
		})
	}
}
