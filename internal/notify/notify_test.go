package notify

import "testing"

func TestFormatMessage(t *testing.T) {
	cases := []struct {
		want string
		n    int
	}{
		{"Copied 0 characters", 0},
		{"Copied 1 character", 1},
		{"Copied 2 characters", 2},
		{"Copied 142 characters", 142},
	}
	for _, c := range cases {
		if got := FormatMessage(c.n); got != c.want {
			t.Errorf("FormatMessage(%d) = %q, want %q", c.n, got, c.want)
		}
	}
}
