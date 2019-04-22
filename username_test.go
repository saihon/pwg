package pwg

import (
	"testing"

	rf "github.com/saihon/pwg/data"
)

func TestUsername(t *testing.T) {
	p := &password{
		options: new(Options),
	}

	data := []struct {
		length     int
		capitalize bool
	}{
		{1, true},
		{2, true},
		{3, true},
		{4, true},
		{5, true},
		{6, true},
		{7, true},
		{8, true},
		{9, true},
		{10, true},
		{1, false},
		{2, false},
		{3, false},
		{4, false},
		{5, false},
		{6, false},
		{7, false},
		{8, false},
		{9, false},
		{10, false},
	}

	for i, v := range data {
		p.options.Capitalize = v.capitalize
		a := p.Username(rf.Data, v.length)
		if len(a) != v.length {
			t.Errorf("%d:\ngot : %d, want: %d\n", i, len(a), v.length)
			break
		}
		if v.capitalize {
			if a[0] < 65 || a[0] > 90 {
				t.Errorf("%d:\ngot : %c, want: %c\n", i, a[0], a[0]-32)
			}
		} else if a[0] < 97 || a[0] > 122 {
			t.Errorf("%d:\ngot : %c\n", i, a[0])
		}
	}
}
