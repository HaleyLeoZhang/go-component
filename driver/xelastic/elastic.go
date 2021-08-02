package xelastic

import (
	"fmt"
	v7 "github.com/olivere/elastic/v7"
	xhttp "net/http"
	"time"
)

type Config struct {
	Addrs              []string      `yaml:"addrs" json:"addrs"`
	Username           string        `yaml:"username" json:"username"`
	Password           string        `yaml:"password" json:"password"`
	HealthcheckEnabled bool          `yaml:"healthcheckEnabled" json:"healthcheckEnabled"`
	SnifferEnabled     bool          `yaml:"snifferEnabled" json:"snifferEnabled"`
	HealthTimeOut      time.Duration `yaml:"healthtimeout" json:"healthtimeout"`
	SnifferTimeout     time.Duration `yaml:"snifferTimeout" json:"snifferTimeout"`
	V7                 struct {
		MaxIdleConnsPerHost int `yaml:"maxIdleConnsPerHost" json:"maxIdleConnsPerHost"`
		MaxIdleConns        int `yaml:"maxIdleConns" json:"maxIdleConns"`
		TimeOut             int `yaml:"timeOut" json:"timeOut"`
		KeepAlive           int `yaml:"keepAlive" json:"keepAlive"`
	} `yaml:"v7" json:"v7"`
}

const (
	es7                = 7
	clientTimeout      = 30
	clientKeepAlive    = 30
	clientMaxIdleConns = 100
)

func NewV7(esConfig *Config, option ...v7.ClientOptionFunc) (es *v7.Client, err error) {

	if len(esConfig.Addrs) == 0 {
		err = fmt.Errorf("addrs is empty")
		return
	}

	var (
		timeout             = esConfig.V7.TimeOut
		keepAlive           = esConfig.V7.KeepAlive
		maxIdleConns        = esConfig.V7.MaxIdleConns
		maxIdleConnsPerHost = esConfig.V7.MaxIdleConnsPerHost
	)

	if esConfig.HealthTimeOut == 0 {
		esConfig.HealthTimeOut = v7.DefaultHealthcheckTimeout
	}

	if esConfig.SnifferTimeout == 0 {
		esConfig.SnifferTimeout = v7.DefaultSnifferTimeout
	}

	if esConfig.SnifferTimeout == 0 {
		esConfig.SnifferTimeout = v7.DefaultSnifferTimeout
	}

	if timeout == 0 {
		timeout = clientTimeout
	}

	if keepAlive == 0 {
		keepAlive = clientKeepAlive
	}

	if maxIdleConns == 0 {
		maxIdleConns = clientMaxIdleConns
	}

	if maxIdleConnsPerHost == 0 {
		maxIdleConnsPerHost = xhttp.DefaultMaxIdleConnsPerHost
	}
	if esConfig.V7.MaxIdleConns == 0 {
		esConfig.SnifferTimeout = v7.DefaultSnifferTimeout
	}

	// 提供些默认必要配置，但还支持elastic更多配置，提供啦注入配置方式
	option = append(option, v7.SetURL(esConfig.Addrs...),
		v7.SetBasicAuth(esConfig.Username, esConfig.Password),
		v7.SetHealthcheck(esConfig.HealthcheckEnabled),
		v7.SetHealthcheckTimeout(esConfig.HealthTimeOut),
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
