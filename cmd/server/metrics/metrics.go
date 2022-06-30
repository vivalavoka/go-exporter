package metrics

import "fmt"

type Gauge float64
type Counter int64

const GaugeType = "gauge"
const CounterType = "counter"

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *Counter `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *Gauge   `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`
}

func (m Metric) String() string {
	switch m.MType {
	case GaugeType:
		return fmt.Sprintf("%s:%s:%f", m.ID, m.MType, *m.Value)
	case CounterType:
		return fmt.Sprintf("%s:%s:%d", m.ID, m.MType, *m.Delta)
	}
	return ""
}
