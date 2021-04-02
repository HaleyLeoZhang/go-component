package httpserver

import (
	"fmt"
	"github.com/HaleyLeoZhang/go-component/driver/xlog"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Config struct {
	Name           string        `yaml:"name" json:"name"` // 用于 Trace 识别
	Ip             string        `yaml:"ip" json:"ip"`
	Port           int           `yaml:"port" json:"port"`
	Pprof          bool          `yaml:"pprof" json:"pprof"` // true 开启  pprof 性能监控路由
	ReadTimeout    time.Duration `yaml:"readTimeout" json:"readTimeout"`
	WriteTimeout   time.Duration `yaml:"writeTimeout" json:"writeTimeout"`
	MaxHeaderBytes int           `yaml:"maxHeaderBytes" json:"maxHeaderBytes"`
}

func Run(c *Config, routersInit *gin.Engine) {
	addrString := fmt.Sprintf("%s:%v", c.Ip, c.Port)

	if c.Pprof {
		// pprof 相关说明 http://www.hlzblog.top/article/74.html
		xlog.Info("Enabled pprof")
		Wrap(routersInit)
	}

	server := &http.Server{
		Addr:           addrString,
		Handler:        routersInit,
		ReadTimeout:    c.ReadTimeout,
		WriteTimeout:   c.WriteTimeout,
		MaxHeaderBytes: c.MaxHeaderBytes,
	}
	xlog.Infof("Start http server listening %s", addrString)
	err := server.ListenAndServe()
	if err != nil {
		xlog.Errorf("HttpServer.Err %+v", err)
	}
}
