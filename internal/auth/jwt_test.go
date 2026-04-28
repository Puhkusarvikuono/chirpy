package auth

import (
	"testing"
	"github.com/google/uuid"
	"time"
	"fmt"
)


func TestJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "passw0rd"
	expiresIn, err := time.ParseDuration("5m")
	if err != nil {
		t.Errorf("Error parsing duration, err: %v\n", err)
	}

	jwt, err := MakeJWT(userID, tokenSecret, expiresIn)
	
	if err != nil {
		t.Errorf("Err is not nil, %v\n", err)
	}

	fmt.Printf("Jwt: %v\nErr: %v\n", jwt, err)
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validSecret := "super-secret"

	// Pre-make a valid token you can reuse across cases
	validToken, err := MakeJWT(userID, validSecret, time.Hour)
	if err != nil {
		t.Fatalf("failed to create valid token: %v", err)
	}

	// Pre-make an expired token
	expiredToken, err := MakeJWT(userID, validSecret, -time.Hour)
	if err != nil {
		t.Fatalf("failed to create expired token: %v", err)
	}

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "valid token returns correct user ID",
			tokenString: validToken,
			tokenSecret: validSecret,
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "expired token is rejected",
			tokenString: expiredToken,
			tokenSecret: validSecret,
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		// add more cases here...
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotID, err := ValidateJWT(tc.tokenString, tc.tokenSecret)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if gotID != tc.wantUserID {
				t.Errorf("ValidateJWT() gotID = %v, want %v", gotID, tc.wantUserID)
			}
		})
	}
}
