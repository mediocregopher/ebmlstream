package ebml

import (
	"bytes"
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

	for in, out := range m {
		if n := numPrecedingZeros(in); n != out {
			t.Fatalf("(%x) %d != %d", in, n, out)
		}
	}
}

func sb(bs ...byte) string {
	b := make([]byte, 0, len(bs))
	b = append(b, bs...)
	return string(b)
}

func TestReadVarInt(t *T) {
	m := map[string]int64{
		sb(0x81):                                           0x01,
		sb(0xC1):                                           0x41,
		sb(0x40, 0x01):                                     0x01,
		sb(0x41, 0x21):                                     0x0121,
		sb(0x20, 0x41, 0x21):                               0x4121,
		sb(0x23, 0x41, 0x21):                               0x034121,
		sb(0x01, 0x41, 0x21, 0x12, 0x34, 0x56, 0x78, 0x9a): 0x4121123456789a,
	}

	for in, out := range m {
		e := NewEBMLReader(bytes.NewBuffer([]byte(in)))
		i, err := e.readVarInt()
		if err != nil {
			t.Fatal(err)
		}
		if i != out {
			t.Fatalf("(%x) %d != %d", in, i, out)
		}
	}
}
