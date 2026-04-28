package auth

import (
	"testing"
)

func TestPasswordHashEqual(t *testing.T) {
	password := "Passw0rd"
	hashedPassword, err := HashPassword(password)

	if err != nil {
		t.Errorf("Hash function returns error. %n\n", err)
	}

	if password == hashedPassword {
		t.Errorf("Expected hashed password")
	}
}

func TestPasswordHashNotEqual(t *testing.T) {
	password := "Passw0rd"
	password2 := "Passw0rd2"
	hashedPassword, err := HashPassword(password)

	if err != nil {
		t.Errorf("Hash function returns error. %n\n", err)
	}

	hashedPassword2, err := HashPassword(password2)

	if err != nil {
		t.Errorf("Hash function returns error. %n\n", err)
	}

	if hashedPassword == hashedPassword2 {
		t.Errorf("Hash should not be equal.")
	}
}

func TestPasswordHashMatchEqual(t *testing.T) {
	password := "Passw0rd"
	hashedPassword, err := HashPassword(password)

	if err != nil {
		t.Errorf("Hash function returns error. %n\n", err)
	}

	ok, err := CheckPasswordHash(password, hashedPassword)

	if !ok || err != nil {
		t.Errorf("Checking matching password & hash. Match is %v. Err: %v\n.", ok, err)
	}
}

func TestPasswordHashMatchNotEqual(t *testing.T) {
	password := "Passw0rd"
	password2 := "raisins"
	hashedPassword, _ := HashPassword(password)
	hashedPassword2, _ := HashPassword(password2)

	ok, err := CheckPasswordHash(password2, hashedPassword)

	if ok || err != nil {
		t.Errorf("Checking mismatch password & hash. Match is %v. Err %v\n%s\n%s\n%s\n", ok, err, password2, hashedPassword, hashedPassword2)
	}

	ok, err = CheckPasswordHash(password, hashedPassword2)

	if ok || err != nil {
		t.Errorf("Checking mismatch password & hash. Match is %v. Err %v\n%s\n%s\n%s\n", ok, err, password, hashedPassword2, hashedPassword)
	}
}
