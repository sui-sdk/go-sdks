package transactions

import (
	"math/big"
	"testing"

	"github.com/sui-sdks/go-sdks/bcs"
	stx "github.com/sui-sdks/go-sdks/sui/transactions"
)

func TestPureU64Encoding(t *testing.T) {
	tx := stx.NewTransaction()
	arg := pureU64(tx, 123456789)
	b := pureBytesFromArg(t, tx, arg)
	r := bcs.NewReader(b)
	v, err := r.Read64()
	if err != nil {
		t.Fatalf("read u64 failed: %v", err)
	}
	if v != 123456789 {
		t.Fatalf("u64 mismatch: got %d", v)
	}
}

func TestPureU128StringEncoding(t *testing.T) {
	tx := stx.NewTransaction()
	arg := pureU128String(tx, "340282366920938463463374607431768211455")
	got := readU128StringFromArg(t, tx, arg)
	if got != "340282366920938463463374607431768211455" {
		t.Fatalf("u128 mismatch: got %s", got)
	}
}

func TestPureVecU128Encoding(t *testing.T) {
	tx := stx.NewTransaction()
	arg := pureVecU128(tx, []string{"1", "2", "3"})
	raw := pureBytesFromArg(t, tx, arg)
	r := bcs.NewReader(raw)
	n, err := r.ReadULEB()
	if err != nil {
		t.Fatalf("read vec len failed: %v", err)
	}
	if n != 3 {
		t.Fatalf("expected vec len 3, got %d", n)
	}
	for i := 1; i <= 3; i++ {
		chunk, err := r.ReadBytes(16)
		if err != nil {
			t.Fatalf("read u128[%d] failed: %v", i, err)
		}
		n := new(big.Int).SetBytes(reverse(chunk))
		if n.Cmp(new(big.Int).SetUint64(uint64(i))) != 0 {
			t.Fatalf("unexpected value at %d: got %s", i, n.String())
		}
	}
}

func TestParseU64Fallback(t *testing.T) {
	if got := parseU64("not-number"); got != 0 {
		t.Fatalf("expected fallback 0, got %d", got)
	}
}
