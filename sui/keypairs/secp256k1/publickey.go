package secp256k1

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"fmt"

	"github.com/sui-sdks/go-sdks/sui/cryptography"
)

const PublicKeySize = 33

type PublicKey struct{ data []byte }

func NewPublicKey(value []byte) (*PublicKey, error) {
	if len(value) != PublicKeySize {
		return nil, fmt.Errorf("invalid public key size: %d", len(value))
	}
	return &PublicKey{data: append([]byte(nil), value...)}, nil
}
func NewPublicKeyFromBase64(value string) (*PublicKey, error) {
	b, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}
	return NewPublicKey(b)
}

func (p *PublicKey) ToRawBytes() []byte { return append([]byte(nil), p.data...) }
func (p *PublicKey) Flag() byte { return cryptography.SignatureSchemeToFlag[cryptography.SchemeSecp256k1] }
func (p *PublicKey) Verify(data, signature []byte) bool {
	curve := elliptic.P256()
	x, y := elliptic.UnmarshalCompressed(curve, p.data)
	if x == nil {
		return false
	}
	pub := ecdsa.PublicKey{Curve: curve, X: x, Y: y}
	return ecdsa.VerifyASN1(&pub, data, signature)
}
func (p *PublicKey) ToSuiAddress() string { return cryptography.ToSuiAddress(p) }
