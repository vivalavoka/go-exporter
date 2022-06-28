package main

import (
	"crypto/sha256"
	"fmt"
)

func GetSHA256(str string) string {
	h := sha256.New()
	h.Write([]byte(config.SHAKey))
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}
