package main

import (
	"encoding/hex"
	"golang.org/x/crypto/ed25519"
	"log"
	"os"
	"project/keystore"
)

func Hex(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

func main() {
	if len(os.Args) < 1 {
		log.Fatal("Action not specified")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "new":
		if len(os.Args) < 3 {
			log.Fatal("Passphrase is missing")
			os.Exit(1)
		}
		address, err := keystore.NewKey(os.Args[2])
		checkErr(err)

		log.Println(address.Hex())
	case "get-key":
		if len(os.Args) < 4 {
			log.Fatal("Usage: unlock ADDRESS PASSPHRASE")
			os.Exit(1)
		}
		key, err := keystore.GetKey(os.Args[2], os.Args[3])
		checkErr(err)

		log.Println(key.Address.Hex())
		log.Println(Hex(key.PublicKey[:]))
	case "sign":
		if len(os.Args) < 5 {
			log.Fatal("Usage: sign ADDRESS PASSPHRASE MESSAGE")
			os.Exit(1)
		}
		key, err := keystore.GetKey(os.Args[2], os.Args[3])
		checkErr(err)

		message := []byte(os.Args[4])
		sig := ed25519.Sign(key.PrivateKey, message)

		log.Printf("Signature is %s", Hex(sig))
		log.Printf("Public Key is %s", Hex(key.PublicKey))
	case "verify":
		message := []byte(os.Args[2])

		signatureDecoded, _ := hex.DecodeString(os.Args[3])
		signature := []byte(signatureDecoded)

		pub, _ := hex.DecodeString(os.Args[4])
		pubKey := ed25519.PublicKey(pub)

		if !ed25519.Verify(pubKey, message, signature) {
			log.Fatal("verification failed")
		}
		log.Println("correct")
	default:
		log.Fatal("Unknown action:", os.Args[1])
		os.Exit(1)
	}
}

func checkErr(err interface{}) {
	if err != nil {
		panic(err)
	}
}
