package seal

import (
	"crypto/rand"
	"fmt"
)

type KemType int

const (
	KemTypeBonehFranklinBLS12381DemCCA KemType = iota
)

type DemType int

const (
	DemTypeAesGcm256 DemType = iota
	DemTypeHmac256Ctr
)

func Encrypt(input struct {
	KeyServers      []KeyServer
	KEMType         KemType
	Threshold       int
	PackageID       string
	ID              string
	EncryptionInput EncryptionInput
}) (encryptedObject []byte, key []byte, err error) {
	if input.Threshold <= 0 || input.Threshold >= MaxU8 || len(input.KeyServers) < input.Threshold || len(input.KeyServers) >= MaxU8 {
		return nil, nil, (&UserError{SealError{Msg: fmt.Sprintf("invalid threshold %d for %d key servers", input.Threshold, len(input.KeyServers))}})
	}
	baseKey, err := input.EncryptionInput.GenerateKey()
	if err != nil {
		return nil, nil, err
	}
	shares, err := Split(baseKey, input.Threshold, len(input.KeyServers))
	if err != nil {
		return nil, nil, err
	}
	fullID := CreateFullID(input.PackageID, input.ID)
	fullIDBytes := []byte(fullID)
	encryptedShares := make([][]byte, len(shares))
	services := make([][2]any, len(shares))
	for i, sh := range shares {
		mask := deriveDigest(input.KeyServers[i].PK, fullIDBytes, []byte{byte(sh.Index)})
		encryptedShares[i] = xorUnchecked(sh.Share, mask[:len(sh.Share)])
		services[i] = [2]any{input.KeyServers[i].ObjectID, sh.Index}
	}

	nonce := make([]byte, 32)
	if _, err := rand.Read(nonce); err != nil {
		return nil, nil, err
	}
	randomness := deriveDigest(baseKey, nonce)
	randomnessKey := DeriveKey(KeyPurposeEncryptedRandomness, baseKey, encryptedShares, input.Threshold, extractObjectIDs(input.KeyServers))
	encryptedRandomness := xorUnchecked(randomness, randomnessKey[:len(randomness)])
	demKey := DeriveKey(KeyPurposeDEM, baseKey, encryptedShares, input.Threshold, extractObjectIDs(input.KeyServers))
	ciphertext, err := input.EncryptionInput.Encrypt(demKey)
	if err != nil {
		return nil, nil, err
	}
	obj := EncryptedObject{
		Version:   0,
		PackageID: input.PackageID,
		ID:        input.ID,
		Services:  services,
		Threshold: input.Threshold,
		EncryptedShares: IBEEncryptions{BonehFranklinBLS12381: &BonehFranklinEncryptedShares{
			EncryptedShares:     encryptedShares,
			EncryptedRandomness: encryptedRandomness,
			Nonce:               nonce,
		}},
		Ciphertext: ciphertext,
	}
	bytes, err := SerializeEncryptedObject(obj)
	if err != nil {
		return nil, nil, err
	}
	return bytes, demKey, nil
}

func extractObjectIDs(servers []KeyServer) []string {
	out := make([]string, len(servers))
	for i, s := range servers {
		out[i] = s.ObjectID
	}
	return out
}
