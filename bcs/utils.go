package bcs

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

type Encoding string

const (
	EncodingBase58 Encoding = "base58"
	EncodingBase64 Encoding = "base64"
	EncodingHex    Encoding = "hex"
)

func EncodeStr(data []byte, encoding Encoding) (string, error) {
	switch encoding {
	case EncodingBase58:
		return ToBase58(data), nil
	case EncodingBase64:
		return base64.StdEncoding.EncodeToString(data), nil
	case EncodingHex:
		return hex.EncodeToString(data), nil
	default:
		return "", fmt.Errorf("unsupported encoding: %s", encoding)
	}
}

func DecodeStr(data string, encoding Encoding) ([]byte, error) {
	switch encoding {
	case EncodingBase58:
		return FromBase58(data)
	case EncodingBase64:
		return base64.StdEncoding.DecodeString(data)
	case EncodingHex:
		data = strings.TrimPrefix(strings.TrimPrefix(data, "0x"), "0X")
		return hex.DecodeString(data)
	default:
		return nil, fmt.Errorf("unsupported encoding: %s", encoding)
	}
}

func SplitGenericParameters(s string, separators [2]rune) []string {
	left, right := separators[0], separators[1]
	if left == 0 || right == 0 {
		left, right = '<', '>'
	}
	tokens := make([]string, 0, 4)
	var b strings.Builder
	nested := 0
	for _, ch := range s {
		if ch == left {
			nested++
		}
		if ch == right {
			nested--
		}
		if nested == 0 && ch == ',' {
			tokens = append(tokens, strings.TrimSpace(b.String()))
			b.Reset()
			continue
		}
		b.WriteRune(ch)
	}
	tokens = append(tokens, strings.TrimSpace(b.String()))
	return tokens
}
