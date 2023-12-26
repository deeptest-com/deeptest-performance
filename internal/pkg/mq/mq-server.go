package mq

import (
	"context"
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
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

func SubServerMsg(callback func(result proto.PerformanceExecResult) error, cancel context.CancelFunc) {
	ch, err := mqClientServer.Subscribe(mqTopicServer)
	if err != nil {
		fmt.Printf("sub mq topic %s failed, err: %s\n", mqTopicServer, err.Error())
		return
	}

	for {
		msg := mqClientServer.GetPayLoad(ch).(proto.PerformanceExecResult)
		fmt.Printf("get mq msg [%s]%s\n", mqTopicServer, msg.Instruction)

		if msg.Instruction == consts.Exit.String() {
			mqClientServer.Unsubscribe(mqTopicServer, ch)
			cancel()
			break

		} else {
			callback(msg)

		}

		time.Sleep(time.Millisecond * 100)
	}
}

func PubServerMsg(data proto.PerformanceExecResult) {
	err := mqClientServer.Publish(mqTopicServer, data)
	if err != nil {
		fmt.Println("pub mq message failed", err)
	}
}
