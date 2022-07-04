package postgresdb

import (
	"context"
	"fmt"

	pgx "github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
	"github.com/vivalavoka/go-exporter/cmd/server/config"
	"github.com/vivalavoka/go-exporter/internal/metrics"
)

type PostgresDB struct {
	config     config.Config
	connection *pgx.Conn
}

func New(cfg config.Config) (*PostgresDB, error) {
	conn, err := pgx.Connect(context.Background(), cfg.DatabaseDSN)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	return &PostgresDB{config: cfg, connection: conn}, nil
}

func (r *PostgresDB) Close() {
	r.connection.Close(context.Background())
}

func (r *PostgresDB) CheckConnection() bool {
	err := r.connection.Ping(context.Background())
	return err != nil
}

func (r *PostgresDB) GetMetrics() (map[string]metrics.Metric, error) {
	var metricMap = map[string]metrics.Metric{}

	rows, err := r.connection.Query(context.Background(), "SELECT id, m_type, value, delta FROM metrics;")
	if err != nil {
		return metricMap, fmt.Errorf("query row failed: %v", err)
	}

	for rows.Next() {
		var mID string
		var mType string
		delta := metrics.Counter(0)
		value := metrics.Gauge(0)

		err := rows.Scan(&mID, &mType, &value, &delta)
		if err != nil {
			log.Error(err)
		}
		metric := metrics.Metric{
			ID:    mID,
			MType: mType,
			Delta: delta,
			Value: value,
		}

		metricMap[metric.ID] = metric
	}

	return metricMap, rows.Err()
}

func (r *PostgresDB) GetMetric(ID string) (metrics.Metric, error) {
	var mID string
	var mType string
	var delta metrics.Counter
	var value metrics.Gauge

	err := r.connection.QueryRow(context.Background(), `
		SELECT id, m_type, delta, value FROM metrics WHERE id = $1;`,
		ID,
	).Scan(&mID, &mType, &value, &delta)
	if err != nil {
		return metrics.Metric{}, err
	}

	metric := metrics.Metric{
		ID:    mID,
		MType: mType,
		Delta: delta,
		Value: value,
	}

	return metric, nil
}

func (r *PostgresDB) Save(metric *metrics.Metric) error {
	_, err := r.connection.Exec(context.Background(), `
		INSERT INTO metrics(id, m_type, value, delta)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET value = $3, delta = $4;`,
		metric.ID, metric.MType, metric.Value, metric.Delta,
	)
	return err
}
