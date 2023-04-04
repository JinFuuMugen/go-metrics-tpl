package main

import (
	"net/http"
)

var MS MemStorage

func main() {
	MS.Init()
	mux := http.NewServeMux()

	mux.HandleFunc(`/update/`, UpdateMetricsHandle)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
