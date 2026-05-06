package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/scrypt"
)

func deriveKey(password string, salt []byte) ([]byte, error) {
	key, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func encrypt(plaintext []byte, password string) ([]byte, error) {
	// Generate random salt
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("generating salt: %w", err)
	}

	// Derive key
	key, err := deriveKey(password, salt)
	if err != nil {
		return nil, fmt.Errorf("deriving key: %w", err)
	}

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("creating cipher: %w", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("creating GCM: %w", err)
	}

	// Generate random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("generating nonce: %w", err)
	}

	// Encrypt
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	// Combine: salt + nonce + ciphertext
	result := make([]byte, 0, len(salt)+len(nonce)+len(ciphertext))
	result = append(result, salt...)
	result = append(result, nonce...)
	result = append(result, ciphertext...)

	return result, nil
}

func decrypt(ciphertext []byte, password string) ([]byte, error) {
	if len(ciphertext) < 16+12+16 {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// Extract salt, nonce, and encrypted data
	salt := ciphertext[:16]
	nonce := ciphertext[16:28]
	encrypted := ciphertext[28:]

	// Derive key
	key, err := deriveKey(password, salt)
	if err != nil {
		return nil, fmt.Errorf("deriving key: %w", err)
	}

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("creating cipher: %w", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("creating GCM: %w", err)
	}

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed (wrong password or corrupted file): %w", err)
	}

	return plaintext, nil
}

func encryptFile(inputPath, outputPath, password string) error {
	// Read input file
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("reading input file: %w", err)
	}

	// Encrypt
	ciphertext, err := encrypt(plaintext, password)
	if err != nil {
		return fmt.Errorf("encrypting: %w", err)
	}

	// Write output file
	if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
		return fmt.Errorf("writing output file: %w", err)
	}

	return nil
}

func decryptFile(inputPath, outputPath, password string) error {
	// Read input file
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("reading input file: %w", err)
	}

	// Decrypt
	plaintext, err := decrypt(ciphertext, password)
	if err != nil {
		return fmt.Errorf("decrypting: %w", err)
	}

	// Write output file
	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("writing output file: %w", err)
	}

	return nil
}
