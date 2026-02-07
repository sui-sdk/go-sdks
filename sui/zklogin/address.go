package zklogin

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/sui-sdks/go-sdks/sui/utils"
)

const (
	MaxHeaderLenB64         = 248
	MaxPaddedUnsignedJWTLen = 64 * 25
)

func ComputeZkLoginAddressFromSeed(seed []byte, iss, aud string) string {
	h := sha256.Sum256(append(append(seed, []byte(iss)...), []byte(aud)...))
	return utils.NormalizeSuiAddress(hex.EncodeToString(h[:]))
}

func LengthChecks(jwt string) error {
	parts := strings.Split(jwt, ".")
	if len(parts) < 2 {
		return fmt.Errorf("invalid jwt")
	}
	if len(parts[0]) > MaxHeaderLenB64 {
		return fmt.Errorf("header too long")
	}
	if len(parts[0])+len(parts[1]) > MaxPaddedUnsignedJWTLen {
		return fmt.Errorf("jwt too long")
	}
	return nil
}

func JWTToAddress(jwt string, userSalt string, legacyAddress bool) (string, error) {
	if err := LengthChecks(jwt); err != nil {
		return "", err
	}
	h := sha256.Sum256([]byte(jwt + ":" + userSalt))
	return utils.NormalizeSuiAddress(hex.EncodeToString(h[:])), nil
}

type ComputeZkLoginAddressOptions struct {
	Iss             string
	Aud             string
	UserSalt        string
	Jwt             string
	LegacyAddress   bool
}

func ComputeZkLoginAddress(opts ComputeZkLoginAddressOptions) (string, error) {
	if opts.Jwt != "" {
		return JWTToAddress(opts.Jwt, opts.UserSalt, opts.LegacyAddress)
	}
	return ComputeZkLoginAddressFromSeed([]byte(opts.UserSalt), opts.Iss, opts.Aud), nil
}
