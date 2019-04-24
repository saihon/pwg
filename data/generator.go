package data

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
)

// RelativeFrequency appearance frequency to character
//
// map[int][][][]int{
//    (key - char code): {
//        { (0 - hi score), (1 - hi score), (2 - hi score), (3 - hi score) }
//    } {
//      { (0 - score), (1 - score), (2 - score), (3 - score), (4 - char code) },
//      { (0 - score), (1 - score), (2 - score), (3 - score), (4 - char code) },
//      { (0 - score), (1 - score), (2 - score), (3 - score), (4 - char code) },
//      ...
//    }
// }
//
// score:
//  0: Frequency of the second character from the beginning
//  1: Frequency of occurrence when the two previous character is a vowel
//  2: Frequency of occurrence when the two previous character is a consonant
//  3: Occurrence frequency of last character
//
// char code:
//  range 97 to 122
//
type RelativeFrequency map[int][][][]int

// IsValidData validation RelativeFrequency type data
func IsValidData(m RelativeFrequency) error {
	for k, v := range m {
		if k < 97 || k > 122 {
			return fmt.Errorf("%d is must be in the range of 97 to 122", k)
		}
		score := v[0][0]

		if len(score) < 4 {
			return fmt.Errorf("%d-0-0: slice must be 4 length or more %v", k, score)
		}

		for i, vvv := range v[1] {
			if len(vvv) < 5 {
				return fmt.Errorf("%d-1-%d: length must be 2 or more", k, i)
			}
			for j, s := range score {
				if vvv[j] > s {
					return fmt.Errorf("%d-1-%d: should not be higher than the high score: %d > %d", k, i, vvv[j], s)
				}
			}
			charcode := vvv[4]
			if charcode < 97 || charcode > 122 {
				return fmt.Errorf("%d-1-%d: %d is must be in the range of 97 to 122", charcode, k, i)
			}
		}
	}

	return nil
}

type intslice [][]int

func (a intslice) Len() int {
	return len(a)
}

func (a intslice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a intslice) Less(i, j int) bool {
	return a[i][4] > a[j][4]
}

func isvowel(n byte) bool {
	switch n {
	case 97, 101, 105, 111, 117:
		return true
	}
	return false
}

// Write
func Write(w io.Writer, m *RelativeFrequency) {
	key := make([]int, len(*m))
	i := 0
	for k := range *m {
		key[i] = k
		i++
	}
	sort.Ints(key)

	fmt.Fprint(w, "{\n")
	for _, k := range key {
		fmt.Fprintf(w, "\t%d: {\n\t\t{\n", k)
		v := (map[int][][][]int(*m))[k]
		for _, vv := range v[0] {
			format := "\t\t\t{%d, %d, %d, %d},\n"
			fmt.Fprintf(w, format, vv[0], vv[1], vv[2], vv[3])
		}
		v1 := intslice(v[1])
		sort.Sort(v1)
		for i, vv := range v1 {
			var format string
			switch i {
			case 0:
				format = "\t\t},\n\t\t{\n\t\t\t{%d, %d, %d, %d, %d},\n"
			case len(v[1]) - 1:
				format = "\t\t\t{%d, %d, %d, %d, %d},\n\t\t},\n\t},\n"
			default:
				format = "\t\t\t{%d, %d, %d, %d, %d},\n"
			}
			fmt.Fprintf(w, format, vv[0], vv[1], vv[2], vv[3], vv[4])
		}
	}
	fmt.Fprint(w, "}\n")
}

// Generate
func Generate(a []string) (RelativeFrequency, error) {
	mm := map[byte]map[byte][]int{}
	for _, v := range a {
		v = strings.ToLower(v)
		for i := 1; i < len(v); i++ {
			if v[i-1] < 97 || v[i-1] > 122 ||
				v[i] < 97 || v[i] > 122 {
				return nil, errors.New("out of range 97 to 122")
			}

			k1, k2 := byte(v[i-1]), byte(v[i])
			_, ok := mm[k1]
			if !ok {
				mm[k1] = map[byte][]int{
					k2: []int{0, 0, 0, 0},
				}
			} else if _, ok = mm[k1][k2]; !ok {
				mm[k1][k2] = []int{0, 0, 0, 0}
			}

			if i == 1 {
				mm[k1][k2][0]++

			} else if len(v)-1 == i {
				mm[k1][k2][3]++

			} else if i-2 >= 0 {
				if isvowel(v[i-2]) {
					mm[k1][k2][1]++
				} else {
					mm[k1][k2][2]++
				}
			}
		}
	}

	m := map[int][][][]int{}
	for k, v := range mm {
		k1 := int(k)
		_, ok := m[k1]
		if !ok {
			m[k1] = [][][]int{
				{{0, 0, 0, 0}}, {},
			}
		}

		for kk, vv := range v {
			k2 := int(kk)
			if vv[0] > m[k1][0][0][0] {
				m[k1][0][0][0] = vv[0]
			}
			if vv[1] > m[k1][0][0][1] {
				m[k1][0][0][1] = vv[1]
			}
			if vv[2] > m[k1][0][0][2] {
				m[k1][0][0][2] = vv[2]
			}
			if vv[3] > m[k1][0][0][3] {
				m[k1][0][0][3] = vv[3]
			}

			m[k1][1] = append(m[k1][1], append(vv, k2))
		}
	}

	return m, nil
}
