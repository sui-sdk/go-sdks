package bcs

import (
	"fmt"
)

type Type struct {
	Name           string
	Read           func(*Reader) (any, error)
	Write          func(any, *Writer) error
	SerializedSize func(any) (int, bool)
	Validate       func(any) error
}

type TransformOptions struct {
	Name     string
	Input    func(any) (any, error)
	Output   func(any) (any, error)
	Validate func(any) error
}

func (t *Type) Serialize(value any, options *WriterOptions) (*Serialized, error) {
	if t.Validate != nil {
		if err := t.Validate(value); err != nil {
			return nil, err
		}
	}
	writer := NewWriter(options)
	if t.Write == nil {
		return nil, fmt.Errorf("bcs type %s has no writer", t.Name)
	}
	if err := t.Write(value, writer); err != nil {
		return nil, err
	}
	return &Serialized{schema: t, bytes: writer.ToBytes()}, nil
}

func (t *Type) Parse(bytes []byte) (any, error) {
	if t.Read == nil {
		return nil, fmt.Errorf("bcs type %s has no reader", t.Name)
	}
	return t.Read(NewReader(bytes))
}

func (t *Type) FromHex(hex string) (any, error) {
	b, err := DecodeStr(hex, EncodingHex)
	if err != nil {
		return nil, err
	}
	return t.Parse(b)
}

func (t *Type) FromBase64(v string) (any, error) {
	b, err := DecodeStr(v, EncodingBase64)
	if err != nil {
		return nil, err
	}
	return t.Parse(b)
}

func (t *Type) FromBase58(v string) (any, error) {
	b, err := DecodeStr(v, EncodingBase58)
	if err != nil {
		return nil, err
	}
	return t.Parse(b)
}

func (t *Type) Transform(opts TransformOptions) *Type {
	name := t.Name
	if opts.Name != "" {
		name = opts.Name
	}
	return &Type{
		Name: name,
		Read: func(r *Reader) (any, error) {
			v, err := t.Read(r)
			if err != nil {
				return nil, err
			}
			if opts.Output == nil {
				return v, nil
			}
			return opts.Output(v)
		},
		Write: func(v any, w *Writer) error {
			if opts.Input != nil {
				nv, err := opts.Input(v)
				if err != nil {
					return err
				}
				v = nv
			}
			return t.Write(v, w)
		},
		SerializedSize: func(v any) (int, bool) {
			if opts.Input != nil {
				nv, err := opts.Input(v)
				if err != nil {
					return 0, false
				}
				v = nv
			}
			if t.SerializedSize == nil {
				return 0, false
			}
			return t.SerializedSize(v)
		},
		Validate: func(v any) error {
			if opts.Validate != nil {
				if err := opts.Validate(v); err != nil {
					return err
				}
			}
			if opts.Input != nil {
				nv, err := opts.Input(v)
				if err != nil {
					return err
				}
				v = nv
			}
			if t.Validate != nil {
				return t.Validate(v)
			}
			return nil
		},
	}
}

type Serialized struct {
	schema *Type
	bytes  []byte
}

func (s *Serialized) ToBytes() []byte {
	out := make([]byte, len(s.bytes))
	copy(out, s.bytes)
	return out
}

func (s *Serialized) ToHex() (string, error) { return EncodeStr(s.bytes, EncodingHex) }
func (s *Serialized) ToBase64() (string, error) { return EncodeStr(s.bytes, EncodingBase64) }
func (s *Serialized) ToBase58() (string, error) { return EncodeStr(s.bytes, EncodingBase58) }

func (s *Serialized) Parse() (any, error) { return s.schema.Parse(s.bytes) }
