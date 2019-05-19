package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

type Hash struct{}

func NewHash() *Hash {
	return &Hash{}
}

func (h *Hash) Generate(s string, c chan string, e chan error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		e <- err
	} else {
		c <- string(hashedBytes[:])
	}
	close(c)
	close(e)
}

func (h *Hash) Compare(hash string, s string) error {
	incoming := []byte(s)
	existing := []byte(hash)
	return bcrypt.CompareHashAndPassword(existing, incoming)
}