package seal

import "encoding/json"

type IBEEncryptions struct {
	BonehFranklinBLS12381 *BonehFranklinEncryptedShares `json:"BonehFranklinBLS12381,omitempty"`
}

type BonehFranklinEncryptedShares struct {
	EncryptedShares     [][]byte `json:"encryptedShares"`
	EncryptedRandomness []byte   `json:"encryptedRandomness"`
	Nonce               []byte   `json:"nonce"`
}

type Aes256GcmCiphertext struct {
	Nonce      []byte `json:"nonce"`
	Ciphertext []byte `json:"ciphertext"`
	AAD        []byte `json:"aad,omitempty"`
}

type Hmac256CtrCiphertext struct {
	IV         []byte `json:"iv"`
	Ciphertext []byte `json:"ciphertext"`
	AAD        []byte `json:"aad,omitempty"`
	MAC        []byte `json:"mac"`
}

type Ciphertext struct {
	Aes256Gcm *Aes256GcmCiphertext `json:"Aes256Gcm,omitempty"`
	Hmac256Ctr *Hmac256CtrCiphertext `json:"Hmac256Ctr,omitempty"`
}

type EncryptedObject struct {
	Version         uint8          `json:"version"`
	PackageID       string         `json:"packageId"`
	ID              string         `json:"id"`
	Services        [][2]any       `json:"services"`
	Threshold       int            `json:"threshold"`
	EncryptedShares IBEEncryptions `json:"encryptedShares"`
	Ciphertext      Ciphertext     `json:"ciphertext"`
}

func SerializeEncryptedObject(obj EncryptedObject) ([]byte, error) {
	return json.Marshal(obj)
}

func ParseEncryptedObject(data []byte) (EncryptedObject, error) {
	var obj EncryptedObject
	err := json.Unmarshal(data, &obj)
	return obj, err
}
