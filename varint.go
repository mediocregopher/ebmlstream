package ebml

import (
	"io"
	"bufio"
)

type EBMLReader struct {
	r io.Reader
	buf *bufio.Reader
}

func NewEBMLReader(r io.Reader) *EBMLReader {
	return &EBMLReader{
		r:   r,
		buf: bufio.NewReader(r),
	}
}

func numPrecedingZeros(b byte) byte {
	for i := byte(0); i < 8; i++ {
		if b & 0x80 > 0 {
			return i
		}
		b <<= 1
	}
	return 8
}

func (e *EBMLReader) readVarInt() (int64, error) {
	b, err := e.buf.ReadByte()
	if err != nil {
		return 0, err
	}

	rem := numPrecedingZeros(b)
	ret := int64(b & (0xFF >> (rem + 1)))
	for ; rem > 0; rem-- {
		b, err = e.buf.ReadByte()
		if err != nil {
			return 0, err
		}
		ret = (ret << 8) | int64(b)
	}

	return ret, nil
}
