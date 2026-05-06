package main

import (
	"os"
	"path/filepath"
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

func TestDecrypt(t *testing.T) {
	plaintext := []byte("Hello, World! This is a test message.")
	password := "strongpassword"

	// Encrypt first
	ciphertext, err := encrypt(plaintext, password)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	// Decrypt
	decrypted, err := decrypt(ciphertext, password)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	// Should match original
	if string(decrypted) != string(plaintext) {
		t.Errorf("decrypted text doesn't match: got %q, want %q", decrypted, plaintext)
	}
}

func TestDecryptWrongPassword(t *testing.T) {
	plaintext := []byte("Secret message")
	password := "correctpassword"

	ciphertext, err := encrypt(plaintext, password)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	// Try decrypt with wrong password
	_, err = decrypt(ciphertext, "wrongpassword")
	if err == nil {
		t.Error("expected error with wrong password, got nil")
	}
}

func TestEncryptFile(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	// Create test file
	inputFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(inputFile, []byte("Hello, World!"), 0644)

	// Encrypt
	outputFile := filepath.Join(tmpDir, "test.txt.enc")
	err := encryptFile(inputFile, outputFile, "password123")
	if err != nil {
		t.Fatalf("encryptFile failed: %v", err)
	}

	// Verify output exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("encrypted file not created")
	}

	// Verify output is different from input
	input, _ := os.ReadFile(inputFile)
	output, _ := os.ReadFile(outputFile)
	if string(input) == string(output) {
		t.Error("encrypted file should be different from input")
	}
}

func TestDecryptFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create and encrypt test file
	inputFile := filepath.Join(tmpDir, "test.txt")
	original := []byte("Hello, World!")
	os.WriteFile(inputFile, original, 0644)

	encryptedFile := filepath.Join(tmpDir, "test.txt.enc")
	encryptFile(inputFile, encryptedFile, "password123")

	// Decrypt
	decryptedFile := filepath.Join(tmpDir, "test.txt.dec")
	err := decryptFile(encryptedFile, decryptedFile, "password123")
	if err != nil {
		t.Fatalf("decryptFile failed: %v", err)
	}

	// Verify content matches
	decrypted, _ := os.ReadFile(decryptedFile)
	if string(decrypted) != string(original) {
		t.Errorf("decrypted content doesn't match: got %q, want %q", decrypted, original)
	}
}
