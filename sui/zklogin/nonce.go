package zklogin

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const NonceLength = 27

func GenerateRandomness() ([]byte, error) {
	out := make([]byte, 16)
	_, err := rand.Read(out)
	return out, err
}

func GenerateNonce(publicKey []byte, maxEpoch uint64, randomness []byte) (string, error) {
	if len(publicKey) == 0 || len(randomness) == 0 {
		return "", fmt.Errorf("invalid inputs")
	}
	buf := append([]byte{}, publicKey...)
	buf = append(buf, byte(maxEpoch), byte(maxEpoch>>8), byte(maxEpoch>>16), byte(maxEpoch>>24))
	buf = append(buf, randomness...)
	out := base64.RawURLEncoding.EncodeToString(buf)
	if len(out) > NonceLength {
		out = out[:NonceLength]
	}
	return out, nil
}
