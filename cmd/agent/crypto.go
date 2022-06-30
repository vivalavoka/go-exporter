package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
)

func GetSHA256(str string) string {
	hash := hmac.New(sha256.New, []byte(config.SHAKey))
	io.WriteString(hash, str)
	return fmt.Sprintf("%x", hash.Sum(nil))
}
