package probe

import (
	"time"

	"github.com/freedomkk-qfeng/n9e-probe/probe/core"
)

func Collect() {

	for _, v := range Mappers {
		for _, f := range v.Fs {
			go collect(int64(v.Interval), f)
		}
	}
}

func collect(sec int64, fn func() []*core.MetricValue) {
	t := time.NewTicker(time.Second * time.Duration(sec))
	defer t.Stop()

	for {
		<-t.C

		metricValues := []*core.MetricValue{}
		now := time.Now().Unix()

		items := fn()
		if items == nil || len(items) == 0 {
			continue
		}

		for _, item := range items {
			item.Step = sec
			item.Timestamp = now
			metricValues = append(metricValues, item)
		}
		core.Push(metricValues)
	}
}
