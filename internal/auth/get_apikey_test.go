package auth

import (
	"net/http"
	"testing"
)

func TestGetAPIKey(t *testing.T) {
	tests := []struct {
		name    string
		header  http.Header
		want    string
		wantErr bool
	}{
		{
			name:    "valid token",
			header:  http.Header{"Authorization": {"ApiKey abc123"}},
			want:    "abc123",
			wantErr: false,
		},
		{
			name:    "valid token with extra spaces",
			header:  http.Header{"Authorization": {"  ApiKey   xyz789  "}},
			want:    "xyz789",
			wantErr: false,
		},
		{
			name:    "missing header",
			header:  http.Header{},
			want:    "",
			wantErr: true,
		},
		{
			name:    "wrong scheme",
			header:  http.Header{"Authorization": {"Basic abc123"}},
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty APIkey value",
			header:  http.Header{"Authorization": {"ApiKey "}},
			want:    "",
			wantErr: true,
		},
		{
			name:    "multiple header values - first used",
			header:  http.Header{"Authorization": {"ApiKey first", "ApiKey second"}},
			want:    "first",
			wantErr: false,
		},
		{
			name:    "non-apikey first value but later apikey",
			header:  http.Header{"Authorization": {"Basic foo", "ApiKey later"}},
			want:    "later",
			wantErr: false,
		},
		{
			name:    "mixed-case scheme",
			header:  http.Header{"Authorization": {"aPiKeY MixedCase"}},
			want:    "MixedCase",
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GetAPIKey(tc.header)

			if (err != nil) != tc.wantErr {
				t.Fatalf("GetAPIKey() error = %v, wantErr = %v", err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("GetAPIKey() = %q, want %q", got, tc.want)
			}
		})
	}
}
