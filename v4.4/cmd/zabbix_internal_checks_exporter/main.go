package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/moznion/zabbix_internal_checks_exporter/internal"
	exporter "github.com/moznion/zabbix_internal_checks_exporter/v4.4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const defaultPort = -1
const defaultZabbixURL = ""
const defaultZabbixUserName = ""
const defaultZabbixPassword = ""
const defaultIntervalSec = 30

func parseCommandLineOptions() (int64, string, string, string, time.Duration, error) {
	var port int64
	var zabbixURL, zabbixUserName, zabbixPassword string
	var intervalSec uint64
	var showVersion bool
	flag.Int64Var(&port, "port", defaultPort, "[mandatory] a port number of exporter listens")
	flag.StringVar(&zabbixURL, "zabbix-url", defaultZabbixURL, "[mandatory] a Zabbix server URL to collect the metrics")
	flag.StringVar(&zabbixUserName, "zabbix-user", defaultZabbixUserName, "[mandatory] a Zabbix server user name for authentication to use API")
	flag.StringVar(&zabbixPassword, "zabbix-password", defaultZabbixPassword, "[mandatory] a Zabbix server password for authentication to use API")
	flag.Uint64Var(&intervalSec, "interval-sec", defaultIntervalSec, "[optional] an interval seconds of collecting the metrics")
	flag.BoolVar(&showVersion, "version", false, "show version information")
	flag.Parse()

	if showVersion { // XXX dirty!
		fmt.Printf("%s\n", internal.GetVersions())
		os.Exit(0)
	}

	if port == defaultPort {
		return 0, "", "", "", 0, errors.New(`"--port" is mandatory parameter, but that is missing`)
	}
	if zabbixURL == defaultZabbixURL {
		return 0, "", "", "", 0, errors.New(`"--zabbix-url" is mandatory parameter, but that is missing`)
	}
	if zabbixUserName == defaultZabbixUserName {
		return 0, "", "", "", 0, errors.New(`"--zabbix-user" is mandatory parameter, but that is missing`)
	}
	if zabbixPassword == defaultZabbixPassword {
		return 0, "", "", "", 0, errors.New(`"--zabbix-password" is mandatory parameter, but that is missing`)
	}
	return port, zabbixURL, zabbixUserName, zabbixPassword, time.Duration(intervalSec) * time.Second, nil
}

func main() {
	port, zabbixURL, zabbixUserName, zabbixPassword, interval, err := parseCommandLineOptions()
	if err != nil {
		flag.Usage()
		log.Fatal(err)
	}

	exporter.NewMetricsCollector(
		&internal.JSONRPCClient{ZabbixBaseURL: zabbixURL},
		zabbixUserName,
		zabbixPassword,
		interval,
	).StartCollecting()

	http.Handle("/metrics", promhttp.Handler())

	log.Printf("[info] start listening :%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
