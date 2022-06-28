package crypto

import (
	"crypto/sha256"
	"fmt"

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
	h := sha256.New()
	h.Write([]byte(s.key))
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}
