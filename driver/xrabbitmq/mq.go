package xrabbitmq

import (
	"fmt"
	"github.com/HaleyLeoZhang/go-component/driver/xlog"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"sync"
)

type AMQP struct {
	// 当前驱动配置项
	Exchange      string
	Queue         string
	PullLimit     int // 拉取数据条数限
	ConsumerLimit int // 消费者并发处理数限制
	// 单例连接
	Conn      *amqp.Connection
	closeFlag bool // 关闭链接
	wg        sync.WaitGroup
	// 读取配置信息
	Conf *Config
}

func (a *AMQP) Push(exchange string, routingKey string, payload []byte) error {
	conn := a.Conn

	ch, err := conn.Channel()
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	defer ch.Close()

	err = ch.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        payload,
		})
	if err != nil {
		err = errors.WithStack(err)
		return err
	}

	return nil
}

func (a *AMQP) Pull(callback func([]byte) error) error {
	conn := a.Conn

	ch, err := conn.Channel()
	if err != nil {
		xlog.Errorf("RabbitMq.Channel.Err(%+v).Exchange(%v).Queue(%v)", err, a.Exchange, a.Queue)
		return err
	}
	defer ch.Close()

	err = ch.Qos(a.PullLimit, 0, false)
	if err != nil {
		xlog.Errorf("RabbitMq.Channel.Qos.Err(%+v).Exchange(%v).Queue(%v)", err, a.Exchange, a.Queue)
		return err
	}

	delivery, err := ch.Consume(
		a.Queue,     // name
		a.Conf.Name, // consumerTag,
		false,       // noAck
		false,       // exclusive
		false,       // noLocal
		false,       // noWait
		nil,         // arguments
	)

	if a.PullLimit == 0 {
		a.PullLimit = 1
	}
	if a.ConsumerLimit == 0 {
		a.ConsumerLimit = 1
	}
	pool := make(chan int, a.PullLimit)
	consumerPool := make(chan int, a.ConsumerLimit)
	for {
		if a.closeFlag == true { // 应用主动关闭的时候
			return nil
		}
		select {
		case d, ok := <-delivery:
			if !ok { // 通道关闭时，退出函数
				a.closeFlag = true
				return nil
			}
			pool <- 1
			go func() {
				defer func() {
					if r := recover(); r != nil {
						xlog.Errorf("Panic %+v", r)
					}
				}()
				a.wg.Add(1)
				consumerPool <- 1
				a.handle(d, callback, pool)
				<-consumerPool
				a.wg.Done()
			}()
		}
	}
}

func (a *AMQP) handle(d amqp.Delivery, callback func([]byte) error, pool chan int) error {
	err := callback(d.Body)

	defer func() {
		<-pool
	}()
	if err != nil {
		xlog.Errorf("RabbitMq.Callback.Err(%+v).Exchange(%v).Queue(%v).Body(%v)", err, a.Exchange, a.Queue, string(d.Body))
		return err
	}

	err = d.Ack(false)
	if err != nil {
		xlog.Errorf("RabbitMq.Ack.Err(%+v).Exchange(%v).Queue(%v).Body(%v)", err, a.Exchange, a.Queue, string(d.Body))
		return err
	}
	xlog.Infof("RabbitMq.Consumer.success.Exchange(%v).Queue(%v).Body(%v)", a.Exchange, a.Queue, string(d.Body))
	return nil
}

func (a *AMQP) Start() error {
	dial := fmt.Sprintf("amqp://%v:%v@%v:%v/", a.Conf.UserName, a.Conf.Password, a.Conf.Host, a.Conf.Port)
	iniConn, err := amqp.Dial(dial)

	if err != nil {
		xlog.Errorf("RabbitMq.Connect.Err(%+v).Conf(%+v)", err, a.Conf)
		return err
	}
	a.Conn = iniConn
	return err
}

func (a *AMQP) QueueDeclare() error {
	// 初始化时，声明 Queue
	chs, err := a.Conn.Channel()
	if _, err := chs.QueueDeclare(a.Queue, true, false, false, false, nil); err != nil {
		xlog.Warnf("queue.declare (%v) err: %s", a.Queue, err)
	}
	defer chs.Close()
	return err
}

func (a *AMQP) BindRoutingKey(routingKey string) error {
	// 初始化时，通过routingKet绑定
	chs, err := a.Conn.Channel()
	defer chs.Close()

	err = chs.QueueBind(a.Queue, routingKey, a.Exchange, false, nil)
	if err != nil {
		xlog.Warnf("queue.bind queue(%v) routingKey(%v) exchange(%v) err: %s", a.Queue, routingKey, a.Exchange, err)
	}
	return err

}
func (a *AMQP) Close() error {
	xlog.Infof("rabbitmq 关闭.Exchange(%v).Queue(%v)", a.Exchange, a.Queue)
	a.closeFlag = true
	a.wg.Wait() // 平滑关闭
	err := a.Conn.Close()
	if err != nil {
		xlog.Errorf("RabbitMq.Close.Err(%+v).Conf(%+v)", err, a.Conf)
	}
	return err
}
