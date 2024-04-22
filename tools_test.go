package toolkit

import "testing"

func TestTools_RandomString(t *testing.T) {
	tools := Tools{}
	s := tools.RandomString(10)
	if len(s) != 10 {
		t.Errorf("Expected length of 10, got %d", len(s))
	}
}
