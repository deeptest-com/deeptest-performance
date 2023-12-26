package mq

import (
	"github.com/aaronchen2k/deeptest/pkg/core/mq"
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
