package usecase

import (
	"testing"
)

func TestGenerateHash_WithoutSalt(t *testing.T) {
	t.Parallel()

	argon2Hash := NewArgon2idHash(1, 32, 64*1024, 1, 128)

	hashSalt, err := argon2Hash.GenerateHash([]byte("testpassword"), nil)
	if err != nil {
		t.Fatalf("GenerateHash() unexpected error: %v", err)
	}

	if hashSalt.Hash == nil {
		t.Error("GenerateHash() hash should not be nil")
	}

	if hashSalt.Salt == nil {
		t.Error("GenerateHash() salt should not be nil")
	}

	if len(hashSalt.Salt) != 32 {
		t.Errorf("GenerateHash() salt length = %d, want %d", len(hashSalt.Salt), 32)
	}
}

func TestGenerateHash_WithSalt(t *testing.T) {
	t.Parallel()

	argon2Hash := NewArgon2idHash(1, 32, 64*1024, 1, 128)

	providedSalt := []byte("this-is-a-fixed-salt-for-testing")

	hashSalt, err := argon2Hash.GenerateHash([]byte("testpassword"), providedSalt)
	if err != nil {
		t.Fatalf("GenerateHash() unexpected error: %v", err)
	}

	if hashSalt.Hash == nil {
		t.Error("GenerateHash() hash should not be nil")
	}

	if string(hashSalt.Salt) != string(providedSalt) {
		t.Errorf("GenerateHash() salt = %v, want %v", hashSalt.Salt, providedSalt)
	}
}

func TestCompare_Matching(t *testing.T) {
	t.Parallel()

	argon2Hash := NewArgon2idHash(1, 32, 64*1024, 1, 128)

	password := []byte("testpassword")

	hashSalt, err := argon2Hash.GenerateHash(password, nil)
	if err != nil {
		t.Fatalf("GenerateHash() unexpected error: %v", err)
	}

	if err := argon2Hash.Compare(hashSalt.Hash, hashSalt.Salt, password); err != nil {
		t.Errorf("Compare() should match, got error: %v", err)
	}
}

func TestCompare_NonMatching(t *testing.T) {
	t.Parallel()

	argon2Hash := NewArgon2idHash(1, 32, 64*1024, 1, 128)

	hashSalt, err := argon2Hash.GenerateHash([]byte("testpassword"), nil)
	if err != nil {
		t.Fatalf("GenerateHash() unexpected error: %v", err)
	}

	err = argon2Hash.Compare(hashSalt.Hash, hashSalt.Salt, []byte("wrongpassword"))
	if err == nil {
		t.Error("Compare() should return error for non-matching password")
	}
}
