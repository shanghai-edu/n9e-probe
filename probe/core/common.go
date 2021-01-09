package core

func GenProbeMetric(nid, metric string, val interface{}, tags map[string]string) *MetricValue {
	mv := MetricValue{
		Nid:          nid,
		Metric:       metric,
		ValueUntyped: val,
		CounterType:  GAUGE,
		TagsMap:      map[string]string{},
	}

	for k, v := range tags {
		mv.TagsMap[k] = v
	}

	return &mv
}
