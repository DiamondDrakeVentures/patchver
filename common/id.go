package common

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// ID generates an ID string with the default length.
func ID() string {
	id, err := GenID(8)
	if err != nil {
		panic(fmt.Errorf("cannot generate ID: %v", err))
	}

	return id
}

// GenID generates an ID string with length len.
func GenID(len int) (string, error) {
	if len < 2 {
		len = 2
	} else {
		len /= 2
	}

	return genID(len)
}

func genID(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
