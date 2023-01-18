package test

import (
	"testing"
)

func TestMakeAlert(t *testing.T) {
	type args struct {
		u string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "make alert",
			args: args{
				u: "http://192.168.10.101:8001/notify",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := MakeAlert(tt.args.u)

			if (err != nil) != tt.wantErr {
				t.Errorf("MakeAlert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
