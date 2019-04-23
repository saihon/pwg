package pwg

import (
	"math/rand"

	"github.com/saihon/pwg/data"
)

func isvowel(n int) bool {
	switch n {
	case 97, 101, 105, 111, 117:
		return true
	}
	return false
}

// Username
func (p *password) Username(m data.RelativeFrequency, length int) []byte {
	if p.options.Random {
		length = rand.Intn(10-5) + 5
	}

	if length < 1 {
		return []byte{}
	}

	a := make([]byte, length, length)
	a[0] = byte(rand.Intn(123-97) + 97)

	for i := 1; i < length; i++ {
		key := int(a[i-1])

		// most frequent
		mf := 0
		switch {
		case i == length-1:
			mf = 3
		case i > 1:
			if isvowel(int(a[i-2])) {
				mf = 1
			} else {
				mf = 2
			}
		}
		max := m[key][0][0][mf]

		var seed []byte
		if max > 0 {
			// baseline
			r := rand.Intn(max)

			// append characters greater than or equal
			// the baseline to the slice as a seed
			// v[0] char code
			// v[1] score
			v := m[key]
			for _, vv := range v[1] {
				if vv[mf] == 0 || vv[mf] < r {
					continue
				}
				seed = append(seed, byte(vv[4]))
			}
		}

		if len(seed) == 0 {
			seed = []byte("abcdefghijklmnopqrstuvwxyz")
		}

		n := 0
	Again:
		c := byte(seed[rand.Intn(len(seed))])

		// if the same letter continues three times, try again
		if len(seed) > 1 && i-2 >= 0 {
			if a[i-1] == c && a[i-2] == c {
				if n < 3 {
					n++
					goto Again
				}
			}
		}

		a[i] = c
	}

	if p.options.Capitalize {
		a[0] = byte(int(a[0]) - 32)
	}
	return a
}

func (p *password) users(m data.RelativeFrequency, ch chan<- []byte) {
	length := p.options.Length
	if length <= 1 {
		length = 5
	}

	for i := 0; i < p.options.Generate; i++ {
		ch <- p.Username(m, length)
	}
	close(ch)
}

// Users
func (p *password) Users(m data.RelativeFrequency) <-chan []byte {
	ch := make(chan []byte, p.options.Generate)
	go p.users(m, ch)
	return ch
}
