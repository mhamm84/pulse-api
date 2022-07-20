package validator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {

	val := New()
	if val == nil {
		t.Errorf("call to New should always return a non nil pointer to a Validator")
	}
}

func TestValidator_AddError(t *testing.T) {
	type args struct {
		key     string
		message string
		want    bool
	}
	tests := []struct {
		name   string
		fields map[string]string
		args   args
	}{
		{name: "Add Error", fields: map[string]string{}, args: struct {
			key     string
			message string
			want    bool
		}{key: "id", message: "must be an integer value", want: false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Validator{
				Errors: tt.fields,
			}
			v.AddError(tt.args.key, tt.args.message)
			actual := v.Valid()
			assert.Equal(t, tt.args.want, actual)
		})
	}
}
