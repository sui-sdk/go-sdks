package ed25519

import (
	"testing"
)

func TestEd25519SignVerify(t *testing.T) {
	kp, err := Generate()
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}
	msg := []byte("hello")
	sig, err := kp.Sign(msg)
	if err != nil {
		t.Fatalf("sign failed: %v", err)
	}
	if !kp.GetPublicKey().Verify(msg, sig) {
		t.Fatalf("verify failed")
	}
	if kp.ToSuiAddress() == "" {
		t.Fatalf("expected address")
	}
	encoded := kp.GetSecretKey()
	restored, err := FromSecretKeyString(encoded)
	if err != nil {
		t.Fatalf("restore failed: %v", err)
	}
	sig2, _ := restored.Sign(msg)
	if !restored.GetPublicKey().Verify(msg, sig2) {
		t.Fatalf("restored verify failed")
	}
}
