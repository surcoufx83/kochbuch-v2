package services

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

func hash(in string) string {
	return Sha256(in, 0)
}

func Sha256(input1 string, cutafter int) string {
	hash := sha256.New()
	hash.Write([]byte(input1))
	hashstr := fmt.Sprintf("%x", hash.Sum(nil))
	if cutafter <= 0 || cutafter > len(hashstr) {
		return hashstr
	}
	return hashstr[:cutafter]
}

func Sha512(input1 string, input2 string) string {
	hash := sha512.New()
	hash.Write([]byte(input1))
	if input2 == "" {
		return fmt.Sprintf("%x", hash.Sum(nil))
	}
	return fmt.Sprintf("%x", hash.Sum([]byte(input2)))
}
