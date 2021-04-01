package xkafka

import (
	"github.com/Shopify/sarama"
	"time"
)

// 消费模型1
// 批量获取kafka消息
// - 超过指定秒数拉取不完，自动退出拉取
func IteratorBatchFetch(session *ConsumerSession, msgs <-chan *sarama.ConsumerMessage, num int, timeout time.Duration) func() ([]*sarama.ConsumerMessage, bool) {
	var (
		lastMsg *sarama.ConsumerMessage
	)
	return func() ([]*sarama.ConsumerMessage, bool) {
		msg := make([]*sarama.ConsumerMessage, 0)
		var ticker *time.Ticker
		ticker = time.NewTicker(time.Second * timeout)
		defer ticker.Stop()
		if lastMsg != nil {
			session.Commit(lastMsg)
			lastMsg = nil
		}
		for {
			select {
			case tempMsg, ok := <-msgs:
				if !ok {
					return msg, ok
				}
				lastMsg = tempMsg
				msg = append(msg, tempMsg)
				if len(msg) >= num {
					return msg, true
				}
			case <-ticker.C:
				if len(msg) > 0 {
					return msg, true
				}
			}
		}
	}
}
