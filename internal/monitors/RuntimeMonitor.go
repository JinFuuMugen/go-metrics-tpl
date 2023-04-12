package monitors

import (
	"math/rand"
	"runtime"
)

func NewMonitor(GaugeMap *map[string]float64) {
	var rtm runtime.MemStats
	RandomValue := 1000 * rand.Float64()

	// Read full mem stats
	runtime.ReadMemStats(&rtm)

	(*GaugeMap)["BuckHashSys"] = float64(rtm.BuckHashSys)   //uint64
	(*GaugeMap)["Alloc"] = float64(rtm.Alloc)               //uint64
	(*GaugeMap)["Frees"] = float64(rtm.Frees)               //uint64
	(*GaugeMap)["GCCPUFraction"] = rtm.GCCPUFraction        //float64
	(*GaugeMap)["GCSys"] = float64(rtm.GCSys)               //uint64
	(*GaugeMap)["HeapAlloc"] = float64(rtm.HeapAlloc)       //uint64
	(*GaugeMap)["HeapAlloc"] = float64(rtm.HeapIdle)        //uint64
	(*GaugeMap)["HeapInuse"] = float64(rtm.HeapInuse)       //uint64
	(*GaugeMap)["HeapObjects"] = float64(rtm.HeapObjects)   //uint64
	(*GaugeMap)["HeapReleased"] = float64(rtm.HeapReleased) //uint64
	(*GaugeMap)["HeapSys"] = float64(rtm.HeapSys)           //uint64
	(*GaugeMap)["LastGC"] = float64(rtm.LastGC)             //uint64
	(*GaugeMap)["Lookups"] = float64(rtm.Lookups)           //uint64
	(*GaugeMap)["MCacheInuse"] = float64(rtm.MCacheInuse)   //uint64
	(*GaugeMap)["MCacheSys"] = float64(rtm.MCacheSys)       //uint64
	(*GaugeMap)["MSpanInuse"] = float64(rtm.MSpanInuse)     //uint64
	(*GaugeMap)["MSpanSys"] = float64(rtm.MSpanSys)         //uint64
	(*GaugeMap)["Mallocs"] = float64(rtm.Mallocs)           //uint64
	(*GaugeMap)["NextGC"] = float64(rtm.NextGC)             //uint64
	(*GaugeMap)["NumForcedGC"] = float64(rtm.NumForcedGC)   //uint32
	(*GaugeMap)["NumGC"] = float64(rtm.NumGC)               //uint32
	(*GaugeMap)["OtherSys"] = float64(rtm.OtherSys)         //uint64
	(*GaugeMap)["PauseTotalNs"] = float64(rtm.PauseTotalNs) //uint64
	(*GaugeMap)["StackInuse"] = float64(rtm.StackInuse)     //uint64
	(*GaugeMap)["StackSys"] = float64(rtm.StackSys)         //uint64
	(*GaugeMap)["Sys"] = float64(rtm.Sys)                   //uint64
	(*GaugeMap)["TotalAlloc"] = float64(rtm.TotalAlloc)     //uint64
	(*GaugeMap)["RandomValue"] = RandomValue                //float64
}
