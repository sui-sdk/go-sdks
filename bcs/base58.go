package bcs

import (
	"fmt"
	"strings"
)

const b58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

var b58Indexes = func() map[rune]int {
	m := make(map[rune]int, len(b58Alphabet))
	for i, r := range b58Alphabet {
		m[r] = i
	}
	return m
}()

func ToBase58(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	zeros := 0
	for zeros < len(data) && data[zeros] == 0 {
		zeros++
	}
	buf := append([]byte(nil), data...)
	encoded := make([]byte, 0, len(data)*138/100+1)
	for start := zeros; start < len(buf); {
		carry := 0
		for i := start; i < len(buf); i++ {
			v := int(buf[i]) + carry*256
			buf[i] = byte(v / 58)
			carry = v % 58
		}
		encoded = append(encoded, b58Alphabet[carry])
		for start < len(buf) && buf[start] == 0 {
			start++
		}
	}
	for i := 0; i < zeros; i++ {
		encoded = append(encoded, '1')
	}
	for i, j := 0, len(encoded)-1; i < j; i, j = i+1, j-1 {
		encoded[i], encoded[j] = encoded[j], encoded[i]
	}
	return string(encoded)
}

func FromBase58(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return []byte{}, nil
	}
	zeros := 0
	for zeros < len(s) && s[zeros] == '1' {
		zeros++
	}
	decoded := make([]byte, 0, len(s))
	for _, r := range s {
		idx, ok := b58Indexes[r]
		if !ok {
			return nil, fmt.Errorf("invalid base58 char: %q", r)
		}
		carry := idx
		for i := 0; i < len(decoded); i++ {
			v := int(decoded[i])*58 + carry
			decoded[i] = byte(v & 0xff)
			carry = v >> 8
		}
		for carry > 0 {
			decoded = append(decoded, byte(carry&0xff))
			carry >>= 8
		}
	}
	for i := 0; i < zeros; i++ {
		decoded = append(decoded, 0)
	}
	for i, j := 0, len(decoded)-1; i < j; i, j = i+1, j-1 {
		decoded[i], decoded[j] = decoded[j], decoded[i]
	}
	return decoded, nil
}
