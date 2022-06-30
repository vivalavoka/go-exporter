package crypto

import (
	"crypto/sha256"
	"fmt"

	log "github.com/sirupsen/logrus"
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
	log.Info(str)
	log.Info(s.key)
	h := sha256.New()
	h.Write([]byte(str))
	h.Write([]byte(s.key))
	hash := fmt.Sprintf("%x", h.Sum(nil))
	log.Info(hash)
	log.Info()
	return hash
}
