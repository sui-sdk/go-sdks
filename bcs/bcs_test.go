package bcs

import (
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
