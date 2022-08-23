package data

import "testing"

func TestPermissions_Included(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name string
		p    Permissions
		args args
		want bool
	}{
		{name: "T1", p: []string{"economic:all"}, args: args{"economic:all"}, want: true},
		{name: "T2", p: []string{}, args: args{"economic:all"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Included(tt.args.code); got != tt.want {
				t.Errorf("Included() = %v, want %v", got, tt.want)
			}
		})
	}
}
