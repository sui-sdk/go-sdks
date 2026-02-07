package ed25519

import (
	cryptoed25519 "crypto/ed25519"
	"crypto/rand"
	"fmt"

	"github.com/sui-sdks/go-sdks/sui/cryptography"
)

const DefaultDerivationPath = "m/44'/784'/0'/0'/0'"

type Keypair struct {
	secretKey cryptoed25519.PrivateKey
	publicKey *PublicKey
}

func Generate() (*Keypair, error) {
	pub, prv, err := cryptoed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	pk, _ := NewPublicKey(pub)
	return &Keypair{secretKey: prv, publicKey: pk}, nil
}

func FromSecretKey(secretKey []byte) (*Keypair, error) {
	if len(secretKey) != cryptography.PrivateKeySize {
		return nil, fmt.Errorf("wrong secret key size: %d", len(secretKey))
	}
	prv := cryptoed25519.NewKeyFromSeed(secretKey)
	pk, _ := NewPublicKey(prv[32:])
	return &Keypair{secretKey: prv, publicKey: pk}, nil
}

func FromSecretKeyString(secretKey string) (*Keypair, error) {
	decoded, err := cryptography.DecodeSuiPrivateKey(secretKey)
	if err != nil {
		return nil, err
	}
	if decoded.Scheme != cryptography.SchemeED25519 {
		return nil, fmt.Errorf("expected ED25519 keypair, got %s", decoded.Scheme)
	}
	return FromSecretKey(decoded.SecretKey)
}

func (k *Keypair) Sign(bytes []byte) ([]byte, error) { return cryptoed25519.Sign(k.secretKey, bytes), nil }
func (k *Keypair) GetKeyScheme() cryptography.SignatureScheme { return cryptography.SchemeED25519 }
func (k *Keypair) GetPublicKey() cryptography.PublicKey { return k.publicKey }
func (k *Keypair) ToSuiAddress() string { return k.publicKey.ToSuiAddress() }
func (k *Keypair) SignWithIntent(bytes []byte, intent cryptography.IntentScope) (cryptography.SignatureWithBytes, error) {
	return cryptography.SignWithIntent(k, bytes, intent)
}
func (k *Keypair) SignTransaction(bytes []byte) (cryptography.SignatureWithBytes, error) {
	return cryptography.SignTransaction(k, bytes)
}
func (k *Keypair) SignPersonalMessage(bytes []byte) (cryptography.SignatureWithBytes, error) {
	return cryptography.SignPersonalMessage(k, bytes)
}
func (k *Keypair) GetSecretKey() string {
	seed := k.secretKey.Seed()
	v, _ := cryptography.EncodeSuiPrivateKey(seed, k.GetKeyScheme())
	return v
}
