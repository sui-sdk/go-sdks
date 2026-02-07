package seal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
)

type KeyPurpose uint8

const (
	KeyPurposeDEM KeyPurpose = iota
	KeyPurposeEncryptedRandomness
	KeyPurposeShareMask
)

func DeriveKey(purpose KeyPurpose, baseKey []byte, encryptedShares [][]byte, threshold int, services []string) []byte {
	h := hmac.New(sha256.New, baseKey)
	h.Write([]byte{byte(purpose)})
	for _, s := range encryptedShares {
		h.Write(s)
	}
	var th [4]byte
	binary.BigEndian.PutUint32(th[:], uint32(threshold))
	h.Write(th[:])
	SortStrings(services)
	for _, s := range services {
		h.Write([]byte(s))
	}
	return h.Sum(nil)
}
