package main
import (
  "os"
  "log"
  "bytes"
  "encoding/hex"
  "crypto/sha256"
  "golang.org/x/crypto/ed25519"
)

func Hash(bytes []byte) ([]byte) {
  hasher := sha256.New()
  hasher.Write(bytes)
  return hasher.Sum(nil)
}

func HashString(str string) ([]byte)  {
  return Hash([]byte(str))
}

func Hex(bytes []byte) (string) {
  return hex.EncodeToString(bytes)
}

func main()  {
  if len(os.Args) < 2 {
    log.Fatal("Argument #1 is empty")
  }

  public, private, _ := ed25519.GenerateKey(bytes.NewReader(HashString(os.Args[1])))

  message := []byte("Hello World")

  sig := ed25519.Sign(private, message)

  log.Printf("Signature is %s", Hex(sig))
  log.Printf("Public Key is %s", Hex(public))

  if ! ed25519.Verify(public, message, sig) {
    log.Fatal("valid signature rejected")
  }

  wrongMessage := []byte("Wrong Message")

  if ed25519.Verify(public, wrongMessage, sig) {
    log.Fatal("signature of differnet message accepted")
  }

  log.Print("success")
}
