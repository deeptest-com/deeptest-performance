package logs

import (
	"github.com/aaronchen2k/deeptest/proto"
)

type GrpcSender struct {
	Stream *proto.PerformanceService_ExecServer
}

func NewGrpcSender(stream *proto.PerformanceService_ExecServer) MessageSender {
	ret := GrpcSender{
		Stream: stream,
	}

	return ret
}

func (s GrpcSender) Send(result proto.PerformanceExecResult) (err error) {
	Count(result)

	//mqData := mq.MqMsg{
	//	Event:  "result",
	//	Result: result,
	//}
	//mq.PubAgentMsg(mqData)

	(*s.Stream).Send(&result)

	return
}
