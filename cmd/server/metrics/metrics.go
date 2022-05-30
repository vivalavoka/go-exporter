package metrics

type Gauge float64
type Counter int64

const GaugeType = "gauge"
const CounterType = "counter"

type GagueMetric struct {
	Name  string
	Value Gauge
}

type CounterMetric struct {
	Name  string
	Value Counter
}
