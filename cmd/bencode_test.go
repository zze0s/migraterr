package cmd

import "testing"

func Test_validateReplacement(t *testing.T) {
	type args struct {
		replacement string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test", args: args{replacement: "/path/|/new/path/"}, wantErr: false},
		{name: "fail", args: args{replacement: "/path//new/path/"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateReplacement(tt.args.replacement); (err != nil) != tt.wantErr {
				t.Errorf("validateReplacement() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
