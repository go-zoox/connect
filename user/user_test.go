package user

import (
	"testing"

	"github.com/go-zoox/crypto/jwt"
)

func TestEncodeDecodePermissions(t *testing.T) {
	signer := jwt.New("test-secret-key")

	encoded, err := (&User{
		ID:          "1",
		Nickname:    "Zero",
		Email:       "zero@gozoox.com",
		Username:    "zero",
		Permissions: []string{"ADMIN", "DEVELOPER"},
	}).Encode(signer)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}

	decoded := &User{}
	if err := decoded.Decode(signer, encoded); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if len(decoded.Permissions) != 2 || decoded.Permissions[0] != "ADMIN" || decoded.Permissions[1] != "DEVELOPER" {
		t.Fatalf("permissions want [ADMIN DEVELOPER], got %v", decoded.Permissions)
	}
}

func TestDecodeWithoutPermissions(t *testing.T) {
	signer := jwt.New("test-secret-key")

	encoded, err := (&User{ID: "1", Username: "zero"}).Encode(signer)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}

	decoded := &User{}
	if err := decoded.Decode(signer, encoded); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if len(decoded.Permissions) != 0 {
		t.Fatalf("permissions want empty, got %v", decoded.Permissions)
	}
	if decoded.Username != "zero" {
		t.Fatalf("username want zero, got %q", decoded.Username)
	}
}
