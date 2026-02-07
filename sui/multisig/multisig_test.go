package multisig

import (
	"encoding/base64"
	"testing"

	"github.com/sui-sdks/go-sdks/sui/cryptography"
	edkp "github.com/sui-sdks/go-sdks/sui/keypairs/ed25519"
)

type partial struct{ inner *edkp.Keypair }

func (p partial) Sign(data []byte) ([]byte, error) { return p.inner.Sign(data) }
func (p partial) GetPublicKey() cryptography.PublicKey { return p.inner.GetPublicKey() }

func TestMultiSigSerializeParse(t *testing.T) {
	k1, _ := edkp.Generate()
	k2, _ := edkp.Generate()
	ms := MultiSigPublicKey{
		PublicKeys: []WeightedPublicKey{
			{Scheme: cryptography.SchemeED25519, PublicKey: base64.StdEncoding.EncodeToString(k1.GetPublicKey().ToRawBytes()), Weight: 1},
			{Scheme: cryptography.SchemeED25519, PublicKey: base64.StdEncoding.EncodeToString(k2.GetPublicKey().ToRawBytes()), Weight: 1},
		},
		Threshold: 2,
	}
	signer := Signer{PublicKey: ms, Signers: []PartialSigner{partial{k1}, partial{k2}}}
	sig, err := signer.Sign([]byte("data"))
	if err != nil {
		t.Fatalf("sign failed: %v", err)
	}
	parsed, err := ParseMultiSig(sig)
	if err != nil || len(parsed.Signatures) != 2 {
		t.Fatalf("parse failed")
	}
	if !ms.Verify([]byte("data"), sig) {
		t.Fatalf("verify failed")
	}
}
