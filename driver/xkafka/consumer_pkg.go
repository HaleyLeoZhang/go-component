package xkafka

import (
	"context"
	"errors"
	"github.com/HaleyLeoZhang/go-component/driver/xlog"
	"github.com/Shopify/sarama"
	"strings"
	"sync"
)

var (
	NoInitErr = errors.New("未初始化消费方法或未初始化topic") // NoInitErr 未初始化错误
)

type GroupConsumerHandler func(session *ConsumerSession, msgs <-chan *sarama.ConsumerMessage) error
type GroupConsumerSessionHandler func(session sarama.ConsumerGroupSession) error

type GroupConsumer struct {
	topics []string
	group  string

	client   sarama.Client
	consumer sarama.ConsumerGroup
	handler  GroupConsumerHandler
	setup    GroupConsumerSessionHandler
	cleanup  GroupConsumerSessionHandler

	ctx    context.Context
	cancel context.CancelFunc
	waiter sync.WaitGroup
	isInit bool
}

//************
// ConsumerSession
//************
type ConsumerSession struct {
	session sarama.ConsumerGroupSession
}

func (s *ConsumerSession) Context() context.Context {
	if s.session != nil {
		return s.session.Context()
	}
	return nil
}

func (cs *ConsumerSession) Commit(msg *sarama.ConsumerMessage) {
	cs.session.MarkMessage(msg, "")
}

func (cs *ConsumerSession) MarkOffset(topic string, partition int32, offset int64, metadata string) {
	cs.session.MarkOffset(topic, partition, offset, metadata)
}

func (cs *ConsumerSession) Nack(msg *sarama.ConsumerMessage) {

}

//************
// GroupConsumer
//************
func (g *GroupConsumer) Setup(x sarama.ConsumerGroupSession) error {
	if g.setup != nil {
		return g.setup(x)
	}
	return nil
}

func (g *GroupConsumer) Cleanup(x sarama.ConsumerGroupSession) error {
	if g.cleanup != nil {
		return g.cleanup(x)
	}

	return nil
}

//注册handler
func (g *GroupConsumer) RegisterHandler(handler GroupConsumerHandler) {
	g.handler = handler
}

//注册Cleanup handler
func (g *GroupConsumer) RegisterCleanupHandler(handler GroupConsumerSessionHandler) {
	g.cleanup = handler
}

//注册Setup handler
func (g *GroupConsumer) RegisterSetupHandler(handler GroupConsumerSessionHandler) {
	g.setup = handler
}

func (g *GroupConsumer) Errors() <-chan error {
	if g.consumer != nil {
		return g.consumer.Errors()
	}
	return nil
}

func (g *GroupConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	s := &ConsumerSession{
		session: session,
	}
	return g.handler(s, claim.Messages())
}

func (g *GroupConsumer) Start() (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			return
		}
	}()
	if g.handler == nil || g.consumer == nil || len(g.topics) == 0 {
		return NoInitErr
	}
	g.isInit = true

	g.waiter.Add(1)

	go func() {
		defer g.waiter.Done()
		for {
			xlog.Infof("Listening Kafka Group(%v) Topic(%v)",  g.group, strings.Join(g.topics, ","))
			if err := g.consumer.Consume(g.ctx, g.topics, g); err != nil {
				xlog.Errorf("Error from consumer: %v", err)
				continue
			}
			if g.ctx.Err() != nil {
				return
			}
		}
	}()
	return nil
}

func (g *GroupConsumer) Close() error {
	if !g.isInit {
		return nil
	}

	g.cancel()
	g.waiter.Wait()
	if err := g.consumer.Close(); err != nil {
		return err
	}
	return nil
}

//InitGroupConsumer
// topics 消费的 topic 数组, eg. []string{"test1","test2"}
// group 消费组 eg. "test" 这里消费组可以自己定义
// conf kafka 的 conf
func InitGroupConsumer(cfg *Config, topicList []string, group string, handler GroupConsumerHandler) (*GroupConsumer, error) {
	var err error
	c := &GroupConsumer{}
	c.client, err = sarama.NewClient(cfg.BrokersAddr, cfg.GetSaramaConf())
	if err != nil {
		return nil, err
	}
	c.consumer, err = sarama.NewConsumerGroupFromClient(group, c.client)
	if err != nil {
		return nil, err
	}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.handler = handler
	c.topics = topicList
	c.group = group
	return c, nil
}
