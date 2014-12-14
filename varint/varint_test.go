package varint

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	. "testing"
)

func TestNumPrecedingZeros(t *T) {
	m := map[byte]byte{
		0x80: 0,
		0x40: 1,
		0x20: 2,
		0x10: 3,
		0x08: 4,
		0x04: 5,
		0x02: 6,
		0x01: 7,
		0x00: 8,
	}

	assert := assert.New(t)
	for in, out := range m {
		assert.Equal(out, numPrecedingZeros(in), "input: %x", in)
	}
}

func sb(bs ...byte) string {
	b := make([]byte, 0, len(bs))
	b = append(b, bs...)
	return string(b)
}

var m = map[string]int64{
	sb(0x81):                                           0x01,
	sb(0xC1):                                           0x41,
	sb(0x41, 0x21):                                     0x0121,
	sb(0x53, 0xac):                                     0x13ac,
	sb(0x20, 0x41, 0x21):                               0x4121,
	sb(0x23, 0x41, 0x21):                               0x034121,
	sb(0x03, 0x21, 0x12, 0x34, 0x56, 0x78, 0x9a):       0x0121123456789a,
	sb(0x01, 0x41, 0x21, 0x12, 0x34, 0x56, 0x78, 0x9a): 0x4121123456789a,
}

func TestReadVarInt(t *T) {
	assert := assert.New(t)
	for in, out := range m {
		i, err := ReadVarInt(bytes.NewBuffer([]byte(in)))
		assert.Nil(err, "input: 0x%x", in)
		assert.Equal(out, i, "input: 0x%x", in)
	}
}

func TestWriteVarInt(t *T) {
	assert := assert.New(t)
	for out, in := range m {
		w := bytes.NewBuffer([]byte{})
		_, err := WriteVarInt(in, w)
		require.Nil(t, err, "input: 0x%x", in)
		assert.Equal(out, w.String(), "input 0x%x out 0x%x expected 0x%x", in, w.String(), out)
	}
}
