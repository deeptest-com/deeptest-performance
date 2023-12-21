package mq

import (
	"fmt"
	"github.com/aaronchen2k/deeptest/pkg/core/mq"
	"github.com/aaronchen2k/deeptest/proto"
	"time"
)

var (
	mqTopicAgent  = "MQ_TOPIC_AGENT"
	mqClientAgent *mq.Client
)

func InitAgentMq() {
	mqClientAgent = mq.NewClient()
	//defer mqClientAgent.Close()
	mqClientAgent.SetConditions(10000000)
}

func SubAgentMsgWithStream(callback func(MqMsg, *proto.PerformanceService_ExecServer) error,
	stream *proto.PerformanceService_ExecServer) {

	ch, err := mqClientAgent.Subscribe(mqTopicAgent)
	if err != nil {
		fmt.Printf("sub mq topic %s failed, err: %s\n", mqTopicAgent, err.Error())
		return
	}

	for {
		msg := mqClientAgent.GetPayLoad(ch).(MqMsg)
		fmt.Printf("get mq msg [%s]%s\n", mqTopicAgent, msg.Event)

		if msg.Event == "exit" {
			mqClientAgent.Unsubscribe(mqTopicAgent, ch)
			break
		} else {
			callback(msg, stream)
		}

		time.Sleep(time.Millisecond * 100)
	}
}

func PubAgentMsg(data MqMsg) {
	err := mqClientAgent.Publish(mqTopicAgent, data)
	if err != nil {
		fmt.Println("pub mq message failed", err)
	}
}
