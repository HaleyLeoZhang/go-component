package xkafka

import (
	"context"
	"encoding/json"
	"github.com/HaleyLeoZhang/go-component/driver/xlog"
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

const TopicName = "email_tester"

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
	ctx, cancel := context.WithCancel(context.Background())
	err := StartKafkaConsumer(ctx, ConsumerOption{
		Conf:        config.Kafka,        // Kafka 配置
		Topic:       []string{TopicName}, // 消费Topic列表
		Group:       "email_consumer",    // Consumer group name
		Batch:       10,                  // 每次拉取的消息数
		Procs:       2,                   // 并发处理消息的数量
		PollTimeout: 3 * time.Second,
		Handler:     testHandler, // 处理单条消息的函数
		Mode:        ModeBatch,   // 消费模式
	})
	if err != nil {
		t.Fatalf("%+v", err)
	}
	// 因为上面 Start() 方法不阻塞，为了消费者正常消费，请不要让主进程退出
	xlog.Infof("consumer going to shutdown")
	<-time.After(1 * time.Minute)
	cancel()
	<-time.After(4 * time.Second)
	xlog.Infof("consumer done")
}

func testHandler(ctx context.Context, message *sarama.ConsumerMessage) error {
	xlog.Infof("offset(%v) message(%v) ", message.Offset, string(message.Value))
	return nil
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
	testTopic := TopicName
	d := NewProducer(config.Kafka)
	one := &TestSmtp{}
	one.SenderName = "测试同步"
	bs, _ := json.Marshal(one)
	err := d.SendMsg(testTopic, bs)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	one.SenderName = "测试Key"
	bs, _ = json.Marshal(one)
	err = d.SendMsgByKey(testTopic, "key2333", bs)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	one.SenderName = "测试异步"
	bs, _ = json.Marshal(one)
	err = d.SendMsgAsync(testTopic, bs)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	one.SenderName = "测试异步key"
	bs, _ = json.Marshal(one)
	err = d.SendMsgAsyncByKey(testTopic, "key2333", bs)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	<-time.After(3 * time.Second) // 等待异步刷入
}
