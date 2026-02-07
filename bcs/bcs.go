package bcs

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"sort"
	"unicode/utf8"
)

type Option struct {
	Name     string
	Validate func(any) error
}

const (
	// Matches Rust bcs crate limits.
	MaxSequenceLength = 1<<31 - 1
	MaxContainerDepth = 500
)

type Field struct {
	Name string
	Type *Type
}

func CompareBCSBytes(a, b []byte) int {
	return bytes.Compare(a, b)
}

func fixedSizeType(name string, size int, read func(*Reader) (any, error), write func(any, *Writer) error, validate func(any) error) *Type {
	return &Type{
		Name:           name,
		Read:           read,
		Write:          write,
		SerializedSize: func(any) (int, bool) { return size, true },
		Validate:       validate,
	}
}

func dynamicType(name string, read func(*Reader) (any, error), write func(any, *Writer) error, size func(any) (int, bool), validate func(any) error) *Type {
	return &Type{Name: name, Read: read, Write: write, SerializedSize: size, Validate: validate}
}

type bcsBuilder struct{}

var BCS = bcsBuilder{}

func (bcsBuilder) U8(opts ...Option) *Type {
	op := pickOption(opts)
	name := coalesce(op.Name, "u8")
	return fixedSizeType(name, 1,
		func(r *Reader) (any, error) { v, e := r.Read8(); return v, e },
		func(v any, w *Writer) error {
			n, ok := toUint64(v)
			if !ok || n > math.MaxUint8 {
				return fmt.Errorf("invalid %s value: %v", name, v)
			}
			return w.Write8(uint8(n))
		},
		op.Validate,
	)
}

func (bcsBuilder) U16(opts ...Option) *Type {
	op := pickOption(opts)
	name := coalesce(op.Name, "u16")
	return fixedSizeType(name, 2,
		func(r *Reader) (any, error) { v, e := r.Read16(); return v, e },
		func(v any, w *Writer) error {
			n, ok := toUint64(v)
			if !ok || n > math.MaxUint16 {
				return fmt.Errorf("invalid %s value: %v", name, v)
			}
			return w.Write16(uint16(n))
		},
		op.Validate,
	)
}

func (bcsBuilder) U32(opts ...Option) *Type {
	op := pickOption(opts)
	name := coalesce(op.Name, "u32")
	return fixedSizeType(name, 4,
		func(r *Reader) (any, error) { v, e := r.Read32(); return v, e },
		func(v any, w *Writer) error {
			n, ok := toUint64(v)
			if !ok || n > math.MaxUint32 {
				return fmt.Errorf("invalid %s value: %v", name, v)
			}
			return w.Write32(uint32(n))
		},
		op.Validate,
	)
}

func (bcsBuilder) U64(opts ...Option) *Type {
	op := pickOption(opts)
	name := coalesce(op.Name, "u64")
	return fixedSizeType(name, 8,
		func(r *Reader) (any, error) {
			v, e := r.Read64()
			if e != nil {
				return nil, e
			}
			return fmt.Sprintf("%d", v), nil
		},
		func(v any, w *Writer) error {
			n, err := toBigInt(v)
			if err != nil || n.Sign() < 0 || n.BitLen() > 64 {
				return fmt.Errorf("invalid %s value: %v", name, v)
			}
			return w.Write64(n.Uint64())
		},
		op.Validate,
	)
}

func (bcsBuilder) U128(opts ...Option) *Type {
	op := pickOption(opts)
	name := coalesce(op.Name, "u128")
	return fixedSizeType(name, 16,
		func(r *Reader) (any, error) {
			v, err := r.Read128()
			if err != nil {
				return nil, err
			}
			bi := littleToBigInt(v[:])
			return bi.String(), nil
		},
		func(v any, w *Writer) error {
			bi, err := toBigInt(v)
			if err != nil || bi.Sign() < 0 || bi.BitLen() > 128 {
				return fmt.Errorf("invalid %s value: %v", name, v)
			}
			return w.WriteBytes(bigIntToLittleFixed(bi, 16))
		},
		op.Validate,
	)
}

func (bcsBuilder) U256(opts ...Option) *Type {
	op := pickOption(opts)
	name := coalesce(op.Name, "u256")
	return fixedSizeType(name, 32,
		func(r *Reader) (any, error) {
			v, err := r.Read256()
			if err != nil {
				return nil, err
			}
			bi := littleToBigInt(v[:])
			return bi.String(), nil
		},
		func(v any, w *Writer) error {
			bi, err := toBigInt(v)
			if err != nil || bi.Sign() < 0 || bi.BitLen() > 256 {
				return fmt.Errorf("invalid %s value: %v", name, v)
			}
			return w.WriteBytes(bigIntToLittleFixed(bi, 32))
		},
		op.Validate,
	)
}

func (bcsBuilder) Bool(opts ...Option) *Type {
	op := pickOption(opts)
	name := coalesce(op.Name, "bool")
	return fixedSizeType(name, 1,
		func(r *Reader) (any, error) {
			v, e := r.Read8()
			if e != nil {
				return nil, e
			}
			if v == 0 {
				return false, nil
			}
			if v == 1 {
				return true, nil
			}
			return nil, fmt.Errorf("invalid bool byte: %d", v)
		},
		func(v any, w *Writer) error {
			b, ok := v.(bool)
			if !ok {
				return fmt.Errorf("invalid %s value: %v", name, v)
			}
			if b {
				return w.Write8(1)
			}
			return w.Write8(0)
		},
		op.Validate,
	)
}

func (bcsBuilder) ULEB128(opts ...Option) *Type {
	op := pickOption(opts)
	name := coalesce(op.Name, "uleb128")
	return dynamicType(name,
		func(r *Reader) (any, error) { v, e := r.ReadULEB(); return v, e },
		func(v any, w *Writer) error {
			n, ok := toUint64(v)
			if !ok {
				return fmt.Errorf("invalid %s value: %v", name, v)
			}
			return w.WriteULEB(n)
		},
		func(v any) (int, bool) {
			n, ok := toUint64(v)
			if !ok {
				return 0, false
			}
			return len(ULEBEncode(n)), true
		},
		op.Validate,
	)
}

func (bcsBuilder) Bytes(size int, opts ...Option) *Type {
	op := pickOption(opts)
	name := coalesce(op.Name, fmt.Sprintf("bytes[%d]", size))
	return fixedSizeType(name, size,
		func(r *Reader) (any, error) { return r.ReadBytes(size) },
		func(v any, w *Writer) error {
			bs, err := asBytes(v)
			if err != nil || len(bs) != size {
				return fmt.Errorf("invalid %s size: %v", name, v)
			}
			return w.WriteBytes(bs)
		},
		op.Validate,
	)
}

func (bcsBuilder) ByteVector(opts ...Option) *Type {
	op := pickOption(opts)
	name := coalesce(op.Name, "vector<u8>")
	return dynamicType(name,
		func(r *Reader) (any, error) {
			ln, err := r.ReadULEB()
			if err != nil {
				return nil, err
			}
			if ln > MaxSequenceLength {
				return nil, fmt.Errorf("sequence length %d exceeds limit %d", ln, MaxSequenceLength)
			}
			return r.ReadBytes(int(ln))
		},
		func(v any, w *Writer) error {
			bs, err := asBytes(v)
			if err != nil {
				return err
			}
			if len(bs) > MaxSequenceLength {
				return fmt.Errorf("sequence length %d exceeds limit %d", len(bs), MaxSequenceLength)
			}
			if err := w.WriteULEB(uint64(len(bs))); err != nil {
				return err
			}
			return w.WriteBytes(bs)
		},
		func(v any) (int, bool) {
			bs, err := asBytes(v)
			if err != nil {
				return 0, false
			}
			return len(ULEBEncode(uint64(len(bs)))) + len(bs), true
		},
		op.Validate,
	)
}

func (bcsBuilder) String(opts ...Option) *Type {
	op := pickOption(opts)
	name := coalesce(op.Name, "string")
	return dynamicType(name,
		func(r *Reader) (any, error) {
			ln, err := r.ReadULEB()
			if err != nil {
				return nil, err
			}
			if ln > MaxSequenceLength {
				return nil, fmt.Errorf("sequence length %d exceeds limit %d", ln, MaxSequenceLength)
			}
			b, err := r.ReadBytes(int(ln))
			if err != nil {
				return nil, err
			}
			if !utf8.Valid(b) {
				return nil, fmt.Errorf("invalid utf-8 string")
			}
			return string(b), nil
		},
		func(v any, w *Writer) error {
			s, ok := v.(string)
			if !ok {
				return fmt.Errorf("invalid %s value: %v", name, v)
			}
			bs := []byte(s)
			if len(bs) > MaxSequenceLength {
				return fmt.Errorf("sequence length %d exceeds limit %d", len(bs), MaxSequenceLength)
			}
			if err := w.WriteULEB(uint64(len(bs))); err != nil {
				return err
			}
			return w.WriteBytes(bs)
		},
		func(v any) (int, bool) {
			s, ok := v.(string)
			if !ok {
				return 0, false
			}
			return len(ULEBEncode(uint64(len(s)))) + len(s), true
		},
		op.Validate,
	)
}

func (bcsBuilder) FixedArray(size int, t *Type, opts ...Option) *Type {
	op := pickOption(opts)
	name := coalesce(op.Name, fmt.Sprintf("%s[%d]", t.Name, size))
	return dynamicType(name,
		func(r *Reader) (any, error) {
			res := make([]any, size)
			for i := 0; i < size; i++ {
				v, err := t.Read(r)
				if err != nil {
					return nil, err
				}
				res[i] = v
			}
			return res, nil
		},
		func(v any, w *Writer) error {
			arr, err := asAnySlice(v)
			if err != nil || len(arr) != size {
				return fmt.Errorf("expected fixed array length %d", size)
			}
			for _, el := range arr {
				if err := t.Write(el, w); err != nil {
					return err
				}
			}
			return nil
		},
		func(v any) (int, bool) {
			arr, err := asAnySlice(v)
			if err != nil || len(arr) != size {
				return 0, false
			}
			total := 0
			for _, el := range arr {
				s, ok := t.SerializedSize(el)
				if !ok {
					return 0, false
				}
				total += s
			}
			return total, true
		},
		op.Validate,
	)
}

func (bcsBuilder) Option(t *Type) *Type {
	enum := BCS.Enum(fmt.Sprintf("Option<%s>", t.Name), []Field{{Name: "None", Type: nil}, {Name: "Some", Type: t}})
	return enum.Transform(TransformOptions{
		Input: func(v any) (any, error) {
			if v == nil {
				return map[string]any{"None": true}, nil
			}
			return map[string]any{"Some": v}, nil
		},
		Output: func(v any) (any, error) {
			m := v.(map[string]any)
			if m["$kind"] == "Some" {
				return m["Some"], nil
			}
			return nil, nil
		},
	})
}

func (bcsBuilder) Vector(t *Type, opts ...Option) *Type {
	op := pickOption(opts)
	name := coalesce(op.Name, fmt.Sprintf("vector<%s>", t.Name))
	return dynamicType(name,
		func(r *Reader) (any, error) {
			ln, err := r.ReadULEB()
			if err != nil {
				return nil, err
			}
			if ln > MaxSequenceLength {
				return nil, fmt.Errorf("sequence length %d exceeds limit %d", ln, MaxSequenceLength)
			}
			res := make([]any, 0, ln)
			for i := 0; i < int(ln); i++ {
				v, err := t.Read(r)
				if err != nil {
					return nil, err
				}
				res = append(res, v)
			}
			return res, nil
		},
		func(v any, w *Writer) error {
			arr, err := asAnySlice(v)
			if err != nil {
				return err
			}
			if len(arr) > MaxSequenceLength {
				return fmt.Errorf("sequence length %d exceeds limit %d", len(arr), MaxSequenceLength)
			}
			if err := w.WriteULEB(uint64(len(arr))); err != nil {
				return err
			}
			for _, el := range arr {
				if err := t.Write(el, w); err != nil {
					return err
				}
			}
			return nil
		},
		func(v any) (int, bool) {
			arr, err := asAnySlice(v)
			if err != nil {
				return 0, false
			}
			total := len(ULEBEncode(uint64(len(arr))))
			for _, el := range arr {
				s, ok := t.SerializedSize(el)
				if !ok {
					return 0, false
				}
				total += s
			}
			return total, true
		},
		op.Validate,
	)
}

func (bcsBuilder) Tuple(types []*Type, opts ...Option) *Type {
	op := pickOption(opts)
	name := op.Name
	if name == "" {
		for i, t := range types {
			if i == 0 {
				name = "("
			} else {
				name += ", "
			}
			name += t.Name
		}
		name += ")"
	}
	return dynamicType(name,
		func(r *Reader) (any, error) {
			out := make([]any, 0, len(types))
			for _, t := range types {
				v, err := t.Read(r)
				if err != nil {
					return nil, err
				}
				out = append(out, v)
			}
			return out, nil
		},
		func(v any, w *Writer) error {
			arr, err := asAnySlice(v)
			if err != nil || len(arr) != len(types) {
				return fmt.Errorf("expected tuple of length %d", len(types))
			}
			for i, t := range types {
				if err := t.Write(arr[i], w); err != nil {
					return err
				}
			}
			return nil
		},
		func(v any) (int, bool) {
			arr, err := asAnySlice(v)
			if err != nil || len(arr) != len(types) {
				return 0, false
			}
			total := 0
			for i, t := range types {
				s, ok := t.SerializedSize(arr[i])
				if !ok {
					return 0, false
				}
				total += s
			}
			return total, true
		},
		op.Validate,
	)
}

func (bcsBuilder) Struct(name string, fields []Field, opts ...Option) *Type {
	op := pickOption(opts)
	if op.Name != "" {
		name = op.Name
	}
	return dynamicType(name,
		func(r *Reader) (any, error) {
			if err := r.enterContainer(); err != nil {
				return nil, err
			}
			defer r.exitContainer()
			res := make(map[string]any, len(fields))
			for _, f := range fields {
				v, err := f.Type.Read(r)
				if err != nil {
					return nil, err
				}
				res[f.Name] = v
			}
			return res, nil
		},
		func(v any, w *Writer) error {
			if err := w.enterContainer(); err != nil {
				return err
			}
			defer w.exitContainer()
			m, ok := v.(map[string]any)
			if !ok {
				return fmt.Errorf("expected object for struct %s", name)
			}
			for _, f := range fields {
				if err := f.Type.Write(m[f.Name], w); err != nil {
					return err
				}
			}
			return nil
		},
		func(v any) (int, bool) {
			m, ok := v.(map[string]any)
			if !ok {
				return 0, false
			}
			total := 0
			for _, f := range fields {
				s, ok := f.Type.SerializedSize(m[f.Name])
				if !ok {
					return 0, false
				}
				total += s
			}
			return total, true
		},
		op.Validate,
	)
}

func (bcsBuilder) Enum(name string, fields []Field, opts ...Option) *Type {
	op := pickOption(opts)
	if op.Name != "" {
		name = op.Name
	}
	index := make(map[string]int, len(fields))
	for i, f := range fields {
		index[f.Name] = i
	}
	return dynamicType(name,
		func(r *Reader) (any, error) {
			if err := r.enterContainer(); err != nil {
				return nil, err
			}
			defer r.exitContainer()
			i, err := r.ReadULEB()
			if err != nil {
				return nil, err
			}
			if int(i) >= len(fields) {
				return nil, fmt.Errorf("unknown enum index %d for %s", i, name)
			}
			f := fields[int(i)]
			res := map[string]any{"$kind": f.Name}
			if f.Type == nil {
				res[f.Name] = true
				return res, nil
			}
			v, err := f.Type.Read(r)
			if err != nil {
				return nil, err
			}
			res[f.Name] = v
			return res, nil
		},
		func(v any, w *Writer) error {
			if err := w.enterContainer(); err != nil {
				return err
			}
			defer w.exitContainer()
			m, ok := v.(map[string]any)
			if !ok {
				return fmt.Errorf("expected enum object for %s", name)
			}
			var kind string
			if k, ok := m["$kind"].(string); ok {
				kind = k
			} else {
				for k := range m {
					if k != "$kind" {
						if _, ok := index[k]; ok {
							kind = k
							break
						}
					}
				}
			}
			i, ok := index[kind]
			if !ok {
				return fmt.Errorf("invalid enum variant %q", kind)
			}
			if err := w.WriteULEB(uint64(i)); err != nil {
				return err
			}
			f := fields[i]
			if f.Type != nil {
				return f.Type.Write(m[f.Name], w)
			}
			return nil
		},
		func(v any) (int, bool) {
			m, ok := v.(map[string]any)
			if !ok {
				return 0, false
			}
			kind, _ := m["$kind"].(string)
			if kind == "" {
				for k := range m {
					if _, ok := index[k]; ok {
						kind = k
						break
					}
				}
			}
			i, ok := index[kind]
			if !ok {
				return 0, false
			}
			total := len(ULEBEncode(uint64(i)))
			if fields[i].Type == nil {
				return total, true
			}
			s, ok := fields[i].Type.SerializedSize(m[kind])
			if !ok {
				return 0, false
			}
			return total + s, true
		},
		op.Validate,
	)
}

func (bcsBuilder) Map(keyType, valueType *Type) *Type {
	name := fmt.Sprintf("Map<%s, %s>", keyType.Name, valueType.Name)
	return dynamicType(name,
		func(r *Reader) (any, error) {
			ln, err := r.ReadULEB()
			if err != nil {
				return nil, err
			}
			if ln > MaxSequenceLength {
				return nil, fmt.Errorf("sequence length %d exceeds limit %d", ln, MaxSequenceLength)
			}
			out := make([][2]any, 0, ln)
			var prev []byte
			for i := 0; i < int(ln); i++ {
				k, err := keyType.Read(r)
				if err != nil {
					return nil, err
				}
				sk, err := keyType.Serialize(k, nil)
				if err != nil {
					return nil, err
				}
				keyBytes := sk.ToBytes()
				if i > 0 && CompareBCSBytes(prev, keyBytes) >= 0 {
					return nil, fmt.Errorf("map keys are not strictly increasing at index %d", i)
				}
				prev = keyBytes
				v, err := valueType.Read(r)
				if err != nil {
					return nil, err
				}
				out = append(out, [2]any{k, v})
			}
			return out, nil
		},
		func(v any, w *Writer) error {
			entries, err := asEntries(v)
			if err != nil {
				return err
			}
			if len(entries) > MaxSequenceLength {
				return fmt.Errorf("sequence length %d exceeds limit %d", len(entries), MaxSequenceLength)
			}
			type kv struct {
				k  []byte
				vv any
			}
			serialized := make([]kv, 0, len(entries))
			for _, e := range entries {
				sk, err := keyType.Serialize(e[0], nil)
				if err != nil {
					return err
				}
				serialized = append(serialized, kv{k: sk.ToBytes(), vv: e[1]})
			}
			sort.Slice(serialized, func(i, j int) bool {
				return CompareBCSBytes(serialized[i].k, serialized[j].k) < 0
			})
			if err := w.WriteULEB(uint64(len(serialized))); err != nil {
				return err
			}
			for _, item := range serialized {
				if err := w.WriteBytes(item.k); err != nil {
					return err
				}
				if err := valueType.Write(item.vv, w); err != nil {
					return err
				}
			}
			return nil
		},
		func(v any) (int, bool) {
			entries, err := asEntries(v)
			if err != nil {
				return 0, false
			}
			total := len(ULEBEncode(uint64(len(entries))))
			for _, e := range entries {
				sk, ok := keyType.SerializedSize(e[0])
				if !ok {
					return 0, false
				}
				sv, ok := valueType.SerializedSize(e[1])
				if !ok {
					return 0, false
				}
				total += sk + sv
			}
			return total, true
		},
		nil,
	)
}

func (bcsBuilder) Lazy(cb func() *Type) *Type {
	var cached *Type
	get := func() *Type {
		if cached == nil {
			cached = cb()
		}
		return cached
	}
	return &Type{
		Name:           "lazy",
		Read:           func(r *Reader) (any, error) { return get().Read(r) },
		Write:          func(v any, w *Writer) error { return get().Write(v, w) },
		SerializedSize: func(v any) (int, bool) { return get().SerializedSize(v) },
		Validate: func(v any) error {
			if get().Validate != nil {
				return get().Validate(v)
			}
			return nil
		},
	}
}

func pickOption(opts []Option) Option {
	if len(opts) == 0 {
		return Option{}
	}
	return opts[0]
}
func coalesce(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}

func toUint64(v any) (uint64, bool) {
	switch t := v.(type) {
	case uint8:
		return uint64(t), true
	case uint16:
		return uint64(t), true
	case uint32:
		return uint64(t), true
	case uint64:
		return t, true
	case uint:
		return uint64(t), true
	case int:
		if t < 0 {
			return 0, false
		}
		return uint64(t), true
	case int64:
		if t < 0 {
			return 0, false
		}
		return uint64(t), true
	case float64:
		if t < 0 || t != math.Trunc(t) {
			return 0, false
		}
		return uint64(t), true
	case string:
		bi, ok := new(big.Int).SetString(t, 10)
		if !ok || bi.Sign() < 0 || bi.BitLen() > 64 {
			return 0, false
		}
		return bi.Uint64(), true
	default:
		return 0, false
	}
}

func toBigInt(v any) (*big.Int, error) {
	switch t := v.(type) {
	case *big.Int:
		return new(big.Int).Set(t), nil
	case uint64:
		return new(big.Int).SetUint64(t), nil
	case int64:
		return big.NewInt(t), nil
	case int:
		return big.NewInt(int64(t)), nil
	case string:
		bi, ok := new(big.Int).SetString(t, 10)
		if !ok {
			return nil, fmt.Errorf("invalid bigint string: %q", t)
		}
		return bi, nil
	default:
		return nil, fmt.Errorf("invalid bigint type: %T", v)
	}
}

func littleToBigInt(le []byte) *big.Int {
	be := make([]byte, len(le))
	for i := range le {
		be[len(le)-1-i] = le[i]
	}
	return new(big.Int).SetBytes(be)
}

func bigIntToLittleFixed(v *big.Int, size int) []byte {
	be := v.Bytes()
	le := make([]byte, size)
	for i := 0; i < len(be) && i < size; i++ {
		le[i] = be[len(be)-1-i]
	}
	return le
}

func asBytes(v any) ([]byte, error) {
	switch t := v.(type) {
	case []byte:
		return append([]byte(nil), t...), nil
	case string:
		return DecodeStr(t, EncodingHex)
	case [16]byte:
		return t[:], nil
	case [32]byte:
		return t[:], nil
	default:
		return nil, fmt.Errorf("expected bytes")
	}
}

func asAnySlice(v any) ([]any, error) {
	switch t := v.(type) {
	case []any:
		return t, nil
	case []string:
		out := make([]any, len(t))
		for i := range t {
			out[i] = t[i]
		}
		return out, nil
	case []int:
		out := make([]any, len(t))
		for i := range t {
			out[i] = t[i]
		}
		return out, nil
	case [][]byte:
		out := make([]any, len(t))
		for i := range t {
			out[i] = t[i]
		}
		return out, nil
	default:
		return nil, fmt.Errorf("expected slice")
	}
}

func asEntries(v any) ([][2]any, error) {
	switch t := v.(type) {
	case [][2]any:
		return t, nil
	case map[string]any:
		out := make([][2]any, 0, len(t))
		for k, vv := range t {
			out = append(out, [2]any{k, vv})
		}
		return out, nil
	default:
		return nil, fmt.Errorf("expected map entries")
	}
}
