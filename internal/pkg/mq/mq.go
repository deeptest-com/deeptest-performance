package mq

import (
	"fmt"
	"github.com/aaronchen2k/deeptest/pkg/core/mq"
	"github.com/aaronchen2k/deeptest/proto"
	"time"
)

var (
	mqTopic  = "MQ_WebsocketTopic"
	mqClient *mq.Client
)

func InitMq() {
	mqClient = mq.NewClient()
	//defer mqClient.Close()
	mqClient.SetConditions(10000000)
}

func SubMsg(callback func(mqMsg MqMsg) error) {
	ch, err := mqClient.Subscribe(mqTopic)
	if err != nil {
		fmt.Printf("sub mq topic %s failed, err: %s\n", mqTopic, err.Error())
		return
	}

	for {
		msg := mqClient.GetPayLoad(ch).(MqMsg)
		fmt.Printf("get mq msg [%s]%s\n", mqTopic, msg.Event)

		if msg.Event == "exit" {
			mqClient.Unsubscribe(mqTopic, ch)
			break
		} else {
			callback(msg)
		}

		time.Sleep(time.Millisecond * 100)
	}
}

func PubMsg(data MqMsg) {
	err := mqClient.Publish(mqTopic, data)
	if err != nil {
		fmt.Println("pub mq message failed", err)
	}
}

type MqMsg struct {
	Event  string                      `json:"event"`
	Result proto.PerformanceExecResult `json:"result"`
}
