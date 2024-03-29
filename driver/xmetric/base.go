package xmetric

import "github.com/prometheus/client_golang/prometheus"

// 实例化Counter 并 注册指标
func NewCounterVec(cfg *prometheus.CounterOpts, labels []string) (metricItem *prometheus.CounterVec) {
	if cfg == nil {
		return
	}
	metricItem = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: cfg.Namespace, // 指标命名空间
		Subsystem: cfg.Subsystem, // 子系统名
		Name:      cfg.Name,      // 指标名称
		Help:      cfg.Help,      // 指标简介
	}, labels)
	prometheus.MustRegister(metricItem)
	return
}

// 实例化Histogram 并 注册指标
func NewHistogramVec(cfg *prometheus.HistogramOpts, labels []string) (metricItem *prometheus.HistogramVec) {
	if cfg == nil {
		return
	}
	metricItem = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: cfg.Namespace, // 指标命名空间
		Subsystem: cfg.Subsystem, // 子系统名
		Name:      cfg.Name,      // 指标名称
		Help:      cfg.Help,      // 指标简介
		Buckets:   cfg.Buckets,
	}, labels)
	prometheus.MustRegister(metricItem)
	return
}
