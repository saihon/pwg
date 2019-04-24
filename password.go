package pwg

import (
	"crypto/rand"
	"math/big"
)

var instance *password

type password struct {
	options *Options
	letters []string
	max     *big.Int
}

// New new struct passowrd. rand.Seed runs only once
func New(o *Options) *password {
	if instance == nil {
		instance = new(password)
		instance.max = new(big.Int)
	}
	instance.options = o
	return instance.Init()
}

// Init initialize struct passowrd
func (p *password) Init() *password {
	if p.options.All {
		p.letters = []string{
			"0123456789",
			"abcdefghijklmnopqrstuvwxyz",
			"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			"!\"#$%&'()-=^~\\|@`[{;+:*]},<.>/?_",
		}
	} else {
		p.letters = nil

		if p.options.Number {
			p.letters = append(p.letters, "0123456789")
		}
		if p.options.LowerCase {
			p.letters = append(p.letters, "abcdefghijklmnopqrstuvwxyz")
		}
		if p.options.UpperCase {
			p.letters = append(p.letters, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
		}
		if p.options.Symbol {
			p.letters = append(p.letters, "!\"#$%&'()-=^~\\|@`[{;+:*]},<.>/?_")
		}
	}

	if p.options.Custom != "" {
		p.letters = append(p.letters, p.options.Custom)
	}

	if len(p.letters) == 0 {
		p.letters = append(p.letters, "0123456789")
	}

	if p.options.Length < 1 || p.options.Length < len(p.letters) {
		p.options.Length = len(p.letters)
	}

	if p.options.Generate < 1 {
		p.options.Generate = 1
	}

	return p
}

// Generate Generate one password.
// return value should be shuffled.
// b := p.Generate()
// pwg.Shuffle(b)
func (p *password) Generate() []byte {
	adds, index, size := 0, 0, 0
	remain := p.options.Length % len(p.letters)
	b := make([]byte, p.options.Length, p.options.Length)

	for i := len(p.letters) - 1; i >= 0; i-- {
		if i > 0 {
			if p.options.Evenly {
				size = p.options.Length / len(p.letters)
				if remain > 0 && i < remain {
					size++
				}
			} else {
				p.max.SetInt64(int64((p.options.Length - adds - i) - 1))
				r, _ := rand.Int(rand.Reader, p.max)
				size = int(r.Int64() + 1)
			}
			adds += size
		} else {
			size = p.options.Length - adds
		}

		p.max.SetInt64(int64(len(p.letters[i])))
		for j := 0; j < size; j++ {
			r, _ := rand.Int(rand.Reader, p.max)
			b[index] = p.letters[i][r.Int64()]
			index++
		}
	}
	return b
}

// Shuffle wrapped rand.Shuffle
func (p *password) Shuffle(b []byte) {
	p.max.SetInt64(int64(len(b)))
	for i := 0; i < len(b); i++ {
		r, _ := rand.Int(rand.Reader, p.max)
		j := r.Int64()
		b[i], b[j] = b[j], b[i]
	}
}

func (p *password) gen(ch chan<- []byte) {
	for i := 0; i < p.options.Generate; i++ {
		b := p.Generate()
		p.Shuffle(b)
		ch <- b
	}
	close(ch)
}

// Gen return iterable channel
func (p *password) Gen() <-chan []byte {
	ch := make(chan []byte, p.options.Generate)
	go p.gen(ch)
	return ch
}
