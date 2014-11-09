package edtd

import (
	"math"
	"strings"
	. "testing"
)

func testRangeTok(t *T, typ Type, rangeStr string, expected *rangeParam) {
	rparts := strings.Split(rangeStr, ",")
	toks := make([]*token, len(rparts))
	for i := range rparts {
		toks[i] = &token{alphaNum, rparts[i]}
	}

	out, err := parseRangeParams(typ, toks)
	if err != nil {
		t.Fatalf("'%s' %s", rangeStr, err)
	}

	outCur, expectedCur := out, expected
	for i := 0; ; i++ {
		outMore, expectedMore := outCur.more, expectedCur.more
		outCur.more, expectedCur.more = nil, nil

		if *outCur != *expectedCur {
			t.Fatalf("'%s' gave %#v for '%s'", rangeStr, outCur, rparts[i])
		}
		if outMore == nil && expectedMore == nil {
			return
		}
		if outMore == nil || expectedMore == nil {
			t.Fatalf("'%s' gave %#v for '%s'", rangeStr, outCur, rparts[i])
		}
		outCur, expectedCur = outMore, expectedMore
	}
}

func TestIntRanges(t *T) {
	rtyp := Int
	rstr := "0..1,..1,-5..,0"
	rtok :=
		&rangeParam{
			loweri: 0,
			upperi: 1,
			more:
		&rangeParam{
			loweri: int64(math.MinInt64),
			upperi: 1,
			more:
		&rangeParam{
			loweri: -5,
			upperi: int64(math.MaxInt64),
			more:
		&rangeParam{
			loweri: 0,
			upperi: 0,
		}}}}
	testRangeTok(t, rtyp, rstr, rtok)
}

func TestUintRanges(t *T) {
	rtyp := Uint
	rstr := "0..1,2..,0"
	rtok :=
		&rangeParam{
			lowerui: 0,
			upperui: 1,
			more:
		&rangeParam{
			lowerui: 2,
			upperui: uint64(math.MaxUint64),
			more:
		&rangeParam{
			lowerui: 0,
			upperui: 0,
		}}}
	testRangeTok(t, rtyp, rstr, rtok)
}

func TestFloatRanges(t *T) {
	rtyp := Float
	rstr := ">=1.0,<2.0,-6.5<..<7.2,-6.5<=..<7.2,-6.5<..<=7.2,-6.5<=..<=7.2"
	rtok :=
		&rangeParam{
			lowerf: 1.0,
			upperf: math.MaxFloat64,
			exLower: false,
			exUpper: false,
			more:
		&rangeParam{
			lowerf: -1 * math.MaxFloat64,
			upperf: 2.0,
			exLower: true,
			exUpper: true,
			more:
		&rangeParam{
			lowerf: -6.5,
			upperf: 7.2,
			exLower: true,
			exUpper: true,
			more:
		&rangeParam{
			lowerf: -6.5,
			upperf: 7.2,
			exLower: false,
			exUpper: true,
			more:
		&rangeParam{
			lowerf: -6.5,
			upperf: 7.2,
			exLower: true,
			exUpper: false,
			more:
		&rangeParam{
			lowerf: -6.5,
			upperf: 7.2,
			exLower: false,
			exUpper: false,
		}}}}}}
	testRangeTok(t, rtyp, rstr, rtok)
}
