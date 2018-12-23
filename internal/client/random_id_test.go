package client

import "testing"

func TestRandomIDLen(t *testing.T) {
	tests := []struct {
		n int
	}{
		{1},
		{5},
		{100},
	}
	for _, tt := range tests {
		act := len(RandomID(tt.n))
		if act != tt.n {
			t.Errorf("false key length %d, expected %d", act, tt.n)
		}
	}
}
