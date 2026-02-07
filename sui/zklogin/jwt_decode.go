package zklogin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

type JWTHeader map[string]any
type JWTPayload map[string]any

type DecodeOptions struct {
	Header bool
}

func JWTDecode(token string, options DecodeOptions) (any, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid token")
	}
	idx := 1
	if options.Header {
		idx = 0
	}
	payloadPart := parts[idx]
	payloadPart += strings.Repeat("=", (4-len(payloadPart)%4)%4)
	b, err := base64.URLEncoding.DecodeString(payloadPart)
	if err != nil {
		return nil, err
	}
	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}
