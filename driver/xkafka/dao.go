package xkafka

import (
	"context"
	"fmt"
	"github.com/HaleyLeoZhang/go-component/driver/xlog"
	"github.com/HaleyLeoZhang/go-component/errgroup"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"strings"
	"time"
)

// -----------------------------------
// 集成Kafka消费者
// -----------------------------------

//Consumer
type Consumer struct {
	Consumer *GroupConsumer
	option   ConsumerOption
	ctx      context.Context
	handler  GroupConsumerHandler
}

type ConsumeMode int

// 消费模式
const (
	ModeBatch ConsumeMode = 1 // 批量拉取
)

type Handler func(ctx context.Context, message *sarama.ConsumerMessage) error

type ConsumerOption struct {
	Conf        *Config       // Kafka 配置
	Topic       []string      // 消费Topic列表
	Group       string        // Consumer group name
	Batch       int           // 每次拉取的消息数
	Procs       int           // 并发处理消息的数量
	PollTimeout time.Duration // 每次拉取消息的超时时间
	Handler     Handler       // 处理单条消息的函数
	Mode        ConsumeMode   // 消费模式
	//PreHandler  *GroupConsumerHandler             // 预处理消息的函数: 合并消息，然后分组去重等场景 TODO
}

// 初始化并开启消费者
// 注: 这里入参 ctx 得是 withcancel 的那种
func StartKafkaConsumer(ctx context.Context, option ConsumerOption) (err error) {
	if len(option.Topic) == 0 {
		return errors.WithStack(fmt.Errorf("topic不能为空"))
	}
	if len(option.Group) == 0 {
		return errors.WithStack(fmt.Errorf("consumerGroup不能为空"))
	}
	d := &Consumer{
		ctx:    ctx,
		option: option,
	}
	if d.Consumer, err = InitGroupConsumer(option.Conf, option.Topic, option.Group, nil); err != nil {
		panic(err)
	}
	d.Consumer.RegisterSetupHandler(d.Consumer.setupHandler)
	d.Consumer.RegisterCleanupHandler(d.Consumer.cleanupHandler)
	switch option.Mode {
	case ModeBatch:
		d.handler = d.handleBatch
	default:
		return errors.WithStack(fmt.Errorf("请选择消费模型mode"))
	}
	d.Consumer.RegisterHandler(d.handler)
	err = d.Consumer.Start() // 请注意：内部消息者是异步的，此方法不会产生阻塞
	if err != nil {
		return
	}
	go func() {
		<-d.ctx.Done()
		infoTmp := fmt.Sprintf("Stop Kafka Group(%v) Topic(%v)", d.Consumer.group, strings.Join(d.Consumer.topics, ","))
		fmt.Println(infoTmp)
		_ = d.Consumer.Close()
	}()
	return
}

// 批量消费模型
func (d *Consumer) handleBatch(session *ConsumerSession, msgs <-chan *sarama.ConsumerMessage) (errKafka error) {
	fun := iteratorBatchFetch(session, msgs, d.option.Batch, d.option.PollTimeout)
	for {
		messages, ok := fun()
		if !ok {
			return
		}
		eg := &errgroup.Group{}
		eg.GOMAXPROCS(d.option.Procs)
		for _, msgTmp := range messages {
			msg := msgTmp
			eg.Go(func(context.Context) error {
				_ = d.option.Handler(d.ctx, msg) // 抛弃当前Error
				return nil
			})
		}
		err := eg.Wait()
		if err != nil {
			xlog.Errorf("kafka consumer Error(%+v)", err)
		}
	}
}

func (d *Consumer) Start() {
	err := d.Consumer.Start()
	if err == nil {
		go func() {
			<-d.Consumer.ctx.Done()
			_ = d.Consumer.Close()
		}()
	}
}
