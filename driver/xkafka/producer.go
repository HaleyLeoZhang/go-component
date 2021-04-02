package xkafka

import (
	"encoding/json"
	"github.com/HaleyLeoZhang/go-component/driver/xlog"
	"github.com/Shopify/sarama"
)

// -----------------------------------
// 集成Kafka生产者
// -----------------------------------

// 生产者
type Producer struct {
	syncProducer  *SyncProducer
	asyncProducer *AsyncProducer
}

// 生产者实例
func NewProducer(cfg *Config) (p *Producer) {
	var (
		err error
	)
	p = &Producer{}
	p.syncProducer, err = InitSyncProducer(cfg)
	if err != nil {
		panic(err)
	}
	p.asyncProducer, err = InitAsyncProducer(cfg)
	if err != nil {
		panic(err)
	}
	return
}

// ----------------------------------------
//   同步发送消息
// ----------------------------------------
//   - 适合记录用户行为，无消息丢失
// ----------------------------------------

// - 轮询节点插入
func (p *Producer) SendMsg(topic string, message interface{}) (err error) {
	bs, err := json.Marshal(message)
	if err != nil {
		xlog.Errorf("SendMsg To topic(%s) error(%+v) message(%+v)", topic, err, message)
		return
	}
	partition, commitId, err := p.syncProducer.SendMessage(topic, bs)
	if err != nil {
		xlog.Errorf("SendMsg To topic(%s) error(%v) partition(%d),commitId(%d)", topic, err, partition, commitId)
		return
	}
	xlog.Infof("SendMsg to topic(%s) message(%s) partition(%d),commitId(%d)", topic, string(bs), partition, commitId)
	return
}

// - 依据key计算分区刷入
func (p *Producer) SendMsgByKey(topic string, key string, message interface{}) (err error) {
	bs, err := json.Marshal(message)
	if err != nil {
		xlog.Errorf("SendMsg To topic(%s) error(%+v) message(%+v)", topic, err, message)
		return
	}
	saraMsg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(bs),
	}
	err = p.syncProducer.SendMessages([]*sarama.ProducerMessage{saraMsg})
	if err != nil {
		xlog.Errorf("SendMsg To topic(%s) error(%v)", topic, err)
		return
	}
	xlog.Infof("SendMsg to topic(%s) message(%s)", topic, string(bs))
	return
}

// ----------------------------------------
//   异步发送消息
// ----------------------------------------
//   - 适合刷数据时，批量推送消息
// ----------------------------------------

// - 轮询节点插入
func (p *Producer) SendMsgAsync(topic string, message interface{}) (err error) {
	bs, err := json.Marshal(message)
	if err != nil {
		xlog.Errorf("SendMsgAsync To topic(%s) error(%+v) message(%+v)", topic, err, message)
		return
	}
	err = p.asyncProducer.AsyncPublish(topic, bs)
	if err != nil {
		xlog.Errorf("SendMsgAsync To topic(%s) error(%v)", topic, err)
		return
	}
	xlog.Infof("SendMsgAsync to topic(%s) message(%s) ", topic, string(bs))
	return
}

// - 依据key计算分区刷入
func (p *Producer) SendMsgAsyncByKey(topic string, key string, message interface{}) (err error) {
	bs, err := json.Marshal(message)
	if err != nil {
		xlog.Errorf("SendMsgAsync To topic(%s) error(%+v) message(%+v)", topic, err, message)
		return
	}
	saraMsg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(bs),
	}
	err = p.asyncProducer.PublishOriginMessage(topic, saraMsg)
	if err != nil {
		xlog.Errorf("SendMsgAsyncByKey To topic(%s) error(%v)", topic, err)
		return
	}
	xlog.Infof("SendMsgAsyncByKey to topic(%s) message(%s) ", topic, string(bs))
	return
}
