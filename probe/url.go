package probe

import (
	"github.com/shanghai-edu/n9e-probe/config"
	"github.com/shanghai-edu/n9e-probe/probe/core"
	"github.com/shanghai-edu/n9e-probe/probe/models"
	"github.com/toolkits/pkg/logger"
)

func UrlMetrics() []*core.MetricValue {
	ret := []*core.MetricValue{}

	for nid, urls := range config.Get().Url {
		results := models.UrlProbe(urls, config.Get().Probe.Headers, config.Get().Probe.Limit, config.Get().Probe.Timeout)
		for _, res := range results {
			logger.Debugf("url result %+v", res)
			tags := map[string]string{
				"scheme": res.Url.Scheme,
				"host":   res.Url.Host,
			}
			if res.Url.Path == "" {
				tags["path"] = "/"
			} else {
				tags["path"] = res.Url.Path
			}
			if config.Get().Probe.Region != "" {
				tags["region"] = config.Get().Probe.Region
			}
			ret = append(ret, core.GenProbeMetric(nid, "url.latency", res.Latency, tags))
			ret = append(ret, core.GenProbeMetric(nid, "url.cert", res.Cert, tags))
			ret = append(ret, core.GenProbeMetric(nid, "url.status_code", res.HttpStatusCode, tags))
			if res.Cert == 1 {
				ret = append(ret, core.GenProbeMetric(nid, "url.cert_remaining_day", res.CertRemainingDay, tags))
			}
		}
	}
	return ret
}
