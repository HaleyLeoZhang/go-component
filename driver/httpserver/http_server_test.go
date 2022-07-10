package httpserver

import (
	"fmt"
	"github.com/HaleyLeoZhang/go-component/driver/xgin"
	"github.com/HaleyLeoZhang/go-component/driver/xmetric"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"
)

type TestConfig struct {
	HttpServer *Config      `yaml:"httpServer"`
	Gin        *xgin.Config `yaml:"gin"`
}

var (
	cfg = &TestConfig{}
)

func TestRun(t *testing.T) {
	var yamlFile string
	yamlFile, err := filepath.Abs("./app.yaml") // 示例的kafka配置文件请看这个文件
	if err != nil {
		panic(err)
	}
	yamlRead, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlRead, cfg)
	if err != nil {
		panic(err)
	}
	// --
	ginEngine := xgin.New(cfg.Gin)
	ginEngine.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "pong")
	})
	// 启动前注册指标
	if cfg.HttpServer.Metrics {
		metricsTest()
	}

	Run(cfg.HttpServer, ginEngine)
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
	}, []string{"router"})
	// - 注册指标
	prometheus.MustRegister(httpMsMetrics)
	//reg := prometheus.NewRegistry()
	//reg.MustRegister(httpMs)
	// - bucket 打点
	for i := 0; i < 1000; i++ {
		var f = float64(i + 50)
		go func() {
			<-time.After(2 * time.Second)
			httpMsMetrics.WithLabelValues("comic/detail_by_id").Observe(f)
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
