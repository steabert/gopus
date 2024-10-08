package binary

import (
	"encoding/binary"
	"io"
)

type reader struct {
	r   io.Reader
	err error
}

func NewReader(r io.Reader) *reader {
	return &reader{r: r}
}

func (r *reader) Err() error {
	return r.err
}

func (r *reader) ReadUint8() uint8 {
	if r.err != nil {
		return 0
	}
	b := [1]byte{}
	_, err := io.ReadFull(r.r, b[:])
	if err != nil {
		r.err = err
		return 0
	}
	return uint8(b[0])
}

func (r *reader) ReadUint16() uint16 {
	if r.err != nil {
		return 0
	}
	b := [2]byte{}
	_, err := io.ReadFull(r.r, b[:])
	if err != nil {
		r.err = err
		return 0
	}
	return binary.LittleEndian.Uint16(b[:])
}

func (r *reader) ReadUint32() uint32 {
	if r.err != nil {
		return 0
	}
	b := [4]byte{}
	_, err := io.ReadFull(r.r, b[:])
	if err != nil {
		r.err = err
		return 0
	}
	return binary.LittleEndian.Uint32(b[:])
}

func (r *reader) ReadUint64() uint64 {
	if r.err != nil {
		return 0
	}
	b := [8]byte{}
	_, err := io.ReadFull(r.r, b[:])
	if err != nil {
		r.err = err
		return 0
	}
	return binary.LittleEndian.Uint64(b[:])
}
