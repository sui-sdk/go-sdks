package verify

import (
	"testing"

	edkp "github.com/sui-sdks/go-sdks/sui/keypairs/ed25519"
)

func TestVerifyHelpers(t *testing.T) {
	kp, err := edkp.Generate()
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}
	msg := []byte("hello")
	sig, _ := kp.Sign(msg)
	if !VerifySignature(kp.GetPublicKey(), msg, sig) {
		t.Fatalf("verify signature failed")
	}
}
