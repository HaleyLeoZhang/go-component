package xgin

import (
	"github.com/HaleyLeoZhang/go-component/driver/xmetric"
	"github.com/gin-gonic/gin"
	"time"
)

func HttpMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求前时间
		nowTime := time.Now()
		// 请求处理
		c.Next()
		// 计算时间消耗
		takeTime := time.Since(nowTime)
		url := c.Request.URL.String()
		method := c.Request.Method
		// 记录指标
		xmetric.MetricHttpResponse.WithLabelValues(method, url).Observe(float64(takeTime.Milliseconds()))
		xmetric.MetricHttpRequestCount.WithLabelValues(method, url).Inc()
	}
}
