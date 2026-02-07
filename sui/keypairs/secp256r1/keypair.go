package secp256r1

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/sui-sdks/go-sdks/sui/cryptography"
)

const DefaultDerivationPath = "m/74'/784'/0'/0/0"

type Keypair struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *PublicKey
}

func Generate() (*Keypair, error) {
	prv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	comp := elliptic.MarshalCompressed(elliptic.P256(), prv.PublicKey.X, prv.PublicKey.Y)
	pk, _ := NewPublicKey(comp)
	return &Keypair{privateKey: prv, publicKey: pk}, nil
}

func FromSeed(seed []byte) (*Keypair, error) {
	if len(seed) != cryptography.PrivateKeySize {
		return nil, fmt.Errorf("wrong secret key size: %d", len(seed))
	}
	curve := elliptic.P256()
	d := new(big.Int).SetBytes(seed)
	n := new(big.Int).Sub(curve.Params().N, big.NewInt(1))
	d.Mod(d, n)
	d.Add(d, big.NewInt(1))
	x, y := curve.ScalarBaseMult(d.Bytes())
	prv := &ecdsa.PrivateKey{PublicKey: ecdsa.PublicKey{Curve: curve, X: x, Y: y}, D: d}
	comp := elliptic.MarshalCompressed(curve, x, y)
	pk, _ := NewPublicKey(comp)
	return &Keypair{privateKey: prv, publicKey: pk}, nil
}

func FromSecretKey(secretKey []byte) (*Keypair, error) { return FromSeed(secretKey) }

func FromSecretKeyString(secretKey string) (*Keypair, error) {
	decoded, err := cryptography.DecodeSuiPrivateKey(secretKey)
	if err != nil {
		return nil, err
	}
	if decoded.Scheme != cryptography.SchemeSecp256r1 {
		return nil, fmt.Errorf("expected Secp256r1 keypair, got %s", decoded.Scheme)
	}
	return FromSecretKey(decoded.SecretKey)
}

func (k *Keypair) Sign(bytes []byte) ([]byte, error) { return ecdsa.SignASN1(rand.Reader, k.privateKey, bytes) }
func (k *Keypair) GetKeyScheme() cryptography.SignatureScheme { return cryptography.SchemeSecp256r1 }
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
	b := k.privateKey.D.Bytes()
	seed := make([]byte, cryptography.PrivateKeySize)
	copy(seed[cryptography.PrivateKeySize-len(b):], b)
	v, _ := cryptography.EncodeSuiPrivateKey(seed, k.GetKeyScheme())
	return v
}
