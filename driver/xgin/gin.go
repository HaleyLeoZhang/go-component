package xgin

import (
	"context"
	"github.com/HaleyLeoZhang/go-component/driver/xlog"
	"github.com/gin-gonic/gin"
	"time"
)

// 超时时间
var timeout time.Duration

// 在原来的包上，增加功能
type Gin struct {
	GinContext *gin.Context // xgin context
	C          context.Context
	Cancel     context.CancelFunc
}

func New(c *Config) *gin.Engine {
	gin.SetMode("release")
	if true == c.Debug {
		gin.SetMode("debug")
	}
	timeout = c.Timeout

	r := gin.New()
	if true == c.Debug {
		r.Use(gin.Logger())
	}
	r.Use(HttpMetrics()) // 默认 Metrics 打点
	//r.Use(gin.Recovery()) // 外部 recovery
	return r
}

// 获取单个Gin
func NewGin(c *gin.Context) *Gin {
	var (
		ctx       = context.Background()
		cancelFun context.CancelFunc
	)
	ctx = xlog.WithLogID(ctx, xlog.GenerateLogID())
	ctx, cancelFun = context.WithTimeout(ctx, timeout)
	//defer cancelFun()

	o := &Gin{
		GinContext: c,
		C:          ctx,
		Cancel:     cancelFun,
	}
	return o
}
