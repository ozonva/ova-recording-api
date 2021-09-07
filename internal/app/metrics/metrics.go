package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics interface {
	IncSuccessCreateAppointmentCounter()
	IncFailCreateAppointmentCounter()
	IncSuccessMultiCreateAppointmentCounter()
	IncFailMultiCreateAppointmentCounter()
	IncSuccessUpdateAppointmentCounter()
	IncFailUpdateAppointmentCounter()
	IncSuccessRemoveAppointmentCounter()
	IncFailRemoveAppointmentCounter()
}

const (
	successResultLabel = "success"
	failResultLabel    = "fail"
)

var labels = []string{"result"}

type metrics struct {
	createAppointmentCounter *prometheus.CounterVec
	multiCreateAppointmentCounter *prometheus.CounterVec
	updateAppointmentCounter *prometheus.CounterVec
	removeAppointmentCounter *prometheus.CounterVec
}

func NewApiMetrics() Metrics {
	m := &metrics{
		createAppointmentCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "create_appointment_request_count",
			Help: "number of created appointments",
		},
		labels),
		multiCreateAppointmentCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "multi_create_appointment_request_count",
			Help: "number of multi created appointments",
		},
		labels),
		updateAppointmentCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "update_appointment_request_count",
			Help: "number of updated appointments",
		},
		labels),
		removeAppointmentCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "remove_appointment_request_count",
			Help: "number of removed appointments",
		},
		labels),
	}
	return m
}

func (m *metrics) IncSuccessCreateAppointmentCounter() {
	m.createAppointmentCounter.WithLabelValues(successResultLabel).Inc()
}

func (m *metrics) IncFailCreateAppointmentCounter() {
	m.createAppointmentCounter.WithLabelValues(failResultLabel).Inc()
}

func (m *metrics) IncSuccessMultiCreateAppointmentCounter() {
	m.multiCreateAppointmentCounter.WithLabelValues(successResultLabel).Inc()
}

func (m *metrics) IncFailMultiCreateAppointmentCounter() {
	m.multiCreateAppointmentCounter.WithLabelValues(failResultLabel).Inc()
}

func (m *metrics) IncSuccessUpdateAppointmentCounter() {
	m.updateAppointmentCounter.WithLabelValues(successResultLabel).Inc()
}

func (m *metrics) IncFailUpdateAppointmentCounter() {
	m.updateAppointmentCounter.WithLabelValues(failResultLabel).Inc()
}

func (m *metrics) IncSuccessRemoveAppointmentCounter() {
	m.removeAppointmentCounter.WithLabelValues(successResultLabel).Inc()
}

func (m *metrics) IncFailRemoveAppointmentCounter() {
	m.removeAppointmentCounter.WithLabelValues(failResultLabel).Inc()
}
