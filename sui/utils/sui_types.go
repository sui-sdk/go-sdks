package utils

import (
	"encoding/hex"
	"fmt"
	"strings"
)

const SuiAddressLength = 32

func NormalizeSuiAddress(addr string) string {
	a := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(addr)), "0x")
	if len(a) > SuiAddressLength*2 {
		a = a[len(a)-SuiAddressLength*2:]
	}
	if len(a) < SuiAddressLength*2 {
		a = strings.Repeat("0", SuiAddressLength*2-len(a)) + a
	}
	return "0x" + a
}

func NormalizeSuiObjectID(id string) string {
	return NormalizeSuiAddress(id)
}

func IsValidSuiAddress(addr string) bool {
	a := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(addr)), "0x")
	if len(a) == 0 || len(a) > SuiAddressLength*2 {
		return false
	}
	_, err := hex.DecodeString(a)
	return err == nil
}

func IsValidSuiObjectID(id string) bool {
	return IsValidSuiAddress(id)
}

func MustNormalizeSuiAddress(addr string) string {
	if !IsValidSuiAddress(addr) {
		panic(fmt.Sprintf("invalid Sui address: %s", addr))
	}
	return NormalizeSuiAddress(addr)
}
