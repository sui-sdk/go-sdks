package cryptography

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

const (
	PrivateKeySize       = 32
	LegacyPrivateKeySize = 64
	SuiPrivateKeyPrefix  = "suiprivkey"
)

type ParsedKeypair struct {
	Scheme    SignatureScheme
	SecretKey []byte
}

type SignatureWithBytes struct {
	Bytes     string `json:"bytes"`
	Signature string `json:"signature"`
}

type Signer interface {
	Sign(bytes []byte) ([]byte, error)
	GetKeyScheme() SignatureScheme
	GetPublicKey() PublicKey
	ToSuiAddress() string
	SignWithIntent(bytes []byte, intent IntentScope) (SignatureWithBytes, error)
	SignTransaction(bytes []byte) (SignatureWithBytes, error)
	SignPersonalMessage(bytes []byte) (SignatureWithBytes, error)
}

type Keypair interface {
	Signer
	GetSecretKey() string
}

func SignWithIntent(signer Signer, bytes []byte, intent IntentScope) (SignatureWithBytes, error) {
	intentMessage := MessageWithIntent(intent, bytes)
	digest := sha256.Sum256(intentMessage)
	sig, err := signer.Sign(digest[:])
	if err != nil {
		return SignatureWithBytes{}, err
	}
	serialized, err := ToSerializedSignature(SerializeSignatureInput{
		SignatureScheme: signer.GetKeyScheme(),
		Signature:       sig,
		PublicKey:       signer.GetPublicKey(),
	})
	if err != nil {
		return SignatureWithBytes{}, err
	}
	return SignatureWithBytes{Bytes: base64.StdEncoding.EncodeToString(bytes), Signature: serialized}, nil
}

func SignTransaction(signer Signer, bytes []byte) (SignatureWithBytes, error) {
	return SignWithIntent(signer, bytes, IntentTransactionData)
}

func SignPersonalMessage(signer Signer, bytes []byte) (SignatureWithBytes, error) {
	return SignWithIntent(signer, bytes, IntentPersonalMessage)
}

func DecodeSuiPrivateKey(value string) (ParsedKeypair, error) {
	if !strings.HasPrefix(value, SuiPrivateKeyPrefix+":") {
		return ParsedKeypair{}, fmt.Errorf("invalid private key prefix")
	}
	payload := strings.TrimPrefix(value, SuiPrivateKeyPrefix+":")
	raw, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return ParsedKeypair{}, err
	}
	if len(raw) < 1 {
		return ParsedKeypair{}, fmt.Errorf("invalid private key bytes")
	}
	scheme, ok := SignatureFlagToScheme[raw[0]]
	if !ok {
		return ParsedKeypair{}, fmt.Errorf("unknown signature scheme flag")
	}
	return ParsedKeypair{Scheme: scheme, SecretKey: append([]byte(nil), raw[1:]...)}, nil
}

func EncodeSuiPrivateKey(bytes []byte, scheme SignatureScheme) (string, error) {
	if len(bytes) != PrivateKeySize {
		return "", fmt.Errorf("invalid bytes length")
	}
	flag, ok := SignatureSchemeToFlag[scheme]
	if !ok {
		return "", fmt.Errorf("unsupported signature scheme")
	}
	payload := make([]byte, 1+len(bytes))
	payload[0] = flag
	copy(payload[1:], bytes)
	return SuiPrivateKeyPrefix + ":" + base64.StdEncoding.EncodeToString(payload), nil
}
