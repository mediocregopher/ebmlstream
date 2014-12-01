package ebmlstream

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"time"

	"github.com/mediocregopher/ebmlstream/varint"
)

type Elem struct {
	r   io.Reader
	buf *bufio.Reader
	data []byte

	Id   int64
	Size int64
}

func RootElem(r io.Reader) *Elem {
	return &Elem{
		r:   r,
		buf: bufio.NewReader(r),
	}
}

func (e *Elem) Next() (*Elem, error) {
	id, err := varint.ReadVarInt(e.buf)
	if err != nil {
		return nil, err
	}

	size, err := varint.ReadVarInt(e.buf)
	if err != nil {
		return nil, err
	}

	return &Elem{
		r:    e.r,
		buf:  e.buf,
		Id:   id,
		Size: size,
	}, nil
}

func (e *Elem) fillBuffer(total int64) error {
	if e.data == nil {
		if total == -1 {
			total = e.Size
		}
		e.data = make([]byte, total)
		n, err := io.ReadFull(e.buf, e.data[total-e.Size:])
		if int64(n) == e.Size {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

func (e *Elem) Int() (int64, error) {
	if e.Size == 0 {
		return 0, nil
	} else if err := e.fillBuffer(8); err != nil {
		return 0, err
	}

	var ret int64
	buf := bytes.NewBuffer(e.data)
	if err := binary.Read(buf, binary.BigEndian, &ret); err != nil {
		return 0, err
	}

	return ret, nil
}

func (e *Elem) Uint() (uint64, error) {
	if e.Size == 0 {
		return 0, nil
	}
	if err := e.fillBuffer(8); err != nil {
		return 0, err
	}

	var ret uint64
	buf := bytes.NewBuffer(e.data)
	if err := binary.Read(buf, binary.BigEndian, &ret); err != nil {
		return 0, err
	}

	return ret, nil
}

var timeStart = time.Date(
	2001, time.January, 1,
	0, 0, 0, 0,
	time.UTC,
)

func (e *Elem) Date() (time.Time, error) {
	i, err := e.Int()
	if err != nil {
		return time.Time{}, err
	}
	return timeStart.Add(time.Duration(i)), nil
}

func (e *Elem) Float() (float64, error) {
	if e.Size == 0 {
		return 0, nil
	} else if e.Size == 4 {
		f, err := e.f32()
		return float64(f), err
	} else if err := e.fillBuffer(8); err != nil {
		return 0, err
	}

	var ret float64
	buf := bytes.NewBuffer(e.data)
	if err := binary.Read(buf, binary.BigEndian, &ret); err != nil {
		return 0, err
	}

	return ret, nil
}

func (e *Elem) f32() (float32, error) {
	if err := e.fillBuffer(4); err != nil {
		return 0, err
	}

	var ret float32
	buf := bytes.NewBuffer(e.data)
	if err := binary.Read(buf, binary.BigEndian, &ret); err != nil {
		return 0, err
	}

	return ret, nil
}

func (e *Elem) String() (string, error) {
	if e.Size == 0 {
		return "", nil
	} else if err := e.fillBuffer(-1); err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(e.data)
	ret, err := buf.ReadString(0)
	if err != nil {
		return ret, nil
	} else {
		return ret[:len(ret)-1], nil
	}
}

func (e *Elem) Bytes() ([]byte, error) {
	if e.Size == 0 {
		return []byte{}, nil
	} else if err := e.fillBuffer(-1); err != nil {
		return nil, err
	}

	return e.data, nil
}
