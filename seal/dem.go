package seal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

type EncryptionInput interface {
	GenerateKey() ([]byte, error)
	Encrypt(key []byte) (Ciphertext, error)
}

type AesGcm256 struct {
	Data []byte
	AAD  []byte
}

func (a AesGcm256) GenerateKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	return key, err
}

func (a AesGcm256) Encrypt(key []byte) (Ciphertext, error) {
	block, err := aes.NewCipher(key[:32])
	if err != nil {
		return Ciphertext{}, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return Ciphertext{}, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return Ciphertext{}, err
	}
	ct := gcm.Seal(nil, nonce, a.Data, a.AAD)
	return Ciphertext{Aes256Gcm: &Aes256GcmCiphertext{Nonce: nonce, Ciphertext: ct, AAD: a.AAD}}, nil
}

func AesGcmDecrypt(key []byte, ct *Aes256GcmCiphertext) ([]byte, error) {
	block, err := aes.NewCipher(key[:32])
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return gcm.Open(nil, ct.Nonce, ct.Ciphertext, ct.AAD)
}

type Hmac256Ctr struct {
	Data []byte
	AAD  []byte
}

func (h Hmac256Ctr) GenerateKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	return key, err
}

func (h Hmac256Ctr) Encrypt(key []byte) (Ciphertext, error) {
	block, err := aes.NewCipher(key[:32])
	if err != nil {
		return Ciphertext{}, err
	}
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return Ciphertext{}, err
	}
	stream := cipher.NewCTR(block, iv)
	enc := make([]byte, len(h.Data))
	stream.XORKeyStream(enc, h.Data)
	mac := hmac.New(sha256.New, key[:32])
	mac.Write(iv)
	mac.Write(h.AAD)
	mac.Write(enc)
	tag := mac.Sum(nil)
	return Ciphertext{Hmac256Ctr: &Hmac256CtrCiphertext{IV: iv, Ciphertext: enc, AAD: h.AAD, MAC: tag}}, nil
}

func HmacCtrDecrypt(key []byte, ct *Hmac256CtrCiphertext) ([]byte, error) {
	mac := hmac.New(sha256.New, key[:32])
	mac.Write(ct.IV)
	mac.Write(ct.AAD)
	mac.Write(ct.Ciphertext)
	if !hmac.Equal(mac.Sum(nil), ct.MAC) {
		return nil, fmt.Errorf("invalid mac")
	}
	block, err := aes.NewCipher(key[:32])
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(block, ct.IV)
	plain := make([]byte, len(ct.Ciphertext))
	stream.XORKeyStream(plain, ct.Ciphertext)
	return plain, nil
}
