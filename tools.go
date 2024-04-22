package toolkit

import (
	"crypto/rand"
)

const randomStringSource = "abcdefghijklmnopqrstuwwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_+"

// Tools is a struct that contains utility functions
type Tools struct{}

// RandomString generates a random string of a given length using the randomStringSource
// constant as the source of characters
func (t *Tools) RandomString(length int) string {
	s, r := make([]rune, length), []rune(randomStringSource)
	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))
		s[i] = r[x%y]
	}
	return string(s)
}
