package prompt

import (
	"reflect"
	"testing"
)

func TestToSelectItems(t *testing.T) {
	type args struct {
		res string
	}
	tests := []struct {
		name string
		args args
		want []SelectItem
	}{
		{
			"valid: multiple suggestions",
			args{
				`Command: ls -a
Summary: summary1
Description: description1

Command: ls
Summary: summary2
Description: description2
`,
			},
			[]SelectItem{
				{
					Command:     "ls -a",
					Summary:     "summary1",
					Description: "description1",
				},
				{
					Command:     "ls",
					Summary:     "summary2",
					Description: "description2",
				},
			},
		},
		{
			"valid: reverse order",
			args{
				`Summary: summary1
Command: ls -a
Description: description1

Summary: summary2
Description: description1
Command: ls
`,
			},
			[]SelectItem{
				{
					Command:     "ls -a",
					Summary:     "summary1",
					Description: "description1",
				},
				{
					Command:     "ls",
					Summary:     "summary2",
					Description: "description1",
				},
			},
		},
		{
			"invalid",
			args{
				`Command: ls -a`,
			},
			[]SelectItem{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSelectItems(tt.args.res); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCompletion() = %v, want %v", got, tt.want)
			}
		})
	}
}
