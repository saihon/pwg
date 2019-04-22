package pwg

import (
	"time"

	mrand "math/rand"
)

var instance *password

type password struct {
	options *Options
	letters []string
}

// New new struct passowrd. rand.Seed runs only once
func New(o *Options) *password {
	if instance == nil {
		mrand.Seed(time.Now().UnixNano())
		instance = new(password)
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
				size = mrand.Intn((p.options.Length-adds-i)-1) + 1
			}
			adds += size
		} else {
			size = p.options.Length - adds
		}

		for j := 0; j < size; j++ {
			b[index] = p.letters[i][mrand.Intn(len(p.letters[i]))]
			index++
		}
	}
	return b
}

// Shuffle wrapped rand.Shuffle
func Shuffle(b []byte) {
	mrand.Shuffle(len(b), func(i, j int) {
		b[i], b[j] = b[j], b[i]
	})
}

func (p *password) gen(ch chan<- []byte) {
	for i := 0; i < p.options.Generate; i++ {
		b := p.Generate()
		Shuffle(b)
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
