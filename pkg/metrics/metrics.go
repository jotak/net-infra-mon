package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

const prefix = "netinfra_"

var log = logrus.WithField("module", "metrics")

type MetricDefinition struct {
	Name   string
	Help   string
	Type   metricType
	Labels []string
}

type metricType string

const (
	TypeCounter   metricType = "counter"
	TypeGauge     metricType = "gauge"
	TypeHistogram metricType = "histogram"
)

var allMetrics = []MetricDefinition{}

func Declare(name, help string, t metricType, labels ...string) MetricDefinition {
	def := MetricDefinition{
		Name:   name,
		Help:   help,
		Type:   t,
		Labels: labels,
	}
	allMetrics = append(allMetrics, def)
	return def
}

func (def *MetricDefinition) mapLabels(labels []string) prometheus.Labels {
	if len(labels) != len(def.Labels) {
		log.Errorf("Could not map labels, length differ in def %s [%v / %v]", def.Name, def.Labels, labels)
	}
	labelsMap := prometheus.Labels{}
	for i, label := range labels {
		labelsMap[def.Labels[i]] = label
	}
	return labelsMap
}

func verifyMetricType(def *MetricDefinition, t metricType) {
	if def.Type != t {
		log.Errorf("metric %q is of type %q but is being registered as %q", def.Name, def.Type, t)
	}
}

// register will register against the default registry
func register(c prometheus.Collector, name string) {
	err := prometheus.DefaultRegisterer.Register(c)
	if err != nil {
		log.Errorf("metrics registration error [%s]: %v", name, err)
	}
}

func NewCounter(def *MetricDefinition, labels ...string) prometheus.Counter {
	verifyMetricType(def, TypeCounter)
	fullName := prefix + def.Name
	c := prometheus.NewCounter(prometheus.CounterOpts{
		Name:        fullName,
		Help:        def.Help,
		ConstLabels: def.mapLabels(labels),
	})
	register(c, fullName)
	return c
}

func NewCounterVec(def *MetricDefinition) *prometheus.CounterVec {
	verifyMetricType(def, TypeCounter)
	fullName := prefix + def.Name
	c := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: fullName,
		Help: def.Help,
	}, def.Labels)
	register(c, fullName)
	return c
}

func NewGauge(def *MetricDefinition, labels ...string) prometheus.Gauge {
	verifyMetricType(def, TypeGauge)
	fullName := prefix + def.Name
	c := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        fullName,
		Help:        def.Help,
		ConstLabels: def.mapLabels(labels),
	})
	register(c, fullName)
	return c
}

func NewGaugeVec(def *MetricDefinition) *prometheus.GaugeVec {
	verifyMetricType(def, TypeGauge)
	fullName := prefix + def.Name
	g := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: fullName,
		Help: def.Help,
	}, def.Labels)
	register(g, fullName)
	return g
}

func NewHistogram(def *MetricDefinition, buckets []float64, labels ...string) prometheus.Histogram {
	verifyMetricType(def, TypeHistogram)
	fullName := prefix + def.Name
	c := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:        fullName,
		Help:        def.Help,
		Buckets:     buckets,
		ConstLabels: def.mapLabels(labels),
	})
	register(c, fullName)
	return c
}
