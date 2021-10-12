package xgin

import "time"

type Config struct {
	Name    string        `yaml:"name" json:"name"` // 用于 Trace 识别
	Debug   bool          `yaml:"debug" json:"debug"`
	Timeout time.Duration `yaml:"timeout" json:"timeout"`
}

