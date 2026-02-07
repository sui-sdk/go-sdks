package zklogin

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestJWTDecodeAndNonce(t *testing.T) {
	head := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"iss":"https://example.com","sub":"u1"}`))
	jwt := strings.Join([]string{head, payload, "sig"}, ".")
	decoded, err := DecodeJWT(jwt)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if decoded["sub"].(string) != "u1" {
		t.Fatalf("unexpected payload")
	}
	rand, _ := GenerateRandomness()
	nonce, err := GenerateNonce([]byte{1, 2, 3}, 10, rand)
	if err != nil || nonce == "" {
		t.Fatalf("generate nonce failed")
	}
}

func TestZkLoginSignatureRoundTrip(t *testing.T) {
	sig, err := GetZkLoginSignature(ZkLoginSignatureExtended{Inputs: map[string]any{"a": 1}, MaxEpoch: 10, UserSignature: "abc"})
	if err != nil {
		t.Fatalf("get signature failed: %v", err)
	}
	parsed, err := ParseZkLoginSignature(sig)
	if err != nil || parsed.UserSignature != "abc" {
		t.Fatalf("parse signature failed")
	}
}
