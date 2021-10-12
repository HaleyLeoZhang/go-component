package xredis

import "time"

//  DB使用默认0库
type Config struct {
	Name         string        `yaml:"name" json:"name"` // 用于 Trace 识别
	Proto        string        `yaml:"proto" json:"proto"`
	Addr         string        `yaml:"addr" json:"addr"`
	Auth         string        `yaml:"auth" json:"auth"`
	DialTimeout  time.Duration `yaml:"dialTimeout" json:"dialTimeout"`
	ReadTimeout  time.Duration `yaml:"readTimeout" json:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout" json:"writeTimeout"`
	SlowLog      time.Duration `yaml:"slowLog" json:"slowLog"`
	Pool         PoolConfig    `yaml:"pool" json:"pool"`
}

type PoolConfig struct {
	MaxActive   int           `yaml:"maxActive" json:"maxActive"`
	MaxIdle     int           `yaml:"maxIdle" json:"maxIdle"`
	IdleTimeout time.Duration `yaml:"idleTimeout" json:"idleTimeout"`
}
