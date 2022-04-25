package consumer

import (
	"context"
	"testing"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/chuan-fu/Common/zlog"
)

func TestInit(t *testing.T) {
	err := NewConsumer([]ConsumerCfg{
		NewAAAConsumer(),
	}).InitConsumer()
	if err != nil {
		log.Error(err)
		return
	}

	select {}
}

func NewAAAConsumer() ConsumerCfg {
	return ConsumerCfg{
		Topic: "",
		Group: "",
		Options: []consumer.Option{
			consumer.WithInstance("UpdateStatisticStock"),
			consumer.WithNameServer([]string{"192.168.4.34:9876"}),
			consumer.WithGroupName(""),
			consumer.WithConsumerModel(consumer.Clustering),
			consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset),
			// rqconsumer.WithConsumerOrder(true),

			// 不允许批量消费
			consumer.WithConsumeMessageBatchMaxSize(1),
		},
		ConsumerFunc: func(ctx context.Context, ext ...*primitive.MessageExt) (r consumer.ConsumeResult, err error) {
			return
		},
	}
}
