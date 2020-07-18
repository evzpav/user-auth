package errors_test

import (
	"testing"

	"gitlab.com/evzpav/user-auth/pkg/errors"
)

func TestInfo_Error(t *testing.T) {
	tests := []struct {
		name string
		info errors.Info
		want string
	}{
		{
			name: "error full",
			info: errors.Info{
				Code:    "ERROR_CODE",
				Message: "error message",
				Args: map[string]interface{}{
					"errorArg1": "errorArg1Value",
				},
			},
			want: "<ERROR_CODE> error message (errorArg1: errorArg1Value)",
		},
		{
			name: "error without code",
			info: errors.Info{
				Message: "error message",
				Args: map[string]interface{}{
					"errorArg1": "errorArg1Value",
				},
			},
			want: "error message (errorArg1: errorArg1Value)",
		},
		{
			name: "error without message",
			info: errors.Info{
				Code: "ERROR_CODE",
			},
			want: "<ERROR_CODE> ",
		},
		{
			name: "error without arguments",
			info: errors.Info{
				Code:    "ERROR_CODE",
				Message: "error message",
			},
			want: "<ERROR_CODE> error message",
		},
		{
			name: "error without code and arguments",
			info: errors.Info{
				Message: "error message",
			},
			want: "error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.info.Error(); got != tt.want {
				t.Errorf("Info.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
