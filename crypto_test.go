package main

import (
	"testing"
)

func TestDeriveKey(t *testing.T) {
	password := "testpassword123"
	salt := make([]byte, 16)
	for i := range salt {
		salt[i] = byte(i)
	}

	key1, err := deriveKey(password, salt)
	if err != nil {
		t.Fatalf("deriveKey failed: %v", err)
	}
	key2, err := deriveKey(password, salt)
	if err != nil {
		t.Fatalf("deriveKey failed: %v", err)
	}

	if len(key1) != 32 {
		t.Errorf("expected key length 32, got %d", len(key1))
	}

	// Same password + salt should produce same key
	for i := range key1 {
		if key1[i] != key2[i] {
			t.Error("same password+salt should produce same key")
			break
		}
	}

	// Different salt should produce different key
	salt2 := make([]byte, 16)
	for i := range salt2 {
		salt2[i] = byte(i + 1)
	}
	key3, err := deriveKey(password, salt2)
	if err != nil {
		t.Fatalf("deriveKey failed: %v", err)
	}

	same := true
	for i := range key1 {
		if key1[i] != key3[i] {
			same = false
			break
		}
	}
	if same {
		t.Error("different salt should produce different key")
	}
}
