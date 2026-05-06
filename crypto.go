package main

import (
	"golang.org/x/crypto/scrypt"
)

func deriveKey(password string, salt []byte) []byte {
	key, _ := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	return key
}
