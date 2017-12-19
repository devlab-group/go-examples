package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/sha3"
	"log"
	"os"
)

// Hash bytes with sha256
func Hash(bytes []byte) []byte {
	hasher := sha256.New()
	hasher.Write(bytes)
	return hasher.Sum(nil)
}

// Hash string with sha256
func HashFromString(str string) []byte {
	return Hash([]byte(str))
}

// Convert bytes to hex encoded string
func Hex(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Argument #1 is empty")
	}

	public, private, _ := ed25519.GenerateKey(bytes.NewReader(HashFromString(os.Args[1])))

	message := make([]byte, 64)
	sha3.ShakeSum256(message, []byte("Hello World"))

	sig := ed25519.Sign(private, message)

	log.Printf("Signature is %s", Hex(sig))
	log.Printf("Public Key is %s", Hex(public))

	if !ed25519.Verify(public, message, sig) {
		log.Fatal("valid signature rejected")
	}

	wrongMessage := make([]byte, 64)
	sha3.ShakeSum256(wrongMessage, []byte("Wrong Message"))

	if ed25519.Verify(public, wrongMessage, sig) {
		log.Fatal("signature of differnet message accepted")
	}

	log.Print("success")
}
