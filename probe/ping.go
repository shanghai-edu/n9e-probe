package probe

import (
	"github.com/freedomkk-qfeng/n9e-probe/config"
	"github.com/freedomkk-qfeng/n9e-probe/probe/core"
	"github.com/freedomkk-qfeng/n9e-probe/probe/models"
	"github.com/toolkits/pkg/logger"
)

func PingMetrics() []*core.MetricValue {
	ret := []*core.MetricValue{}

	for nid, ips := range config.Get().Ping {
		results := models.PingProbe(ips, config.Get().Probe.Limit, config.Get().Probe.Timeout)
		for _, res := range results {
			logger.Debugf("ping result %+v", res)
			tags := map[string]string{
				"ip": res.Ip,
			}
			if config.Get().Probe.Region != "" {
				tags["region"] = config.Get().Probe.Region
			}
			ret = append(ret, core.GenProbeMetric(nid, "ping.latency", res.Ping, tags))

		}
	}
	return ret
}
