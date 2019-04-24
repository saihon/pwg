package data

import (
	"bufio"
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	data := []string{
		"foo",
		"bar",
		"baz",
	}

	actual, err := Generate(data)
	if err != nil {
		t.Errorf("\nshould be not an error: %v\n", err)
	}

	expect := RelativeFrequency{
		97:  {{{0, 0, 0, 1}}, {{0, 0, 0, 1, 114}, {0, 0, 0, 1, 122}}},
		98:  {{{2, 0, 0, 0}}, {{2, 0, 0, 0, 97}}},
		102: {{{1, 0, 0, 0}}, {{1, 0, 0, 0, 111}}},
		111: {{{0, 0, 0, 1}}, {{0, 0, 0, 1, 111}}},
	}

	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("\ngot : %v\nwant: %v\n", actual, expect)
	}
}

func TestWrite(t *testing.T) {
	data := RelativeFrequency{
		97:  {{{0, 0, 0, 1}}, {{0, 0, 0, 1, 114}, {0, 0, 0, 1, 122}}},
		98:  {{{2, 0, 0, 0}}, {{2, 0, 0, 0, 97}}},
		102: {{{1, 0, 0, 0}}, {{1, 0, 0, 0, 111}}},
		111: {{{0, 0, 0, 1}}, {{0, 0, 0, 1, 111}}},
	}

	expect := `{
	97: {
		{
			{0, 0, 0, 1},
		},
		{
			{0, 0, 0, 1, 122},
			{0, 0, 0, 1, 114},
		},
	},
	98: {
		{
			{2, 0, 0, 0},
		},
		{
			{2, 0, 0, 0, 97},
	102: {
		{
			{1, 0, 0, 0},
		},
		{
			{1, 0, 0, 0, 111},
	111: {
		{
			{0, 0, 0, 1},
		},
		{
			{0, 0, 0, 1, 111},
}
`

	var buff bytes.Buffer
	Write(&buff, &data)

	actual := buff.String()
	if actual != expect {

		e := bufio.NewScanner(strings.NewReader(expect))
		a := bufio.NewScanner(strings.NewReader(actual))

		for i := 0; e.Scan() && a.Scan(); i++ {
			if strings.Compare(e.Text(), a.Text()) != 0 {
				t.Errorf("\nline: %d:\n%v\n%v\n", i, e.Text(), a.Text())
				return
			}
		}
		if err := a.Err(); err != nil {
			t.Errorf("\nerror actual-scanner: %v\n", err)
		}
		if err := e.Err(); err != nil {
			t.Errorf("\nerror expect-scanner: %v\n", err)
		}
	}
}
