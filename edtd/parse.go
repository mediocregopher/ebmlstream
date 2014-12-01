package edtd

import (
	"container/list"
	"fmt"
	"io"

	ebml "github.com/mediocregopher/go.ebml"
)

type Parser struct {
	edtd     *Edtd
	lastElem *ebml.Elem
	buffer   *list.List
}

type Elem struct {
	ebml.Elem
	Type
	Name  string
	Level uint64
}

func (e *Edtd) NewParser(r io.Reader) *Parser {
	return &Parser{
		edtd:     e,
		lastElem: ebml.RootElem(r),
		buffer:   list.New(),
	}
}

func (p *Parser) Next() (*Elem, error) {
	if f := p.buffer.Front(); f != nil {
		return p.buffer.Remove(f).(*Elem), nil
	}

	e, err := p.lastElem.Next()
	if err != nil {
		return nil, err
	}

	etpl, ok := p.edtd.elements[elementID(e.Id)]
	if !ok {
		return nil, fmt.Errorf("unknown id: %x", e.Id)
	}

	switch etpl.typ {
	case Int:
		_, err = e.Int()
	case Uint:
		_, err = e.Uint()
	case Float:
		_, err = e.Float()
	case Date:
		_, err = e.Date()
	case String:
		_, err = e.String()
	case Binary:
		_, err = e.Bytes()
	}
	if err != nil {
		return nil, err
	}

	p.lastElem = e

	return &Elem{
		Elem:  *e,
		Type:  etpl.typ,
		Name:  etpl.name,
		Level: etpl.level,
	}, nil
}