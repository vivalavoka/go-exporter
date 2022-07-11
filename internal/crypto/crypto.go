package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
)

type SHA256 struct {
	Enable bool
	key    string
}

func New(key string) *SHA256 {
	isEmpty := key == ""

	return &SHA256{
		Enable: !isEmpty,
		key:    key,
	}
}

func (s *SHA256) GetSum(str string) string {
	hash := hmac.New(sha256.New, []byte(s.key))
	io.WriteString(hash, str)
	return fmt.Sprintf("%x", hash.Sum(nil))
}
