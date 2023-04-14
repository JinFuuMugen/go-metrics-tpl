package main

import (
	"flag"
	"fmt"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/monitors"
	"github.com/caarlos0/env"
	"github.com/go-resty/resty/v2"
	"strconv"
	"time"
)

type metricType interface {
	float64 | int64
}

type Config struct {
	Addr           string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func sendPost[T metricType](serverAddr string, metricKind string, metricName string, metricValue T, client *resty.Client) (*resty.Response, error) {

	var url string
	if metricKind == "gauge" {
		url = "http://" + serverAddr + "/update/" + metricKind + "/" + metricName + "/" + strconv.FormatFloat(float64(metricValue), 'E', -1, 64)
	} else {
		url = "http://" + serverAddr + "/update/" + metricKind + "/" + metricName + "/" + strconv.FormatInt(int64(metricValue), 10)
	}

	resp, err := client.R().SetHeader("Content-Type", "text/plain").Post(url)
	return resp, err
}

func main() {
	var cfg Config
	envParseError := env.Parse(&cfg)
	if envParseError != nil {
		panic(envParseError)
	}
	var serverAddr *string
	var poll *int
	var report *int

	if cfg.Addr != "" {
		serverAddr = &cfg.Addr
	} else {
		serverAddr = flag.String("a", "localhost:8080", "server address")
	}

	if cfg.PollInterval != 0 {
		poll = &cfg.PollInterval
	} else {
		poll = flag.Int("p", 2, "poll interval")
	}

	if cfg.ReportInterval != 0 {
		report = &cfg.ReportInterval
	} else {
		report = flag.Int("r", 10, "report interval")
	}
	flag.Parse()

	pollInterval := time.Duration(*poll) * time.Second
	pollTimer := time.NewTimer(pollInterval)

	reportInterval := time.Duration(*report) * time.Second
	reportTimer := time.NewTimer(reportInterval)

	GaugeMap := make(map[string]float64)
	CounterMap := make(map[string]int64)
	CounterMap["PollCounter"] = 1

	client := resty.New()

	for {
		select {
		case <-pollTimer.C:
			monitors.NewMonitor(&GaugeMap)
			pollTimer.Reset(pollInterval)
		case <-reportTimer.C:
			for k, v := range GaugeMap {
				resp, _ := sendPost(*serverAddr, "gauge", k, v, client)
				if resp != nil {
					fmt.Println(resp.StatusCode())
				}
			}
			for k, v := range CounterMap {
				resp, _ := sendPost(*serverAddr, "counter", k, v, client)
				if resp != nil {
					fmt.Println(resp.StatusCode())
				}
			}
			reportTimer.Reset(reportInterval)
		}
	}
}
