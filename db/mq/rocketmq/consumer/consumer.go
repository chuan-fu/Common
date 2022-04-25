package consumer

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/chuan-fu/Common/zlog"
)

type Consumer struct {
	ConsumerInfos []ConsumerCfg
	pushConsumers map[string]rocketmq.PushConsumer
}

type ConsumerCfg struct {
	Topic        string
	Group        string
	Options      []consumer.Option
	ConsumerFunc func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error)
}

func NewConsumer(infos []ConsumerCfg) *Consumer {
	return &Consumer{
		ConsumerInfos: infos,
		pushConsumers: make(map[string]rocketmq.PushConsumer),
	}
}

func (c *Consumer) InitConsumer() (err error) {
	// 消费者客户端
	for k := range c.ConsumerInfos {
		v := &c.ConsumerInfos[k]
		var pc rocketmq.PushConsumer
		pc, err = rocketmq.NewPushConsumer(
			v.Options...,
		)
		if err != nil {
			log.Error(err)
			return
		}
		// 订阅消费
		if err = pc.Subscribe(v.Topic, consumer.MessageSelector{}, v.ConsumerFunc); err != nil {
			log.Error(err)
			return err
		}

		if err = pc.Start(); err != nil {
			log.Error(err)
			return err
		}
		log.Infof("start consumer topic %s,group %s", v.Topic, v.Group)
		c.pushConsumers[v.Group] = pc
	}
	return nil
}
