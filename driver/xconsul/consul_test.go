package xconsul

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"
)

type TestConfig struct {
	Consul *Config `yaml:"consul"`
}

var (
	cfg = &TestConfig{}
	clt *Client
	ctx = context.Background()
	err error
)

const (
	SERVICE_NAME = "comic.pre.hlzblog.top"
	HTTP_PORT    = 4211
)

func TestMain(m *testing.M) {
	err := loadConfig()
	if err != nil {
		panic(err)
	}
	iniClient()
	m.Run()
}

func loadConfig() (err error) {
	var yamlFile string
	yamlFile, err = filepath.Abs("./app.yml")
	if err != nil {
		return
	}
	yamlRead, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(yamlRead, cfg)
	if err != nil {
		return
	}
	return
}

func iniClient() {
	clt, err = NewClient(cfg.Consul)
	if err != nil {
		panic(err)
	}
}

func TestHttp(t *testing.T) {
	ServerLoad()
}

//Handler 3001
func handlerHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("you are visiting health check api"))
}

//ServerLoad 启动
func ServerLoad() {
	//go func() {
	//	// 注销
	//	<-time.After(time.Second * 10)
	//	_ = clt.Deregister()
	//}()
	// 注册
	err = clt.HttpRegister()
	if err != nil {
		panic(errors.WithStack(err))
	}
	// 定义一个http接口
	http.HandleFunc(cfg.Consul.HealthCheckRouter, handlerHealth)
	err = http.ListenAndServe(fmt.Sprintf(":%v", cfg.Consul.HttpPort), nil)
	if err != nil {
		panic(errors.WithStack(err))
	}

}
