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

func TestEncrypt(t *testing.T) {
	plaintext := []byte("Hello, World! This is a test message.")
	password := "strongpassword"

	ciphertext, err := encrypt(plaintext, password)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	// ciphertext should be longer than plaintext (salt + nonce + tag)
	if len(ciphertext) <= len(plaintext) {
		t.Error("ciphertext should be longer than plaintext")
	}

	// First 16 bytes should be salt
	// Next 12 bytes should be nonce
	// Rest should be encrypted data + tag
	if len(ciphertext) != 16+12+len(plaintext)+16 {
		t.Errorf("unexpected ciphertext length: %d", len(ciphertext))
	}
}
