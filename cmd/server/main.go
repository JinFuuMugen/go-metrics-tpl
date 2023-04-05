package main

import (
	"github.com/JinFuuMugen/go-metrics-tpl.git/cmd/server/handlers"
	"net/http"
)

func main() {
	handlers.MS.Init()
	mux := http.NewServeMux()

	mux.HandleFunc(`/update/`, handlers.UpdateMetricsHandle)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
