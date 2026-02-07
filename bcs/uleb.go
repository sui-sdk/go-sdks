package bcs

import "fmt"

func ULEBEncode(num uint64) []byte {
	if num == 0 {
		return []byte{0}
	}
	out := make([]byte, 0, 10)
	for num > 0 {
		b := byte(num & 0x7f)
		num >>= 7
		if num > 0 {
			b |= 0x80
		}
		out = append(out, b)
	}
	return out
}

func ULEBDecode(arr []byte) (value uint64, length int, err error) {
	var total uint64
	var shift uint
	for i, b := range arr {
		total |= uint64(b&0x7f) << shift
		if b&0x80 == 0 {
			return total, i + 1, nil
		}
		shift += 7
		if shift >= 64 {
			return 0, 0, fmt.Errorf("uleb decode overflow")
		}
	}
	return 0, 0, fmt.Errorf("uleb decode error: buffer overflow")
}
