package xkafka

// -----------------------------------
// 集成Kafka消费者
// -----------------------------------

//Consumer
type Consumer struct {
	Consumer *GroupConsumer
}

// 消费者实例
func NewConsumer(cfg *Config, topicList []string, consumerGroup string) (d *Consumer) {
	var (
		err error
	)
	d = &Consumer{}
	if d.Consumer, err = InitGroupConsumer(cfg, topicList, consumerGroup, nil); err != nil {
		panic(err)
	}
	return
}

func (d *Consumer) Start() {
	d.Consumer.Start()
}

// Close close the resource.
func (d *Consumer) Close() {
	d.Consumer.Close()
}
