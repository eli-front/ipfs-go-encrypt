package main

import (
	// "context"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

func main() {
	// Connect to your local IPFS deamon running in the background.

	// Where your local node is running on localhost:5001
	sh := shell.NewShell("localhost:5001")

	// Get input for the file to be added to IPFS from the command line
	filepath := os.Args[1]

	// Read the file from the path provided
	file_data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get the key
	key, err := os.ReadFile("key")
	if err != nil {
		log.Fatal(err)
	}

	// Encrypt the file data with the key
	encrypted_data := Encrypt(key, file_data)

	// Add the encrypted file to IPFS
	cid, err := sh.Add(strings.NewReader(string(encrypted_data)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}
	fmt.Printf("added %s\n", cid)

	// Get the encrypted file from IPFS
	data, err := sh.Cat(cid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}

	// Use a buffer to read the data from the reader
	buf := new(bytes.Buffer)
	buf.ReadFrom(data)
	encrypted_string := buf.String()
	fmt.Printf("encrypted data %s", encrypted_string)

	// Decrypt the data
	decrypted_data := Decrypt(key, buf.Bytes())

	fmt.Printf("decrypted data %s", decrypted_data)
}

func Encrypt(key []byte, data []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("cipher err: %v", err.Error())
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalf("cipher GCM err: %v", err.Error())
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatalf("nonce  err: %v", err.Error())
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext
}

func Decrypt(key []byte, data []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("cipher err: %v", err.Error())
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalf("cipher GCM err: %v", err.Error())
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Fatalf("cipher GCM Open err: %v", err.Error())
	}

	return plaintext
}
