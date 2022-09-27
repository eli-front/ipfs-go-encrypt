# ipfs-go-encrypt

A simple CLI written in Go to both encrypt + upload files to `ipfs` and decrypt + download them.

## Install

1. Run `go install` inside the working dir
2. Make sure the go bin directory is added to your path. You can check this by restarting your terminal and running `ipfs-go-encrypt`. If it's added to your path then you will see a summary of the tool.

## Using the Tool

If IPFS is not running run the following

```
ipfs daemon // or 'ipfs daemon &' to run in the background
```

Create a private key for encryption
You can generate the key like this

```
openssl enc -aes-128-cbc -k secret -P -md sha1
// or
ipfs-go-encrypt keygen
```

Run the script with the test file or your file of choice

```
ipfs-go-encrypt upload --key=<private key> test.txt // uploads encrypted version & logs the cuid
ipfs-go-encrypt download --key<private key> <cuid> // downloads & decrypts (optional --out flag for output file)
```
