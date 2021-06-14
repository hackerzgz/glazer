package cmd

import (
	"reflect"
	"testing"

	"github.com/tidwall/gjson"
)

func Test_generateFakerArray(t *testing.T) {
	type args struct {
		raw gjson.Result
	}
	tests := []struct {
		name      string
		args      args
		wantFaker []interface{}
	}{
		{
			name: "string array",
			args: args{
				raw: gjson.Result{
					Type: gjson.JSON, Raw: `["test", "string"]`,
				},
			},
			wantFaker: []interface{}{"test", "string"},
		},
		{
			name: "integer array",
			args: args{
				raw: gjson.Result{
					Type: gjson.JSON, Raw: `[1, 2, 3]`,
				},
			},
			wantFaker: []interface{}{
				1.0, 2.0, 3.0,
			},
		},
		{
			name: "struct array",
			args: args{
				raw: gjson.Result{
					Type: gjson.JSON, Raw: `[{"k1": "v1"}]`,
				},
			},
			wantFaker: []interface{}{map[string]interface{}{"k1": "v1"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotFaker := generateFakerArray(tt.args.raw); !reflect.DeepEqual(gotFaker, tt.wantFaker) {
				t.Errorf("generateFakerArray() = %v, want %v", gotFaker, tt.wantFaker)
			}
		})
	}
}
