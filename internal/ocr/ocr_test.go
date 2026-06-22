package ocr

import "testing"

func TestClean(t *testing.T) {
	cases := []struct {
		raw  string
		want string
	}{
		{"hello", "hello"},
		{"  hello  ", "hello"},
		{"\n\thello world\n", "hello world"},
		{"   \n  ", ""},
		{"", ""},
	}
	for _, c := range cases {
		if got := Clean(c.raw); got != c.want {
			t.Errorf("Clean(%q) = %q, want %q", c.raw, got, c.want)
		}
	}
}
