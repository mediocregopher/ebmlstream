package edtd

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	. "testing"
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
	id:   0xa45dfa3,
	typ:  Container,
	name: "EBML",
	card: oneOrMore,
	kids: []tplElement{
		{
			id:   0x286,
			typ:  Uint,
			name: "EBMLVersion",
			def:  mustDefDataBytes(uint64(1)),
		},
		{
			id:   0x2f7,
			typ:  Uint,
			name: "EBMLReadVersion",
			def:  mustDefDataBytes(uint64(1)),
		},
		{
			id:   0x2f2,
			typ:  Uint,
			name: "EBMLMaxIDLength",
			def:  mustDefDataBytes(uint64(4)),
		},
		{
			id:   0x2f3,
			typ:  Uint,
			name: "EBMLMaxSizeLength",
			def:  mustDefDataBytes(uint64(8)),
		},
		{
			id:     0x282,
			typ:    String,
			name:   "DocType",
			ranges: &rangeParam{loweri: 32, upperi: 126},
		},
		{
			id:   0x287,
			typ:  Uint,
			name: "DocTypeVersion",
			def:  mustDefDataBytes(uint64(1)),
		},
		{
			id:   0x285,
			typ:  Uint,
			name: "DocTypeReadVersion",
			def:  mustDefDataBytes(uint64(1)),
		},
	},
}

var implicitCRC32 = &tplElement{
	id:   0x43,
	typ:  Container,
	name: "CRC32",
	card: zeroOrMore,
	kids: []tplElement{
		{
			id:   0x2fe,
			typ:  Binary,
			name: "CRC32Value",
			size: 4,
		},
	},
}

var implicitVoid = &tplElement{
	id:   0x6c,
	typ:  Binary,
	name: "Void",
	card: zeroOrMore,
}

func TestParseImplicitElements(t *T) {
	m := elementMap{}
	tm := typesMap{}
	lex := newLexer(bytes.NewBufferString(implicitElements))
	_, err := parseElements(lex, m, tm, false)
	require.Nil(t, err)

	ebml := m[elementID(0xa45dfa3)]
	assert.Equal(t, implicitEBML, ebml)

	crc32 := m[elementID(0x43)]
	assert.Equal(t, implicitCRC32, crc32)

	void := m[elementID(0x6c)]
	assert.Equal(t, implicitVoid, void)
}

func TestParseTypes(t *T) {

	test := `
        define types {
            bool := uint [ range:0..1; ]
            ascii := string [ range:32..126; ]
        }

        define elements {
		    Foo := 53ab bool [ def:1; ]
			Bar := 53ac bool [ card:?; ]
		}
	`

	e, err := NewEdtd(bytes.NewBufferString(test))
	require.Nil(t, err)

	boolRange := &rangeParam{
		lowerui: 0,
		upperui: 1,
	}

	foo := &tplElement{
		id:     0x13ab,
		typ:    Uint,
		name:   "Foo",
		def:    mustDefDataBytes(uint64(1)),
		ranges: boolRange,
	}
	assert.Equal(t, foo, e.elements[0x13ab])

	bar := &tplElement{
		id:     0x13ac,
		typ:    Uint,
		name:   "Bar",
		card:   zeroOrOnce,
		ranges: boolRange,
	}
	assert.Equal(t, bar, e.elements[0x13ac])
}
