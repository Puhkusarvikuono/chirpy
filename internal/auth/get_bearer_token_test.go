package auth

import (
	"testing"
	"net/http"
)


func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		header    http.Header
		want      string
		wantErr   bool
	}{
		{
			name:    "valid token",
			header:  http.Header{"Authorization": {"Bearer abc123"}},
			want:    "abc123",
			wantErr: false,
		},
		{
			name:    "valid token with extra spaces",
			header:  http.Header{"Authorization": {"  Bearer   xyz789  "}},
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
			name:    "empty bearer value",
			header:  http.Header{"Authorization": {"Bearer "}},
			want:    "",
			wantErr: true,
		},
		{
			name:    "multiple header values - first used",
			header:  http.Header{"Authorization": {"Bearer first", "Bearer second"}},
			want:    "first",
			wantErr: false,
		},
		{
			name:    "non-bearer first value but later bearer",
			header:  http.Header{"Authorization": {"Basic foo", "Bearer later"}},
			want:    "later",
			wantErr: false,
		},
		{
			name:    "mixed-case scheme",
			header:  http.Header{"Authorization": {"bEaReR MixedCase"}},
			want:    "MixedCase",
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GetBearerToken(tc.header)

			if (err != nil) != tc.wantErr {
				t.Fatalf("GetBearerToken() error = %v, wantErr = %v", err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("GetBearerToken() = %q, want %q", got, tc.want)
			}
		})
	}
}

