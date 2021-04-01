package conf

import "time"

type Metadata struct {
	Retry            MetadataRetry `yaml:"retry"`
	RefreshFrequency time.Duration `yaml:"refresh_frequency"`
	Full             bool          `yaml:"full"`
	Timeout          time.Duration `yaml:"timeout"`
}

type MetadataRetry struct {
	Max     int           `yaml:"max"`
	Backoff time.Duration `yaml:"back_off"`
}
