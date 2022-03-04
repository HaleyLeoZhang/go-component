package conf

import (
	"time"
)

type Consumer struct {
	Group             Group          `yaml:"group"`
	Retry             ConsumerRetry  `yaml:"retry"`
	Fetch             Fetch          `yaml:"fetch"`
	MaxWaitTime       time.Duration  `yaml:"max_wait_time"`
	MaxProcessingTime time.Duration  `yaml:"max_processing_time"`
	Return            ConsumerReturn `yaml:"return"`
	Offsets           Offsets        `yaml:"offsets"`
}

type Group struct {
	Session   Session   `yaml:"session"`
	Heartbeat Heartbeat `yaml:"heart_beat"`
	Rebalance Rebalance `yaml:"rebalance"`
}

type ConsumerRetry struct {
	Backoff time.Duration `yaml:"backoff"`
}

type Fetch struct {
	Min     int32 `yaml:"min"`
	Default int32 `yaml:"default"`
	Max     int32 `yaml:"max"`
}

type ConsumerReturn struct {
	Errors bool `yaml:"errors"`
}

type Offsets struct {
	AutoCommit     AutoCommit    `yaml:"auto_commit"`
	CommitInterval time.Duration `yaml:"commit_interval"`
	Initial        int64         `yaml:"initial"`
	Retention      time.Duration `yaml:"retention"`
	Retry          OffsetsRetry  `yaml:"retry"`
}

type AutoCommit struct {
	// Whether or not to auto-commit updated offsets back to the broker.
	// (default enabled).
	Enable bool `yaml:"enable"`

	// How frequently to commit updated offsets. Ineffective unless
	// auto-commit is enabled (default 1s)
	Interval time.Duration `yaml:"interval"`
}

type OffsetsRetry struct {
	Max int `yaml:"max"`
}

type Session struct {
	TimeOut time.Duration `yaml:"timeout"`
}

type Heartbeat struct {
	Interval time.Duration `yaml:"interval"`
}

type Rebalance struct {
	Timeout time.Duration  `yaml:"timeout"`
	Retry   RebalanceRetry `yaml:"retry"`
}

type RebalanceRetry struct {
	Max     int           `yaml:"max"`
	Backoff time.Duration `yaml:"backoff"`
}
