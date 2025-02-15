package services

import (
	"crypto/sha256"
	"fmt"
)

func hash(in string) string {
	hash := sha256.New()
	hash.Write([]byte(in))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
