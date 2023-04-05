package main

import (
	"fmt"
	"github.com/JinFuuMugen/go-metrics-tpl.git/cmd/agent/monitors"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type metricType interface {
	float64 | int64
}

func sendPost[T metricType](metricKind string, metricName string, metricValue T, client *http.Client) (*http.Response, error) {
	if client == nil {
		client = &http.Client{}
	}
	var url string
	if metricKind == "gauge" {
		url = "http://localhost:8080/update/" + metricKind + "/" + metricName + "/" + strconv.FormatFloat(float64(metricValue), 'E', -1, 64)
	} else {
		url = "http://localhost:8080/update/" + metricKind + "/" + metricName + "/" + strconv.FormatInt(int64(metricValue), 10)
	}
	req, _ := http.NewRequest(http.MethodPost, url, nil)
	req.Header.Add("Content-Type", "text/plain")
	resp, err := client.Do(req)
	if err == nil {
		io.Copy(os.Stdout, resp.Body)
	} else {
		fmt.Println("Connection error")
	}
	return resp, err
}

func main() {
	pollInterval := 2 * time.Second
	reportInterval := 10 * time.Second

	GaugeMap := make(map[string]float64)
	CounterMap := make(map[string]int64)

	CounterMap["PollCounter"] = 1
	ticks := 0

	client := &http.Client{
		Timeout: time.Second * 1,
	}
	for {
		<-time.After(pollInterval)
		monitors.NewMonitor(&GaugeMap)
		ticks++
		if ticks == int(reportInterval/pollInterval) {
			ticks = 0
			for k, v := range GaugeMap {
				resp, _ := sendPost("gauge", k, v, client)
				if resp != nil {
					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()
				}
			}
			for k, v := range CounterMap {
				resp, _ := sendPost("counter", k, v, client)
				if resp != nil {
					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()
				}
			}
		}
	}
}
