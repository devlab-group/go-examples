package keystore

import (
	"encoding/hex"
)

const AddressLength = 20

type Address [AddressLength]byte

// Sets the address to the value of b. If b is larger than len(a) it will panic
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

func (a Address) Hex() string {
	return hex.EncodeToString(a[:])
}
