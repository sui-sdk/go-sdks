package seal

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"time"
)

type Certificate struct {
	Address   string `json:"address"`
	PackageID string `json:"packageId"`
	ExpiresAt int64  `json:"expiresAt"`
	PubKey    []byte `json:"pubKey"`
	Signature []byte `json:"signature"`
}

type ExportedSessionKey struct {
	SecretKey []byte      `json:"secretKey"`
	Cert      Certificate `json:"certificate"`
}

type SessionKey struct {
	secretKey ed25519.PrivateKey
	publicKey ed25519.PublicKey
	cert      Certificate
}

func NewSessionKey(packageID, address string, ttl time.Duration) (*SessionKey, error) {
	pub, prv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	cert := Certificate{Address: address, PackageID: packageID, ExpiresAt: time.Now().Add(ttl).Unix(), PubKey: pub}
	payload, _ := json.Marshal(map[string]any{"address": address, "packageId": packageID, "expiresAt": cert.ExpiresAt, "pubKey": base64.StdEncoding.EncodeToString(pub)})
	cert.Signature = ed25519.Sign(prv, payload)
	return &SessionKey{secretKey: prv, publicKey: pub, cert: cert}, nil
}

func (s *SessionKey) GetPackageID() string { return s.cert.PackageID }
func (s *SessionKey) GetCertificate() (Certificate, error) { return s.cert, nil }

func (s *SessionKey) CreateRequestParams(txBytes []byte) (map[string]any, error) {
	sig := ed25519.Sign(s.secretKey, txBytes)
	cert, _ := s.GetCertificate()
	return map[string]any{"txBytes": base64.StdEncoding.EncodeToString(txBytes), "signature": base64.StdEncoding.EncodeToString(sig), "certificate": cert}, nil
}

func (s *SessionKey) Export() ExportedSessionKey {
	sk := make([]byte, len(s.secretKey))
	copy(sk, s.secretKey)
	return ExportedSessionKey{SecretKey: sk, Cert: s.cert}
}

func ImportSessionKey(in ExportedSessionKey) *SessionKey {
	prv := ed25519.PrivateKey(in.SecretKey)
	pub := prv.Public().(ed25519.PublicKey)
	return &SessionKey{secretKey: prv, publicKey: pub, cert: in.Cert}
}
