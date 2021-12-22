package xelastic

import (
	"fmt"
	v7 "github.com/olivere/elastic/v7"
)

func NewV7(esConfig *Config, option ...v7.ClientOptionFunc) (es *v7.Client, err error) {

	if len(esConfig.Addrs) == 0 {
		err = fmt.Errorf("addrs is empty")
		return
	}

	if esConfig.HealthTimeout == 0 {
		esConfig.HealthTimeout = v7.DefaultHealthcheckTimeout
	}

	if esConfig.SnifferTimeout == 0 {
		esConfig.SnifferTimeout = v7.DefaultSnifferTimeout
	}

	if esConfig.SnifferTimeout == 0 {
		esConfig.SnifferTimeout = v7.DefaultSnifferTimeout
	}

	// 提供些默认必要配置，但还支持elastic更多配置，提供啦注入配置方式
	option = append(option, v7.SetURL(esConfig.Addrs...),
		v7.SetBasicAuth(esConfig.Username, esConfig.Password),
		v7.SetHealthcheck(esConfig.HealthCheckEnabled),
		v7.SetHealthcheckTimeout(esConfig.HealthTimeout),
		v7.SetSnifferTimeout(esConfig.SnifferTimeout),
		v7.SetSniff(esConfig.SnifferEnabled))
	es, err = v7.NewClient(option...)
	return
}

// ------------------------------------------------------------------------------------------
// 以下为打印查询语句时需要
// ------------------------------------------------------------------------------------------

type TraceLog struct{}

func (TraceLog) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}
