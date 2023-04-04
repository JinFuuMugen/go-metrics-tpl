package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func main() {
	GaugeMap := make(map[string]float64)
	CounterMap := make(map[string]int64)

	CounterMap["PollCounter"] = 0
	for {
		<-time.After(2 * time.Second)
		NewMonitor(&GaugeMap)
		CounterMap["PollCounter"]++
		<-time.After(8 * time.Second)
		for k, v := range GaugeMap {
			hc := &http.Client{}
			url := "http://localhost:8080/update/gauge/" + k + "/" + strconv.FormatFloat(v, 'E', -1, 64)
			req, _ := http.NewRequest(http.MethodPost, url, nil)
			resp, _ := hc.Do(req)
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err.Error())
				break
			}
			fmt.Println(string(body))
		}
		for k, v := range CounterMap {
			hc := &http.Client{}
			url := "http://localhost:8080/update/counter/" + k + "/" + strconv.FormatInt(v, 10)
			req, _ := http.NewRequest(http.MethodPost, url, nil)
			resp, _ := hc.Do(req)
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err.Error())
				break
			}
			fmt.Println(string(body))
		}
	}
}
