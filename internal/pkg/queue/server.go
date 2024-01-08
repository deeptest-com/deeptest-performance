package queue

import (
	"context"
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/pkg/core/queue"
	"github.com/aaronchen2k/deeptest/proto"
	"time"
)

var (
	queueTopicOfServer  = "QUEUE_TOPIC_OF_SERVER"
	queueClientOfServer *queue.Client
)

func InitServerQueue() {
	queueClientOfServer = queue.NewClient()
	//defer queueClientOfServer.Close()
	queueClientOfServer.SetConditions(10000000)
}

func PubServerMsg(data proto.PerformanceExecResult) {
	err := queueClientOfServer.Publish(queueTopicOfServer, data)
	if err != nil {
		fmt.Println("pub mq message failed", err)
	}
}

func SubServerMsg(callback func(result proto.PerformanceExecResult) error, cancel context.CancelFunc) {
	ch, err := queueClientOfServer.Subscribe(queueTopicOfServer)
	if err != nil {
		fmt.Printf("sub mq topic %s failed, err: %s\n", queueTopicOfServer, err.Error())
		return
	}

	for {
		msg := queueClientOfServer.GetPayLoad(ch).(proto.PerformanceExecResult)
		fmt.Printf("get mq msg [%s]%s\n", queueTopicOfServer, msg.Instruction)

		if msg.Instruction == consts.Exit.String() {
			queueClientOfServer.Unsubscribe(queueTopicOfServer, ch)
			cancel()
			break
		} else {
			callback(msg)
		}

		time.Sleep(time.Millisecond * 100)
	}
}
