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

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/urfave/cli/v2"
)

func main() {

	// Where your local node is running on localhost:5001

	sh := shell.NewShell("localhost:5001")

	app := &cli.App{
		Name:  "ipfs-go-encrypt",
		Usage: "Upload and download encrypted files from IPFS",
		Commands: []*cli.Command{
			{
				Name:    "upload",
				Aliases: []string{"a"},
				Usage:   "upload a file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "key",
						Value:    "key",
						Usage:    "encyrption key (use aes-128)",
						Required: true,
					},
				},
				Action: func(cCtx *cli.Context) error {
					// Get the file
					file, err := os.ReadFile(cCtx.Args().First())
					if err != nil {
						log.Fatal(err)
					}

					key := []byte(cCtx.String("key"))

					// Encrypt the file
					encryptedFile := Encrypt(key, file)

					// Upload the file
					cid, err := sh.Add(bytes.NewReader(encryptedFile))

					if err != nil {
						log.Fatal(err)
					}

					fmt.Println(cid)

					return nil
				},
			},
			{
				Name:    "download",
				Aliases: []string{"d"},
				Usage:   "download a file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "out",
						Value: "out",
						Usage: "output file",
					},
					&cli.StringFlag{
						Name:     "key",
						Value:    "key",
						Usage:    "encyrption key (use aes-128)",
						Required: true,
					},
				},

				Action: func(cCtx *cli.Context) error {
					data, err := sh.Cat(cCtx.Args().First())
					if err != nil {
						fmt.Fprintf(os.Stderr, "error: %s", err)
						os.Exit(1)
					}

					// Use a buffer to read the data from the reader
					buf := new(bytes.Buffer)
					buf.ReadFrom(data)

					key := []byte(cCtx.String("key"))

					// Decrypt the data
					decrypted := Decrypt(key, buf.Bytes())

					// Write the data to the output file
					err = os.WriteFile(cCtx.String("out"), decrypted, 0644)

					if err != nil {
						fmt.Fprintf(os.Stderr, "error: %s", err)
						os.Exit(1)
					}

					return nil
				},
			},
			{
				Name:    "keygen",
				Aliases: []string{"k"},
				Usage:   "generate a key",
				Action: func(cCtx *cli.Context) error {
					key := make([]byte, 16)
					_, err := rand.Read(key)
					if err != nil {
						log.Fatal(err)
					}

					fmt.Println(string(key))

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
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
