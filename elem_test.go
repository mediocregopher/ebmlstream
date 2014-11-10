package ebml

import (
	"bytes"
	"io"
	"os"
	. "testing"
	"github.com/stretchr/testify/assert"
	
	"github.com/mediocregopher/go.ebml/varint"
)

func sb(bs ...byte) string {
	b := make([]byte, 0, len(bs))
	b = append(b, bs...)
	return string(b)
}

func TestIntElem(t *T) {
	m := map[string]int64{
		sb(0x80, 0x80):                   0,
		sb(0x80, 0x81, 0x01):             0x01,
		sb(0x80, 0x82, 0x02, 0x01):       0x0201,
		sb(0x80, 0x83, 0x03, 0x02, 0x01): 0x030201,
	}
	assert := assert.New(t)
	for in, out := range m {
		b := []byte(in)
		e, err := RootElem(bytes.NewBuffer(b)).Next()
		assert.Nil(err, "input: %x", in)

		i, err := e.Int()
		assert.Nil(err, "input: %x", in)
		assert.Exactly(out, i, "input: %x", in)

		ui, err := e.Uint()
		assert.Nil(err, "input: %x", in)
		assert.Exactly(uint64(out), ui, "input: %x", in)
	}
}

func TestStringElem(t *T) {
	m := map[string]string{
		sb(0x80, 0x80):                      "",
		sb(0x80, 0x83, 'f', 'o', 'o'):       "foo",
		sb(0x80, 0x85, 'f', 'o', 'o', 0, 0): "foo",
	}
	assert := assert.New(t)
	for in, out := range m {
		b := []byte(in)
		e, err := RootElem(bytes.NewBuffer(b)).Next()
		assert.Nil(err, "input: %x", in)

		i, err := e.String()
		assert.Nil(err, "input: %x", in)
		assert.Equal(out, i, "input: %x", in)
	}
}

func getTestFile(t *T) io.ReadCloser {
	f, err := os.Open("test.webm")
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func TestFile1(t *T) {
	f := getTestFile(t)
	defer f.Close()
	assert := assert.New(t)

	e := RootElem(f)
	id, err := varint.ReadVarInt(e.buf)
	assert.Nil(err)
	assert.Equal(0xa45dfa3, id)

	size, err := varint.ReadVarInt(e.buf)
	assert.Nil(err)
	assert.Equal(0x23, size)
}
