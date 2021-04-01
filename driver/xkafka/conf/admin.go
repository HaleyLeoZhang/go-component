package conf

import "time"

type Admin struct {
	TimeOut time.Duration `yaml:"timeout"`
}
