package keystore

import (
	crand "crypto/rand"
	"encoding/hex"
	"github.com/pborman/uuid"
	"golang.org/x/crypto/ed25519"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"project/crypto"
)

const keysDir = "keys"

type Key struct {
	Id         uuid.UUID
	Address    Address
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
}

type cryptoJSON struct {
	Cipher       string                 `json:"cipher"`
	CipherText   string                 `json:"ciphertext"`
	CipherParams cipherparamsJSON       `json:"cipherparams"`
	KDF          string                 `json:"kdf"`
	KDFParams    map[string]interface{} `json:"kdfparams"`
	MAC          string                 `json:"mac"`
}

type cipherparamsJSON struct {
	IV string `json:"iv"`
}

type encryptedKeyJSON struct {
	Address   string     `json:"address"`
	PublicKey string     `json:"pub_key"`
	Crypto    cryptoJSON `json:"crypto"`
	Id        string     `json:"id"`
}

func PubkeyToAddress(p ed25519.PublicKey) Address {
	pubkeyHash := crypto.Keccak256(p[1:])[12:]
	var a Address
	a.SetBytes(pubkeyHash)
	return a
}

func NewKey(passphrase string) (Address, error) {
	key, err := storeNewKey(crand.Reader, passphrase)
	if err != nil {
		return Address{}, err
	}
	return key.Address, nil
}

func newKey(rand io.Reader) (*Key, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand)
	log.Println(hex.EncodeToString(publicKey[:]))
	if err != nil {
		return nil, err
	}
	address := PubkeyToAddress(publicKey)
	id := uuid.NewRandom()
	key := &Key{
		Id:         id,
		Address:    address,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	return key, nil
}

func storeNewKey(rand io.Reader, passphrase string) (*Key, error) {
	key, err := newKey(rand)
	if err != nil {
		return nil, err
	}

	err = StoreKey(key, keyFileName(key.Address), passphrase)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func writeKeyToFile(filename string, content []byte) error {
	filename = filepath.Join(keysDir, filename)

	const dirPerm = 0755
	if err := os.MkdirAll(filepath.Dir(filename), dirPerm); err != nil {
		return err
	}

	f, err := ioutil.TempFile(filepath.Dir(filename), "."+filepath.Base(filename)+".tmp")
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err := f.Write(content); err != nil {
		os.Remove(f.Name())
		return err
	}

	return os.Rename(f.Name(), filename)
}

// Set key file name as key address for a while
func keyFileName(keyAddr Address) string {
	return keyAddr.Hex()
}
