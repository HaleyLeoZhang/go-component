package xredis

import (
	"github.com/gomodule/redigo/redis"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type TestConfig struct {
	Redis *Config `yaml:"redis"`
}

var (
	config    = &TestConfig{}
	redisConn *redis.Pool
)

func TestMain(m *testing.M) {
	InitConfig()
	os.Exit(m.Run())
}

func InitConfig() {
	var yamlFile string
	yamlFile, err := filepath.Abs("./app.yml") // 示例的kafka配置文件请看这个文件
	if err != nil {
		panic(err)
	}
	yamlRead, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlRead, config)
	if err != nil {
		panic(err)
	}
	redisConn, err = NewPool(config.Redis)
	if err != nil {
		panic(err)
	}
}

// 测试查询
func TestGet(t *testing.T) {
	redisKey := "test:redisGetCommand:v1:cached"
	conn := redisConn.Get()
	defer conn.Close()
	bytes, err := conn.Do("get", redisKey)
	if err != nil {
		t.Fatalf("Err(%+v)", err)
	}
	if bytes == nil {
		t.Logf("数据不存在")
	}
	t.Logf("redisKey(%v)value(%v)", redisKey, bytes)
}
