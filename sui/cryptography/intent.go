package cryptography

import (
	"bytes"
	"encoding/binary"
)

type IntentScope string

const (
	IntentTransactionData IntentScope = "TransactionData"
	IntentPersonalMessage IntentScope = "PersonalMessage"
)

func MessageWithIntent(scope IntentScope, message []byte) []byte {
	// A compact Go-friendly intent prefix, preserving domain separation semantics.
	var scopeFlag byte
	switch scope {
	case IntentTransactionData:
		scopeFlag = 0
	case IntentPersonalMessage:
		scopeFlag = 3
	default:
		scopeFlag = 255
	}
	buf := bytes.NewBuffer(make([]byte, 0, len(message)+8))
	buf.WriteByte(scopeFlag)
	buf.WriteByte(0) // version V0
	buf.WriteByte(0) // app Sui
	var ln [4]byte
	binary.LittleEndian.PutUint32(ln[:], uint32(len(message)))
	buf.Write(ln[:])
	buf.Write(message)
	return buf.Bytes()
}
