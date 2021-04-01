package xkafka

import (
	"errors"
	"github.com/Shopify/sarama"
)

type SyncProducer struct {
	producer sarama.SyncProducer
}

func InitSyncProducer(config *Config) (producer *SyncProducer, err error) {
	if config == nil {
		err = errors.New("conf 不能为空")
		return
	}

	if len(config.BrokersAddr) == 0 {
		err = errors.New("config中BrokersAddr 不能为空")
		return
	}

	syncProducer, err := sarama.NewSyncProducer(config.BrokersAddr, config.GetSaramaConf())
	if err != nil {
		return
	}

	producer = new(SyncProducer)
	producer.producer = syncProducer

	return
}

//如果channel满，就抛错
func (producer *SyncProducer) SendMessage(topic string, msg []byte) (partition int32, offset int64, err error) {
	if producer.producer == nil {
		return 0, 0, errors.New("kafka failed to initialize")
	}

	if topic == "" {
		return 0, 0, errors.New("topic 不能为空")
	}

	partition, offset, err = producer.producer.SendMessage(&sarama.ProducerMessage{Topic: topic, Value: sarama.ByteEncoder(msg)})
	return
}

// 如果channel满，就抛错 需要配合config.producer.partitioner = manual 使用
func (producer *SyncProducer) SendMessageByPartition(topic string, msg []byte, partition int32) (offset int64, err error) {
	if producer.producer == nil {
		return 0, errors.New("kafka failed to initialize")
	}

	if topic == "" {
		return 0, errors.New("topic 不能为空")
	}

	_, offset, err = producer.producer.SendMessage(&sarama.ProducerMessage{Topic: topic, Value: sarama.ByteEncoder(msg), Partition: partition})
	return
}

// 同步发送原生消息, msg内部可指定hash key 和 topic
func (producer *SyncProducer) SendOriginMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	if producer.producer == nil {
		return 0, 0, errors.New("kafka failed to initialize")
	}

	partition, offset, err = producer.producer.SendMessage(msg)
	return
}

func (producer *SyncProducer) SendMessages(message []*sarama.ProducerMessage) (err error) {
	if producer.producer == nil {
		return errors.New("kafka failed to initialize")
	}

	if len(message) == 0 {
		return errors.New("message 不能为空")
	}

	err = producer.producer.SendMessages(message)

	return
}

func (producer *SyncProducer) Close() (err error) {
	err = producer.producer.Close()

	return
}
