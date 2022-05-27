package producer

import (
	"context"
	"fmt"
	"testing"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/chuan-fu/Common/zlog"
)

func TestProducer1(t *testing.T) {
	err := ConnectProducer(ProducerConf{
		NameServer: []string{"192.168.4.34:9876"},
		Retry:      3,
		GroupName:  "testProducer",
	})
	if err != nil {
		log.Error(err)
		return
	}

	data, err := GetProducer().SendSync(context.TODO(),
		primitive.NewMessage("testTopic_F", []byte(`{"a":5,"b":"bb"}`)))
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Printf("%+v", data)
}
