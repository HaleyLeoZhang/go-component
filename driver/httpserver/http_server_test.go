package httpserver

import (
	"github.com/HaleyLeoZhang/go-component/driver/xgin"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"testing"
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

	Run(cfg.HttpServer, ginEngine)
}
