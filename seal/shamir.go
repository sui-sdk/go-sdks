package seal

import (
	"crypto/rand"
	"fmt"
)

type Share struct {
	Index int    `json:"index"`
	Share []byte `json:"share"`
}

func gfMul(a, b byte) byte {
	var p byte
	for b > 0 {
		if b&1 == 1 {
			p ^= a
		}
		hi := a & 0x80
		a <<= 1
		if hi > 0 {
			a ^= 0x1b
		}
		b >>= 1
	}
	return p
}

func gfPow(a byte, n int) byte {
	if n == 0 {
		return 1
	}
	out := byte(1)
	for i := 0; i < n; i++ {
		out = gfMul(out, a)
	}
	return out
}

func gfInv(a byte) byte {
	if a == 0 {
		return 0
	}
	// a^254 in GF(2^8)
	return gfPow(a, 254)
}

func evalPoly(coeffs []byte, x byte) byte {
	out := coeffs[len(coeffs)-1]
	for i := len(coeffs) - 2; i >= 0; i-- {
		out = gfMul(out, x) ^ coeffs[i]
	}
	return out
}

func Split(secret []byte, threshold, total int) ([]Share, error) {
	if threshold < 1 || total < threshold || total >= MaxU8 {
		return nil, fmt.Errorf("invalid threshold/total")
	}
	shares := make([]Share, total)
	for i := 0; i < total; i++ {
		shares[i] = Share{Index: i + 1, Share: make([]byte, len(secret))}
	}
	for pos := range secret {
		coeffs := make([]byte, threshold)
		coeffs[0] = secret[pos]
		if _, err := rand.Read(coeffs[1:]); err != nil {
			return nil, err
		}
		for i := 0; i < total; i++ {
			x := byte(i + 1)
			shares[i].Share[pos] = evalPoly(coeffs, x)
		}
	}
	return shares, nil
}

func Combine(shares []Share) ([]byte, error) {
	if len(shares) == 0 {
		return nil, fmt.Errorf("no shares")
	}
	ln := len(shares[0].Share)
	for _, s := range shares {
		if len(s.Share) != ln {
			return nil, fmt.Errorf("inconsistent share lengths")
		}
	}
	secret := make([]byte, ln)
	for pos := 0; pos < ln; pos++ {
		var result byte
		for i, si := range shares {
			xi := byte(si.Index)
			num := byte(1)
			den := byte(1)
			for j, sj := range shares {
				if i == j {
					continue
				}
				xj := byte(sj.Index)
				num = gfMul(num, xj)
				den = gfMul(den, xi^xj)
			}
			li := gfMul(num, gfInv(den))
			result ^= gfMul(si.Share[pos], li)
		}
		secret[pos] = result
	}
	return secret, nil
}
