package monitors

import "github.com/JinFuuMugen/go-metrics-tpl.git/internal/sender"

type Monitor interface {
	Collect()
	Dump() error
	SetProcessor(p sender.Sender)
}

type RuntimeMonitor interface {
	Monitor
	CollectRuntimeMetrics()
}

type GopsutilMonitor interface {
	Monitor
	CollectGopsutil()
}
