package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mediocregopher/go.ebml/edtd"
)

func main() {
	edtdf, err := os.Open("matroska.edtd")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("parsing matroska.edtd")
	e, err := edtd.NewEdtd(edtdf)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open("test.webm")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("starting parswer for test.webm")
	p := e.NewParser(f)
	for {
		el, err := p.NextShallow()
		if err != nil {
			log.Fatal(err)
		}

		prefix := fmt.Sprintf("%x %d %s", el.Id, el.Elem.Size, el.Name)
		var line string
		var thing interface{}
		switch el.Type {
		case edtd.Int:
			thing, err = el.Int()
			line = fmt.Sprintf("%s - %d", prefix, thing)
		case edtd.Uint:
			thing, err = el.Uint()
			line = fmt.Sprintf("%s - %d", prefix, thing)
		case edtd.Float:
			thing, err = el.Float()
			line = fmt.Sprintf("%s - %f", prefix, thing)
		case edtd.String:
			thing, err = el.String()
			line = fmt.Sprintf("%s - %s", prefix, thing)
		default:
			line = prefix
		}

		if err != nil {
			log.Fatalf("line: %s, err: %s", line, err)
		}

		log.Println(line)
	}
}
