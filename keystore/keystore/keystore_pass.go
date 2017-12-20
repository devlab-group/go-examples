package keystore

import (
	"bytes"
	"crypto/aes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/pborman/uuid"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/scrypt"
	"io/ioutil"
	"path/filepath"
	"project/crypto"
)

const (
	keyHeaderKDF = "scrypt"

	scryptN     = 1 << 18
	scryptP     = 1
	scryptR     = 8
	scryptDKLen = 32
)

func GetKey(filename, passphrase string) (*Key, error) {
	filename = filepath.Join(keysDir, filename)
	keyjson, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	key, err := DecryptKey(keyjson, passphrase)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func StoreKey(key *Key, filename, passphrase string) error {
	keyjson, err := EncryptKey(key, passphrase)
	if err != nil {
		return err
	}

	return writeKeyToFile(filename, keyjson)
}

func EncryptKey(key *Key, passphrase string) ([]byte, error) {
	passphraseArray := []byte(passphrase)
	salt := crypto.GetEntropyCSPRNG(32)
	derivedKey, err := scrypt.Key(passphraseArray, salt, scryptN, scryptR, scryptP, scryptDKLen)
	if err != nil {
		return nil, err
	}

	// For AES-128
	encryptKey := derivedKey[:16]
	iv := crypto.GetEntropyCSPRNG(aes.BlockSize)

	cipherText, err := crypto.AesCTRXOR(encryptKey, key.PrivateKey, iv)
	if err != nil {
		return nil, err
	}
	mac := crypto.Keccak256(derivedKey[16:32], cipherText)

	scryptParamsJSON := make(map[string]interface{}, 5)
	scryptParamsJSON["n"] = scryptN
	scryptParamsJSON["r"] = scryptR
	scryptParamsJSON["p"] = scryptP
	scryptParamsJSON["dklen"] = scryptDKLen
	scryptParamsJSON["salt"] = hex.EncodeToString(salt)

	cipherparamsJSON := cipherparamsJSON{
		IV: hex.EncodeToString(iv),
	}

	cryptoJSON := cryptoJSON{
		Cipher:       "aes-128-ctr",
		CipherText:   hex.EncodeToString(cipherText),
		CipherParams: cipherparamsJSON,
		KDF:          keyHeaderKDF,
		KDFParams:    scryptParamsJSON,
		MAC:          hex.EncodeToString(mac),
	}

	encryptedKeyJSON := encryptedKeyJSON{
		key.Address.Hex(),
		hex.EncodeToString(key.PublicKey[:]),
		cryptoJSON,
		key.Id.String(),
	}

	return json.Marshal(encryptedKeyJSON)
}

func DecryptKey(keyjson []byte, passphrase string) (*Key, error) {
	keyProtected := new(encryptedKeyJSON)
	if err := json.Unmarshal(keyjson, keyProtected); err != nil {
		return nil, err
	}

	mac, err := hex.DecodeString(keyProtected.Crypto.MAC)
	if err != nil {
		return nil, err
	}

	iv, err := hex.DecodeString(keyProtected.Crypto.CipherParams.IV)
	if err != nil {
		return nil, err
	}

	cipherText, err := hex.DecodeString(keyProtected.Crypto.CipherText)
	if err != nil {
		return nil, err
	}

	derivedKey, err := getKDFKey(keyProtected.Crypto, passphrase)
	if err != nil {
		return nil, err
	}

	calculatedMAC := crypto.Keccak256(derivedKey[16:32], cipherText)
	if !bytes.Equal(calculatedMAC, mac) {
		return nil, errors.New("could not decrypt key with given passphrase")
	}

	keyBytes, err := crypto.AesCTRXOR(derivedKey[:16], cipherText, iv)
	if err != nil {
		return nil, err
	}

	pubKeyDecoded, err := hex.DecodeString(keyProtected.PublicKey)
	if err != nil {
		return nil, err
	}

	addressDecoded, err := hex.DecodeString(keyProtected.Address)
	if err != nil {
		return nil, err
	}

	var address Address
	address.SetBytes(addressDecoded)

	privateKey := ed25519.PrivateKey(keyBytes)
	publicKey := ed25519.PublicKey(pubKeyDecoded)

	key := &Key{
		Id:         uuid.Parse(keyProtected.Id),
		Address:    address,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	return key, nil
}

func getKDFKey(cryptoJSON cryptoJSON, passphrase string) ([]byte, error) {
	passphraseArray := []byte(passphrase)

	salt, err := hex.DecodeString(cryptoJSON.KDFParams["salt"].(string))
	if err != nil {
		return nil, err
	}

	n := ensureInt(cryptoJSON.KDFParams["n"])
	r := ensureInt(cryptoJSON.KDFParams["r"])
	p := ensureInt(cryptoJSON.KDFParams["p"])
	dkLen := ensureInt(cryptoJSON.KDFParams["dklen"])

	return scrypt.Key(passphraseArray, salt, n, r, p, dkLen)
}

func ensureInt(x interface{}) int {
	res, ok := x.(int)
	if !ok {
		res = int(x.(float64))
	}

	return res
}
