package xmetric

import "github.com/prometheus/client_golang/prometheus"

const _metricNamespaceForCache = "cache"

// 本文件功能: 记录命中未命中情况

var (
	// 示例 promSQL  以15秒为时间间隔，采集数据，计算每个段内的命中率
	// avg(
	//	cache_hit_total{job="mlf-k8s-prd-pods"}[15s] / (
	//		cache_hit_total{job="mlf-k8s-prd-pods"}[15s] + cache_miss_total{job="mlf-k8s-prd-pods"}[15s]
	//	)
	//) by (name)

	MetricHit = NewCounterVec(&prometheus.CounterOpts{
		Namespace: _metricNamespaceForCache,
		Subsystem: "",
		Name:      "hit_total",
		Help:      "cache hit total.",
	}, []string{"name"})

	MetricMiss = NewCounterVec(&prometheus.CounterOpts{
		Namespace: _metricNamespaceForCache,
		Subsystem: "",
		Name:      "miss_total",
		Help:      "cache miss total.",
	}, []string{"name"})
)

// 使用示例 有一个 缓存key 为 comic_detail 的数据命中缓存了
// xmetric.MetricHit.WithLabelValues("comic_detail").Inc() // 使用现成指标
