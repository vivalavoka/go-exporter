CREATE TABLE IF NOT EXISTS metrics (
  id VARCHAR PRIMARY KEY,
  m_type VARCHAR,
  value DOUBLE PRECISION DEFAULT 0,
  delta INT DEFAULT 0
);
