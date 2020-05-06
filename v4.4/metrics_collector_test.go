package exporter

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func TestSanitizePrometheusExporterName(t *testing.T) {
	sanitized := sanitizePrometheusExporterName("zabbix_internal_checks_zabbix[wcache,values,float]")
	assert.Equal(t, "zabbix_internal_checks_zabbix__wcache:values:float__", sanitized)
	promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "test",
		Name:      sanitized,
		Help:      "test help",
	})

	sanitized = sanitizePrometheusExporterName("zabbix[wcache,values,not supported]")
	assert.Equal(t, "zabbix__wcache:values:not_supported__", sanitized)
	promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "test",
		Name:      sanitized,
		Help:      "test help",
	})

	sanitized = sanitizePrometheusExporterName("zabbix[process,self-monitoring,avg,busy]")
	assert.Equal(t, "zabbix__process:self_monitoring:avg:busy__", sanitized)
	promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "test",
		Name:      sanitized,
		Help:      "test help",
	})

	sanitized = sanitizePrometheusExporterName("zabbix[stats,{$ADDRESS},{$PORT},queue]")
	assert.Equal(t, "zabbix__stats:__ADDRESS__:__PORT__:queue__", sanitized)
	promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "test",
		Name:      sanitized,
		Help:      "test help",
	})
}
