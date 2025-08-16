package httpserver

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Config struct {
	Name  string `yaml:"name" json:"name"` // 用于 Trace 识别
	Ip    string `yaml:"ip" json:"ip"`
	Port  int    `yaml:"port" json:"port"`
	Pprof bool   `yaml:"pprof" json:"pprof"` // true 开启  pprof 性能监控路由
	// 注: 网关层请不要让外部访问到 /metrics 这个路由
	Metrics bool `yaml:"metrics" json:"metrics"` // true 开启  metrics 打点，支持 prometheus 主动来拉数据
	// -
	ReadTimeout    time.Duration `yaml:"readTimeout" json:"readTimeout"`
	WriteTimeout   time.Duration `yaml:"writeTimeout" json:"writeTimeout"`
	MaxHeaderBytes int           `yaml:"maxHeaderBytes" json:"maxHeaderBytes"`
}

func Run(c *Config, routersInit *gin.Engine) {
	addrString := fmt.Sprintf("%s:%v", c.Ip, c.Port)

	if c.Pprof {
		// pprof 相关说明 http://www.hlzblog.top/article/74.html
		fmt.Println("Enabled pprof")
		Wrap(routersInit)
	}
	if c.Metrics {
		// prometheus 相关说明 https://prometheus.io/docs/guides/go-application/
		fmt.Println("Enabled metrics")
		WrapPrometheus(routersInit)
	}

	server := &http.Server{
		Addr:           addrString,
		Handler:        routersInit,
		ReadTimeout:    c.ReadTimeout,
		WriteTimeout:   c.WriteTimeout,
		MaxHeaderBytes: c.MaxHeaderBytes,
	}
	fmt.Println("Start http server listening ", addrString)
	err := server.ListenAndServe()
	if err != nil {
		msg := fmt.Sprintf("HttpServer.Err %+v", err)
		fmt.Println(msg)
	}
}
