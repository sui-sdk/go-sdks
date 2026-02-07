package transactions

import (
	"math"
	"math/big"
	"strconv"

	"github.com/sui-sdks/go-sdks/bcs"
	stx "github.com/sui-sdks/go-sdks/sui/transactions"
)

func pureBool(tx *stx.Transaction, v bool) stx.Argument {
	b := byte(0)
	if v {
		b = 1
	}
	return tx.PureBytes([]byte{b})
}

func pureU8(tx *stx.Transaction, v uint8) stx.Argument {
	return tx.PureBytes([]byte{v})
}

func pureU64(tx *stx.Transaction, v uint64) stx.Argument {
	w := bcs.NewWriter(nil)
	_ = w.Write64(v)
	return tx.PureBytes(w.ToBytes())
}

func pureU128String(tx *stx.Transaction, s string) stx.Argument {
	n := new(big.Int)
	if _, ok := n.SetString(s, 10); !ok {
		n.SetUint64(0)
	}
	out := make([]byte, 16)
	b := n.Bytes()
	for i := 0; i < len(b) && i < 16; i++ {
		out[i] = b[len(b)-1-i]
	}
	return tx.PureBytes(out)
}

func pureVecU128(tx *stx.Transaction, values []string) stx.Argument {
	w := bcs.NewWriter(nil)
	_ = w.WriteULEB(uint64(len(values)))
	for _, s := range values {
		n := new(big.Int)
		if _, ok := n.SetString(s, 10); !ok {
			n.SetUint64(0)
		}
		out := make([]byte, 16)
		b := n.Bytes()
		for i := 0; i < len(b) && i < 16; i++ {
			out[i] = b[len(b)-1-i]
		}
		_ = w.WriteBytes(out)
	}
	return tx.PureBytes(w.ToBytes())
}

func boolToByte(v bool) byte {
	if v {
		return 1
	}
	return 0
}

func parseU64(s string) uint64 {
	return uint64(math.Round(parseFloat(s)))
}

func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
