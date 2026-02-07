package ed25519

import (
	cryptoed25519 "crypto/ed25519"
	"encoding/base64"
	"fmt"

	"github.com/sui-sdks/go-sdks/sui/cryptography"
)

const PublicKeySize = 32

type PublicKey struct {
	data []byte
}

func NewPublicKey(value []byte) (*PublicKey, error) {
	if len(value) != PublicKeySize {
		return nil, fmt.Errorf("invalid public key size: %d", len(value))
	}
	out := append([]byte(nil), value...)
	return &PublicKey{data: out}, nil
}

func NewPublicKeyFromBase64(value string) (*PublicKey, error) {
	b, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}
	return NewPublicKey(b)
}

func (p *PublicKey) ToRawBytes() []byte { return append([]byte(nil), p.data...) }
func (p *PublicKey) Flag() byte { return cryptography.SignatureSchemeToFlag[cryptography.SchemeED25519] }
func (p *PublicKey) Verify(data, signature []byte) bool {
	return cryptoed25519.Verify(cryptoed25519.PublicKey(p.data), data, signature)
}
func (p *PublicKey) ToSuiAddress() string { return cryptography.ToSuiAddress(p) }
