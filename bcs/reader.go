package bcs

import (
	"encoding/binary"
	"fmt"
)

type Reader struct {
	data            []byte
	pos             int
	containerBudget int
}

func NewReader(data []byte) *Reader {
	return &Reader{data: data, containerBudget: MaxContainerDepth}
}

func (r *Reader) Shift(bytes int) *Reader {
	r.pos += bytes
	return r
}

func (r *Reader) remaining() int {
	return len(r.data) - r.pos
}

func (r *Reader) Remaining() int {
	return r.remaining()
}

func (r *Reader) enterContainer() error {
	if r.containerBudget == 0 {
		return fmt.Errorf("bcs: exceeded max container depth")
	}
	r.containerBudget--
	return nil
}

func (r *Reader) exitContainer() {
	r.containerBudget++
}

func (r *Reader) Read8() (uint8, error) {
	if r.remaining() < 1 {
		return 0, fmt.Errorf("bcs: out of range")
	}
	v := r.data[r.pos]
	r.pos++
	return v, nil
}

func (r *Reader) Read16() (uint16, error) {
	if r.remaining() < 2 {
		return 0, fmt.Errorf("bcs: out of range")
	}
	v := binary.LittleEndian.Uint16(r.data[r.pos:])
	r.pos += 2
	return v, nil
}

func (r *Reader) Read32() (uint32, error) {
	if r.remaining() < 4 {
		return 0, fmt.Errorf("bcs: out of range")
	}
	v := binary.LittleEndian.Uint32(r.data[r.pos:])
	r.pos += 4
	return v, nil
}

func (r *Reader) Read64() (uint64, error) {
	if r.remaining() < 8 {
		return 0, fmt.Errorf("bcs: out of range")
	}
	v := binary.LittleEndian.Uint64(r.data[r.pos:])
	r.pos += 8
	return v, nil
}

func (r *Reader) Read128() ([16]byte, error) {
	var out [16]byte
	if r.remaining() < 16 {
		return out, fmt.Errorf("bcs: out of range")
	}
	copy(out[:], r.data[r.pos:r.pos+16])
	r.pos += 16
	return out, nil
}

func (r *Reader) Read256() ([32]byte, error) {
	var out [32]byte
	if r.remaining() < 32 {
		return out, fmt.Errorf("bcs: out of range")
	}
	copy(out[:], r.data[r.pos:r.pos+32])
	r.pos += 32
	return out, nil
}

func (r *Reader) ReadBytes(n int) ([]byte, error) {
	if n < 0 || r.remaining() < n {
		return nil, fmt.Errorf("bcs: out of range")
	}
	out := make([]byte, n)
	copy(out, r.data[r.pos:r.pos+n])
	r.pos += n
	return out, nil
}

func (r *Reader) ReadULEB() (uint64, error) {
	v, n, err := ULEBDecode(r.data[r.pos:])
	if err != nil {
		return 0, err
	}
	r.pos += n
	return v, nil
}

func (r *Reader) ReadVec(cb func(*Reader, int, int) (any, error)) ([]any, error) {
	ln, err := r.ReadULEB()
	if err != nil {
		return nil, err
	}
	if ln > MaxSequenceLength {
		return nil, fmt.Errorf("sequence length %d exceeds limit %d", ln, MaxSequenceLength)
	}
	res := make([]any, 0, ln)
	for i := 0; i < int(ln); i++ {
		v, err := cb(r, i, int(ln))
		if err != nil {
			return nil, err
		}
		res = append(res, v)
	}
	return res, nil
}
