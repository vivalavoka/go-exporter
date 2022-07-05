package storage

import (
	"github.com/vivalavoka/go-exporter/cmd/server/config"
	"github.com/vivalavoka/go-exporter/cmd/server/metrics"
	memorydb "github.com/vivalavoka/go-exporter/cmd/server/storage/repositories/in-memory"
	postgresdb "github.com/vivalavoka/go-exporter/cmd/server/storage/repositories/in-memory/postgres"
)

type MetricsRepoInterface interface {
	Close()
	CheckConnection() bool
	GetMetrics() (map[string]metrics.Metric, error)
	GetMetric(ID string) (metrics.Metric, error)
	Save(metric *metrics.Metric) error
}

type Storage struct {
	config config.Config
	Repo   MetricsRepoInterface
}

var stg *Storage

func New(config config.Config) (*Storage, error) {
	var repo MetricsRepoInterface
	var err error

	if config.DatabaseDSN == "" {
		repo, err = memorydb.New(config)
	} else {
		repo, err = postgresdb.New(config)
	}

	if err != nil {
		return nil, err
	}

	stg = &Storage{
		config: config,
		Repo:   repo,
	}

	return stg, nil
}

func GetStorage() *Storage {
	return stg
}

func (s *Storage) Close() {
	s.Repo.Close()
}
