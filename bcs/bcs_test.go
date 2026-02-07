package bcs

import (
	"fmt"
	"testing"
)

func TestULEBRoundTrip(t *testing.T) {
	values := []uint64{0, 1, 127, 128, 255, 300, 16384, 1<<32 - 1}
	for _, v := range values {
		enc := ULEBEncode(v)
		dec, n, err := ULEBDecode(enc)
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		if dec != v {
			t.Fatalf("expected %d got %d", v, dec)
		}
		if n != len(enc) {
			t.Fatalf("expected consumed %d got %d", len(enc), n)
		}
	}
}

func TestULEBRejectsNonCanonicalAndOverflow(t *testing.T) {
	_, _, err := ULEBDecode([]byte{0x80, 0x00})
	if err == nil {
		t.Fatalf("expected non-canonical uleb decode error")
	}

	// 2^32 encoded in ULEB128 exceeds BCS u32 limit.
	tooLarge := ULEBEncode(uint64(1) << 32)
	_, _, err = ULEBDecode(tooLarge)
	if err == nil {
		t.Fatalf("expected u32 overflow decode error")
	}
}

func TestBase58RoundTrip(t *testing.T) {
	in := []byte{0, 1, 2, 3, 4, 5, 255}
	s := ToBase58(in)
	out, err := FromBase58(s)
	if err != nil {
		t.Fatalf("from base58 failed: %v", err)
	}
	if string(out) != string(in) {
		t.Fatalf("roundtrip mismatch")
	}
}

func TestBCSStructEnumVectorMap(t *testing.T) {
	myEnum := BCS.Enum("MyEnum", []Field{{Name: "None", Type: nil}, {Name: "Num", Type: BCS.U8()}})
	myStruct := BCS.Struct("MyStruct", []Field{
		{Name: "name", Type: BCS.String()},
		{Name: "values", Type: BCS.Vector(BCS.U16())},
		{Name: "kind", Type: myEnum},
	})

	input := map[string]any{
		"name":   "alice",
		"values": []any{uint16(1), uint16(2), uint16(3)},
		"kind":   map[string]any{"Num": uint8(7)},
	}

	ser, err := myStruct.Serialize(input, nil)
	if err != nil {
		t.Fatalf("serialize failed: %v", err)
	}
	parsedAny, err := myStruct.Parse(ser.ToBytes())
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	parsed := parsedAny.(map[string]any)
	if parsed["name"].(string) != "alice" {
		t.Fatalf("expected alice")
	}

	mapType := BCS.Map(BCS.String(), BCS.U8())
	m := map[string]any{"b": uint8(2), "a": uint8(1)}
	mapSer, err := mapType.Serialize(m, nil)
	if err != nil {
		t.Fatalf("map serialize failed: %v", err)
	}
	if len(mapSer.ToBytes()) == 0 {
		t.Fatalf("expected bytes")
	}
}

func TestBoolParseRejectsNonCanonicalByte(t *testing.T) {
	_, err := BCS.Bool().Parse([]byte{2})
	if err == nil {
		t.Fatalf("expected bool parse error for byte 2")
	}
}

func TestParseRejectsTrailingBytes(t *testing.T) {
	_, err := BCS.U8().Parse([]byte{1, 2})
	if err == nil {
		t.Fatalf("expected trailing bytes parse error")
	}
}

func TestMapParseRejectsUnsortedKeys(t *testing.T) {
	mapType := BCS.Map(BCS.String(), BCS.U8())
	// len=2, first key "b", second key "a" (unsorted).
	invalid := []byte{
		0x02,
		0x01, 'b', 0x01,
		0x01, 'a', 0x02,
	}
	_, err := mapType.Parse(invalid)
	if err == nil {
		t.Fatalf("expected map parse error for unsorted keys")
	}
}

func TestSequenceLengthLimitEnforced(t *testing.T) {
	byteVector := BCS.ByteVector()
	tooLong := uint64(MaxSequenceLength) + 1
	input := ULEBEncode(tooLong)
	_, err := byteVector.Parse(input)
	if err == nil {
		t.Fatalf("expected sequence length limit error")
	}
}

func TestContainerDepthLimitEnforced(t *testing.T) {
	var node *Type
	node = BCS.Enum("Node", []Field{
		{Name: "Leaf", Type: nil},
		{Name: "Child", Type: BCS.Lazy(func() *Type { return node })},
	})

	v := map[string]any{"$kind": "Leaf", "Leaf": true}
	for i := 0; i < MaxContainerDepth+1; i++ {
		v = map[string]any{"$kind": "Child", "Child": v}
	}

	_, err := node.Serialize(v, nil)
	if err == nil {
		t.Fatalf("expected container depth error")
	}
	if got := err.Error(); got == "" {
		t.Fatalf("expected non-empty error")
	}
}

func TestMapParseRejectsDuplicateKeys(t *testing.T) {
	mapType := BCS.Map(BCS.String(), BCS.U8())
	// len=2, key "a" twice (not strictly increasing).
	invalid := []byte{
		0x02,
		0x01, 'a', 0x01,
		0x01, 'a', 0x02,
	}
	_, err := mapType.Parse(invalid)
	if err == nil {
		t.Fatalf("expected map parse error for duplicate keys")
	}
}

func TestULEBDecodeErrorMessages(t *testing.T) {
	_, _, err := ULEBDecode([]byte{})
	if err == nil {
		t.Fatalf("expected empty input error")
	}
	if err != nil && err.Error() == fmt.Sprintf("%v", nil) {
		t.Fatalf("unexpected nil-like error")
	}
}
