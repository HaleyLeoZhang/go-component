package xmetric

import "github.com/prometheus/client_golang/prometheus"

const _metricNamespaceForKafka = "kafka"

// 本文件功能: kafka 不同Topic消息 生产、消费数量打点
// - 文末有使用示例

var (
	// 示例 promSQL   以1分钟为采样周期，预测整个采样周期内每秒的增长率，通过增长率汇总1分钟内可能发送的量
	// sum (rate (kafka_consumer{job="mlf-k8s-prd-pods"}[1m])) by (topic)

	MetricConsumer = NewCounterVec(&prometheus.CounterOpts{
		Namespace: _metricNamespaceForKafka,
		Subsystem: "",
		Name:      "consumer",
		Help:      "kafka consumer speed",
	}, []string{"topic", "partition"})

	// 示例 promSQL
	// sum (rate (kafka_producer{job="mlf-k8s-prd-pods"}[1m])) by (topic)

	MetricProducer = NewCounterVec(&prometheus.CounterOpts{
		Namespace: _metricNamespaceForKafka,
		Subsystem: "",
		Name:      "producer",
		Help:      "kafka producer speed.",
	}, []string{"topic"})
)

// 使用示例 有一条 Topic 为 blog_search 的消息被发送了
// xmetric.MetricProducer.WithLabelValues("blog_search").Inc() // 使用现成指标
