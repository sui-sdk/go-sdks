package seal

import "fmt"

func Decrypt(opts struct {
	EncryptedObject     EncryptedObject
	Keys                map[KeyCacheKey][]byte
	CheckLEEncoding     bool
}) ([]byte, error) {
	enc := opts.EncryptedObject.EncryptedShares.BonehFranklinBLS12381
	if enc == nil {
		return nil, (&UnsupportedFeatureError{UserError{SealError{Msg: "encryption mode not supported"}}})
	}
	fullID := CreateFullID(opts.EncryptedObject.PackageID, opts.EncryptedObject.ID)
	shares := make([]Share, 0, len(opts.EncryptedObject.Services))
	services := make([]string, 0, len(opts.EncryptedObject.Services))
	for i, srv := range opts.EncryptedObject.Services {
		objectID, _ := srv[0].(string)
		idx, _ := toInt(srv[1])
		services = append(services, objectID)
		k, ok := opts.Keys[KeyCacheKey(fullID+":"+objectID)]
		if !ok {
			continue
		}
		if i >= len(enc.EncryptedShares) {
			return nil, (&InvalidCiphertextError{UserError{SealError{Msg: "mismatched share count"}}})
		}
		share := xorUnchecked(enc.EncryptedShares[i], k[:len(enc.EncryptedShares[i])])
		shares = append(shares, Share{Index: idx, Share: share})
	}
	if len(shares) < opts.EncryptedObject.Threshold {
		return nil, fmt.Errorf("not enough shares, please fetch more keys")
	}
	baseKey, err := Combine(shares[:opts.EncryptedObject.Threshold])
	if err != nil {
		return nil, err
	}
	demKey := DeriveKey(KeyPurposeDEM, baseKey, enc.EncryptedShares, opts.EncryptedObject.Threshold, services)
	if opts.EncryptedObject.Ciphertext.Aes256Gcm != nil {
		return AesGcmDecrypt(demKey, opts.EncryptedObject.Ciphertext.Aes256Gcm)
	}
	if opts.EncryptedObject.Ciphertext.Hmac256Ctr != nil {
		return HmacCtrDecrypt(demKey, opts.EncryptedObject.Ciphertext.Hmac256Ctr)
	}
	return nil, (&InvalidCiphertextError{UserError{SealError{Msg: "invalid ciphertext type"}}})
}

func toInt(v any) (int, bool) {
	switch t := v.(type) {
	case int:
		return t, true
	case float64:
		return int(t), true
	default:
		return 0, false
	}
}
