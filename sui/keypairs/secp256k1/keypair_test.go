package secp256k1

import "testing"

func TestSecp256k1SignVerify(t *testing.T) {
	kp, err := Generate()
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}
	digest := []byte("0123456789abcdef0123456789abcdef")
	sig, err := kp.Sign(digest)
	if err != nil {
		t.Fatalf("sign failed: %v", err)
	}
	if !kp.GetPublicKey().Verify(digest, sig) {
		t.Fatalf("verify failed")
	}
}
