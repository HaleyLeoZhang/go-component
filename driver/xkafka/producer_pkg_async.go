package xkafka

import (
	"context"
	"errors"
	"fmt"
	"github.com/HaleyLeoZhang/go-component/driver/xlog"
	"github.com/Shopify/sarama"
	"sync/atomic"
)

type AsyncProducer struct {
	Data     chan *sarama.ProducerMessage
	producer sarama.AsyncProducer
	runCount int64 //生产统计
	delayAvg int64 //平均延迟
}

func InitAsyncProducer(config *Config) (producer *AsyncProducer, err error) {
	if config == nil {
		err = errors.New("conf 不能为空")
		return
	}

	if len(config.BrokersAddr) == 0 {
		err = errors.New("config中BrokersAddr 不能为空")
		return
	}

	saramaAsyncProducer, err := sarama.NewAsyncProducer(config.BrokersAddr, config.GetSaramaConf())
	if err != nil {
		return
	}

	//实例化Producer
	producer = new(AsyncProducer)
	producer.producer = saramaAsyncProducer

	//实例化send buffered
	producer.Data = make(chan *sarama.ProducerMessage, config.ChannelBufferSize)

	go producer.send()
	go producer.handleError()

	return
}

//异步发送消息
func (producer *AsyncProducer) AsyncPublish(topic string, msg []byte) (err error) {
	if producer.producer == nil {
		return errors.New("KafkaProducer ")
	}

	if topic == "" {
		return errors.New("topic 不能为空")
	}

	if msg == nil {
		return errors.New("msg 不能为空")
	}

	select {
	case producer.Data <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(msg),
	}:
		return nil
	default:
		return errors.New("kafka channel is full")
	}
}

//异步发送消息 需要配合config.producer.partitioner = custom 使用
func (producer *AsyncProducer) AsyncPublishByPartition(topic string, msg []byte, partition int32) (err error) {
	if producer.producer == nil {
		return errors.New("KafkaProducer ")
	}

	if topic == "" {
		return errors.New("topic 不能为空")
	}

	if msg == nil {
		return errors.New("msg 不能为空")
	}

	select {
	case producer.Data <- &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.ByteEncoder(msg),
		Partition: partition,
	}:
		return nil
	default:
		return errors.New("kafka channel is full")
	}
}

func (producer *AsyncProducer) PublishOriginMessage(topic string, msg *sarama.ProducerMessage) (err error) {
	if producer.producer == nil {
		return errors.New("KafkaProducer")
	}

	if topic == "" {
		return errors.New("topic 不能为空")
	}

	if msg == nil {
		return errors.New("msg 不能为空")
	}

	select {
	case producer.Data <- msg:
		return nil
	default:
		return errors.New("kafka channel is full")
	}
}

func (producer *AsyncProducer) Close() (err error) {
	if producer.producer == nil {
		return fmt.Errorf("kafka failed to initialize")
	}

	return producer.producer.Close()
}

func (producer *AsyncProducer) send() {
	for {
		select {
		case msg := <-producer.Data:
			producer.producer.Input() <- msg
		}
	}
}

func (producer *AsyncProducer) RunCount() int64 {
	return producer.runCount
}
func (producer *AsyncProducer) DelayAvg() int64 {
	return 0
}

func (producer *AsyncProducer) handleError() {
	var (
		err *sarama.ProducerError
		ok  bool
		ctx = context.Background()
	)
	for {
		select {
		case err, ok = <-producer.producer.Errors():
			if err != nil {
				xlog.Errorf(ctx, "producer message error, partition:%d offset:%d key:%v value:%s error(%v)", err.Msg.Partition, err.Msg.Offset, err.Msg.Key, err.Msg.Value, err.Err)
			}

			if !ok {
				xlog.Warn(ctx, "producer ProducerError has be closed, break the handleError goroutine")
				return
			} else {
				atomic.AddInt64(&producer.runCount, 1)
			}

			if err != nil {
				fmt.Printf("producer message error, partition:%d offset:%d key:%v valus:%s error(%v)\n", err.Msg.Partition, err.Msg.Offset, err.Msg.Key, err.Msg.Value, err.Err)
			}

		case <-producer.producer.Successes():
		}

	}
}
