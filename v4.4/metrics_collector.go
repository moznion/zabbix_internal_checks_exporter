package exporter

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/moznion/zabbix_internal_checks_exporter/internal"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const enabledItemStatus = 0

const namespace = "zabbix_internal_checks"
const metricsCollectedAtEpochMillis = "metrics_collected_at"

// MetricsCollector is a struct that represents a collector for "Zabbix internal checks".
type MetricsCollector struct {
	metrics       map[string]prometheus.Gauge
	jsonRPCClient *internal.JSONRPCClient
	userName      string
	password      string
	interval      time.Duration
}

// NewMetricsCollector makes a MetricCollector according to the given parameters.
func NewMetricsCollector(jsonRPCClient *internal.JSONRPCClient, userName string, password string, interval time.Duration) *MetricsCollector {
	return &MetricsCollector{
		metrics:       make(map[string]prometheus.Gauge),
		jsonRPCClient: jsonRPCClient,
		userName:      userName,
		password:      password,
		interval:      interval,
	}
}

// StartCollecting starts a goroutine to collect the "Zabbix internal checks" metrics.
func (mc *MetricsCollector) StartCollecting() {
	go func() {
		mc.initMetrics()

		userLoginResponse, err := mc.jsonRPCClient.UserLogin(mc.userName, mc.password)
		if err != nil {
			log.Fatalf("failed user.login request on init: %s", err)
		}
		authToken := userLoginResponse.AuthToken

		mutChan := make(chan struct{}, 1)
		tick := time.Tick(mc.interval)

		for range tick {
			func() {
				select {
				case mutChan <- struct{}{}:
					// NOP
				default:
					// a lock has already been acquired, skip
					return
				}

				defer func() {
					<-mutChan
				}()

				itemResponse, err := mc.jsonRPCClient.GetItem(authToken, "zabbix[*]")
				if err != nil {
					if errors.Is(err, internal.ErrJSONRPCClientRequestError) {
						// if it comes here, the reason is either one of the followings:
						//   - not auth token was expired
						//   - a request was completely failed
						// So if it is in this clause, it attempts signing-in at first;
						// it that passed authentication it can proceed RPC process,
						// otherwise, it should be failed.
						userLoginResponse, err := mc.jsonRPCClient.UserLogin(mc.userName, mc.password)
						if err != nil {
							log.Printf("[error] failed user.login request on re-authentication: %s; continue", err)
						}
						authToken = userLoginResponse.AuthToken
						return
					}

					log.Fatalf("[error] failed get.item request: %s", err)
				}

				for _, metricItem := range itemResponse.Result {
					status, err := strconv.ParseInt(metricItem.Status, 10, 64)
					if err != nil {
						log.Printf("[warn] failed to parse a `status` of item: %s; skip", err)
						continue
					}
					if status != enabledItemStatus {
						continue
					}

					lastValueStr, ok := metricItem.LastValue.(string)
					if !ok {
						continue
					}
					lastValue, err := strconv.ParseFloat(lastValueStr, 64)
					if err != nil {
						continue
					}

					if mc.metrics[metricItem.Key] == nil {
						mc.metrics[metricItem.Key] = promauto.NewGauge(prometheus.GaugeOpts{
							Namespace: namespace,
							Name:      sanitizePrometheusExporterName(metricItem.Key),
							Help:      metricItem.Name,
						})
					}
					mc.metrics[metricItem.Key].Set(lastValue)
				}
				mc.setCollectedAt(time.Now().UnixNano() / int64(time.Millisecond))
			}()
		}
	}()
}

func (mc *MetricsCollector) setCollectedAt(epochMillis int64) {
	mc.metrics[metricsCollectedAtEpochMillis].Set(float64(epochMillis))
}

func (mc *MetricsCollector) initMetrics() {
	// put a timestamp on data collected
	mc.metrics[metricsCollectedAtEpochMillis] = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      metricsCollectedAtEpochMillis,
		Help:      "Epoch milliseconds at Zabbix internal checks metrics were collected",
	})
	mc.setCollectedAt(-1)
}

func sanitizePrometheusExporterName(name string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(
					strings.ReplaceAll(
						name, "[", "__",
					), ",", ":",
				), "]", "",
			), " ", "_",
		), "-", "_",
	)
}
