package rocketmq

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/chuan-fu/Common/zlog"
	"github.com/pkg/errors"
)

type CustomerConf struct {
	NameServer       []string `required:"true" json:"nameServer" yaml:"nameServer"`
	Retry            int      `json:"retry" yaml:"retry"`
	InstanceName     string   `json:"instanceName" yaml:"instanceName"`
	NameServerDomain string   `json:"nameServerDomain" yaml:"nameServerDomain"`
	Namespace        string   `json:"namespace" yaml:"namespace"`
	GroupName        string   `json:"groupName" yaml:"groupName"`
}

func ConnectCustomer(conf CustomerConf) error {
	if err := connectCustomer(&conf); err != nil {
		err = errors.Wrap(err, "RocketMQ消费者连接错误")
		log.Error(err)
		return err
	}
	return nil
}

func connectCustomer(conf *CustomerConf) error {
	if len(conf.NameServer) == 0 {
		return errors.New("connectCustomer: conf.NameServer 为空")
	}

	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(conf.GroupName),
		consumer.WithNameServer(conf.NameServer),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithConsumeMessageBatchMaxSize(10),
	)

	/*
		if len(conf.NameServer) == 0 {
			return errors.New("connectProducer: conf.NameServer 为空")
		}

		opts := []producer.Option{
			producer.WithNameServer(conf.NameServer),
			producer.WithRetry(conf.Retry),
		}
		if conf.InstanceName != "" {
			opts = append(opts, producer.WithInstanceName(conf.InstanceName))
		}
		if conf.NameServerDomain != "" {
			opts = append(opts, producer.WithNameServerDomain(conf.NameServerDomain))
		}
		if conf.Namespace != "" {
			opts = append(opts, producer.WithNamespace(conf.Namespace))
		}
		if conf.GroupName != "" {
			opts = append(opts, producer.WithGroupName(conf.GroupName))
		}

		p, err := rocketmq.NewProducer(opts...)
		if err != nil {
			log.Error(err)
			return err
		}
	*/
}
