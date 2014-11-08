// The lex package implements a lexical reader which can take in any io.Reader.
// It does not care about the meaning or logical validity of the tokens it
// parses out, it simply does its job.
package edtd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"unicode"
)

type TokenType int

const (
	AlphaNum TokenType = iota
	Control
	QuotedString
	Err
	EOF
)

// Token represents a single set of characters which *could* be a valid token of
// the given type
type Token struct {
	Type TokenType
	Val  string
}

// Returns the token's value as an error, or nil if the token is not of type
// Err. If the token is nil returns io.EOF, since that is the ostensible meaning
func (t *Token) AsError() error {
	if t.Type == EOF {
		return io.EOF
	} else if t.Type == Err {
		return errors.New(t.Val)
	}
	return nil
}

var (
	errInvalidUTF8 = errors.New("invalid utf8 character")
)

// Lexer reads through an io.Reader and emits Tokens from it.
type Lexer struct {
	r      *bufio.Reader
	outbuf *bytes.Buffer
	ch     chan *Token
}

// NewLexer constructs a new Lexer struct and returns it. r is internally
// wrapped with a bufio.Reader, unless it already is one. This will spawn a
// go-routine which reads from r until it hits an error, at which point it will
// end execution.
func NewLexer(r io.Reader) *Lexer {
	var br *bufio.Reader
	var ok bool
	if br, ok = r.(*bufio.Reader); !ok {
		br = bufio.NewReader(r)
	}

	l := Lexer{
		r:      br,
		ch:     make(chan *Token),
		outbuf: bytes.NewBuffer(make([]byte, 0, 1024)),
	}

	go l.spin()

	return &l
}

func (l *Lexer) spin() {
	f := lexWhitespace
	for {
		f = f(l)
		if f == nil {
			return
		}
	}
}

// Returns the next available token. This method should not be called after any
// Err or EOF tokens
func (l *Lexer) Next() *Token {
	return <-l.ch
}

func (l *Lexer) emit(t TokenType) {
	str := l.outbuf.String()
	l.ch <- &Token{
		Type: t,
		Val:  str,
	}
	l.outbuf.Reset()
}

func (l *Lexer) peek() (rune, error) {
	r, err := l.readRune()
	if err != nil {
		return 0, err
	}
	if err = l.r.UnreadRune(); err != nil {
		return 0, err
	}
	return r, nil
}

func (l *Lexer) readRune() (rune, error) {
	r, i, err := l.r.ReadRune()
	if err != nil {
		return 0, err
	} else if r == unicode.ReplacementChar && i == 1 {
		return 0, errInvalidUTF8
	}
	return r, nil
}

func (l *Lexer) err(err error) lexerFunc {
	if err == io.EOF {
		l.ch <- &Token{EOF, ""}
	} else {
		l.ch <- &Token{Err, err.Error()}
	}
	close(l.ch)
	return nil
}

func (l *Lexer) errf(format string, args ...interface{}) lexerFunc {
	s := fmt.Sprintf(format, args...)
	l.ch <- &Token{Err, s}
	close(l.ch)
	return nil
}

type lexerFunc func(*Lexer) lexerFunc

func lexWhitespace(l *Lexer) lexerFunc {
	r, err := l.readRune()
	if err != nil {
		return l.err(err)
	}

	if unicode.IsSpace(r) {
		return lexWhitespace
	} else if r == '/' {
		return lexComment
	}

	l.outbuf.WriteRune(r)

	if r == ':' {
		return lexColon
	} else if r == '"' {
		return lexQuotedString
	} else if unicode.IsLetter(r) && unicode.IsNumber(r) {
		return lexAlphaNum
	} else {
		l.emit(Control)
		return lexWhitespace
	}
}

func lexComment(l *Lexer) lexerFunc {
	r, err := l.peek()
	// There is a / in the buffer. If there's an error or not another / after, I
	// guess it's a control character?
	if err != nil {
		l.emit(Control)
		return l.err(err)
	} else if r != '/' {
		l.emit(Control)
		return lexWhitespace
	}

	l.readRune()
	return lexCommentRest
}

// Reads until a newline, since everything after a // is tossed away
func lexCommentRest(l *Lexer) lexerFunc {
	r, err := l.readRune()
	if err != nil {
		return l.err(err)
	}

	if r == '\n' {
		return lexWhitespace
	}
	return lexCommentRest
}

func lexAlphaNum(l *Lexer) lexerFunc {
	r, err := l.peek()
	if err != nil {
		return l.err(err)
	}

	// Underscore is allowed by the spec, . isn't but it's convenient to have
	// here because it's used in places like range and decimal numbers where it
	// acts as part of a single thing
	if r == '_' || r == '.' || unicode.IsLetter(r) || unicode.IsNumber(r) {
		l.readRune()
		l.outbuf.WriteRune(r)
		return lexAlphaNum
	}
	l.emit(AlphaNum)

	return lexWhitespace
}

func lexColon(l *Lexer) lexerFunc {
	r, err := l.peek()
	if err != nil {
		return l.err(err)
	}

	if r == '=' {
		l.readRune()
		l.emit(Control)
	}

	return lexWhitespace
}

func lexQuotedString(l *Lexer) lexerFunc {
	r, err := l.readRune()
	if err != nil {
		l.emit(QuotedString)
		return l.err(err)
	}
	l.outbuf.WriteRune(r)

	if r == '\\' {
		r, err := l.readRune()
		if err != nil {
			l.emit(QuotedString)
			return l.err(err)
		}
		l.outbuf.WriteRune(r)
	} else if r == '"' {
		return lexWhitespace
	}

	return lexQuotedString
}
