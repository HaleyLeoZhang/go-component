package httpserver

import (
	"github.com/HaleyLeoZhang/go-component/driver/xgin"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"math"
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
	// 尝试写入
	ginEngine.GET("/write", func(ctx *gin.Context) {
		addMetrics()
		ctx.String(200, "pong")
	})

	c2 := &Config{}
	c2.Metrics = true
	c2.Pprof = true
	c2.Name = "testHttp"
	c2.Ip = "0.0.0.0"
	c2.Port = 80
	Run(c2, ginEngine)
}

func addMetrics() {
	const NAMESPACE = "test"
	var httpMs = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: NAMESPACE,
		Name:      "http_response_ms",
		Help:      "响应时间，毫秒",
		Buckets:   []float64{50, 100, 200, 300, 500, 1000, 3000},
	}, []string{"service", "action"})

	for i := 0; i < 1000; i++ {
		httpMs.WithLabelValues("comic_service", "api").Observe(50 + math.Floor(120*math.Sin(float64(i)*0.1))/10)
	}

	reg := prometheus.NewRegistry()
	reg.MustRegister(httpMs)
}
