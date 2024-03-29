package id

import (
	"crypto/rand"
	"fmt"
)

const idSize = 16

func New() string {
	b := make([]byte, idSize)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", b)
}
