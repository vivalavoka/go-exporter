package postgresdb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/vivalavoka/go-exporter/cmd/server/config"
	"github.com/vivalavoka/go-exporter/internal/metrics"
)

type PostgresDB struct {
	config     config.Config
	connection *sqlx.DB
	insertStmt *sql.Stmt
}

func New(cfg config.Config) (*PostgresDB, error) {
	conn, err := sqlx.Open("postgres", cfg.DatabaseDSN)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	postgres := PostgresDB{config: cfg, connection: conn}

	err = postgres.migration()

	if err != nil {
		return nil, fmt.Errorf("migration failed: %v", err)
	}

	postgres.insertStmt, err = postgres.connection.Prepare(`
		INSERT INTO metrics (id, m_type, value, delta)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET value = $3, delta = metrics.delta + $4;`,
	)

	if err != nil {
		return nil, err
	}

	return &postgres, nil
}

func (r *PostgresDB) migration() error {
	rows, err := r.connection.Query(`
		CREATE TABLE IF NOT EXISTS metrics (
			id VARCHAR PRIMARY KEY,
			m_type VARCHAR,
			value DOUBLE PRECISION DEFAULT 0,
			delta INT DEFAULT 0
		);`,
	)
	if rows.Err() != nil {
		return rows.Err()
	}
	return err
}

func (r *PostgresDB) Close() {
	r.connection.Close()
}

func (r *PostgresDB) CheckConnection() bool {
	err := r.connection.Ping()
	return err == nil
}

func (r *PostgresDB) GetMetrics() (map[string]metrics.Metric, error) {
	var metricList []metrics.Metric
	var metricMap = map[string]metrics.Metric{}

	err := r.connection.Select(&metricList, "SELECT id, m_type, value, delta FROM metrics;")
	if err != nil {
		return metricMap, fmt.Errorf("query row failed: %v", err)
	}

	for _, metric := range metricList {
		metricMap[metric.ID] = metric
	}

	return metricMap, nil
}

func (r *PostgresDB) GetMetric(ID string, MType string) (metrics.Metric, error) {
	var metric metrics.Metric

	err := r.connection.Get(&metric, `SELECT id, m_type, delta, value FROM metrics WHERE id = $1 AND m_type = $2;`, ID, MType)

	if err != nil {
		return metrics.Metric{}, fmt.Errorf("there is no metric by name: %s", ID)
	}

	return metric, nil
}

func (r *PostgresDB) Save(metric *metrics.Metric) error {
	log.Info(metric)
	_, err := r.insertStmt.Exec(metric.ID, metric.MType, metric.Value, metric.Delta)
	return err
}

func (r *PostgresDB) SaveBatch(metricList []*metrics.Metric) error {
	tx, err := r.connection.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := tx.StmtContext(context.Background(), r.insertStmt)
	defer stmt.Close()

	for _, metric := range metricList {
		if _, err = stmt.Exec(metric.ID, metric.MType, metric.Value, metric.Delta); err != nil {
			if err = tx.Rollback(); err != nil {
				log.Fatalf("update drivers: unable to rollback: %v", err)
			}
			return err
		}
	}
	return tx.Commit()
}
