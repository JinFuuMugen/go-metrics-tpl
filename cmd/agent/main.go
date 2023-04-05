package main

import (
	"github.com/JinFuuMugen/go-metrics-tpl.git/cmd/agent/monitors"
	"net/http"
	"strconv"
	"time"
)

type metricType interface {
	float64 | int64
}

func sendPost[T metricType](metricKind string, metricName string, metricValue T) (*http.Response, error) {
	hc := &http.Client{}
	var url string
	if metricKind == "gauge" {
		url = "http://localhost:8080/update/" + metricKind + "/" + metricName + "/" + strconv.FormatFloat(float64(metricValue), 'E', -1, 64)
	} else {
		url = "http://localhost:8080/update/" + metricKind + "/" + metricName + "/" + strconv.FormatInt(int64(metricValue), 10)
	}
	req, _ := http.NewRequest(http.MethodPost, url, nil)
	req.Header.Add("Content-Type", "text/plain")
	resp, err := hc.Do(req)
	return resp, err
}

func main() {
	pollInterval := 2 * time.Second
	reportInterval := 10 * time.Second
	GaugeMap := make(map[string]float64)
	CounterMap := make(map[string]int64)

	CounterMap["PollCounter"] = 1
	ticks := 0
	for {
		<-time.After(pollInterval)
		monitors.NewMonitor(&GaugeMap)
		//fmt.Println("REFRESHED", ticks)
		ticks++
		if ticks == int(reportInterval/pollInterval) {
			ticks = 0
			for k, v := range GaugeMap {
				_, err := sendPost("gauge", k, v)
				if err != nil {
					//fmt.Println("Connection error")
					continue
				}
				//bodyString, _ := io.ReadAll(resp.Body)
				//fmt.Println(string(bodyString))
			}
			for k, v := range CounterMap {
				_, err := sendPost("counter", k, v)
				if err != nil {
					//fmt.Println("Connection error")
					continue
				}
				//bodyString, _ := io.ReadAll(resp.Body)
				//fmt.Println(string(bodyString))
			}
		}
	}
}
