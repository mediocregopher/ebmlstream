package edtd

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"strings"
	
	"github.com/mediocregopher/go.ebml/varint"
)

type Type int
const (
	Int Type = iota
	Uint
	Float
	String
	Date
	Binary
	Container
)

type card int
const (
	zeroOrOnce card = iota
	zeroOrMore
	exactlyOnce
	oneOrMore
)

type (
	elementID    int64
	elementIndex map[elementID]*tplElement
)

type tplElement struct {
	id   elementID
	typ  Type
	name string
	kids []tplElement
	def  []byte
	card
}

// Parser is generated from an edtd specification. It will read in streams of
// ebml data and attempt to parse them based on its edtd.
type Parser struct {
	elementIndex
}

var implicitHeader = `
   define elements {
     EBML := 1a45dfa3 container [ card:+; ] {
       EBMLVersion := 4286 uint [ def:1; ]
       EBMLReadVersion := 42f7 uint [ def:1; ]
       EBMLMaxIDLength := 42f2 uint [ def:4; ]
       EBMLMaxSizeLength := 42f3 uint [ def:8; ]
       DocType := 4282 string [ range:32..126; ]
       DocTypeVersion := 4287 uint [ def:1; ]
       DocTypeReadVersion := 4285 uint [ def:1; ]
     }
 
     //TODO this nonsense
     //CRC32 := c3 container [ level:1..; card:*; ] {
     //  %children;
     //  CRC32Value := 42fe binary [ size:4; ]
     //}
     Void  := ec binary [ level:1..; card:*; ]
   }
`

var (
	defineTok    = token{alphaNum, "define"}
	elementsTok  = token{alphaNum, "elements"}
	assignTok    = token{control,  ":="}
	openCurlyTok = token{control, "{"}
	colonTok     = token{control, ":"}
)

// Pulls the next token from the lexer and checks if it matches any of the
// tokens, returning the one it matches or an error
func expect(lex *lexer, tok ...*token) (*token, error) {
	nextTok := lex.next()
	for i := range tok {
		if *nextTok == *tok[i] {
			return tok[i], nil
		}
	}
	return nil, fmt.Errorf("expected one of %v but found '%s'", tok, nextTok)
}

// Pulls the next token from the lexer and checks if it is of the given type,
// returning it if it is or an error
func expectType(lex *lexer, typ ...tokentyp) (*token, error) {
	nextTok := lex.next()
	for i := range typ {
		if nextTok.typ == typ[i] {
			return nextTok, nil
		}
	}
	return nil, fmt.Errorf("unexpected token '%s' found", nextTok)
}

func parseAsRoot(r io.Reader) (elementIndex, error) {
	lex := newLexer(r)
	m := elementIndex{}
	
	if _, err := expect(lex, &defineTok); err != nil {
		return nil, err
	}

	defWhat, err := expect(lex, &elementsTok)
	if err != nil {
		return nil, err
	}

	if _, err = expect(lex, &openCurlyTok); err != nil {
		return nil, err
	}

	switch defWhat.val {
	case "elements":
		parseElements(lex, m)
	}
	return nil, nil // TODO obviously take this out
}

func parseElements(lex *lexer, m elementIndex) ([]tplElement, error) {
	elems := make([]tplElement, 0, 8)
	for {
		elem, err, done := parseElement(lex, m)
		if err != nil {
			return nil, err
		} else if done {
			return elems, nil
		}
		elems = append(elems, elem)
	}
}

// Parses a single element out and returns it. Since this is expected to run
// within a container to parse out the elements within a container, it also
// handles running into a closing curly brace (the end of the container). In
// thise case it returns the third argument as true and doesn't parse anything
// out
func parseElement(lex *lexer, m elementIndex) (tplElement, error, bool) {
	nameTok := lex.next()
	if err := nameTok.asError(); err != nil {
		return tplElement{}, err, false
	} else if nameTok.typ == control && nameTok.val == "}" {
		return tplElement{}, nil, true
	} else if nameTok.typ != alphaNum {
		return tplElement{}, fmt.Errorf("unexpected '%s' found", nameTok), false
	}

	if _, err := expect(lex, &assignTok); err != nil {
		return tplElement{}, err, false
	}

	idTok, err := expectType(lex, alphaNum)
	if err != nil {
		return tplElement{}, err, false
	}

	id, err := strToID(idTok.val)
	if err != nil {
		return tplElement{}, err, false
	}

	typTok, err := expectType(lex, alphaNum)
	if err != nil {
		return tplElement{}, err, false
	}

	typ, err := strToType(typTok.val)
	if err != nil {
		return tplElement{}, err, false
	}

	elem := tplElement{
		id:   id,
		typ:  typ,
		name: nameTok.val,
	}

	controlTok, err := expectType(lex, control)
	if err != nil {
		return elem, err, false
	}

	// TODO parse params

	if controlTok.val != ";" {
		return elem, fmt.Errorf("found unexpected '%s'", controlTok), false
	}

	// TODO children

	m[elem.id] = &elem
	return elem, nil, false
}

func strToID(s string) (elementID, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return 0, err
	}
	i, err := varint.VarInt(b)
	return elementID(i), err
}

func strToType(s string) (Type, error) {
	switch strings.ToLower(s) {
	case "int":
		return Int, nil
	case "uint":
		return Uint, nil
	case "float":
		return Float, nil
	case "string":
		return String, nil
	case "date":
		return Date, nil
	case "binary":
		return Binary, nil
	case "container":
		return Container, nil
	default:
		return 0, fmt.Errorf("unknown type '%s'", s)
	}
}

func parseParams(lex *lexer, elem *tplElement) error {
	pnameTok, err := expectType(lex, alphaNum)
	if err != nil {
		return err
	}

	if _, err = expect(lex, &colonTok); err != nil {
		return err
	}

	pvalTok, err := expectType(lex, alphaNum, quotedString)
	if err != nil {
		return err
	}

	switch pnameTok.val {
	case "card":
		return parseCardParam(lex, elem, pvalTok)
	case "def":
		return parseDefParam(lex, elem, pvalTok)
	}

	return nil
}

func parseCardParam(lex *lexer, elem *tplElement, pvalTok *token) error {
	switch pvalTok.val {
	case "*":
		elem.card = zeroOrMore
	case "?":
		elem.card = zeroOrOnce
	case "1":
		elem.card = exactlyOnce
	case "+":
		elem.card = oneOrMore
	default:
		return fmt.Errorf("unknown cardinality '%s'", pvalTok.val)
	}
	return nil
}

func parseDefParam(lex *lexer, elem *tplElement, pvalTok *token) error {
	switch elem.typ {
	case Int:
		i, err := strconv.ParseInt(pvalTok.val, 10, 64)
		if err != nil {
			return err
		}
		return setDefDataRaw(elem, &i)
	case Uint:
		i, err := strconv.ParseUint(pvalTok.val, 10, 64)
		if err != nil {
			return err
		}
		return setDefDataRaw(elem, &i)
	case Float:
		f, err := strconv.ParseFloat(pvalTok.val, 64)
		if err != nil {
			return err
		}
		return setDefDataRaw(elem, &f)
	case String, Binary:
		if pvalTok.val[:2] == "0x" {
			s, err := hex.DecodeString(pvalTok.val[2:])
			if err != nil {
				return err
			}
			elem.def = []byte(s)
			return nil
		}
		if pvalTok.typ != quotedString {
			return fmt.Errorf("quoted string expected")
		}
		s, err := strconv.Unquote(pvalTok.val)
		if err != nil {
			return err
		}
		elem.def = []byte(s)
		return nil
	case Date:
		return fmt.Errorf("Default date not supported yet")
	default:
		return fmt.Errorf("Found default on unsupported type")
	}
}

func setDefDataRaw(elem *tplElement, d interface{}) error {
	buf := bytes.NewBuffer(make([]byte, 0, 8))
	if err := binary.Write(buf, binary.BigEndian, d); err != nil {
		return err
	}
	elem.def = buf.Bytes()
	return nil
}
