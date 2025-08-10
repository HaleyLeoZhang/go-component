package xkafka

import (
	"github.com/HaleyLeoZhang/go-component/driver/xkafka/conf"
	"github.com/Shopify/sarama"
)

type Config struct {
	Admin             conf.Admin    `yaml:"admin"`
	Consumer          conf.Consumer `yaml:"consumer"`
	Metadata          conf.Metadata `yaml:"metadata"`
	Net               conf.Net      `yaml:"net"`
	Producer          conf.Producer `yaml:"producer"`
	ClientID          string        `yaml:"client_id"`
	ChannelBufferSize int           `yaml:"channel_buffer_size"`
	Version           int           `yaml:"version"`
	BrokersAddr       []string      `yaml:"brokers_addr,flow"`
}

func (config *Config) GetSaramaConf() *sarama.Config {
	//配置详情见 当前驱动下的 app.yml
	defaultConfig := sarama.NewConfig()

	if config.Admin.TimeOut != 0 {
		defaultConfig.Admin.Timeout = config.Admin.TimeOut
	}

	if config.Net.MaxOpenRequests != 0 {
		defaultConfig.Net.MaxOpenRequests = config.Net.MaxOpenRequests
	}
	if config.Net.DialTimeout != 0 {
		defaultConfig.Net.DialTimeout = config.Net.DialTimeout
	}

	if config.Net.ReadTimeout != 0 {
		defaultConfig.Net.ReadTimeout = config.Net.ReadTimeout
	}

	if config.Net.WriteTimeout != 0 {
		defaultConfig.Net.WriteTimeout = config.Net.WriteTimeout
	}

	if config.Net.SASL.Enable != false {
		defaultConfig.Net.SASL.Enable = config.Net.SASL.Enable
	}

	if config.Net.SASL.Handshake != false {
		defaultConfig.Net.SASL.Handshake = config.Net.SASL.Handshake
	}
	if config.Net.SASL.User != "" {
		defaultConfig.Net.SASL.User = config.Net.SASL.User
	}
	if config.Net.SASL.Password != "" {
		defaultConfig.Net.SASL.Password = config.Net.SASL.Password
	}
	if config.Net.TLS.Enable != false {
		defaultConfig.Net.TLS.Enable = config.Net.TLS.Enable
	}
	if config.Net.KeepAlive != 0 {
		defaultConfig.Net.KeepAlive = config.Net.KeepAlive
	}

	if config.Metadata.Retry.Max != 0 {
		defaultConfig.Metadata.Retry.Max = config.Metadata.Retry.Max
	}

	if config.Metadata.Retry.Backoff != 0 {
		defaultConfig.Metadata.Retry.Backoff = config.Metadata.Retry.Backoff
	}

	if config.Metadata.RefreshFrequency != 0 {
		defaultConfig.Metadata.RefreshFrequency = config.Metadata.RefreshFrequency
	}
	if config.Metadata.Full != false {
		defaultConfig.Metadata.Full = config.Metadata.Full
	}

	if config.Producer.MaxMessageBytes != 0 {
		defaultConfig.Producer.MaxMessageBytes = config.Producer.MaxMessageBytes
	}

	if config.Producer.Timeout != 0 {
		defaultConfig.Producer.Timeout = config.Producer.Timeout
	}
	if config.Producer.RequiredAcks != 0 {
		defaultConfig.Producer.RequiredAcks = sarama.RequiredAcks(config.Producer.RequiredAcks)

	}
	if config.Producer.Retry.Max != 0 {
		defaultConfig.Producer.Retry.Max = config.Producer.Retry.Max
	}
	if config.Producer.Retry.Backoff != 0 {
		defaultConfig.Producer.Retry.Backoff = config.Producer.Retry.Backoff
	}
	if config.Producer.Return.Errors != false {
		defaultConfig.Producer.Return.Errors = config.Producer.Return.Errors
	}
	if config.Producer.Return.Successes != false {
		defaultConfig.Producer.Return.Successes = config.Producer.Return.Successes
	}
	if config.Producer.CompressionLevel != 0 {
		defaultConfig.Producer.CompressionLevel = config.Producer.CompressionLevel
	}

	if config.Producer.Partitioner != "" {
		switch config.Producer.Partitioner {
		case "rr":
			defaultConfig.Producer.Partitioner = sarama.NewRoundRobinPartitioner
		case "hash":
			defaultConfig.Producer.Partitioner = sarama.NewHashPartitioner
		case "random":
			defaultConfig.Producer.Partitioner = sarama.NewRandomPartitioner
		case "manual":
			defaultConfig.Producer.Partitioner = sarama.NewManualPartitioner
		}
	}

	if config.Consumer.Fetch.Min != 0 {
		defaultConfig.Consumer.Fetch.Min = config.Consumer.Fetch.Min
	}

	if config.Consumer.Fetch.Max != 0 {
		defaultConfig.Consumer.Fetch.Max = config.Consumer.Fetch.Max
	}
	if config.Consumer.Fetch.Default != 0 {
		defaultConfig.Consumer.Fetch.Default = config.Consumer.Fetch.Default
	}
	if config.Consumer.Retry.Backoff != 0 {
		defaultConfig.Consumer.Retry.Backoff = config.Consumer.Retry.Backoff
	}
	if config.Consumer.MaxWaitTime != 0 {
		defaultConfig.Consumer.MaxWaitTime = config.Consumer.MaxWaitTime
	}
	if config.Consumer.MaxProcessingTime != 0 {
		defaultConfig.Consumer.MaxProcessingTime = config.Consumer.MaxProcessingTime
	}
	if config.Consumer.Return.Errors != false {
		defaultConfig.Consumer.Return.Errors = config.Consumer.Return.Errors
	}
	if config.Consumer.Offsets.CommitInterval != 0 {
		defaultConfig.Consumer.Offsets.CommitInterval = config.Consumer.Offsets.CommitInterval
	}
	if config.Consumer.Offsets.Initial != 0 {
		defaultConfig.Consumer.Offsets.Initial = config.Consumer.Offsets.Initial
	}
	if config.Consumer.Offsets.Retention != 0 {
		defaultConfig.Consumer.Offsets.Retention = config.Consumer.Offsets.Retention
	}
	if config.Consumer.Offsets.Retry.Max != 0 {
		defaultConfig.Consumer.Offsets.Retry.Max = config.Consumer.Offsets.Retry.Max
	}
	if config.Consumer.Group.Session.TimeOut != 0 {
		defaultConfig.Consumer.Group.Session.Timeout = config.Consumer.Group.Session.TimeOut
	}
	if config.Consumer.Group.Heartbeat.Interval != 0 {
		defaultConfig.Consumer.Group.Heartbeat.Interval = config.Consumer.Group.Heartbeat.Interval
	}
	if config.Consumer.Group.Rebalance.Timeout != 0 {
		defaultConfig.Consumer.Group.Rebalance.Timeout = config.Consumer.Group.Rebalance.Timeout
	}
	if config.Consumer.Group.Rebalance.Retry.Max != 0 {
		defaultConfig.Consumer.Group.Rebalance.Retry.Max = config.Consumer.Group.Rebalance.Retry.Max
	}
	if config.Consumer.Group.Rebalance.Retry.Backoff != 0 {
		defaultConfig.Consumer.Group.Rebalance.Retry.Backoff = config.Consumer.Group.Rebalance.Retry.Backoff

	}

	if config.ClientID != "" {
		defaultConfig.ClientID = config.ClientID
	}

	if config.ChannelBufferSize != 0 {
		defaultConfig.ChannelBufferSize = config.ChannelBufferSize
	}

	if config.Version != 0 {
		defaultConfig.Version = sarama.SupportedVersions[config.Version]
	}
	return defaultConfig
}
