package cryptography

import (
	"encoding/base64"
	"testing"
)

type mockPK struct{ raw []byte }

func (m mockPK) ToRawBytes() []byte { return append([]byte(nil), m.raw...) }
func (m mockPK) Flag() byte { return SignatureSchemeToFlag[SchemeED25519] }
func (m mockPK) Verify(data, signature []byte) bool { return len(data) > 0 && len(signature) > 0 }

func TestEncodeDecodeSuiPrivateKey(t *testing.T) {
	seed := make([]byte, PrivateKeySize)
	for i := range seed {
		seed[i] = byte(i)
	}
	encoded, err := EncodeSuiPrivateKey(seed, SchemeED25519)
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}
	parsed, err := DecodeSuiPrivateKey(encoded)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if parsed.Scheme != SchemeED25519 {
		t.Fatalf("scheme mismatch")
	}
	if base64.StdEncoding.EncodeToString(parsed.SecretKey) != base64.StdEncoding.EncodeToString(seed) {
		t.Fatalf("secret key mismatch")
	}
}

func TestSerializeSignature(t *testing.T) {
	pk := mockPK{raw: make([]byte, 32)}
	for i := range pk.raw {
		pk.raw[i] = byte(i)
	}
	serialized, err := ToSerializedSignature(SerializeSignatureInput{SignatureScheme: SchemeED25519, Signature: []byte{1, 2, 3}, PublicKey: pk})
	if err != nil {
		t.Fatalf("serialize failed: %v", err)
	}
	parsed, err := ParseSerializedSignature(serialized)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if parsed.SignatureScheme != SchemeED25519 {
		t.Fatalf("scheme mismatch")
	}
	if len(parsed.PublicKey) != 32 {
		t.Fatalf("public key size mismatch")
	}
}
