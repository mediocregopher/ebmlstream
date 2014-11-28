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
	Name string
	Kids []*Elem
}

func (e *Edtd) NewParser(r io.Reader) *Parser {
	return &Parser{
		edtd:     e,
		lastElem: ebml.RootElem(r),
		buffer:   list.New(),
	}
}

func (p *Parser) NextShallow() (*Elem, error) {
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

	// If the element isn't a Container we call Bytes to force the ebml parser
	// to read in the body of the element and store it in the Elem's buffer
	if etpl.typ != Container {
		if _, err := e.Bytes(); err != nil {
			return nil, err
		}
	}

	p.lastElem = e

	return &Elem{
		Elem: *e,
		Type: etpl.typ,
		Name: etpl.name,
	}, nil
}

//func (p *Parser) Next() (*Elem, error) {
//	el, err := p.NextShallow()
//	if err != nil {
//		return nil, err
//	}
//
//	if el.Type != Container {
//		return el, nil
//	}
//
//	parentTPL := p.edtd.elements[elementID(el.Elem.Id)]
//	el.Kids = make([]*Elem, 0, len(parentTPL.kids))
//	if len(parentTPL.kids) == 0 {
//		return el, nil
//	}
//
//	for {
//		kid, err := p.lastElem.Next()
//		if err != nil {
//			return nil, err
//		}
//
//		etpl, ok := p.edtd.elements[elementID(kid.Id)]
//		if !ok {
//			return nil, fmt.Errorf("unknown id: %x", kid.Id)
//		}
//
//		
//	}
//}
