package multisig

import (
	"encoding/base64"
	"fmt"

	"github.com/sui-sdks/go-sdks/sui/cryptography"
)

type PartialSigner interface {
	Sign(data []byte) ([]byte, error)
	GetPublicKey() cryptography.PublicKey
}

type Signer struct {
	PublicKey MultiSigPublicKey
	Signers   []PartialSigner
}

func (s *Signer) Sign(data []byte) ([]byte, error) {
	entries := make([]MultiSigEntry, 0, len(s.Signers))
	bitmap := make([]int, 0, len(s.Signers))
	for i, signer := range s.Signers {
		sig, err := signer.Sign(data)
		if err != nil {
			return nil, err
		}
		entries = append(entries, MultiSigEntry{
			PublicKey: base64.StdEncoding.EncodeToString(signer.GetPublicKey().ToRawBytes()),
			Signature: base64.StdEncoding.EncodeToString(sig),
		})
		bitmap = append(bitmap, i)
	}
	return SerializeMultiSig(MultiSigSerialized{Signatures: entries, Bitmap: bitmap, Threshold: s.PublicKey.Threshold})
}

func (s *Signer) GetKeyScheme() cryptography.SignatureScheme { return cryptography.SchemeMultiSig }
func (s *Signer) GetPublicKey() cryptography.PublicKey       { return s.PublicKey }
func (s *Signer) ToSuiAddress() string                       { return s.PublicKey.ToSuiAddress() }

func (s *Signer) SignWithIntent(bytes []byte, intent cryptography.IntentScope) (cryptography.SignatureWithBytes, error) {
	msg := cryptography.MessageWithIntent(intent, bytes)
	sig, err := s.Sign(msg)
	if err != nil {
		return cryptography.SignatureWithBytes{}, err
	}
	serialized, err := cryptography.ToSerializedSignature(cryptography.SerializeSignatureInput{
		SignatureScheme: cryptography.SchemeMultiSig,
		Signature:       sig,
		PublicKey:       s.PublicKey,
	})
	if err != nil {
		return cryptography.SignatureWithBytes{}, err
	}
	return cryptography.SignatureWithBytes{Bytes: base64.StdEncoding.EncodeToString(bytes), Signature: serialized}, nil
}

func (s *Signer) SignTransaction(bytes []byte) (cryptography.SignatureWithBytes, error) {
	return s.SignWithIntent(bytes, cryptography.IntentTransactionData)
}

func (s *Signer) SignPersonalMessage(bytes []byte) (cryptography.SignatureWithBytes, error) {
	return s.SignWithIntent(bytes, cryptography.IntentPersonalMessage)
}

func (s *Signer) GetSecretKey() string {
	return fmt.Sprintf("multisig:%d", len(s.Signers))
}
