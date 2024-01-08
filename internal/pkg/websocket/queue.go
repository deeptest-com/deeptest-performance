package websocketHelper

import (
	"fmt"
	"github.com/aaronchen2k/deeptest/pkg/core/queue"
	"github.com/aaronchen2k/deeptest/pkg/domain"
	"time"
)

var (
	queueTopicOfWebSocket  = "QUEUE_TOPIC_OF_WEBSOCKET"
	queueClientOfWebSocket *queue.Client
)

func InitMq() {
	queueClientOfWebSocket = queue.NewClient()
	//defer queueClientOfWebSocket.Close()
	queueClientOfWebSocket.SetConditions(10000)

	go SubMsg()
}

func SubMsg() {
	ch, err := queueClientOfWebSocket.Subscribe(queueTopicOfWebSocket)
	if err != nil {
		fmt.Printf("sub mq topic %s failed\n", queueTopicOfWebSocket)
		return
	}

	for {
		msg := queueClientOfWebSocket.GetPayLoad(ch).(_domain.MqMsg)
		fmt.Printf("%s get mq msg '%#v'\n", queueTopicOfWebSocket, msg.Content)

		if msg.Content == "exit" {
			queueClientOfWebSocket.Unsubscribe(queueTopicOfWebSocket, ch)
			break
		} else {
			Broadcast(msg.Namespace, msg.Room, msg.Event, msg.Content)
		}

		time.Sleep(time.Millisecond * 100)
	}
}

func PubMsg(data _domain.MqMsg) {
	err := queueClientOfWebSocket.Publish(queueTopicOfWebSocket, data)
	if err != nil {
		fmt.Println("pub mq message failed")
	}
}
