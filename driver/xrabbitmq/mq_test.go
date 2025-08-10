package xrabbitmq

import (
	"context"
	"github.com/HaleyLeoZhang/go-component/driver/xlog"
	"testing"
	"time"
)

var app = &AMQP{}
var config = &Config{}
var routineKey = "email.sender"

func TestMain(m *testing.M) {
	ctx := context.Background()
	config.Host = "192.168.56.110"
	config.Port = 5672
	config.Name = "email_server_local"
	config.UserName = "guest"
	config.Password = "guest"
	app.Conf = config
	// 初始化配置
	app.PullLimit = 3          // 每次最多拉多少条
	app.ConsumerLimit = 2      // 每次最多 多少个消费者
	app.Exchange = "amq.topic" // 交换机名
	app.Queue = "email_sender" // 消费队列
	app.Start()
	app.QueueDeclare()
	app.BindRoutingKey(routineKey) // 初始化约定要绑定的 routing_key
	xlog.Infof(ctx, "RabbitMQ.Init.Exchange (%v) Queue (%v)", app.Exchange, app.Queue)
	m.Run()
	app.Close()
}

func TestService_Push(t *testing.T) {
	app.Push(app.Exchange, routineKey, []byte(`{"title":"For test","content":"unit test","sender_name":"local","receiver":"haleyleozhang@sohu.com","receiver_name":""}`))
}

func TestService_Pull(t *testing.T) {
	app.Pull(func(bytes []byte) error {
		<-time.After(1 * time.Second)
		t.Logf("success")
		return nil
	})
}
