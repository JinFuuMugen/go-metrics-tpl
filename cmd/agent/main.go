package main

import (
	"flag"
	"fmt"
	"github.com/JinFuuMugen/go-metrics-tpl.git/cmd/agent/monitors"
	"github.com/go-resty/resty/v2"
	"strconv"
	"time"
)

type metricType interface {
	float64 | int64
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
	serverAddr := flag.String("a", "localhost:8080", "server address")
	poll := flag.Int("p", 2, "poll interval")
	report := flag.Int("r", 10, "report interval")
	flag.Parse()

	pollInterval := time.Duration(*poll) * time.Second
	reportInterval := time.Duration(*report) * time.Second

	GaugeMap := make(map[string]float64)
	CounterMap := make(map[string]int64)

	CounterMap["PollCounter"] = 1
	ticks := 0

	client := resty.New()

	for {
		<-time.After(pollInterval)
		monitors.NewMonitor(&GaugeMap)
		ticks++
		if ticks == int(reportInterval/pollInterval) {
			ticks = 0
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
		}
	}
}
