package main

import (
	"crypto/rand"
	"testing"
)

func TestTextEncryption(t *testing.T) {
	// generate test key
	test_key := make([]byte, 32)
	rand.Read(test_key)

	test_text := "Hello World123&*^%$#@!"

	// convert to bytes
	test_bytes := []byte(test_text)

	// encrypt
	encrypted_text := Encrypt(test_key, test_bytes)

	// decrypt
	decrypted_text := Decrypt(test_key, encrypted_text)

	// check if decrypted text is the same as the original
	if string(decrypted_text) != test_text {
		t.Errorf("Encryption/Decryption failed")
	}
}
