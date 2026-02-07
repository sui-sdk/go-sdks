package verify

import "github.com/sui-sdks/go-sdks/sui/cryptography"

func VerifySignature(pk cryptography.PublicKey, message, signature []byte) bool {
	return pk.Verify(message, signature)
}

func VerifyPersonalMessage(pk cryptography.PublicKey, message, signature []byte) bool {
	return cryptography.VerifyPersonalMessage(pk, message, signature)
}

func VerifyTransaction(pk cryptography.PublicKey, txBytes, signature []byte) bool {
	return cryptography.VerifyTransaction(pk, txBytes, signature)
}
