package bcs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type rustVector struct {
	ID    string `json:"id"`
	Kind  string `json:"kind"`
	Type  string `json:"type"`
	Value any    `json:"value"`
	Hex   string `json:"hex"`
}

func loadRustVectors(t *testing.T) []rustVector {
	t.Helper()
	path := filepath.Join("testdata", "rust_official_vectors.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read vectors failed: %v", err)
	}
	var out []rustVector
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal vectors failed: %v", err)
	}
	if len(out) == 0 {
		t.Fatalf("no vectors loaded")
	}
	return out
}

func resolveType(t *testing.T, name string) *Type {
	t.Helper()
	switch name {
	case "bool":
		return BCS.Bool()
	case "u16":
		return BCS.U16()
	case "u32":
		return BCS.U32()
	case "u64":
		return BCS.U64()
	case "string":
		return BCS.String()
	case "vector<u16>":
		return BCS.Vector(BCS.U16())
	case "fixed_array_3_u16":
		return BCS.FixedArray(3, BCS.U16())
	default:
		t.Fatalf("unsupported vector type %q", name)
		return nil
	}
}

func TestRustOfficialVectors(t *testing.T) {
	vectors := loadRustVectors(t)
	for _, v := range vectors {
		v := v
		t.Run(v.ID, func(t *testing.T) {
			switch v.Kind {
			case "uleb":
				runULEBVector(t, v)
			case "type_serialize":
				runTypeSerializeVector(t, v)
			default:
				t.Fatalf("unsupported vector kind %q", v.Kind)
			}
		})
	}
}

func runULEBVector(t *testing.T, v rustVector) {
	t.Helper()
	n, ok := toUint64(v.Value)
	if !ok {
		t.Fatalf("invalid uleb value: %#v", v.Value)
	}
	gotHex, err := EncodeStr(ULEBEncode(n), EncodingHex)
	if err != nil {
		t.Fatalf("encode uleb hex failed: %v", err)
	}
	if gotHex != v.Hex {
		t.Fatalf("uleb serialize mismatch: got %s want %s", gotHex, v.Hex)
	}
	raw, err := DecodeStr(v.Hex, EncodingHex)
	if err != nil {
		t.Fatalf("decode hex failed: %v", err)
	}
	decoded, consumed, err := ULEBDecode(raw)
	if err != nil {
		t.Fatalf("uleb decode failed: %v", err)
	}
	if consumed != len(raw) {
		t.Fatalf("uleb consumed mismatch: got %d want %d", consumed, len(raw))
	}
	if decoded != n {
		t.Fatalf("uleb decode mismatch: got %d want %d", decoded, n)
	}
}

func runTypeSerializeVector(t *testing.T, v rustVector) {
	t.Helper()
	typ := resolveType(t, v.Type)
	serialized, err := typ.Serialize(v.Value, nil)
	if err != nil {
		t.Fatalf("serialize failed for type %s: %v", v.Type, err)
	}
	gotHex, err := serialized.ToHex()
	if err != nil {
		t.Fatalf("to hex failed: %v", err)
	}
	if gotHex != v.Hex {
		t.Fatalf("serialize mismatch for %s: got %s want %s", v.Type, gotHex, v.Hex)
	}

	decodedAny, err := typ.Parse(serialized.ToBytes())
	if err != nil {
		t.Fatalf("parse failed for %s: %v", v.Type, err)
	}
	reSerialized, err := typ.Serialize(decodedAny, nil)
	if err != nil {
		t.Fatalf("re-serialize failed for %s: %v", v.Type, err)
	}
	reHex, err := reSerialized.ToHex()
	if err != nil {
		t.Fatalf("re-hex failed for %s: %v", v.Type, err)
	}
	if reHex != v.Hex {
		t.Fatalf("roundtrip mismatch for %s: got %s want %s", v.Type, reHex, v.Hex)
	}
}

func TestRustVectorFileHasUniqueIDs(t *testing.T) {
	vectors := loadRustVectors(t)
	seen := map[string]bool{}
	for _, v := range vectors {
		if v.ID == "" {
			t.Fatalf("vector id cannot be empty")
		}
		if seen[v.ID] {
			t.Fatalf("duplicate vector id %q", v.ID)
		}
		seen[v.ID] = true
		if v.Hex == "" {
			t.Fatalf("vector %s missing hex", v.ID)
		}
		if _, err := DecodeStr(v.Hex, EncodingHex); err != nil {
			t.Fatalf("vector %s has invalid hex: %v", v.ID, err)
		}
	}
}

func TestRustVectorTypeResolverCoverage(t *testing.T) {
	vectors := loadRustVectors(t)
	for _, v := range vectors {
		if v.Kind != "type_serialize" {
			continue
		}
		name := v.Type
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("unsupported type in vector file: %s (%v)", name, r)
				}
			}()
			_ = resolveType(t, name)
		}()
	}
}

func Example_rustVectorFormat() {
	fmt.Println(`{"id":"u16_4660","kind":"type_serialize","type":"u16","value":4660,"hex":"3412"}`)
	// Output: {"id":"u16_4660","kind":"type_serialize","type":"u16","value":4660,"hex":"3412"}
}
