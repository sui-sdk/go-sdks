package seal

import (
	"bytes"
	"testing"
	"time"
)

func TestShamirSplitCombine(t *testing.T) {
	secret := []byte("0123456789abcdef0123456789abcdef")
	shares, err := Split(secret, 3, 5)
	if err != nil {
		t.Fatalf("split failed: %v", err)
	}
	recovered, err := Combine([]Share{shares[0], shares[2], shares[4]})
	if err != nil {
		t.Fatalf("combine failed: %v", err)
	}
	if !bytes.Equal(secret, recovered) {
		t.Fatalf("recovered secret mismatch")
	}
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	servers := []KeyServer{
		{ObjectID: "s1", PK: bytes.Repeat([]byte{1}, 32)},
		{ObjectID: "s2", PK: bytes.Repeat([]byte{2}, 32)},
		{ObjectID: "s3", PK: bytes.Repeat([]byte{3}, 32)},
	}
	plain := []byte("hello seal")
	enc, _, err := Encrypt(struct {
		KeyServers      []KeyServer
		KEMType         KemType
		Threshold       int
		PackageID       string
		ID              string
		EncryptionInput EncryptionInput
	}{
		KeyServers:      servers,
		KEMType:         KemTypeBonehFranklinBLS12381DemCCA,
		Threshold:       2,
		PackageID:       "0x2",
		ID:              "obj-1",
		EncryptionInput: AesGcm256{Data: plain},
	})
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}
	obj, err := ParseEncryptedObject(enc)
	if err != nil {
		t.Fatalf("parse encrypted object failed: %v", err)
	}

	fullID := CreateFullID(obj.PackageID, obj.ID)
	keys := map[KeyCacheKey][]byte{}
	for i, srv := range servers {
		idx := obj.Services[i][1].(float64)
		mask := deriveDigest(srv.PK, []byte(fullID), []byte{byte(int(idx))})
		keys[KeyCacheKey(fullID+":"+srv.ObjectID)] = mask[:len(obj.EncryptedShares.BonehFranklinBLS12381.EncryptedShares[i])]
	}
	dec, err := Decrypt(struct {
		EncryptedObject EncryptedObject
		Keys            map[KeyCacheKey][]byte
		CheckLEEncoding bool
	}{EncryptedObject: obj, Keys: keys})
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}
	if !bytes.Equal(dec, plain) {
		t.Fatalf("plaintext mismatch")
	}
}

func TestSessionKeyExportImport(t *testing.T) {
	sk, err := NewSessionKey("0x2", "0xabc", time.Hour)
	if err != nil {
		t.Fatalf("new session key failed: %v", err)
	}
	exported := sk.Export()
	imported := ImportSessionKey(exported)
	if imported.GetPackageID() != "0x2" {
		t.Fatalf("package id mismatch")
	}
}
