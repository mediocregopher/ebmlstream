package edtd

import (
	"bytes"
	. "testing"
	"github.com/stretchr/testify/assert"
)

// TODO this needs tests
func mustDefDataBytes(d interface{}) []byte {
	b, err := defDataBytes(d)
	if err != nil {
		panic(err)
	}
	return b
}

// This is what implicitHeader should parse to
var implicitEBML = &tplElement{
	id: 0xa45dfa3,
	typ: Container,
	name: "EBML",
	card: oneOrMore,
	kids: []tplElement{
		{
			id: 0x286,
			typ: Uint,
			name: "EBMLVersion",
			def: mustDefDataBytes(uint64(1)),
		},
		{
			id: 0x2f7,
			typ: Uint,
			name: "EBMLReadVersion",
			def: mustDefDataBytes(uint64(1)),
		},
		{
			id: 0x2f2,
			typ: Uint,
			name: "EBMLMaxIDLength",
			def: mustDefDataBytes(uint64(4)),
		},
		{
			id: 0x2f3,
			typ: Uint,
			name: "EBMLMaxSizeLength",
			def: mustDefDataBytes(uint64(8)),
		},
		{
			id: 0x282,
			typ: String,
			name: "DocType",
			ranges: &rangeParam{loweri: 32, upperi: 126},
		},
		{
			id: 0x287,
			typ: Uint,
			name: "DocTypeVersion",
			def: mustDefDataBytes(uint64(1)),
		},
		{
			id: 0x285,
			typ: Uint,
			name: "DocTypeReadVersion",
			def: mustDefDataBytes(uint64(1)),
		},
	},
}

var implicitCRC32 = &tplElement{
	id: 0x43,
	typ: Container,
	name: "CRC32",
	card: zeroOrMore,
	kids: []tplElement{
		{
			id: 0x2fe,
			typ: Binary,
			name: "CRC32Value",
			size: 4,
		},
	},
}

var implicitVoid = &tplElement{
	id: 0x6c,
	typ: Binary,
	name: "Void",
	card: zeroOrMore,
}

func TestParseImplicitElements(t *T) {
	i := elementIndex{}
	lex := newLexer(bytes.NewBufferString(implicitElements))
	_, err := parseElements(lex, i)
	assert.Nil(t, err)

	ebml := i[elementID(0xa45dfa3)]
	assert.Equal(t, implicitEBML, ebml)

	crc32 := i[elementID(0x43)]
	assert.Equal(t, implicitCRC32, crc32)

	void := i[elementID(0x6c)]
	assert.Equal(t, implicitVoid, void)
}
