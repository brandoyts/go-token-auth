package hash

import "golang.org/x/crypto/bcrypt"

type Hash struct{}

func New() *Hash {
	return &Hash{}
}

func (h *Hash) Generate(value string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (h *Hash) Compare(hashed string, value string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(value))
}
