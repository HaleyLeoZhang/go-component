package xelastic

import (
	"time"
)

type Config struct {
	Addrs              []string      `yaml:"addrs" `
	Username           string        `yaml:"username" `
	Password           string        `yaml:"password" `
	HealthCheckEnabled bool          `yaml:"healthCheckEnabled" `
	SnifferEnabled     bool          `yaml:"snifferEnabled" `
	HealthTimeout      time.Duration `yaml:"healthTimeout" `
	SnifferTimeout     time.Duration `yaml:"snifferTimeout" `
}
