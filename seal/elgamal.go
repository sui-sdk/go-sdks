package seal

import (
	"crypto/rand"
	"crypto/sha256"
)

func ElgamalDecrypt(sk []byte, c [2][]byte) []byte {
	h := sha256.Sum256(append(sk, c[0]...))
	out := make([]byte, len(c[1]))
	for i := range c[1] {
		out[i] = c[1][i] ^ h[i%len(h)]
	}
	return out
}

func GenerateSecretKey() ([]byte, error) {
	sk := make([]byte, 32)
	_, err := rand.Read(sk)
	return sk, err
}

func ToPublicKey(sk []byte) []byte {
	h := sha256.Sum256(append([]byte("pk"), sk...))
	return h[:]
}

func ToVerificationKey(sk []byte) []byte {
	h := sha256.Sum256(append([]byte("vk"), sk...))
	return h[:]
}
