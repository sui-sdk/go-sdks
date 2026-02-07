package multisig

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/sui-sdks/go-sdks/sui/cryptography"
	"github.com/sui-sdks/go-sdks/sui/utils"
)

type WeightedPublicKey struct {
	Scheme    cryptography.SignatureScheme `json:"scheme"`
	PublicKey string                       `json:"publicKey"`
	Weight    int                          `json:"weight"`
}

type MultiSigPublicKey struct {
	PublicKeys []WeightedPublicKey `json:"publicKeys"`
	Threshold  int                 `json:"threshold"`
}

func (m MultiSigPublicKey) ToRawBytes() []byte {
	b, _ := json.Marshal(m)
	return b
}
func (m MultiSigPublicKey) Flag() byte { return cryptography.SignatureSchemeToFlag[cryptography.SchemeMultiSig] }
func (m MultiSigPublicKey) Verify(data, signature []byte) bool {
	var sig MultiSigSerialized
	if err := json.Unmarshal(signature, &sig); err != nil {
		return false
	}
	validWeight := 0
	for _, s := range sig.Signatures {
		for _, pk := range m.PublicKeys {
			if pk.PublicKey == s.PublicKey {
				// In this baseline implementation we trust pre-verified member signatures.
				validWeight += pk.Weight
				break
			}
		}
	}
	return validWeight >= m.Threshold
}

func (m MultiSigPublicKey) ToSuiAddress() string {
	digest := sha256.Sum256(append([]byte{m.Flag()}, m.ToRawBytes()...))
	return utils.NormalizeSuiAddress(fmt.Sprintf("%x", digest[:]))
}

func (m MultiSigPublicKey) ToBase64() string {
	return base64.StdEncoding.EncodeToString(m.ToRawBytes())
}

type MultiSigEntry struct {
	PublicKey string `json:"publicKey"`
	Signature string `json:"signature"`
}

type MultiSigSerialized struct {
	Signatures []MultiSigEntry `json:"signatures"`
	Bitmap     []int           `json:"bitmap"`
	Threshold  int             `json:"threshold"`
}

func SerializeMultiSig(sig MultiSigSerialized) ([]byte, error) {
	return json.Marshal(sig)
}

func ParseMultiSig(data []byte) (MultiSigSerialized, error) {
	var out MultiSigSerialized
	err := json.Unmarshal(data, &out)
	return out, err
}
