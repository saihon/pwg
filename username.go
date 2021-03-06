package pwg

import (
	"crypto/rand"

	"github.com/saihon/pwg/data"
)

var (
	Vowel     = []byte{'a', 'e', 'i', 'o', 'u'}
	Consonant = []byte{'b', 'c', 'd', 'f', 'g', 'h', 'j', 'k', 'l', 'm', 'n', 'p', 'q', 'r', 's', 't', 'v', 'w', 'x', 'y', 'z'}
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
		p.max.SetInt64(10 - 5)
		r, _ := rand.Int(rand.Reader, p.max)
		length = int(r.Int64() + 5)
	}

	if length < 1 {
		return []byte{}
	}

	a := make([]byte, length, length)
	p.max.SetInt64(123 - 97)
	r, _ := rand.Int(rand.Reader, p.max)
	a[0] = byte(int(r.Int64() + 97))

	for i := 1; i < length; i++ {
		var seed []byte
		key := int(a[i-1])

		var ok bool
		var mm [][][]int
		if m != nil {
			mm, ok = m[key]
		}

		if ok && len(mm) > 1 && len(mm[0]) > 0 && len(mm[0][0]) == 4 {
			// most frequent index
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
			max := mm[0][0][mf]

			if max > 0 {
				p.max.SetInt64(int64(max))
				r, _ := rand.Int(rand.Reader, p.max)
				n := int(r.Int64())
				for _, vv := range mm[1] {
					if vv[mf] == 0 || vv[mf] < n {
						continue
					}
					seed = append(seed, byte(vv[4]))
				}
			}
		}

		if len(seed) == 0 {
			if isvowel(key) {
				seed = Consonant
			} else {
				seed = Vowel
			}
		}

		n := 0
	Again:
		p.max.SetInt64(int64(len(seed)))
		r, _ := rand.Int(rand.Reader, p.max)
		c := byte(seed[int(r.Int64())])

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
