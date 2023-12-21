package mq

import (
	"fmt"
	"github.com/aaronchen2k/deeptest/pkg/core/mq"
	"github.com/aaronchen2k/deeptest/proto"
	"time"
)

var (
	mqTopicServer  = "MQ_TOPIC_SERVER"
	mqClientServer *mq.Client
)

func InitServerMq() {
	mqClientServer = mq.NewClient()
	//defer mqClientServer.Close()
	mqClientServer.SetConditions(10000000)
}

func SubServerMsg(callback func(mqMsg MqMsg) error) {
	ch, err := mqClientServer.Subscribe(mqTopicServer)
	if err != nil {
		fmt.Printf("sub mq topic %s failed, err: %s\n", mqTopicServer, err.Error())
		return
	}

	for {
		msg := mqClientServer.GetPayLoad(ch).(MqMsg)
		fmt.Printf("get mq msg [%s]%s\n", mqTopicServer, msg.Event)

		if msg.Event == "exit" {
			mqClientServer.Unsubscribe(mqTopicServer, ch)
			break
		} else {
			callback(msg)
		}

		time.Sleep(time.Millisecond * 100)
	}
}
func PubServerMsg(data MqMsg) {
	err := mqClientServer.Publish(mqTopicServer, data)
	if err != nil {
		fmt.Println("pub mq message failed", err)
	}
}

type MqMsg struct {
	Event  string                      `json:"event"`
	Result proto.PerformanceExecResult `json:"result"`
}
