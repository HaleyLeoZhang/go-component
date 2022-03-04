package conf

import "time"

type Producer struct {
	MaxMessageBytes  int           `yaml:"max_message_bytes"`
	RequiredAcks     int16         `yaml:"required_acks"`
	Timeout          time.Duration `yaml:"timeout"`
	CompressionLevel int           `yaml:"compression_level"`
	Idempotent       bool          `yaml:"idempotent"`
	Retry            ProducerRetry `yaml:"retry"`
	Return           Return        `yaml:"return"`
	Partitioner      string        `yaml:"partitioner"` // 支持 hash, manual, rr, random
}

type ProducerRetry struct {
	Max     int           `yaml:"max"`
	Backoff time.Duration `yaml:"backoff"`
}

type Return struct {
	Errors    bool `yaml:"errors"`
	Successes bool `yaml:"successes"`
}
