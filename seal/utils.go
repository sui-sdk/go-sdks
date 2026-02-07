package seal

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

const (
	MaxU8                = 255
	SuiAddressLength     = 32
	EncryptedShareLength = 32
	KeyLength            = 32
)

func Xor(a, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("length mismatch")
	}
	return xorUnchecked(a, b), nil
}

func xorUnchecked(a, b []byte) []byte {
	out := make([]byte, len(a))
	for i := range a {
		out[i] = a[i] ^ b[i]
	}
	return out
}

func CreateFullID(packageID, innerID string) string {
	return strings.TrimPrefix(packageID, "0x") + ":" + innerID
}

func Flatten(arrays [][]byte) []byte {
	total := 0
	for _, a := range arrays {
		total += len(a)
	}
	out := make([]byte, 0, total)
	for _, a := range arrays {
		out = append(out, a...)
	}
	return out
}

func Count[T comparable](array []T, value T) int {
	count := 0
	for _, v := range array {
		if v == value {
			count++
		}
	}
	return count
}

func HasDuplicates(array []int) bool {
	seen := map[int]struct{}{}
	for _, v := range array {
		if _, ok := seen[v]; ok {
			return true
		}
		seen[v] = struct{}{}
	}
	return false
}

func AllEqual(array []int) bool {
	if len(array) < 2 {
		return true
	}
	first := array[0]
	for _, v := range array[1:] {
		if v != first {
			return false
		}
	}
	return true
}

func Equals(a, b []byte) bool {
	return bytes.Equal(a, b)
}

type Version struct{ raw string }

func NewVersion(v string) Version { return Version{raw: v} }
func (v Version) String() string { return v.raw }

func (v Version) Compare(other Version) int {
	pa := parseVersion(v.raw)
	pb := parseVersion(other.raw)
	for i := 0; i < 3; i++ {
		if pa[i] < pb[i] {
			return -1
		}
		if pa[i] > pb[i] {
			return 1
		}
	}
	return 0
}

func parseVersion(v string) [3]int {
	var out [3]int
	parts := strings.Split(v, ".")
	for i := range out {
		if i < len(parts) {
			fmt.Sscanf(parts[i], "%d", &out[i])
		}
	}
	return out
}

func deriveDigest(parts ...[]byte) []byte {
	h := sha256.New()
	for _, p := range parts {
		h.Write(p)
	}
	return h.Sum(nil)
}

func NormalizeHex(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	s = strings.TrimPrefix(s, "0x")
	return s
}

func HexToBytes(s string) ([]byte, error) {
	return hex.DecodeString(NormalizeHex(s))
}

func SortStrings(ss []string) {
	sort.Strings(ss)
}
