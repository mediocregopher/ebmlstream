package edtd

import (
	"bytes"
	. "testing"
)

// I'm being lazy and just doing a single giant fucking string. yolo
func TestLexer(t *T) {
	testStr := `
		define types {
			bool := uint [ range:0..1; ]
			ascii := string [ range:32..126; ]
		}
		define elements {
			
			// Snip... farther down the list
			TimecodeScale := 2ad7b1 uint [ def:1000000; ]
			Duration := 4489 float [ range:>0.0; ]

			// Snip
			Language := 22b59c string [ def:"eng"; range:32..126; ]
	`

	output := []Token{
		{AlphaNum, "define"},
		{AlphaNum, "types"},
		{Control,  "{"},

		{AlphaNum, "bool"},
		{Control, ":="},
		{AlphaNum, "uint"},
		{Control, "["},
		{AlphaNum, "range"},
		{Control, ":"},
		{AlphaNum, "0..1"},
		{Control, ";"},
		{Control, "]"},

		{AlphaNum, "ascii"},
		{Control, ":="},
		{AlphaNum, "string"},
		{Control, "["},
		{AlphaNum, "range"},
		{Control, ":"},
		{AlphaNum, "32..126"},
		{Control, ";"},
		{Control, "]"},
		
		{Control, "}"},

		{AlphaNum, "define"},
		{AlphaNum, "elements"},
		{Control, "{"},

		{AlphaNum, "TimecodeScale"},
		{Control, ":="},
		{AlphaNum, "2ad7b1"},
		{AlphaNum, "uint"},
		{Control, "["},
		{AlphaNum, "def"},
		{Control, ":"},
		{AlphaNum, "1000000"},
		{Control, ";"},
		{Control, "]"},

		{AlphaNum, "Duration"},
		{Control, ":="},
		{AlphaNum, "4489"},
		{AlphaNum, "float"},
		{Control, "["},
		{AlphaNum, "range"},
		{Control, ":"},
		{Control, ">"},
		{AlphaNum, "0.0"},
		{Control, ";"},
		{Control, "]"},

		{AlphaNum, "Language"},
		{Control, ":="},
		{AlphaNum, "22b59c"},
		{AlphaNum, "string"},
		{Control, "["},
		{AlphaNum, "def"},
		{Control, ":"},
		{QuotedString, "\"eng\""},
		{Control, ";"},
		{AlphaNum, "range"},
		{Control, ":"},
		{AlphaNum, "32..126"},
		{Control, ";"},
		{Control, "]"},
		{EOF, ""},
	}

	buf := bytes.NewBufferString(testStr)
	l := NewLexer(buf)

	for i := range output {
		tok := l.Next()
		t.Logf("Checking for %#v", output[i])
		if *tok != output[i] {
			t.Fatalf("Found %#v instead of %#v", *tok, output[i])
		}
	}
}
