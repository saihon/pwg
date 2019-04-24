package pwg

import (
	"math/big"
	"math/rand"
	"reflect"
	"regexp"
	"testing"
	"time"
)

func TestLetters(t *testing.T) {
	data := []struct {
		options *Options
		expect  []string
	}{
		{options: &Options{
			All: true,
		}, expect: []string{
			"0123456789",
			"abcdefghijklmnopqrstuvwxyz",
			"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			"!\"#$%&'()-=^~\\|@`[{;+:*]},<.>/?_",
		}},
		{options: &Options{
			Number:    true,
			LowerCase: true,
			UpperCase: true,
			Symbol:    true,
		}, expect: []string{
			"0123456789",
			"abcdefghijklmnopqrstuvwxyz",
			"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			"!\"#$%&'()-=^~\\|@`[{;+:*]},<.>/?_",
		}},
		{&Options{}, []string{"0123456789"}}, // if all of false, only number letters
		{&Options{Number: true}, []string{"0123456789"}},
		{&Options{LowerCase: true}, []string{"abcdefghijklmnopqrstuvwxyz"}},
		{&Options{UpperCase: true}, []string{"ABCDEFGHIJKLMNOPQRSTUVWXYZ"}},
		{&Options{Symbol: true}, []string{"!\"#$%&'()-=^~\\|@`[{;+:*]},<.>/?_"}},
		{&Options{Custom: "hello"}, []string{"hello"}},
		{options: &Options{
			Number:    true,
			LowerCase: true,
		}, expect: []string{
			"0123456789",
			"abcdefghijklmnopqrstuvwxyz",
		}},
		{options: &Options{
			LowerCase: true,
			UpperCase: true,
		}, expect: []string{
			"abcdefghijklmnopqrstuvwxyz",
			"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		}},
		{options: &Options{
			UpperCase: true,
			Symbol:    true,
		}, expect: []string{
			"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			"!\"#$%&'()-=^~\\|@`[{;+:*]},<.>/?_",
		}},
	}

	for i, v := range data {
		p := (&password{
			options: v.options,
		}).Init()

		if !reflect.DeepEqual(p.letters, v.expect) {
			t.Errorf("%d\n  got : %v\n  want: %v\n", i, p.letters, v.expect)
			break
		}
	}
}

func TestInit(t *testing.T) {
	data := []struct {
		options  *Options
		generate int
		length   int
	}{
		{&Options{Generate: 0, Length: 0}, 1, 1},
		{&Options{All: true, Generate: 0, Length: 0}, 1, 4},
		{&Options{Number: true, UpperCase: true, Generate: 0, Length: 0}, 1, 2},
	}

	var expect, actual int
	for i, v := range data {
		p := (&password{options: v.options}).Init()

		actual = p.options.Generate
		expect = v.generate
		if actual != expect {
			t.Errorf("%d: options.Generate:\ngot : %d, want: %d\n", i, actual, expect)
			break
		}

		actual = p.options.Length
		expect = v.length
		if actual != expect {
			t.Errorf("%d: options.Digits:\ngot : %d, want: %d\n", i, actual, expect)
		}
	}
}

func TestGenerate(t *testing.T) {
	data := []struct {
		options *Options
		re      *regexp.Regexp
	}{
		{options: &Options{
			Number: true,
			Length: 30,
		},
			re: regexp.MustCompile(`^\d{30}$`),
		},
		{options: &Options{
			Number:    true,
			LowerCase: true,
			Length:    50,
		},
			re: regexp.MustCompile(`^([a-z]|\d){50}$`),
		},
		{options: &Options{
			LowerCase: true,
			UpperCase: true,
			Length:    40,
		},
			re: regexp.MustCompile(`^[a-zA-Z]{40}$`),
		},
		{options: &Options{
			UpperCase: true,
			Symbol:    true,
			Length:    50,
		},
			re: regexp.MustCompile("^[A-Z!\"#$%&'()\\-=^~|@`[{;+:*\\]},<.>/?_\\\\]{50}$"),
		},
	}

	rand.Seed(time.Now().UnixNano())

	for i, v := range data {
		p := (&password{
			options: v.options,
			max:     new(big.Int),
		}).Init()
		b := p.Generate()

		if !v.re.Match(b) {
			t.Errorf("%d:\nnot matched: %s\n", i, string(b))
			break
		}
	}
}

type count struct {
	upper  int
	lower  int
	number int
	symbol int
}

func counter(b []byte) count {
	var c count
	for i := 0; i < len(b); i++ {
		switch {
		case 48 <= b[i] && b[i] <= 57:
			c.number++
		case 65 <= b[i] && b[i] <= 90:
			c.upper++
		case 97 <= b[i] && b[i] <= 122:
			c.lower++
		case (b[i] >= 33 && b[i] <= 47) ||
			(b[i] >= 58 && b[i] <= 64) ||
			(b[i] >= 91 && b[i] <= 96) ||
			(b[i] >= 123 && b[i] <= 126):
			c.symbol++
		}
	}
	return c
}

func TestEvenly(t *testing.T) {
	data := []struct {
		options *Options
		count   count
	}{
		{&Options{
			Length: 20,
		}, count{number: 20}},
		{&Options{
			Number: true,
			Length: 30,
		}, count{number: 30}},
		{&Options{
			UpperCase: true,
			Length:    20,
		}, count{upper: 20}},
		{&Options{
			LowerCase: true,
			Length:    40,
		}, count{lower: 40}},
		{&Options{
			Symbol: true,
			Length: 10,
		}, count{symbol: 10}},
		{&Options{
			All:    true,
			Length: 40,
		}, count{10, 10, 10, 10}},
		{&Options{
			Number:    true,
			LowerCase: true,
			Length:    40,
		}, count{number: 20, lower: 20}},
		{&Options{
			LowerCase: true,
			UpperCase: true,
			Length:    40,
		}, count{upper: 20, lower: 20}},
		{&Options{
			UpperCase: true,
			Symbol:    true,
			Length:    40,
		}, count{upper: 20, symbol: 20}},
		{&Options{
			Number:    true,
			UpperCase: true,
			Symbol:    true,
			Length:    30,
		}, count{number: 10, upper: 10, symbol: 10}},
	}

	for i, v := range data {
		v.options.Evenly = true
		p := (&password{
			options: v.options,
			max:     new(big.Int),
		}).Init()

		b := p.Generate()
		c := counter(b)
		if !reflect.DeepEqual(c, v.count) {
			t.Errorf("%d:\ngot : %#v\nwant: %#v\n", i, c, v.count)
		}
	}
}
