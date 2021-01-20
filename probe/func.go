package probe

import (
	"github.com/shanghai-edu/n9e-probe/config"
	"github.com/shanghai-edu/n9e-probe/probe/core"
)

type FuncsAndInterval struct {
	Fs       []func() []*core.MetricValue
	Interval int
}

var Mappers []FuncsAndInterval

func BuildMappers() {
	interval := config.Get().Probe.Interval
	Mappers = []FuncsAndInterval{
		{
			Fs: []func() []*core.MetricValue{
				PingMetrics,
				UrlMetrics,
			},
			Interval: interval,
		},
	}

}
