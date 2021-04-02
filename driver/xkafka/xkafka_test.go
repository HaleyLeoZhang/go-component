package xkafka

import (
	"context"
	"encoding/json"
	"github.com/HaleyLeoZhang/go-component/driver/xlog"
	"github.com/HaleyLeoZhang/go-component/errgroup"
	"github.com/Shopify/sarama"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type TestConfig struct {
	Kafka *Config `yaml:"kafka"`
}

var config = &TestConfig{}

//var ctx = context.Background()

func TestMain(m *testing.M) {
	InitConfig()
	os.Exit(m.Run())
}

func InitConfig() {
	var yamlFile string
	yamlFile, err := filepath.Abs("./app.yaml") // 示例的kafka配置文件请看这个文件
	if err != nil {
		panic(err)
	}
	yamlRead, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlRead, config)
	if err != nil {
		panic(err)
	}
}

type TestSmtp struct {
	Subject      string   `json:"subject"`
	SenderName   string   `json:"sender_name"`
	Body         string   `json:"body"`
	Receiver     []string `json:"receiver"`
	ReceiverName []string `json:"receiver_name"`
	Attachment   []string `json:"attachment"`
	Remark       []string `json:"remark"`
}

// 消费者组测试
func TestConsumer(t *testing.T) {
	consumerGroup := "email_consumer"
	topicList := []string{"biz_email"}
	d := NewConsumer(config.Kafka, topicList, consumerGroup)
	d.Consumer.RegisterHandler(handlerExampleStart)
	err := d.Consumer.Start() // 请注意：内部消息者是异步的，此方法不会产生阻塞
	if err != nil {
		xlog.Errorf("handlerExampleStart  Err(%+v)", err)
	}
	// 因为上面 Start() 方法不阻塞，为了消费者正常消费，请不要让主进程退出
	a := make(chan int)
	<-a
	xlog.Infof("consumer done")
}

func handlerExampleStart(session *ConsumerSession, msgs <-chan *sarama.ConsumerMessage) (errKafka error) {
	//s.wg.Add(1)
	//defer s.wg.Done() // 杀进程的时候，等待下面停止消费
	fun := IteratorBatchFetch(session, msgs, 10, 1)
	for {
		kafkaMessages, ok := fun()
		if !ok {
			xlog.Infof("ok(%v)", ok)
			return
		}
		messages := mergeEmailMessageHandler(kafkaMessages)
		eg := &errgroup.Group{}
		eg.GOMAXPROCS(1)
		for _, business := range messages {
			tmp := business
			eg.Go(func(context.Context) error {
				xlog.Infof("当前数据 (%v)", tmp)
				return nil
			})
		}
		_ = eg.Wait()
	}
}

// 消息转对应结构体--并去重
func mergeEmailMessageHandler(msgs []*sarama.ConsumerMessage) (batchList []*TestSmtp) {
	batchList = make([]*TestSmtp, 0, len(msgs))
	for _, msg := range msgs {
		batchInfo := &TestSmtp{}
		err := json.Unmarshal(msg.Value, &batchInfo)
		if err != nil {
			xlog.Errorf("kafkaMergeMsgs topic(%s) val(%s) Err(%+v) ", msg.Topic, string(msg.Value), err)
			continue
		}
		batchList = append(batchList, batchInfo)
		xlog.Infof("kafkaMergeMsgs message(%+v) success topic(%+v) partition(%d) offset(%d)", batchInfo, msg.Topic, msg.Partition, msg.Offset)
	}
	return
}

// 生产者测试
func TestProducer(t *testing.T) {
	testTopic := "biz_email"
	d := NewProducer(config.Kafka)
	one := &TestSmtp{}
	one.SenderName = "测试同步"
	err := d.SendMsg(testTopic, one)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	one.SenderName = "测试Key"
	err = d.SendMsgByKey(testTopic, "key2333", one)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	one.SenderName = "测试异步"
	err = d.SendMsgAsync(testTopic, one)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	one.SenderName = "测试异步key"
	err = d.SendMsgAsyncByKey(testTopic, "key2333", one)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	<-time.After(3 * time.Second) // 等待异步刷入
}
