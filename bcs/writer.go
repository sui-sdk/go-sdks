package bcs

import (
	"encoding/binary"
	"fmt"
)

type WriterOptions struct {
	InitialSize  int
	MaxSize      int
	AllocateSize int
}

type Writer struct {
	buf          []byte
	pos          int
	maxSize      int
	allocateSize int
}

func NewWriter(opts *WriterOptions) *Writer {
	initial := 1024
	maxSize := int(^uint(0) >> 1)
	alloc := 1024
	if opts != nil {
		if opts.InitialSize > 0 {
			initial = opts.InitialSize
		}
		if opts.MaxSize > 0 {
			maxSize = opts.MaxSize
		}
		if opts.AllocateSize > 0 {
			alloc = opts.AllocateSize
		}
	}
	return &Writer{buf: make([]byte, initial), maxSize: maxSize, allocateSize: alloc}
}

func (w *Writer) ensure(n int) error {
	req := w.pos + n
	if req <= len(w.buf) {
		return nil
	}
	next := len(w.buf)
	for next < req {
		next += maxInt(w.allocateSize, req)
		if next > w.maxSize {
			next = w.maxSize
			break
		}
	}
	if req > next {
		return fmt.Errorf("attempting to serialize to BCS, required size %d exceeds max size %d", req, w.maxSize)
	}
	nb := make([]byte, next)
	copy(nb, w.buf)
	w.buf = nb
	return nil
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (w *Writer) Shift(bytes int) *Writer {
	w.pos += bytes
	return w
}

func (w *Writer) Write8(v uint8) error {
	if err := w.ensure(1); err != nil {
		return err
	}
	w.buf[w.pos] = v
	w.pos++
	return nil
}

func (w *Writer) WriteBytes(bs []byte) error {
	if err := w.ensure(len(bs)); err != nil {
		return err
	}
	copy(w.buf[w.pos:], bs)
	w.pos += len(bs)
	return nil
}

func (w *Writer) Write16(v uint16) error {
	if err := w.ensure(2); err != nil {
		return err
	}
	binary.LittleEndian.PutUint16(w.buf[w.pos:], v)
	w.pos += 2
	return nil
}

func (w *Writer) Write32(v uint32) error {
	if err := w.ensure(4); err != nil {
		return err
	}
	binary.LittleEndian.PutUint32(w.buf[w.pos:], v)
	w.pos += 4
	return nil
}

func (w *Writer) Write64(v uint64) error {
	if err := w.ensure(8); err != nil {
		return err
	}
	binary.LittleEndian.PutUint64(w.buf[w.pos:], v)
	w.pos += 8
	return nil
}

func (w *Writer) Write128(v [16]byte) error { return w.WriteBytes(v[:]) }
func (w *Writer) Write256(v [32]byte) error { return w.WriteBytes(v[:]) }

func (w *Writer) WriteULEB(v uint64) error {
	for _, b := range ULEBEncode(v) {
		if err := w.Write8(b); err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) WriteVec(vec []any, cb func(*Writer, any, int, int) error) error {
	if err := w.WriteULEB(uint64(len(vec))); err != nil {
		return err
	}
	for i, el := range vec {
		if err := cb(w, el, i, len(vec)); err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) ToBytes() []byte {
	out := make([]byte, w.pos)
	copy(out, w.buf[:w.pos])
	return out
}

func (w *Writer) ToString(encoding Encoding) (string, error) {
	return EncodeStr(w.ToBytes(), encoding)
}
