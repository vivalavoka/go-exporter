package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/vivalavoka/go-exporter/cmd/server/config"
)

type SHA256 struct {
	Enable bool
	key    string
}

func New(cfg config.Config) *SHA256 {
	isEmpty := cfg.SHAKey == ""

	return &SHA256{
		Enable: !isEmpty,
		key:    cfg.SHAKey,
	}
}

func (s *SHA256) GetSum(str string) string {
	hash := hmac.New(sha256.New, []byte(s.key))
	io.WriteString(hash, str)
	return fmt.Sprintf("%x", hash.Sum(nil))
}
