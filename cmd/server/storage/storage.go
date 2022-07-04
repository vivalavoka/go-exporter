package storage

import (
	"github.com/vivalavoka/go-exporter/cmd/server/config"
	memorydb "github.com/vivalavoka/go-exporter/cmd/server/storage/repositories/in-memory"
	"github.com/vivalavoka/go-exporter/internal/metrics"
)

type MetricsRepoInterface interface {
	Close()
	GetMetrics() (map[string]metrics.Metric, error)
	GetMetric(ID string) (metrics.Metric, error)
	Save(metric *metrics.Metric) error
}

type Storage struct {
	config config.Config
	Repo   MetricsRepoInterface
}

func New(config config.Config) (*Storage, error) {
	var repo MetricsRepoInterface
	var err error
	repo, err = memorydb.New(config)

	if err != nil {
		return nil, err
	}

	storage := &Storage{
		config: config,
		Repo:   repo,
	}

	return storage, nil
}

func (s *Storage) Close() {
	s.Repo.Close()
}
