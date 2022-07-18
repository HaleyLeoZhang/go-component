package xmetric

import "github.com/prometheus/client_golang/prometheus"

const _metricNamespaceForHttp = "http"

// 本文件功能: 记录http请求情况

var (
	// 示例 promSQL  以15秒为时间间隔，采集数据，计算每个段内的命中率
	MetricHttpResponse = NewHistogramVec(&prometheus.HistogramOpts{
		Namespace: _metricNamespaceForHttp,
		Subsystem: "",
		Name:      "response_ms",
		Help:      "http requests duration",
		Buckets:   []float64{10, 25, 50, 100, 150, 200, 300, 500, 1000, 3000},
	}, []string{"method", "path"})

	MetricHttpRequestCount = NewCounterVec(&prometheus.CounterOpts{
		Namespace: _metricNamespaceForCache,
		Subsystem: "",
		Name:      "request_total",
		Help:      "http requests  total.",
	}, []string{"method", "path"})
)

// 使用示例
//xmetric.MetricHttpResponse.WithLabelValues(method, url).Observe(23.11))

// 使用示例
//xmetric.MetricHttpRequestCount.WithLabelValues(method, url).Inc()
