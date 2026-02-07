package cryptography

import (
	"encoding/base64"
	"fmt"
)

type SerializeSignatureInput struct {
	SignatureScheme SignatureScheme
	Signature       []byte
	PublicKey       PublicKey
}

func ToSerializedSignature(input SerializeSignatureInput) (string, error) {
	if input.PublicKey == nil {
		return "", fmt.Errorf("publicKey is required")
	}
	flag, ok := SignatureSchemeToFlag[input.SignatureScheme]
	if !ok {
		return "", fmt.Errorf("unsupported signature scheme")
	}
	pk := input.PublicKey.ToRawBytes()
	out := make([]byte, 1+len(input.Signature)+len(pk))
	out[0] = flag
	copy(out[1:], input.Signature)
	copy(out[1+len(input.Signature):], pk)
	return base64.StdEncoding.EncodeToString(out), nil
}

func ParseSerializedSignature(serialized string) (ParsedSerializedSignature, error) {
	return ParseSerializedKeypairSignature(serialized)
}
