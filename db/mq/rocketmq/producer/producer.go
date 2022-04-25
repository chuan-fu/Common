package producer

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/chuan-fu/Common/zlog"
	"github.com/pkg/errors"
)

type ProducerConf struct {
	NameServer       []string `required:"true" json:"nameServer" yaml:"nameServer"`
	Retry            int      `json:"retry" yaml:"retry"`
	InstanceName     string   `json:"instanceName" yaml:"instanceName"`
	NameServerDomain string   `json:"nameServerDomain" yaml:"nameServerDomain"`
	Namespace        string   `json:"namespace" yaml:"namespace"`
	GroupName        string   `json:"groupName" yaml:"groupName"`
}

var globalProducer rocketmq.Producer

func GetProducer() rocketmq.Producer {
	return globalProducer
}

func ConnectProducer(conf ProducerConf) error {
	if err := connectProducer(&conf); err != nil {
		err = errors.Wrap(err, "RocketMQ生产者连接错误")
		log.Error(err)
		return err
	}
	return nil
}

func ReloadProducer(conf ProducerConf) error {
	oldProducer := globalProducer
	if err := connectProducer(&conf); err != nil {
		err = errors.Wrap(err, "RocketMQ生产者重连错误")
		log.Error(err)
		return err
	}
	_ = CloseProducer(oldProducer)
	return nil
}

func connectProducer(conf *ProducerConf) error {
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
	err = p.Start()
	if err != nil {
		log.Error(err)
		return err
	}

	globalProducer = p
	return nil
}

func CloseProducer(p rocketmq.Producer) (err error) {
	if p != nil {
		err = p.Shutdown()
		if err != nil {
			log.Error(err)
		}
	}
	return
}
