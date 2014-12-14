// Implementation of the variable-sized integer for the ebml document type.
// varints are integers which are encoded in such a way as to take up only as
// many bytes as needed (smaller integers take up fewer bytes). This is done
// through a utf-8 style method where the number of preceding zero bits in the
// first byte of the integer indicates how many further bytes need to be read to
// encompass the full integer
package varint

import (
	"bytes"
	"io"
)

func numPrecedingZeros(b byte) byte {
	for i := byte(0); i < 8; i++ {
		if b&0x80 > 0 {
			return i
		}
		b <<= 1
	}
	return 8
}

func readByte(r io.Reader) (byte, error) {
	b := make([]byte, 1)
	_, err := io.ReadFull(r, b)
	return b[0], err
}

// Reads a variable integer from the given reader, reading only as many bytes as
// necessary
func ReadVarInt(r io.Reader) (int64, error) {

	b, err := readByte(r)
	if err != nil {
		return 0, err
	}

	rem := numPrecedingZeros(b)
	ret := int64(b & (0xFF >> (rem + 1)))
	for ; rem > 0; rem-- {
		b, err = readByte(r)
		if err != nil {
			return 0, err
		}
		ret = (ret << 8) | int64(b)
	}

	return ret, nil
}

// Reads a varint from the given slice of bytes. The slice of bytes must have
// enough bytes to encompass the full varint, but having more bytes than
// necessary is ok
func VarInt(b []byte) (int64, error) {
	buf := bytes.NewBuffer(b)
	return ReadVarInt(buf)
}

func WriteVarInt(i int64, w io.Writer) (int, error) {

	// Count the number of bits actually being used by the number
	bits := byte(0)
	for ic := i; ; {
		ic >>= 1
		bits++
		if ic == 0 {
			break
		}
	}

	// Find the number of bytes the raw number uses. We round up using the cool
	// addition hack
	bytes := (bits + 7) / 8

	// If we have enough bits available at the front of the raw number to simply
	// drop the size bits there we don't need to add another byte to the front
	// of our encoded number
	bitsAvailable := (bytes * 8) - bits
	bitsNeeded := bytes
	if bitsNeeded > bitsAvailable {
		bytes++
	}

	one := 0x80 >> (bytes - 1)
	shifted := int64(one) << ((bytes - 1) * 8)
	newI := i | shifted

	// newI is the encoded form of our number, now to write it to the io.Writer
	// using the minimum number of bytes
	out := make([]byte, 0, bytes)
	for j := bytes - 1; j >= 0 && j < 255; j-- {
		b := byte((newI >> (j * 8)) & 0xff)
		out = append(out, b)
	}
	return w.Write(out)
}

func ToVarInt(i int64) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 8))
	if _, err := WriteVarInt(i, buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
