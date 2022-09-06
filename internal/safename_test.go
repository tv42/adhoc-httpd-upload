package internal_test

import (
	"testing"

	"eagain.net/go/adhoc-httpd-upload/internal"
)

func TestSafename(t *testing.T) {
	for idx, test := range []struct {
		name string
		want bool
	}{
		{"good", true},
		{".", false},
		{".bad", false},
		{"", false},
		{"/", false},
		{"//", false},
		{"../bad", false},
		{"/bad", false},
	} {
		got := internal.IsSafeName(test.name)
		if g, e := got, test.want; g != e {
			t.Errorf("#%d %q: %v != %v\n", idx, test.name, g, e)
		}
	}
}
