# ipfs-go-encrypt

If IPFS is not running run the following

```
ipfs daemon // or 'ipfs daemon &' to run in the background
```

Create a private key for encryption in a file called `key` in the file directory.
You can generate the key like this

```
openssl enc -aes-128-cbc -k secret -P -md sha1
```

Run the script with the test file or your file of choice

```
go run . upload test.txt // uploads encrypted version & logs the cuid
go run . download <cuid> // downloads & decrypts (optional --out flag for output file)
```
