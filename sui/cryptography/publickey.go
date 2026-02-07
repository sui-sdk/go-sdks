package cryptography

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/sui-sdks/go-sdks/sui/utils"
)

type PublicKey interface {
	ToRawBytes() []byte
	Flag() byte
	Verify(data, signature []byte) bool
}

func PublicKeyEquals(a, b PublicKey) bool {
	ab := a.ToRawBytes()
	bb := b.ToRawBytes()
	if len(ab) != len(bb) {
		return false
	}
	for i := range ab {
		if ab[i] != bb[i] {
			return false
		}
	}
	return true
}

func ToBase64(pk PublicKey) string {
	return base64.StdEncoding.EncodeToString(pk.ToRawBytes())
}

func ToSuiBytes(pk PublicKey) []byte {
	raw := pk.ToRawBytes()
	out := make([]byte, len(raw)+1)
	out[0] = pk.Flag()
	copy(out[1:], raw)
	return out
}

func ToSuiPublicKey(pk PublicKey) string {
	return base64.StdEncoding.EncodeToString(ToSuiBytes(pk))
}

func ToSuiAddress(pk PublicKey) string {
	digest := sha256.Sum256(ToSuiBytes(pk))
	return utils.NormalizeSuiAddress(hex.EncodeToString(digest[:])[:utils.SuiAddressLength*2])
}

func VerifyWithIntent(pk PublicKey, bytes, signature []byte, intent IntentScope) bool {
	msg := MessageWithIntent(intent, bytes)
	digest := sha256.Sum256(msg)
	return pk.Verify(digest[:], signature)
}

func VerifyPersonalMessage(pk PublicKey, msg, signature []byte) bool {
	return VerifyWithIntent(pk, msg, signature, IntentPersonalMessage)
}

func VerifyTransaction(pk PublicKey, tx, signature []byte) bool {
	return VerifyWithIntent(pk, tx, signature, IntentTransactionData)
}

type ParsedSerializedSignature struct {
	SerializedSignature string
	SignatureScheme     SignatureScheme
	Signature           []byte
	PublicKey           []byte
	Bytes               []byte
}

func ParseSerializedKeypairSignature(serialized string) (ParsedSerializedSignature, error) {
	bytes, err := base64.StdEncoding.DecodeString(serialized)
	if err != nil {
		return ParsedSerializedSignature{}, err
	}
	if len(bytes) < 2 {
		return ParsedSerializedSignature{}, fmt.Errorf("invalid signature bytes")
	}
	scheme, ok := SignatureFlagToScheme[bytes[0]]
	if !ok {
		return ParsedSerializedSignature{}, fmt.Errorf("unsupported signature scheme")
	}
	pkSize := SignatureSchemeToSize[scheme]
	if len(bytes) < 1+pkSize {
		return ParsedSerializedSignature{}, fmt.Errorf("invalid signature length")
	}
	sig := bytes[1 : len(bytes)-pkSize]
	pk := bytes[len(bytes)-pkSize:]
	return ParsedSerializedSignature{
		SerializedSignature: serialized,
		SignatureScheme:     scheme,
		Signature:           append([]byte(nil), sig...),
		PublicKey:           append([]byte(nil), pk...),
		Bytes:               bytes,
	}, nil
}
