package httpserver

import (
	"github.com/HaleyLeoZhang/go-component/driver/xgin"
	"github.com/gin-gonic/gin"
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
	c2.Pprof = true
	c2.Name = "testHttp"
	c2.Ip = "0.0.0.1"
	c2.Port = 80
	Run(c2, ginEngine)
}
