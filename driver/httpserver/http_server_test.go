package httpserver

import (
	"fmt"
	"github.com/HaleyLeoZhang/go-component/driver/xgin"
	"github.com/HaleyLeoZhang/go-component/driver/xmetric"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	c := &xgin.Config{}
	c.Debug = true
	c.Name = "testHttp"
	c.Timeout = 3 * time.Second
	ginEngine := xgin.New(c)
	ginEngine.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "pong")
	})

	c2 := &Config{}
	c2.Metrics = true
	c2.Pprof = true
	c2.Name = "testHttp"
	c2.Ip = "0.0.0.0"
	c2.Port = 80

	// 启动前注册指标
	if c2.Metrics {
		metricsTest()
	}

	Run(c2, ginEngine)
}

func metricsTest() {
	// 启动前注册指标
	// Case 1 --- 分桶计数
	var httpMsMetrics = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "__http_buckets",
		Subsystem: "",
		Name:      "buckets",
		Help:      "响应时间，毫秒",
		Buckets:   []float64{10, 25, 50, 100, 150, 200, 300, 500, 1000, 3000},
	}, []string{"service", "action"})
	// - 注册指标
	prometheus.MustRegister(httpMsMetrics)
	//reg := prometheus.NewRegistry()
	//reg.MustRegister(httpMs)
	// - bucket 打点
	for i := 0; i < 1000; i++ {
		var f = float64(i + 50)
		go func() {
			<-time.After(2 * time.Second)
			httpMsMetrics.WithLabelValues("api").Observe(f)
		}()
	}
	fmt.Println("Case 1 --- 分桶计数  Start")
	// Case 2 增长情况
	// 业务指标
	for i := 0; i < 324; i++ {
		go func() {
			<-time.After(2 * time.Second)
			xmetric.MetricProducer.WithLabelValues("blog_search").Inc() // 使用现成指标
		}()
	}
	fmt.Println("Case 2 --- 增长情况  Start")
}
