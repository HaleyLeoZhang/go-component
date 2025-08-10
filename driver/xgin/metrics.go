package xgin

import (
	"fmt"
	"github.com/HaleyLeoZhang/go-component/driver/xmetric"
	"github.com/gin-gonic/gin"
	"net/url"
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
		// 解析路由
		l, _ := url.Parse(c.Request.URL.String())
		path := l.Path
		method := c.Request.Method
		httpStatus := fmt.Sprintf("%v", c.Writer.Status()) // 响应状态码
		// 记录指标
		xmetric.MetricHttpResponse.WithLabelValues(method, path, httpStatus).Observe(float64(takeTime.Milliseconds()))
		xmetric.MetricHttpRequestCount.WithLabelValues(method, path).Inc()
	}
}
